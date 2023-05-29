<template>
    <div class="review-details" :style="{ 'z-index': modalTransferIndex }">
        <div class="review-details-box">
            <h2 class="review-details-title">
                <span class="review-details-title-span">{{ $L(datas.procDefName) }}</span>
                <el-tag v-if="datas.state == 0" type="warning">{{ $L('待审批') }}</el-tag>
                <el-tag v-if="datas.state == 1" type="warning">{{ $L('审批中') }}</el-tag>
                <el-tag v-if="datas.state == 2" type="success">{{ $L('已通过') }}</el-tag>
                <el-tag v-if="datas.state == 3" type="danger">{{ $L('已拒绝') }}</el-tag>
                <el-tag v-if="datas.state == 4" type="danger">{{ $L('已撤回') }}</el-tag>
            </h2>
            <h3 class="review-details-subtitle">
                <el-avatar :src="datas.userimg" :size="24" /><span>{{ datas.startUserName }}</span>
            </h3>
            <h3 class="review-details-subtitle"><span>{{ $L('提交于') }} {{ datas.startTime }}</span></h3>
            <el-divider />
            <div class="review-details-text" v-if="(datas.procDefName || '').indexOf('班') == -1">
                <h4>{{ $L('假期类型') }}</h4>
                <p>{{ $L(datas.var?.type) }}</p>
            </div>
            <div class="review-details-text">
                <h4>{{ $L('开始时间') }}</h4>
                <p>{{ datas.var?.startTime }}</p>
            </div>
            <div class="review-details-text">
                <h4>{{ $L('结束时间') }}</h4>
                <p>{{ datas.var?.endTime }}</p>
            </div>
            <div class="review-details-text">
                <h4>{{ $L('时长') }}（{{ getTimeDifference(datas.var?.startTime, datas.var?.endTime)['unit'] }}）</h4>
                <p>{{ datas.var?.startTime ? getTimeDifference(datas.var?.startTime, datas.var?.endTime)['time'] : 0 }}
                </p>
            </div>
            <div class="review-details-text">
                <h4>{{ $L('请假事由') }}</h4>
                <p>{{ datas.var?.description }}</p>
            </div>
            <div class="review-details-text" v-if="datas.var?.other">
                <h4>{{ $L('图片') }}</h4>
                <div class="img-body">
                    <div v-for="(src, key) in (datas.var.other).split(',') " @click="onViewPicture(src)">
                        <img :src="src" :key="key" class="img-view" />
                    </div>
                </div>
            </div>
            <el-divider />
            <h3 class="review-details-subtitle">{{ $L('审批记录') }}</h3>
            <el-timeline style="margin-top: 20px;">
                <template v-for="(item, key) in datas.nodeInfos">

                    <!-- 提交 -->
                    <el-timeline-item :key="key" v-if="item.type == 'starter'" color="#19be6b">
                        <p class="timeline-title">{{ $L('提交') }}</p>
                        <div style="display: flex;">
                            <el-avatar :src="data.userimg || datas.userimg" :size="38" />
                            <div style="margin-left: 10px;flex: 1;">
                                <p class="review-process-name">{{ data.startUserName || datas.startUserName }}</p>
                                <p class="review-process-state">{{ $L('已提交') }}</p>
                            </div>
                            <div class="review-process-right">
                                <p v-if="parseInt(getTimeAgo(item.claimTime)) < showTimeNum">{{ getTimeAgo(item.claimTime)
                                }}</p>
                                <p>{{ item.claimTime?.substr(0, 16) }}</p>
                            </div>
                        </div>
                    </el-timeline-item>

                    <!-- 审批 -->
                    <el-timeline-item :key="key" v-if="item.type == 'approver' && item._show"
                        :color="item.identitylink ? (item.identitylink?.state > 1 ? '#f03f3f' : '#19be6b') : '#ccc'">
                        <p class="timeline-title">{{ $L('审批') }}</p>
                        <div style="display: flex;">
                            <el-avatar :src="(item.nodeUserList && item.nodeUserList[0]?.userimg) || item.userimg"
                                :size="38" />
                            <div style="margin-left: 10px;flex: 1;">
                                <p class="review-process-name">{{ item.approver }}</p>
                                <p class="review-process-state" style="color: #6d6d6d;" v-if="!item.identitylink">待审批</p>
                                <p class="review-process-state" v-if="item.identitylink">
                                    <span v-if="item.identitylink.state == 0" style="color:#496dff;">{{ $L('审批中') }}</span>
                                    <span v-if="item.identitylink.state == 1">{{ $L('已通过') }}</span>
                                    <span v-if="item.identitylink.state == 2" style="color:#f03f3f;">{{ $L('已拒绝') }}</span>
                                    <span v-if="item.identitylink.state == 3" style="color:#f03f3f;">{{ $L('已撤回') }}</span>
                                </p>
                            </div>
                            <div class="review-process-right">
                                <p v-if="parseInt(getTimeAgo(item.claimTime)) < showTimeNum">
                                    {{ item.identitylink?.state == 0 ?
                                        ($L('已等待') + " " + getTimeAgo(datas.nodeInfos[key - 1].claimTime, 2)) :
                                        (item.claimTime ? getTimeAgo(item.claimTime) : '')
                                    }}
                                </p>
                                <p>{{ item.claimTime?.substr(0, 16) }}</p>
                            </div>
                        </div>
                        <p class="comment" v-if="item.identitylink?.comment"><span>“{{ item.identitylink?.comment }}”</span>
                        </p>
                    </el-timeline-item>

                    <!-- 抄送 -->
                    <el-timeline-item :key="key" :color="item.isFinished ? '#19be6b' : '#ccc'"
                        v-if="item.type == 'notifier' && item._show">
                        <p class="timeline-title">{{ $L('抄送') }}</p>
                        <div style="display: flex;">
                            <el-avatar src="@/assets/images/default_bot.png" :size="38" />
                            <div style="margin-left: 10px;flex: 1;">
                                <p class="review-process-name">{{ $L('系统') }}</p>
                                <p style="font-size: 12px;">{{ $L('自动抄送') }}
                                    <span style="color: #486fed;">
                                        {{ item.nodeUserList?.map(h => h.name).join(',') }}
                                        {{ $L('等' + item.nodeUserList?.length + '人') }}
                                    </span>
                                </p>
                            </div>
                        </div>
                    </el-timeline-item>

                    <!-- 结束 -->
                    <el-timeline-item :key="key" :color="item.isFinished ? '#19be6b' : '#ccc'"
                        v-if="item.aproverType == 'end'">
                        <p class="timeline-title">{{ $L('结束') }}</p>
                        <div style="display: flex;">
                            <el-avatar src="@/assets/images/default_bot.png" :size="38" />
                            <div style="margin-left: 10px;flex: 1;">
                                <p class="review-process-name">{{ $L('系统') }}</p>
                                <p style="font-size: 12px;"> {{ datas.isFinished ? $L('已结束') : $L('未结束') }}</p>
                            </div>
                        </div>
                    </el-timeline-item>

                </template>

            </el-timeline>

            <template v-if="datas.globalComments">
                <el-divider />
                <h3 class="review-details-subtitle">{{ $L('全文评论') }}</h3>
                <div class="review-record-comment">

                    <div class="review-record-box" v-for="(item, key) in datas.globalComments" :key="key">

                        <div class="top">
                            <el-avatar :src="item.avatar" :size="38" />
                            <div>
                                <p>{{ item.nickName }}</p>
                                <p class="time">{{ item.createAt }}</p>
                            </div>
                            <span>{{ getTimeAgo(item.createAt, 2) }}</span>
                        </div>
                        <div class="content">
                            {{ getContent(item.content) }}
                        </div>
                        <div class="content" style="display: flex; gap: 10px;">
                            <div v-for="(src, k) in getPictures(item.content)" :key="k" @click="onViewPicture(src)">
                                <img :src="src" class="img-view" />
                            </div>
                        </div>

                    </div>

                </div>
            </template>
        </div>
        <div class="review-operation">
            <div style="flex: 1;"></div>
            <el-button type="success" v-if="(datas.candidate || '').split(',').indexOf(userId + '') != -1"
                @click="approve(true)">{{
                    $L('同意') }}</el-button>
            <el-button type="danger" v-if="(datas.candidate || '').split(',').indexOf(userId + '') != -1"
                @click="approve(false)">{{
                    $L('拒绝') }}</el-button>
            <el-button type="warning" v-if="isShowWarningBtn" @click="revocation">{{ $L('撤销') }}</el-button>
            <el-button @click="comment" type="success" ghost>+{{ $L('添加评论') }}</el-button>
        </div>
    </div>
