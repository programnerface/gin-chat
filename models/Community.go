package models

import (
	"fmt"
	"ginchat/utils"
	"gorm.io/gorm"
)

// 群聊关系实体类
type Community struct {
	gorm.Model
	Name    string //群名称
	OwnerId uint   //群主
	Img     string //图片
	Desc    string //
}

func CreateCommunity(community Community) (int, string) {
	//群名称不能太短或为空
	if len(community.Name) == 0 {
		return -1, "群名称不能为空"
	}
	//群主id不能为空
	if community.OwnerId == 0 {
		return -1, "请先登录"
	}
	//插入数据库比那个判断异常
	if err := utils.DB.Create(&community).Error; err != nil {
		fmt.Println(err)
		return 1, "建群失败"
	}
	return 0, "建群成功"
}

func LoadCommunity(ownerId uint) ([]*Community, string) {
	data := make([]*Community, 10)
	utils.DB.Where(""+
		""+
		"owner_id=?", ownerId).Find(&data)
	//for循环打印出data集合
	for _, v := range data {
		fmt.Println(v)
	}
	//utils.DB.Where()
	return data, "查询成功"
}

// 加入群
func JoinGroup(userId uint, comId string) (int, string) {
	contact := Contact{}
	contact.OwnerId = userId
	contact.Type = 2
	community := Community{}
	utils.DB.Where("id=? or name=?", comId, comId).Find(&community)
	if community.Name == "" || community.ID == 0 {
		return -1, "没有找到群"
	}
	utils.DB.Where("owner_id=? and target_id=? and type=2", userId, comId).Find(&contact)
	//if contact.ID != 0 {
	if !contact.CreatedAt.IsZero() {
		return -1, "您已加入此群"
	} else {
		contact.TargetId = community.ID
		utils.DB.Create(&contact)
		return 0, "加群成功"
	}
}
