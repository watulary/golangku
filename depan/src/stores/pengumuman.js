import { defineStore } from "pinia";
import api from "../services/api";

export const usePengumumanStore = defineStore("pengumuman", {
  state: () => ({
    pengumuman: [],
    semuaPengumuman: [],
    detailPengumuman: null,
    loading: false,

    // cache
    loaded: false,
    loadedAll: false,

    // admin
    adminPengumuman: [],
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
    sortBy: "tanggal",
    sortDir: "desc",
  }),

  actions: {
    async fetchPengumuman(force = false) {
      if (this.loaded && !force) return;

      this.loading = true;

      try {
        const response = await api.get("/public/pengumuman");

        this.pengumuman = response.data.data || [];
        this.loaded = true;
      } catch (error) {
        console.error("Gagal mengambil pengumuman:", error);
      } finally {
        this.loading = false;
      }
    },

    async fetchAllPengumuman(force = false) {
      if (this.loadedAll && !force) return;

      this.loading = true;

      try {
        const response = await api.get("/public/pengumuman/list");

        this.semuaPengumuman = response.data.data || [];
        this.loadedAll = true;
      } catch (error) {
        console.error("Gagal mengambil pengumuman:", error);
      } finally {
        this.loading = false;
      }
    },

    async fetchDetailPengumuman(slug) {
      this.loading = true;

      try {
        const response = await api.get(`/public/pengumuman/detil/${slug}`);

        this.detailPengumuman = response.data.data;
      } catch (error) {
        console.error("Gagal mengambil detil pengumuman:", error);
      } finally {
        this.loading = false;
      }
    },

    async incrementHits(idPengumuman) {
      try {
        await api.post(`/public/pengumuman/hits/${idPengumuman}`);
      } catch (error) {
        console.error("Gagal menambah hits pengumuman:", error);
      }
    },

    async fetchAdminPengumuman() {
      this.adminLoading = true;

      try {
        const response = await api.get("/pengumuman", {
          params: {
            search: this.filters.q,
            page: this.pagination.page,
            limit: this.pagination.limit,
            sort_by: this.sortBy,
            sort_dir: this.sortDir,
          },
        });

        this.adminPengumuman = response.data.data || [];
        this.pagination.total = response.data.total || 0;
        this.pagination.page = response.data.page || 1;
        this.pagination.limit = response.data.limit || 10;
        this.pagination.total_page = Math.max(
          1,
          Math.ceil(this.pagination.total / this.pagination.limit),
        );
      } catch (error) {
        console.error("Gagal mengambil pengumuman admin:", error);
      } finally {
        this.adminLoading = false;
      }
    },

    changeLimit(limit) {
      this.pagination.limit = Number(limit);
      this.pagination.page = 1;
      this.fetchAdminPengumuman();
    },

    async fetchPengumumanById(id) {
      this.adminLoading = true;

      try {
        const response = await api.get(`/pengumuman/${id}`);

        this.adminDetail = response.data.data;

        return response.data.data;
      } catch (error) {
        console.error("Gagal mengambil detil pengumuman:", error);
        throw error;
      } finally {
        this.adminLoading = false;
      }
    },

    async createPengumuman(formData) {
      this.adminLoading = true;

      try {
        const response = await api.post("/pengumuman", formData, {
          headers: { "Content-Type": "multipart/form-data" },
        });

        return response.data;
      } catch (error) {
        console.error("Gagal membuat pengumuman:", error);
        throw error;
      } finally {
        this.adminLoading = false;
      }
    },

    async updatePengumuman(id, formData) {
      this.adminLoading = true;

      try {
        const response = await api.put(`/pengumuman/${id}`, formData, {
          headers: { "Content-Type": "multipart/form-data" },
        });

        return response.data;
      } catch (error) {
        console.error("Gagal memperbarui pengumuman:", error);
        throw error;
      } finally {
        this.adminLoading = false;
      }
    },

    async deletePengumuman(id) {
      this.adminLoading = true;

      try {
        const response = await api.delete(`/pengumuman/${id}`);

        this.adminPengumuman = this.adminPengumuman.filter(
          (pm) => pm.id_pengumuman !== id,
        );

        return response.data;
      } catch (error) {
        console.error("Gagal menghapus pengumuman:", error);
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
      this.fetchAdminPengumuman();
    },

    clearCache() {
      this.loaded = false;
      this.loadedAll = false;
    },
  },
});
