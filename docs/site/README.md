This directory contains files for building the GitHub Pages site with
[Hugo][h] and the [Hextra][x] theme.  The site is published on release
via GitHub Actions ([workflow](/.github/workflows/site.yaml)).

[h]: https://gohugo.io/
[x]: https://imfing.github.io/hextra/

## Previewing the site

 * Install Hugo (<https://gohugo.io/installation/>)

 * Run `make serve` to prepare the necessary files and launch
   [Hugo's web server](https://gohugo.io/commands/hugo_server/)

## Components

 * [hugo.yaml](./hugo.yaml): site configuration

   * [Hugo docs](https://gohugo.io/getting-started/configuration/)

   * [Hextra docs](https://imfing.github.io/hextra/docs/guide/configuration/)

 * [go.mod](./go.mod): defines Hextra as a Hugo module

   * [Hugo docs](https://gohugo.io/hugo-modules/)

   * [Hextra docs](https://imfing.github.io/hextra/docs/getting-started/#setup-hextra-as-hugo-module)

 * [content/][c]: Markdown content of the site

   Most of the content is pulled from other spots in the repository by
   the [Makefile][m], including [NEWS.md](/NEWS.md) and the
   command-line Markdown documentation at
   [docs/commands/](/docs/commands).

   * [Hugo docs](https://gohugo.io/content-management/)

 * [layouts/](./layouts): templates for rendering the site pages

   Most of the layouts are defined by the Hextra theme.  The files
   here override particular templates.  They are derived from the
   corresponding Hextra template.  The header of each template
   documents what was changed.

   * [Hugo docs](https://gohugo.io/templates/)

 * [scripts/](./scripts): custom scripts used by the [Makefile][m] to
   process the Markdown content before writing it to [content/][c]

   This is responsible for things like adding PR links in the NEWS
   page rendered on the site.

 * [static/](./static): this directory contains the logos and any
   other images used on the site

[c]: ./content
[m]: ./Makefile
