package v1

import (
	"IMConnection/pkg/util"
	"IMConnection/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

func Chat(c *gin.Context) {
	claim, _ := util.ParseToken(c.GetHeader("Authorization"))
	sender := claim.Username
	receiver := c.Param("receiver")

	// 将 http 协议升级为 websocket 协议
	conn, err := (&websocket.Upgrader{
		// CheckOrigin 解决跨域问题
		CheckOrigin: func(r *http.Request) bool {
			fmt.Println("升级协议", "ua:", r.Header["User-Agent"], "referer:", r.Header["Referer"])
			return true
		}}).Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		http.NotFound(c.Writer, c.Request)
		return
	}
	log.Println("webSocket 建立连接:", conn.RemoteAddr().String())

	// Create A new Client Object
	client := &service.Client{
		SID:    sender,
		RID:    receiver,
		Socket: conn,
		Send:   make(chan []byte),
	}

	println("SID: ", sender, "; RID: ", receiver)
	println(client.SID, "; ", client.RID, "; ", client.Socket, "; ", client.Send)
	println(client.Socket.LocalAddr().String())

	// 用户注册到用户管理上
	service.Manager.Register <- client
	go client.Read()
	go client.Write()
}
