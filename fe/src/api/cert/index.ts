import { request } from '@/api/request';

/**
 * 获取证书列表
 */
export function list() {
    return request<string[]>(
        {
            url: '/IoTDevTest',
            method: 'get'
        }
    );
}

/**
 * 撤销证书
 */
export function revoke(params: API.RevokeParams) {
    return request<API.RevokeResponse>(
        {
            url: '/RevokeABSCertificate',
            method: 'get',
            params
        }
    );
}

/**
 * 查询证书详细信息
 */
export function info(params: API.InfoParams) {
    return request<API.Cert>(
        {
            url: '/GetCertificateByUID',
            method: 'get',
            params
        }
    );
}

/**
 * 申请证书
 */
export function apply(params: API.ApplyParams) {
    return request<any>(
        {
            url: '/ApplyForABSCertificate',
            method: 'get',
            params
        }
    );
}
