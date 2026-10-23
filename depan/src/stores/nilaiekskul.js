import { defineStore } from "pinia";
import api from "../services/api";

function snapshot(list) {
  return list.map((s) => ({
    nis: s.nis,
    nilai: (s.nilai || "").trim(),
    keterangan: (s.keterangan || "").trim(),
  }));
}

export const useNilaiEkskulStore = defineStore("nilaiekskul", {
  state: () => ({
    daftar: [],
    original: "",
    loading: false,
    submitting: false,
    importing: false,
    templateLoading: false,

    // ekskul aktif untuk dropdown
    ekskulOptions: [],
  }),

  actions: {
    async fetchEkskul() {
      try {
        const response = await api.get("/nilai-ekskul/ekskul");
        this.ekskulOptions = response.data.data || [];
      } catch (error) {
        console.error("Gagal mengambil ekstrakurikuler:", error);
        this.ekskulOptions = [];
      }
    },

    async fetchDaftar(params) {
      this.loading = true;
      try {
        const response = await api.get("/nilai-ekskul/daftar", { params });
        this.daftar = response.data.data || [];
        this.original = JSON.stringify(snapshot(this.daftar));
      } catch (error) {
        this.daftar = [];
        this.original = "";
        throw error;
      } finally {
        this.loading = false;
      }
    },

    clear() {
      this.daftar = [];
      this.original = "";
    },

    hasChanges() {
      return JSON.stringify(snapshot(this.daftar)) !== this.original;
    },

    // pesan kesalahan pertama, atau "" kalau semua valid
    validate() {
      for (const s of this.daftar) {
        if ((s.nilai || "").trim().length > 20) {
          return `Nilai ${s.nama_siswa} maksimal 20 karakter.`;
        }
        if ((s.keterangan || "").trim().length > 255) {
          return `Keterangan ${s.nama_siswa} maksimal 255 karakter.`;
        }
      }
      return "";
    },

    async submitBulk(payload) {
      this.submitting = true;
      try {
        const items = this.daftar.map((s) => ({
          nis: s.nis,
          nilai: s.nilai || "",
          keterangan: s.keterangan || "",
        }));
        const response = await api.post("/nilai-ekskul/bulk", {
          ...payload,
          items,
        });
        return response.data;
      } finally {
        this.submitting = false;
      }
    },

    async downloadTemplate(params) {
      this.templateLoading = true;
      try {
        const response = await api.get("/nilai-ekskul/template", {
          params,
          responseType: "blob",
        });
        const url = window.URL.createObjectURL(response.data);
        const a = document.createElement("a");
        a.href = url;
        a.download = `template-nilai-ekskul-${params.semester}.xlsx`;
        document.body.appendChild(a);
        a.click();
        a.remove();
        window.URL.revokeObjectURL(url);
      } finally {
        this.templateLoading = false;
      }
    },

    async importExcel(file, params) {
      this.importing = true;
      try {
        const fd = new FormData();
        fd.append("file", file);
        Object.entries(params).forEach(([k, v]) => fd.append(k, v));
        const response = await api.post("/nilai-ekskul/import", fd, {
          headers: { "Content-Type": "multipart/form-data" },
        });
        return response.data;
      } finally {
        this.importing = false;
      }
    },
  },
});
