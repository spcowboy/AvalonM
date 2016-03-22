package models

type GameStatus struct {
	GameRound   int               //1,2,3,4,5
	GameMode    int               //5,6,7,8,9,10
	PlayerRoles map[string]string //uname:role
	ReadyStatus map[string]bool   //uname :true|false

}
type GameRoom struct {
	ActiveRoom
	Status  *GameStatus
	Leader  string //room leader
	Process chan int
}

func (this *GameRoom) run(statuschan chan GameStatus) {
	this.Process <- C_PreparePhase
	for {
		switch <-this.Process {
		case C_PreparePhase:
			PrepareProcess(this.Process, statuschan)
		case C_AllGetReady:
			StartProcess(this.Process, statuschan)
		case C_GameStart:
			AssignRoleProcess(this.Process, statuschan)
		case C_AssignRoleReady:
			LightOffAndShowRoleProcess(this.Process, statuschan)
		case C_LightOffAndShowRoleReady: //开始round1,把round1的标志码作为参数传给round处理函数
			RoundStartProcess()
		case C_Round1End:
			RoundStartProcess()
		case C_Round2End:
			RoundStartProcess()
		case C_Round3End:
			RoundStartProcess()
		case C_Round4End:
			RoundStartProcess()
		case C_Round5End: //round5结束进入刺杀阶段
			AssassinateProcess()
		case C_AssassinateEnd: //刺杀结束进入游戏结果公布阶段
			ShowGameResultProcess()

		}
	}
}

//广播ready状态，判断是否全部ready，全部ready时返回标识码C_AllGetReady->process
func PrepareProcess() {

}

//告诉leader可以开始，等待leader的开始信号，得到信号后返回C_GameStart->process
func StartProcess() {

}

//分配角色，把角色信息发给对应的玩家，返回C_AssignRoleReady
func AssignRoleProcess() {

}

//天黑阶段，向对应角色玩家展示对应的玩家角色，返回C_LightOffAndShowRoleReady
func LightOffAndShowRoleProcess() {

}

//天亮，开始Round-X,返回对应round的结束标志C_Round1End|C_Round2End...
func RoundStartProcess(round int) {
	phasetag := make(chan int)
	phasetag <- C_SelectMissioner
	for {
		switch <-phasetag {
		case C_SelectMissioner: //开始选人阶段
			SelectMissionerProcess()
		case C_SelectMissionerEnd: //选人结束，开始发言阶段
			SpeechProcess()
		case C_SpeechEnd: //发言结束，开始投票
			VoteProcess()
		case C_VoteSucess: //投票成功，进入任务执行阶段
			MissionStartProcess()
		case C_VoteFail: //投票失败,判断延迟次数，决定进入重新选人（更换leader）还是执行任务
			SpeechProcess() //待完善
		case C_MissionEnd: //任务结束，本轮结束,返回本轮结束标志
			return C_Round1End|C_Round2End...//伪代码
		}
	}
}

//刺杀阶段处理，通知刺客选人，等待杀手选人，广播杀人结果，返回C_AssassinateEnd
func AssassinateProcess() {

}

//公布结果阶段，返回C_PreparePhase
func ShowGameResultProcess() {

}

//开始round，告诉leader开始选人，等待选人结果，广播结果，返回C_SelectMissionerEnd
func SelectMissionerProcess() {

}

//开始发言阶段，提醒对应玩家发言，等待玩家发言结束标志，提醒下一个玩家发言，循环，判断最后一个玩家发言完毕，返回C_SpeechEnd
func SpeechProcess() {

}

//进入投票阶段,通知玩家开始投票，接收所有投票并广播，判断进入任务执行阶段或者延迟任务，返回C_VoteSucess|C_VoteFail
func VoteProcess() {

}

//执行任务，发送任务开始通知给对应玩家，接收任务结果，保存结果，广播结果，返回C_MissionEnd
func MissionStartProcess() {

}
