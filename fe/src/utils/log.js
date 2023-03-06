
/**
 * 打印带有调用栈信息的数据
 * @param {String} msg 需要打印的字符
 * @param {Number} depth 调用栈深度
 */
function log(depth = 3, ...msg) {
    // 只在开发环境下打印
    if (process.env.NODE_ENV != "development") return
    // 获取函数的调用栈
    let err = new Error()
    let stack = err.stack
        .split("\n")
        .slice(2, depth + 2)
        .map(line => line.match(/at (.*) /g)[0])
        .reverse()
        .map(at_f => at_f.split(" ")[1])

    if (msg.length > 1) {
        console.log(`[${stack.join(" -> ")}] :`);
        console.log(...msg);
    }
    else {
        console.log(`[${stack.join(" -> ")}] :`, ...msg);
    }
}

/**
 * 返回一个固定第一个参数的函数
 * @param {Function} f 待 curried 的函数
 * @returns 
 */
function curry(f) {
    return function (fst) {
        return function (...args) {
            f.apply(this, [fst, ...args])
        }
    }
}

export const devLog = curry(log)(3)