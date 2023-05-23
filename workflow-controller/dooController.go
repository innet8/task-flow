package controller

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"workflow/util"
	"workflow/workflow-engine/service"
)

// GetDooRobot 获取机器人信息 by name
func GetDooRobot(w http.ResponseWriter, r *http.Request) {
	dooRobotSvc := service.NewDooService()
	resp, err := dooRobotSvc.GetDooRobot("")
	if err != nil {
		handleError(w, err)
		return
	}
	fmt.Fprint(w, string(resp))
}

// SendDooRobot 机器人发送信息
func SendDooRobot(w http.ResponseWriter, r *http.Request) {
	dooRobotSvc := service.NewDooService()

	// 获取模板内容
	var context string
	// context, _ = dooRobotSvc.GetContentTemplate()

	updateId := 0
	dialogId := 50
	text := context
	sender := 10
	resp, err := dooRobotSvc.SendDooRobot(updateId, dialogId, text, sender)
	if err != nil {
		handleError(w, err)
		return
	}
	fmt.Fprint(w, string(resp))
}

// GetDialog 获取/创建会话
func GetDialog(w http.ResponseWriter, r *http.Request) {
	botId := 10
	userId := 4
	dooRobotSvc := service.NewDooService()
	resp, err := dooRobotSvc.GetDialog(botId, userId)
	if err != nil {
		handleError(w, err)
		return
	}
	fmt.Fprint(w, string(resp))
}

// 流程启动发送通知
func SendNotification(w http.ResponseWriter, r *http.Request) {
	service.NewDooService().HandleProcInfoMsg(20, "")
	// resp, err := dooRobotSvc.SendNotifierNotification(userId)
	// resp, err := dooRobotSvc.SendSubmitterNotification(userId)
	// if err != nil {
	// 	handleError(w, err)
	// 	return
	// }
	// fmt.Fprint(w, resp)
}

// 导出Excel文件
func Export(w http.ResponseWriter, r *http.Request) {
	// 打印调试信息
	fmt.Println("开始导出Excel文件...")
	// 数据标题
	headings := []string{
		"申请编号",
		"标题",
		"申请状态",
		"发起时间",
		"完成时间",
		"发起人工号",
		"发起人User ID",
		"发起人姓名",
		"发起人部门",
		"发起人部门ID",
		"部门负责人",
		"历史审批人",
		"历史办理人",
		"审批记录",
		"当前处理人",
		"审批节点",
		"审批人数",
		"审批耗时",
		"假期类型",
		"开始时间",
		"结束时间",
		"时长",
		"请假事由",
		"请假单位",
	}
	// 读取数据记录
	var receiver = service.GetDefaultProcessPageReceiver()
	// todo r.Body赋固定值
	r.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(`{"isFinished":1,"procName":"","state":0,"startTime":"","endTime":""}`)))
	err := util.Body2Struct(r, &receiver)
	if err != nil {
		util.ResponseErr(w, err)
		return
	}
	data, err := service.NewDooService().GetProcExportData(receiver)
	if err != nil {
		util.ResponseErr(w, err)
		return
	}
	// 调用导出Excel文件的函数
	filename := "审批记录_" + time.Now().Format("2006-01-02_15-04-05") + ".xlsx"
	err = service.NewDooService().ExportExcel(filename, headings, data)
	if err != nil {
		handleError(w, err)
	}
	// 延迟一分钟删除文件
	time.AfterFunc(time.Minute, func() {
		os.Remove(filename)
	})
	// 下载文件
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", "attachment; filename="+filename)
	http.ServeFile(w, r, "./"+filename)
}

func handleError(w http.ResponseWriter, err error) {
	fmt.Fprint(w, "Request error: ", err)
}
