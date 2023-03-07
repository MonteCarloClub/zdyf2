<script lang="ts" setup>
import { ref } from "vue";
import { message } from "ant-design-vue";
import { apply } from "@/api/cert";
import { download } from "@/utils/files";

const success = ref(false);
const placeholder = "请输入 UID";
const certApplyUID = ref("");

const successCert = ref<API.Cert>();
let content = "";

function createCert() {
  const uid = certApplyUID.value;
  if (uid === undefined || uid === "") {
    message.error("UID 不能为空");
    return;
  }

  apply({ uid })
    .then((res) => {
      success.value = true;
      successCert.value = res.certificate;
      content = JSON.stringify(res)
    })
    .catch((err) => {
      message.error(err);
    });
}

function reCreate() {
  success.value = false;
  certApplyUID.value = "";
  successCert.value = undefined;
}

function downloadCert() {
  const name = `${successCert.value?.serialNumber}.cert`
  message.info(`开始下载 ${name}`)
  download(content, name);
}
</script>

<template>
  <div class="content">
    <a-result
      v-if="success"
      status="success"
      title="已成功创建证书!"
      :sub-title="`证书编号为: ${successCert?.serialNumber}`"
    >
      <template #extra>
        <a-button key="console" type="primary" @click="downloadCert">下载该证书</a-button>
        <a-button key="buy" @click="reCreate">再次创建</a-button>
      </template>
    </a-result>
    <div v-else>
      <div class="title">创建证书</div>
      <div class="input-container">
        <input
          type="text"
          @keydown.enter="createCert"
          v-model="certApplyUID"
          :placeholder="placeholder"
        />
      </div>
      <a-button class="btn" block size="large" @click="createCert">创建</a-button>
    </div>
  </div>
</template>

<style scoped>
.content {
  margin-top: 8%;
}

.title {
  font-size: 26px;
  margin: 0 auto 64px;
}

.btn {
  margin-top: 48px;
  width: 100px;
}

input {
  border-style: none;
  background: transparent;
  outline: none;
}

.input-container {
  width: 100%;
  margin: 0 auto;
  max-width: 520px;
  position: relative;
  border-radius: 2px;
  padding: 12px 16px;
  background: #eeeeee55;
}

.input-container input {
  width: 100%;
  font-size: 20px;
  text-align: center;
  vertical-align: middle;
}
.input-container input::-webkit-input-placeholder {
  color: #7881a1;
}

.input-container:after {
  content: "";
  position: absolute;
  left: 0px;
  right: 0px;
  bottom: 0px;
  height: 2px;
  z-index: 999;
  border-bottom-left-radius: 2px;
  border-bottom-right-radius: 2px;
  background-position: 0% 0%;
  background: linear-gradient(
    to right,
    #b294ff,
    #57e6e6,
    #feffb8,
    #57e6e6,
    #b294ff,
    #57e6e6
  );
  background-size: 500% auto;
  -webkit-animation: gradient 3s linear infinite;
  animation: gradient 3s linear infinite;
}

@-webkit-keyframes gradient {
  0% {
    background-position: 0 0;
  }
  100% {
    background-position: 100% 0;
  }
}

@keyframes gradient {
  0% {
    background-position: 0 0;
  }
  100% {
    background-position: 100% 0;
  }
}
</style>
