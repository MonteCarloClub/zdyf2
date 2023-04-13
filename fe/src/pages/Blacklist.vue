<script setup lang="ts">
import { ref, computed } from "vue";
import CertDetail from "@/components/CertDetail.vue";
import SearchInput from "@/components/SearchInput.vue";
// import { list, revoke } from "@/api/cert";
import { blacklist, removeFromBlacklist, addToBlacklist } from "@/api/user";
import { Modal, message } from "ant-design-vue";
import { getStorage, setStorage } from "@/utils/storage";

const users = ref<API.UserParams[]>([]);

blacklist().then((res) => {
  console.log("res: ", res)
  if (res.certificates) {
    console.log("certificates: ", res.certificates);
    const list = res.certificates.map(uid => ({ "uid": uid }));
    users.value = list;
    console.log("users: ", list);
  }
});

const table = computed(() =>
  users.value?.map((data, index) => {
    const row: API.DemoParams = data;
    if (data.hasOwnProperty("key")) {
      row.key = index.toString();
    }
    return row;
  })
);

console.log(table)

interface Column {
  key: string;
  title: string;
  dataIndex: keyof API.Cert | string;
  width?: number;
  ellipsis?: boolean;
  align?: "left" | "center" | "right";
  sorter?: (a: API.Cert, b: API.Cert) => number;
  customFilterDropdown?: boolean;
  /**
   * Callback executed when the confirm filter button is clicked, Use as a filter event when using template or jsx
   * @type Function
   */
  onFilter?: (value: any, record: any) => boolean;

  /**
   * Callback executed when filterDropdownVisible is changed, Use as a filterDropdownVisible event when using template or jsx
   * @type Function
   */
  onFilterDropdownVisibleChange?: (visible: boolean) => void;
}

const columns: Column[] = [
  {
    key: "UID",
    dataIndex: "UID",
    width: 100,
    title: "用户",
  },
  {
    key: "operation",
    dataIndex: "operation",
    width: 80,
    title: "",
    align: "right",
  },
];

const addBlackFormVisible = ref(false);
const uidAddBlack = ref<string>("");
const addUid = ref("")
function addUserToBlacklist(record: string) {
  addUid.value = uidAddBlack.value;
  addToBlacklist({
    uid: record
  })
    .then((res) => {
      if (res === "Add to blacklist success.") {
        message.success("添加成功");
        addBlackFormVisible.value = false;
      }
    })
    .catch((err) => {
      message.error(err);
    })
    .finally(() => {
      addUid.value = ""
    });
}

function deleteRow(record: API.UserParams) {
  Modal.confirm({
    title: "不可逆操作",
    content: `该操作会移出黑名单用户 ${record.uid}`,
    okText: "确认移出",
    onOk() {
      rmFromBlacklist(record);
    },
    cancelText: "取消",
  });
}

const rmUid = ref("");
function rmFromBlacklist(record: API.UserParams) {
  rmUid.value = record.uid;
  removeFromBlacklist({
    uid: record.uid,
  })
    .then((res) => {
      if (res === "Revoke OK.") {
        message.success("撤销成功");
        users.value = users.value.filter(
          (user) => user.uid != record.uid
        );
      }
    })
    .catch((err) => {
      message.error(err);
    })
    .finally(() => {
      rmUid.value = "";
    });
}

const DEFAULT_PAGE_SIZE = 'DEFAULT_PAGE_SIZE';
const dpagesize = getStorage<string>(DEFAULT_PAGE_SIZE) || '10';
const defaultPageSize = ref(parseInt(dpagesize));

function tableChanged(pagination: any) {
  setStorage(DEFAULT_PAGE_SIZE, pagination.pageSize);
}

</script>

<template>
  <div class="panel">
    <div class="flex-expand nav">
      <router-link to="/"> 首页 </router-link>
    </div>
    <div class="flex-expand nav" style="text-align: right;">
      <a-button @click="addBlackFormVisible = true"> 添加黑名单 </a-button>
      <a-modal v-model:visible="addBlackFormVisible" title="添加用户">
        <a-input v-model:value="uidAddBlack" placeholder="请输入UID" />
        <template #footer>
          <a-button key="back" @click="addBlackFormVisible = false">取消</a-button>
          <a-button key="submit" type="primary" :loading="addUid === uidAddBlack" @click="addUserToBlacklist(uidAddBlack)"
            :disabled="uidAddBlack === ''">确认</a-button>
        </template>
      </a-modal>
    </div>

  </div>
  <div class="content">
    <a-table :columns="columns" :data-source="table" size="middle"
      :pagination="{ position: ['bottomCenter'], hideOnSinglePage: true, defaultPageSize }" @change="tableChanged">
      <template #bodyCell="{ column, text, record }">
        <template v-if="column.key === 'operation'">
          <span>
            <a-button type="link" danger @click="deleteRow(record)" :loading="record.uid === rmUid">
              移出黑名单
            </a-button>
          </span>
        </template>
      </template>
    </a-table>
  </div>
</template>

<style scoped>
.panel {
  display: flex;
  text-align: left;
  align-items: center;
  margin: 0 32px 22px;
}

.search-container {
  width: 600px;
}

.flex-expand {
  flex: 1;
}

.content {
  margin: 0 32px;
}

.nav a {
  font-size: 22px;
  color: black;

}
</style>
