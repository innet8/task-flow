<template>
  <div class="page-review">
    <div class="review-wrapper" ref="fileWrapper">
      <div class="review-head">
        <div class="review-nav">
          <h1>{{ $L('审批中心') }}</h1>
        </div>
        <el-button type="primary" @click="addApply">{{ $L("添加申请") }}</el-button>
        <!-- <Button v-for="(item,key) in procdefList" :loading="loadIng > 0" :key="key" type="primary" @click="initiate(item)" style="margin-right:10px;">{{$L(item.name)}}</Button> -->
      </div>

      <el-tabs v-model="tabsValue" @tab-change="tabsClick" style="margin: 0 20px;height: 100%;">
        <el-tab-pane :label="$L('待办') + (backlogTotal > 0 ? ('(' + backlogTotal + ')') : '')" name="backlog"
          style="height: 100%;">
          <div class="review-main-search">
            <div style="display: flex;gap: 10px;">
              <el-select v-model="approvalType" @change="tabsClick('')" style="width: 150px;">
                <el-option v-for="item in approvalList" :value="item.value" :label="item.label" :key="item.value" />
              </el-select>
            </div>
          </div>
          <div class="noData" v-if="backlogList.length == 0">{{ $L('暂无数据') }}</div>
          <div v-else class="review-mains">
            <div class="review-main-left">
              <div class="review-main-list">
                <div @click.stop="clickList(item, key)" v-for="(item, key) in backlogList">
                  <list :class="{ 'review-list-active': item._active }" :data="item"></list>
                </div>
              </div>
            </div>
            <div class="review-main-right">
              <listDetails v-if="!detailsShow && tabsValue == 'backlog'" :data="details" @approve="tabsClick('backlog')"
                @revocation="tabsClick('backlog')"></listDetails>
            </div>
          </div>
        </el-tab-pane>
        <el-tab-pane :label="$L('已办')" name="done">
          <div class="review-main-search">
            <div style="display: flex;gap: 10px;">
              <el-select v-model="approvalType" @change="tabsClick('')" style="width: 150px;">
                <el-option v-for="item in approvalList" :value="item.value" :label="item.label" :key="item.value" />
              </el-select>
            </div>
          </div>
          <div v-if="doneList.length == 0" class="noData">{{ $L('暂无数据') }}</div>
          <div v-else class="review-mains">
            <div class="review-main-left">
              <div class="review-main-list">
                <div @click.stop="clickList(item, key)" v-for="(item, key) in doneList">
                  <list :class="{ 'review-list-active': item._active }" :data="item"></list>
                </div>
              </div>
            </div>
            <div class="review-main-right">
              <listDetails v-if="!detailsShow && tabsValue == 'done'" :data="details" @approve="tabsClick('done')"
                @revocation="tabsClick('done')"></listDetails>
            </div>
          </div>
        </el-tab-pane>
        <el-tab-pane :label="$L('抄送我')" name="notify">
          <div class="review-main-search">
            <div class="review-main-search">
              <div style="display: flex;gap: 10px;">
                <el-select v-model="approvalType" @change="tabsClick('')" style="width: 150px;">
                  <el-option v-for="item in approvalList" :value="item.value" :label="item.label" :key="item.value" />
                </el-select>
              </div>
            </div>
          </div>
          <div class="noData" v-if="notifyList.length == 0">{{ $L('暂无数据') }}</div>
          <div v-else class="review-mains">
            <div class="review-main-left">
              <div class="review-main-list">
                <div @click.stop="clickList(item, key)" v-for="(item, key) in notifyList">
                  <list :class="{ 'review-list-active': item._active }" :data="item"></list>
                </div>
              </div>
            </div>
            <div class="review-main-right">
              <listDetails v-if="!detailsShow && tabsValue == 'notify'" :data="details" @approve="tabsClick('notify')"
                @revocation="tabsClick('notify')"></listDetails>
            </div>
          </div>
        </el-tab-pane>
        <el-tab-pane :label="$L('已发起')" name="initiated">
          <div class="review-main-search">
            <div style="display: flex;gap: 10px;">
              <el-select v-model="approvalType" @change="tabsClick('')" style="width: 150px;">
                <el-option v-for="item in approvalList" :value="item.value" :label="item.label" :key="item.value" />
              </el-select>
              <el-select v-model="searchState" @change="tabsClick('')" style="width: 150px;">
                <el-option v-for="item in searchStateList" :value="item.value" :label="item.label" :key="item.value" />
              </el-select>
            </div>
          </div>
          <div class="noData" v-if="initiatedList.length == 0">{{ $L('暂无数据') }}</div>
          <div v-else class="review-mains">
            <div class="review-main-left">
              <div class="review-main-list">
                <div @click.stop="clickList(item, key)" v-for="(item, key) in initiatedList">
                  <list :class="{ 'review-list-active': item._active }" :data="item"></list>
                </div>
              </div>
            </div>
            <div class="review-main-right">
              <listDetails v-if="!detailsShow && tabsValue == 'initiated'" :data="details"
                @approve="tabsClick('initiated')" @revocation="tabsClick('notify')"></listDetails>
            </div>
          </div>
        </el-tab-pane>
      </el-tabs>

    </div>

    <!--详情-->
    <!--        <DrawerOverlay v-model="detailsShow"  placement="right" :size="600">-->
    <!--            <listDetails v-if="detailsShow" :data="details" @approve="tabsClick" @revocation="tabsClick" style="height: 100%;border-radius: 10px;"></listDetails>-->
    <!--        </DrawerOverlay>-->

    <!--        &lt;!&ndash;发起&ndash;&gt;-->
    <!--        <Modal v-model="addShow" :title="$L(addTitle)" :mask-closable="false" class="page-review-initiate">-->
    <!--            <Form ref="initiateRef" :model="addData" :rules="addRule" label-width="auto" @submit.native.prevent>-->
    <!--                <FormItem v-if="departmentList.length>1" prop="department_id" :label="$L('选择部门')">-->
    <!--                    <el-select v-model="addData.department_id" :placeholder="$L('请选择部门')">-->
    <!--                        <el-option v-for="(item, index) in departmentList" :value="item.id" :key="index">{{ item.name }}</el-option>-->
    <!--                    </el-select>-->
    <!--                </FormItem>-->
    <!--                <FormItem prop="applyType" :label="$L('申请类型')">-->
    <!--                    <el-select v-model="addData.applyType" :placeholder="$L('请选择申请类型')">-->
    <!--                        <el-option v-for="(item, index) in procdefList" :value="item.name" :key="index">{{ item.name }}</el-option>-->
    <!--                    </el-select>-->
    <!--                </FormItem>-->
    <!--                <FormItem v-if="(addData.applyType || '').indexOf('请假') !== -1" prop="type" :label="$L('假期类型')">-->
    <!--                    <el-select v-model="addData.type" :placeholder="$L('请选择假期类型')">-->
    <!--                        <el-option v-for="(item, index) in selectTypes" :value="item" :key="index">{{ $L(item) }}</el-option>-->
    <!--                    </el-select>-->
    <!--                </FormItem>-->
    <!--                <FormItem prop="startTime" :label="$L('开始时间')">-->
    <!--                    <div style="display: flex;gap: 3px;">-->
    <!--                        <DatePicker type="date" format="yyyy-MM-dd"-->
    <!--                            v-model="addData.startTime"-->
    <!--                            :editable="false"-->
    <!--                            @on-change="(e)=>{ addData.startTime = e }"-->
    <!--                            :placeholder="$L('请选择开始时间')"-->
    <!--                            style="flex: 1;min-width: 122px;"-->
    <!--                        ></DatePicker>-->
    <!--                        <el-select v-model="addData.startTimeHour" style="max-width: 100px;">-->
    <!--                            <el-option v-for="(item,index) in 24" :value="item-1 < 10 ? '0'+(item-1) : item-1 " :key="index">{{item-1 < 10 ? '0' : ''}}{{item-1}}</el-option>-->
    <!--                        </el-select>-->
    <!--                        <el-select v-model="addData.startTimeMinute" style="max-width: 100px;">-->
    <!--                            <el-option value="00">00</el-option>-->
    <!--                            <el-option value="30">30</el-option>-->
    <!--                        </el-select>-->
    <!--                    </div>-->
    <!--                </FormItem>-->
    <!--                <FormItem prop="endTime" :label="$L('结束时间')">-->
    <!--                    <div style="display: flex;gap: 3px;">-->
    <!--                        <DatePicker type="date" format="yyyy-MM-dd"-->
    <!--                            v-model="addData.endTime"-->
    <!--                            :editable="false"-->
    <!--                            @on-change="(e)=>{ addData.endTime = e }"-->
    <!--                            :placeholder="$L('请选择结束时间')"-->
    <!--                            style="flex: 1;min-width: 122px;"-->
    <!--                        ></DatePicker>-->
    <!--                        <el-select v-model="addData.endTimeHour" style="max-width: 100px;">-->
    <!--                            <el-option v-for="(item,index) in 24" :value="item-1 < 10 ? '0'+(item-1) : ((item-1)+'') " :key="index">{{item-1 < 10 ? '0' : ''}}{{item-1}}</el-option>-->
    <!--                        </el-select>-->
    <!--                        <el-select v-model="addData.endTimeMinute" style="max-width: 100px;">-->
    <!--                            <el-option value="00">00</el-option>-->
    <!--                            <el-option value="30">30</el-option>-->
    <!--                        </el-select>-->
    <!--                    </div>-->
    <!--                </FormItem>-->
    <!--                <FormItem prop="description" :label="$L('事由')">-->
    <!--                    <Input type="textarea" v-model="addData.description"></Input>-->
    <!--                </FormItem>-->
    <!--            </Form>-->
    <!--            <div slot="footer" class="adaption">-->
    <!--                <Button type="default" @click="addShow=false">{{$L('取消')}}</Button>-->
    <!--                <Button type="primary" :loading="loadIng > 0" @click="onInitiate">{{$L('确认')}}</Button>-->
    <!--            </div>-->
    <!--        </Modal>-->

  </div>
