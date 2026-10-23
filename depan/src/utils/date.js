/*
|--------------------------------------------------------------------------
| Date Utilities
|--------------------------------------------------------------------------
| Format tanggal dan waktu untuk seluruh aplikasi.
| Locale   : Indonesia
| Timezone : Asia/Jakarta
|--------------------------------------------------------------------------
*/

const DATE_LOCALE = "id-ID";
const TIME_ZONE = "Asia/Jakarta";

/**
 * True kalau d adalah Date yang valid (bukan "Invalid Date").
 * Dipakai di semua fungsi format supaya data tanggal rusak (mis. "0000-00-00",
 * format aneh dari import data lama) tidak bikin Intl.DateTimeFormat throw
 * dan meng-crash seluruh halaman — cukup fallback ke "-".
 */
function isValidDate(d) {
  return d instanceof Date && !isNaN(d.getTime());
}

/**
 * 1 Juli 2026
 */
export function formatDate(date) {
  if (!date) return "-";

  const d = new Date(date);
  if (!isValidDate(d)) return "-";

  return new Intl.DateTimeFormat(DATE_LOCALE, {
    day: "numeric",
    month: "long",
    year: "numeric",
    timeZone: TIME_ZONE,
  }).format(d);
}

/**
 * 21.30 WIB
 */
export function formatTime(date) {
  if (!date) return "-";

  const d = new Date(date);
  if (!isValidDate(d)) return "-";

  return (
    new Intl.DateTimeFormat(DATE_LOCALE, {
      hour: "2-digit",
      minute: "2-digit",
      hour12: false,
      timeZone: TIME_ZONE,
    }).format(d) + " WIB"
  );
}

/**
 * 1 Juli 2026 • 21.30 WIB
 */
export function formatDateTime(date) {
  if (!date) return "-";

  const d = new Date(date);
  if (!isValidDate(d)) return "-";

  return `${formatDate(date)} • ${formatTime(date)}`;
}

/**
 * 01/07/2026
 */
export function formatDateShort(date) {
  if (!date) return "-";

  const d = new Date(date);
  if (!isValidDate(d)) return "-";

  return new Intl.DateTimeFormat(DATE_LOCALE, {
    day: "2-digit",
    month: "2-digit",
    year: "numeric",
    timeZone: TIME_ZONE,
  }).format(d);
}

/**
 * 01-07-2026
 */
export function formatDateNumeric(date) {
  if (!date) return "-";

  const d = new Date(date);
  if (!isValidDate(d)) return "-";

  const parts = new Intl.DateTimeFormat(DATE_LOCALE, {
    day: "2-digit",
    month: "2-digit",
    year: "numeric",
    timeZone: TIME_ZONE,
  }).formatToParts(d);

  const get = (type) => parts.find((p) => p.type === type)?.value;

  return `${get("day")}-${get("month")}-${get("year")}`;
}

/**
 * 01/07/2026 • 21.30 WIB
 */
export function formatDateTimeShort(date) {
  if (!date) return "-";

  const d = new Date(date);
  if (!isValidDate(d)) return "-";

  return `${formatDateShort(date)} • ${formatTime(date)}`;
}

/**
 * 30-12-2026 06:10:48
 */
export function formatDateTimeNumeric(date) {
  if (!date) return "-";

  const d = new Date(date);
  if (!isValidDate(d)) return "-";

  const pad = (n) => String(n).padStart(2, "0");

  const parts = new Intl.DateTimeFormat(DATE_LOCALE, {
    day: "2-digit",
    month: "2-digit",
    year: "numeric",
    hour: "2-digit",
    minute: "2-digit",
    second: "2-digit",
    hour12: false,
    timeZone: TIME_ZONE,
  }).formatToParts(d);

  const get = (type) => parts.find((p) => p.type === type)?.value;

  return `${pad(get("day"))}-${pad(get("month"))}-${get("year")} ${pad(get("hour"))}:${pad(get("minute"))}:${pad(get("second"))}`;
}
