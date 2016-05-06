//=======================================================================
//define
var socket;
localStorage.Me=$('#uname').text()
localStorage.Leader=""
localStorage.Currentstatus = 0
localStorage.Round = 0
localStorage.Role = -1
localStorage.numberofselect=0
var Users = [localStorage.Me]
var ctx1=$("#user1").getContext("2d");
var ctx2=$("#user2").getContext("2d");
var ctx3=$("#user3").getContext("2d");
var ctx4=$("#user4").getContext("2d");
var ctx5=$("#user5").getContext("2d");
var ctx6=$("#user6").getContext("2d");
var ctx7=$("#user7").getContext("2d");
var ctx8=$("#user8").getContext("2d");
var ctx9=$("#user9").getContext("2d");
var ctx10=$("#user10").getContext("2d");
var img=document.getElementById("icon-user");

img.onload = function(){
    ctx1.drawImage(img,0,0,50,50);
    
    ctx2.drawImage(img,0,0,50,50);
    ctx3.drawImage(img,0,0,50,50);
    ctx4.drawImage(img,0,0,50,50);
    ctx5.drawImage(img,0,0,50,50);
    ctx6.drawImage(img,0,0,50,50);
    ctx7.drawImage(img,0,0,50,50);
    ctx8.drawImage(img,0,0,50,50);
    ctx9.drawImage(img,0,0,50,50);
    ctx10.drawImage(img,0,0,50,50);
    
}
//========================================================================
$(document).ready(function () {
    
    // Create a socket
    socket = new WebSocket('ws://' + window.location.host+ '/ws/join?uname=' +{{.UserName}});
});
//chat
function addMessage(textMessage) {

                //$("#msg-template .userpic").html("<img src='" + textMessage.UserInfo.Gravatar + "'>")
                $("#msg-template .msg-time").html(textMessage.Time);
                $("#msg-template .content").html(textMessage.Content);
                $("#chat-messages").append($("#msg-template").html());
                $('#chat-column')[0].scrollTop = $('#chat-column')[0].scrollHeight;
}
//control msg
function addControlMsg(textMessage){
    //$("#msg-template .userpic").html("<img src='" + textMessage.UserInfo.Gravatar + "'>")
                $("#msg-template .msg-time").html(textMessage.Time);
                $("#msg-template .content").html(textMessage.Content);
                $("#chat-messages").append($("#msg-template").html());
                $('#chat-column')[0].scrollTop = $('#chat-column')[0].scrollHeight;
            
}
function notifyMessage(msg) {
    $("#msg-template .content").html(msg);
    $("#chat-messages").append($("#msg-template").html());
    $('#chat-column')[0].scrollTop = $('#chat-column')[0].scrollHeight;
}
function updateReadyState(UserStatus){
    if(UserStatus.Me.Name!=localStorage.Me){
        //todo
        //update user ready state
        //canvas boarder :green
        if (UserStatus.CurrentStatus==101){
            for(x in Users){
                if(UserStatus.Me.Name==Users[x]){
                    switch(x){
                        $("#user"+x).css("boarder","red");
                    }
                }
            }
            
            addControlMsg(UserStatus.Me.Name+"Cancel Ready!");
        }else{
            for(x in Users){
                if(UserStatus.Me.Name==Users[x]){
                    switch(x){
                        $("#user"+x).css("boarder","green");
                    }
                }
            }
            addControlMsg(UserStatus.Me.Name+"Ready!");
        }
        
    }
    
    
}
function showUserSlected(UserSelected){
        
        //todo
        //update user select state
        var name;
        for(i in UserSelected){
            name=name+","+UserSelected[i];
        }

        addControlMsg("队长已选出人物人选！"+name);
    
}
function showVoteResult(UserStatus){
    //change leader
    for (x in Users){
        if(Users[x]==localStorage.Leader){
            $("#user"+x).css("boarder","1px");
        }
        if(Users[x]==UserStatus.UserStatus.Me.Leader){
            $("#user"+x).css("boarder","3px");
        }
    }

    
    localStorage.Leader=UserStatus.UserStatus.Me.Leader;
    //show vote result
    for( x in UserStatus.VoteResultMap){
        for(i in Users){
            if(x==Users[i]){
                if(UserStatus.VoteResultMap[x]){
                    $("#user"+i).css("boarder","yellow");   //approve : yellow //black
                }else{
                    $("#user"+i).css("boarder","black");
                }
                
            }
        }
        

    }

    if(UserStatus.VoteResult){
        $("#VoteResult_Popup").html("赞成过半,投票通过！");
        $("#VoteResult_Popup").click();
        addControlMsg("赞成过半,投票通过！");
    }else{
        $("#VoteResult_Popup").html("赞成不过半,投票未通过！");
        $("#VoteResult_Popup").click();
        addControlMsg("赞成不过半,投票未通过！");
    }
   
}
function showMissionResult(UserStatus){
    
    //show mission result
    //count down
    var s=0,f=0;
    for (x in UserStatus.MissionResultArray){
        if (UserStatus.MissionResultArray[x]){
            s++;
        }else{
            f++;
        }
    }
    addControlMsg(s+" 票成功, "+f+" 票失败!")

    if (UserStatus.MissionResult){
        $("#MissionResult_Popup").html("任务成功!");
        $("#MissionResult_Popup").click();
        addControlMsg("任务成功!");
    }else{
        $("#MissionResult_Popup").html("任务失败!");
        $("#MissionResult_Popup").click();
        addControlMsg("任务失败!");
    }
}
function showGameResult(data){
    
    //show game result
    //go to prepare phase
    addMessage(data.TextMessage);
    showRoles(data.UserStatus.ShowRole);

}
function getRole(introle){
    switch(introle){
        case 1:
            role="梅林"
            break;
        case 2:
        role="派西维尔"
            break;
        case 3:
        role="亚瑟的忠臣"
            break;
        case 4:
        role="莫德雷德"
            break;
        case 5:
        role="莫甘娜"
            break;
        case 6:
        role="奥伯伦"
            break;
        case 7:
        role="刺客"
            break;
        case 8:
        role="莫德雷德的爪牙"
            break;
        case 9:
        role="未知"
            break;
        default:
            break;
    }
    return role;
}
function leaderCanStart(UserStatus){
    
    //tell leader can start game
    if(UserStatus.Me.Name==localStorage.Me){
        $("#btn-start").parent("div").css('display','block');
    }
    addControlMsg("All users ready,"+UserStatus.Me.Name+" 请点击 start.");
}
//* 5人：2-3-2-3-3（均为出现一个任务失败就判定为任务失败）
//* 6人：2-3-4-3-4（均为出现一个任务失败就判定为任务失败）
//* 7人：2-3-3-4-4（第一个4人任务需要出现两个任务失败才判定为失败，其余只需要一个）
//* 8-10人：3-4-4-5-5（第一个5人任务需要出现两个任务失败才判定为失败，其余只需要一个）
function initMissionerSelect(){
    var number = Users.length;
    number = getMiissionerNumber();
    localStorage.numberofselect=number;
    $("#field-selectmissioner").empty();
    $("#field-selectmissioner").append("<label for='missioner'>您可以选择"+number+"名人选做任务:</label>");
    $("#field-selectmissioner").append("<select name='missioner' id='missioner' multiple='multiple' data-native-menu='false'>");
    for (x in Users){
        
        $("#missioner").append("<option value='"+Users[x]+"'>"+Users[x]+"</option>");
    }
    $("#field-selectmissioner").append("</select></fieldset>");
    $("#field-selectmissioner").append("<input id='submit-missioner' type='submit' data-inline='true' value='提交'><input id='reset-missiner' type='reset' data-inline='true' value='重置'>");
}
function getMiissionerNumber(){
    var ret=0;
    switch(Users.length){
        case 5:
            switch(localStorage.Round){
                case 1:
                    ret = 2;
                    break;
                case 2:
                ret = 3;
                    break;
                case 3:
                ret = 2;
                    break;
                case 4:
                ret = 3;
                    break;
                case 5:
                ret = 3;
                    break;
                default:
                    break;
            }
        case 6:
            switch(localStorage.Round){
                case 1:
                    ret = 2;
                    break;
                case 2:
                ret = 3;
                    break;
                case 3:
                ret = 4;
                    break;
                case 4:
                ret = 3;
                    break;
                case 5:
                ret = 4;
                    break;
                default:
                    break;
            }
        case 7:
            switch(localStorage.Round){
                case 1:
                    ret = 2;
                    break;
                case 2:
                ret = 3;
                    break;
                case 3:
                ret = 3;
                    break;
                case 4:
                ret = 4;
                    break;
                case 5:
                ret = 4;
                    break;
                default:
                    break;
            }
        default:    //8-10
            switch(localStorage.Round){
                case 1:
                    ret = 3;
                    break;
                case 2:
                ret = 4;
                    break;
                case 3:
                ret = 4;
                    break;
                case 4:
                ret = 5;
                    break;
                case 5:
                ret = 5;
                    break;
                default:
                    break;
            }
    }
    return ret;
}
function notifyPhaseChange(UserStatus){
    
    //notifyPhaseChange
    switch(data.UserStatus.CurrentStatus){
        case 104:
            addControlMsg("leader started the game");
            //hide ready/ start button
            $("#btn-ready").parent("div").css("display","none");
            $("#btn-start").parent("div").css("display","none");
            break;
        case 105:
            //根据人数初始化人物人选选择界面
            addControlMsg("天黑了，请梅林确认坏人，派西维尔确认梅林和莫甘娜，坏人相互确认！");

            break;
        case 110:
            addControlMsg("天亮了，开始第一回合！");
            initMissionerSelect();
            break; 
        case 111:
        addControlMsg("第二回合开始！");
        initMissionerSelect();
        break;
        case 112:
        addControlMsg("第三回合开始！");
        initMissionerSelect();
        break;
        case 113:
        addControlMsg("第四回合开始！");
        initMissionerSelect();
        break;
        case 114:
        addControlMsg("第五回合开始！");
        initMissionerSelect();
        break;
        case 122:
        addControlMsg("推迟已经达到4次直接开始任务！");
        break;
        default:
            break;
    }
}
function showRoles(ShowRole){
    
    //showRoles
    var role;
    //show role
    for (x in ShowRole){
        for(y in Users){
            if(x==Users[y]){
                switch(y){
                    case 0:
                        ctx1.fillText(getRole(data.UserStatus.ShowRole[x]),0,70);
                        break;
                    case 1:
                    ctx2.fillText(getRole(data.UserStatus.ShowRole[x]),0,70);
                        break;
                    case 2:
                    ctx3.fillText(getRole(data.UserStatus.ShowRole[x]),0,70);
                        break;
                    case 3:
                    ctx4.fillText(getRole(data.UserStatus.ShowRole[x]),0,70);
                        break;
                    case 4:
                    ctx5.fillText(getRole(data.UserStatus.ShowRole[x]),0,70);
                        break;
                    case 5:
                    ctx6.fillText(getRole(data.UserStatus.ShowRole[x]),0,70);
                        break;
                    case 6:
                    ctx7.fillText(getRole(data.UserStatus.ShowRole[x]),0,70);
                        break;
                    case 7:
                    ctx8.fillText(getRole(data.UserStatus.ShowRole[x]),0,70);
                        break;
                    case 8:
                    ctx9.fillText(getRole(data.UserStatus.ShowRole[x]),0,70);
                        break;
                    case 9:
                    ctx10.fillText(getRole(data.UserStatus.ShowRole[x]),0,70);
                        break;
                    default:
                    break;
                }
            }
        }
        
    }
    }
    
}
function popupMessage(msg){
    $("#Popup-message").empty();
    $("#Popup-message").append("<a href='#' data-rel='back' class='ui-btn ui-corner-all ui-shadow ui-btn ui-icon-delete ui-btn-icon-notext ui-btn-right'>关闭</a>");
    $("#Popup-message").append("<p><b>"+msg+"</b></p>");
    $("#btn-ppupmessage").click();
}
function countDown(UserStatus){
    
    //countDown
    
}
function assassinate(UserStatus){
    
    //assassinate
    addControlMsg("刺客开始刺杀!");
    if (UserStatus.Me.Name==localStorage.Me){
        $("#assassinate").click()
    }
}
function leaderstartSelect(UserStatus){
    
    //leaderstartSelect
    addControlMsg("队长开始选人!");
    if(UserStatus.Me.Name==localStorage.Me){
        $("#btn-missionerselect").click()
    }
}
function starttoSpeech(UserStatus){
    
    //starttoSpeech
    addControlMsg("开始发言!");
    addControlMsg("请 "+UserStatus.Me.Name+" 发言!");
    //countdown

}
function startToVote(UserStatus){
    
    //startToVote
    addControlMsg("开始对任务的人选进行投票!");
    $("#btn-vote").css('display','block');

}
function startToMission(UserStatus){
    
    //startToMission
    addControlMsg("开始做任务!");
    $("#btn-mission").css('display','block');
    
}
function gameControl(data) {
    switch (data.UserStatus.CurrentStatus){
        //C_PreparePhase，进入准备阶段，可以选择ready,unready
        case 0:
            break;
            //C_UnReady
        case 101:
            updateReadyState(data.UserStatus);
            break;
        //C_Ready,更新用户准备状态
        case 102:
            updateReadyState(data.UserStatus);
            break;
            //C_LeaderStart
        case 104:
            notifyPhaseChange(data.UserStatus);
            break;
        //C_LeaderFinSelect,//leader选人完毕
        case 116:
            showUserSlected(data.UserStatus.UserSelected);
            break;
        //C_VoteSuccAndStartMission
        case 121:
            showVoteResult(data.UserStatus);
            break;
        //C_VoteFailAndReSelect
        case 120:
            showVoteResult(data.UserStatus);
            break;
        //C_MissionEnd
        case 87:
            showMissionResult(data.UserStatus);
            break;
        //C_ShowGameResult
        case 9:
            showGameResult(data);
            break;
        //C_LeaderCanStart
        case 103:
            leaderCanStart(data.UserStatus);
            break;
        //C_GoIntoLightOffPhase
        case 105:
            notifyPhaseChange(data.UserStatus);
            break;
        //C_RolesToMerlin
        case 106:
            showRoles(data.UserStatus.ShowRole);
            break;
        //C_RolesToPercival
        case 108:
            showRoles(data.UserStatus.ShowRole);
            break;
        //C_RolesBadToBad
        case 107:
            showRoles(data.UserStatus.ShowRole);
            break;
        //C_CountDown
        case 109:
            countDown(data.UserStatus);
            break;
        //C_LightOnandR1start
        case 110:
            notifyPhaseChange(data.UserStatus);
            break;
        //C_R2start
        case 111:
        notifyPhaseChange(data.UserStatus);
            break;
        //C_R3start
        case 112:
        notifyPhaseChange(data.UserStatus);
            break;
        //C_R4start
        case 113:
        notifyPhaseChange(data.UserStatus);
            break;
        //C_R5start
        case 114:
        notifyPhaseChange(data.UserStatus);
            break;
        //C_DelayOver4TimesAndNoSpeech
        case 122:
        notifyPhaseChange(data.UserStatus);
            break;
        //C_Assassinate
        case 125:
            assassinate(data.UserStatus);
            break;
        //C_LeaderstartSelect
        case 115:
            leaderstartSelect(data.UserStatus);
            break;
        //C_StarttoSpeech
        case 117:
            starttoSpeech(data.UserStatus);
            break;
        //C_StartToVote
        case 118:
            startToVote(data.UserStatus);
            break;
        //C_StartToMission
        case 123:
            startToMission(data.UserStatus);
            break;
        default:
            break;
        //
    }
}
//用户加入房间消息
function updateGame(data) {
    if(typeof(Storage)!=="undefined")
      {
      // Yes! localStorage and sessionStorage support!
      // Some code.....
      }
    else
      {
      // Sorry! No web storage support..
      }
      if data.UserStatus.Me.Name!=localStorage.Me{
            addUser(data.UserStatus.Me.Name);
            freshUser();//fresh user list
            //print xxx join room
            //todo
      }else{
            //me is leader
            if (data.UserStatus.Me.Leader==localStorage.Me){
                //$("#btn-start").parent("div").css("display","block")
            }else{

                $("#btn-ready").parent("div").css("display","block");
            }
            initUsers(data.UserStatus.Users);
            freshUser();
      }
      //更新用户列表

      //todo
}
function freshUser(){
    var img_user = $("#icon-user");
    for (int x=0;x!=Users.length;x++)
    {
        switch(x){
            case 0:
                ctx1.drawImage(img_user,1,1,50,50);
                break;
            case 1:
                ctx2.drawImage(img_user,1,1,50,50);
                break;
            case 2:
                ctx3.drawImage(img_user,1,1,50,50);
                break;
            case 3:
                ctx4.drawImage(img_user,1,1,50,50);
                break;
            case 4:
                ctx5.drawImage(img_user,1,1,50,50);
                break;
            case 5:
                ctx6.drawImage(img_user,1,1,50,50);
                break;
            case 6:
                ctx7.drawImage(img_user,1,1,50,50);
                break;
            case 7:
                ctx8.drawImage(img_user,1,1,50,50);
                break;
            case 8:
                ctx9.drawImage(img_user,1,1,50,50);
                break;
            case 9:
                ctx10.drawImage(img_user,1,1,50,50);
                break;
        }
        $("#user-name"+x).html(Users[x]);
        
    }
}
function addUser(name){
    Users.push(name);
}
function initUsers(Users){
    Users.push(Users);
}
//click start
$("#btn-start").click(function(){
    msg = {MType:"control_mtype",UserStatus:{Me:{Name:localStorage.Me},CurrentStatus:104}};
    sendJson(msg);
});

    // Message received on the socket
