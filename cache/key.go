package cache

import (
	"fmt"
	"strconv"
)

const (
	FriendList = "imc:friends:list"
)

func GetUserOnline(uid int) string {
	return fmt.Sprintf("imc:user:online:%s", strconv.Itoa(uid))
}
