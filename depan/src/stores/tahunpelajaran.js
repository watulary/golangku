import { defineStore } from "pinia";
import api from "../services/api";

export const useTahunPelajaranStore = defineStore("tahunpelajaran", {
  state: () => ({
    list: [],
    loading: false,

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

    sortBy: "nama",
    sortDir: "desc",
  }),

  getters: {
    aktif: (state) => state.list.find((t) => t.aktif === 1) || null,
  },

  actions: {
    async fetchAll() {
      this.loading = true;

      try {
        const response = await api.get("/tahunpelajaran", {
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
        console.error("Gagal mengambil tahun pelajaran:", error);
      } finally {
        this.loading = false;
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
        const response = await api.get(`/tahunpelajaran/${id}`);

        this.detail = response.data.data;

        return response.data.data;
      } catch (error) {
        console.error("Gagal mengambil detil tahun pelajaran:", error);
        throw error;
      } finally {
        this.detailLoading = false;
      }
    },

    async create(payload) {
      this.loading = true;

      try {
        const response = await api.post("/tahunpelajaran", payload);

        await this.fetchAll();

        return response.data;
      } catch (error) {
        console.error("Gagal membuat tahun pelajaran:", error);
        throw error;
      } finally {
        this.loading = false;
      }
    },

    async update(id, payload) {
      this.loading = true;

      try {
        const response = await api.put(`/tahunpelajaran/${id}`, payload);

        await this.fetchAll();

        return response.data;
      } catch (error) {
        console.error("Gagal memperbarui tahun pelajaran:", error);
        throw error;
      } finally {
        this.loading = false;
      }
    },

    async remove(id) {
      this.loading = true;

      try {
        const response = await api.delete(`/tahunpelajaran/${id}`);

        await this.fetchAll();

        return response.data;
      } catch (error) {
        console.error("Gagal menghapus tahun pelajaran:", error);
        throw error;
      } finally {
        this.loading = false;
      }
    },

    async setActive(id) {
      this.loading = true;

      try {
        const response = await api.put(`/tahunpelajaran/${id}/aktifkan`);

        await this.fetchAll();

        return response.data;
      } catch (error) {
        console.error("Gagal mengaktifkan tahun pelajaran:", error);
        throw error;
      } finally {
        this.loading = false;
      }
    },
  },
});
