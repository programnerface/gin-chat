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
