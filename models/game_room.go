package models

import (
	"code.google.com/p/go.net/websocket"
	"github.com/astaxie/beego"
	//"strings"
	"container/list"
	"time"
	"strconv"
)

var RunningGameRoom *GameRoom = &GameRoom{}

func Init() {
	beego.Info("game inited.")
	InitGameRoom()
	go RunningGameRoom.rungameprocess()
	go RunningGameRoom.runroomhub()
	go RunningGameRoom.runwebsocketserver()
	go RunningGameRoom.rungameserver()
}



//每轮的信息
type RoundStatus struct {
	Leader           string   //本轮当前的leader
	Gamers           []string //本轮当前的玩家发言顺序
	VoteResult       bool     //投票结果
	MissionResult    bool     //任务结果
	UserSelected     []string //leader选择的人选集合
	CurrentStatus    int      //本轮游戏当前状态
	DelayTimes       int
	MissionResultMap map[string]bool
	VoteResultMap    map[string]bool
}
func InitGameRoom() {
	RunningGameRoom = &GameRoom{
		Gamers:      list.New(),
		GamerStatus: make(map[string]*OnlineGamer),
		UserConn:    make(map[string]*websocket.Conn),
		ReadyStatus:  make(map[string]bool),
		GamerRoles:   make(map[string]int),    //name:role
		RoomStatus:   make(map[int]*RoundStatus), //房间的轮次信息
		Leader: "",
		subscribe:     make(chan *Subscriber,256),
		Broadcast:     make(chan Message, 256),
		ProcessTag:    make(chan int, 256),
		MsgFromClient: make(chan Message, 256),
		ProcessQueue:  make(chan Message, 256),
		SingleMessage: make(chan Message, 256),
		phasetag:      make(chan int, 256),
		
	}
	//go runningGameRoom.run()
}
type GameRoom struct {
	Gamers      *list.List //[]name
	GamerStatus map[string]*OnlineGamer
	UserConn    map[string]*websocket.Conn
	ReadyStatus map[string]bool
	GamerRoles  map[string]int       //name:role
	RoomStatus  map[int]*RoundStatus //房间的轮次信息

	subscribe         chan *Subscriber
	Broadcast         chan Message
	SingleMessage     chan Message
	CloseSign         chan bool
	Status            int
	Leader            string //room leader
	ProcessTag        chan int
	phasetag          chan int //每一回合的tag
	MsgFromClient     chan Message
	ProcessQueue      chan Message //异步
	GameMode          int          //游戏模式,人数
	Speechqueue       []string     //发言顺序
	Round             int          //轮次
	AssassinateTarget string       //刺杀目标

}
type OnlineGamer struct {
	InRoom     *GameRoom
	Connection *websocket.Conn
	UserInfo   Gamer
	Send       chan Message
}
type Gamer struct {
	Name   string
	Role   int
	Leader string
}
type Message struct {
	MType       string //chat|notify|controll
	TextMessage TextMessage
	UserStatus  UserStatus
}
type TextMessage struct {
	Content  string
	UserInfo Gamer
	Time     string
}
type UserStatus struct {
	Me Gamer

	Users              []string
	Role               int  //玩家角色
	CurrentStatus      int  //玩家当前状态,指示客户端执行相应操作
	VoteResult         bool //投票结果
	VoteResultMap      map[string]bool
	ReadyStatus        map[string]bool
	MissionResult      bool           //任务结果
	MissionResultArray []bool         //任务结果列表
	UserSelected       []string       //leader选择的人选集合
	AssassinateResult  string         //杀人目标
	ShowRole           map[string]int //向玩家显示其他对应玩家的角色,name:role(int,由客户端和服务端查询角色对应表)
	CountDown          int            //让客户端进行倒计时
	Round              int
}
type Subscriber struct {
	Name string
	Conn *websocket.Conn
}

func (this *GameRoom) GetOnlineUsers() (users []*Gamer) {
	for _, online := range this.GamerStatus {
		users = append(users, &online.UserInfo)
	}
	return
}

