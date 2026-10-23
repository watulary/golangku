import { defineStore } from "pinia";
import api from "../services/api";

export const useFotoEkskulStore = defineStore("fotoekskul", {
  state: () => ({
    fotoekskul: [],
    loading: false,

    // cache
    loaded: false,
    loadedAll: false,

    // admin
    adminFotoEkskul: [],
    adminLoading: false,
    adminDetail: null,
    pagination: {
      page: 1,
      limit: 10,
      total: 0,
      total_page: 1,
    },
    filters: {
      q: "",
    },

    // sort
    sortBy: "waktu",
    sortDir: "desc",
  }),

  actions: {
    async fetchFotoEkskul(force = false) {
      if (this.loaded && !force) return;

      this.loading = true;

      try {
        const response = await api.get("/public/fotoekskul");

        this.fotoekskul = response.data.data || [];
        this.loaded = true;
      } catch (error) {
        console.error("Gagal mengambil foto ekskul:", error);
      } finally {
        this.loading = false;
      }
    },

    async fetchAllFotoEkskul(force = false) {
      if (this.loadedAll && !force) return;

      this.loading = true;

      try {
        const response = await api.get("/public/fotoekskul/list");

        this.fotoekskul = response.data.data || [];
        this.loadedAll = true;
      } catch (error) {
        console.error("Gagal mengambil semua foto ekskul:", error);
      } finally {
        this.loading = false;
      }
    },

    async fetchAdminFotoEkskul() {
      this.adminLoading = true;

      try {
        const response = await api.get("/fotoekskul", {
          params: {
            search: this.filters.q,
            page: this.pagination.page,
            limit: this.pagination.limit,
            sort_by: this.sortBy,
            sort_dir: this.sortDir,
          },
        });

        this.adminFotoEkskul = response.data.data || [];
        this.pagination.total = response.data.total || 0;
        this.pagination.page = response.data.page || 1;
        this.pagination.limit = response.data.limit || 10;
        this.pagination.total_page = Math.max(
          1,
          Math.ceil(this.pagination.total / this.pagination.limit),
        );
      } catch (error) {
        console.error("Gagal mengambil foto ekskul admin:", error);
      } finally {
        this.adminLoading = false;
      }
    },

    changeLimit(limit) {
      this.pagination.limit = Number(limit);
      this.pagination.page = 1;
      this.fetchAdminFotoEkskul();
    },

    async fetchFotoEkskulById(id) {
      this.adminLoading = true;

      try {
        const response = await api.get(`/fotoekskul/${id}`);

        this.adminDetail = response.data.data;

        return response.data.data;
      } catch (error) {
        console.error("Gagal mengambil detil foto ekskul:", error);
        throw error;
      } finally {
        this.adminLoading = false;
      }
    },

    async createFotoEkskul(formData) {
      this.adminLoading = true;

      try {
        const response = await api.post("/fotoekskul", formData, {
          headers: { "Content-Type": "multipart/form-data" },
        });

        return response.data;
      } catch (error) {
        console.error("Gagal membuat foto ekskul:", error);
        throw error;
      } finally {
        this.adminLoading = false;
      }
    },

    async updateFotoEkskul(id, formData) {
      this.adminLoading = true;

      try {
        const response = await api.put(`/fotoekskul/${id}`, formData, {
          headers: { "Content-Type": "multipart/form-data" },
        });

        return response.data;
      } catch (error) {
        console.error("Gagal memperbarui foto ekskul:", error);
        throw error;
      } finally {
        this.adminLoading = false;
      }
    },

    async deleteFotoEkskul(id) {
      this.adminLoading = true;

      try {
        const response = await api.delete(`/fotoekskul/${id}`);

        this.adminFotoEkskul = this.adminFotoEkskul.filter((f) => f.id !== id);

        return response.data;
      } catch (error) {
        console.error("Gagal menghapus foto ekskul:", error);
        throw error;
      } finally {
        this.adminLoading = false;
      }
    },

    toggleSort(column) {
      if (this.sortBy === column) {
        this.sortDir = this.sortDir === "asc" ? "desc" : "asc";
      } else {
        this.sortBy = column;
        this.sortDir = "asc";
      }

      this.pagination.page = 1;
      this.fetchAdminFotoEkskul();
    },

    clearCache() {
      this.loaded = false;
      this.loadedAll = false;
    },
  },
});
