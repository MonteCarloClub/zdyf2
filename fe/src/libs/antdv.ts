import type { App } from "vue";
import Antd from "ant-design-vue";
import "ant-design-vue/dist/antd.less";

export function setupAntd(app: App<Element>): void {
  app.use(Antd);
}