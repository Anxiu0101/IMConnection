package service

import (
	"IMConnection/cache"
	"IMConnection/pkg/logging"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"strconv"
	"time"
)

type Client struct {
	SID    string
	RID    string
	Socket *websocket.Conn
	Send   chan []byte
}

// Msg 发送消息的类型
type Msg struct {
	SID     uint   `json:"sid"`
	RID     uint   `json:"rid"`
	Type    int    `json:"type"`
	Content string `json:"content"`
}

// Broadcast 广播类，包括广播内容和源用户
type Broadcast struct {
	Client  *Client
	Message []byte
	Type    int
}

// Type of Message
const (
	SingleChat = 1
	GroupChat  = 2
	History    = 3
)

// ClientManager 用户管理
type ClientManager struct {
	Clients    map[string]*Client // 全部的连接
	Broadcast  chan *Broadcast    // 广播
	Reply      chan *Client
	Register   chan *Client // 连接连接处理
	Unregister chan *Client // 断开连接处理
}

var Manager = ClientManager{
	Clients:    make(map[string]*Client), // 参与连接的用户，出于性能的考虑，需要设置最大连接数
	Broadcast:  make(chan *Broadcast),
	Register:   make(chan *Client),
	Reply:      make(chan *Client),
	Unregister: make(chan *Client),
}

func (client *Client) Read() {
	defer func() {
		Manager.Unregister <- client
		_ = client.Socket.Close()
	}()

	// 无限循环保持连接 这个算轮询吗？

	for {
		client.Socket.PongHandler()
		msg := new(Msg)
		// _,msg,_:=c.Socket.ReadMessage()
		err := client.Socket.ReadJSON(&msg) // 读取json格式，如果不是json格式，会报错
		if err != nil {
			logging.Info("数据格式不正确", err)
			Manager.Unregister <- client
			_ = client.Socket.Close()
			break
		}
		if msg.Type == 1 {
			r1, _ := cache.RedisClient.Get(cache.Ctx, client.SID).Result()
			r2, _ := cache.RedisClient.Get(cache.Ctx, client.RID).Result()
			if r1 >= "3" && r2 == "" { // 限制单聊
				replyMsg := ReplyMsg{
					Code:    e.WebsocketLimit,
					Content: "达到限制",
				}
				msg, _ := json.Marshal(replyMsg)
				_ = client.Socket.WriteMessage(websocket.TextMessage, msg)
				_, _ = cache.RedisClient.Expire(cache.Ctx, client.SID, time.Hour*24*30).Result() // 防止重复骚扰，未建立连接刷新过期时间一个月
				continue
			} else {
				cache.RedisClient.Incr(cache.Ctx, client.SID)
				_, _ = cache.RedisClient.Expire(cache.Ctx, client.SID, time.Hour*24*30*3).Result() // 防止过快“分手”，建立连接三个月过期
			}
			log.Println(client.SID, "发送消息", msg.Content)
			Manager.Broadcast <- &Broadcast{
				Client:  client,
				Message: []byte(msg.Content),
			}
		} else if msg.Type == 2 { //拉取历史消息
			timeT, err := strconv.Atoi(sendMsg.Content) // 传送来时间
			if err != nil {
				timeT = 999999999
			}
			results, _ := FindMany(conf.MongoDBName, c.SendID, c.ID, int64(timeT), 10)
			if len(results) > 10 {
				results = results[:10]
			} else if len(results) == 0 {
				replyMsg := ReplyMsg{
					Code:    e.WebsocketEnd,
					Content: "到底了",
				}
				msg, _ := json.Marshal(replyMsg)
				_ = client.Socket.WriteMessage(websocket.TextMessage, msg)
				continue
			}
			for _, result := range results {
				replyMsg := ReplyMsg{
					From:    result.From,
					Content: fmt.Sprintf("%s", result.Msg),
				}
				msg, _ := json.Marshal(replyMsg)
				_ = c.Socket.WriteMessage(websocket.TextMessage, msg)
			}
		} else if msg.Type == 3 {
			results, err := FirsFindtMsg(conf.MongoDBName, c.SendID, c.ID)
			if err != nil {
				log.Println(err)
			}
			for _, result := range results {
				replyMsg := ReplyMsg{
					From:    result.From,
					Content: fmt.Sprintf("%s", result.Msg),
				}
				msg, _ := json.Marshal(replyMsg)
				_ = c.Socket.WriteMessage(websocket.TextMessage, msg)
			}
		}
	}
}

func (client *Client) Write() {
	//defer func() {
	//	_ = c.Socket.Close()
	//}()
	//for {
	//	select {
	//	case message, ok := <-c.Send:
	//		if !ok {
	//			_ = c.Socket.WriteMessage(websocket.CloseMessage, []byte{})
	//			return
	//		}
	//		log.Println(c.ID, "接受消息:", string(message))
	//		replyMsg := ReplyMsg{
	//			Code:    e.WebsocketSuccessMessage,
	//			Content: fmt.Sprintf("%s", string(message)),
	//		}
	//		msg, _ := json.Marshal(replyMsg)
	//		_ = c.Socket.WriteMessage(websocket.TextMessage, msg)
	//	}
	//}
}
