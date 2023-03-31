package flow

import (
	"container/list"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"workflow/util"
)

// ActionConditionType 条件类型
type ActionConditionType int

const (
	RANGE ActionConditionType = iota // RANGE 条件类型: 范围
	VALUE                            // VALUE 条件类型： 值
)

// ActionConditionTypes 所有条件类型
var ActionConditionTypes = [...]string{RANGE: "dingtalk_actioner_range_condition", VALUE: "dingtalk_actioner_value_condition"}

// NodeType 节点类型
type NodeType int

const (
	// START 类型start
	START NodeType = iota
	ROUTE
	CONDITION
	APPROVER
	NOTIFIER
)

// ActionRuleType 审批人类型
type ActionRuleType int

const (
	MANAGER ActionRuleType = iota
	LABEL
)

// NodeTypes 节点类型
var NodeTypes = [...]string{START: "start", ROUTE: "route", CONDITION: "condition", APPROVER: "approver", NOTIFIER: "notifier"}
var actionRuleTypes = [...]string{MANAGER: "target_management", LABEL: "target_label"}

// 级别类型
var DirectorLevelTypes = [...]string{1: "直接主管", 2: "第二级主管", 3: "第三级主管", 4: "第四级主管"}
var ExamineEndDirectorLevelTypes = [...]string{1: "最高层级主管", 2: "第二级主管", 3: "第三级主管", 4: "第四级主管"}

type NodeInfoType int

const (
	STARTER NodeInfoType = iota
)

var NodeInfoTypes = [...]string{STARTER: "starter"}

// Node represents a specific logical unit of processing and routing
// in a workflow.
// 流程中的一个节点
type Node struct {
	Name                    string               `json:"name,omitempty"`                    // 节点名称
	Type                    string               `json:"type,omitempty"`                    // 节点类型，0 发起人 1审批 2抄送 3条件 4路由
	NodeID                  string               `json:"nodeId,omitempty"`                  // 节点id
	PrevID                  string               `json:"prevId,omitempty"`                  // 父级id
	ChildNode               *Node                `json:"childNode,omitempty"`               // 子节点
	ConditionNodes          []*Node              `json:"conditionNodes,omitempty"`          // 条件节点
	ConditionList           []*NodeConditionList `json:"conditionList,omitempty"`           // 条件列表
	Properties              *NodeProperties      `json:"properties,omitempty"`              // 属性
	Settype                 int                  `json:"settype,omitempty"`                 // 审批人设置 1指定成员 2主管 4发起人自选 5发起人自己 7连续多级主管
	SelectMode              int                  `json:"selectMode,omitempty"`              // 审批人数 1选一个人 2选多个人
	SelectRange             int                  `json:"selectRange,omitempty"`             // 选择范围 1.全公司 2指定成员 2指定角色
	PriorityLevel           int                  `json:"priorityLevel,omitempty"`           // 优先级
	DirectorLevel           int                  `json:"directorLevel,omitempty"`           // 主管级别，1直接主管， 2第二级主管，3第三级主管，4第四级主管，Settype=3才有效
	ExamineEndDirectorLevel int                  `json:"examineEndDirectorLevel,omitempty"` // 最终主管级别
	ExamineMode             int                  `json:"examineMode,omitempty"`             // 多人审批时采用的审批方式 1依次审批，2会签(须所有审批人同意)
	NoHanderAction          int                  `json:"noHanderAction,omitempty"`          // 审批人为空时 1自动审批通过/不允许发起, 2转交给审核管理员
	NodeUserList            []*NodeUser          `json:"nodeUserList,omitempty"`            // 节点用户列表
	CcSelfSelectFlag        int                  `json:"ccSelfSelectFlag"`                  // 允许发起人自选抄送人
}

