<script setup lang="ts">
import { ref, watch } from "vue";
import { useRoute } from "vue-router";
import { DownloadOutlined, DeleteOutlined } from "@ant-design/icons-vue";
import { message } from "ant-design-vue";
import CertDetail from "@/components/CertDetail.vue";
import SearchInput from "@/components/SearchInput.vue";
import { SEARCH_PLACE_HOLDER } from "@/common/constants";
import { query, revoke } from "@/api/cert";
import { download } from "@/utils/files";

const route = useRoute();

/**相同路由下，地址改变需要重新获取地址信息 */
watch(
  () => route.params.no,
  async (no) => {
    if (no === undefined) {
      return;
    }
    init(typeof no === "string" ? no : no[0]);
  }
);

const certDetail = ref<API.Cert>();
const loading = ref(true);
const hasResult = ref(false);
let responseText = "";

function init(no: string) {
  query({
    no,
  })
    .then((res) => {
      certDetail.value = res.certificate;
      responseText = JSON.stringify(res);
      hasResult.value = true;
    })
    .catch((err) => {
      hasResult.value = false;
    })
    .finally(() => {
      loading.value = false;
    });
}

function downloadCert() {
  const fileName = `${certDetail.value?.serialNumber}.cert`;
  message.info(`开始下载${fileName}`);
  download(responseText, fileName);
}

const no = ref(route.params.no);
init(typeof no.value === "string" ? no.value : no.value[0]);

const revoking = ref(false);

function revokeCert() {
  const no = certDetail.value?.serialNumber;
  if (no === undefined) {
    return
  }
  revoking.value = true;
  revoke({ no })
    .then((res) => {
      if (res === "Revoke OK.") {
        message.success("撤销成功");
      }
    })
    .catch((err) => {
      message.error(err);
    })
    .finally(() => {
      revoking.value = false;
    });
}
</script>

<template>
  <div class="panel">
    <SearchInput :dynamic-placeholder="false" :placeholder="SEARCH_PLACE_HOLDER"/>
  </div>
  <div v-if="loading" class="loading-area">
    <a-spin tip="正在检索"/>
  </div>
  <div v-else>
    <a-card v-if="hasResult" hoverable class="panel">
      <a-card-meta title="证书信息" />
      <CertDetail :cert="certDetail" />
      <template #actions>
        <a-button type="link" @click="revokeCert" :loading="revoking" danger>
          <template #icon>
            <DeleteOutlined />
          </template>
          撤销
        </a-button>
        <a-button type="link" @click="downloadCert">
          <template #icon>
            <DownloadOutlined />
          </template>
          下载完整版证书
        </a-button>
      </template>
    </a-card>
    <div v-else>
      <a-empty style="margin-top: 200px">
        <template #description>
          <span>未搜索到编号为</span>
          <a-tag> {{ no }} </a-tag>
          <span>的证书</span>
        </template>
      </a-empty>
    </div>
  </div>
</template>

<style scoped>
.panel {
  text-align: left;
  width: 600px;
  margin: 0 auto 22px;
}

.panel .ant-card-meta {
  margin-bottom: 22px;
}

.loading-area {
  margin-top: 120px;
}
</style>
