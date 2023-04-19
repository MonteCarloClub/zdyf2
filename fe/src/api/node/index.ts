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


/**
 * 获取节点名称
 */
export function caName() {
    return request<string>(
        {
            url: '/getCAName',
            method: 'get'
        }
    );
}