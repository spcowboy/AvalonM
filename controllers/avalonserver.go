package controllers

import (
	. "AvalonM/models"
	"code.google.com/p/go.net/websocket"
	"github.com/astaxie/beego"
	//"strings"
)

// WebSocketController handles WebSocket requests.
type WebSocketController struct {
	baseController
}

// Get method handles GET requests for WebSocketController.
func (this *WebSocketController) Get() {
	// Safe check.
	beego.Info("websocket redirect.")
	uname := this.GetString("uname")
	if len(uname) == 0 {
		this.Redirect("/", 302)
		return
	}

	this.TplName = "test.html"
	this.Data["IsWebSocket"] = true
	this.Data["UserName"] = uname

}

func BuildConnection(ws *websocket.Conn) {
	uname := ws.Request().URL.Query().Get("uname")

	if uname == "" {
		beego.Info("uname is nil.")
		return
	}
	beego.Info("new user "+uname+" join the room. ")
	RunningGameRoom.Join(uname, ws)

}

