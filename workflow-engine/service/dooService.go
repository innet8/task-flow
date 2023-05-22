package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"text/template"
	"time"

	"workflow/util"
	"workflow/workflow-engine/model"
)

const (
	apiBaseURL = "http://192.168.100.219:2222/api/"
)

type Notification struct {
	Type   string `json:"type"`
	Data   *Data  `json:"data"`
	Action string `json:"action"` // pass, refuse, withdraw
}

type Data struct {
	ID          int    `json:"id"`
	ProcDefName string `json:"proc_def_name"`
	Nickname    string `json:"nickname"`
	Department  string `json:"department"`
	Type        string `json:"type"` // 假期类型
	StartTime   string `json:"start_time"`
	EndTime     string `json:"end_time"`
	IsFinished  bool   `json:"is_finished"`
}

// DooService doo服务层
type DooService struct {
	client *util.HTTPClient
	mu     sync.Mutex
}

// NewDooService 创建doo服务层实例
func NewDooService() *DooService {
	return &DooService{
		client: util.NewHTTPClient(5 * time.Second),
	}
}

// GetContentTemplate 获取通知模板
func (s *DooService) GetContentTemplate(character string, action string, data map[string]interface{}) (string, error) {
	n := &Notification{
		Type: character,
		Data: &Data{
			ID:          1,
			ProcDefName: data["ProcDefName"].(string),
			Nickname:    data["Nickname"].(string),
			Department:  data["Department"].(string),
			Type:        data["Type"].(string),
			StartTime:   data["StartTime"].(string),
			EndTime:     data["EndTime"].(string),
			IsFinished:  data["IsFinished"].(bool),
		},
		Action: action, // pass, refuse, withdraw
	}

	content, err := handleNotification(n)
	if err != nil {
		fmt.Println("Error:", err)
		return "", err
	}
	return content, nil
}

// 处理通知模板
func handleNotification(n *Notification) (string, error) {
	var contentTemplate string

	switch n.Type {
	case "notifier":
		contentTemplate = `
		<span class="open-review-details" data-id="{{.Data.ID}}">
		<b>抄送{{.Data.Nickname}}提交的「{{.Data.ProcDefName}}」记录</b>
		<div class="cause">
			<span>申请人：<span style="color:#84c56a">{{"@"}}{{.Data.Nickname}}</span> {{.Data.Department}}</span>
			<p><b>审批事由</b></p>
			{{if .Data.Type}}<span>假期类型：{{.Data.Type}}</span>{{end}}
			<span>开始时间：{{.Data.StartTime}}</span>
			<span>结束时间：{{.Data.EndTime}}</span>
		</div>
		<div class="btn-raw">
			<Button type="button" class="ivu-btn" style="flex: 1;">{{if eq .Data.IsFinished true}}已同意{{else}}查看详情{{end}}</Button>
		</div>
		</span>`
	case "submitter":
		contentTemplate = `
		<span class="open-review-details" data-id="{{.Data.ID}}">
		<b>您发起的「{{.Data.ProcDefName}}」{{if eq .Action "pass"}}已通过{{else if eq .Action "refuse"}}被{{.Data.Nickname}}拒绝{{else if eq .Action "withdraw"}}已撤销{{end}}</b>
		<div class="cause">
			<span>申请人：<span style="color:#84c56a">{{"@"}}{{.Data.Nickname}}</span> {{.Data.Department}}</span>
			<p><b>审批事由</b></p>
			{{if .Data.Type}}<span>假期类型：{{.Data.Type}}</span>{{end}}
			<span>开始时间：{{.Data.StartTime}}</span>
			<span>结束时间：{{.Data.EndTime}}</span>
		</div>
		<div class="btn-raw">
			<Button type="button" class="ivu-btn" style="flex: 1;">{{if eq .Action "pass"}}已同意{{else if eq .Action "refuse"}}已拒绝{{else if eq .Action "withdraw"}}已撤销{{end}}</Button>
		</div>
		</span>`
	case "reviewer":
		contentTemplate = `
		<span class="open-review-details" data-id="{{.Data.ID}}">
		<b>{{.Data.Nickname}}提交的「{{.Data.ProcDefName}}」待你审批</b>
		<div class="cause">
			<span>申请人：<span style="color:#84c56a">{{"@"}}{{.Data.Nickname}}</span> {{.Data.Department}}</span>
			<p><b>审批事由</b></p>
			{{if .Data.Type}}<span>假期类型：{{.Data.Type}}</span>{{end}}
			<span>开始时间：{{.Data.StartTime}}</span>
			<span>结束时间：{{.Data.EndTime}}</span>
		</div>
		<div class="btn-raw">
			{{if eq .Action "pass" "refuse" "withdraw"}}
				<Button type="button" class="ivu-btn" style="flex: 1;">{{if eq .Action "pass"}}已同意{{else if eq .Action "refuse"}}已拒绝{{else if eq .Action "withdraw"}}已撤销{{end}}</Button>
			{{else}}
				<Button type="button" class="ivu-btn ivu-btn-primary" style="flex: 1;">同意</Button>
				<Button type="button" class="ivu-btn ivu-btn-error" style="flex: 1;">拒绝</Button>
			{{end}}
		</div>
		</span>`
	}

	tmpl, err := template.New("content").Parse(contentTemplate)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, n)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

