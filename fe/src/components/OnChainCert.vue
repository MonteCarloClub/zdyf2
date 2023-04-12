<script setup lang="ts">
import { ref } from "vue";
import { message } from "ant-design-vue";
import CertDetail from "@/components/CertDetail.vue";
import SearchInput from "@/components/SearchInput.vue";
import { statusOnChain } from "@/api/cert";
import { SEARCH_PLACE_HOLDER } from "@/common/constants";
import { isObject } from "@/utils/is";

const modalVisible = ref(false);

const cert = ref<API.Cert>();
const inputKey = ref('')
const first = ref(false)
const loading = ref(false);
const hasResult = ref(false);
const revokedTime = ref('');

const searchOnChain = (searchKey: string) => {
    const s = searchKey;
    if (s.length < 1) {
        message.error('请输入证书编号')
        return;
    }

    first.value = true;
    inputKey.value = s;

    loading.value = true;
    statusOnChain({
        no: searchKey
    }).then((res) => {
        hasResult.value = true;
        // 如果未撤销，是 object
        if (isObject(res)) {
            revokedTime.value = ''
            cert.value = res
            return;
        }
        // 如果是 string
        const regexp = /\[Revoked at (.*)\] (.*)/g;
        const matches = Array.from(res.matchAll(regexp));
        if (matches.length > 0 && matches[0].length === 3) {
            // 撤销时间
            revokedTime.value = matches[0][1];
            // 证书信息
            cert.value = JSON.parse(JSON.parse(matches[0][2]))
        }
        else {
            hasResult.value = false;
        }
    }).catch(() => {
        hasResult.value = false;
    }).finally(() => {
        loading.value = false;
    })
}
</script>

<template>
    <a-modal v-model:visible="modalVisible" centered :footer="null" title="查询证书的链上信息" @ok="modalVisible = false">
        <div class="panel">
            <SearchInput :placeholder="SEARCH_PLACE_HOLDER" @search="searchOnChain" />
            <div v-if="first">
                <div v-if="loading" class="loading-area">
                    <a-spin tip="正在检索" />
                </div>
                <div v-else>
                    <a-card v-if="hasResult" class="cert-card" title="证书详情">
                        <template #extra>
                            <a-tooltip v-if="revokedTime" placement="left" color="#f50">
                                <template #title>
                                    <span> 撤销时间为：{{ revokedTime }} </span>
                                </template>
                                <a-tag color="#f50">已撤销</a-tag>
                            </a-tooltip>
                            <a-tag v-else color="#87d068">有效</a-tag>
                        </template>
                        <CertDetail :cert="cert" />
                    </a-card>
                    <div v-else>
                        <a-empty style="margin-top: 100px">
                            <template #description>
                                <span>未搜索到编号为</span>
                                <a-tag> {{ inputKey }} </a-tag>
                                <span>的证书</span>
                            </template>
                        </a-empty>
                    </div>
                </div>
            </div>
        </div>
    </a-modal>
    <a-button class="link-btn" type="link" size="large" @click="modalVisible = true"> 查询链上信息 </a-button>
</template>

<style scoped>
.link-btn {
    padding: 0;
}

.panel {
    margin: 0 auto;
}

.cert-card {
    margin-top: 22px;
    text-align: left;
}

.panel .ant-card-meta {
    margin-bottom: 22px;
}

.loading-area {
    margin-top: 120px;
    text-align: center;
}
</style>