// 活动规则
type ActionerRule struct {
	Type        string `json:"type,omitempty"`
	LabelNames  string `json:"labelNames,omitempty"`
	Labels      int    `json:"labels,omitempty"`
	IsEmpty     bool   `json:"isEmpty,omitempty"`
	MemberCount int8   `json:"memberCount,omitempty"` // 表示需要通过的人数 如果是会签
	ActType     string `json:"actType,omitempty"`     // and 表示会签 or表示或签，默认为或签
	Level       int8   `json:"level,omitempty"`
	AutoUp      bool   `json:"autoUp,omitempty"`
}

// 节点属性
type NodeProperties struct {
	ActivateType       string             `json:"activateType,omitempty"` // ONE_BY_ONE 代表依次审批
	AgreeAll           bool               `json:"agreeAll,omitempty"`
	Conditions         [][]*NodeCondition `json:"conditions,omitempty"`
	ActionerRules      []*ActionerRule    `json:"actionerRules,omitempty"`
	NoneActionerAction string             `json:"noneActionerAction,omitempty"`
}

// 节点条件
type NodeCondition struct {
	Type            string      `json:"type,omitempty"`
	ParamKey        string      `json:"paramKey,omitempty"`
	ParamLabel      string      `json:"paramLabel,omitempty"`
	IsEmpty         bool        `json:"isEmpty,omitempty"`
	LowerBound      string      `json:"lowerBound,omitempty"` // 类型为range
	LowerBoundEqual string      `json:"lowerBoundEqual,omitempty"`
	UpperBoundEqual string      `json:"upperBoundEqual,omitempty"`
	UpperBound      string      `json:"upperBound,omitempty"`
	BoundEqual      string      `json:"boundEqual,omitempty"`
	Unit            string      `json:"unit,omitempty"`
	ParamValues     []string    `json:"paramValues,omitempty"` // 类型为 value
	OriValue        []string    `json:"oriValue,omitempty"`
	Conds           []*NodeCond `json:"conds,omitempty"`
}

// 节点条件列表
type NodeConditionList struct {
	ColumnId     int    `json:"columnId"`
	Type         int    `json:"type,omitempty"`
	ShowType     string `json:"showType,omitempty"` //显示类型
	ShowName     string `json:"showName,omitempty"` //发起人
	OptType      string `json:"optType,omitempty"`  //选择类型
	Zdy1         string `json:"zdy1,omitempty"`
	Opt1         string `json:"opt1,omitempty"`
	Zdy2         string `json:"zdy2,omitempty"`
	Opt2         string `json:"opt2,omitempty"`
	ColumnDbname string `json:"columnDbname,omitempty"`
	ColumnType   string `json:"columnType,omitempty"`
}

// 节点气孔导度
type NodeCond struct {
	Type  string    `json:"type,omitempty"`
	Value string    `json:"value,omitempty"`
	Attrs *NodeUser `json:"attrs,omitempty"`
}

// 节点用户
type NodeUser struct {
	Name     string `json:"name,omitempty"`
	Avatar   string `json:"avatar,omitempty"`
	TargetId string `json:"targetId,omitempty"`
	Type     int    `json:"type,omitempty"`
}

// NodeInfo 节点信息
type NodeInfo struct {
	NodeID                  string `json:"nodeId"`
	Type                    string `json:"type"`
	Settype                 int    `json:"settype,omitempty"`
	Aprover                 string `json:"approver"`
	AproverType             string `json:"aproverType"`
	MemberCount             int8   `json:"memberCount"`
	Level                   int8   `json:"level"`
	ActType                 string `json:"actType"`
	DirectorLevel           int    `json:"directorLevel,omitempty"`           // 主管级别，1直接主管， 2第二级主管，3第三级主管，4第四级主管，Settype=3才有效
	ExamineEndDirectorLevel int    `json:"examineEndDirectorLevel,omitempty"` // 最终主管级别 1直接主管， 2第二级主管，3第三级主管，4第四级主管，Settype=3才有效
}

