<template>
  <div class="content">
    <div class="wer">
      <div class="wer-left">
        <el-input placeholder="设备名称" v-model="searchInput" />
      </div>
      <div>
        <el-button @click="searchCert"> 查询证书 </el-button>
        <el-button type="success" @click="applyFormVisible = true"> 证书申请 </el-button>
      </div>
    </div>

    <el-table :data="certificates" :empty-text="emptyText">
      <el-table-column show-overflow-tooltip label="设备名称" prop="ABSUID" />
      <el-table-column show-overflow-tooltip label="证书序号" prop="serialNumber" />
      <el-table-column label="操作" align="right">
        <template slot-scope="scope">
          <el-button size="mini" type="warning" @click="revoke(scope)" plain> 撤销证书 </el-button>
          <el-button size="mini" type="info" @click="detailInfo(scope.row)"> 详细信息 </el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog title="详细信息" :visible.sync="detailVisible">
      <el-descriptions :column="1">
        <el-descriptions-item label="版本">{{ detail.version }}</el-descriptions-item>
        <el-descriptions-item label="设备">{{ detail.ABSUID }} </el-descriptions-item>
        <el-descriptions-item label="证书序号"> {{ detail.serialNumber }} </el-descriptions-item>
        <el-descriptions-item label="标签"> {{ detail.ABSAttribute }} </el-descriptions-item>
        <el-descriptions-item label="签名人">{{ detail.issuer }}</el-descriptions-item>
        <el-descriptions-item label="签名">{{ detail.signatureName }}</el-descriptions-item>
        <el-descriptions-item label="优先级">{{ detail.validityPeriod }}</el-descriptions-item>
      </el-descriptions>

      <div slot="footer">
        <el-button type="primary" @click="detailVisible = false">确 定</el-button>
      </div>
    </el-dialog>

    <el-dialog title="证书申请" :visible.sync="applyFormVisible">
      <el-form ref="applyForm" :rules="applyRules" :model="applyCert" label-width="80px">
        <el-form-item prop="uid" label="UID">
          <el-input v-model="applyCert.uid"></el-input>
        </el-form-item>
        <el-form-item prop="attribute" label="标签">
          <el-input v-model="applyCert.attribute"></el-input>
        </el-form-item>
        <div v-if="applyRes">
          <el-form-item size="mini" label="版本">{{ applyRes.version }}</el-form-item>
          <el-form-item size="mini" label="设备">{{ applyRes.ABSUID }} </el-form-item>
          <el-form-item size="mini" label="证书序号"> {{ applyRes.serialNumber }} </el-form-item>
          <el-form-item size="mini" label="标签"> {{ applyRes.ABSAttribute }} </el-form-item>
          <el-form-item size="mini" label="签名人">{{ applyRes.issuer }}</el-form-item>
          <el-form-item size="mini" label="签名">{{ applyRes.signatureName }}</el-form-item>
          <el-form-item size="mini" label="优先级">{{ applyRes.validityPeriod }}</el-form-item>
        </div>
      </el-form>

      <div slot="footer">
        <el-button @click="applyFormVisible = false">取 消</el-button>
        <el-button type="primary" @click="applyForCert">申 请</el-button>
      </div>
    </el-dialog>
  </div>
</template>

<script>
import { certApi } from "@/api/certificates";

export default {
  name: "List",
  components: {},
  data() {
    return {
      certificates: [],
      emptyText: "暂无数据",
      detailVisible: false,
      detail: {},

      applyFormVisible: false,
      applyCert: {},

      applyRules: {
        uid: [{ required: true, trigger: "blur", message: "UID 不能为空" }],
        attribute: [{ required: true, trigger: "blur", message: "标签不能为空" }],
      },

      searchInput: "",
      applyRes: false,
    };
  },

  mounted() {
    this.emptyText = "正在加载...";
    this.getlist();
  },

  methods: {
    getlist() {
      certApi
        .list()
        .then((res) => {
          this.certificates = res.map((item) => {
            item = JSON.parse(item);
            // item.certificate["absSignature"] = item.absSignature;
            return item;
          });
          this.emptyText = "暂无证书";
        })
        .catch(console.log);
    },
    detailInfo(cert) {
      this.detailVisible = true;
      this.detail = cert;
    },

    revoke(scope) {
      const no = scope.row.serialNumber;
      certApi
        .revoke({ no })
        .then((res) => {
          if (res === "Revoke OK.") {
            this.$message({
              message: "撤销成功",
              duration: 2 * 1000,
              type: "success",
            });
            this.certificates.splice(scope.$index, 1);
          } else {
            this.$message({
              message: res,
              type: "error",
            });
          }
        })
        .catch((e) => {
          this.$message({
            message: e.message,
            type: "error",
          });
        });
    },

    applyForCert() {
      this.$refs.applyForm.validate((valid) => {
        if (!valid) return;
        const { uid, attribute } = this.applyCert;
        this.applyRes = false;

        certApi
          .apply(uid, attribute)
          .then((item) => {
            item.certificate["absSignature"] = item.absSignature;
            this.applyRes = item.certificate;
            this.applyFormVisible = false;
            this.getlist()
            this.$message({
              message: "申请成功",
              duration: 2 * 1000,
              type: "success",
            });
          })
          .catch((e) => {
            this.$message({
              message: e.message,
              type: "error",
            });
          });
      });
    },

    searchCert() {
      const a = this.searchInput;
      if (a) {
        certApi
          .certInfo(a)
          .then((item) => {
            item.certificate["absSignature"] = item.absSignature;
            this.certificates = [item.certificate];
          })
          .catch((e) => {
            this.$message({
              message: e.message,
              type: "error",
            });
          });
      }
    },
  },
};
</script>

<style scoped>
.content {
  padding: 16px;
  background-color: var(--bg-color-0, white);
}

.wer {
  display: flex;
  gap: 16px;
  padding: 8px 0;
}

.wer-left {
  flex: 1;
}
</style>
