# IMConnection

## 目标

这是一个作业 demo，目的是创造一个轻量级的 IM 沟通软件后端。

## 主要依赖

```shell
# web structure
$ go get -u github.com/gin-gonic/gin

# database
$ go get -u gorm.io/gorm
$ go get -u gorm.io/driver/postgres

# util
$ go get -u github.com/go-ini/ini # config file reader
$ go get -u golang.org/x/crypto/bcrypt # transform Ciphertext
$ go get -u golang.org/x/oauth2 # token generate and check
```

## 问题

- [ ] 一对一消息传输 私聊
- [ ] 一对多消息传输 群聊
- [ ] 离线消息队列
- [ ] 邮箱注册？
- [ ] 聊天图片
  - [ ] 先用 redis 缓存，三天后存入数据库
  - [ ] 使用 MD5 判断文件是否相同？
    - [ ] 如果面临哈希碰撞，这样的情况下解决方案是什么？
- [ ] 聊天文件
- [ ] 个人信息展示
- [ ] 添加好友
- [ ] 好友关系
- [ ] 搜索群聊和用户
  - [ ] 条件筛选 (用户 Tag 和 群组 Tag)
  - [ ] SQL 在同一行内插入多个值？比如一列 Tag，每一行里可能 Tag 列里面有数个值
- [ ] 群聊关注消息？
- [ ] 用户好友关系 (好友备注)
- [ ] 即时通信，那么离线时是否推送消息？
- [ ] 离线时显示，user 离开群聊



## 需求分析

即时通讯系统 (IM System) 在许多领域都有运用，这里只是一个作业 demo，所以计划里只有两个功能，私信和群聊。

### 用户操作

#### JWT

在这个 demo 中，使用 JWT(JSON Web Token) 验证作为登录验证，使用 access token 和 refresh token 这样的一个 token 组。服务器使用 access token 进行验证操作，当 access token 过期但是 refresh token 未过期时，服务器会签发一个新的 access token 给客户端。

在这里使用的是 oauth2 lib，将用户 ID 和用户名作为 access token 的信息，以便于用户鉴权。而 refresh token 持续时间长，使用频率低，用于向服务端刷新 access token，包含了用户名和密码。

#### 好友关系

好友关系在用户未登录时，储存在 SQL 数据库中，当用户从客户端发送了登录请求后，用户 ID 会被服务器加载到内存，该用户的好友关系会被缓存到 redis 数据库中，以提高访问速度。

### 私信