</template>

<script setup>
import { ref, onMounted, getCurrentInstance, nextTick } from "vue";
import { getBacklogData, getNotifyData , getDoneData ,getInitiatedData} from "@/plugins/api.js";
import { ElMessage } from 'element-plus';
import list from "./list.vue";
import listDetails from "./details.vue";
const { proxy } = getCurrentInstance()
let minDate = ref(new Date(2020, 0, 1))
let maxDate = ref(new Date(2025, 10, 1))
let currentDate = (new Date(2021, 0, 17))

let userID = ref(1)

let procdefList = ref([])
let page = ref(1)
let pageSize = ref(250)
let total = ref(0)
let noText = ref('')
let loadIng = ref(false)

let tabsValue = ref("")
//
let approvalType = ref("all")
let approvalList = ref([
  { value: "all", label: proxy.$L("全部审批") },
  { value: "请假", label: proxy.$L("请假") },
  { value: "加班申请", label: proxy.$L("加班申请") },
])
let searchState = ref(0)
let searchStateList = ref([
  { value: 0, label: proxy.$L("全部状态") },
  { value: 1, label: proxy.$L("审批中") },
  { value: 2, label: proxy.$L("已通过") },
  { value: 3, label: proxy.$L("已拒绝") },
  { value: 4, label: proxy.$L("已撤回") }
])
//
let backlogTotal = ref(0)
let backlogList = ref([])
let doneList = ref([])
let notifyList = ref([])
let initiatedList = ref([])
//
let details = ref({})
let detailsShow = ref(false)
//
let addTitle = ref('')
let addShow = ref(false)
let startTimeOpen = ref(false)
let endTimeOpen = ref(false)
let addData = ref({
  department_id: 0,
  applyType: '',
  type: '',
  startTime: "2023-04-20",
  startTimeHour: "09",
  startTimeMinute: "00",
  endTime: "2023-04-20",
  endTimeHour: "18",
  endTimeMinute: "00",

})
let addRule = ref({
  department_id: { type: 'number', required: true, message: proxy.$L('请选择部门！'), trigger: 'change' },
  applyType: { type: 'string', required: true, message: proxy.$L('请选择申请类型！'), trigger: 'change' },
  type: { type: 'string', required: true, message: proxy.$L('请选择假期类型！'), trigger: 'change' },
  startTime: { type: 'string', required: true, message: proxy.$L('请选择开始时间！'), trigger: 'change' },
  endTime: { type: 'string', required: true, message: proxy.$L('请选择结束时间！'), trigger: 'change' },
  description: { type: 'string', required: true, message: proxy.$L('请输入事由！'), trigger: 'change' },
})
let selectTypes = ref(["年假", "事假", "病假", "调休", "产假", "陪产假", "婚假", "丧假", "哺乳假"])

