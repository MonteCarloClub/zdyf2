<template>
  <div class="group">
    <input
      type="text"
      required
      v-model="searchKey"
      @keydown.enter="onSearch(searchKey)"
      :placeholder="dynamicPlaceholder ? '' : placeholder"
    />
    <span class="bar"></span>
    <a-button
      type="link"
      class="search-icon"
      size="large"
      @click="onSearch(searchKey)"
    >
      <template #icon>
        <SearchOutlined />
      </template>
    </a-button>
    <label v-if="dynamicPlaceholder">{{ placeholder }}</label>
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { onSearch } from "@/composition/useSearch";
import { SearchOutlined } from "@ant-design/icons-vue";

defineProps({
  placeholder: {
    type: String,
    default: "请在此输入",
  },
  dynamicPlaceholder: {
    type: Boolean,
    default: true,
  },
});

const searchKey = ref("");
</script>

<style scoped>
.group {
  position: relative;
}

input {
  width: 100%;
  border: none;
  display: block;
  font-size: 18px;
  padding: 10px 10px 10px 0;
  background-color: transparent;
  border-bottom: 1px solid #757575;
}

input:focus {
  outline: none;
}

label {
  left: 0;
  top: 10px;
  color: #999;
  font-size: 18px;
  position: absolute;
  font-weight: normal;
  pointer-events: none;
  transition: 0.2s ease all;
  -moz-transition: 0.2s ease all;
  -webkit-transition: 0.2s ease all;
}

/* active state */
input:focus ~ label,
input:valid ~ label {
  top: -20px;
  font-size: 14px;
  color: #646464;
}

.bar {
  width: 100%;
  display: block;
  position: relative;
}

.bar:before {
  left: 0%;
  width: 0;
  bottom: 1px;
  content: "";
  height: 2px;
  position: absolute;
  background: #646464;
  transition: 0.2s ease all;
  -moz-transition: 0.2s ease all;
  -webkit-transition: 0.2s ease all;
}

/* active state */
input:focus ~ .bar:before {
  width: 100%;
}

.search-icon {
  position: absolute;
  right: 0;
  bottom: 3px;
  color: black;
}

.search-icon:hover {
  color: #1890ff;
}
</style>
