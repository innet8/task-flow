<template>
    <div class="review-list">
        <h2>
            <span class="list-name">{{$L(data.procDefName)}}</span> 
            <el-tag v-if="data.state == 0" type="warning">{{$L('待审批')}}</el-tag>
            <el-tag v-if="data.state == 1" type="warning">{{$L('审批中')}}</el-tag>
            <el-tag v-if="data.state == 2" type="success">{{$L('已通过')}}</el-tag>
            <el-tag v-if="data.state == 3" type="danger">{{$L('已拒绝')}}</el-tag>
            <el-tag v-if="data.state == 4" type="danger">{{$L('已撤回')}}</el-tag>
        </h2>
        <p v-if="data.var?.type">{{$L('假期类型')}}：<span>{{$L(data.var?.type)}}</span></p>
        <p>{{$L('开始时间')}}：<span>{{data.var?.startTime}}</span></p>
        <p>{{$L('结束时间')}}：<span>{{data.var?.endTime}}</span></p>
        <div class="list-member">
            <span>
                <el-avatar :src="data.userimg" :size="20"/>
                {{ data.startUserName }}
            </span>
            <span>
                {{$L('发起时间')}}：{{data.startTime}}
            </span>
        </div>
    </div>
</template>

<script setup>

let props = defineProps({
    data:{
        type: Object,
        default: {},
    },
})


</script>

