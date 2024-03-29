package service

import (
	"IMConnection/cache"
	"IMConnection/model"
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
	SID     string `json:"sid"`
	RID     string `json:"rid"`
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
	// 当用户与服务器断开连接是，关闭 websocket 连接
	defer func() {
		Manager.Unregister <- client
		_ = client.Socket.Close()
	}()

	// 无限循环保持连接 这个算轮询吗？
	for {
		client.Socket.PongHandler()
		msg := new(Msg)
		// _,msg,_:=c.Socket.ReadMessage()
		// 将用户发送的 JSON 消息进行参数绑定
		if err := client.Socket.ReadJSON(msg); err != nil {
			logging.Info("非 JSON 格式，数据格式不正确", err)
			Manager.Unregister <- client
			_ = client.Socket.Close()
			break
		}
		msg.SID = client.SID
		msg.RID = client.RID

		// 信息类型为私信
		if msg.Type == SingleChat {
			// 查看该联系 ID 的连接个数
			r1, _ := cache.RedisClient.Get(cache.Ctx, client.SID).Result()
			r2, _ := cache.RedisClient.Get(cache.Ctx, client.RID).Result()
			// 限制单聊未回应消息个数
			if r1 >= "3" && r2 == "" {
				replyMsg := Msg{
					Content: "达到限制",
				}
				println("r1: ", r1)
				println("r2: ", r2)
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
			log.Println("1. ", client.SID, "发送消息", msg.Content)
			rid, _ := strconv.Atoi(msg.RID)
			sid, _ := strconv.Atoi(msg.SID)
			message := model.Message{
				RID:     uint(rid),
				SID:     uint(sid),
				Type:    msg.Type,
				Content: msg.Content,
			}
			if err := model.DB.Create(&message).Error; err != nil {
				log.Println("Save Message Error")
			}
			Manager.Broadcast <- &Broadcast{
				Type:    1,
				Client:  client,
				Message: []byte(msg.Content),
			}
			// 信息类型为群聊
		} else if msg.Type == GroupChat {
			log.Println(client.SID, "发送消息", msg.Content)
			//data := model.GetUserList(client.RID, conf.AppSetting.PageSize)
			//client.RID = list

			var group model.Group
			model.DB.Model(model.Group{}).Where("name = ?", client.RID).Find(&group)

			var members []uint
			model.DB.Model(&group).Select("user_id").Association("Members").Find(&members)

			var message = model.GroupMessage{
				SID:     msg.SID,
				RID:     msg.RID,
				Type:    msg.Type,
				Content: msg.Content,
			}
			model.DB.Model(model.GroupMessage{}).Create(&message)
			Manager.Broadcast <- &Broadcast{
				Client:  client,
				Message: []byte(msg.Content),
				Type:    2,
			}
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
			replyMsg := Msg{
				SID:     client.SID,
				RID:     client.RID,
				Type:    0,
				Content: fmt.Sprintf("%s", string(message)),
			}
			msg, _ := json.Marshal(replyMsg)
			_ = client.Socket.WriteMessage(websocket.TextMessage, msg)
		}
	}
}