func (this *GameRoom) Join(user string, ws *websocket.Conn) {
	beego.Info("send user and ws")
	

	beego.Info("new gamer "+user)

	if !isUserExist(this.Gamers, user) {
		onlineGamer := &OnlineGamer{
			InRoom:     this,
			Connection: ws,
			Send:       make(chan Message, 256),
			UserInfo:   Gamer{Name: user}}
			
		this.AddUser(onlineGamer)
		
		//publish <- newEvent(models.EVENT_JOIN, a.Name, "")
		beego.Info("New user:", user, ";WebSocket:", ws != nil)
	} else {
		beego.Info("Old user:", user, ";WebSocket:", ws != nil)
	}

	//this.subscribe <- &Subscriber{Name: user, Conn: ws}
}

func isUserExist(subscribers *list.List, user string) bool {
	for sub := subscribers.Front(); sub != nil; sub = sub.Next() {
		if sub.Value.(string) == user {
			return true
		}
	}
	return false
}

func (this *GameRoom) AddUser(onlineGamer *OnlineGamer) {
	this.Gamers.PushBack(onlineGamer.UserInfo.Name) // Add user to the end of list.
	this.ReadyStatus[onlineGamer.UserInfo.Name] = false
	this.UserConn[onlineGamer.UserInfo.Name] = onlineGamer.Connection
	if this.Leader == "" {
		this.Leader = onlineGamer.UserInfo.Name
		onlineGamer.UserInfo.Leader = this.Leader
	}
	this.GamerStatus[onlineGamer.UserInfo.Name] = onlineGamer

	m := Message{
		MType: UPDATE_MTYPE, //更新用户加入信息
		TextMessage: TextMessage{
			Content: onlineGamer.UserInfo.Name+" join this room.",
			Time:humanCreatedAt(),
			},
		UserStatus: UserStatus{
			CurrentStatus: this.Status, //告诉用户当前处于状态
			ReadyStatus:   this.ReadyStatus,
			Me: Gamer{
				Leader: this.Leader,
				Name: onlineGamer.UserInfo.Name,
			},
			Users: toArray(this.Gamers),
		},
	}
	
	
	go onlineGamer.PushToClient()
	this.Broadcast <- m
	onlineGamer.PullFromClient()
	
}
//接收客户端信息
func (this *OnlineGamer) PullFromClient() {
	for {
		var msg Message
		err := websocket.JSON.Receive(this.Connection, &msg)
		// If user closes or refreshes the browser, a err will occur
		if err != nil {
			beego.Info("PullFromClient "+this.UserInfo.Name+" with websocket msg error :"+err.Error())
			break
		}
		//beego.Info(msg.MType)
		msg.TextMessage.UserInfo = this.UserInfo
		msg.TextMessage.Time = humanCreatedAt()

		this.InRoom.MsgFromClient <- msg
	}
	defer this.killUserResource()
}

//发送信息到客户端
func (this *OnlineGamer) PushToClient() {
	for b := range this.Send {
		//beego.Info("send websocket msg to "+this.UserInfo.Name)
		err := websocket.JSON.Send(this.Connection, b) //注意这里采用根据socket发送失败的判断来退出goroutine,如果直接关闭chan并不能使本goroutine退出
		if err != nil {
			beego.Info("send websocket msg to "+this.UserInfo.Name+" error :"+err.Error())
			break
		}
	}
}
func (this *OnlineGamer) killUserResource() {
	this.InRoom.DelUser(this)
	this.Connection.Close()

	close(this.Send)

	m := Message{
		MType: UPDATE_MTYPE,
		TextMessage: TextMessage{
			Content: this.UserInfo.Name+" leave this room.",
			Time:humanCreatedAt(),
			},
		UserStatus: UserStatus{
		
			Me: Gamer{
				Leader: this.InRoom.Leader,
				Name: this.UserInfo.Name,
			},
			Users: toArray(this.InRoom.Gamers),
		},
	}
	
	RunningGameRoom.Broadcast <- m
}
func (this *GameRoom) DelUser(onlineGamer *OnlineGamer) {

	if isUserExist(this.Gamers, onlineGamer.UserInfo.Name) {
		if contain, e := Contains(this.Gamers, onlineGamer.UserInfo.Name); contain {
			this.Gamers.Remove(e)
		}

		delete(this.ReadyStatus, onlineGamer.UserInfo.Name)
		if onlineGamer.UserInfo.Leader == onlineGamer.UserInfo.Name {
			this.Leader = onlineGamer.UserInfo.Leader
		}
		delete(this.GamerStatus, onlineGamer.UserInfo.Name)
		delete(this.UserConn, onlineGamer.UserInfo.Name)

		ws := onlineGamer.Connection
		if ws != nil {
			ws.Close()
			beego.Error("WebSocket closed:", onlineGamer.UserInfo.Name)
		}
	}
	if this.Leader==onlineGamer.UserInfo.Name{
		if(this.Gamers.Len()!=0){
				this.Leader=this.Gamers.Front().Value.(string)
			}else{
				this.Leader=""
			}
		
	}
}

