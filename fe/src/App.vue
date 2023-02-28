<template>
  <div>
    <router-link to="/">home</router-link>|
    <router-link to="/hello">hello</router-link>
  </div>
  <img alt="Vue logo" src="./assets/logo.png" />
  <div>{{ count }}</div>
  <router-view></router-view>
</template>

<script setup lang="ts">
import { useStore } from "@/store";
import { computed } from 'vue'

const store = useStore();
const count = computed(() => store.state.count)
</script>

<style>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  margin-top: 60px;
}
</style>
