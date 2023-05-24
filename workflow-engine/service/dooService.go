package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"text/template"
	"time"

	"github.com/xuri/excelize/v2"

	"workflow/util"
	"workflow/workflow-engine/model"
)

const (
	apiBaseURL = "http://192.168.100.219:2222/api/"
)

var globalToken string

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

// 校验token
func (s *DooService) checkToken(token string) ([]byte, error) {
	return s.client.PostToken(apiBaseURL+"plugin/verifyToken", "", token)
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

// 校验token
func (s *DooService) getUserInfoByToken(token string) (map[string]interface{}, error) {
	user, err := s.checkToken(token)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info by token: %w", err)
	}
	return s.unmarshalAndCheckResponse(user)
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

// 获取申请状态描述
func (s *DooService) getStateDescription(state int) string {
	stateMap := map[int]string{
		0: "全部",
		1: "审批中",
		2: "通过",
		3: "拒绝",
		4: "撤回",
	}
	if desc, ok := stateMap[state]; ok {
		return desc
	}
	return ""
}

type Process struct {
	StartTime string
}

type Participant struct {
	Type     string
	Step     int
	Username string
	Comment  string
}

func handleParticipant(process Process, participant []Participant) map[string]interface{} {
	if len(participant) == 0 {
		return make(map[string]interface{})
	}
	res := make(map[string]interface{})
	historicalApprover := make([]string, 0)
	approvedNode := 0
	approvedNum := 0
	var approvalRecord strings.Builder
	for _, val := range participant {
		if val.Type == "participant" {
			if val.Step != 0 {
				if val.Comment == "" || contains(historicalApprover, val.Username) {
					continue
				}
				historicalApprover = appendUnique(historicalApprover, strings.Split(val.Username, ","))
				approvedNode++
				approvedNum++
			}
			name := val.Username + "|"
			call := ""
			if val.Step == 0 {
				call = "发起审批" + "|"
			} else {
				call = "同意" + "|"
			}
			timeString := ""
			if val.Step == 0 {
				timeString = process.StartTime + "|"
			}
			comment := ""
			if val.Step != 0 {
				comment = val.Comment + "|"
			}
			approvalRecord.WriteString(name + call + timeString + comment)
		}
	}
	res["approval_record"] = approvalRecord.String()
	res["historical_approver"] = strings.Trim(strings.Join(historicalApprover, ";"), ";")
	res["approved_node"] = approvedNode
	res["approved_num"] = approvedNum
	res["historical_agent"] = res["historical_approver"]
	return res
}

func contains(slice []string, value string) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

func appendUnique(slice []string, values []string) []string {
	for _, value := range values {
		if !contains(slice, value) {
			slice = append(slice, value)
		}
	}
	return slice
}

// GetProcExportData 获取申请的数据
func (s *DooService) GetProcExportData(receiver *ProcessPageReceiver) ([][]string, error) {
	fmt.Printf("receiver: %+v\n", receiver)
	var exportData [][]string
	datas, _, _ := model.FindAllProcIns(receiver.UserID, receiver.ProcName, receiver.State, receiver.StartTime, receiver.EndTime, receiver.IsFinished)
	result, _ := AllVar2Json(datas)
	// 处理数据
	for _, v := range result {
		var ret []string
		ret = append(ret, strconv.Itoa(v.ID))             //申请编号
		ret = append(ret, v.ProcDefName)                  //申请标题
		ret = append(ret, s.getStateDescription(v.State)) //申请状态
		ret = append(ret, v.StartTime)                    //申请时间
		ret = append(ret, v.EndTime)                      //结束时间
		ret = append(ret, "")                             // 发起人工号
		ret = append(ret, v.StartUserID)                  // 发起人ID
		ret = append(ret, v.StartUserName)                // 发起人姓名
		ret = append(ret, v.Department)                   // 发起人部门
		ret = append(ret, strconv.Itoa(v.DepartmentId))   // 发起人部门ID
		ret = append(ret, "")                             // 部门负责人
		// 查找所有与流程实例相关的参与者
		data2s, _ := model.FindParticipantAllByProcInstID(v.ID)
		participant := make([]Participant, 0)
		for _, data2 := range data2s {
			participant = append(participant, Participant{
				Type:     data2.Type,
				Step:     data2.Step,
				Username: data2.UserName,
				Comment:  data2.Comment,
			})
		}
		process := Process{
			StartTime: v.StartTime,
		}
		ParticipantMap := handleParticipant(process, participant)
		ret = append(ret, ParticipantMap["historical_approver"].(string))      // 历史审批人
		ret = append(ret, "")                                                  // 历史办理人
		ret = append(ret, ParticipantMap["approval_record"].(string))          // 审批记录
		ret = append(ret, "")                                                  // 当前处理人
		ret = append(ret, strconv.Itoa(ParticipantMap["approved_node"].(int))) // 审批节点
		ret = append(ret, strconv.Itoa(ParticipantMap["approved_num"].(int)))  // 审批人数
		// 计算审批耗时 单位小时
		startTime, _ := time.Parse("2006-01-02 15:04:05", v.StartTime)
		endTime, _ := time.Parse("2006-01-02 15:04:05", v.EndTime)
		if endTime.Before(startTime) {
			endTime = startTime
		}
		duration := endTime.Sub(startTime)
		ret = append(ret, strconv.FormatFloat(duration.Hours(), 'f', 1, 64)) // 审批耗时
		// 把var字段指针转换为map
		varMap := *v.Var
		ret = append(ret, varMap.Type)      // 假期类型
		ret = append(ret, varMap.StartTime) // 开始时间
		ret = append(ret, varMap.EndTime)   // 结束时间
		// 计算请假时长 单位小时
		startTime2, _ := time.Parse("2006-01-02 15:04", varMap.StartTime)
		endTime2, _ := time.Parse("2006-01-02 15:04", varMap.EndTime)
		duration2 := endTime2.Sub(startTime2)
		ret = append(ret, strconv.FormatFloat(duration2.Hours(), 'f', 1, 64)) // 时长
		ret = append(ret, varMap.Description)                                 // 请假事由
		ret = append(ret, "小时")                                               // 请假单位

		exportData = append(exportData, ret)
	}
	return exportData, nil
}

// 创建一个Excel文件
func createExcelFile() (*excelize.File, error) {
	xlsx := excelize.NewFile()
	xlsx.NewSheet("Sheet1")
	xlsx.SetColWidth("Sheet1", "A", "Z", 20)
	return xlsx, nil
}

// 设置Excel文件的标题行
func setExcelTitleRow(xlsx *excelize.File, title []string) error {
	boldFont, _ := xlsx.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
			Size: 12,
		},
	})
	xlsx.SetCellStyle("Sheet1", "A1", "Z1", boldFont)
	xlsx.SetSheetRow("Sheet1", "A1", &title)
	return nil
}

