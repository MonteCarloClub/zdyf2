<template>
  <div class="sign">
    <div class="sign-container">
      <el-card class="sign-card">
        <div slot="header" class="clearfix">
          <span>注册</span>
          <router-link :to="{ name: 'login' }" class="sign-card-header--button">
            前往登录
          </router-link>
        </div>

        <el-form ref="signupForm" :model="signup" :rules="signupRules">
          <el-form-item label="所在通道" prop="channel">
            <el-input v-model="signup.channel" placeholder="请输入用户所在通道"></el-input>
          </el-form-item>

          <el-form-item label="用户名" prop="name">
            <el-input v-model="signup.name" placeholder="请输入用户名"></el-input>
          </el-form-item>

          <el-form-item label="密码" prop="password">
            <el-input v-model="signup.password" placeholder="请输入密码" show-password></el-input>
          </el-form-item>

          <el-form-item label="用户角色" prop="role" style="text-align: right">
            <el-radio-group v-model="signup.role">
              <el-radio v-for="role in userRoles" :key="role.name" :label="role.name">
                {{ role.label }}
              </el-radio>
            </el-radio-group>
          </el-form-item>

          <el-form-item>
            <el-button
              :disabled="!agree"
              style="width: 100%"
              type="primary"
              @click="onSignupSubmit"
              :loading="loading"
            >
              注册
            </el-button>
          </el-form-item>
        </el-form>
        <el-checkbox v-model="agree">我自愿承担使用本系统过程中可能出现的风险</el-checkbox>
      </el-card>
    </div>
  </div>
</template>

<script>
import { userApi } from "@/api/user.js";

export default {
  name: "SignUp",
  components: {},

  data() {
    return {
      signup: {
        name: "",
        password: "",
        role: "user", // 用户默认角色
        channel: "", // 请输入用户所在通道
      },
      signupRules: {
        name: [{ required: true, trigger: "blur", message: "用户名不能为空" }],
        password: [{ required: true, trigger: "blur", message: "密码不能为空" }],
        role: [{ required: true, trigger: "blur", message: "请勾选用户角色" }],
        channel: [{ required: true, trigger: "blur", message: "请输入用户所在通道" }],
      },
      userRoles: [
        {
          label: "普通用户",
          name: "user",
        },
        {
          label: "机构用户",
          name: "org",
        },
      ],
      agree: false,
      loading: false,
    };
  },

  methods: {
    onSignupSubmit() {
      this.$refs.signupForm.validate((valid) => {
        if (!valid) return;
        // 开始加载
        this.loading = true;
        userApi
          .signup(this.signup)
          .then((data) => {
            this.$message({
              message: "注册成功",
              type: "success",
            });
            // 下载证书内容
            data.certificate && this.download(JSON.stringify(data.certificate));
            
            this.$router.push({
              name: "login",
              query: { name: this.signup.name },
            });
          })
          .catch((err) => {
            this.$message({
              message: err,
              type: "error",
            });
          })
          .finally(() => {
            // 结束加载
            this.loading = false;
          });
      });
    },

    download(content) {
      console.log(content);
      const fileName = 'certificate.cert'
      var a = document.createElement("a");
      var file = new Blob([content], { type: 'text/plain' });
      a.href = URL.createObjectURL(file);
      a.download = fileName;
      a.click();
    },
  },
};
</script>

<style scoped>
.sign-container {
  max-width: 600px;
  margin: 60px auto;
}

.sign-card {
  width: 400px;
  margin: auto;
}

.sign-card-header--button {
  float: right;
  text-decoration: none;
  color: #79bbff;
}

.sign-card-header--button:hover {
  color: #409eff;
}

.tips {
  line-height: 19px;
  font-size: 14px;
  color: #606266;
  font-weight: 500;
}
</style>
