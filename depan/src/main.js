import { createApp } from "vue";
import { createPinia } from "pinia";
import App from "./App.vue";
import router from "./router";
import { setupInterceptors } from "./services/setupInterceptors";
import { useAuthStore } from "./stores/auth";
import "./style.css";

const app = createApp(App);

// Buat Pinia terlebih dahulu
const pinia = createPinia();

app.use(pinia);
app.use(router);

// Ambil auth store setelah Pinia terpasang
const authStore = useAuthStore(pinia);

// Pasang interceptor global
setupInterceptors(authStore, router);

app.mount("#app");
