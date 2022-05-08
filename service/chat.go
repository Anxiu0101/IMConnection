package service

import (
	"IMConnection/cache"
	"IMConnection/pkg/e"
	"IMConnection/pkg/logging"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

type Client struct {
	SID    string
	RID    string
	Socket *websocket.Conn
	Send   chan []byte
}

// MsgContent 发送消息的类型
type MsgContent struct {
	Type    int    `json:"type"`
	Content string `json:"content"`
	Code    int    `json:"code"`
}

type StatusReply struct {
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

	println("Here")

	// 无限循环保持连接 这个算轮询吗？
	for {
		client.Socket.PongHandler()
		msg := new(MsgContent)
		fmt.Println("Msg: ", msg)
		// _,msg,_:=c.Socket.ReadMessage()
		err := client.Socket.ReadJSON(&msg) // 读取json格式，如果不是json格式，会报错
		if err != nil {
			logging.Info("数据格式不正确", err)
			Manager.Unregister <- client
			_ = client.Socket.Close()
			break
		}

		// 信息类型为私信
		if msg.Type == SingleChat {

			println("Here")
			// 查看该联系 ID 的连接个数
			r1, _ := cache.RedisClient.Get(cache.Ctx, client.SID).Result()
			r2, _ := cache.RedisClient.Get(cache.Ctx, client.RID).Result()
			// 限制单聊个数
			if r1 >= "3" && r2 == "" {
				replyMsg := MsgContent{
					Code:    e.Error,
					Content: "达到限制",
				}
				msg, _ := json.Marshal(replyMsg)
				_ = client.Socket.WriteMessage(websocket.TextMessage, msg)
				// 设置 key 过期时间，防止重复骚扰
				_, _ = cache.RedisClient.Expire(cache.Ctx, client.SID, time.Hour*24*30).Result()
				continue
			} else {
				// 单聊个数未超过上限，key 为 ID 的 value 自增一
				cache.RedisClient.Incr(cache.Ctx, client.SID)
				// 设置过期时间为一日
				_, _ = cache.RedisClient.Expire(cache.Ctx, client.SID, time.Hour*24).Result()
			}
			log.Println(client.SID, "发送消息", msg.Content)
			Manager.Broadcast <- &Broadcast{
				Client:  client,
				Message: []byte(msg.Content),
			}
			// 信息类型为群聊
		} else if msg.Type == GroupChat {

		} else if msg.Type == History {

		}
	}
}

func (client *Client) Write() {
	defer func() {
		_ = client.Socket.Close()
	}()

	for {
		select {
		case message, ok := <-client.Send:
			if !ok {
				_ = client.Socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			log.Println(client.SID, "接受消息:", string(message))
			replyMsg := MsgContent{
				Code:    e.Error,
				Content: fmt.Sprintf("%s", string(message)),
			}
			msg, _ := json.Marshal(replyMsg)
			_ = client.Socket.WriteMessage(websocket.TextMessage, msg)
		}
	}
}