//游戏逻辑服务器
func (this *GameRoom) rungameserver() {
	//this.ProcessTag <- C_PreparePhase
	for {
		this.Status=<-this.ProcessTag
		switch  this.Status{

		case C_PreparePhase:
			go this.PrepareProcess()
		case C_AllGetReady:
			this.StartProcess()
		case C_GameStart:
			this.Speechqueue = make([]string,this.Gamers.Len())
			//beego.Info(len(this.Speechqueue))
			//beego.Info(this.GameMode)
			this.AssignRoleProcess()
		case C_AssignRoleReady: //角色分配完成，进入天黑阶段
			this.LightOffAndShowRoleProcess()
		case C_LightOffAndShowRoleReady: //开始round1,把round1的标志码作为参数传给round处理函数
			go this.RoundStartProcess(1)
		case C_RoundEnd:
			if this.Round > 2 { //判断是否可以直接进入刺杀阶段
				suc := 0
				fail := 0
				for i := 1; i != this.Round+1; i++ {
					if this.RoomStatus[i].MissionResult {
						suc++
					} else {
						fail++
					}
				}
				if suc > 2 {
					this.AssassinateProcess() //好人胜利进入刺杀阶段
					break
				} else if fail > 2 {
					this.ShowGameResult() //坏人胜利提前结束
					break
				}

			}
			if this.Round != 5 {
				this.Round++
				go this.RoundStartProcess(this.Round)
			} else { //round5结束进入刺杀阶段
				this.AssassinateProcess()
			}

		case C_AssassinateEnd: //刺杀结束进入游戏结果公布阶段
			ShowGameResultProcess(this.ProcessTag, this.ProcessQueue)

		}
	}
}

//消息路由
func (this *GameRoom) runroomhub() {
	for {
		select {
		case a := <-this.MsgFromClient:
			switch a.MType {
			case CHAT_MTYPE:
				this.Broadcast <- a
			case CONTROL_MTYPE:
				this.ProcessQueue <- a

			}

		}
	}
}

