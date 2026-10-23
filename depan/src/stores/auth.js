import { defineStore } from "pinia";
import api from "../services/api";
import { getErrorMessage } from "../utils/error";
import { useStatistikStore } from "./statistik";

export const useAuthStore = defineStore("auth", {
  state: () => ({
    user: null,
    token: localStorage.getItem("token"),
    loading: false,

    // cache
    loaded: false,

    // error
    error: null,

    changingPassword: false,
  }),

  getters: {
    isAuthenticated: (state) => !!state.token,

    // Getter utama untuk seluruh sistem role
    role: (state) => state.user?.level ?? null,

    displayRole: (state) => {
      const roles = {
        Admin: "Super Admin",
        Staf: "Admin",
        Pegawai: "Pegawai",
        Pesertadidik: "Peserta Didik",
      };

      return roles[state.user?.level] ?? state.user?.level ?? "";
    },

    // ===============================
    // Password
    // ===============================

    passwordChanged: (state) => state.user?.password_changed ?? true,

    mustChangePassword: (state) =>
      !!state.user && state.user.password_changed === false,

    isAdmin: (state) => state.user?.level === "Admin",

    isStaf: (state) => state.user?.level === "Staf",

    isPegawai: (state) => state.user?.level === "Pegawai",

    isPesertadidik: (state) => state.user?.level === "Pesertadidik",
  },

  actions: {
    async login(email, password) {
      this.loading = true;
      this.error = null;

      try {
        const response = await api.post("/public/login", {
          email,
          password,
        });

        this.token = response.data.token;
        localStorage.setItem("token", this.token);

        // Ambil data user terbaru
        await this.fetchMe(true);

        return {
          success: true,
        };
      } catch (error) {
        this.logout();

        this.error = getErrorMessage(error);

        return {
          success: false,
          message: this.error,
        };
      } finally {
        this.loading = false;
      }
    },

    async fetchMe(force = false) {
      if (this.loaded && !force) return;

      if (!this.token) return;

      this.loading = true;
      this.error = null;

      try {
        const response = await api.get("/public/me");

        this.user = response.data.user;
        this.loaded = true;
      } catch (error) {
        this.error = getErrorMessage(error);

        throw error;
      } finally {
        this.loading = false;
      }
    },

    async changePassword(payload) {
      this.changingPassword = true;
      this.error = null;

      try {
        const { data } = await api.post("/pengguna/change-password", payload);

        // Refresh data user agar password_changed ikut berubah
        await this.fetchMe();

        return data;
      } catch (error) {
        this.error =
          error.response?.data?.message || "Gagal mengubah password.";

        throw error;
      } finally {
        this.changingPassword = false;
      }
    },

    async updateFotoProfil(formData) {
      this.error = null;

      try {
        const { data } = await api.put("/pengguna/foto", formData);

        // Refresh data user agar foto_url ikut berubah
        await this.fetchMe(true);

        return data;
      } catch (error) {
        this.error =
          error.response?.data?.message || "Gagal mengubah foto profil.";

        throw error;
      }
    },

    async deleteFotoProfil() {
      this.error = null;

      try {
        const { data } = await api.delete("/pengguna/foto");

        // Refresh data user agar foto_url ikut berubah
        await this.fetchMe(true);

        return data;
      } catch (error) {
        this.error =
          error.response?.data?.message || "Gagal menghapus foto profil.";

        throw error;
      }
    },

    logout() {
      this.user = null;
      this.token = null;
      this.loaded = false;
      this.error = null;

      localStorage.removeItem("token");

      // Statistik dashboard difilter per pengguna (kecuali Super Admin) —
      // cache-nya harus ikut dibersihkan saat logout, supaya akun berikutnya
      // yang login (tanpa full refresh browser) tidak melihat angka statistik
      // sisa dari akun sebelumnya.
      useStatistikStore().clearCache();
    },

    clearCache() {
      this.loaded = false;
      this.error = null;
    },
  },
});
