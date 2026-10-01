# Carrel website

A one-page static site: plain HTML, CSS and one small script. No framework, no npm, no CDN. It is assembled from partials by a small Go tool in `tools/sitebuild`.

```
make site          # builds site/src into site/dist
make site-serve    # builds, serves on http://127.0.0.1:8080, rebuilds on change
go run ./tools/sitebuild -release   # same build, but fails while placeholders remain
sh scripts/test-install.sh          # tests install.sh against a fake local release
```

## Layout

```
src/index.html          page shell; each <!-- include: partials/x.html --> line is replaced by that file
src/404.html            not-found page
src/site.config.json    domain, licence and repository (see below)
src/css/*.css           joined in a fixed order into dist/site.css
src/js/site.js          theme, accent, copy buttons, pausing loops off screen
src/partials/*.html     one file per section
public/                 copied as is (fonts, icons, og.png, robots.txt, sitemap.xml)
../scripts/install.*    copied to dist/ with tokens replaced
```

## Placeholders

`{{domain}}`, `{{licence}}`, `{{repo}}`, `{{repo_slug}}` and `{{site_url}}` are replaced in every page, stylesheet, script and text file. Their values live in `src/site.config.json`. If a value is left empty the page shows a visible placeholder such as `[your-domain]`, and `-release` refuses to build.

## Deploying

The site is published to https://saurabhraghuvanshii.github.io/carrel by `.github/workflows/site.yml` on every push to `main` that touches the site, the install scripts or the build tool. It can also be run by hand from the Actions tab. Pages must be set to source "GitHub Actions" in the repository settings.

Make sure the latest release has `carrel_<os>_<arch>.tar.gz` (`.zip` for Windows) and `checksums.txt`, because the install scripts download those names.

## Open Graph image

`public/og.png` is 1200 by 630: Paper background, the icon, the eyebrow in JetBrains Mono and the headline in Source Serif 4. To change it, write a small HTML page with those elements and render it with `google-chrome --headless=new --window-size=1200,630 --screenshot=og.png page.html`.
