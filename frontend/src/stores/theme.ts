import { computed, ref } from "vue";
import { defineStore } from "pinia";
import { theme as antdPresetTheme } from "antdv-next";

const THEME_KEY = "lykn.user.theme";

type AppTheme = "light" | "dark";

const readStoredTheme = (): AppTheme =>
  window.localStorage.getItem(THEME_KEY) === "dark" ? "dark" : "light";

export const useThemeStore = defineStore("theme", () => {
  const theme = ref<AppTheme>(readStoredTheme());
  const rootThemeClass = computed(() => `lykn-theme-${theme.value}`);
  const antdTheme = computed(() => ({
    algorithm:
      theme.value === "dark" ? antdPresetTheme.darkAlgorithm : antdPresetTheme.defaultAlgorithm,
  }));

  const applyTheme = (value: AppTheme) => {
    document.documentElement.dataset.theme = value;
  };

  applyTheme(theme.value);

  const setTheme = (value: AppTheme) => {
    theme.value = value;
    window.localStorage.setItem(THEME_KEY, value);
    applyTheme(value);
  };

  return {
    theme,
    antdTheme,
    rootThemeClass,
    setTheme,
  };
});