//游戏控制消息处理服务器
func (this *GameRoom) rungameprocess() {
	for {
		select {
		case a := <-this.ProcessQueue:
			if a.MType == CONTROL_MTYPE {
				if a.UserStatus.CurrentStatus == C_Ready {
					//beego.Info("user ready: "+a.UserStatus.Me.Name)
					this.ReadyStatus[a.UserStatus.Me.Name] = true //更新准备状态
					// beego.Info("user  num is : "+strconv.Itoa(this.Gamers.Len()))
					// if this.AllReady(){
					// 	beego.Info("all ready : true")
					// }
					
					if this.Gamers.Len()>4&&this.AllReady() {
						beego.Info("user > 4 and all ready.")
						this.ProcessTag <- C_AllGetReady
					}
					this.Broadcast <- a
				} else if a.UserStatus.CurrentStatus == C_UnReady {
					this.ReadyStatus[a.UserStatus.Me.Name] = false
					this.Broadcast <- a
				} else if a.UserStatus.CurrentStatus == C_LeaderStart { //leader开始游戏
					this.Broadcast <- a
					this.ProcessTag <- C_GameStart
				} else if a.UserStatus.CurrentStatus == C_ReStart { //重新开始游戏

					this.Broadcast <- a
					this.ProcessTag <- C_GameStart
				} else if a.UserStatus.CurrentStatus == C_LeaderFinSelect { //leader选人完毕
					//广播选人结果
					beego.Info(a.UserStatus.Me.Name+" select over :")
					this.RoomStatus[this.Round].UserSelected = a.UserStatus.UserSelected
					this.Broadcast <- a
					this.phasetag <- C_SelectMissionerEnd
				} else if a.UserStatus.CurrentStatus == C_VoteFin {
					this.RoomStatus[this.Round].VoteResultMap[a.UserStatus.Me.Name] = a.UserStatus.VoteResult
					if a.UserStatus.VoteResult{
						beego.Info(a.UserStatus.Me.Name+" is aprove")
					}else{
						beego.Info(a.UserStatus.Me.Name+" is oppose")
					}
					
					//判断是否全部投票完毕
					if this.AllVote() {
						beego.Info("allvote fin.")
						//广播投票结果
						if this.VoteSucc() {
							this.RoomStatus[this.Round].VoteResult = true
							a.UserStatus.CurrentStatus = C_VoteSuccAndStartMission
							a.UserStatus.VoteResult = true
							a.UserStatus.VoteResultMap = this.RoomStatus[this.Round].VoteResultMap
							this.LeaderTrans()
							a.UserStatus.Me.Leader = this.Leader

							this.Broadcast <- a
							this.phasetag <- C_VoteSucess
						} else {
							this.RoomStatus[this.Round].VoteResult = false
							a.UserStatus.CurrentStatus = C_VoteFailAndReSelect
							a.UserStatus.VoteResult = false
							this.LeaderTrans()
							a.UserStatus.Me.Leader = this.Leader

							this.Broadcast <- a
							this.phasetag <- C_VoteFail
						}
					}

				} else if a.UserStatus.CurrentStatus ==  C_MissionEnd {
					this.RoomStatus[this.Round].MissionResultMap[a.UserStatus.Me.Name] = a.UserStatus.MissionResult
					suc, allfin := this.AllMissionFin()
					if allfin {
						//任务结果
						m := Message{
							MType:CONTROL_MTYPE,
							UserStatus: UserStatus{
								CurrentStatus:C_MissionResult,
								MissionResultArray:make([]bool,getMissionerNumber(this.GameMode,this.Round)),
							},
						}
						
						

						if suc {
							m.UserStatus.MissionResult = true
						} else {
							m.UserStatus.MissionResult = false
						}
						i := 0
						for _, r := range this.RoomStatus[this.Round].MissionResultMap {
							m.UserStatus.MissionResultArray[i] = r
							i++
						}
						m.UserStatus.CountDown = 10
						this.Broadcast <- m
						time.Sleep(time.Second * 11)
						this.phasetag <- C_MissionEnd
					}

				} else if a.UserStatus.CurrentStatus == C_AssassinateResult { //收到刺杀信息,广播刺杀对象，公布游戏结果
					this.AssassinateTarget = a.UserStatus.AssassinateResult
					GameResult := this.GetGameResult()
					a.TextMessage.Content = GameResult
					a.TextMessage.Time = humanCreatedAt()
					a.UserStatus.CurrentStatus = C_ShowGameResult
					a.UserStatus.ShowRole = this.GamerRoles
					this.Broadcast <- a
					this.ProcessTag <- C_PreparePhase
				}

			}

		}
	}
}

// Core function of room
func (this *GameRoom) runwebsocketserver() {
	//beego.Info("runwebsocketserver init.")
	for {
		select {
		// case a := <-this.subscribe:
		// 	beego.Info("new gamer "+a.Name)

		// 	if !isUserExist(this.Gamers, a.Name) {
		// 		onlineGamer := &OnlineGamer{
		// 			InRoom:     this,
		// 			Connection: a.Conn,
		// 			Send:       make(chan Message, 256),
		// 			UserInfo:   &Gamer{Name: a.Name}}

		// 		this.AddUser(onlineGamer)

		// 		//publish <- newEvent(models.EVENT_JOIN, a.Name, "")
		// 		beego.Info("New user:", a.Name, ";WebSocket:", a.Conn != nil)
		// 	} else {
		// 		beego.Info("Old user:", a.Name, ";WebSocket:", a.Conn != nil)
		// 	}
		case b := <-this.Broadcast:
			for _, online := range this.GamerStatus {
				online.Send <- b
			}
		case d := <-this.SingleMessage:
			this.GamerStatus[d.UserStatus.Me.Name].Send <- d
		case c := <-this.CloseSign:
			if c == true {
				close(this.Broadcast)
				close(this.CloseSign)
				close(this.ProcessTag)
				close(this.MsgFromClient)
				close(this.ProcessQueue)
				return
			}
		}
	}
}



