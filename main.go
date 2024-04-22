package main

import (
	"ginchat/router"
	"ginchat/utils"
)

func main() {
	//将mysql的配置信息写在app.yml，在使用InitConfig()方法
	//去获取到配置信息然后在main.go进行调用--这个就是封装
	utils.InitConfig()
	utils.InitMySQL()
	utils.InitRedis()
	r := router.Router()
	r.Run(":8081")
	//r := gin.Default()
	//r.GET("/ping", func(c *gin.Context) {

	//	c.JSON(http.StatusOK, gin.H{
	//		"message": "pong",
	//	})
	//})
	//r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
