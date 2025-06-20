import { defineConfig } from "@lingui/cli";

export default defineConfig({
  sourceLocale: "zh-CN",
  locales: ["zh-CN", "en-US"],
  catalogs: [
    {
      path: "src/locales/{locale}/messages",
      include: ["src/components/**", "src/features/**"],
      exclude: ["**/node_modules/**"],
    },
  ],
});