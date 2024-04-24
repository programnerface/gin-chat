package models

import (
	"fmt"
	"ginchat/utils"
	"gorm.io/gorm"
)

// 人员关系
type Contact struct {
	gorm.Model
	OwnerId  uint   //谁的关系信息
	TargetId uint   //对应的谁
	Type     int    //对应的类型 群聊等  1好友  2群组  3
	Desc     string //描述信息
}

func (table *Contact) TableName() string {
	return "contact"
}

// 查找好友
func SearchFriend(userId uint) []UserBasic {
	contacts := make([]Contact, 0)
	objIds := make([]uint64, 0)
	//查找contacts对象中ower_id等于userId且type=1的数据
	utils.DB.Where("owner_id = ? and type=1", userId).Find(&contacts)
	for _, v := range contacts {
		fmt.Println(">>>>>>>>>>>", v)
		objIds = append(objIds, uint64(v.TargetId))
	}
	users := make([]UserBasic, 0)
	utils.DB.Where("id in ?", objIds).Find(&users)
	return users
}

// 添加好友
func AddFriend(userId uint, targetId uint) int {
	user := UserBasic{}
	if targetId != 0 {
		user = FindByID(targetId)
		fmt.Println(targetId, " >>>>>", userId)
		//user.Identity 修改成 user.Salt 用户需要登录才会有Identity
		if user.Salt != "" {
			contact := Contact{}
			contact.OwnerId = userId
			contact.TargetId = targetId
			contact.Type = 1
			utils.DB.Create(&contact)
			return 0
		}
		return -1
	}
	return -1
}
