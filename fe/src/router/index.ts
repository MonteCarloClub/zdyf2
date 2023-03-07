import type { App } from 'vue'
import { createWebHistory, createRouter } from "vue-router";

const router = createRouter({
    history : createWebHistory(),
    routes: [
        {
            path: '/',
            component: () => import("@/pages/Home.vue")
        },
        {
            path: '/certs',
            component: () => import("@/pages/Certs.vue")
        },
        {
            path: '/query/:no',
            component: () => import("@/pages/Query.vue")
        },
    ]
})

export function setupRouter(app: App<Element>) {
    app.use(router);
}

export default router