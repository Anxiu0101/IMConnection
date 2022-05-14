package service

import (
	"IMConnection/model"
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
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
					Content: "服务器连接已断开",
				}
				msg, _ := json.Marshal(replyMsg)
				_ = conn.Socket.WriteMessage(websocket.TextMessage, msg)
				close(conn.Send)
				delete(Manager.Clients, conn.SID)
			}

		/* 消息广播 */
		case broadcast := <-Manager.Broadcast:
			println("2. List Online Users")
			for _, client := range Manager.Clients {
				println(client.SID)
			}

			println("broadcast.Client.RID: ", broadcast.Client.RID)
			println("Content: ", string(broadcast.Message))
			println("Type: ", broadcast.Type)
			//RID := model.GetUserList(broadcast.Client.RID, conf.AppSetting.PageSize)
			flag := false // 默认对方不在线
			i := 0

			if broadcast.Type == 1 {
				for id, conn := range Manager.Clients {
					if id != broadcast.Client.RID {
						continue
					}
					select {
					case conn.Send <- broadcast.Message:
						flag = true
						println("Flag become true here")
					default:
						close(conn.Send)
						println("conn SID: ", conn.SID)
						delete(Manager.Clients, conn.SID)
					}
				}
			} else if broadcast.Type == 2 {
				var group model.Group
				model.DB.Model(model.Group{}).Where("name = ?", broadcast.Client.RID).Find(&group)
				var members []model.User
				model.DB.Model(&group).Select("username").Association("Members").Find(&members)
				for id, conn := range Manager.Clients {
					// 判断用户是否为需要的用户，不是则寻找下一个用户
					// TODO 遍历的性能开销大，尝试使用其他方法
					println(len(Manager.Clients))
					println("ID:", id)
					println("conn.SID:", conn.SID)
					if conn.SID != members[i].UserName {
						println("i:", i)
						println("username:", members[i].UserName)
						i++
						continue
					}
					select {
					case conn.Send <- broadcast.Message:
						// 信息广播给指定用户
						flag = true
						println("3. Flag become true here")
					default:
						close(conn.Send)
						println("conn SID: ", conn.SID)
						delete(Manager.Clients, conn.SID)
					}
				}
			}

			// 服务器应答
			if flag {
				log.Println("对方在线应答")
				replyMsg := &Msg{
					SID:     "Client",
					RID:     broadcast.Client.SID,
					Content: "对方在线应答",
				}
				msg, _ := json.Marshal(replyMsg)
				_ = broadcast.Client.Socket.WriteMessage(websocket.TextMessage, msg)
			} else {
				log.Println("对方不在线")
				replyMsg := Msg{
					SID:     "Client",
					RID:     broadcast.Client.SID,
					Content: "对方不在线应答",
				}
				msg, _ := json.Marshal(replyMsg)
				_ = broadcast.Client.Socket.WriteMessage(websocket.TextMessage, msg)
			}
		}
	}
}
