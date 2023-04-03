<script setup lang="ts">
import { ref, onMounted } from "vue";
import { history } from "@/api/cert";

const modalVisible = ref(false);
const loadingMore = ref(false);
const hasMore = ref(true);

const list = ref<string[]>([])

let index = 0;
const count = 10;
function fetchNext() {
    history({
        count,
        index
    }).then((res) => {
        const { certificates } = res;
        list.value = list.value.concat(certificates);
        index += count;
        // 不够的时候，认为后续没有了
        if (certificates.length < count) {
            hasMore.value = false;
        }
    }).catch(console.log)
}

onMounted(() => {
    fetchNext();
});
</script>

<template>
    <a-modal v-model:visible="modalVisible" centered :footer="null" title="证书颁发记录" @ok="modalVisible = false">
        <div class="scrollable">
            <p v-for="cert in list" :key="cert">{{ cert }}</p>
            <div style="text-align: center;">
                <a-button v-if="hasMore" type="link" :loading="loadingMore" @click="fetchNext">获取更多</a-button>
                <p v-else style="color: gray">没有更多了</p>
            </div>
        </div>
    </a-modal>
    <a-button class="link-btn" type="link" size="large" @click="modalVisible = true"> 历史颁发记录 </a-button>
</template>

<style scoped>
.link-btn {
    padding: 0;
}

.scrollable {
    height: 300px;
    overflow-y: auto;
}

.scrollable::-webkit-scrollbar {
    width: 5px;
    height: 5px;
}

.scrollable::-webkit-scrollbar-thumb {
    background-color: lightgray;
    border-radius: 2px;
}

.scrollable::-webkit-scrollbar-track {
    background-color: transparent;
}

.scrollable-hidden {
    overflow: auto;
}

.scrollable-hidden::-webkit-scrollbar {
    width: 0px;
    height: 0px;
}

.scrollable-hidden::-webkit-scrollbar-track {
    background-color: transparent;
}
</style>
