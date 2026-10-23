import { defineStore } from "pinia";
import api from "../services/api";

export const useWaliKelasTingkatStore = defineStore("walikelastingkat", {
  state: () => ({
    siswa: [],
    loading: false,
    saving: false,
    importing: false,
    pegawaiOptions: [],
    lastImport: null, // { tersimpan, dilewati[] }
  }),

  actions: {
    async fetchSiswa(idRombel, tingkat) {
      this.loading = true;
      try {
        const res = await api.get("/bukuinduk/walikelas-tingkat", {
          params: { id_rombel: idRombel, tingkat },
        });
        this.siswa = res.data.data || [];
      } finally {
        this.loading = false;
      }
    },

    async searchPegawai(q) {
      const res = await api.get("/bukuinduk/walikelas-tingkat/pegawai", {
        params: { q },
      });
      this.pegawaiOptions = res.data.data || [];
      return this.pegawaiOptions;
    },

    async saveRombel(payload) {
      this.saving = true;
      try {
        const res = await api.put(
          "/bukuinduk/walikelas-tingkat/rombel",
          payload,
        );
        return res.data;
      } finally {
        this.saving = false;
      }
    },

    async clearRombel(idRombel, tingkat) {
      const res = await api.delete("/bukuinduk/walikelas-tingkat/rombel", {
        params: { id_rombel: idRombel, tingkat },
      });
      return res.data;
    },

    async importExcel(file, tingkat) {
      this.importing = true;
      try {
        const fd = new FormData();
        fd.append("file", file);
        fd.append("tingkat", String(tingkat));
        const res = await api.post("/bukuinduk/walikelas-tingkat/import", fd);
        this.lastImport = res.data;
        return res.data;
      } finally {
        this.importing = false;
      }
    },

    async downloadTemplate() {
      const res = await api.get("/bukuinduk/walikelas-tingkat/template", {
        responseType: "blob",
      });
      const url = URL.createObjectURL(res.data);
      const a = document.createElement("a");
      a.href = url;
      a.download = "template-walikelas-tingkat.xlsx";
      a.click();
      URL.revokeObjectURL(url);
    },
  },
});
