import { defineStore } from "pinia";
import api from "../services/api";

// Nilai akhir dikirim sebagai angka, atau null kalau kosong.
function snapshot(list) {
  return list.map((s) => ({
    nis: s.nis,
    nilai: s.nilai || "",
    elemen: (s.elemen || "").trim(),
    sub_elemen: (s.sub_elemen || "").trim(),
    target: (s.target || "").trim(),
  }));
}

export const useP5Store = defineStore("p5", {
  state: () => ({
    daftar: [],
    original: "[]",
    loading: false,
    submitting: false,
    importing: false,
    templateLoading: false,
  }),

  actions: {
    async fetchDaftar(params) {
      this.loading = true;
      try {
        const response = await api.get("/p5/daftar", { params });
        this.daftar = response.data.data || [];
        this.original = JSON.stringify(snapshot(this.daftar));
      } catch (error) {
        this.daftar = [];
        this.original = "[]";
        throw error;
      } finally {
        this.loading = false;
      }
    },

    clear() {
      this.daftar = [];
      this.original = "[]";
    },

    hasChanges() {
      return JSON.stringify(snapshot(this.daftar)) !== this.original;
    },

    async submitBulk(payload) {
      this.submitting = true;
      try {
        const items = this.daftar.map((s) => ({
          nis: s.nis,
          nilai: s.nilai || "",
          elemen: s.elemen || "",
          sub_elemen: s.sub_elemen || "",
          target: s.target || "",
        }));
        const response = await api.post("/p5/bulk", { ...payload, items });
        return response.data;
      } finally {
        this.submitting = false;
      }
    },

    async downloadTemplate(params) {
      this.templateLoading = true;
      try {
        const response = await api.get("/p5/template", {
          params,
          responseType: "blob",
        });
        const url = window.URL.createObjectURL(response.data);
        const a = document.createElement("a");
        a.href = url;
        a.download = `template-p5-${params.semester}.xlsx`;
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
        const response = await api.post("/p5/import", fd, {
          headers: { "Content-Type": "multipart/form-data" },
        });
        return response.data;
      } finally {
        this.importing = false;
      }
    },
  },
});
