package main

import (
	"AvalonM/controllers"
	"code.google.com/p/go.net/websocket"
	"github.com/astaxie/beego"
	"github.com/beego/i18n"
)

const APP_VER = "v0.1"

func main() {
	beego.Info(beego.BConfig.AppName, APP_VER)

	// Register routers.
	beego.Router("/", &controllers.AppController{})
	// Indicate AppController.Join method to handle POST requests.
	beego.Router("/join", &controllers.AppController{}, "post:Join")

	// WebSocket.
	beego.Router("/ws", &controllers.WebSocketController{})
	//beego.Router("/ws/join", &controllers.WebSocketController{}, "get:Join")

	beego.Handler("/ws/join", websocket.Handler(controllers.BuildConnection))

	// Register template functions.
	beego.AddFuncMap("i18n", i18n.Tr)
	//go controllers.InitChatRoom()
	beego.Run()
}
