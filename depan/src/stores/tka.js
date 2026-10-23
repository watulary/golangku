import { defineStore } from "pinia";
import api from "../services/api";

function triggerDownload(res) {
  const disposition = res.headers["content-disposition"] || "";
  const match = disposition.match(/filename="?([^"]+)"?/);
  const filename = match ? match[1] : "data.xlsx";
  const url = window.URL.createObjectURL(new Blob([res.data]));
  const link = document.createElement("a");
  link.href = url;
  link.setAttribute("download", filename);
  document.body.appendChild(link);
  link.click();
  link.remove();
  window.URL.revokeObjectURL(url);
}

export const useTkaStore = defineStore("tka", {
  state: () => ({
    // ==== Data TKA (admin-only, import-only — tanpa create/update/delete per-nis) ====
    adminTka: [],
    adminLoading: false,
    pagination: {
      page: 1,
      limit: 10,
      total: 0,
      total_page: 1,
    },
    filters: {
      q: "",
    },
    sortBy: "nama_rombel",
    sortDir: "asc",
    importLoading: false,
    importErrors: [],

    // ==== Pengaturan Pengumuman (tka_waktu, single-row) ====
    waktu: null,
    waktuLoading: false,

    // ==== pesertadidik ====
    siswaData: null,
    siswaLoading: false,
  }),
  actions: {
    // ==== Data TKA ====
    async fetchAdminTka() {
      this.adminLoading = true;
      try {
        const res = await api.get("/tka", {
          params: {
            search: this.filters.q,
            sort_by: this.sortBy,
            sort_dir: this.sortDir,
            limit: this.pagination.limit,
            page: this.pagination.page,
          },
        });
        this.adminTka = res.data.data || [];
        this.pagination.total = res.data.pagination?.total || 0;
        this.pagination.total_page = Math.max(
          1,
          Math.ceil(this.pagination.total / this.pagination.limit),
        );
        if (this.pagination.page > this.pagination.total_page) {
          this.pagination.page = this.pagination.total_page;
          await this.fetchAdminTka();
          return;
        }
      } catch (error) {
        console.error("Gagal mengambil data TKA", error);
      } finally {
        this.adminLoading = false;
      }
    },

    // ==== Bulk actions ====
    async bulkDelete(nisList) {
      this.adminLoading = true;
      try {
        const res = await api.post("/tka/bulk-delete", {
          nis: nisList,
        });
        await this.fetchAdminTka();
        return res.data;
      } catch (error) {
        console.error("Gagal menghapus data terpilih", error);
        throw error;
      } finally {
        this.adminLoading = false;
      }
    },

    // ==== Import / Export / Template ====
    async importTka(file) {
      this.importLoading = true;
      this.importErrors = [];
      try {
        const formData = new FormData();
        formData.append("file_excel", file);
        const res = await api.post("/tka/import", formData, {
          headers: { "Content-Type": "multipart/form-data" },
        });
        await this.fetchAdminTka();
        return res.data;
      } catch (error) {
        if (error.response?.data?.errors) {
          this.importErrors = error.response.data.errors;
        }
        console.error("Gagal mengimpor data TKA", error);
        throw error;
      } finally {
        this.importLoading = false;
      }
    },
    async downloadTemplate() {
      const res = await api.get("/tka/template", { responseType: "blob" });
      triggerDownload(res);
    },
    async exportTka() {
      const res = await api.get("/tka/export", {
        params: { search: this.filters.q },
        responseType: "blob",
      });
      triggerDownload(res);
    },

    // ==== Sort & filter helper ====
    toggleSort(column) {
      if (this.sortBy === column) {
        this.sortDir = this.sortDir === "asc" ? "desc" : "asc";
      } else {
        this.sortBy = column;
        this.sortDir = "asc";
      }
      this.pagination.page = 1;
      this.fetchAdminTka();
    },
    changeLimit(limit) {
      this.pagination.limit = Number(limit);
      this.pagination.page = 1;
      this.fetchAdminTka();
    },

    // ==== Pengaturan Pengumuman (tka_waktu) ====
    async fetchWaktu() {
      this.waktuLoading = true;
      try {
        const res = await api.get("/tka/waktu");
        this.waktu = res.data.data;
        return this.waktu;
      } catch (error) {
        console.error("Gagal mengambil pengaturan pengumuman TKA", error);
        throw error;
      } finally {
        this.waktuLoading = false;
      }
    },
    async updateWaktu(payload) {
      this.waktuLoading = true;
      try {
        const res = await api.put("/tka/waktu", payload);
        await this.fetchWaktu();
        return res.data;
      } catch (error) {
        console.error("Gagal memperbarui pengaturan pengumuman TKA", error);
        throw error;
      } finally {
        this.waktuLoading = false;
      }
    },

    // ==== pesertadidik ====
    async fetchTkaPesertadidik() {
      this.siswaLoading = true;
      try {
        const res = await api.get("/tka/me");
        this.siswaData = res.data.data;
        return this.siswaData;
      } catch (error) {
        console.error("Gagal mengambil data TKA peserta didik", error);
        throw error;
      } finally {
        this.siswaLoading = false;
      }
    },
  },
});
