package models

import (
	"fmt"
	"ginchat/utils"
	"gorm.io/gorm"
	"time"
)

type UserBasic struct {
	gorm.Model
	//用户名
	Name string
	//密码
	PassWord string
	//手机号
	Phone string `valid:"matches(^1[3-9]{1}\d{9}$)"`
	//邮箱
	Email string `valid:"email"`
	//客户端ip
	ClentIp string
	//唯一标识
	Identity string
	//客户端口号
	ClientPort string
	//随机数
	Salt string
	//登陆时间
	LoginTime time.Time `gorm:"time.now()"`
	//心跳检测时间
	HeartbeatTime time.Time `gorm:"default:'2024-03-19 00:29:07.602'"`
	//登出时间
	LoginOutTime time.Time `gorm:"column: login_out_time" json:"loginOutTime"`
	//是否登出
	ISLogout bool
	//设备信息
	DeviceInfo string
}

// 'gorm:"column:login_out_time" json:"login_out_time"'
func (table *UserBasic) TableName() string {
	return "user_basic"
}

// 返回一个 []*UserBasic的切片
func GetUserList() []*UserBasic {
	//大小为10的data集合并赋值给 data
	data := make([]*UserBasic, 10)
	utils.DB.Find(&data)
	//for循环打印出data集合
	for _, v := range data {
		fmt.Println(v)
	}

	return data
}

// 登录
func FindUserByNameAndPwd(name string, password string) UserBasic {
	user := UserBasic{}
	utils.DB.Where("name = ? and pass_word=?", name, password).First(&user)

	//token加密
	str := fmt.Sprintf("%d", time.Now().Unix())
	temp := utils.MD5Encode(str)
	utils.DB.Model(&user).Where("id = ?", user.ID).Update("identity", temp)
	return user
}

// 通过名字查找用户
func FindUserByName(name string) UserBasic {
	user := UserBasic{}
	utils.DB.Where("name = ? ", name).First(&user)
	return user
}

func FindUserByPhone(phone string) *gorm.DB {
	user := UserBasic{}
	return utils.DB.Where("phone = ?", phone).First(&user)
}

func FindUserByEmail(email string) *gorm.DB {
	user := UserBasic{}
	return utils.DB.Where("email = ?", email).First(&user)
}

// 新增用户
func CreateUser(user UserBasic) *gorm.DB {

	return utils.DB.Create(&user)
}

// 删除用户
func DeleteUser(user UserBasic) *gorm.DB {
	return utils.DB.Delete(&user)
}

// 修改用户
func UpdateUser(user UserBasic) *gorm.DB {
	return utils.DB.Model(&user).Updates(UserBasic{
		Name:     user.Name,
		PassWord: user.PassWord,
		Phone:    user.Phone,
		Email:    user.Email,
	})
}

// 查找某个用户
func FindByID(id uint) UserBasic {
	user := UserBasic{}
	utils.DB.Where("id = ?", id).First(&user)
	return user
}
