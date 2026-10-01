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
