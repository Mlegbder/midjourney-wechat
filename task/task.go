package task

import (
	"bytes"
	"fmt"
	"github.com/eatmoreapple/openwechat"
	"golang.org/x/image/webp"
	"image/png"
	"io"
	"log"
	"mj-wechat-bot/api"
	"mj-wechat-bot/bot"
	"mj-wechat-bot/conn"
	"mj-wechat-bot/consts"
	"mj-wechat-bot/model"
	"mj-wechat-bot/replay"
	"mj-wechat-bot/types/response"
	"mj-wechat-bot/utils"
	"path"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

var (
	taskIds = sync.Map{}
	wg      = sync.WaitGroup{}
)

type ImageMsg struct {
	taskId       string
	fromUserName string
	reader       io.Reader
	url          string
}

func RunTask() {
	go Looper()
	go ImageSender()
}

var (
	count   = int64(0)
	msgChan = make(chan ImageMsg, 100)
	test    = sync.RWMutex{}
)

// AddTask 添加任务
func AddTask(msg *openwechat.Message, taskId string) {
	log.Printf("添加任务:%s", taskId)
	atomic.AddInt64(&count, 1)
	taskIds.Store(taskId, msg)
}

func ImageSender() {
	for {
		select {
		case imageMsg := <-msgChan:
			log.Printf("收到发送图片任务，开始发送图片")
			sendImage(imageMsg)
			//log.Printf("%v", msg)
			// 发送图片消息
			time.Sleep(5 * time.Second)
		}
	}
}

// Looper 任务循环
func Looper() {
	log.Printf("开始启动任务循环")
	for {
		log.Printf("任务数量:%d", count)
		taskIds.Range(func(taskId, _ any) bool {
			wg.Add(1)
			// 查询任务状态
			go QueryTaskStatus(taskId.(string))
			return true
		})
		wg.Wait()
		time.Sleep(5 * time.Second)
	}
}

// QueryTaskStatus 查询任务状态并发送图片消息
func QueryTaskStatus(taskId string) {

	// 查询任务状态
	ok, data := api.QueryTaskStatus(taskId)
	value, ok1 := taskIds.Load(taskId)
	if !ok1 {
		wg.Done()
		return
	}
	userMsg := value.(*openwechat.Message)
	fromUserName := userMsg.FromUserName

	name, err := utils.GetUserName(userMsg)
	if err == nil {

	}
	info := replay.Info{
		NickName: name,
		TaskId:   taskId,
	}
	if ok {
		// 判断是否完成
		switch data.Status {
		case "finish":
		case "SUCCESS":
			go func() {
				var url string
				var reader io.Reader
				if data.Action != "DESCRIBE" {
					url = data.ImageUrl
					info.Url = url
					ok := false

					failCount := 0
					for !ok {
						//转码失败3次
						if failCount > 3 {
							//发送失败消息
							failTask(taskId, fromUserName, info.GenrateMessage(replay.TaskTransImgErrMsg))
							return
						}
						reader, ok = webp2png(url)
						failCount++
						time.Sleep(1 * time.Second)
					}
					addImageMsgChan(ImageMsg{
						taskId:       taskId,
						fromUserName: fromUserName,
						reader:       reader,
						url:          url,
					})
				} else {
					info.Describe = data.Prompt
					atomic.AddInt64(&count, -1)
				}
				typeName, exist := userMsg.Get("type")
				if exist {
					if typeName.(string) == "main" {
						userMsg.ReplyText(info.GenrateMessage(replay.TaskMainFinishMsg))
					} else if strings.HasPrefix(typeName.(string), "V") {
						userMsg.ReplyText(info.GenrateMessage(replay.TaskSubVFinishMsg))
					} else if strings.HasPrefix(typeName.(string), "U") {
						userMsg.ReplyText(info.GenrateMessage(replay.TaskSubUFinishMsg))
					} else if strings.HasPrefix(typeName.(string), "DESCRIBE") {
						userMsg.ReplyText(info.GenrateMessage(replay.TaskSubDFinishMsg))
					} else if strings.HasPrefix(typeName.(string), "BLEND") {
						userMsg.ReplyText(info.GenrateMessage(replay.TaskSubBFinishMsg))
					}
				}

			}()
			// 删除任务
			taskIds.Delete(taskId)
			break
		case "IN_PROGRESS":
			// 任务未完成
			break
		case "invalid params":
			// 任务参数错误
			failTask(taskId, fromUserName, info.GenrateMessage(replay.TaskParamsErrMsg))
			break
		case "invalid link":
			// 任务参数错误
			failTask(taskId, fromUserName, info.GenrateMessage(replay.TaskLinkErrMsg))
			break
		case "banned":
			// 任务被封禁
			// 任务参数错误
			failTask(taskId, fromUserName, info.GenrateMessage(replay.TaskBannedErrMsg))
			break
		case "FAILURE":
			// 任务被封禁
			// 任务参数错误
			failTask(taskId, fromUserName, info.GenrateMessage(replay.TaskErrMsg)+"💰余额已自动退回 \n❎错误原因: "+data.FailReason)
			rollbackBalance(data, userMsg)
			break
		}

	} else {
		failTask(taskId, fromUserName, info.GenrateMessage(replay.TaskErrMsg1))
	}
	wg.Done()
}

func rollbackBalance(data *response.TaskResp, userMsg *openwechat.Message) {
	var (
		db = conn.DB
	)
	userId, err := utils.GetUserId(userMsg)
	if err != nil {
		log.Println(err)
		return
	}
	var wUser model.WUser
	err = db.Model(&model.WUser{}).Where("wechat_id = ?", userId).Find(&wUser).Error
	if err != nil {
		log.Println(err)
		return
	}
	switch data.Action {
	case "IMAGINE":
		wUser.Balance = wUser.Balance + consts.IMAGINE
		break
	case "UPSCALE":
		wUser.Balance = wUser.Balance + consts.UPSCALE
		break
	case "VARIATION":
		wUser.Balance = wUser.Balance + consts.VARIATION
		break
	case "REROLL":
		wUser.Balance = wUser.Balance + consts.REROLL
		break
	case "DESCRIBE":
		wUser.Balance = wUser.Balance + consts.DESCRIBE
		break
	case "BLEND":
		wUser.Balance = wUser.Balance + consts.BLEND
		break
	}
	err = db.Save(&wUser).Error
	if err != nil {
		log.Println(err)
		return
	}
}

func addImageMsgChan(msg ImageMsg) {
	msgChan <- msg
}
func failTask(taskId string, fromUserName string, msg string) {
	req := bot.Bot.Storage.Request
	info := bot.Bot.Storage.LoginInfo
	log.Printf("req:%v,info:%v,bot.CurrentUser:%s,fromUserName:%s", req, info, bot.CurrentUser, fromUserName)
	// 获取登陆的用户
	CurrentUser := bot.CurrentUser
	_, err := bot.Bot.Caller.WebWxSendMsg(&openwechat.SendMessage{
		FromUserName: CurrentUser.UserName,
		ToUserName:   fromUserName,
		Content:      msg,
	}, info, req)
	if err != nil {
		fmt.Println(err)
		return
	}
	log.Printf("任务失败(%s),删除任务", taskId)
	//删除任务
	taskIds.Delete(taskId)
	atomic.AddInt64(&count, -1)
}
func webp2png(url string) (io.Reader, bool) {
	// 发送图片消息
	ok, reader := utils.GetImageUrlData(url)
	if !ok {
		return nil, false
	}
	// 通过 path.Ext 函数解析链接地址中的后缀名
	ext := path.Ext(url)
	// 根据后缀名判断是否是 webp 格式的图片
	if ext == ".webp" {
		image, err := webp.Decode(reader)
		if err != nil {
			fmt.Println(err)
			return nil, false
		}
		// 创建一个 PNG 格式的 io.Reader
		var pngReader io.Reader
		buf := new(bytes.Buffer)
		if err := png.Encode(buf, image); err != nil {
			fmt.Printf("pngReader: %v", err)
			return nil, false
		}
		pngReader = bytes.NewReader(buf.Bytes())
		reader = pngReader

	}
	return reader, ok
}

// 发送图片消息
func sendImage(imageMsg ImageMsg) {

	// 发送图片消息
	req := bot.Bot.Storage.Request
	info := bot.Bot.Storage.LoginInfo
	//log.Printf("req:%v,info:%v,bot.CurrentUser:%s,fromUserName:%s\n", req, info, bot.CurrentUser, fromUserName)
	// 获取登陆的用户
	CurrentUser, err := bot.Bot.GetCurrentUser()
	if err != nil {
		fmt.Printf("获取当前登陆用户失败:%s", err)
		addImageMsgChan(imageMsg)
		return
	}
	_, err = bot.Bot.Caller.WebWxSendImageMsg(imageMsg.reader, req, info, CurrentUser.UserName, imageMsg.fromUserName)
	if err != nil {
		fmt.Println(err)
		addImageMsgChan(imageMsg)
		return
	}
	//完成任务
	log.Printf("发送图片完成,删除任务:%s", imageMsg.taskId)

	atomic.AddInt64(&count, -1)
}
