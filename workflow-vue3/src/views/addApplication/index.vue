<template>
    <div class="page-review">
        <div class="review-wrapper" ref="fileWrapper">
            <div class="review-head">
                <div class="review-nav">
                    <h1>{{ $L('添加申请') }}</h1>
                </div>

                <el-form ref="initiateRef" :model="addData" :rules="addRule" label-width="auto" @submit.native.prevent>
                    <el-el-form-item v-if="departmentList.length > 1" prop="department_id" :label="$L('选择部门')">
                        <el-select v-model="addData.department_id" :placeholder="$L('请选择部门')">
                            <el-option v-for="(item, index) in departmentList" :value="item.id" :key="index">{{ item.name
                            }}</el-option>
                        </el-select>
                    </el-el-form-item>
                    <el-el-form-item prop="applyType" :label="$L('申请类型')">
                        <el-select v-model="addData.applyType" :placeholder="$L('请选择申请类型')">
                            <el-option v-for="(item, index) in procdefList" :value="item.name" :key="index">{{ item.name
                            }}</el-option>
                        </el-select>
                    </el-el-form-item>
                    <el-el-form-item v-if="(addData.applyType || '').indexOf('请假') !== -1" prop="type" :label="$L('假期类型')">
                        <el-select v-model="addData.type" :placeholder="$L('请选择假期类型')">
                            <el-option v-for="(item, index) in selectTypes" :value="item" :key="index">{{ $L(item)
                            }}</el-option>
                        </el-select>
                    </el-el-form-item>
                    <el-el-form-item prop="startTime" :label="$L('开始时间')">
                        <div style="display: flex;gap: 3px;">
                            <DatePicker type="date" el-format="yyyy-MM-dd" v-model="addData.startTime" :editable="false"
                                @on-change="(e) => { addData.startTime = e }" :placeholder="$L('请选择开始时间')"
                                style="flex: 1;min-width: 122px;"></DatePicker>
                            <el-select v-model="addData.startTimeHour" style="max-width: 100px;">
                                <el-option v-for="(item, index) in 24" :value="item - 1 < 10 ? '0' + (item - 1) : item - 1"
                                    :key="index">{{ item - 1 < 10 ? '0' : '' }}{{ item - 1 }}</el-option>
                            </el-select>
                            <el-select v-model="addData.startTimeMinute" style="max-width: 100px;">
                                <el-option value="00">00</el-option>
                                <el-option value="30">30</el-option>
                            </el-select>
                        </div>
                    </el-el-form-item>
                    <el-el-form-item prop="endTime" :label="$L('结束时间')">
                        <div style="display: flex;gap: 3px;">
                            <DatePicker type="date" el-format="yyyy-MM-dd" v-model="addData.endTime" :editable="false"
                                @on-change="(e) => { addData.endTime = e }" :placeholder="$L('请选择结束时间')"
                                style="flex: 1;min-width: 122px;"></DatePicker>
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
                    </el-el-form-item>
                    <el-el-form-item prop="description" :label="$L('事由')">
                        <Input type="textarea" v-model="addData.description"></Input>
                    </el-el-form-item>
                    <el-el-form-item class="adaption">
                        <el-button type="default" @click="addShow = false">{{ $L('取消') }}</el-button>
                        <el-button type="primary" :loading="loadIng > 0" @click="onInitiate">{{ $L('确认') }}</el-button>
                    </el-el-form-item>
                </el-form>

            </div>



        </div>



    </div>
</template>
  
<script setup lang="ts">
import { ref, onMounted, getCurrentInstance, nextTick } from "vue";
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
</script>
  
<style lang="less">
@import "../../../css/review.less";
</style>
  