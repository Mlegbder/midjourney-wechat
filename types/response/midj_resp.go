package response

type (
	TaskResp struct {
		Action      string   `json:"action"`      //可用值:IMAGINE,UPSCALE,VARIATION,REROLL,DESCRIBE,BLEND
		Description string   `json:"description"` //任务描述
		FailReason  string   `json:"failReason"`  //	失败原因
		FinishTime  int      `json:"finishTime"`  //结束时间
		Id          string   `json:"id"`          //任务ID
		ImageUrl    string   `json:"imageUrl"`    //图片url
		Progress    string   `json:"progress"`    //任务进度
		Prompt      string   `json:"prompt"`      //	提示词
		PromptEn    string   `json:"promptEn"`    //提示词-英文
		Properties  struct{} `json:"properties"`
		StartTime   int      `json:"startTime"`  //开始执行时间
		State       string   `json:"state"`      //自定义参数
		Status      string   `json:"status"`     //任务状态,可用值:NOT_START,SUBMITTED,IN_PROGRESS,FAILURE,SUCCESS
		SubmitTime  int      `json:"submitTime"` //	提交时间
	}

	MidjResp struct {
		Code        int      `json:"code"`
		Description string   `json:"description"`
		Properties  struct{} `json:"properties"`
		Result      string   `json:"result"`
	}
)