func (n *Node) add2ExecutionList(list *list.List) {
	// var NodeTypes = [...]string{START: "start", ROUTE: "route", CONDITION: "condition", APPROVER: "approver", NOTIFIER: "notifier"}
	switch n.Type {
	case NodeTypes[APPROVER], NodeTypes[NOTIFIER]: //审核人，抄送人
		var aprover string
		var memberCount int8
		memberCount = 1
		// 审批人设置 1指定成员 2主管 4发起人自选 5发起人自己 7连续多级主管
		if n.Settype == 1 {
			aprover = "指定成员"
			memberCount = int8(len(n.NodeUserList))
		} else if n.Settype == 2 {
			aprover = DirectorLevelTypes[n.DirectorLevel]
		} else if n.Settype == 4 {
			aprover = "发起人自选"
		} else if n.Settype == 5 {
			aprover = "发起人自己"
		} else if n.Settype == 7 {
			aprover = ExamineEndDirectorLevelTypes[n.ExamineEndDirectorLevel]
		}

		// 多人审批时采用的审批方式
		var actType string
		if n.ExamineMode == 1 {
			actType = "any"
		} else {
			actType = "or"
		}

		list.PushBack(NodeInfo{
			NodeID:                  n.NodeID,
			Type:                    n.Type,
			Settype:                 n.Settype,
			Aprover:                 aprover,
			AproverType:             n.Type,
			MemberCount:             memberCount,
			ActType:                 actType,
			DirectorLevel:           n.DirectorLevel,
			ExamineEndDirectorLevel: n.ExamineEndDirectorLevel,
		})
		break
	default:
	}
}

// IfProcessConifgIsValid 检查流程配置是否有效
func IfProcessConifgIsValid(node *Node) error {
	// 节点名称是否有效
	if len(node.NodeID) == 0 {
		return errors.New("节点的【nodeId】不能为空！！")
	}
	// 检查类型是否有效
	if len(node.Type) == 0 {
		return errors.New("节点【" + node.NodeID + "】的类型【type】不能为空")
	}
	var flag = false
	for _, val := range NodeTypes {
		if val == node.Type {
			flag = true
			break
		}
	}
	if !flag {
		str, _ := util.ToJSONStr(NodeTypes)
		return errors.New("节点【" + node.NodeID + "】的类型为【" + node.Type + "】，为无效类型,有效类型为" + str)
	}
	// 当前节点是否设置有审批人
	if node.Type == NodeTypes[APPROVER] || node.Type == NodeTypes[NOTIFIER] {
		if node.Properties == nil || node.Properties.ActionerRules == nil {
			return errors.New("节点【" + node.NodeID + "】的Properties属性不能为空，如：`\"properties\": {\"actionerRules\": [{\"type\": \"target_label\",\"labelNames\": \"人事\",\"memberCount\": 1,\"actType\": \"and\"}],}`")
		}
	}
	// 条件节点是否存在
	if node.ConditionNodes != nil { // 存在条件节点
		if len(node.ConditionNodes) == 1 {
			return errors.New("节点【" + node.NodeID + "】条件节点下的节点数必须大于1")
		}
		// 根据条件变量选择节点索引
		err := CheckConditionNode(node.ConditionNodes)
		if err != nil {
			return err
		}
	}

	// 子节点是否存在
	if node.ChildNode != nil && node.ChildNode.Name != "" {
		return IfProcessConifgIsValid(node.ChildNode)
	}
	return nil
}

// CheckConditionNode 检查条件节点
func CheckConditionNode(nodes []*Node) error {
	for _, node := range nodes {
		if node.Properties == nil {
			return errors.New("节点【" + node.NodeID + "】的Properties对象为空值！！")
		}
		// if len(node.Properties.Conditions) == 0 {
		// 	return errors.New("节点【" + node.NodeID + "】的Conditions对象为空值！！")
		// }
		err := IfProcessConifgIsValid(node)
		if err != nil {
			return err
		}
	}
	return nil
}