</template>

<script setup>
import { ref, onMounted, watch, nextTick, computed, getCurrentInstance, defineEmits } from 'vue';
import { getProcessData, agreeOrRefuse, revocationTask ,getUserInfo } from "@/plugins/api.js";
import { useRoute, useRouter } from 'vue-router'
const { proxy } = getCurrentInstance()
import { ElMessage, ElMessageBox } from 'element-plus';
import 'element-plus/es/components/message/style/css'; // 
import 'element-plus/es/components/message-box/style/css';
let modalTransferIndex = ref(window.modalTransferIndex)
let datas = ref({})
let showTimeNum = ref(24)
let userId = ref(0)

let props = defineProps({
    data: {
        type: Object,
        default: {},
    },
})


watch(() => props.data, (newValue) => {
    if (newValue.id) {
        const queryString = window.location.search.slice(1);
    const params = {};
    const pairs = queryString.split('&');
    for (let i = 0; i < pairs.length; i++) {
        const pair = pairs[i].split('=');
        params[pair[0]] = decodeURIComponent(pair[1] || '');
    }
    if (params.token) {
        getUserData(params.token)
    }
  
    }
})

const isShowWarningBtn = computed(() => {
    let is = userId.value == datas.value.startUserId;
    (datas.value.nodeInfos || []).map(h => {
        if (h.type != 'starter' && h.isFinished == true && h.identitylink?.userid != userId.value) {
            is = false;
        }
    })
    return is;
})

