import { request } from '@/api/request';

/**
 * 获取用户的详细信息
 * @param data 请求参数
 * @returns 用户详细信息
 */
export function userInfo(data: API.UserInfoParams) {
    return request<API.UserInfoResponse>(
        {
            url: '/user/info',
            method: 'post',
            data
        }
    );
}
