import { defineStore } from "pinia";
import api from "../services/api";

export const useBeritaStore = defineStore("berita", {
  state: () => ({
    berita: [],
    semuaBerita: [],
    detailBerita: null,
    informasiDokumen: [],
    informasiBerkala: [],
    informasiSetiapSaat: [],
    informasiSertaMerta: [],
    informasiKontak: [],
    loading: false,

    // cache
    loaded: false,
    loadedAll: false,
    loadedInformasiDokumen: false,
    loadedInformasiBerkala: false,
    loadedInformasiSetiapSaat: false,
    loadedInformasiSertaMerta: false,
    loadedInformasiKontak: false,

    // admin
    adminBerita: [],
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
    async fetchBerita(force = false) {
      if (this.loaded && !force) return;

      this.loading = true;

      try {
        const response = await api.get("/public/berita");

        this.berita = response.data.data || [];
        this.loaded = true;
      } catch (error) {
        console.error("Gagal mengambil berita:", error);
      } finally {
        this.loading = false;
      }
    },

    async fetchAllBerita(force = false) {
      if (this.loadedAll && !force) return;

      this.loading = true;

      try {
        const response = await api.get("/public/berita/list");

        this.semuaBerita = response.data.data || [];
        this.loadedAll = true;
      } catch (error) {
        console.error("Gagal mengambil berita:", error);
      } finally {
        this.loading = false;
      }
    },

    async fetchDetailBerita(slug) {
      this.loading = true;

      try {
        const response = await api.get(`/public/berita/detil/${slug}`);

        this.detailBerita = response.data.data;
      } catch (error) {
        console.error("Gagal mengambil detil berita:", error);
      } finally {
        this.loading = false;
      }
    },

    async incrementHits(idBerita) {
      try {
        await api.post(`/public/berita/hits/${idBerita}`);
      } catch (error) {
        console.error("Gagal menambah hits berita:", error);
      }
    },

    async fetchInformasiDokumen(force = false) {
      if (this.loadedInformasiDokumen && !force) return;

      this.loading = true;

      try {
        const response = await api.get("/public/informasi-dokumen");

        this.informasiDokumen = response.data.data || [];

        this.loadedInformasiDokumen = true;
      } catch (error) {
        console.error("Gagal mengambil informasi dokumen:", error);
      } finally {
        this.loading = false;
      }
    },

    async fetchInformasiSetiapSaat(force = false) {
      if (this.loadedInformasiSetiapSaat && !force) return;

      this.loading = true;

      try {
        const response = await api.get("/public/informasi-setiap-saat");

        this.informasiSetiapSaat = response.data.data || [];

        this.loadedInformasiSetiapSaat = true;
      } catch (error) {
        console.error("Gagal mengambil informasi setiap saat:", error);
      } finally {
        this.loading = false;
      }
    },

    async fetchInformasiBerkala(force = false) {
      if (this.loadedInformasiBerkala && !force) return;

      this.loading = true;

      try {
        const response = await api.get("/public/informasi-berkala");

        this.informasiBerkala = response.data.data || [];

        this.loadedInformasiBerkala = true;
      } catch (error) {
        console.error("Gagal mengambil informasi berkala:", error);
      } finally {
        this.loading = false;
      }
    },

    async fetchInformasiSertaMerta(force = false) {
      if (this.loadedInformasiSertaMerta && !force) return;

      this.loading = true;

      try {
        const response = await api.get("/public/informasi-serta-merta");

        this.informasiSertaMerta = response.data.data || [];

        this.loadedInformasiSertaMerta = true;
      } catch (error) {
        console.error("Gagal mengambil informasi serta merta:", error);
      } finally {
        this.loading = false;
      }
    },

    async fetchInformasiKontak(force = false) {
      if (this.loadedInformasiKontak && !force) return;

      this.loading = true;

      try {
        const response = await api.get("/public/informasi-kontak");

        this.informasiKontak = response.data.data || [];

        this.loadedInformasiKontak = true;
      } catch (error) {
        console.error("Gagal mengambil informasi kontak:", error);
      } finally {
        this.loading = false;
      }
    },

    async fetchAdminBerita() {
      this.adminLoading = true;

      try {
        const response = await api.get("/berita", {
          params: {
            search: this.filters.q,
            page: this.pagination.page,
            limit: this.pagination.limit,
            sort_by: this.sortBy,
            sort_dir: this.sortDir,
          },
        });

        this.adminBerita = response.data.data || [];
        this.pagination.total = response.data.total || 0;
        this.pagination.page = response.data.page || 1;
        this.pagination.limit = response.data.limit || 10;
        this.pagination.total_page = Math.max(
          1,
          Math.ceil(this.pagination.total / this.pagination.limit),
        );
      } catch (error) {
        console.error("Gagal mengambil berita admin:", error);
      } finally {
        this.adminLoading = false;
      }
    },

    changeLimit(limit) {
      this.pagination.limit = Number(limit);
      this.pagination.page = 1;
      this.fetchAdminBerita();
    },

    async fetchBeritaById(id) {
      this.adminLoading = true;

      try {
        const response = await api.get(`/berita/${id}`);

        this.adminDetail = response.data.data;

        return response.data.data;
      } catch (error) {
        console.error("Gagal mengambil detil berita:", error);
        throw error;
      } finally {
        this.adminLoading = false;
      }
    },

    async createBerita(formData) {
      this.adminLoading = true;

      try {
        const response = await api.post("/berita", formData, {
          headers: { "Content-Type": "multipart/form-data" },
        });

        this.pagination.page = 1;

        return response.data;
      } catch (error) {
        console.error("Gagal membuat berita:", error);
        throw error;
      } finally {
        this.adminLoading = false;
      }
    },

    async updateBerita(id, formData) {
      this.adminLoading = true;

      try {
        const response = await api.put(`/berita/${id}`, formData, {
          headers: { "Content-Type": "multipart/form-data" },
        });

        return response.data;
      } catch (error) {
        console.error("Gagal memperbarui berita:", error);
        throw error;
      } finally {
        this.adminLoading = false;
      }
    },

    async deleteBerita(id) {
      this.adminLoading = true;

      try {
        const response = await api.delete(`/berita/${id}`);

        this.adminBerita = this.adminBerita.filter((b) => b.id_berita !== id);

        return response.data;
      } catch (error) {
        console.error("Gagal menghapus berita:", error);
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
      this.fetchAdminBerita();
    },

    clearCache() {
      this.loaded = false;
      this.loadedAll = false;
      this.loadedInformasiDokumen = false;
      this.loadedInformasiBerkala = false;
      this.loadedInformasiSetiapSaat = false;
      this.loadedInformasiSertaMerta = false;
      this.loadedInformasiKontak = false;
    },
  },
});