//===================================================================游戏方法====================================================
//所有人准备完毕
func (this *GameRoom) AllReady() bool {
	readyusers:=0
	for _, r := range this.ReadyStatus {
		if r == true {
			readyusers++
		}
	}
	if readyusers>=this.Gamers.Len()-1{
		return true
	}
	return false
}

//所有人投票完毕
func (this *GameRoom) AllVote() bool {
	if len(this.RoomStatus[this.Round].VoteResultMap) == this.GameMode {
		return true
	}
	return false
}

//判断投票是否成功
func (this *GameRoom) VoteSucc() bool {
	suc := 0
	for _, r := range this.RoomStatus[this.Round].VoteResultMap {
		if r == true {
			suc++
		}
	}

	if suc > this.GameMode/2 {
		
		return true
	}
	return false
}

//判断所有人任务结束
func (this *GameRoom) AllMissionFin() (suc, allfin bool) {
	if len(this.RoomStatus[this.Round].MissionResultMap) == len(this.RoomStatus[this.Round].UserSelected) {
		allfin = true

	} else {
		allfin = false
		return
	}
	sucm := 0
	for _, r := range this.RoomStatus[this.Round].MissionResultMap {
		if r == true {
			sucm++
		}
	}
	if sucm < len(this.RoomStatus[this.Round].MissionResultMap) {
		suc = false
		this.RoomStatus[this.Round].MissionResult = false
	} else {
		suc = true
		this.RoomStatus[this.Round].MissionResult = true
	}
	return
}

//leader传递
func (this *GameRoom) LeaderTrans() {
	tempqueue:=this.Speechqueue[1:]
	
	tempqueue = append(tempqueue,this.Speechqueue[0])
	this.Speechqueue = tempqueue
	this.Leader = this.Speechqueue[this.GameMode-1]

}

//整理并返回结果
func (this *GameRoom) GetGameResult() string {
	for name, r := range this.GamerRoles {
		if r == R_MERLIN {
			if this.AssassinateTarget == name {
				return "Merlin been killed"
			} else {
				suc := 0
				for i := 1; i != this.Round+1; i++ {
					if this.RoomStatus[i].MissionResult {
						suc++
					}
				}
				if suc > 2 {
					return "Good win"
				} else {
					return "Bad win"
				}
			}
		}
	}
	return ""
}

//提前发表游戏结果,坏人胜利
func (this *GameRoom) ShowGameResult() {
	a := Message{}
	a.TextMessage.Content = "Bad win"
	a.TextMessage.Time = humanCreatedAt()
	a.UserStatus.CurrentStatus = C_ShowGameResult
	a.UserStatus.ShowRole = this.GamerRoles
	this.Broadcast <- a
	this.ProcessTag <- C_PreparePhase
}

//分配角色
func (this *GameRoom) AssignRole(send bool) {
	var rolelist []int
	switch this.GameMode {
	case 5:
		rolelist = random([]int{R_MERLIN, R_Percival, R_loyalist, R_Morgana, R_Assassin})
	case 6: //梅林、派西维尔、忠臣*2  vs 莫甘娜、刺客
		rolelist = random([]int{R_MERLIN, R_Percival, R_loyalist, R_loyalist, R_Morgana, R_Assassin})
	case 7: //梅林、派西维尔、忠臣*2  vs 莫甘娜、奥伯伦、刺客
		rolelist = random([]int{R_MERLIN, R_Percival, R_loyalist, R_loyalist, R_Morgana, R_Oberon, R_Assassin})
	case 8: //梅林、派西维尔、忠臣*3  vs 莫甘娜、刺客、爪牙
		rolelist = random([]int{R_MERLIN, R_Percival, R_loyalist, R_loyalist, R_loyalist, R_Morgana, R_Assassin, R_Minion})
	case 9: //梅林、派西维尔、忠臣*4  vs 莫德雷德、莫甘娜、刺客
		rolelist = random([]int{R_MERLIN, R_Percival, R_loyalist, R_loyalist, R_loyalist, R_loyalist, R_Mordred, R_Morgana, R_Assassin})
	case 10: //梅林、派西维尔、忠臣*4  vs 莫德雷德、莫甘娜、奥伯伦、刺客
		rolelist = random([]int{R_MERLIN, R_Percival, R_loyalist, R_loyalist, R_loyalist, R_loyalist, R_Mordred, R_Morgana, R_Oberon, R_Assassin})
	}
	i := 0

	for _, online := range this.GamerStatus {
		online.UserInfo.Role = rolelist[i]
		this.GamerRoles[online.UserInfo.Name] = online.UserInfo.Role
		beego.Info(online.UserInfo.Name+" is "+strconv.Itoa(online.UserInfo.Role))
		if send == true {
			m := Message{
				MType: CONTROL_MTYPE,
				UserStatus: UserStatus{
					CurrentStatus: C_AssignRole,
					Me: Gamer{
						Name: online.UserInfo.Name,
						Role: online.UserInfo.Role,
					},

					Role: online.UserInfo.Role,
				},
			}
			online.Send<-m
			//this.SingleMessage <- m
		}
		i++
	}
}

