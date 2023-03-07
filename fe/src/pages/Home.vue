<script setup lang="ts">
import SearchInput from "@/components/SearchInput.vue";
import VerifyCert from "@/components/VerifyCert.vue";
import { TITLE, SEARCH_PLACE_HOLDER } from "@/common/constants";
import { ArrowRightOutlined } from "@ant-design/icons-vue";
</script>

<template>
  <div class="background"></div>
  <div class="content">
    <div class="main">
      <div class="title">{{ TITLE }}</div>
      <SearchInput :placeholder="SEARCH_PLACE_HOLDER" :dark="true"/>
      <div class="ops">
        <router-link to="/create">
          创建证书
        </router-link>
        <VerifyCert />
      </div>
      <div class="link">
        <router-link to="/certs">
          浏览所有证书
          <arrow-right-outlined />
        </router-link>
      </div>
    </div>
  </div>
</template>

<style scoped>
.content {
  margin-top: 18%;
  padding: 20px 0;
  backdrop-filter: blur(1px);
  background: radial-gradient(#00000044, transparent);
}

.main {
  width: 520px;
  margin: auto;
}

.main .title {
  font-size: 32px;
  font-weight: bold;
  margin-bottom: 54px;
  color: white;
}

.main .ops {
  gap: 22px;
  display: flex;
  font-size: 16px;
  margin-top: 22px;
  line-height: 40px;
  justify-content: center;
}

.main .link {
  margin-top: 32px;
  font-size: 18px;
  font-weight: bold;
}

.background {
  position: absolute;
  top: 0;
  bottom: 0;
  left: 0;
  background-image: url(/cert/src/assets/bg.jpg);
  right: 0;
  z-index: -1;
  opacity: 0.95;
  background-position: center;
  background-repeat: no-repeat;
  background-size: cover;
}
</style>
