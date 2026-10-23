import { defineStore } from "pinia";
import api from "../services/api";

export const useEkskulStore = defineStore("ekskul", {
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
      aktif: "",
    },

    sortBy: "ekskul",
    sortDir: "asc",

    pegawaiOptions: [],
    pegawaiSearching: false,

    pesertaDidikOptions: [],
    pesertaDidikSearching: false,

    anggotaList: [],
    anggotaLoading: false,
    anggotaPagination: {
      page: 1,
      limit: 10,
      total: 0,
      total_page: 1,
    },
    anggotaFilters: {
      q: "",
      idKelas: 0,
      idRombel: 0,
    },

    siswaList: [],
    siswaLoading: false,
  }),

  actions: {
    async fetchAll() {
      this.loading = true;

      try {
        const response = await api.get("/ekskul", {
          params: {
            search: this.filters.q,
            aktif: this.filters.aktif !== "" ? this.filters.aktif : undefined,
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
        console.error("Gagal mengambil ekskul:", error);
      } finally {
        this.loading = false;
      }
    },

    async exportEkskul() {
      const response = await api.get("/ekskul/export", {
        responseType: "blob",
      });

      const disposition = response.headers["content-disposition"];
      let filename = "daftar-ekskul.xlsx";

      if (disposition) {
        const match = disposition.match(/filename="?([^"]+)"?/);
        if (match) {
          filename = match[1];
        }
      }

      const url = window.URL.createObjectURL(new Blob([response.data]));
      const link = document.createElement("a");

      link.href = url;
      link.download = filename;

      document.body.appendChild(link);
      link.click();
      link.remove();

      window.URL.revokeObjectURL(url);
    },

    async exportAnggota() {
      const response = await api.get("/ekskul/anggota/export", {
        responseType: "blob",
      });

      const disposition = response.headers["content-disposition"];
      let filename = "daftar-anggota-ekskul.xlsx";

      if (disposition) {
        const match = disposition.match(/filename="?([^"]+)"?/);
        if (match) {
          filename = match[1];
        }
      }

      const url = window.URL.createObjectURL(new Blob([response.data]));
      const link = document.createElement("a");

      link.href = url;
      link.download = filename;

      document.body.appendChild(link);
      link.click();
      link.remove();

      window.URL.revokeObjectURL(url);
    },

    async exportAnggotaByEkskul(idEkskul) {
      const response = await api.get(`/ekskul/${idEkskul}/anggota/export`, {
        params: {
          search: this.anggotaFilters.q,
          id_kelas: this.anggotaFilters.idKelas || undefined,
          id_rombel: this.anggotaFilters.idRombel || undefined,
        },
        responseType: "blob",
      });

      const disposition = response.headers["content-disposition"];
      let filename = "daftar-anggota-ekskul.xlsx";

      if (disposition) {
        const match = disposition.match(/filename="?([^"]+)"?/);
        if (match) {
          filename = match[1];
        }
      }

      const url = window.URL.createObjectURL(new Blob([response.data]));
      const link = document.createElement("a");

      link.href = url;
      link.download = filename;

      document.body.appendChild(link);
      link.click();
      link.remove();

      window.URL.revokeObjectURL(url);
    },

    changeLimit(limit) {
      this.pagination.limit = Number(limit);
      this.pagination.page = 1;
      this.fetchAll();
    },

    changeFilterAktif(aktif) {
      this.filters.aktif = aktif;
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
        const response = await api.get(`/ekskul/${id}`);

        this.detail = response.data.data;

        return response.data.data;
      } catch (error) {
        console.error("Gagal mengambil detil ekskul:", error);
        throw error;
      } finally {
        this.detailLoading = false;
      }
    },

    // ================= PESERTA DIDIK =================

    async fetchMe() {
      this.siswaLoading = true;

      try {
        const response = await api.get("/ekskul/me");

        this.siswaList = response.data.data || [];
      } catch (error) {
        console.error("Gagal mengambil ekskul saya:", error);
      } finally {
        this.siswaLoading = false;
      }
    },

    async searchPegawai(query) {
      this.pegawaiSearching = true;

      try {
        const response = await api.get("/ekskul/pegawai/search", {
          params: { q: query },
        });

        this.pegawaiOptions = response.data.data || [];
      } catch (error) {
        console.error("Gagal mencari pegawai:", error);
      } finally {
        this.pegawaiSearching = false;
      }
    },

    async create(payload) {
      this.loading = true;

      try {
        const response = await api.post("/ekskul", payload);

        await this.fetchAll();

        return response.data;
      } catch (error) {
        console.error("Gagal membuat ekskul:", error);
        throw error;
      } finally {
        this.loading = false;
      }
    },

    async update(id, payload) {
      this.loading = true;

      try {
        const response = await api.put(`/ekskul/${id}`, payload);

        await this.fetchAll();

        return response.data;
      } catch (error) {
        console.error("Gagal memperbarui ekskul:", error);
        throw error;
      } finally {
        this.loading = false;
      }
    },

    async remove(id) {
      this.loading = true;

      try {
        const response = await api.delete(`/ekskul/${id}`);

        await this.fetchAll();

        return response.data;
      } catch (error) {
        console.error("Gagal menghapus ekskul:", error);
        throw error;
      } finally {
        this.loading = false;
      }
    },

    // ================= ANGGOTA EKSKUL =================

    async fetchAnggota(idEkskul) {
      this.anggotaLoading = true;

      try {
        const response = await api.get(`/ekskul/${idEkskul}/anggota`, {
          params: {
            search: this.anggotaFilters.q,
            page: this.anggotaPagination.page,
            limit: this.anggotaPagination.limit,
            id_kelas: this.anggotaFilters.idKelas || undefined,
            id_rombel: this.anggotaFilters.idRombel || undefined,
          },
        });

        this.anggotaList = response.data.data || [];
        this.anggotaPagination.total = response.data.total || 0;
        this.anggotaPagination.page = response.data.page || 1;
        this.anggotaPagination.limit = response.data.limit || 10;
        this.anggotaPagination.total_page = Math.max(
          1,
          Math.ceil(
            this.anggotaPagination.total / this.anggotaPagination.limit,
          ),
        );
      } catch (error) {
        console.error("Gagal mengambil anggota ekskul:", error);
      } finally {
        this.anggotaLoading = false;
      }
    },

    async changeLimitAnggota(limit, idEkskul) {
      this.anggotaPagination.limit = Number(limit);
      this.anggotaPagination.page = 1;
      await this.fetchAnggota(idEkskul);
    },

    async changeFilterKelasAnggota(idKelas, idEkskul) {
      this.anggotaFilters.idKelas = Number(idKelas);
      this.anggotaPagination.page = 1;
      await this.fetchAnggota(idEkskul);
    },

    async changeFilterRombelAnggota(idRombel, idEkskul) {
      this.anggotaFilters.idRombel = Number(idRombel);
      this.anggotaPagination.page = 1;
      await this.fetchAnggota(idEkskul);
    },

    async bulkRemoveAnggota(ids, idEkskul) {
      const { data } = await api.post("/ekskul/anggota/bulk-delete", { ids });

      await this.fetchAnggota(idEkskul);

      return data;
    },

    async downloadTemplateAnggota() {
      const response = await api.get("/ekskul/anggota/template", {
        responseType: "blob",
      });

      const disposition = response.headers["content-disposition"];
      let filename = "template-anggota-ekskul.xlsx";

      if (disposition) {
        const match = disposition.match(/filename="?([^"]+)"?/);
        if (match) {
          filename = match[1];
        }
      }

      const url = window.URL.createObjectURL(new Blob([response.data]));
      const link = document.createElement("a");

      link.href = url;
      link.download = filename;

      document.body.appendChild(link);
      link.click();
      link.remove();

      window.URL.revokeObjectURL(url);
    },

    async importAnggota(idEkskul, file) {
      const formData = new FormData();
      formData.append("file", file);

      const { data } = await api.post(
        `/ekskul/${idEkskul}/anggota/import`,
        formData,
        {
          headers: { "Content-Type": "multipart/form-data" },
        },
      );

      // Refresh tabel setelah upload berhasil
      await this.fetchAnggota(idEkskul);

      return data;
    },

    async searchPesertaDidik(query) {
      this.pesertaDidikSearching = true;

      try {
        const response = await api.get("/ekskul/pesertadidik/search", {
          params: { q: query },
        });

        this.pesertaDidikOptions = response.data.data || [];
      } catch (error) {
        console.error("Gagal mencari peserta didik:", error);
      } finally {
        this.pesertaDidikSearching = false;
      }
    },

    async addAnggota(idEkskul, nis) {
      this.anggotaLoading = true;

      try {
        const response = await api.post(`/ekskul/${idEkskul}/anggota`, {
          nis,
        });

        await this.fetchAnggota(idEkskul);

        return response.data;
      } catch (error) {
        console.error("Gagal menambah anggota ekskul:", error);
        throw error;
      } finally {
        this.anggotaLoading = false;
      }
    },

    async removeAnggota(idEkskul, idAnggota) {
      this.anggotaLoading = true;

      try {
        const response = await api.delete(`/ekskul/anggota/${idAnggota}`);

        await this.fetchAnggota(idEkskul);

        return response.data;
      } catch (error) {
        console.error("Gagal menghapus anggota ekskul:", error);
        throw error;
      } finally {
        this.anggotaLoading = false;
      }
    },
  },
});
