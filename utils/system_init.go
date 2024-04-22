package utils

//放一些初始化的信息
//viper包 用于读取yaml文件
import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

var (
	DB  *gorm.DB
	Red *redis.Client
)

func InitConfig() {
	viper.SetConfigName("app")
	viper.AddConfigPath("config")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
	}
	//如果没有报错就打印配置信息
	//fmt.Println("config app:", viper.Get("app"))
	//fmt.Println("config app:", viper.Get("mysql"))
	fmt.Println("config app inited ....")
}

func InitMySQL() {
	//自定义日志模板 打印sql语句
	//创建一个新的日志记录器
	//logger.New()第一个参数是输出的位置，第二个参数是日志的配置
	newLogger := logger.New(
		//在New函数中第一个参数是os.Stdout的接口实例表示日志输出到标准输出流(控制台)
		//第二个参数是日志的分隔符表示换行
		//第三个参数 log.LstdFlags 是日志的格式标志表示日志中显示时间日期和文件信息
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second, //慢sql阈值
			LogLevel:      logger.Info, //级别
			Colorful:      true,        //彩色
		},
	)
	//Logger: newLogger使用刚刚创建好的日志记录器
	DB, _ = gorm.Open(mysql.Open(viper.GetString("mysql.dns")), &gorm.Config{Logger: newLogger})
	fmt.Println("MySQL inited....")
	//user := models.UserBasic{}
	//DB.Find(&user)
	//fmt.Println(user)
}

// 初始化redis
func InitRedis() {
	Red = redis.NewClient(&redis.Options{
		Addr:         viper.GetString("redis.addr"),
		Password:     viper.GetString("redis.password"),
		DB:           viper.GetInt("redis.DB"),
		PoolSize:     viper.GetInt("redis.poolSize"),
		MinIdleConns: viper.GetInt("redis.minIdleConn"),
	})
	//pong, err := Red.Ping().Result()
	//if err != nil {
	//	//如果连接失败 打印err
	//	fmt.Println("init redis err.....", err)
	//} else {
	//	//连接成功 打印pong
	//	fmt.Println("Redis inited........", pong)
	//}
}

const (
	PublishKey = "websocket"
)

// Publish 发布消息到Redis
func Publish(ctx context.Context, channel string, msg string) error {
	var err error
	fmt.Println("Publish........", msg)
	//调用了Publish方法，将消息发送到channel，ctx表示上下文
	err = Red.Publish(ctx, channel, msg).Err()
	if err != nil {
		fmt.Println(err)
	}
	return err
}

// Subscribe 订阅Redis消息
func Subscribe(ctx context.Context, channel string) (string, error) {
	//订阅指定的redis通道，返回*redis.PubSub对象
	sub := Red.Subscribe(ctx, channel)
	fmt.Println("Subscribe.........", ctx)
	ctxStr := fmt.Sprintf("%+v", ctx)
	fmt.Println("Subscribe Context:", ctxStr)
	//用于从订阅的通道接收消息，
	msg, err := sub.ReceiveMessage(ctx)
	// 检查上一步操作是否发生了错误
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	//打印接收到的消息内容
	fmt.Println("Subscribe.........", msg.Payload)
	fmt.Println("2")
	fmt.Printf("%v\n", msg.Payload)
	//返回接收到的消息内容和可能发生的错误信息
	return msg.Payload, err
}
