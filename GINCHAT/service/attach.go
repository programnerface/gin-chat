package service

import (
	"fmt"
	"ginchat/utils"
	"github.com/gin-gonic/gin"
	"io"
	"math/rand"
	"os"
	"strings"
	"time"
)

// 图片发送
func Upload(c *gin.Context) {
	w := c.Writer
	req := c.Request
	srcFile, head, err := req.FormFile("file") //此处要与前端页面一致
	if err != nil {
		utils.RespFail(w, err.Error())
		fmt.Println(err)
	}
	suffix := ".png"                     //默认图片格式
	ofilName := head.Filename            //拿到当前文件名称
	temp := strings.Split(ofilName, ".") //通过.分割
	if len(temp) > 1 {
		suffix = "." + temp[len(temp)-1] //拿到后缀名
	}
	fileName := fmt.Sprintf("%d%04d%s", time.Now().Unix(), rand.Int31(), suffix)
	dstFile, err := os.Create("./asset/upload/" + fileName)
	if err != nil {
		utils.RespFail(w, err.Error())
	}
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {

	}
	url := "./asset/upload/" + fileName
	utils.RespOK(w, url, "发送图片成功")
}
