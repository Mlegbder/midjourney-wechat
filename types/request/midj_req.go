package request

type (
	//ImagineReq 画图参数
	ImagineReq struct {
		Base64Array *[]interface{} `json:"base64Array"` //	垫图base64数组
		NotifyHook  *string        `json:"notifyHook"`  // 	回调地址, 为空时使用全局notifyHook
		Prompt      string         `json:"prompt"`      //	提示词,示例值(Cat)
		State       *string        `json:"state"`       //	自定义参数
	}
	//ChangeReq  up 放大变换
	ChangeReq struct {
		Action     string `json:"action"`     //	UPSCALE(放大); VARIATION(变换); REROLL(重新生成),可用值:UPSCALE,VARIATION,REROLL,示例值(UPSCALE)
		Index      int    `json:"index"`      //	序号(1~4), action为UPSCALE,VARIATION时必传,示例值(1)
		NotifyHook string `json:"notifyHook"` //	回调地址, 为空时使用全局notifyHook
		State      string `json:"state"`      // 	自定义参数
		TaskId     string `json:"taskId"`     // 任务ID,示例值(1320098173412546)
	}

	// DescribeReq 反推
	DescribeReq struct {
		Base64     string `json:"base64"`
		NotifyHook string `json:"notifyHook"`
		State      string `json:"state"`
	}

	// BlendReq 混图
	BlendReq struct {
		Base64Array []string `json:"base64Array"`
		Dimensions  string   `json:"dimensions"`
		NotifyHook  string   `json:"notifyHook"`
		State       string   `json:"state"`
	}
)
