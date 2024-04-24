package service

//业务逻辑层
import (
	"fmt"
	"ginchat/models"
	"ginchat/utils"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"math/rand"

	"net/http"
	"strconv"
	"time"
)

// GetUserList
// @Summary 新增用户
// @Tags 用户列表
// @Success 200 {string} json{"code","message"}
// @Router /user/getUserList [get]
func GetUserList(c *gin.Context) {
	//返回值
	data := make([]*models.UserBasic, 10)
	//从models 层拿到这个方法
	data = models.GetUserList()
	c.JSON(http.StatusOK, gin.H{
		"code":    0, //0成功，-1失败
		"message": "用户列表",
		"data":    data,
	})
}

// CreateUser
// @Summary 新增用户
// @Tags 用户模块
// @param name query string false "用户名"
// @param password query string false "密码"
// @param repassword query string false "确认密码"
// @Success 200 {string} json{"code","message"}
// @Router /user/createUser [get]
func CreateUser(c *gin.Context) {
	user := models.UserBasic{}
	//user.Name = c.Query("name")
	//password := c.Query("password")
	//repassword := c.Query("repassword")
	user.Name = c.Request.FormValue("name")
	password := c.Request.FormValue("password")
	repassword := c.Request.FormValue("Identity")
	fmt.Println(user.Name, ">>>>>>>>>", password, repassword)
	salt := fmt.Sprintf("%06d", rand.Int31())

	data := models.FindUserByName(user.Name)
	//用户名密码不能为空
	if user.Name == "" || password == "" || repassword == "" {
		c.JSON(200, gin.H{
			"code":    -1, //0成功，-1失败
			"message": "用户名或密码不能为空",
			"data":    user,
		})
		return
	}

	if data.Name != "" {
		c.JSON(200, gin.H{
			"code":    -1, //0成功，-1失败
			"message": "用户名已注册",
			"data":    user,
		})
		return
	}
	if password != repassword {
		c.JSON(200, gin.H{
			"code":    -1, //0成功，-1失败
			"message": "两次密码不一致",
			"data":    user,
		})
		return
	}
	//user.PassWord = password
	user.PassWord = utils.MakePassword(password, salt)
	user.Salt = salt //表更新看字段
	fmt.Println(user.PassWord)
	user.LoginTime = time.Now()
	user.LoginOutTime = time.Now()
	models.CreateUser(user)
	c.JSON(200, gin.H{
		"code":    0, //0成功，-1失败
		"message": "新增用户成功！",
		"data":    user,
	})
}

// FindUserByNameAndPwd
// @Summary 登录
// @Tags 用户列表
// @param name query string false "用户名"
// @param password query string false "密码"
// @Success 200 {string} json{"code","message"}
// @Router /user/findUserByNameAndPwd [post]
func FindUserByNameAndPwd(c *gin.Context) {
	data := models.UserBasic{}

	//name := c.Query("name")
	//password := c.Query("password")
	//获取表单中的信息
	name := c.Request.FormValue("name")
	password := c.Request.FormValue("password")
	fmt.Println(name, password)
	user := models.FindUserByName(name)
	if user.Name == "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    -1, //0成功，-1失败
			"message": "该用户不存在",
			"data":    data,
		})
		return
	}
	fmt.Println(user)
	flag := utils.ValidPassword(password, user.Salt, user.PassWord)
	if !flag {
		c.JSON(http.StatusOK, gin.H{
			"code":    -1, //0成功，-1失败
			"message": "密码不正确",
			"data":    data,
		})
		return
	}
	pwd := utils.MakePassword(password, user.Salt)
	data = models.FindUserByNameAndPwd(name, pwd)
	c.JSON(http.StatusOK, gin.H{
		"code":    0, //0成功，-1失败
		"message": "登录成功！",
		"data":    data,
	})
}

// DeleteUser
// @Summary 删除用户
// @Tags 用户模块
// @param id query string false "id"
// @Success 200 {string} json{"code","message"}
// @Router /user/deleteUser [get]
func DeleteUser(c *gin.Context) {
	user := models.UserBasic{}
	id, _ := strconv.Atoi(c.Query("id"))
	user.ID = uint(id)
	models.DeleteUser(user)
	c.JSON(200, gin.H{
		"code":    0, //0成功，-1失败
		"message": "删除用户成功！",
		"data":    user,
	})
}

