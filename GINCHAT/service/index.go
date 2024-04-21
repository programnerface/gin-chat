package service

//业务逻辑层
import (
	"fmt"
	"ginchat/models"
	"github.com/gin-gonic/gin"
	"html/template"
	"strconv"
)

// GetIndex
// @Tags 首页
// @Success 200 {string} success
// @Router /index [get]
func GetIndex(c *gin.Context) {
	ind, err := template.ParseFiles("index.html", "views/chat/head.html")
	if err != nil {
		panic(err)
		fmt.Println(err)
	}
	ind.Execute(c.Writer, "index")
	//c.JSON(http.StatusOK, gin.H{
	//	"message": "success",
	//})
}

// 注册页面
func ToRegister(c *gin.Context) {
	ind, err := template.ParseFiles("views/user/register.html")
	if err != nil {
		panic(err)
		fmt.Println(err)
	}
	ind.Execute(c.Writer, "register")
	//c.JSON(http.StatusOK, gin.H{
	//	"message": "success",
	//})
}

// 聊天
func ToChat(c *gin.Context) {
	ind, err := template.ParseFiles("views/chat/index.html",
		"views/chat/head.html",
		"views/chat/tabmenu.html",
		"views/chat/concat.html",
		"views/chat/group.html",
		"views/chat/profile.html",
		"views/chat/main.html",
		"views/chat/createcom.html",
		"views/chat/userinfo.html",
		"views/chat/foot.html",
	)
	if err != nil {
		panic(err)
		fmt.Println(err)
	}
	userId, _ := strconv.Atoi(c.Query("userId"))
	token := c.Query("token")
	user := models.UserBasic{}
	user.ID = uint(userId)
	user.Identity = token
	fmt.Println("ToChat>>>>>>>>", user)
	ind.Execute(c.Writer, user)
	//c.JSON(http.StatusOK, gin.H{
	//	"message": "success",
	//})
}

// 调用models中的Chat方法初始化websocket
func Chat(c *gin.Context) {
	models.Chat(c.Writer, c.Request)
}
