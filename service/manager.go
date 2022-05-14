package service

import (
	"IMConnection/model"
	"IMConnection/pkg/e"
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/gorilla/websocket"
)

func (manager *ClientManager) Listen() {
	for {
		println("=======================Listening Channel Connection=======================")
		select {
		/* 建立连接 */
		case conn := <-Manager.Register:
			log.Printf("Build New Connect: %v", conn.SID)
			Manager.Clients[conn.SID] = conn
			replyMsg := &Msg{
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
				replyMsg := &Msg{
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
			for _, client := range Manager.Clients {
				println(client.SID)
			}

			message := broadcast.Message
			RID := broadcast.Client.RID
			println("broadcast.Client.RID: ", broadcast.Client.RID)
			println("Content: ", string(broadcast.Message))
			//RID := model.GetUserList(broadcast.Client.RID, conf.AppSetting.PageSize)
			flag := false // 默认对方不在线
			for id, conn := range Manager.Clients {
				if id != RID {
					continue
				}
				select {
				case conn.Send <- message:
					flag = true
					println("Flag become true here")
				default:
					close(conn.Send)
					println("conn SID: ", conn.SID)
					delete(Manager.Clients, conn.SID)
				}
			}
			sid, _ := strconv.Atoi(broadcast.Client.SID)
			rid, _ := strconv.Atoi(broadcast.Client.RID)
			if flag {
				log.Println("对方在线应答")
				replyMsg := &Msg{
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
				if err := model.DB.Create(&msg).Error; err != nil {
					fmt.Println("InsertOneMsg Err", err)
				} else {
					println("here")
				}
			} else {
				log.Println("对方不在线")
				replyMsg := Msg{
					SID:     "Client",
					RID:     broadcast.Client.SID,
					Code:    e.Success,
					Content: "对方不在线应答",
				}
				msg, _ := json.Marshal(replyMsg)
				_ = broadcast.Client.Socket.WriteMessage(websocket.TextMessage, msg)
			}
		}
	}
}
