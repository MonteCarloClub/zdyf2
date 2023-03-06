<script setup lang="ts">
import { ref, computed } from "vue";
import CertDetail from "@/components/CertDetail.vue";
import { list, revoke } from "@/api/cert";
import { Modal } from "ant-design-vue";
import { message } from "ant-design-vue";

const certs = ref<API.Cert[]>([]);

list().then((res) => {
  if (res) {
    const list = res.map((certJson) => JSON.parse(certJson) as API.Cert);
    certs.value = list;
  }
});

const table = computed(() =>
  certs.value?.map((data, index) => {
    const row: API.Cert & {
      key?: string;
    } = data;
    if (!Object.hasOwn(data, "key")) {
      row.key = index.toString();
    }
    return row;
  })
);

interface Column {
  key: string;
  title: string;
  dataIndex: keyof API.Cert | string;
  width?: number;
  ellipsis?: boolean;
  align?: "left" | "center" | "right";
}
const columns: Column[] = [
  { key: "ABSUID", dataIndex: "ABSUID", width: 100, title: "设备名称" },
  {
    key: "serialNumber",
    dataIndex: "serialNumber",
    width: 160,
    title: "证书序号",
  },
  { key: "version", dataIndex: "version", width: 80, title: "版本" },
  {
    key: "operation",
    dataIndex: "operation",
    width: 80,
    title: "",
    align: "right",
  },
];

const certDetailFormVisible = ref(false);
const certDetail = ref<API.Cert>();
function viewDetail(record: API.Cert) {
  certDetail.value = record;
  certDetailFormVisible.value = true;
}

function deleteRow(record: API.Cert) {
  Modal.confirm({
    title: "不可逆操作",
    content: `该操作会撤销证书 ${record.serialNumber}`,
    okText: "确认撤销",
    onOk() {
      revokeCert(record);
    },
    cancelText: "取消",
  });
}

const revokingCertNumber = ref("");
function revokeCert(record: API.Cert) {
  revokingCertNumber.value = record.serialNumber;
  revoke({
    no: record.serialNumber,
  })
    .then((res) => {
      if (res === "Revoke OK.") {
        message.success("撤销成功");
        certs.value = certs.value.filter(
          (cert) => cert.serialNumber != record.serialNumber
        );
      }
    })
    .catch(err => {
        message.error(err);
    })
    .finally(() => {
      revokingCertNumber.value = "";
    });
}
</script>

<template>
  <a-table :columns="columns" :data-source="table" size="middle">
    <template #bodyCell="{ column, text, record }">
      <template v-if="column.key === 'operation'">
        <span>
          <a-button
            danger
            @click="deleteRow(record)"
            :loading="record.serialNumber === revokingCertNumber"
          >
            撤销证书
          </a-button>
          <a-button type="link" @click="viewDetail(record)"> 详情 </a-button>
        </span>
      </template>
    </template>
  </a-table>

  <a-modal v-model:visible="certDetailFormVisible" title="证书详情">
    <template #footer>
      <a-button
        key="submit"
        type="primary"
        @click="certDetailFormVisible = false"
      >
        确认
      </a-button>
    </template>
    <CertDetail :cert="certDetail" />
  </a-modal>
</template>

<style scoped></style>
