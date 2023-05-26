<template>
    <el-upload v-model:file-list="fileList"  :multiple="true" :action="Url" list-type="picture-card" :limit="5" :on-exceed="handleExceed" :on-preview="handlePictureCardPreview"
        :on-remove="handleRemove" :on-success="handleUpData">
        <el-icon >
            <Plus />
        </el-icon>
    </el-upload>

    <el-dialog v-model="dialogVisible">
        <img w-full :src="dialogImageUrl"  />
    </el-dialog>
</template>
  
<script lang="ts" setup>
import { ref , getCurrentInstance} from 'vue'
import { Plus } from '@element-plus/icons-vue'
import type { UploadProps, UploadUserFile } from 'element-plus'
import { ElMessage } from 'element-plus';
const { proxy } = getCurrentInstance()
let emits = defineEmits(['update:img'])

const Url = ref<string>(`${window.location.origin}/api/v1/workflow/dootask/upload`)
//获取请假类型


const fileList = ref<UploadUserFile[]>([
])

const dialogImageUrl = ref('')
const dialogVisible = ref(false)

const handleRemove: UploadProps['onRemove'] = (uploadFile, uploadFiles) => {
    emits("update:img", fileList)
}

const handleUpData: UploadProps['onSuccess'] = () => {
    emits("update:img", fileList)
}

const handleExceed : UploadProps['onExceed'] = () => {
    ElMessage.error(proxy.$L("超出上传限制！"))
}

const handlePictureCardPreview: UploadProps['onPreview'] = (uploadFile) => {
    dialogImageUrl.value = uploadFile.url!
    if(window.parent){
        parent.postMessage({
            "action":'openvalie',
            "url":dialogImageUrl.value
        },'*')
    }else{
        dialogVisible.value = true
    }
    
}
</script>
<style lang="less" scoped>
/deep/ .el-upload-list__item, /deep/ .el-upload{
    width: 80px;
    height: 80px;
}
/deep/ .el-dialog{
    min-width: 90% !important;
}
</style>