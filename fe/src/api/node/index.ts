import { request } from '@/api/request';

/**
 * 获取节点分数
 */
export function score(params: API.NodeScoreParams) {
    return request<API.NodeScoreResponse>(
        {
            url: '/getScore',
            method: 'get',
            params
        }
    );
}
