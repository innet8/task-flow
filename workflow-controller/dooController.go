package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"workflow/util"
	"workflow/workflow-engine/service"
)

// 导出Excel文件
func Export(w http.ResponseWriter, r *http.Request) {
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
		util.ResponseErr(w, err)
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

// VerifyToken 验证token
func VerifyToken(w http.ResponseWriter, r *http.Request) {
	token := service.GetToken(r)
	fmt.Println("token: ", token)
	// 验证token
	dooRobotSvc := service.NewDooService()
	resp, err := dooRobotSvc.ValidateToken(token)
	if err != nil {
		util.ResponseErr(w, err)
		return
	}
	respStr, err := json.Marshal(resp)
	if err != nil {
		util.ResponseErr(w, err)
	}
	respStrPlain := string(respStr)
	util.ResponseData(w, respStrPlain)
}

// 上传文件
func Upload(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("file")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		util.ResponseErr(w, "Bad Request")
		return
	}
	defer file.Close()

	uploadDir := filepath.Join(".", "workflow-vue3", "public", "uploads")
	filename, err := util.UploadFile(file, header, uploadDir)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		util.ResponseErr(w, "Internal Server Error")
		return
	}
	// 返回路径去掉workflow-vue3
	filename = filename[14:]
	util.ResponseData(w, filename)
}
