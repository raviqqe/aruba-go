import sitemap from "@astrojs/sitemap";
import starlight from "@astrojs/starlight";
import { defineConfig } from "astro/config";

export default defineConfig({
  base: "/aruba-go",
  integrations: [
    sitemap(),
    starlight({
      title: "Muffy",
      customCss: ["./src/index.css"],
      favicon: "/icon.svg",
      head: [
        {
          tag: "link",
          attrs: {
            rel: "manifest",
            href: "/aruba-go/manifest.json",
          },
        },
        {
          tag: "meta",
          attrs: {
            property: "og:image",
            content: "/aruba-go/icon.svg",
          },
        },
      ],
      logo: {
        src: "./src/icon.svg",
      },
      social: [
        {
          icon: "github",
          label: "GitHub",
          href: "https://github.com/raviqqe/aruba-go",
        },
      ],
      sidebar: [
        {
          label: "Home",
          link: "/",
        },
        {
          label: "Install",
          link: "/install",
        },
        {
          label: "Usage",
          link: "/usage",
        },
      ],
    }),
  ],
  prefetch: { prefetchAll: true },
  site: "https://raviqqe.github.io/aruba-go",
});
