<template>
    <div class="page-review">
        <div class="review-wrapper" ref="fileWrapper">
            <div class="add-box">
                <el-form ref="ruleFormRef" :model="addData" :rules="addRule" label-width="auto" @submit.native.prevent>
                    <el-form-item v-if="departmentList.length > 1" prop="department_id" :label="$L('选择部门')">
                        <el-select style="width: 100%;" v-model="addData.department_id" :placeholder="$L('请选择部门')">
                            <el-option v-for="(item, index) in departmentList" :value="item.id" :label="item.name"
                                :key="index">{{ item.name
                                }}</el-option>
                        </el-select>
                    </el-form-item>
                    <el-form-item prop="applyType" :label="$L('申请类型')">
                        <el-select style="width: 100%;" v-model="addData.applyType" :placeholder="$L('请选择申请类型')">
                            <el-option v-for="(item, index) in procdefList" :value="item.name" :label="item.name"
                                :key="index">{{ item.name
                                }}</el-option>
                        </el-select>
                    </el-form-item>
                    <el-form-item v-if="(addData.applyType || '').indexOf('请假') !== -1" prop="type" :label="$L('假期类型')">
                        <el-select style="width: 100%;" v-model="addData.type" :placeholder="$L('请选择假期类型')">
                            <el-option v-for="(item, index) in selectTypes" :value="item" :key="index">{{ $L(item)
                            }}</el-option>
                        </el-select>
                    </el-form-item>
                    <el-form-item prop="startTime" :label="$L('开始时间')">
                        <div style="display: flex;gap: 3px;">
                            <el-date-picker type="date" value-format="YYYY-MM-DD" v-model="addData.startTime"
                                :editable="false" @on-change="(e) => { addData.startTime = e }" :placeholder="$L('请选择开始时间')"
                                style="flex: 1;min-width: 122px;"></el-date-picker>
                            <el-select v-model="addData.startTimeHour" style="max-width: 100px;">
                                <el-option v-for="(item, index) in 24" :value="item - 1 < 10 ? '0' + (item - 1) : item - 1"
                                    :key="index">{{ item - 1 < 10 ? '0' : '' }}{{ item - 1 }}</el-option>
                            </el-select>
                            <el-select v-model="addData.startTimeMinute" style="max-width: 100px;">
                                <el-option value="00">00</el-option>
                                <el-option value="30">30</el-option>
                            </el-select>
                        </div>
                    </el-form-item>
                    <el-form-item prop="endTime" :label="$L('结束时间')">
                        <div style="display: flex;gap: 3px;">
                            <el-date-picker type="date" value-format="YYYY-MM-DD" v-model="addData.endTime"
                                :editable="false" @on-change="(e) => { addData.endTime = e }" :placeholder="$L('请选择结束时间')"
                                style="flex: 1;min-width: 122px;"></el-date-picker>
                            <el-select v-model="addData.endTimeHour" style="max-width: 100px;">
                                <el-option v-for="(item, index) in 24"
                                    :value="item - 1 < 10 ? '0' + (item - 1) : ((item - 1) + '')" :key="index">{{ item - 1 <
                                        10 ? '0' : '' }}{{ item - 1 }}</el-option>
                            </el-select>
                            <el-select v-model="addData.endTimeMinute" style="max-width: 100px;">
                                <el-option value="00">00</el-option>
                                <el-option value="30">30</el-option>
                            </el-select>
                        </div>
                    </el-form-item>
                    <el-form-item prop="description" :label="$L('事由')">
                        <el-input type="textarea" v-model="addData.description"></el-input>
                    </el-form-item>
                    <el-form-item prop="description" :label="$L('图片')">
                        <ImgUpload @update:img="upImg"></ImgUpload>
                    </el-form-item>
                    <el-form-item class="adaption">
                        <el-button type="default" @click="">{{ $L('取消') }}</el-button>
                        <el-button type="primary" :loading="loadIng > 0" @click="onInitiate">{{ $L('确认') }}</el-button>
                    </el-form-item>
                </el-form>
            </div>


        </div>



    </div>
</template>
  
<script setup >
import { ref, reactive, onMounted, getCurrentInstance, nextTick } from "vue";
import { getUserInfo, getProcdefFindAll, submitAnApplication } from "@/plugins/api.js";
import ImgUpload from '@/components/ImgUpload.vue';
import { ElForm } from 'element-plus';
const { proxy } = getCurrentInstance()


onMounted(() => {
    const queryString = window.location.search.slice(1);
    const params = {};
    const pairs = queryString.split('&');
    for (let i = 0; i < pairs.length; i++) {
        const pair = pairs[i].split('=');
        params[pair[0]] = decodeURIComponent(pair[1] || '');
    }
    if (params.token) {
        getUserData(params.token)
        getProcdef();
    }
});
let userID = ref(0)
// 获取用户信息
const getUserData = async (token) => {
    const upToken = token;
    let { data, status, message } = await getUserInfo(upToken)
    if (status != 200) {
        ElMessage.error(message)
        return;
    }
    if (data) {
        userID.value = data.userid;
        departmentList.value = data.departmentLists;
    }
}
//获取请假类型
const getProcdef = async () => {
    const upData = { name: '' };
    let { data, status, message } = await getProcdefFindAll(upData)
    if (status != 200) {
        ElMessage.error(message)
        return;
    }
    if (data) {
        procdefList.value = data.rows
    }
}

let addData = reactive({
    department_id: '',
    applyType: '',
    type: '',
    startTime: "",
    startTimeHour: "09",
    startTimeMinute: "00",
    endTime: "",
    endTimeHour: "18",
    endTimeMinute: "00",
    other: [],
})

const upImg = (data) => {
    addData.other = []
    if (data.value) {
        data.value.map((item) => {
            addData.other.push(item?.response?.data)
        })
    }
}

let loadIng = ref(0)

const ruleFormRef = ref()
const addRule = reactive({
    department_id: { type: 'number', required: true, message: proxy.$L('请选择部门！'), trigger: 'change' },
    applyType: { type: 'string', required: true, message: proxy.$L('请选择申请类型！'), trigger: 'change' },
    type: { type: 'string', required: true, message: proxy.$L('请选择假期类型！'), trigger: 'change' },
    startTime: { type: 'string', required: true, message: proxy.$L('请选择开始时间！'), trigger: 'change' },
    endTime: { type: 'string', required: true, message: proxy.$L('请选择结束时间！'), trigger: 'change' },
    description: { type: 'string', required: true, message: proxy.$L('请输入事由！'), trigger: 'change' },
})



let departmentList = ref([])
let procdefList = ref([])
let selectTypes = ref(["年假", "事假", "病假", "调休", "产假", "陪产假", "婚假", "丧假", "哺乳假"])

// 提交发起
const onInitiate = async () => {

    const obj = JSON.parse(JSON.stringify(addData))
    obj.startTime = obj.startTime + " " + obj.startTimeHour + ":" + obj.startTimeMinute;
    obj.endTime = obj.endTime + " " + obj.endTimeHour + ":" + obj.endTimeMinute;
    if (addData.other) {
        obj.other = addData.other.map((o) => {
            return o
        }).join(',')
    }
    const upData = {
        userID: userID.value.toString(),
        procName: obj.applyType,
        departmentId: obj.department_id,
        var: obj,
    }
    let { data, status, message } = await submitAnApplication(upData)
    if (status != 200) {
        ElMessage.error(message)
        return;
    }

}


</script>
  
<style lang="less">
@import "../../../css/review.less";
</style>
  