import { defineStore } from "pinia";
import api from "../services/api";

function snapshot(list) {
  return list.map((s) => ({
    nis: s.nis,
    status: s.status || "",
    no_ijazah: (s.no_ijazah || "").trim(),
    no_skhus: (s.no_skhus || "").trim(),
    tanggal: s.tanggal || "",
  }));
}

export const useKeputusanStore = defineStore("keputusanakhir", {
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
        const response = await api.get("/keputusan-akhir/daftar", { params });
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
          status: s.status || "",
          no_ijazah: s.no_ijazah || "",
          no_skhus: s.no_skhus || "",
          tanggal: s.tanggal || "",
        }));
        const response = await api.post("/keputusan-akhir/bulk", {
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
        const response = await api.get("/keputusan-akhir/template", {
          params,
          responseType: "blob",
        });
        const url = window.URL.createObjectURL(response.data);
        const a = document.createElement("a");
        a.href = url;
        a.download = `template-keputusan-akhir.xlsx`;
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
        const response = await api.post("/keputusan-akhir/import", fd, {
          headers: { "Content-Type": "multipart/form-data" },
        });
        return response.data;
      } finally {
        this.importing = false;
      }
    },
  },
});
