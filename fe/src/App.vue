<template>
  <div id="app">
    <router-view></router-view>
  </div>
</template>

<script>
export default {
  name: "App",

  components: {},

  data() {
    return {};
  },

  computed: {},

  methods: {},
};
</script>

<style scoped>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  min-width: 900px;
  margin: auto;
  padding: 16px;
}

.el-header {
  top: 0;
  z-index: 100;
  position: sticky;
  padding: 0 !important;
  background-color: var(--bg-color-0, white);
  transition: ease 0.6s all;
  box-shadow: 0 4px 12px #00000014, 0 0 1px #0100001a;
}
</style>
