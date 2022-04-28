package v1

import (
	"IMConnection/pkg/util"
	"IMConnection/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"strconv"
)

func Chat(c *gin.Context) {
	claim, _ := util.ParseToken(c.GetHeader("Authorization"))
	sender := int(claim.ID)
	receiver := c.Query("receiver")

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
	fmt.Println("webSocket 建立连接:", conn.RemoteAddr().String())

	// Create A new Client Object
	client := &service.Client{
		SID:    strconv.Itoa(sender),
		RID:    receiver,
		Socket: conn,
		Send:   make(chan []byte),
	}

	// 用户注册到用户管理上
	service.Manager.Register <- client
	go client.Read()
	go client.Write()
}
