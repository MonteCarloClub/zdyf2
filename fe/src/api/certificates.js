import axios from 'axios'

const CERT_HOST = "/dpki/"
const certService = axios.create({
    baseURL: process.env.NODE_ENV === "development" ? '/cert' : CERT_HOST,
})

// request interceptor
certService.interceptors.request.use(
    config => {
        config.headers['Access-Control-Allow-Origin'] = "*"
        config.headers['Access-Control-Allow-Methods'] = "GET, POST, PATCH, PUT, DELETE, OPTIONS"
        config.headers['Access-Control-Allow-Headers'] = "Origin, Content-Type, X-Auth-Token"
        config.crossdomain = true
        return config
    },
    error => {
        // do something with request error
        console.log('[request error]', error)
        return Promise.reject(error)
    }
)

export const certApi = {
    /**
     * 获取证书列表
     * @param {*} _data 
     * @returns Promise
     */
    list: function () {
        return new Promise((resolve, reject) => {
            certService.request({
                url: '/IoTDevTest',
                method: 'get',
            }).then(response => {
                if (response.status === 200) {
                    resolve(response.data)
                }
                else {
                    reject(response)
                }
            }).catch(reject)
        })
    },

    /**
     * 撤销证书
     * @param {*} _data 
     * @returns Promise
     */
    revoke: function ({ no }) {
        return new Promise((resolve, reject) => {
            certService.request({
                url: '/RevokeABSCertificate',
                method: 'get',
                params: { no }
            }).then(response => {
                if (response.status === 200) {
                    resolve(response.data)
                }
                else {
                    reject(response)
                }
            }).catch(reject)
        })
    },

    /**
     * 查询证书详细信息
     * @param {*} _data 
     * @returns Promise
     */
    certInfo: function (uid) {
        return new Promise((resolve, reject) => {
            certService.request({
                url: '/GetCertificateByUID',
                method: 'get',
                params: {uid}
            }).then(response => {
                if (response.status === 200) {
                    resolve(response.data)
                }
                else {
                    reject(response)
                }
            }).catch(reject)
        })
    },

    /**
     * 申请证书
     * @param {*} _data 
     * @returns Promise
     */
    apply: function (uid, attribute) {
        return new Promise((resolve, reject) => {
            certService.request({
                url: '/ApplyForABSCertificate',
                method: 'get',
                params: {uid, attribute}
            }).then(response => {
                if (response.status === 200) {
                    resolve(response.data)
                }
                else {
                    reject(response)
                }
            }).catch(reject)
        })
    }
}