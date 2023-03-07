import router from "@/router";
import { message } from 'ant-design-vue';

/**
 * 搜索框复用
 * @param searchKey 搜索关键词
 */
export function onSearch(searchKey: string) {
  const s = searchKey;
  if (s.length > 0) {
    router.push(`/query/${s}`);
  }
  else {
    message.error('请输入证书编号')
  }
}
