import { ref, watch } from "vue";

const theme = ref("dark");

/*
|--------------------------------------------------------------------------
| Apply Theme
|--------------------------------------------------------------------------
*/

function applyTheme(value) {
  document.documentElement.classList.toggle("dark", value === "dark");
}

/*
|--------------------------------------------------------------------------
| Init Theme
|--------------------------------------------------------------------------
*/

function initTheme() {
  const saved = localStorage.getItem("theme");

  // Dashboard default Dark
  theme.value = saved || "dark";

  applyTheme(theme.value);
}

/*
|--------------------------------------------------------------------------
| Toggle Theme
|--------------------------------------------------------------------------
*/

function toggleTheme() {
  theme.value = theme.value === "dark" ? "light" : "dark";
}

/*
|--------------------------------------------------------------------------
| Watch Theme
|--------------------------------------------------------------------------
*/

watch(theme, (value) => {
  localStorage.setItem("theme", value);
  applyTheme(value);
});

/*
|--------------------------------------------------------------------------
| Export
|--------------------------------------------------------------------------
*/

export function useTheme() {
  return {
    theme,
    initTheme,
    toggleTheme,
  };
}