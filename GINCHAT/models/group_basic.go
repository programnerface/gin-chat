package models

import (
	"gorm.io/gorm"
)

// 群信息
type GroupBasic struct {
	gorm.Model
	Name    string //群聊名称
	OwnerId uint   //群主
	Icon    string //图片
	Type    int    //类型
	Desc    string //描述
}

func (table *GroupBasic) TableName() string {
	return "group_basic"
}
