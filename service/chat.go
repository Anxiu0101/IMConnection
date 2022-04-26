package service

import (
	"github.com/gorilla/websocket"
)

type Client struct {
	SID    string
	RID    string
	Socket *websocket.Conn
	Send   chan []byte
}

// Broadcast 广播类，包括广播内容和源用户
type Broadcast struct {
	Client  *Client
	Message []byte
	Type    int
}

// ClientManager 用户管理
type ClientManager struct {
	Clients    map[string]*Client
	Broadcast  chan *Broadcast
	Reply      chan *Client
	Register   chan *Client
	Unregister chan *Client
}

var Manager = ClientManager{
	Clients:    make(map[string]*Client), // 参与连接的用户，出于性能的考虑，需要设置最大连接数
	Broadcast:  make(chan *Broadcast),
	Register:   make(chan *Client),
	Reply:      make(chan *Client),
	Unregister: make(chan *Client),
}

//func CreateClientID(SID int, RID int) string {
//	sid := strconv.Itoa(SID)
//	rid := strconv.Itoa(RID)
//	return sid + "->" + rid
//}

func (c *Client) Read() {
	//defer func() { // 避免忘记关闭，所以要加上close
	//	Manager.Unregister <- c
	//	_ = c.Socket.Close()
	//}()
	//for {
	//	c.Socket.PongHandler()
	//	sendMsg := new(SendMsg)
	//	// _,msg,_:=c.Socket.ReadMessage()
	//	err := c.Socket.ReadJSON(&sendMsg) // 读取json格式，如果不是json格式，会报错
	//	if err != nil {
	//		log.Println("数据格式不正确", err)
	//		Manager.Unregister <- c
	//		_ = c.Socket.Close()
	//		break
	//	}
	//	if sendMsg.Type == 1 {
	//		r1, _ := cache.RedisClient.Get(c.ID).Result()
	//		r2, _ := cache.RedisClient.Get(c.SendID).Result()
	//		if r1 >= "3" && r2 == "" { // 限制单聊
	//			replyMsg := ReplyMsg{
	//				Code:    e.WebsocketLimit,
	//				Content: "达到限制",
	//			}
	//			msg, _ := json.Marshal(replyMsg)
	//			_ = c.Socket.WriteMessage(websocket.TextMessage, msg)
	//			_, _ = cache.RedisClient.Expire(c.ID, time.Hour*24*30).Result() // 防止重复骚扰，未建立连接刷新过期时间一个月
	//			continue
	//		} else {
	//			cache.RedisClient.Incr(c.ID)
	//			_, _ = cache.RedisClient.Expire(c.ID, time.Hour*24*30*3).Result() // 防止过快“分手”，建立连接三个月过期
	//		}
	//		log.Println(c.ID, "发送消息", sendMsg.Content)
	//		Manager.Broadcast <- &Broadcast{
	//			Client:  c,
	//			Message: []byte(sendMsg.Content),
	//		}
	//	} else if sendMsg.Type == 2 { //拉取历史消息
	//		timeT, err := strconv.Atoi(sendMsg.Content) // 传送来时间
	//		if err != nil {
	//			timeT = 999999999
	//		}
	//		results, _ := FindMany(conf.MongoDBName, c.SendID, c.ID, int64(timeT), 10)
	//		if len(results) > 10 {
	//			results = results[:10]
	//		} else if len(results) == 0 {
	//			replyMsg := ReplyMsg{
	//				Code:    e.WebsocketEnd,
	//				Content: "到底了",
	//			}
	//			msg, _ := json.Marshal(replyMsg)
	//			_ = c.Socket.WriteMessage(websocket.TextMessage, msg)
	//			continue
	//		}
	//		for _, result := range results {
	//			replyMsg := ReplyMsg{
	//				From:    result.From,
	//				Content: fmt.Sprintf("%s", result.Msg),
	//			}
	//			msg, _ := json.Marshal(replyMsg)
	//			_ = c.Socket.WriteMessage(websocket.TextMessage, msg)
	//		}
	//	} else if sendMsg.Type == 3 {
	//		results, err := FirsFindtMsg(conf.MongoDBName, c.SendID, c.ID)
	//		if err != nil {
	//			log.Println(err)
	//		}
	//		for _, result := range results {
	//			replyMsg := ReplyMsg{
	//				From:    result.From,
	//				Content: fmt.Sprintf("%s", result.Msg),
	//			}
	//			msg, _ := json.Marshal(replyMsg)
	//			_ = c.Socket.WriteMessage(websocket.TextMessage, msg)
	//		}
	//	}
	//}
}

func (c *Client) Write() {
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
