package cache

const (
	FriendList = "imc:friends:list"
)

func GetUserOnline(uid int) bool {
	result, _ := RedisClient.Get(Ctx, string(uid)).Result()
	return result != "0"
}
