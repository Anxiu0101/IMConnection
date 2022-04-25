package v1

import (
	"IMConnection/pkg/util"
	"IMConnection/service"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/unknwon/com"
	"net/http"
	"strconv"
)

func SingleChat(c *gin.Context) {
	claim, _ := util.ParseToken(c.GetHeader("Authorization"))
	sid := claim.ID
	rid := com.StrTo(c.Query("rid")).MustInt()

	// 将 http 协议升级为 websocket 协议
	conn, err := (&websocket.Upgrader{
		// CheckOrigin 解决跨域问题
		CheckOrigin: func(r *http.Request) bool {
			return true
		}}).Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		http.NotFound(c.Writer, c.Request)
		return
	}

	// Create A new Client Object
	client := &service.Client{
		SID:    service.CreateClientID(int(sid), rid),
		RID:    service.CreateClientID(rid, int(sid)),
		Socket: conn,
		Send:   make(chan []byte),
	}

	// 用户注册到用户管理上
	service.Manager.Register <- client
	go client.Read()
	go client.Write()
}

func GroupChat(c *gin.Context) {
	claim, _ := util.ParseToken(c.GetHeader("Authorization"))
	sid := claim.ID
	group := c.Query("group")

	// 将 http 协议升级为 websocket 协议
	conn, err := (&websocket.Upgrader{
		// CheckOrigin 解决跨域问题
		CheckOrigin: func(r *http.Request) bool {
			return true
		}}).Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		http.NotFound(c.Writer, c.Request)
		return
	}

	// Create A new Client Object
	client := &service.Client{
		SID:    strconv.Itoa(int(sid)),
		RID:    group,
		Socket: conn,
		Send:   make(chan []byte),
	}

	// 用户注册到用户管理上
	service.Manager.Register <- client
	go client.Read()
	go client.Write()
}
