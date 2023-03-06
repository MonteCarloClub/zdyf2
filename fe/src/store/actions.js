import { mutations } from "./store";
import { userApi } from "../api/user"
import { setToken, removeToken, getToken } from '../utils/token'
import { devLog } from "../utils/log";

// 与全局状态（store/state）有关的业务接口，即会改变全局状态的操作
export const actions = {
    /**
     * 登录，成功后会修改 store 中的全局状态
     * @param {Object} userInfo 登录所需的用户信息
     * @returns promise
     */
    login(userInfo, useCert) {
        const { name, password, cert } = userInfo
        return new Promise((resolve, reject) => {
            userApi.login({ name, password, cert }, useCert)
                .then(user => {
                    setToken(`${name}-${password}`)
                    mutations.setUser(user)
                    resolve(user)
                })
                .catch(error => {
                    reject(error)
                })
        })
    },

    // 自动静默登录
    silentLogin() {
        // 从本地持久存储中获得上一次的登录信息
        let token = getToken()
        devLog(token);
        if (!token) return Promise.reject("")
        // 尝试登录
        let userInfo = token
        return this.login(userInfo)
    },

    // 退出登录
    logout() {
        return new Promise((resolve) => {
            setToken('')
            removeToken()
            // resetRouter()
            resolve()
        })
    },
}
