/**防抖 */
interface DebouncedFunc<T extends (...args: any[]) => any> {
    (...args: Parameters<T>): void;
}

export function debounce<T extends (...args: any[]) => any>(callback: T, delay: number): DebouncedFunc<T> {
    if (delay === void 0) { delay = 200; }
    let timer: ReturnType<typeof setTimeout>

    return (...args: Parameters<T>) => {
        clearTimeout(timer)
        timer = setTimeout(() => {
            callback(...args)
        }, delay)
    }
}