import { createApp } from 'vue'
import App from './App.vue'

const app = createApp(App);

import router, { setupRouter } from "@/router";
setupRouter(app);

import { setupAntd } from "@/libs/antdv";
setupAntd(app);

import { setupStore } from "@/store";
setupStore(app)

// the router has resolved all async enter hooks 
// and async components that are associated with the initial route.
router.isReady().then(() => {
    app.mount("#app");
});