//
//广播PrepareProcess,用户进入准备状态,客户端ready按键显示
func (this *GameRoom) PrepareProcess() {
	m := Message{
		MType: CONTROL_MTYPE, 

		UserStatus: UserStatus{
			CurrentStatus: C_PreparePhase, 
			Me: Gamer{
				Leader: this.Leader,
			},
		},
	}
	beego.Info("send C_PreparePhase to every one.")
	this.Broadcast <- m
}

//告诉leader可以开始，等待leader的开始信号，得到信号后返回C_GameStart->process
func (this *GameRoom) StartProcess() {
	m := Message{
		MType: CONTROL_MTYPE, //通知leader可以开始了

		UserStatus: UserStatus{
			CurrentStatus: C_LeaderCanStart, //通知leader可以开始了
			Me: Gamer{
				Name: this.Leader,
			},
		},
	}
	beego.Info("send C_LeaderCanStart to every one.")
	this.Broadcast <- m
}

//游戏开始，根据房间人数决定游戏模式（人数超过10人的处理以后再做）=>分配角色，把角色信息发给对应的玩家，返回C_AssignRoleReady
func (this *GameRoom) AssignRoleProcess() {
	number := this.Gamers.Len()
	this.GameMode = number
	i := 0
	for g := this.Gamers.Front(); g != nil; g=g.Next() {
		//beego.Info(len(this.Speechqueue))
		//beego.Info(this.GameMode)
		this.Speechqueue[i] = g.Value.(string)
		i++
	}
	tempqueue:=this.Speechqueue[1:]
	
	tempqueue = append(tempqueue,this.Speechqueue[0])
	this.Speechqueue = tempqueue
	this.AssignRole(true)
	
	time.Sleep(time.Second*5)
	this.ProcessTag <- C_AssignRoleReady
}

//天黑阶段，向对应角色玩家展示对应的玩家角色，返回C_LightOffAndShowRoleReady:梅林=莫德雷德之外的坏人;除了奥伯伦以外的坏人睁眼，互相辨认同伙;派西维尔=梅林和莫甘娜
func (this *GameRoom) LightOffAndShowRoleProcess() {
	//先发控制消息告诉客户端进入天黑阶段
	m := Message{
		MType: CONTROL_MTYPE,

		UserStatus: UserStatus{
			CurrentStatus: C_GoIntoLightOffPhase,
		},
	}
	this.Broadcast <- m
	time.Sleep(time.Second*5)
	//构建控制消息给对应的角色玩家
	//梅林
	Roles := make(map[string]int)
	mtomerlin := Message{}
	for name, role := range this.GamerRoles {
		if role == R_MERLIN {
			//beego.Info(name)
			mtomerlin.MType = CONTROL_MTYPE
			mtomerlin.UserStatus.Me.Name= name
			mtomerlin.UserStatus.CurrentStatus = C_RolesToMerlin

		} else if role == R_Morgana || role == R_Oberon || role == R_Assassin || role == R_Minion {
			Roles[name] = role
		}
	}
	mtomerlin.UserStatus.ShowRole = Roles
	this.SingleMessage <- mtomerlin
	//badguy
	Rolesbad := make(map[string]int)

	mtobadguy := Message{}
	for name, role := range this.GamerRoles {
		if role == R_Morgana || role == R_Mordred || role == R_Assassin || role == R_Minion {
			Rolesbad[name] = role
		}
	}
	mtobadguy.MType = CONTROL_MTYPE
	mtobadguy.UserStatus.CurrentStatus = C_RolesBadToBad
	mtobadguy.UserStatus.ShowRole = Rolesbad
	for name, role := range Rolesbad {
		mtobadguy.UserStatus.Me.Name = name
		mtobadguy.UserStatus.Me.Role = role
		this.SingleMessage <- mtobadguy //可能会有问题,没有重新建message对象
	}

	//派西维尔

	mtoPercival := Message{}
	RolestoPercival := make(map[string]int)
	for name, role := range this.GamerRoles {
		if role == R_Percival {
			mtoPercival.MType = CONTROL_MTYPE
			mtoPercival.UserStatus.CurrentStatus = C_RolesToPercival
			mtoPercival.UserStatus.Me.Name = name
		} else if role == R_MERLIN || role == R_Morgana {
			RolestoPercival[name] = R_Unknown
		}
	}
	mtoPercival.UserStatus.ShowRole = RolestoPercival
	this.SingleMessage <- mtoPercival

	//广播等待20秒，客户端进行倒数，服务器sleep21秒
	m2 := Message{}
	m2.MType = CONTROL_MTYPE
	m2.UserStatus.CountDown = 20
	m2.UserStatus.CurrentStatus = C_CountDown
	this.Broadcast <- m2

	time.Sleep(time.Second * 21)
	this.ProcessTag <- C_LightOffAndShowRoleReady
}