// 导出Excel文件
func (s *DooService) ExportExcel(filename string, title []string, data [][]string) error {
	xlsx, err := createExcelFile()
	if err != nil {
		return err
	}
	err = setExcelTitleRow(xlsx, title)
	if err != nil {
		return err
	}
	for i, row := range data {
		xlsx.SetSheetRow("Sheet1", fmt.Sprintf("A%d", i+2), &row)
	}
	// 保存Excel文件
	err = xlsx.SaveAs(filename)
	if err != nil {
		return err
	}
	return nil
}

// 校验token
func (s *DooService) ValidateToken(tokenString string) (map[string]interface{}, error) {
	// 检验 token 是否有效
	user, err := s.getUserInfoByToken(tokenString)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("Invalid token")
	}
	departments := user["department"].([]interface{})
	var lists []model.UserDepartments
	if len(user["department"].([]interface{})) >= 1 {
		var numbers []int
		for _, department := range departments {
			// 将department转换为int类型，并添加到numbers切片中
			number, err := strconv.Atoi(fmt.Sprintf("%v", department))
			if err != nil {
				return nil, err
			}
			numbers = append(numbers, number)
		}
		list, _ := model.GetDepByDeptIds(numbers)
		listsBytes, _ := json.Marshal(list)
		err := json.Unmarshal(listsBytes, &lists)
		if err != nil {
			return nil, err
		}
	}
	user["departmentLists"] = lists
	return user, nil
}

// 获取token 支持从url参数、post参数、header中获取
func GetToken(r *http.Request) string {
	token := r.URL.Query().Get("token")
	if token != "" {
		return token
	}
	token = r.PostFormValue("token")
	if token != "" {
		return token
	}
	token = r.Header.Get("Authorization")
	if token != "" {
		return strings.TrimPrefix(token, "Bearer ")
	}
	return ""
}
