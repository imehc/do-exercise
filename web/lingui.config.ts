import { defineConfig } from "@lingui/cli";

export default defineConfig({
  sourceLocale: "zh-CN",
  locales: ["zh-CN", "en-US"],
  catalogs: [
    {
      path: "<rootDir>/src/locales/{locale}/messages",
      include: ["src/components/**", "src/features/**", "src/hooks/**", "src/commons/**"],
      exclude: ["**/node_modules/**"],
    },
  ],
});
