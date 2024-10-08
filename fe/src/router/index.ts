import type { App } from 'vue'
import { createWebHistory, createRouter } from "vue-router";

const router = createRouter({
    history : createWebHistory(import.meta.env.BASE_URL),
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
            path: '/create',
            component: () => import("@/pages/Create.vue")
        },
        {
            path: '/query/:no',
            component: () => import("@/pages/Query.vue")
        },
        {
            path: '/blacklist',
            component: () => import("@/pages/Blacklist.vue")
        },
        {
            path: '/cas',
            component: () => import("@/pages/CAs.vue")
        }
    ]
})

export function setupRouter(app: App<Element>) {
    app.use(router);
}

export default router