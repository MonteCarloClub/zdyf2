import { defineConfig } from "vite";
import vue from "@vitejs/plugin-vue";
import { resolve } from "path";

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      "@": resolve(__dirname, "./src"),
    },
  },
  // 打包后的根路径
  base: "/",
  // 打包产物存放的目录
  // build: {
  //   outDir: '../static'
  // },
  // 开发时的代理
  server: {
    proxy: {
      "/dev": {
        target: 'http://10.176.40.48/dpki',
        changeOrigin: true,
        secure: false,
        rewrite: (path) => path.replace(/^\/dev/, ""),
      },
    },
  },
  css: {
    preprocessorOptions: {
      less: {
        // all less variables could be found in:
        // https://github.com/vueComponent/ant-design-vue/blob/main/components/style/themes/default.less
        modifyVars: {
          // "primary-color": "#000000",
          // "link-color": "#0057a8",
          // "select-item-selected-bg": "#AAAAAA",
        },
        javascriptEnabled: true,
      },
    },
  },
});
