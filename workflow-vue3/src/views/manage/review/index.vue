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

  </div>
</template>

<script setup>
import { ref, onMounted, getCurrentInstance, nextTick } from "vue";
import { getBacklogData, getNotifyData , getDoneData ,getInitiatedData ,getUserInfo} from "@/plugins/api.js";
import { ElMessage } from 'element-plus';
import list from "./list.vue";
import listDetails from "./details.vue";
const { proxy } = getCurrentInstance()
let minDate = ref(new Date(2020, 0, 1))
let maxDate = ref(new Date(2025, 10, 1))
let currentDate = (new Date(2021, 0, 17))

let userID = ref(0)

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
  const queryString = window.location.search.slice(1);
  const params = {};
  const pairs = queryString.split('&');
  for (let i = 0; i < pairs.length; i++) {
    const pair = pairs[i].split('=');
    params[pair[0]] = decodeURIComponent(pair[1] || '');
  }
  if(params.token){
    getUserData(params.token)
  }
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

// 获取用户信息
const getUserData = async (token) => {

const upToken = token;
let { data, status, message } = await getUserInfo(upToken)
if (status != 200) {
  ElMessage.error(message)
  return;
}
  if(data){
    userID.value = data.userid;
    tabsValue.value = "backlog"
    getBacklogList()
  }
}


// 获取待办列表
const getBacklogList = async () => {


  const getDate = {
    userID:userID.value.toString(),
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
    userID:userID.value.toString(),
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
