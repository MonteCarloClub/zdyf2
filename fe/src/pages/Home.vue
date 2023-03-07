<script setup lang="ts">
import SearchInput from "@/components/SearchInput.vue";
import CreateCert from "@/components/CreateCert.vue";
import VerifyCert from "@/components/VerifyCert.vue";
import { TITLE, SEARCH_PLACE_HOLDER } from "@/common/constants";
import { ArrowRightOutlined } from "@ant-design/icons-vue";
</script>

<template>
  <div class="background"></div>
  <div class="main">
    <div class="title">{{ TITLE }}</div>
    <SearchInput :placeholder="SEARCH_PLACE_HOLDER" />
    <div class="ops">
      <CreateCert />
      <VerifyCert />
    </div>
    <div class="link">
      <router-link to="/certs"
        >浏览所有证书
        <arrow-right-outlined />
      </router-link>
    </div>
  </div>
</template>

<style scoped>
.main {
  width: 520px;
  margin: 20% auto 0;
}

.main .title {
  font-size: 32px;
  font-weight: bold;
  margin-bottom: 54px;
}

.main .ops {
  margin-top: 22px;
  display: flex;
  justify-content: center;
  gap: 22px;
}
.main .link {
  margin-top: 32px;
  font-size: 18px;
  font-weight: bold;
}
</style>