//
let showDateTime = ref(false)



// watch: {
//     '$route' (to, from) {
//         if(to.name == 'manage-review'){
//             this.tabsClick()
//         }
//     },
//     wsMsg: {
//         handler(info) {
//             const {type, action} = info;
//             switch (type) {
//                 case 'workflow':
//                     if (action == 'backlog') {
//                         this.tabsClick()
//                     }
//                     break;
//             }
//         },
//         deep: true,
//     },
// },

onMounted(() => {
  tabsValue.value = "backlog"
  tabsClick()
  getBacklogList()
});

// mounted() {
//     this.tabsValue = "backlog"
//     this.tabsClick()
//     this.getBacklogList()
//     this.addData.department_id = this.userInfo.department[0] || 0;
//     this.addData.startTime = this.addData.endTime = this.getCurrentDate();
// },

//     getCurrentDate() {
//         const today = new Date();
//         const year = today.getFullYear();
//         const month = String(today.getMonth() + 1).padStart(2, '0');
//         const date = String(today.getDate()).padStart(2, '0');
//         return `${year}-${month}-${date}`;
//     },
// tab切换事件
const tabsClick = (val) => {
  tabsValue.value = val || tabsValue.value
  if (val != "") {
    approvalType.value = "all"
    searchState.value = 0
  }
  if (tabsValue.value == 'backlog') {
    getBacklogList();
  }
  if (tabsValue.value == 'done') {
    getDoneList();
  }
  if (tabsValue.value == 'notify') {
    getNotifyList();
  }
  if (tabsValue.value == 'initiated') {
    getInitiatedList();
  }
}

