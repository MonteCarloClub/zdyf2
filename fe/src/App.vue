<template>
  <router-view></router-view>
</template>

<script setup lang="ts">
import { TITLE } from "@/common/constants";
document.title = TITLE;
</script>

<style>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  padding: 20px;
}
</style>
