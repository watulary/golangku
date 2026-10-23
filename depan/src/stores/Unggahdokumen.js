import { defineStore } from "pinia";
import api from "../services/api";

export const useUnggahdokumenStore = defineStore("unggahdokumen", {
  state: () => ({
    // Publik
    data: [],
    loading: false,
    loaded: false,

    // Admin
    adminLoading: false,
    adminUnggahdokumen: [],
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
    sortBy: "waktu",
    sortDir: "desc",
  }),

  actions: {
    async fetchData(force = false) {
      if (this.loaded && !force) return;

      this.loading = true;

      try {
        const res = await api.get("/public/unduh");

        this.data = res.data.data || [];
        this.loaded = true;
      } finally {
        this.loading = false;
      }
    },

    async fetchAdminUnggahdokumen() {
      this.adminLoading = true;

      try {
        const res = await api.get("/unggahdokumen", {
          params: {
            q: this.filters.q,
            sortBy: this.sortBy,
            sortDir: this.sortDir,
            page: this.pagination.page,
            limit: this.pagination.limit,
          },
        });

        this.adminUnggahdokumen = res.data.data || [];
        this.pagination.total = res.data.pagination?.total || 0;
      } finally {
        this.adminLoading = false;
      }
    },

    changeLimit(limit) {
      this.pagination.limit = limit;
      this.pagination.page = 1;
      this.fetchAdminUnggahdokumen();
    },

    async fetchUnggahdokumenById(id) {
      this.adminLoading = true;

      try {
        const res = await api.get(`/unggahdokumen/${id}`);

        this.adminDetail = res.data.data;
        return this.adminDetail;
      } finally {
        this.adminLoading = false;
      }
    },

    async createUnggahdokumen(data) {
      const res = await api.post("/unggahdokumen", data);
      return res.data;
    },

    async updateUnggahdokumen(id, data) {
      const res = await api.put(`/unggahdokumen/${id}`, data);
      return res.data;
    },

    async deleteUnggahdokumen(id) {
      const res = await api.delete(`/unggahdokumen/${id}`);

      this.adminUnggahdokumen = this.adminUnggahdokumen.filter(
        (item) => item.id !== id,
      );

      return res.data;
    },

    toggleSort(column) {
      if (this.sortBy === column) {
        this.sortDir = this.sortDir === "asc" ? "desc" : "asc";
      } else {
        this.sortBy = column;
        this.sortDir = "asc";
      }
      this.fetchAdminUnggahdokumen();
    },

    clearCache() {
      this.data = [];
      this.loaded = false;
    },
  },
});