// GetUsers 获取用户信息
func (s *DooService) GetUsers(userId int) ([]byte, error) {
	post := map[string]int{
		"userid": userId,
	}
	return s.client.Post(apiBaseURL+"plugin/getUsers", post)
}

// GetDooRobot 获取机器人信息 by name
func (s *DooService) GetDooRobot(name string) ([]byte, error) {
	// 如果name为空，则返回默认名称
	if name == "" {
		name = "审批机器人"
	}
	post := map[string]string{
		"name": name,
	}
	return s.client.Post(apiBaseURL+"plugin/getBotByName", post)
}

// GetDialog 获取/创建会话
func (s *DooService) GetDialog(botId int, userId int) ([]byte, error) {
	post := map[string]int{
		"bot_id": botId,
		"userid": userId,
	}
	return s.client.Post(apiBaseURL+"plugin/getDialog", post)
}

// SendDooRobot 机器人发送信息
func (s *DooService) SendDooRobot(updateId int, dialogId int, text string, sender int) ([]byte, error) {
	// 构建请求参数
	userData := map[string]interface{}{
		"update_id": updateId,
		"dialog_id": dialogId,
		"text":      text,
		"sender":    sender,
	}
	return s.client.Post(apiBaseURL+"plugin/msg/sendtext", userData)
}

// 获取用户信息
func (s *DooService) getUserInfo(userId int) (map[string]interface{}, error) {
	user, err := s.GetUsers(userId)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	return s.unmarshalAndCheckResponse(user)
}

// 获取机器人信息
func (s *DooService) getBotInfo() (map[string]interface{}, error) {
	bot, err := s.GetDooRobot("")
	if err != nil {
		return nil, fmt.Errorf("failed to get bot info: %w", err)
	}
	return s.unmarshalAndCheckResponse(bot)
}

// 获取会话信息
func (s *DooService) getDialogInfo(botID, userID int) (map[string]interface{}, error) {
	dialog, err := s.GetDialog(botID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get dialog info: %w", err)
	}
	return s.unmarshalAndCheckResponse(dialog)
}

// 解码并检查返回数据
func (s *DooService) unmarshalAndCheckResponse(resp []byte) (map[string]interface{}, error) {
	var ret map[string]interface{}
	if err := json.Unmarshal(resp, &ret); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	retCode, ok := ret["ret"].(float64)
	if !ok {
		return nil, fmt.Errorf("invalid response format")
	}
	if retCode != 1 {
		msg, ok := ret["msg"].(string)
		if !ok {
			return nil, fmt.Errorf("failed to request")
		}
		return nil, fmt.Errorf("failed to request: %v", msg)
	}
	data, ok := ret["data"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid data format")
	}
	return data, nil
}

// 获取流程详情,并做发送通知的处理
func (s *DooService) HandleProcInfoMsg(procInstID int, action string) (map[string]interface{}, error) {
	procInfoStr, err := FindProcInstByID(procInstID)
	if err != nil {
		return nil, err
	}
	var procInfoMap map[string]interface{}
	if err := json.Unmarshal([]byte(procInfoStr), &procInfoMap); err != nil {
		return nil, fmt.Errorf("failed to unmarshal procInfo: %w", err)
	}
	if action == "" {
		// 根据state（1审批中，2通过，3拒绝，4撤回）状态判断action
		if procInfoMap["state"].(float64) == 2 {
			action = "pass"
		} else if procInfoMap["state"].(float64) == 3 {
			action = "refuse"
		} else if procInfoMap["state"].(float64) == 4 {
			action = "withdraw"
		}
	}
	// 如果isFinished为true，则发送给提交人
	if procInfoMap["isFinished"] == true {
		// 发送给提交人 只有通过或者拒绝才发送
		if action == "pass" || action == "refuse" {
			startUserIdStr := procInfoMap["startUserId"].(string)
			userID, _ := strconv.Atoi(startUserIdStr)
			s.SendSubmitterNotification(procInfoMap, userID, action)
		}
		// 发送给审核人 改变信息状态
		toProcMsgsUser, _ := model.GetProcMsgsByProcInstID(int(procInfoMap["id"].(float64)))
		for _, val := range toProcMsgsUser {
			s.SendReviewerNotification(procInfoMap, val.UserID, action)
		}
	} else if len(procInfoMap["candidate"].(string)) > 0 {
		candidate := strings.Split(procInfoMap["candidate"].(string), ",")
		// 循环发送给下一个审核人
		for _, val := range candidate {
			userID, _ := strconv.Atoi(val)
			s.SendReviewerNotification(procInfoMap, userID, action)
		}
	}
	// 抄送人
	notifier := s.handleProcNode(procInfoMap)
	if len(notifier) > 0 {
		// 循环发送给抄送人
		// 打印调试
		fmt.Println(notifier)
		for _, val := range notifier {
			targetID := val.(map[string]interface{})["targetId"].(string)
			userID, _ := strconv.Atoi(targetID)
			s.SendNotifierNotification(procInfoMap, userID, action)
		}
	}
	return procInfoMap, nil
}

