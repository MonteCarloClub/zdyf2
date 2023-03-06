/**
 * 异步函数的单例封装
 * 多处调用同一个经过此封装的 api 时，实际上共享同一个 promise 实例，发送一个请求
 * @param {Function} asyncFunc 异步函数，返回一个 Promise
 * @returns 正在执行的 promise 单例 
 */
export function singletonPromise(asyncFunc) {
    // 定义一个实例
    let promiseInstance = undefined;

    return function (...args) {
        // 如果实例存在，就直接返回该实例
        if (promiseInstance) {
            return promiseInstance
        }
        
        // 根据传进来的异步函数新建一个实例
        promiseInstance = asyncFunc.call(this, ...args)
            .then(res => {
                return Promise.resolve(res)
            })
            .catch(err => {
                return Promise.reject(err)
            })
            .finally(() => {
                // 执行完成之后清空单例
                promiseInstance = undefined;
            })

        return promiseInstance
    }
}