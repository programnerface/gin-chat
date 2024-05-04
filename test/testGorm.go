package main

import (
	"ginchat/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

//type Product struct {
//	gorm.Model
//	Code  string
//	Price uint
//}

func main() {
	//连接mysql数据库
	db, err := gorm.Open(mysql.Open("root:123456@tcp(127.0.0.1:3306)/ginchat?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{})
	//db, err := gorm.Open(mysql.Open(""), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema 实体类
	//AutoMigrate() 如果结构体UserBasic在数据库不存在表格，则会去创建它
	db.AutoMigrate(&models.Community{})
	//db.AutoMigrate(&models.UserBasic{})
	//db.AutoMigrate(&models.Message{})
	//db.AutoMigrate(&models.GroupBasic{})
	//db.AutoMigrate(&models.Contact{})

	// Create 创建
	//user := &models.UserBasic{}
	//user.Name = "face"
	//user.HeartbeatTime = time.Now()
	//db.Create(user)
	//create()方法向数据库插入一条数据
	//&Product{Code: "D42", Price: 100} 这个是结构体的实例(分别给字段设置了值)
	//db.Create(&Product{Code: "D42", Price: 100})
	// Read 查询 读取
	//fmt.Println(db.First(user, 1))
	//var product Product  //实例化一个product对象
	//db.First(&product, 1)                  //  根据整数组件查找产品
	//db.First(&product, "code = ?", "D42") // 根据条件找到code=D42的数据

	//db.Model(user).Update("PassWord", "12346")
	// Update - update 将产品的价格更新到200元
	//db.Model(&product).Update("Price", 200)
	// Update - update multiple fields
	//更新多个字段的值
	//db.Model(&product).Updates(Product{Price: 200, Code: "F42"}) // non-zero fields
	//使用map参数来更新更多字段
	//db.Model(&product).Updates(map[string]interface{}{"Price": 200, "Code": "F42"})

	// Delete - delete product
	//删除具有指定整数主键的产品
	//db.Delete(&product, 1)
	//db.Model(user).Delete("Name")
}
