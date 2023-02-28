declare namespace API {
  type ErrorResponse = {
    /** 业务约定的错误码 */
    errorCode: string;
    /** 业务上的错误信息 */
    errorMessage?: string;
    /** 业务上的请求是否成功 */
    success?: boolean;
  };

  /** 与后端约定好自定义的 Response 结构 */
  type Response<T = any> = {
    code: number;
    message: string;
    data: T;
  };
}