// 处理流程节点返回是否有抄送人
func (s *DooService) handleProcNode(proc map[string]interface{}) []interface{} {
	// 获取流程节点
	procNode := proc["nodeInfos"].([]interface{})
	var notifier []interface{}
	for _, val := range procNode {
		node := val.(map[string]interface{})
		if node["type"].(string) == "notifier" {
			notifier = append(notifier, node["nodeUserList"].([]interface{})...)
		}
		// 判断到达的节点
		if proc["nodeID"].(string) == node["nodeId"].(string) {
			break
		}
	}
	return notifier
}

// 获取用户昵称
func (s *DooService) getUserNickname(userID int) (string, error) {
	userInfo, err := s.getUserInfo(userID)
	if err != nil {
		return "", err
	}
	return userInfo["nickname"].(string), nil
}

// 生成通知模板信息
func (s *DooService) generateContentTemplate(procInfoMap map[string]interface{}, nickname string, action string, recipientType string) (string, error) {
	s.mu.Lock() // 互斥锁 避免并发读写map
	defer s.mu.Unlock()
	contentTemplateData := map[string]interface{}{
		"ID":          procInfoMap["id"],
		"ProcDefName": procInfoMap["procDefName"],
		"Nickname":    nickname,
		"Department":  procInfoMap["department"],
		"Type":        procInfoMap["var"].(map[string]interface{})["type"],
		"StartTime":   procInfoMap["var"].(map[string]interface{})["startTime"],
		"EndTime":     procInfoMap["var"].(map[string]interface{})["endTime"],
		"IsFinished":  procInfoMap["isFinished"],
	}
	// submitter-提交的人 pass, refuse, withdraw, reviewer-审核的人, notifier-抄送的人
	content, err := s.GetContentTemplate(recipientType, action, contentTemplateData)
	if err != nil {
		return "", err
	}
	return content, nil
}

// 发送通知
func (s *DooService) sendNotification(procInfoMap map[string]interface{}, userID int, action string, recipientType string) ([]byte, error) {
	procInstID := int(procInfoMap["id"].(float64))
	var getNameUserID int
	if recipientType == "submitter" {
		getNameUserID = userID
	} else {
		startUserIdStr := procInfoMap["startUserId"].(string)
		getNameUserID, _ = strconv.Atoi(startUserIdStr)
	}
	nickname, err := s.getUserNickname(getNameUserID)
	if err != nil {
		return nil, err
	}
	content, err := s.generateContentTemplate(procInfoMap, nickname, action, recipientType)
	if err != nil {
		return nil, err
	}
	botInfo, err := s.getBotInfo()
	if err != nil {
		return nil, err
	}
	dialogInfo, err := s.getDialogInfo(int(botInfo["userid"].(float64)), userID)
	if err != nil {
		return nil, err
	}
	var updateID int
	if action == "pass" || action == "refuse" || action == "withdraw" {
		if recipientType == "submitter" && action != "withdraw" {
			// 跳出if判断，不更改updateID
		} else {
			updateID, _ = model.GetMsgIDByProcInstIDAndUserID(procInstID, userID)
		}
	} else {
		updateID = 0
	}
	log.Printf("action: ", action)
	log.Printf("dialogInfo: %d", dialogInfo["id"])
	log.Printf("updateID: %d", updateID)

	send, err := s.SendDooRobot(updateID, int(dialogInfo["id"].(float64)), content, int(botInfo["userid"].(float64)))
	if err != nil {
		return nil, err
	}
	if updateID == 0 && recipientType == "reviewer" {
		sendInfo, err := s.unmarshalAndCheckResponse(send)
		if err != nil {
			return nil, err
		}
		msgID := int(sendInfo["id"].(float64))
		model.InsertProcMsgs(procInstID, userID, msgID)
	}
	return send, nil
}

// SendReviewerNotification 发送审核人通知
func (s *DooService) SendReviewerNotification(procInfoMap map[string]interface{}, userID int, action string) ([]byte, error) {
	send, err := s.sendNotification(procInfoMap, userID, action, "reviewer")
	if err != nil {
		return nil, err
	}
	return send, nil
}

// SendSubmitterNotification 发送提交人通知
func (s *DooService) SendSubmitterNotification(procInfoMap map[string]interface{}, userID int, action string) ([]byte, error) {
	send, err := s.sendNotification(procInfoMap, userID, action, "submitter")
	if err != nil {
		return nil, err
	}
	return send, nil
}

// SendNotifierNotification 发送抄送人通知
func (s *DooService) SendNotifierNotification(procInfoMap map[string]interface{}, userID int, action string) ([]byte, error) {
	send, err := s.sendNotification(procInfoMap, userID, action, "notifier")
	if err != nil {
		return nil, err
	}
	return send, nil
}
