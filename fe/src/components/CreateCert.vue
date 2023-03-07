<script lang="ts" setup>
import { ref } from "vue";
import { message, Modal } from "ant-design-vue";
import { apply } from "@/api/cert";

const createFormVisible = ref(false);
const certApplyUID = ref("");
const placeholder = "请输入 UID";

function createCert() {
  const uid = certApplyUID.value;
  
  if (uid === undefined || uid === "") {
    message.error("UID 不能为空");
    return;
  }
  apply({ uid })
    .then((res) => {
      Modal.success({
        title: "创建成功",
        content: `证书编号为：${res.certificate.serialNumber}`,
        okText: "好的",
        onOk() {
          createFormVisible.value = false;
        },
      });
    })
    .catch((err) => {
      message.error(err);
    });
}
</script>

<template>
  <a-button @click="createFormVisible = true"> 创建证书 </a-button>
  <a-modal v-model:visible="createFormVisible" title="创建">
    <a-form
      ref="formRef"
      name="custom-validation"
      :label-col="{ span: 6 }"
      :wrapper-col="{ span: 18 }"
      labelAlign="left"
    >
      <a-form-item label="UID">
        <a-input
          v-model:value="certApplyUID"
          :placeholder="placeholder"
        ></a-input>
      </a-form-item>
    </a-form>
    <template #footer>
      <a-button key="submit" type="primary" @click="createCert">
        确认创建
      </a-button>
    </template>
  </a-modal>
</template>