// UpdateUser
// @Summary 修改用户
// @Tags 用户模块
// @param id formData string false "id"
// @param name formData string false "用户名"
// @param password formData string false "密码"
// @param phone formData string false "phone"
// @param email formData string false "email"
// @Success 200 {string} json{"code","message"}
// @Router /user/updateUser [post]
func UpdateUser(c *gin.Context) {
	user := models.UserBasic{}
	id, _ := strconv.Atoi(c.PostForm("id"))
	user.ID = uint(id)
	user.Name = c.PostForm("name")
	user.PassWord = c.PostForm("password")
	user.Phone = c.PostForm("phone")
	user.Email = c.PostForm("email")
	fmt.Println("UpdateUser:", user)

	_, err := govalidator.ValidateStruct(user)
	if err != nil {
		fmt.Println(err)
		c.JSON(200, gin.H{
			"code":    -1, //0成功，-1失败
			"message": "修改参数不正确！",
			"data":    user,
		})
	} else {
		models.UpdateUser(user)
		c.JSON(200, gin.H{
			"code":    -1, //0成功，-1失败
			"message": "修改参数不正确！",
			"data":    user,
		})
	}

}

// 防止跨域站点伪造请求
// 将http连接升级成websocket连接
// websocket.Upgrader 是 WebSocket 包中的一个结构体，用于配置 WebSocket 连接的升级过程
var upGrade = websocket.Upgrader{
	//初始化 websocket.Upgrader 结构体的字段
	//CheckOrigin 用于检查请求的来源是否合法
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// 消息发送
func SendMsg(c *gin.Context) {
	//用于将 HTTP 连接升级为 WebSocket 连接
	//c.Writer HTTP 响应写入器  c.Request HTTP请求对象 nil可选协议头
	ws, err := upGrade.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		//错误检查
		fmt.Println(err)
		return
	}
	//使用 defer 关键字延迟执行一个匿名函数
	defer func(ws *websocket.Conn) {
		//在函数执行完毕后关闭 WebSocket 连接
		err = ws.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(ws)
	//发送消息 用于处理 WebSocket 连接的消息
	MsgHandler(ws, c)
}

// 消息发送 处理 WebSocket 连接接收到的消息，并将接收到的消息发送给客户端
// ws（WebSocket 连接对象）和 c（HTTP 请求上下文）
func MsgHandler(ws *websocket.Conn, c *gin.Context) {
	//持续监听从 Redis 订阅通道接收到的消息，并将其发送给客户端
	for {
		//从 Redis 订阅指定的通道（PublishKey）接收消息
		msg, err := utils.Subscribe(c, utils.PublishKey)
		if err != nil {
			fmt.Println(err)
		}
		//打印接收到的消息内容
		fmt.Println("发送消息:", msg)
		//获取当前时间 并格式化为字符串
		tm := time.Now().Format("2006-01-02 15:04:05")
		//将消息内容和时间消息拼接成完整的字符串
		m := fmt.Sprintf("[ws][%s]:%s", tm, msg) //发送的消息
		// 将格式化后的消息字符串通过 WebSocket 连接发送给客户端
		//第一个参数是消息类型这里使用 1 表示文本消息，第二个参数是消息的内容，需要将消息字符串转换为字节数组。
		err = ws.WriteMessage(1, []byte(m))
		if err != nil {
			fmt.Println(err)
		}
	}
}

func SendUserMsg(c *gin.Context) {
	models.Chat(c.Writer, c.Request)
}

func SearchFriends(c *gin.Context) {
	id, _ := strconv.Atoi(c.Request.FormValue("userId"))
	users := models.SearchFriend(uint(id))
	//c.JSON(200, gin.H{
	//	"code":    0, //0成功，-1失败
	//	"message": "查询好友列表成功！",
	//	"data":    users,
	//})
	//返回请求的类型是封装类型
	utils.RespOKList(c.Writer, users, len(users))
}

func AddFriend(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Request.FormValue("userId"))
	targetId, _ := strconv.Atoi(c.Request.FormValue("targetId"))
	code := models.AddFriend(uint(userId), uint(targetId))
	//c.JSON(200, gin.H{
	//	"code":    0, //0成功，-1失败
	//	"message": "查询好友列表成功！",
	//	"data":    users,
	//})

	if code == 0 {
		utils.RespOK(c.Writer, code, "添加好友成功")
	} else {
		utils.RespFail(c.Writer, "添加失败")
	}

}