// ParseProcessConfig 解析流程定义json数据
func ParseProcessConfig(node *Node, variable *map[string]string) (*list.List, error) {
	// defer fmt.Println("----------解析结束--------")
	list := list.New()
	err := parseProcessConfig(node, variable, list)
	return list, err
}

func parseProcessConfig(node *Node, variable *map[string]string, list *list.List) (err error) {
	// fmt.Printf("nodeId=%s\n", node.NodeID)
	node.add2ExecutionList(list)

	// 存在条件节点
	if node.ConditionNodes != nil {

		// data, _ := util.ToJSONStr(node.ConditionNodes)

		// 如果条件节点只有一个或者条件只有一个，直接返回第一个
		if variable == nil || len(node.ConditionNodes) == 1 {
			err = parseProcessConfig(node.ConditionNodes[0].ChildNode, variable, list)
			if err != nil {
				return err
			}
		} else {
			// fmt.Printf("%s", data)

			// 根据条件变量选择节点索引
			condNode, err := GetConditionNode(node.ConditionNodes, variable)
			if err != nil {
				return err
			}
			if condNode == nil {
				str, _ := util.ToJSONStr(variable)
				return errors.New("节点【" + node.NodeID + "】找不到符合条件的子节点,检查变量【var】值是否匹配," + str)
				// panic(err)
			}

			err = parseProcessConfig(condNode, variable, list)
			if err != nil {
				return err
			}

		}
	}
	// 存在子节点
	if node.ChildNode != nil {
		err = parseProcessConfig(node.ChildNode, variable, list)
		if err != nil {
			return err
		}
	}
	return nil
}

// GetConditionNode 获取条件节点
func GetConditionNode(nodes []*Node, maps *map[string]string) (result *Node, err error) {
	// map2 := *maps
	for _, node := range nodes {
		var flag int
		for _, v := range node.ConditionList {
			flag++
			fmt.Printf("%x", v)
			// paramValue := map2[v.ParamKey]
			// if len(paramValue) == 0 {
			// 	return nil, errors.New("流程启动变量【var】的key【" + v.ParamKey + "】的值不能为空")
			// }
			// yes, err := checkConditions(v, paramValue)
			// if err != nil {
			// 	return nil, err
			// }
			// if yes {
			// 	flag++
			// }
		}
		// fmt.Printf("flag=%d\n", flag)
		// 满足所有条件
		if flag == len(node.ConditionList) {
			result = node
		}
	}
	return result, nil
}

// 获取条件节点
func getConditionNode(nodes []*Node, maps *map[string]string) (result *Node, err error) {
	map2 := *maps
	// 获取所有conditionNodes
	getNodesChan := func() <-chan *Node {
		nodesChan := make(chan *Node, len(nodes))
		go func() {
			// defer fmt.Println("关闭nodeChan通道")
			defer close(nodesChan)
			for _, v := range nodes {
				nodesChan <- v
			}
		}()
		return nodesChan
	}

	//获取所有conditions
	getConditionNode := func(nodesChan <-chan *Node, done <-chan interface{}) <-chan *Node {
		resultStream := make(chan *Node, 2)
		go func() {
			// defer fmt.Println("关闭resultStream通道")
			defer close(resultStream)
			for {
				select {
				case <-done:
					return
				case <-time.After(10 * time.Millisecond):
					fmt.Println("Time out.")
				case node, ok := <-nodesChan:
					if ok {
						// for _, v := range node.Properties.Conditions[0] {
						// 	conStream <- v
						// 	fmt.Printf("接收 condition:%s\n", v.Type)
						// }
						var flag int
						for _, v := range node.Properties.Conditions[0] {
							// fmt.Println(v.ParamKey)
							// fmt.Println(map2[v.ParamKey])
							paramValue := map2[v.ParamKey]
							if len(paramValue) == 0 {
								log.Printf("key:%s的值为空\n", v.ParamKey)
								// nodeAndErr.Err = errors.New("key:" + v.ParamKey + "的值为空")
								break
							}
							yes, err := checkConditions(v, paramValue)
							if err != nil {
								// nodeAndErr.Err = err
								break
							}
							if yes {
								flag++
							}
						}
						// fmt.Printf("flag=%d\n", flag)
						// 满足所有条件
						if flag == len(node.Properties.Conditions[0]) {
							// fmt.Printf("flag=%d\n,send node:%s\n", flag, node.NodeID)
							resultStream <- node
						} else {
							// fmt.Println("条件不完全满足")
						}
					}
				}
			}
		}()
		return resultStream
	}
	done := make(chan interface{})
	// defer fmt.Println("结束所有goroutine")
	defer close(done)
	nodeStream := getNodesChan()
	// for i := len(nodes); i > 0; i-- {
	// 	getConditionNode(resultStream, nodeStream, done)
	// }
	resultStream := getConditionNode(nodeStream, done)
	// for node := range resultStream {
	// 	return node, nil
	// }
	for {
		select {
		case <-time.After(1 * time.Second):
			fmt.Println("Time out")
			return
		case node := <-resultStream:
			// result = node
			return node, nil
		}
	}
	// setResult(resultStream, done)
	// time.Sleep(1 * time.Second)
	// log.Println("----------寻找节点结束--------")
	// return result, err
}

