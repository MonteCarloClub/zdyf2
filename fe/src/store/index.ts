import { InjectionKey } from "vue";
import { createStore, useStore as baseUseStore, Store } from "vuex";
import type { App } from "vue";

export interface State {
    count: number;
}

const store = createStore<State>({
    state() {
        return {
            count: 0,
        };
    },
    mutations: {
        increment(state: State) {
            state.count++;
        },
    },
});

// https://v3.cn.vuejs.org/api/composition-api.html#provide-inject
const key: InjectionKey<Store<State>> = Symbol();

// wrap vuex's useStore
export function useStore() {
    return baseUseStore(key);
}

export function setupStore(app: App<Element>) {
    app.use(store, key);
}

export default store;