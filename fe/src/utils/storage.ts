import { isNull, isObject } from "@/utils/is";
import { Nullable } from "@/utils/types";

export type StorageObj = Object | string | undefined;

/**
 * 保存在 localstorage
 * @param key 键
 * @param obj 值，object 或者 string
 * @returns void
 */
export function setStorage(key: string, obj: Nullable<StorageObj>): void {
    if (obj === undefined || isNull(obj)) {
        localStorage.setItem(key, '')
        return;
    }
    if (isObject(obj)) {
        localStorage.setItem(key, JSON.stringify(obj))
        return;
    }
    localStorage.setItem(key, obj)
}

/**
 * 从 localstorage 中取值
 * @param key 键
 * @returns T
 */
export function getStorage<T>(key: string): Nullable<T> {
    const s =  localStorage.getItem(key);
    if (s?.startsWith('{')) {
        return JSON.parse(s)
    }
    return s as T;
}