// 检查条件
func checkConditions(cond *NodeCondition, value string) (bool, error) {
	// 判断类型
	switch cond.Type {
	case ActionConditionTypes[RANGE]:
		val, err := strconv.Atoi(value)
		if err != nil {
			return false, err
		}
		if len(cond.LowerBound) == 0 && len(cond.UpperBound) == 0 && len(cond.LowerBoundEqual) == 0 && len(cond.UpperBoundEqual) == 0 && len(cond.BoundEqual) == 0 {
			return false, errors.New("条件【" + cond.Type + "】的上限或者下限值不能全为空")
		}
		// 判断下限，lowerBound
		if len(cond.LowerBound) > 0 {
			low, err := strconv.Atoi(cond.LowerBound)
			if err != nil {
				return false, err
			}
			if val <= low {
				// fmt.Printf("val:%d小于lowerBound:%d\n", val, low)
				return false, nil
			}
		}
		if len(cond.LowerBoundEqual) > 0 {
			le, err := strconv.Atoi(cond.LowerBoundEqual)
			if err != nil {
				return false, err
			}
			if val < le {
				// fmt.Printf("val:%d小于lowerBound:%d\n", val, low)
				return false, nil
			}
		}
		// 判断上限,upperBound包含等于
		if len(cond.UpperBound) > 0 {
			upper, err := strconv.Atoi(cond.UpperBound)
			if err != nil {
				return false, err
			}
			if val >= upper {
				return false, nil
			}
		}
		if len(cond.UpperBoundEqual) > 0 {
			ge, err := strconv.Atoi(cond.UpperBoundEqual)
			if err != nil {
				return false, err
			}
			if val > ge {
				return false, nil
			}
		}
		if len(cond.BoundEqual) > 0 {
			equal, err := strconv.Atoi(cond.BoundEqual)
			if err != nil {
				return false, err
			}
			if val != equal {
				return false, nil
			}
		}
		return true, nil
	case ActionConditionTypes[VALUE]:
		if len(cond.ParamValues) == 0 {
			return false, errors.New("条件节点【" + cond.Type + "】的 【paramValues】数组不能为空，值如：'paramValues:['调休','年假']")
		}
		for _, val := range cond.ParamValues {
			if value == val {
				return true, nil
			}
		}
		// log.Printf("key:" + cond.ParamKey + "找不到对应的值")
		return false, nil
	default:
		str, _ := util.ToJSONStr(ActionConditionTypes)
		return false, errors.New("未知的NodeCondition类型【" + cond.Type + "】,正确类型应为以下中的一个:" + str)
	}
}