const router = useRouter()
const route = useRoute()
onMounted(() => {
    modalTransferIndex.value = window.modalTransferIndex = window.modalTransferIndex + 1
});

// 获取用户信息
const getUserData = async (token) => {
    const upToken = token;
    let { data, status, message } = await getUserInfo(upToken)
    if (status != 200) {
        ElMessage.error(message)
        return;
    }
    if (data) {
        userId.value = data.userid;
        getInfo();
    }
}


// 把时间转成几小时前
const getTimeAgo = (time, type) => {
    const currentTime = new Date();
    const timeDiff = (currentTime - new Date((time + '').replace(/-/g, "/"))) / 1000; // convert to seconds
    if (timeDiff < 60) {
        return type == 2 ? "0" + proxy.$L('分钟') : proxy.$L('刚刚');
    } else if (timeDiff < 3600) {
        const minutes = Math.floor(timeDiff / 60);
        return type == 2 ? `${minutes}${proxy.$L('分钟')}` : `${minutes} ${proxy.$L('分钟前')}`;
    } else if (timeDiff < 3600 * 24) {
        const hours = Math.floor(timeDiff / 3600);
        return type == 2 ? `${hours}${proxy.$L('小时')}` : `${hours} ${proxy.$L('小时前')}`;
    } else {
        const days = Math.floor(timeDiff / 3600 / 24);
        return type == 2 ? `${days + 1}${proxy.$L('天')}` : `${days + 1} ${proxy.$L('天')}`;
    }
}
// 获取时间差
const getTimeDifference = (startTime, endTime) => {
    const currentTime = new Date((endTime + '').replace(/-/g, "/"));
    const timeDiff = (currentTime - new Date((startTime + '').replace(/-/g, "/"))) / 1000; // convert to seconds
    if (timeDiff < 60) {
        return { time: timeDiff, unit: proxy.$L('秒') };
    } else if (timeDiff < 3600) {
        const minutes = Math.floor(timeDiff / 60);
        return { time: minutes, unit: proxy.$L('分钟') };
    } else if (timeDiff < 3600 * 24) {
        const hours = Math.floor(timeDiff / 3600);
        return { time: hours, unit: proxy.$L('小时') };
    } else {
        const days = Math.floor(timeDiff / 3600 / 24);
        return { time: days + 1, unit: proxy.$L('天') };
    }
}
// 获取详情
const getInfo = async () => {
    datas.value = props.data
    const id = datas.value.id;
    let { data, status, message } = await getProcessData(id)
    if (status != 200) {
        ElMessage.error(message)
        return;
    }
    if (data) {
        let show = true;
        data.nodeInfos = (data.nodeInfos || []).map(item => {
            item._show = show;
            if (item.identitylink?.state == 2 || item.identitylink?.state == 3) {
                show = false;
            }
            return item;
        })
        datas.value = data
    }
}
// 通过
const approve = (type) => {
    ElMessageBox.prompt(proxy.$L('请输入审批意见'), proxy.$L('审批'), {
        confirmButtonText: type ? proxy.$L('同意') : proxy.$L('拒绝'),
        cancelButtonText: proxy.$L('取消'),
        inputValidator: (val) => {
            if (val === null || val.length < 1) {
                return false;
            }
        },
        inputErrorMessage: proxy.$L('请输入审批意见'),
    })
        .then(({ value }) => {
            const upData = {
                userID: userId.value.toString(),
                taskID: datas.value.taskID,
                pass: type ? 'true' : 'false',
                comment: value,
            };
            audit(upData);
        })
        .catch(() => {

        })
}
const emit = defineEmits(['approve', 'revocation'])
const audit = async (upData) => {
    let { data, status, message } = await agreeOrRefuse(upData)
    if (status != 200) {
        ElMessage.error(message)
        return;
    }
    if (status == 200) {
        emit('approve')
        ElMessage.success(message);
        getInfo()
    }
}

// 撤销
const revocation = () => {
    ElMessageBox.confirm(
        proxy.$L('你确定要撤销吗？'),
        proxy.$L('撤销'),
        {
            confirmButtonText: proxy.$L('确定'),
            cancelButtonText: proxy.$L('取消'),
            type: 'warning',
        }
    )
        .then(() => {
            regret();
        })
        .catch(() => {

        })

}

const regret = async () => {
    const upData = {
        userID: userId.value.toString(),
        taskID: datas.value.taskID,
        procInstId: datas.value.id
    }
    let { data, status, message } = await revocationTask(upData)
    if (status != 200) {
        ElMessage.error(message)
        return;
    }
    if (status == 200) {
        emit('revocation')
        ElMessage.success(message);
        getInfo()
    }
}


// 获取内容
const getContent = (content) => {
    try {
        return JSON.parse(content).content || ''
    } catch (error) {
        return ''
    }
}

// 获取内容
const getPictures = (content) => {
    try {
        return JSON.parse(content).pictures || []
    } catch (error) {
        return ''
    }
}

// 打开图片
const onViewPicture = (currentUrl) => {

}

</script>

<style  lang="less">
@import "../../../css/review.less";
</style>