//天亮，开始Round-X,返回对应round的结束标志C_Round1End|C_Round2End...
func (this *GameRoom) RoundStartProcess(round int) {
	mlighton := Message{}
	mlighton.MType = CONTROL_MTYPE

	// Leader           string   //本轮当前的leader
	// Gamers           []string //本轮当前的玩家发言顺序
	// VoteResult       bool     //投票结果
	// MissionResult    bool     //任务结果
	// UserSelected     []string //leader选择的人选集合
	// CurrentStatus    int      //本轮游戏当前状态
	// DelayTimes       int
	// MissionResultMap map[string]bool
	// VoteResultMap    map[string]bool
	switch round {
	case 1:
		this.Round = 1
		this.RoomStatus[1] = &RoundStatus{
			Leader: this.Leader,      //本轮当前的leader
			Gamers: this.Speechqueue, //本轮当前的玩家发言顺序
			DelayTimes: 0,
			MissionResultMap: make(map[string]bool),
			VoteResultMap: make(map[string]bool),
			CurrentStatus: C_Round1Start, //本轮游戏当前状态
		}
		mlighton.UserStatus.CurrentStatus = C_LightOnandR1start
	case 2:
		this.RoomStatus[2] = &RoundStatus{
			Leader: this.Leader,      //本轮当前的leader
			Gamers: this.Speechqueue, //本轮当前的玩家发言顺序
			DelayTimes: 0,
			MissionResultMap: make(map[string]bool),
			VoteResultMap: make(map[string]bool),
			CurrentStatus: C_Round2Start, //本轮游戏当前状态
		}
		mlighton.UserStatus.CurrentStatus = C_R2start
	case 3:
		this.RoomStatus[3] = &RoundStatus{
			Leader: this.Leader,      //本轮当前的leader
			Gamers: this.Speechqueue, //本轮当前的玩家发言顺序
			DelayTimes: 0,
			MissionResultMap: make(map[string]bool),
			VoteResultMap: make(map[string]bool),
			CurrentStatus: C_Round3Start, //本轮游戏当前状态
		}
		mlighton.UserStatus.CurrentStatus = C_R3start
	case 4:
		this.RoomStatus[4] = &RoundStatus{
			Leader: this.Leader,      //本轮当前的leader
			Gamers: this.Speechqueue, //本轮当前的玩家发言顺序
			DelayTimes: 0,
			MissionResultMap: make(map[string]bool),
			VoteResultMap: make(map[string]bool),
			CurrentStatus: C_Round4Start, //本轮游戏当前状态
		}
		mlighton.UserStatus.CurrentStatus = C_R4Start
	case 5:
		this.RoomStatus[5] = &RoundStatus{
			Leader: this.Leader,      //本轮当前的leader
			Gamers: this.Speechqueue, //本轮当前的玩家发言顺序
			DelayTimes: 0,
			MissionResultMap: make(map[string]bool),
			VoteResultMap: make(map[string]bool),
			CurrentStatus: C_Round5Start, //本轮游戏当前状态
		}
		mlighton.UserStatus.CurrentStatus = C_R5start
	}

	this.Broadcast <- mlighton

	this.phasetag <- C_SelectMissioner
	for {
		switch <-this.phasetag {
		case C_SelectMissioner: //开始选人阶段,
			this.SelectMissionerProcess(false)
		case C_SelectMissionerEnd: //选人结束，开始发言阶段
			if this.RoomStatus[this.Round].DelayTimes > 4 {
				m := Message{}
				m.MType = CONTROL_MTYPE
				m.UserStatus.CurrentStatus = C_DelayOver4TimesAndNoSpeech
				this.Broadcast <- m
				this.MissionStartProcess()
			} else {
				this.SpeechProcess()
			} //延迟不能大于4次

		case C_SpeechEnd: //发言结束，开始投票
			this.VoteProcess()
		case C_VoteSucess: //投票成功，进入任务执行阶段
			this.MissionStartProcess()
		case C_VoteFail: //投票失败,判断延迟次数，决定进入重新选人（更换leader）还是执行任务
			this.SelectMissionerProcess(true) //待完善
		case C_MissionEnd: //任务结束，本轮结束,返回本轮结束标志
			this.ProcessTag <- C_RoundEnd
			return
		}
	}
}