// 列表点击事件
const clickList = (item) => {
  backlogList.value.map(h => { h._active = false; })
  doneList.value.map(h => { h._active = false; })
  notifyList.value.map(h => { h._active = false; })
  initiatedList.value.map(h => { h._active = false; })
  item._active = true;

  if (window.innerWidth < 426) {
    this.goForward({ name: 'manage-review-details', query: { id: item.id } });
    return;
  }
  if (window.innerWidth < 1010) {
    detailsShow.value = true;
  }
  details.value = {}
  nextTick(() => {
    details.value = item
  })
}

// 获取待办列表
const getBacklogList = async () => {


  const getDate = {
    page: page.value,
    pageSize: pageSize.value,
    procName: approvalType.value == 'all' ? '' : approvalType.value,
  };
  let { data, status, message } = await getBacklogData(getDate)
  if (status != 200) {
    ElMessage.error(message)
    return;
  }
  if (data.rows) {
    backlogList.value = data.rows.map((h, index) => {
      h._active = index == 0;
      return h;
    })
    if (approvalType.value == 'all') {
      backlogTotal.value = backlogList.value.length
    }
    if (tabsValue.value == 'backlog') {
      nextTick(() => {
        details.value = backlogList.value[0] || {}
      })
    }
  }
}

// 获取已办列表
const getDoneList = async () => {
  const getDate = {
    page: page.value,
    pageSize: pageSize.value,
    procName: approvalType.value == 'all' ? '' : approvalType.value,
  };
  let { data, status, message } = await getDoneData(getDate)
  if (status != 200) {
    ElMessage.error(message)
    return;
  }
  if (data) {
    doneList.value = data.rows.map((h, index) => {
      h._active = index == 0;
      return h;
    })
    if (tabsValue.value == 'done') {
      nextTick(() => {
        details.value = doneList.value[0] || {}
      })
    }
  }
}

