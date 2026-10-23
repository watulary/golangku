export function formatHits(num) {
  if (num < 1000) return num.toString();

  const units = ["K", "M", "B", "T"];
  let unitIndex = -1;

  while (num >= 1000) {
    num /= 1000;
    unitIndex++;
  }

  return num.toFixed(1).replace(/\.0$/, "") + units[unitIndex];
}
