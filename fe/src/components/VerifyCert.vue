<script lang="ts" setup>
import { verify } from "@/api/cert";
import { readLocalFile } from "@/utils/files";
import { message } from "ant-design-vue";

function verifyCert() {
  readLocalFile(
    (_, content) => {
      verify(content)
        .then((res) => {
          if (res === "True") {
            message.success("验证通过");
          }
        })
        .catch(() => {
          message.error("未能通过验证");
        });
    },
    (err) => {
      console.log(err);
      message.error(err);
    }
  );
}
</script>

<template>
  <a-button class="link-btn" type="link" size="large" @click="verifyCert"> 验证本地证书 </a-button>
</template>

<style scoped>
.link-btn {
  padding: 0;
}
</style>
