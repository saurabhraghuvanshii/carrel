// Copies the Latin woff2 files and their licences into internal/server/web/fonts.
import { copyFileSync, mkdirSync, readFileSync, writeFileSync } from 'node:fs';

const out = '../../internal/server/web/fonts';
const fonts = [
  { pkg: 'ibm-plex-sans', name: 'IBM Plex Sans', weights: [400, 500] },
  { pkg: 'source-serif-4', name: 'Source Serif 4', weights: [400, 600] },
  { pkg: 'jetbrains-mono', name: 'JetBrains Mono', weights: [400, 500] },
];

mkdirSync(out, { recursive: true });
let licences = '# Font licences\n\nThe fonts in this folder are the Latin subsets from the Fontsource packages, used under the SIL Open Font License 1.1.\n';
for (const f of fonts) {
  for (const w of f.weights) {
    const file = `${f.pkg}-latin-${w}-normal.woff2`;
    copyFileSync(`node_modules/@fontsource/${f.pkg}/files/${file}`, `${out}/${file}`);
  }
  const licence = readFileSync(`node_modules/@fontsource/${f.pkg}/LICENSE`, 'utf8').trim();
  licences += `\n## ${f.name}\n\n\`\`\`\n${licence}\n\`\`\`\n`;
}
writeFileSync(`${out}/LICENSES.md`, licences);
