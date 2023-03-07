import axios from "axios";
import { message as $message } from "ant-design-vue";
import type { AxiosRequestConfig } from "axios";
import { UNKNOWN_ERROR_MSG, REQUEST_BASE_URL } from "@/common/constants";

// 创建 axios 实例
const service = axios.create({
  baseURL: REQUEST_BASE_URL,
  timeout: 6000,
});

service.interceptors.request.use(
  (config) => {
    // 不需要鉴权
    return config;
  },
  (error) => {
    Promise.reject(error);
  }
);

service.interceptors.response.use(
  (response) => {
    // Any status code that lie within the range of 2xx cause this function to trigger
    // Do something with response data
    return response;
  },
  (error) => {
    // 处理 422 或者 500 的错误异常提示
    const errMsg = error?.response?.data?.message ?? UNKNOWN_ERROR_MSG;
    // $message.error(errMsg);
    error.message = errMsg;
    return Promise.reject(error);
  }
);

export interface RequestOptions {
  /** 请求成功是提示信息 */
  successMsg?: string;
  /** 请求失败是提示信息 */
  errorMsg?: string;
}

export const request = async <T>(
  config: AxiosRequestConfig,
  options: RequestOptions = {}
): Promise<T> => {
  try {
    // 设置这里的 request 返回的数据格式为 Response<T>
    const res = await service.request<T>(config);

    // 弹窗提示消息
    const { successMsg, errorMsg } = options;
    successMsg && $message.success(successMsg);
    errorMsg && $message.error(errorMsg);
    // 取 axios 的 data
    return res.data;
  } catch (error: any) {
    return Promise.reject(error);
  }
};
