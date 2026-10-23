import { defineStore } from "pinia";
import api from "../services/api";

export const useVideoStore = defineStore("video", {
  state: () => ({
    // Home (3 video terbaru)
    videos: [],

    // Halaman Video (semua video)
    allVideos: [],

    loading: false,

    // cache
    loaded: false,
    loadedAll: false,

    // admin
    adminVideos: [],
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
    async fetchVideos(force = false) {
      if (this.loaded && !force) return;

      this.loading = true;

      try {
        const response = await api.get("/public/video");

        this.videos = response.data.data || [];
        this.loaded = true;
      } catch (error) {
        console.error("Gagal mengambil video:", error);
      } finally {
        this.loading = false;
      }
    },

    async fetchAllVideos(force = false) {
      if (this.loadedAll && !force) return;

      this.loading = true;

      try {
        const response = await api.get("/public/video/list");

        this.allVideos = response.data.data || [];
        this.loadedAll = true;
      } catch (error) {
        console.error("Gagal mengambil semua video:", error);
      } finally {
        this.loading = false;
      }
    },

    async fetchAdminVideos() {
      this.adminLoading = true;

      try {
        const response = await api.get("/video", {
          params: {
            search: this.filters.q,
            page: this.pagination.page,
            limit: this.pagination.limit,
            sort_by: this.sortBy,
            sort_dir: this.sortDir,
          },
        });

        this.adminVideos = response.data.data || [];
        this.pagination.total = response.data.total || 0;
        this.pagination.page = response.data.page || 1;
        this.pagination.limit = response.data.limit || 10;
        this.pagination.total_page = Math.max(
          1,
          Math.ceil(this.pagination.total / this.pagination.limit),
        );
      } catch (error) {
        console.error("Gagal mengambil video admin:", error);
      } finally {
        this.adminLoading = false;
      }
    },

    changeLimit(limit) {
      this.pagination.limit = Number(limit);
      this.pagination.page = 1;
      this.fetchAdminVideos();
    },

    async fetchVideoById(id) {
      this.adminLoading = true;

      try {
        const response = await api.get(`/video/${id}`);

        this.adminDetail = response.data.data;

        return response.data.data;
      } catch (error) {
        console.error("Gagal mengambil detil video:", error);
        throw error;
      } finally {
        this.adminLoading = false;
      }
    },

    async createVideo(payload) {
      this.adminLoading = true;

      try {
        const response = await api.post("/video", payload);

        return response.data;
      } catch (error) {
        console.error("Gagal membuat video:", error);
        throw error;
      } finally {
        this.adminLoading = false;
      }
    },

    async updateVideo(id, payload) {
      this.adminLoading = true;

      try {
        const response = await api.put(`/video/${id}`, payload);

        return response.data;
      } catch (error) {
        console.error("Gagal memperbarui video:", error);
        throw error;
      } finally {
        this.adminLoading = false;
      }
    },

    async deleteVideo(id) {
      this.adminLoading = true;

      try {
        const response = await api.delete(`/video/${id}`);

        this.adminVideos = this.adminVideos.filter((v) => v.id !== id);

        return response.data;
      } catch (error) {
        console.error("Gagal menghapus video:", error);
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
      this.fetchAdminVideos();
    },

    clearCache() {
      this.videos = [];
      this.allVideos = [];

      this.loaded = false;
      this.loadedAll = false;
    },
  },
});
