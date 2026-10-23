export function getErrorMessage(error) {
  if (error.response?.data?.message) {
    return error.response.data.message;
  }

  if (error.response?.status === 401) {
    return "Unauthorized.";
  }

  if (error.response?.status === 403) {
    return "Akses ditolak.";
  }

  if (error.response?.status === 404) {
    return "Data tidak ditemukan.";
  }

  if (error.response?.status >= 500) {
    return "Terjadi kesalahan pada server.";
  }

  if (error.code === "ECONNABORTED") {
    return "Request timeout.";
  }

  if (error.message) {
    return error.message;
  }

  return "Terjadi kesalahan.";
}