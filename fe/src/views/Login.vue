<template>
  <div class="sign">
    <div class="sign-container">
      <el-card class="sign-card">
        <div slot="header" class="clearfix">
          <span>登录</span>
          <router-link :to="{ name: 'signup' }" class="sign-card-header--button">
            前往注册
          </router-link>
        </div>

        <el-form ref="loginForm" :model="login" :rules="loginRules" label-position="top">
          <el-form-item prop="name" label="用户名">
            <el-input v-model="login.name" placeholder="请输入用户名" maxlength="11"></el-input>
          </el-form-item>

          <el-form-item prop="password" label="密码" v-if="useCert === false">
            <el-input
              v-model="login.password"
              placeholder="请输入密码"
              show-password
              @keyup.enter.native="onLoginSubmit"
            ></el-input>
          </el-form-item>

          <el-form-item prop="cert" label="证书" v-else>
            <el-input v-model="login.cert" v-show="false"></el-input>
            <el-button @click="selectFile"> 上传证书文件 </el-button>
            <div class="text">
              {{ certFileName }}
            </div>
          </el-form-item>

          <el-form-item>
            <el-button v-if="useCert === false" type="text" @click="useCert = true"
              >使用证书登录</el-button
            >
            <el-button v-else type="text" @click="useCert = false">使用密码登录</el-button>
          </el-form-item>

          <el-form-item>
            <el-button style="width: 100%" type="primary" @click="onLoginSubmit" :loading="loading">
              登录
            </el-button>
          </el-form-item>
        </el-form>
      </el-card>
    </div>
  </div>
</template>

<script>
import { actions } from "../store/actions";

export default {
  name: "Login",
  components: {},

  data() {
    return {
      login: {
        name: "",
        password: "",
        cert: "",
      },
      loginRules: {
        name: [{ required: true, trigger: "blur", message: "用户名不能为空" }],
        password: [{ required: true, trigger: "blur", message: "密码不能为空" }],
        cert: [{ required: true, trigger: "blur", message: "证书不能为空" }],
      },
      certFileName: "",
      useCert: false, // 是否使用证书登录，默认是用密码登录
      loading: false,
    };
  },

  created() {
    const name = this.$route.query["name"];
    if (name) {
      this.login.name = name;
    }
  },

  methods: {
    onLoginSubmit() {
      this.$refs.loginForm.validate((valid) => {
        if (!valid) return;
        this.loading = true;
        actions
          .login(this.login, this.useCert)
          .then(() => {
            this.$message({
              message: "登录成功",
              type: "success",
            });
            this.jumpTo();
          })
          .catch((err) => {
            this.$message({
              message: err,
              type: "error",
            });
          })
          .finally(() => {
            this.loading = false;
          });
      });
    },

    selectFile() {
      var input = document.createElement("input");
      input.type = "file";

      input.onchange = (e) => {
        var file = e.target.files[0];
        this.certFileName = file.name;
        // setting up the reader
        var reader = new FileReader();
        reader.readAsText(file, "UTF-8");

        // here we tell the reader what to do when it's done reading...
        reader.onload = (readerEvent) => {
          var content = readerEvent.target.result; // this is the content!
          const cert = JSON.parse(content);
          if (cert.serialNumber) {
            this.login.cert = cert.serialNumber.split('-')[1];
          }
        };
      };

      input.click();
    },

    jumpTo() {
      const redirect = this.$route.query["redirect"];
      if (redirect) {
        // 登录成功后回到想看的页面
        this.$router.push({ path: redirect });
      } else this.$router.push({ path: "/" });
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

.text {
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
</style>
