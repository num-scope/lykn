import { defineConfig } from "vite-plus";
import vue from "@vitejs/plugin-vue";
import unocss from "unocss/vite";
import dayjs from "vite-plugin-dayjs";

export default defineConfig({
  fmt: {
    ignorePatterns: ["src/components/admin-kit/**"],
  },
  lint: {
    options: {
      typeAware: true,
      typeCheck: true,
    },
  },
  plugins: [dayjs(), vue(), unocss()],
});