// 获取抄送列表
const getNotifyList = async () => {
  const getDate = {
    userID:userID.value.toString(),
    page: page.value,
    pageSize: pageSize.value,
    procName: approvalType.value == 'all' ? '' : approvalType.value,
  };
  let { data, status, message } = await getNotifyData(getDate)
  if (status != 200) {
    ElMessage.error(message)
    return;
  }
  if(data){
        notifyList.value = data.rows.map((h, index) => {
      h._active = index == 0;
      return h;
    })
    if (tabsValue.value == 'notify') {
      nextTick(() => {
        details.value = notifyList.value[0] || {}
      })
    }
  }
}

// 获取我发起的
const getInitiatedList = async () => {
  const getDate = {
    userID:userID.value.toString(),
    page: page.value,
    pageSize: pageSize.value,
    procName: approvalType.value == 'all' ? '' : approvalType.value,
    state: searchState.value == 0 ? 0 : searchState.value
  };
  let { data, status, message } = await getInitiatedData(getDate)
  if (status != 200) {
    ElMessage.error(message)
    return;
  }
      initiatedList.value = data.rows.map((h, index) => {
      h._active = index == 0;
      return h;
    })
    if (tabsValue.value == 'initiated') {
      nextTick(() => {
        details.value = initiatedList.value[0] || {}
      })
    }
}
    //
    //     // 添加申请
    //     addApply(){
    //         this.$store.dispatch("call", {
    //             url: 'users/basic',
    //             data: {
    //                 userid: [ this.userInfo.userid]
    //             },
    //             skipAuthError: true
    //         }).then(({data}) => {
    //             this.addData.department_id = data[0]?.department[0] || 0;
    //             if( !this.addData.department_id ){
    //                 $A.modalError("您当前未加入任何部门，不能发起！");
    //                 return false;
    //             }
    //             this.$store.dispatch("call", {
    //                 url: 'workflow/procdef/all',
    //                 method: 'post',
    //             }).then(({data}) => {
    //                 this.procdefList = data.rows || [];
    //                 this.addTitle = proxy.$L("添加申请");
    //                 this.addShow = true;
    //             }).catch(({msg}) => {
    //                 $A.modalError(msg);
    //             }).finally(_ => {
    //                 this.loadIng--;
    //             });
    //         }).catch(({msg}) => {
    //             $A.modalError(msg);
    //         }).finally(_ => {
    //             this.loadIng--;
    //         });
    //     },
    //
    //     // 提交发起
    //     onInitiate(){
    //         this.$refs.initiateRef.validate((valid) => {
    //             if (valid) {
    //                 this.loadIng = 1;
    //                 var obj = JSON.parse(JSON.stringify(this.addData))
    //                 // if((addTitle || '').indexOf('班') == -1){
    //                     obj.startTime = obj.startTime +" "+ obj.startTimeHour + ":" + obj.startTimeMinute;
    //                     obj.endTime = obj.endTime +" "+ obj.endTimeHour + ":" + obj.endTimeMinute;
    //                 // }
    //                 this.$store.dispatch("call", {
    //                     url: 'workflow/process/start',
    //                     data: {
    //                         proc_name: obj.applyType,
    //                         department_id: obj.department_id,
    //                         var: JSON.stringify(obj)
    //                     },
    //                     method: 'post',
    //                 }).then(({data, msg}) => {
    //                     $A.messageSuccess(msg);
    //                     this.addShow = false;
    //                     this.$refs.initiateRef.resetFields();
    //                     this.tabsClick();
    //                 }).catch(({msg}) => {
    //                     $A.modalError(msg);
    //                 }).finally(_ => {
    //                     this.loadIng--;
    //                 });
    //             }
    //         });
    //     }
    //
    //

</script>

<style lang="less">
@import "../../../css/review.less";
</style>
