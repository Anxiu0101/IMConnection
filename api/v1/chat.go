package v1

import (
	"github.com/gin-gonic/gin"
)

func Chat(c *gin.Context) {
	//uid := c.Query("uid")     // 自己的id
	//toUid := c.Query("toUid") // 对方的id
	//conn, err := (&websocket.Upgrader{
	//	CheckOrigin: func(r *http.Request) bool { // CheckOrigin解决跨域问题
	//		return true
	//	}}).Upgrade(c.Writer, c.Request, nil) // 升级成ws协议
	//if err != nil {
	//	http.NotFound(c.Writer, c.Request)
	//	return
	//}
	// 创建一个用户实例
	//client := &service.Client{
	//	ID:     createId(uid, toUid),
	//	SendID: createId(toUid, uid),
	//	Socket: conn,
	//	Send:   make(chan []byte),
	//}
	//// 用户注册到用户管理上
	//Manager.Register <- client
	//go client.Read()
	//go client.Write()
}
