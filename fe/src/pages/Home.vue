<script setup lang="ts">
import SearchInput from "@/components/SearchInput.vue";
import VerifyCert from "@/components/VerifyCert.vue";
import CertHistory from "@/components/CertHistory.vue";
import OnChainCert from "@/components/OnChainCert.vue";
import { onSearch } from "@/composition/useSearch";
import { TITLE, SEARCH_PLACE_HOLDER } from "@/common/constants";
import bgPath from '@/assets/bg.jpg'
</script>

<template>
  <div class="background" :style="{ backgroundImage: `url(${bgPath})` }"></div>
  <div class="container">
    <div class="content">
      <div class="main">
        <div class="title">{{ TITLE }}</div>
        <SearchInput :placeholder="SEARCH_PLACE_HOLDER" :dark="true" @search="onSearch"/>
        <div class="ops">
          <router-link to="/create">
            创建证书
          </router-link>
          <VerifyCert />
          <router-link to="/certs">
            浏览所有证书
          </router-link>
          <CertHistory />
        </div>
        <div class="link">
          <OnChainCert />
          <router-link to="/blacklist">
            查看黑名单
          </router-link>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.container {
    position: absolute;
    left: 0;
    right: 0;
    top: 0;
    bottom: 0;
    display: flex;
    flex-direction: column;
    justify-content: center;
    text-align: center;
}

.content {
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
  margin-top: 32px;
  line-height: 40px;
  justify-content: center;
}

.main .link {
  margin-top: 22px;
  gap: 22px;
  display: flex;
  font-size: 16px;
  line-height: 40px;
  justify-content: center;
}

.background {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  z-index: -1;
  opacity: 0.95;
  background-position: center;
  background-repeat: no-repeat;
  background-size: cover;
}
</style>
