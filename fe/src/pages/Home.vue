<script setup lang="ts">
import { ref } from "vue";
import { useStore } from "@/store";

const count = ref(0);

const store = useStore();
const increment = () => {
  count.value++;
  store.commit("increment");
};

</script>

<template>
  <p>
    Recommended IDE setup:
    <a href="https://code.visualstudio.com/" target="_blank">VSCode</a>
    +
    <a href="https://github.com/johnsoncodehk/volar" target="_blank">Volar</a>
  </p>

  <p>
    <a href="https://vitejs.dev/guide/features.html" target="_blank"
      >Vite Docs</a
    >
    |
    <a href="https://v3.vuejs.org/" target="_blank">Vue 3 Docs</a>
  </p>

  <a-button type="primary" @click="increment"
    >click in this page: {{ count }}</a-button
  >
</template>

<style scoped>
a {
  color: #42b983;
}

label {
  margin: 0 0.5em;
  font-weight: bold;
}

code {
  background-color: #eee;
  padding: 2px 4px;
  border-radius: 4px;
  color: #304455;
}
</style>
