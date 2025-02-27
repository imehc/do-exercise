"use client";
import { JSX, useEffect } from "react";
import { scan } from "react-scan";

/** @link https://github.com/aidenybai/react-scan */
export function ReactScan(): JSX.Element {
  useEffect(() => {
    scan({
      enabled: true,
    });
  }, []);

  return <></>;
}