package models

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"gopkg.in/fatih/set.v0"
	"gorm.io/gorm"
	"net"
	"net/http"
	"strconv"
	"sync"
)

// 实体类字段名首字母要大写，不然无法写入到数据库中
// 消息
type Message struct {
	gorm.Model
	Formid   int64  ////从哪里发过来的 发送者
	TargetId int64  //接收者
	Type     int    //发送类型 1私聊 2群聊   3广播
	Media    int    //消息类型 1文字 2表情包 3图片 4音频
	Content  string //消息内容
	Pic      string //图片
	Url      string //URL
	Desc     string //描述
	Amount   int    //其他数字统计 (发送的频率，文件大小)
}

func (table *Message) TableName() string {
	return "message"
}

//引入set包

type Node struct {
	Conn      *websocket.Conn //websocket连接
	DataQueue chan []byte     //广告
	GroupSets set.Interface   //集合
}

// 映射关系 客户端
var clientMaps map[int64]*Node = make(map[int64]*Node)

// 读写锁
var rwLocker sync.RWMutex

// 需要:发送者ID,接收者ID，消息类型，发送的内容， 发送的类型
func Chat(writer http.ResponseWriter, request *http.Request) {
	//1.获取参数并校验token 等合法性
	//token := query.Get("token")
	query := request.URL.Query()
	Id := query.Get("userId")
	userId, _ := strconv.ParseInt(Id, 10, 64)
	//msgTpye := query.Get("type")
	//targetId := query.Get("targetId")
	//context := query.Get("context")
	isvalida := true //checkToke()
	//处理连接
	conn, err := (&websocket.Upgrader{
		//做校验
		CheckOrigin: func(r *http.Request) bool {
			return isvalida
		},
	}).Upgrade(writer, request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	//2.获取conn(连接)
	node := &Node{
		Conn:      conn,
		DataQueue: make(chan []byte, 50), //发送的是这个
		GroupSets: set.New(set.ThreadSafe),
	}
	//3.用户关系
	//4.userid 和node 绑定并加锁
	rwLocker.Lock() //加锁
	clientMaps[userId] = node
	rwLocker.Unlock() //解锁
	//5.完成发送逻辑
	go sendProc(node)
	//6.完成接受逻辑
	go recvProc(node)
	fmt.Println("sendMsg>>>>>>>>userId:", userId)
	sendMsg(userId, []byte("欢迎进入聊天系统"))
}

// 发送消息协程
func sendProc(node *Node) {
	for {
		select {
		case data := <-node.DataQueue:
			fmt.Println("[ws]sendMsg>>>>>>>>msg:", string(data))
			err := node.Conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}

func recvProc(node *Node) {
	for {
		_, data, err := node.Conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}
		dispatch(data)
		broadMsg(data) // todo 广播消息到局域网
		fmt.Println("[ws] recvProc <<<<< ", string(data))
	}
}

// 广道
var udpsendchan chan []byte = make(chan []byte, 1024)

// 广播消息到局域网
func broadMsg(data []byte) {
	udpsendchan <- data
}

// 协程初始化
func init() {
	go udpSendProc()
	go udprecvProc()
	fmt.Println("init goroutine ")
}

// 完成udp的数据发送协程
func udpSendProc() {
	con, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.IPv4(192, 168, 1, 1),
		Port: 3000,
	})
	defer con.Close()
	if err != nil {
		fmt.Println(err)
	}
	for {
		select {
		case data := <-udpsendchan:
			fmt.Println("udpSendProc >>>data:", string(data))
			_, err := con.Write(data)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}

func udprecvProc() {
	con, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4zero, //设置为0，所有人都能接收到
		Port: 3000,
	})
	if err != nil {
		fmt.Println(err)
	}
	defer con.Close()
	for {
		var buf [512]byte
		n, err := con.Read(buf[0:]) //读取buf 从0到最后
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("udprecvProc >>>data:", string(buf[0:n]))
		dispatch(buf[0:n])
	}
}

// 后端调度逻辑处理
func dispatch(data []byte) {
	msg := Message{}                  //初始化msg
	err := json.Unmarshal(data, &msg) //json转化并处理异常
	if err != nil {
		fmt.Println(err)
		return
	}
	//根据消息类型做对于的处理
	switch msg.Type {
	case 1: //私信
		fmt.Println("dispatch  >>data:", string(data))
		sendMsg(msg.TargetId, data)
		//case 2: //群发
		//	sendGroupMsg()
		//case 3: //广播
		//	sendAllMsg()
		//case 4:
		//	//
	}
}

func sendMsg(userId int64, msg []byte) {
	fmt.Println("sendMsg >>>userID: ", userId, " msg:", string(msg))
	rwLocker.RLocker()
	node, ok := clientMaps[userId] //拿到user 拿到node
	rwLocker.RLocker()             //解锁
	if ok {
		node.DataQueue <- msg //往node里面丢消息放在队列里面
	}
}
