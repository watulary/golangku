import { defineStore } from "pinia";
import api from "../services/api";

function downloadBlob(response, fallbackFilename) {
  const disposition = response.headers["content-disposition"] || "";
  const match = disposition.match(/filename="?([^"]+)"?/);
  const filename = match ? match[1] : fallbackFilename;
  const url = window.URL.createObjectURL(new Blob([response.data]));
  const link = document.createElement("a");
  link.href = url;
  link.setAttribute("download", filename);
  document.body.appendChild(link);
  link.click();
  link.remove();
  window.URL.revokeObjectURL(url);
}

export const useBukuIndukStore = defineStore("bukuinduk", {
  state: () => ({
    cetakLoading: false,
  }),
  actions: {
    // sumber: "aktif" | "mutasi" | "arsip"
    // id: nis (aktif/mutasi) atau kode_siswa (arsip)
    async cetakSatu(sumber, id, tahunPelajaran) {
      this.cetakLoading = true;
      try {
        const response = await api.get(`/bukuinduk/cetak/${sumber}/${id}`, {
          params: { tahun_pelajaran: tahunPelajaran },
          responseType: "blob",
        });
        downloadBlob(response, `buku-induk-${id}.pdf`);
      } catch (error) {
        if (!error.response) {
          // network-level error tanpa response (ciri khas IDM mencegat
          // XHR) — file sebenarnya sudah berhasil diunduh lewat IDM,
          // jadi diamkan saja, jangan tampilkan sbg gagal
          console.warn(
            "Cetak buku induk: kemungkinan dicegat IDM, diabaikan.",
            error,
          );
          return;
        }
        console.error("Gagal cetak buku induk:", error);
        throw error;
      } finally {
        this.cetakLoading = false;
      }
    },

    async cetakBulk(sumber, ids, tahunPelajaran) {
      this.cetakLoading = true;
      try {
        const response = await api.post(
          `/bukuinduk/cetak/${sumber}/bulk`,
          { tahun_pelajaran: tahunPelajaran, ids: ids.map(String) },
          { responseType: "blob" },
        );
        downloadBlob(response, "buku-induk-massal.pdf");
      } catch (error) {
        if (!error.response) {
          console.warn(
            "Cetak buku induk massal: kemungkinan dicegat IDM, diabaikan.",
            error,
          );
          return;
        }
        console.error("Gagal cetak buku induk massal:", error);
        throw error;
      } finally {
        this.cetakLoading = false;
      }
    },
  },
});
