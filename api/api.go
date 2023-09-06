package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"mj-wechat-bot/errorhandler"
	"mj-wechat-bot/types/request"
	"mj-wechat-bot/types/response"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type API struct {
	ApiUrl   string `yaml:"api_url"`
	Apikey   string `yaml:"api_key"`
	CheckUrl string `yaml:"check_url"`
}

var config API
var (
	createUrl     string
	taskUrl       string
	taskUpdateUrl string
	describeUrl   string
	blendUrl      string
	parseUrl      string
)

func init() {
	// 注册异常处理函数
	defer errorhandler.HandlePanic()
	// Read configuration file.
	data, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		panic(fmt.Sprintf("读取配置文件失败: %v", err))
	}

	// Unmarshal configuration.

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		panic(fmt.Sprintf("解析配置文件失败: %v", err))
	}
	createUrl = config.ApiUrl + "/submit/imagine"
	taskUrl = config.ApiUrl + "/task/%s/fetch"
	taskUpdateUrl = config.ApiUrl + "/submit/change"
	describeUrl = config.ApiUrl + "/submit/describe"
	blendUrl = config.ApiUrl + "/submit/blend"
	parseUrl = config.ApiUrl + "/sucai/parse"
}

type Response struct {
	Code int                    `json:"code"`
	Data map[string]interface{} `json:"data"`
	Msg  string                 `json:"msg"`
	Time int                    `json:"time"`
}

func CreateMessage(text string) (bool, string) {
	imagineReq := request.ImagineReq{
		Base64Array: nil,
		NotifyHook:  nil,
		Prompt:      text,
		State:       nil,
	}

	body, err := DoPost(createUrl, imagineReq)
	if err != nil {
		fmt.Println(err)
		return false, ""
	}

	var response response.MidjResp
	if err := json.Unmarshal(body, &response); err != nil {
		fmt.Println(err)
		return false, ""
	}
	status := []int{1, 21, 22} //状态码: 1(提交成功), 21(已存在), 22(排队中), other(错误)
	// 将切片中的值映射到一个 map 中
	numberMap := make(map[int]bool)
	for _, num := range status {
		numberMap[num] = true
	}
	if !numberMap[response.Code] {
		fmt.Println(response.Description)
		return false, ""
	}
	return true, response.Result
}

//查询任务状态
func QueryTaskStatus(taskID string) (bool, *response.TaskResp) {
	reqUrl, err := url.Parse(fmt.Sprintf(taskUrl, taskID))
	if err != nil {
		fmt.Println(err)
		return false, nil
	}
	body, err := DoGet(reqUrl)
	if err != nil {
		fmt.Println(err)
		return false, nil
	}
	if err != nil {
		fmt.Println(err)
		return false, nil
	}
	log.Printf("任务【%s】返回结果 -> %s", taskID, body)
	var response response.TaskResp
	if err := json.Unmarshal([]byte(body), &response); err != nil {
		fmt.Println(err)
		return false, nil
	}
	return true, &response
}

func TaskUpdate(taskId string, action string) (bool, string) {
	actionParam := ""
	index := 1
	if strings.Contains(action, "u") {
		actionParam = "UPSCALE"
		index, _ = strconv.Atoi(string(action[1]))
	} else if strings.Contains(action, "v") {
		actionParam = "VARIATION"
		index, _ = strconv.Atoi(string(action[1]))
	}
	changeReq := request.ChangeReq{
		Action:     actionParam,
		Index:      index,
		NotifyHook: "",
		State:      "",
		TaskId:     taskId,
	}
	body, err := DoPost(taskUpdateUrl, changeReq)
	if err != nil {
		fmt.Println(err)
		return false, ""
	}

	var response response.MidjResp
	if err := json.Unmarshal(body, &response); err != nil {
		fmt.Println(err)
		return false, ""
	}
	status := []int{1, 21, 22} //状态码: 1(提交成功), 21(已存在), 22(排队中), other(错误)
	// 将切片中的值映射到一个 map 中
	numberMap := make(map[int]bool)
	for _, num := range status {
		numberMap[num] = true
	}
	if !numberMap[response.Code] {
		fmt.Println(response.Description)
		return false, ""
	}
	return true, response.Result
}

func Describe(base64 string) (bool, string) {
	describeReq := request.DescribeReq{
		Base64:     base64,
		NotifyHook: "",
		State:      "",
	}
	body, err := DoPost(describeUrl, describeReq)
	if err != nil {
		fmt.Println(err)
		return false, ""
	}

	var response response.MidjResp
	if err := json.Unmarshal(body, &response); err != nil {
		fmt.Println(err)
		return false, ""
	}
	status := []int{1, 21, 22} //状态码: 1(提交成功), 21(已存在), 22(排队中), other(错误)
	// 将切片中的值映射到一个 map 中
	numberMap := make(map[int]bool)
	for _, num := range status {
		numberMap[num] = true
	}
	if !numberMap[response.Code] {
		fmt.Println(response.Description)
		return false, ""
	}
	return true, response.Result
}

func Blend(base64 []string) (bool, string) {
	blendReq := request.BlendReq{
		Base64Array: base64,
		Dimensions:  "SQUARE",
		NotifyHook:  "",
		State:       "",
	}
	body, err := DoPost(blendUrl, blendReq)
	if err != nil {
		fmt.Println(err)
		return false, ""
	}

	var response response.MidjResp
	if err := json.Unmarshal(body, &response); err != nil {
		fmt.Println(err)
		return false, ""
	}
	status := []int{1, 21, 22} //状态码: 1(提交成功), 21(已存在), 22(排队中), other(错误)
	// 将切片中的值映射到一个 map 中
	numberMap := make(map[int]bool)
	for _, num := range status {
		numberMap[num] = true
	}
	if !numberMap[response.Code] {
		fmt.Println(response.Description)
		return false, ""
	}
	return true, response.Result
}

func Parse(link string) (bool, string) {
	reqUrl, err := url.Parse(parseUrl + "?url=" + link)
	if err != nil {
		fmt.Println(err)
		return false, ""
	}
	body, err := DoGet(reqUrl)
	if err != nil {
		fmt.Println(err)
		return false, ""
	}
	if err != nil {
		fmt.Println(err)
		return false, ""
	}

	var response response.MidjResp
	if err = json.Unmarshal([]byte(body), &response); err != nil {
		log.Printf(err.Error())
		return false, ""
	}
	if response.Code != 200 {
		return false, response.Description
	}

	return true, response.Result
}

func DoGet(reqUrl *url.URL) (string, error) {
	// 构建 HTTP GET 请求
	req, err := http.NewRequest("GET", reqUrl.String(), nil)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	// 创建一个 HTTP 客户端
	client := &http.Client{
		Timeout: 120 * time.Second,
	}
	// 添加请求头
	req.Header.Add("apikey", config.Apikey)
	// 发送 HTTP GET 请求
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return string(body), nil
}

func DoPost(url string, param interface{}) (body []byte, err error) {
	var (
		resp *http.Response
	)
	jsonData, err := json.Marshal(param)
	if err != nil {
		fmt.Println("JSON编码失败:", err)
		return
	}
	// 创建HTTP请求
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("创建请求失败:", err)
		return
	}
	// 设置请求头，指定JSON内容类型
	req.Header.Set("Content-Type", "application/json")
	// 发送请求
	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		fmt.Println("发送请求失败:", err)
		return
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	return
}