//刺杀阶段处理，通知刺客选人，等待杀手选人，广播杀人结果，返回C_AssassinateEnd
func (this *GameRoom) AssassinateProcess() {
	m := Message{}
	m.MType = CONTROL_MTYPE
	m.UserStatus.CurrentStatus = C_Assassinate
	for name, r := range this.GamerRoles {
		if r == R_Assassin {
			m.UserStatus.Me.Name = name
			break
		}
	}
	m.UserStatus.CountDown = 15
	this.Broadcast <- m

}

//公布结果阶段，返回C_PreparePhase
func ShowGameResultProcess(processTag chan int, processQueue chan Message) {

}

//开始round，告诉leader开始选人，等待选人结果，广播结果，返回C_SelectMissionerEnd
func (this *GameRoom) SelectMissionerProcess(reselect bool) {
	if reselect == true {
		this.RoomStatus[this.Round].DelayTimes++
	}
	mtoleader := Message{}
	mtoleader.MType = CONTROL_MTYPE
	mtoleader.UserStatus.CurrentStatus = C_LeaderstartSelect
	mtoleader.UserStatus.Me.Name = this.Leader
	this.Broadcast <- mtoleader

}

//开始发言阶段，提醒对应玩家发言，等待玩家发言结束标志，提醒下一个玩家发言，循环，判断最后一个玩家发言完毕，返回C_SpeechEnd
func (this *GameRoom) SpeechProcess() {
	
	//this.RoomStatus[this.Round].Gamers[this.GameMode]=this.RoomStatus[this.Round].Gamers[0]
	for _, name := range this.RoomStatus[this.Round].Gamers {
		if name!=""{
			m2 := Message{} //挨个发通知发言，然后等待30秒，换下一个
			m2.MType = CONTROL_MTYPE
			m2.UserStatus.Me.Name = name
			m2.UserStatus.CurrentStatus = C_StarttoSpeech
			m2.UserStatus.CountDown = 30 //客户端根据countdown自动判断倒数
			this.Broadcast <- m2
			time.Sleep(time.Second * 31)
		}
		
	}
	//全部发言完毕
	this.phasetag <- C_SpeechEnd
}

//进入投票阶段,通知玩家开始投票，接收所有投票并广播，判断进入任务执行阶段或者延迟任务，返回C_VoteSucess|C_VoteFail
func (this *GameRoom) VoteProcess() {
	m := Message{}
	m.MType = CONTROL_MTYPE
	m.UserStatus.CurrentStatus = C_StartToVote
	m.UserStatus.CountDown = 10
	this.Broadcast <- m

}

//执行任务，发送任务开始通知给对应玩家，接收任务结果，保存结果，广播结果，返回C_MissionEnd
func (this *GameRoom) MissionStartProcess() {
	m := Message{}
	m.MType = CONTROL_MTYPE
	m.UserStatus.CurrentStatus = C_StartToMission
	m.UserStatus.UserSelected = this.RoomStatus[this.Round].UserSelected
	this.Broadcast <- m
}
