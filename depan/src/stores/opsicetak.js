import { defineStore } from "pinia";
import api from "../services/api";

function kosong() {
  return {
    sampul: true,
    daftar_isi: true,
    pemeriksaan: true,
    nomor_urut: true,
    petunjuk: true,
    lembar_iv: true,
    halaman_akhir: true,
    tingkat: [7, 8, 9],
    kepsek_nama: "",
    kepsek_nip: "",
    walikelas_nama: "",
    walikelas_nip: "",
  };
}

export const useOpsiCetakStore = defineStore("bukuindukOpsiCetak", {
  state: () => ({
    opsi: kosong(),
    loading: false,
    saving: false,
  }),

  actions: {
    async fetchOpsi() {
      this.loading = true;
      try {
        const response = await api.get("/bukuinduk/opsi-cetak");
        this.opsi = { ...kosong(), ...(response.data.data || {}) };
      } finally {
        this.loading = false;
      }
    },

    async saveOpsi() {
      this.saving = true;
      try {
        const response = await api.put("/bukuinduk/opsi-cetak", this.opsi);
        return response.data;
      } finally {
        this.saving = false;
      }
    },
  },
});
