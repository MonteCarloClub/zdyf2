import Vue from 'vue';

const state = Vue.observable({
    user: false,
});

// 暴露给外部用于获取全局状态的接口
export const getters = {
    userName: () => state.user.name,
    isLogin: () => state.user ? true : false,
    // 根据传入的属性列表，返回一组箭头函数
    mapUser: (properties = []) => {
        const computed = {}
        properties.forEach(prop => computed[prop] = () => state.user[prop])
        return computed
    },
    // 根据传入的属性列表，返回这组属性值
    properties: (props = []) => {
        const values = {}
        props.forEach(prop => values[prop] = state.user[prop])
        return values
    },
}

// 暴露给外部用于修改全局状态的接口
export const mutations = {
    setUser: (data) => {
        state.user = {
            name: data.Name,
            password: data.Password,
            role: data.UserType,
            channel: data.Channel,
            ...data
        }
    },
}
