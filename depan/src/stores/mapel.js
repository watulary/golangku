import { defineStore } from "pinia";
import api from "../services/api";

export const useMapelStore = defineStore("mapel", {
  state: () => ({
    list: [],
    loading: false,

    // daftar mapel utama (id_induk NULL) untuk dropdown "Induk" di form
    induk: [],

    detail: null,
    detailLoading: false,

    pagination: {
      page: 1,
      limit: 10,
      total: 0,
      total_page: 1,
    },
    filters: {
      q: "",
    },

    // "urutan" = urutan pohon dari backend (mapel utama + sub-mapelnya)
    sortBy: "urutan",
    sortDir: "asc",
  }),

  actions: {
    async fetchAll() {
      this.loading = true;

      try {
        const response = await api.get("/mapel", {
          params: {
            search: this.filters.q,
            page: this.pagination.page,
            limit: this.pagination.limit,
            sort_by: this.sortBy,
            sort_dir: this.sortDir,
          },
        });

        this.list = response.data.data || [];
        this.pagination.total = response.data.total || 0;
        this.pagination.page = response.data.page || 1;
        this.pagination.limit = response.data.limit || 10;
        this.pagination.total_page = Math.max(
          1,
          Math.ceil(this.pagination.total / this.pagination.limit),
        );
      } catch (error) {
        console.error("Gagal mengambil mapel:", error);
      } finally {
        this.loading = false;
      }
    },

    async fetchInduk() {
      try {
        const response = await api.get("/mapel/induk");

        this.induk = response.data.data || [];
      } catch (error) {
        console.error("Gagal mengambil daftar mapel induk:", error);
        this.induk = [];
      }
    },

    changeLimit(limit) {
      this.pagination.limit = Number(limit);
      this.pagination.page = 1;
      this.fetchAll();
    },

    toggleSort(column) {
      if (this.sortBy === column) {
        this.sortDir = this.sortDir === "asc" ? "desc" : "asc";
      } else {
        this.sortBy = column;
        this.sortDir = "asc";
      }

      this.pagination.page = 1;
      this.fetchAll();
    },

    async fetchById(id) {
      this.detailLoading = true;

      try {
        const response = await api.get(`/mapel/${id}`);

        this.detail = response.data.data;

        return response.data.data;
      } catch (error) {
        console.error("Gagal mengambil detil mapel:", error);
        throw error;
      } finally {
        this.detailLoading = false;
      }
    },

    async create(payload) {
      this.loading = true;

      try {
        const response = await api.post("/mapel", payload);

        await this.fetchAll();

        return response.data;
      } catch (error) {
        console.error("Gagal membuat mapel:", error);
        throw error;
      } finally {
        this.loading = false;
      }
    },

    async update(id, payload) {
      this.loading = true;

      try {
        const response = await api.put(`/mapel/${id}`, payload);

        await this.fetchAll();

        return response.data;
      } catch (error) {
        console.error("Gagal memperbarui mapel:", error);
        throw error;
      } finally {
        this.loading = false;
      }
    },

    async remove(id) {
      this.loading = true;

      try {
        const response = await api.delete(`/mapel/${id}`);

        await this.fetchAll();

        return response.data;
      } catch (error) {
        console.error("Gagal menghapus mapel:", error);
        throw error;
      } finally {
        this.loading = false;
      }
    },

    async toggleStatus(id) {
      this.loading = true;

      try {
        const response = await api.put(`/mapel/${id}/status`);

        await this.fetchAll();

        return response.data;
      } catch (error) {
        console.error("Gagal mengubah status mapel:", error);
        throw error;
      } finally {
        this.loading = false;
      }
    },
  },
});
