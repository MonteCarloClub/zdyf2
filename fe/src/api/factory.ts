/**
 * 参考：
 * https://staging-cn.vuejs.org/guide/reusability/composables.html
 * https://staging-cn.vuejs.org/guide/typescript/composition-api.html#typing-ref
 */
import { ref, isReactive, watch } from 'vue'
import type { Ref } from 'vue'

/**
 * 工厂函数，根据传入的 api 函数构造相应的组合式函数，便于组件中复用
 * @param fetchMethod api 函数
 * @returns 组合式函数
 */
export function useFetchFactory<T extends Object, R>
    (fetchMethod: (params: T) => Promise<API.Response<R>>) {

    const data: Ref<R | null> = ref(null);
    const error: Ref<string | null> = ref(null);

    return function (params: T) {
        // 取值函数
        async function doFetch() {
            // reset state before fetching..
            data.value = null;
            error.value = null;

            fetchMethod(params)
                .then((res) => data.value = res.data)
                .catch((err) => (error.value = err))
        }

        if (isReactive(params)) {
            // setup reactive re-fetch if input params is reactive
            watch(params, () => {
                doFetch();
            })
        }

        // otherwise, just fetch once
        doFetch()

        return { data, error }
    }
}