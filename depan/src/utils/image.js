/**
 * Mengompres file gambar sampai ukurannya tidak lebih dari maxSizeKB.
 *
 * Kalau file yang dipilih sudah cukup kecil, dikembalikan apa adanya
 * (tanpa re-encode, supaya kualitas tidak turun percuma).
 *
 * Kalau perlu dikompres: hasil akhir SELALU keluar sbg JPEG (format
 * paling efisien utk foto). Strategi bertahap tiap percobaan:
 *   1. Turunkan quality dulu (0.85 -> ... -> 0.35)
 *   2. Kalau quality minimum masih belum cukup, kecilkan dimensi
 *      gambar 20% lalu ulangi dari quality 0.7
 * Maksimal 20 percobaan supaya tidak infinite loop.
 *
 * @param {File} file
 * @param {number} maxSizeKB - default 100 (KB)
 * @returns {Promise<File>}
 */
export function compressImage(file, maxSizeKB = 100) {
  const maxBytes = maxSizeKB * 1024;

  // Sudah cukup kecil, tidak perlu diapa-apakan
  if (file.size <= maxBytes) {
    return Promise.resolve(file);
  }

  return new Promise((resolve, reject) => {
    const img = new Image();
    const objectUrl = URL.createObjectURL(file);

    img.onload = () => {
      let width = img.naturalWidth;
      let height = img.naturalHeight;

      const canvas = document.createElement("canvas");
      const ctx = canvas.getContext("2d");

      let quality = 0.85;
      let attempts = 0;
      const maxAttempts = 20;

      const step = () => {
        attempts++;

        canvas.width = width;
        canvas.height = height;

        ctx.clearRect(0, 0, width, height);
        ctx.drawImage(img, 0, 0, width, height);

        canvas.toBlob(
          (blob) => {
            if (!blob) {
              URL.revokeObjectURL(objectUrl);
              reject(new Error("Gagal memproses gambar."));
              return;
            }

            const done = blob.size <= maxBytes || attempts >= maxAttempts;

            if (done) {
              URL.revokeObjectURL(objectUrl);

              const namaFile = file.name.replace(/\.[^/.]+$/, "") + ".jpg";

              resolve(
                new File([blob], namaFile, {
                  type: "image/jpeg",
                  lastModified: Date.now(),
                }),
              );
              return;
            }

            if (quality > 0.35) {
              quality = Math.max(0.35, quality - 0.1);
            } else {
              // Quality sudah minimum, kecilkan dimensi & reset quality
              width = Math.round(width * 0.8);
              height = Math.round(height * 0.8);
              quality = 0.7;
            }

            step();
          },
          "image/jpeg",
          quality,
        );
      };

      step();
    };

    img.onerror = () => {
      URL.revokeObjectURL(objectUrl);
      reject(new Error("File bukan gambar yang valid."));
    };

    img.src = objectUrl;
  });
}
