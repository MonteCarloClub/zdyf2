import { request } from '@/api/request';

/**
 * 把用户移出黑名单
 */
export function removeFromBlacklist(params: API.UserParams) {
    return request<API.rmBlacklistResponse>(
        {
            url: '/removeFromBlacklist',
            method: 'get',
            params
        }
    );
}

/**
 * 把用户加入黑名单
 */
export function addToBlacklist(params: API.UserParams) {
    return request<API.addBlacklistResponse>(
        {
            url: '/addToBlacklist',
            method: 'get',
            params
        }
    );
}

/**
 * 获取黑名单
 */
export function blacklist() {
    return request<API.BlacklistResponse>(
        {
            url: '/getBlacklist',
            method: 'get'
        }
    );
}