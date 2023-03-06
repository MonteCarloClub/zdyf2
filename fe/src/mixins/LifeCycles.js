/**
 * 生命周期钩子函数插桩
 * 直接在 Component 组件内解构即可
 * 
 * @returns 包含生命周期钩子的 Object
 */
function logLifeCycles() {
    let logs = {};
    [
        "beforeCreate",
        "created",
        "beforeMount",
        "mounted",
        "activated",
        "deactivated",
        "beforeUpdate",
        "updated",
        "destroyed",
    ].map(hook => logs[hook] = function () {
        // this 指向最终调用该函数的组件
        console.log("[lifecyle log]", this.$options.name, hook);
    })
    return logs;
}

// mixin object
export const LifeCycleMixin = {
    ...logLifeCycles()
}