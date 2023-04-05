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
export function query(params: API.QueryParams) {
    return request<API.QueryResponse>(
        {
            url: '/GetCertificate',
            method: 'get',
            params
        }
    );
}

/**
 * 申请证书
 */
export function apply(params: API.ApplyParams) {
    return request<API.QueryResponse>(
        {
            url: '/ApplyForABSCertificate',
            method: 'get',
            params
        }
    );
}

/**
 * 验证证书
 */
export function verify(data: API.VerifyParams) {
    return request<API.VerifyResponse>(
        {
            url: '/VerifyABSCert',
            method: 'post',
            headers: { "Content-Type": "application/json-patch+json" },
            data
        }
    );
}


/**
 * 查询证书链上状态
 */
export function statusOnChain(params: API.QueryParams) {
    return request<API.CertOnChainResponse>(
        {
            url: '/GetCertificateFromFabric',
            method: 'get',
            params
        }
    );
}

/**
 * 查询证书链上状态
 */
export function history(params: API.HistoryParams) {
    return request<API.HistoryResponse>(
        {
            url: '/getCertificates',
            method: 'get',
            params
        }
    );
}