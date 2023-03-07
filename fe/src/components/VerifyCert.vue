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
        .catch(console.log);
    },
    (err) => {
      console.log(err);
      message.error(err);
    }
  );
}
</script>

<template>
  <a-button @click="verifyCert"> 验证本地证书 </a-button>
</template>
