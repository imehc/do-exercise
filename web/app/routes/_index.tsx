import type { MetaFunction } from "@remix-run/node";
import { ModeToggle } from "~/components/mode-toggle";
import { useTranslation } from "react-i18next";

export const meta: MetaFunction = () => {
  return [
    { title: "do-exercise" },
    { name: "description", content: "do-exercise manage system" },
  ];
};

export const handle = { i18n: "home" };

export default function Index() {
  const { t } = useTranslation();
  const { t: t2 } = useTranslation("home");

  return (
    <div className="w-screen h-screen">
      {/* https://remix.run/resources/remix-i18next */}
      <h1>{t("greeting")}</h1>
      <h1>{t2("title")}</h1>
      <ModeToggle />
    </div>
  );
}

