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
$ go get -u github.com/go-redis/redis/v8

# util
$ go get -u github.com/go-ini/ini       # config file reader
$ go get -u golang.org/x/crypto/bcrypt  # transform Ciphertext
$ go get -u github.com/dgrijalva/jwt-go # token generate and check
$ go get -u github.com/unknwon/com      # util package
```

## 问题

- [x] 一对一消息传输 私聊
- [x] 一对多消息传输 群聊
- [ ] 离线消息队列
- [ ] 邮箱注册？
- [ ] 聊天图片
  - [ ] 先用 redis 缓存，三天后存入数据库
  - [ ] 使用 MD5 判断文件是否相同？
    - [ ] 如果面临哈希碰撞，这样的情况下解决方案是什么？
  - [ ] 图片压缩，缩略图
- [ ] 聊天文件
- [x] 个人信息展示
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

在这里使用的是 `oauth2` lib，将用户 ID 和用户名作为 access token 的信息，以便于用户鉴权。而 refresh token 持续时间长，使用频率低，用于向服务端刷新 access token，包含了用户名和密码。

`oauth2` lib 没用懂，`oauth` 框架是区别于 JWT 的用户认证鉴权方法，其客户端包含着 https 认证的重定向网址。故依旧使用 JWT 来进行用户验证，在用户登录时，生成 token pair 发给用户，用户使用 access token 进行请求，当 access 的有效期(10 mins)过了，再次发送请求时，服务器返回 401，前端通过过滤器使用 refresh token 向服务器发送刷新 access token 请求，获取新 token 后重新发送请求。退出登录操作由客户端删除 token 这样的方法来实现，不需要后端进行操作。

#### 好友关系

好友关系在用户未登录时，储存在 SQL 数据库中，当用户从客户端发送了登录请求后，用户 ID 会被服务器加载到内存，该用户的好友关系会被缓存到 redis 数据库中，以提高访问速度。

利用 `gorm` 的 `ManyToMany` 模式在 User 与 User 之间创建一张自引用连接表

### Many2Many

```go
package main

import (
	"ginLearn.com/models"
)

func main() {
	db := models.DB()
	role := models.Role{}
	role.ID = 8
	var permissionSlice []models.Permission
	//根据角色查询所有权限
	db.Model(&role).Related(&permissionSlice, "permission")

	//我们更新角色的所有权限该怎么做呢？
	//1、删除角色的所有的权限
	db.Where("role_id=?", role.ID).Unscoped().Delete(&models.RolePermission{})
	//2、给角色赋予权限
	role.Permission = []models.Permission{permissionSlice[0]}
	//更新角色的权限
	db.Save(&role)
}


```



### 通信

将服务拆分成多个客户端，包括 UserLogin, UserLogout, Broadcast 等，当用户使用了其中的某个服务时，将其推入到 channel 中传递进程信息，已启用服务。

#### 消息队列

所有用户都有他对应的消息队列，用户收到的任何消息都会储存到这个消息队列中，消息队列设置有效期，用户在离线后只能获取最近7日内的消息，而不是全部未读消息。

#### 消息

由于场景的不同，私聊与群聊的消息做的工作实际上是不一样的，私聊是单对单的信息发送，结构体中需要包含发送者与接收者。而群聊是发送者发送到群聊中后，由服务器广播给群聊的其他用户。

![imgchat-mindmap](https://raw.githubusercontent.com/Anxiu0101/LectureNote4Img/master/static/imgchat-mindmap.png)

// TODO 通过修改和定义 SenderType 来实现

消息可以设置消息类型，可以携带不同类型的信息而不仅仅是文本，例如定位信息、表情信息、网址信息等，可以为网址信息添加快照这样的。

#### 群组

群组包含用户列表，如果有人发送了消息，可以将群聊一定时间内的历史消息纪律加载入缓存中，或者为群聊设置一定大小的消息缓存，因为实际上用户爬楼的意愿是比较低的。

再仔细想了想，将群聊消息记录持久化即可，想过头了

@2022-5-12 现在的问题是，单聊和群聊两个功能都使用了 client 的 Broadcast 服务，但是现在的 Broadcast 服务无法将 Group Name 识别为正常的 ID 导致其持续放回对方未在线而阻塞消息。我的目的是识别 Group Name 并将其群成员查询出来，循环发送广播消息。我在返回 ID 这一块纠结了，其实可以直接返回 用户列表，在比较时直接通过用户来获取 ID

#### 文件上传

群聊文件的上传这个功能还没有想好

## 缓存

缓存字段统一以 `imc_` 开头，需要缓存的只有用户的好友关系
