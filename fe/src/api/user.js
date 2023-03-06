import request from '@/utils/request'

export const userApi = {
    /**
     * 注册用户
     * @param {*} _data 来自 ui 界面的参数，用于发送请求，字段不一定和请求字段一致，需要转换一下
     * @returns Promise
     */
    signup: function (_data) {
        // 根据接口文档将字段转换成后台需要的
        // userName 用户名
        // userType 用户类型(user/org)
        // channel  用户所在通道
        // password 密码
        const data = {
            userName: _data.name,
            password: _data.password,
            userType: _data.role,
            channel: _data.channel
        }

        return new Promise((resolve, reject) => {
            request({
                url: '/user/create',
                method: 'post',
                data
            }).then(response => {
                // {
                //     "code":200, //200，成功；其他，失败
                //     "msg":"",
                //     "data":null 
                // }
                if (response.code === 200) {
                    resolve(response.data)
                }
                else {
                    reject(response)
                }
            }).catch(reject)
        })
    },

    /**
     * 登录
     * @param {Object} _data 用于登录的信息用户名密码
     * @returns 登录成功则返回用户信息
     */
    login: function (_data, useCert = false) {
        // fileName 用户名
        // password 密码
        // cert     证书
        const data = {
            fileName: _data.name,
            password: _data.password,
            cert: _data.cert,
        }

        const apiUrl = useCert ? '/dabe/user3' : '/dabe/user2'
        
        return new Promise((resolve, reject) => {
            request({
                url: apiUrl,
                method: 'post',
                data,
                params: data
            }).then(response => {
                // {
                //    "code":200   200，成功；其他，失败
                //    "msg" :null  描述
                //    "data": {
                //         "appliedAttrMap":{},                  用户已申请的属性集合
                //         "privacyAttrMap":{},                  隐私属性集，暂时不用
                //         "APKMap":{},                          用户自身的属性公钥集合
                //         "ASKMap":{},                          用户自身的属性私钥集合
                //         "EGGAlpha":"[206605..., 320061...]",  属性密码相关参数-用户公钥
                //         "Alpha":"907358...",                  属性密码相关参数-用户私钥
                //         "GAlpha":"[569334..., 105875...]",    属性密码相关参数（这些应该不用显示）
                //         "Name":"someone",                     用户名
                //         "OPKMap":{},                          由多个用户组成的组织公钥集合，可遍历该集合获取用户所在组织列表
                //         "OSKMap":{},                          由多个用户组成的组织私钥集合
                //         "Password":"202cb962ac59075b96..",    用户密码 hash
                //         "UserType":"org",                     用户类型
                //         "Channel":"myc"                       用户所在通道
                //     }
                // }
                const { data } = response
                // 做一个字段转换
                resolve(data)
            }).catch(reject)
        })
    }
}