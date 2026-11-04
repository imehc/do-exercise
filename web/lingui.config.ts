import { defineConfig } from "@lingui/cli";

export default defineConfig({
  sourceLocale: "zh-CN",
  locales: ["zh-CN", "en-US"],
  catalogs: [
    {
      path: "<rootDir>/src/locales/{locale}/messages",
      // src/utils 也要扫：菜单翻译键在 src/utils/menu-catalog.ts 里用 msg() 声明，
      // 运行期是动态 id（i18n._({ id })），不声明就不会进 catalog。
      include: ["src/components/**", "src/features/**", "src/hooks/**", "src/commons/**", "src/utils/**"],
      exclude: ["**/node_modules/**"],
    },
  ],
});
