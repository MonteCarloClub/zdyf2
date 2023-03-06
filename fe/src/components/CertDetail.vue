<template>
  <a-form
    ref="formRef"
    name="custom-validation"
    :model="cert"
    :label-col="{ span: 6 }"
    :wrapper-col="{ span: 18 }"
    labelAlign="left"
  >
    <a-form-item label="版本">
      {{ cert?.version }}
    </a-form-item>
    <a-form-item label="设备">
      {{ cert?.ABSUID }}
    </a-form-item>
    <a-form-item label="证书序号">
      {{ cert?.serialNumber }}
    </a-form-item>
    <a-form-item label="标签">
      {{ cert?.ABSAttribute }}
    </a-form-item>
    <a-form-item label="签名人">
      {{ cert?.issuer }}
    </a-form-item>
    <a-form-item label="签名">
      {{ cert?.signatureName }}
    </a-form-item>
    <a-form-item label="优先级">
      {{ cert?.validityPeriod }}
    </a-form-item>
  </a-form>
</template>

<script lang="ts" setup>
import { PropType } from "vue";

defineProps({
  cert: {} as PropType<API.Cert>,
});
</script>

<style scoped></style>
