<template>
    <div>

        <div class="fd-nav">
            <div class="fd-nav-left">
                <!-- <div class="fd-nav-back" @click="toReturn">
                    <i class="anticon anticon-left"></i>
                </div> -->
                <div class="fd-nav-title">{{ workFlowDef.name || '' }}</div>
            </div>
            <div class="fd-nav-right">
                <button type="button" class="ant-btn button-publish" @click="saveSet">
                    <span>发 布</span>
                </button>
            </div>
        </div>
        
        <div class="fd-nav-content" style="top: 60px;">
            <section class="dingflow-design">
                <div class="zoom">
                    <div class="zoom-out" :class="nowVal == 50 && 'disabled'" @click="zoomSize(1)"></div>
                    <span>{{ nowVal }}%</span>
                    <div class="zoom-in" :class="nowVal == 300 && 'disabled'" @click="zoomSize(2)"></div>
                </div>
                <div class="box-scale" :style="`transform: scale(${nowVal / 100});`">
                    <nodeWrap v-model:nodeConfig="nodeConfig" v-model:flowPermission="flowPermission" />
                    <div class="end-node">
                        <div class="end-node-circle"></div>
                        <div class="end-node-text">流程结束</div>
                    </div>
                </div>
            </section>
        </div>
        
        <errorDialog v-model:visible="tipVisible" :list="tipList" />
        <promoterDrawer />
        <approverDrawer :directorMaxLevel="directorMaxLevel" />
        <copyerDrawer />
        <conditionDrawer />
    </div>
</template>

<script setup>
import errorDialog from "@/components/dialog/errorDialog.vue";
import promoterDrawer from "@/components/drawer/promoterDrawer.vue";
import approverDrawer from "@/components/drawer/approverDrawer.vue";
import copyerDrawer from "@/components/drawer/copyerDrawer.vue";
import conditionDrawer from "@/components/drawer/conditionDrawer.vue";
import { ElMessage } from 'element-plus'
import { ref, onMounted } from "vue";
import { useRoute } from 'vue-router'
import { getWorkFlowData, setWorkFlowData } from "@/plugins/api.js";
import { mapMutations } from "@/plugins/lib.js";
let { setTableId, setIsTried } = mapMutations()

let tipList = ref([]);
let tipVisible = ref(false);
let nowVal = ref(100);
let processConfig = ref({});
let nodeConfig = ref({});
let workFlowDef = ref({});
let flowPermission = ref([]);
let directorMaxLevel = ref(0);

onMounted(async () => {
    let route = useRoute()
    let {data,status,message} = await getWorkFlowData({ id: route.query.workFlowDefId })
    // 
    processConfig.value = data || {};
    nodeConfig.value = data?.resource || {
        "name": "发起人",
        "type": "start",
        "nodeId": "sid-startevent",
        "childNode": {},
        "conditionNodes": []
    };
    flowPermission.value = [];
    directorMaxLevel.value = 4; //最高级别
    workFlowDef.value = processConfig.value;
    setTableId(data?.id || 0);
});

const toReturn = () => {
    //window.location.href = ""
};

const reErr = ({ childNode }) => {
    if (childNode) {
        let { type, error, nodeName, conditionNodes } = childNode;
        if (type == 1 || type == 2) {
            if (error) {
                tipList.value.push({
                    name: nodeName,
                    type: ["", "审核人", "抄送人"][type],
                });
            }
            reErr(childNode);
        } else if (type == 3) {
            reErr(childNode);
        } else if (type == 4) {
            reErr(childNode);
            for (var i = 0; i < conditionNodes.length; i++) {
                if (conditionNodes[i].error) {
                    tipList.value.push({ name: conditionNodes[i].nodeName, type: "条件" });
                }
                reErr(conditionNodes[i]);
            }
        }
    } else {
        childNode = null;
    }
};

const saveSet = async () => {
    setIsTried(true);
    tipList.value = [];
    reErr(nodeConfig);
    if (tipList.value.length != 0) {
        tipVisible.value = true;
        return;
    }
    processConfig.value.flowPermission = flowPermission.value;
    processConfig.value.resource = nodeConfig.value;
    // eslint-disable-next-line no-console
    console.log(nodeConfig.value);
    let res = await setWorkFlowData(processConfig.value);
    if (res.status == 200) {
        ElMessage.success("设置成功")
        setTimeout(function () {
            window.location.href = "";
        }, 200);
    }else{
        ElMessage.error(res.message)
    }
};

const zoomSize = (type) => {
    if (type == 1) {
        if (nowVal.value == 50) {
            return;
        }
        nowVal.value -= 10;
    } else {
        if (nowVal.value == 300) {
            return;
        }
        nowVal.value += 10;
    }
};

</script>
<style>
@import "../css/workflow.css";

.error-modal-list {
    width: 455px;
}
</style>