package service

import (
	"IMConnection/model"
	"IMConnection/pkg/e"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"strconv"
)

func (manager *ClientManager) Listen() {
	for {
		println("==============Listening Channel Connection==============")
		select {
		/* 建立连接 */
		case conn := <-Manager.Register:
			log.Printf("Build New Connect: %v", conn.SID)
			Manager.Clients[conn.SID] = conn
			replyMsg := &MsgContent{
				SID:     "Client",
				RID:     conn.SID,
				Code:    e.Success,
				Content: "已连接至服务器",
			}
			msg, _ := json.Marshal(replyMsg)
			_ = conn.Socket.WriteMessage(websocket.TextMessage, msg)

		/* 断开连接 */
		case conn := <-Manager.Unregister:
			log.Printf("Fail to Connect:%v", conn.SID)
			if _, ok := Manager.Clients[conn.SID]; ok {
				replyMsg := &MsgContent{
					SID:     "Client",
					RID:     conn.SID,
					Code:    e.Success,
					Content: "服务器连接已断开",
				}
				msg, _ := json.Marshal(replyMsg)
				_ = conn.Socket.WriteMessage(websocket.TextMessage, msg)
				close(conn.Send)
				delete(Manager.Clients, conn.SID)
			}

		/* 消息广播 */
		case broadcast := <-Manager.Broadcast:
			log.Println("Broad messages to group member")
			println("SID", broadcast.Client.SID)
			println("RID", broadcast.Client.RID)

			for _, client := range Manager.Clients {
				println(client.SID)
			}

			message := broadcast.Message
			RID := broadcast.Client.RID
			flag := false // 默认对方不在线
			for id, conn := range Manager.Clients {
				if id != RID {
					continue
				}
				select {
				case conn.Send <- message:
					flag = true
				default:
					close(conn.Send)
					delete(Manager.Clients, conn.SID)
				}
			}
			sid, _ := strconv.Atoi(broadcast.Client.SID)
			rid, _ := strconv.Atoi(broadcast.Client.RID)
			if flag {
				log.Println("对方在线应答")
				replyMsg := &MsgContent{
					SID:     "Client",
					RID:     broadcast.Client.SID,
					Code:    e.Success,
					Content: "对方在线应答",
				}
				msg, _ := json.Marshal(replyMsg)
				_ = broadcast.Client.Socket.WriteMessage(websocket.TextMessage, msg)
				msgs := model.Message{
					SID:     uint(sid),
					RID:     uint(rid),
					Type:    broadcast.Type,
					Content: msg,
				}
				println("SID: ", msgs.SID)
				println("RID: ", msgs.RID)
				println("Content", string(msgs.Content))
				if err := model.DB.Create(&msgs).Error; err != nil {
					fmt.Println("InsertOneMsg Err", err)
				} else {
					println("here")
				}
			} else {
				log.Println("对方不在线")
				replyMsg := MsgContent{
					SID:     "Client",
					RID:     broadcast.Client.SID,
					Code:    e.Success,
					Content: "对方不在线应答",
				}
				msg, _ := json.Marshal(replyMsg)
				_ = broadcast.Client.Socket.WriteMessage(websocket.TextMessage, msg)
				msgs := model.Message{
					SID:     uint(sid),
					RID:     uint(rid),
					Type:    broadcast.Type,
					Content: msg,
				}
				if err := model.DB.Create(&msgs).Error; err != nil {
					fmt.Println("InsertOneMsg Err", err)
				}
			}
		}
	}
}
