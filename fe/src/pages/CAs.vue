<template>
    <div class="nav-bar">
        <div style="flex: 1">
            <router-link to="/"> 首页 </router-link>
        </div>
        <div> CA信誉列表 </div>
        <div style="flex: 1"></div>
    </div>

    <div v-if="cas.length > 0" class="card-container">
        <div class="card" v-for="ca in cas" :key="ca.name">
            <a-avatar size="large">CA</a-avatar>
            <div class="info">
                <div class="name">{{ ca.name }}</div>
                <div class="desc">{{ ca.score }}</div>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from "vue";
import { score as getCAScore } from "@/api/node";
import { parallelWithLimit } from "@/utils/promises";

const CANames:string[] = [];

for (let index = 0; index < 10; index++) {
    CANames.push('CA-' + index);
}

interface CA {
    name: string;
    score: number;
}
const cas = ref<CA[]>([]);

function updateCAs() {
    // parallelWithLimit(CANames.map(id => getCAScore({
    //     id
    // })), 4, (index, res) => {
    //     console.log(index, res);
    // })
}

let loopTimer = -1;

onMounted(() => {
    loopTimer = window.setInterval(updateCAs, 1000)
})

onUnmounted(() => {
    window.clearInterval(loopTimer)
})
</script>

<style>
.card-container {
    margin: 0 32px;
    display: grid;
    grid-row-gap: 15px;
    grid-column-gap: 15px;
    grid-template-columns: repeat(4, 1fr);
    align-items: start;
    justify-items: stretch;
}

@media (max-width: 1281px) {
    .card-container {
        grid-template-columns: repeat(3, 1fr);
    }
}

@media (max-width: 1023px) {
    .card-container {
        grid-template-columns: repeat(2, 1fr);
    }
}

@media (max-width: 720px) {
    .card-container {
        grid-template-columns: 100%;
    }
}
</style>

<style scoped>
.card-container .card {
    background-color: #fafafa;
    border-radius: 20px;
    cursor: pointer;
    display: flex;
    padding: 25px;
}

.card-container .card:hover {
    background-color: #eaeaea;
}

.card-container .card .info {
    margin-left: 20px;
    text-align: left;
    display: flex;
    flex-direction: column;
}

.card-container .card .name {
    font-weight: bold;
}

.nav-bar {
    display: flex;
    text-align: left;
    font-size: 22px;
    margin: 0 32px 22px;
}

.nav-bar a {
    color: black;
}
</style>