socket.onmessage = function (event) {
    var data = JSON.parse(event.data);
    //console.log(data);
    switch (data.MType) {
    case "chat_mtype": // chat
        addMessage(data.TextMessage);
        break;

    case "notify_mtype": // notify
        notifyMessage(data);
        break;
    case "control_mtype": // game control
        gameControl(data);
        break;
    case "update_mtype": // update
        updateGame(data);
        break;
    default:
        break;
    }
};

    //Send Json
    var sendJson = function(data){
        socket.send(JSON.stringify(data));
        $('#sendbox').val("");
    }

    // Send messages.
    var postConecnt = function () {
        var uname = $('#uname').text();
        var content = $('#sendbox').val();
        socket.send(content);
        $('#sendbox').val("");
    }

    $('#chat-post').click(function () {
        msg = {MType:"chat_mtype",TextMessage:{Content:$("#text_input").val()},UserStatus:{Me:{Name:localStorage.Me}};
        sendJson(msg);
        $("#text_input").val("");
    });

    $('#submit-missioner').click(function () {
        if($("#missioner").val().length!=localStorage.numberofselect){
            popupMessage("选择人数错误，请重新选择！");
            //$("#reset-missiner").click();
        }else{
            msg = {MType:"control_mtype",UserStatus:{Me:{Name:localStorage.Me},CurrentStatus:116,UserSelected:$("#missioner").val()}};
            sendJson(msg);
        }
    });

    $('#btn-ready').click(function(){
        if($('#btn-ready').val()=="Ready"){
            msg = {MType:"control_mtype",UserStatus:{Me:{Name:localStorage.Me},CurrentStatus:102}};
            sendJson(msg);
            $('#btn-ready').val("UnReady");
        }else{
            msg = {MType:"control_mtype",UserStatus:{Me:{Name:localStorage.Me},CurrentStatus:101}};
            sendJson(msg);
            $('#btn-ready').val("Ready");
        }
        
    });
});