# Editor bundle and fonts

This folder is a dev-time tool. It builds `internal/server/web/vendor/editor.js` (CodeMirror 6 with Java and C++) and copies the fonts into `internal/server/web/fonts/`. The built files are committed, so the Go program never needs Node and the app works offline.

To rebuild after changing `entry.js` or a package version:

```
cd tools/vendor
npm ci
npm run build
```

`entry.js` defines `window.DSAEditor.create(parent, { doc, language, onChange, onRun })`, which returns `{ getValue(), setValue(text), focus(), destroy() }`. Editor colours come from CSS variables in `style.css` (`--syn-keyword`, `--syn-type`, `--syn-number`, `--syn-comment`, `--gutter`, `--active-line`, `--selection`).

`copy-fonts.js` copies weights 400 and 500 of IBM Plex Sans and JetBrains Mono and 400 and 600 of Source Serif 4, Latin subset only, and writes `fonts/LICENSES.md`.
