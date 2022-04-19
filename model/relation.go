package model

type Relation struct {
	ID       uint `json:"id" gorm:"primarykey"`
	UserID   uint
	FriendID uint
	Group    string `json:"group" gorm:"size:64;index"`
}

// LoadRelation 将用户好友关系载入缓存
func LoadRelation() error {
	var err error
	return err
}

// UploadRelation 将用户关系从缓存中更新
func UploadRelation() {

}
