//=======================================================================
//define
console.log("hahahahahahha");
var socket;
localStorage.Me={{.UserName}};
localStorage.Leader="";
localStorage.Currentstatus = 0;
localStorage.Round = 0;
localStorage.Role = -1;
localStorage.numberofselect=0;
var Users = [localStorage.Me];
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
    $("#btn-ready").val("hahahaha");
    // Create a socket
    socket = new WebSocket('ws://' + window.location.host+ '/ws/join?uname=' +{{.UserName}});
};