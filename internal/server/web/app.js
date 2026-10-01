(() => {
  'use strict';

  const view = document.getElementById('view');
  const crumb = document.getElementById('crumb');
  const langSel = document.getElementById('lang');

  // ---- helpers ----

  function el(tag, props = {}, ...kids) {
    const n = document.createElement(tag);
    for (const [k, v] of Object.entries(props)) {
      if (v === false || v == null) continue;
      if (k === 'class') n.className = v;
      else if (k === 'text') n.textContent = v;
      else if (k.startsWith('on')) n.addEventListener(k.slice(2), v);
      else n.setAttribute(k, v === true ? '' : v);
    }
    for (const kid of kids.flat()) if (kid != null && kid !== false) n.append(kid);
    return n;
  }

  async function api(method, path, body) {
    const res = await fetch('/api' + path, {
      method,
      headers: { 'Content-Type': 'application/json', 'X-DSA': '1' },
      body: body ? JSON.stringify(body) : undefined,
    });
    const data = await res.json().catch(() => ({}));
    if (!res.ok) throw new Error(data.error || res.statusText);
    return data;
  }

  const cap = (s) => s.charAt(0).toUpperCase() + s.slice(1);
  const pad2 = (n) => String(n).padStart(2, '0');

  // ---- state ----

  let cfg = null;
  let problems = [];
  const ui = { sheet: 'patterns', filter: 'all' };
  let pendingSave = null; // { timer, fn }
  let activeEditor = null;
  let renderToken = 0; // a slow page load must not draw over a newer page

  async function flushSave() {
    if (!pendingSave) return;
    clearTimeout(pendingSave.timer);
    const fn = pendingSave.fn;
    pendingSave = null;
    await fn();
  }

  function applyTheme() {
    let theme = cfg.theme;
    if (theme === 'system') theme = matchMedia('(prefers-color-scheme: dark)').matches ? 'ink' : 'paper';
    document.documentElement.dataset.theme = theme;
    document.documentElement.dataset.accent = cfg.accent;
  }
  matchMedia('(prefers-color-scheme: dark)').addEventListener('change', () => cfg && applyTheme());

  async function saveConfig(patch = {}, aiPatch = {}) {
    const next = { ...cfg, ...patch };
    const body = {
      theme: next.theme,
      accent: next.accent,
      lang: next.lang,
      remindReviews: next.remindReviews,
      ai: {
        provider: cfg.ai.provider,
        model: cfg.ai.model,
        explainOnly: cfg.ai.explainOnly,
        ...aiPatch,
      },
    };
    cfg = await api('PUT', '/config', body);
    applyTheme();
    langSel.value = cfg.lang;
    return cfg;
  }

  // ---- markdown: headings, lists, fenced code, `code` and **bold**, built with textContent only ----

  function inline(text) {
    return text
      .split(/(`[^`]+`|\*\*[^*]+\*\*)/g)
      .filter(Boolean)
      .map((s) => {
        if (s.length > 1 && s.startsWith('`') && s.endsWith('`')) return el('code', { text: s.slice(1, -1) });
        if (s.length > 4 && s.startsWith('**') && s.endsWith('**')) return el('strong', { text: s.slice(2, -2) });
        return document.createTextNode(s);
      });
  }

  function renderMarkdown(src) {
    const nodes = [];
    const lines = src.replace(/\r/g, '').split('\n');
    const fence = (l) => /^\s*(```|~~~)/.test(l);
    const heading = (l) => /^#{1,3} /.test(l);
    const bullet = (l) => /^\s*[-*] /.test(l);
    const numbered = (l) => /^\s*\d+[.)] /.test(l);
    const special = (l) => heading(l) || bullet(l) || numbered(l) || fence(l);
    let i = 0;
    while (i < lines.length) {
      const line = lines[i];
      if (!line.trim()) { i++; continue; }
      if (fence(line)) {
        const mark = line.trim().slice(0, 3);
        const code = [];
        i++;
        while (i < lines.length && !lines[i].trim().startsWith(mark)) { code.push(lines[i]); i++; }
        i++;
        nodes.push(el('pre', { class: 'code-block' }, el('code', { text: code.join('\n') })));
        continue;
      }
      if (heading(line)) { nodes.push(el('h2', {}, ...inline(line.replace(/^#+ /, '')))); i++; continue; }
      if (bullet(line) || numbered(line)) {
        const isNumbered = numbered(line);
        const test = isNumbered ? numbered : bullet;
        const items = [];
        while (i < lines.length && test(lines[i])) {
          items.push(el('li', {}, ...inline(lines[i].replace(/^\s*([-*]|\d+[.)]) /, ''))));
          i++;
        }
        nodes.push(el(isNumbered ? 'ol' : 'ul', {}, ...items));
        continue;
      }
      const para = [];
      while (i < lines.length && lines[i].trim() && !special(lines[i])) { para.push(lines[i]); i++; }
      nodes.push(el('p', {}, ...inline(para.join(' '))));
    }
    return nodes;
  }

  // ---- router ----

  function route() {
    renderToken++;
    if (activeEditor) { activeEditor.destroy(); activeEditor = null; }
    const h = location.hash || '#/';
    if (h.startsWith('#/p/')) return renderPractice(decodeURIComponent(h.slice(4)));
    if (h === '#/settings') return renderSettings();
    return renderSheets();
  }
  window.addEventListener('hashchange', async () => { await flushSave(); route(); });

  // ---- sheets ----

  const statusLabel = (p) => p.review ? 'Review'
    : p.status === 'solved' ? 'Solved'
    : p.status === 'tried' ? 'Tried'
    : p.ready ? 'Ready next' : 'Not started';
  const statusClass = (p) => p.review ? 'review' : p.status || (p.ready ? 'ready' : '');
  const titleOf = (pid) => (problems.find((x) => x.id === pid) || { title: pid }).title;

  async function renderSheets() {
    const token = renderToken;
    crumb.textContent = 'Sheets';
    let next = null;
    let progress = null;
    try {
      [problems, progress, next] = await Promise.all([
        api('GET', '/problems'),
        api('GET', '/progress'),
        api('GET', '/next?sheet=' + ui.sheet).catch(() => null),
      ]);
    } catch (e) {
      if (token === renderToken) view.replaceChildren(el('div', { class: 'page' }, el('p', { text: 'Could not load the sheets: ' + e.message })));
      return;
    }
    if (token !== renderToken) return;

    const sheets = [['patterns', 'Patterns'], ['real', 'Real interviews']];
    const all = problems.filter((p) => p.sheet === ui.sheet);
    const numbers = new Map(all.map((p, i) => [p.id, i + 1]));
    const shown = all.filter((p) => ui.filter === 'all' || p.difficulty === ui.filter);

    const tabs = el('div', { class: 'sheet-tabs' },
      el('div', { class: 'tabset' }, sheets.map(([id, name]) =>
        el('button', { class: id === ui.sheet ? 'on' : '', text: name, onclick: () => { ui.sheet = id; renderSheets(); } }))),
      archiveButtons());

    const chips = el('div', { class: 'chips' }, ['all', 'easy', 'medium', 'hard'].map((f) =>
      el('button', { class: `chip ${f}${ui.filter === f ? ' on' : ''}`, text: cap(f), onclick: () => { ui.filter = f; renderSheets(); } })));

    const list = el('div');
    if (!shown.length) list.append(el('p', { class: 'hint', text: 'No problems here yet.' }));
    let lastGroup = null;
    for (const p of shown) {
      if (p.group !== lastGroup) {
        list.append(el('div', { class: 'group-title', text: `${Math.floor(p.order / 100)} · ${p.group}` }));
        lastGroup = p.group;
      }
      const needs = p.needs && p.needs.length ? 'Needs: ' + p.needs.map(titleOf).join(', ') : null;
      list.append(el('a', { class: 'list-row', href: '#/p/' + encodeURIComponent(p.id), title: needs },
        el('span', { class: 'num', text: pad2(numbers.get(p.id)) }),
        el('span', { text: p.title }),
        el('span', { class: 'dif ' + p.difficulty, text: cap(p.difficulty) }),
        el('span', { class: 'tags', text: (p.tags || []).join(', ') }),
        el('span', { class: 'status' }, el('span', { class: 'dot ' + statusClass(p) }), statusLabel(p))));
    }

    const sp = progress[ui.sheet];
    const share = sp.total ? Math.round((100 * sp.solved) / sp.total) : 0;
    const byDifficulty = ['easy', 'medium', 'hard']
      .filter((d) => sp.byDifficulty[d].total)
      .map((d) => `${cap(d)} ${sp.byDifficulty[d].solved} of ${sp.byDifficulty[d].total}`).join(' · ');
    const sheetName = (id) => (sheets.find(([s]) => s === id) || [id, id])[1];

    const continueBox = next
      ? el('div', { class: 'aside-block' },
        el('div', { class: 'aside-label', text: next.reason === 'continue' ? 'Continue' : 'Up next' }),
        el('div', { class: 'aside-title', text: next.title }),
        el('div', { class: 'hint', text: `${sheetName(next.sheet)} · ${next.group}` }),
        el('a', { class: 'btn primary open', href: '#/p/' + encodeURIComponent(next.id), text: 'Open' }))
      : el('div', { class: 'aside-block' },
        el('div', { class: 'aside-label', text: 'Continue' }),
        el('div', { class: 'aside-title', text: sp.total ? 'Every problem in this sheet is solved.' : 'No problems here yet.' }));

    const aside = el('aside', { class: 'aside' },
      continueBox,
      el('div', { class: 'aside-block ruled' },
        el('div', { class: 'aside-label', text: 'Progress' }),
        el('div', { text: `${sp.solved} of ${sp.total} solved` }),
        el('div', { class: 'bar-track', role: 'progressbar', 'aria-valuemin': '0', 'aria-valuemax': String(sp.total), 'aria-valuenow': String(sp.solved), 'aria-label': 'Solved in this sheet' },
          el('div', { class: 'bar-fill', style: `width:${share}%` })),
        byDifficulty ? el('div', { class: 'hint', text: byDifficulty }) : null),
      el('p', { class: 'aside-note ruled', text: 'Problems are in learning order. Each one uses an idea from the one before it.' }));

    view.replaceChildren(el('div', { class: 'sheets' }, el('div', { class: 'sheets-main' }, importBox(), tabs, chips, list), aside));
  }

  // ---- export and import ----

  let importNotice = null; // { file, report, error }, shown until closed
  const count = (n, one, many) => `${n} ${n === 1 ? one : many}`;

  function archiveButtons() {
    const picker = el('input', { type: 'file', accept: '.zip,application/zip', class: 'visually-hidden', tabindex: '-1', 'aria-hidden': 'true',
      onchange: () => { if (picker.files[0]) importFile(picker.files[0], false); } });
    return el('div', { class: 'archive' },
      picker,
      el('button', { class: 'btn', text: 'Import solutions', onclick: () => picker.click() }),
      el('a', { class: 'btn', href: '/api/export', download: '', text: 'Export all' }));
  }

  async function importFile(file, overwrite) {
    await flushSave();
    try {
      const res = await fetch('/api/import' + (overwrite ? '?overwrite=true' : ''), {
        method: 'POST', headers: { 'Content-Type': 'application/zip', 'X-DSA': '1' }, body: file,
      });
      const data = await res.json().catch(() => ({}));
      if (!res.ok) throw new Error(data.error || res.statusText);
      importNotice = { file, report: data };
    } catch (e) {
      importNotice = { file, error: e.message };
    }
    route();
  }

  function importBox() {
    if (!importNotice) return null;
    const { file, report, error } = importNotice;
    const close = () => { importNotice = null; route(); };
    if (error) {
      const box = el('div', { class: 'import-box', role: 'status' },
        el('div', { class: 'summary fail', text: 'Import failed: ' + error }),
        el('div', { class: 'archive' }, el('button', { class: 'btn', text: 'Close', onclick: close })));
      requestAnimationFrame(() => box.scrollIntoView({ block: 'nearest' }));
      return box;
    }
    const names = (list) => el('ul', {}, list.slice(0, 10).map((x) => el('li', { text: x })),
      list.length > 10 ? el('li', { text: `and ${list.length - 10} more` }) : null);
    const lines = [el('div', { class: report.imported.length ? 'summary pass' : 'summary', text: `Imported ${count(report.imported.length, 'solution', 'solutions')}.` })];
    if (report.progressMerged) lines.push(el('div', { text: `Progress: ${count(report.progressMerged, 'problem', 'problems')} added or updated.` }));
    if (report.rejected.length) {
      lines.push(el('div', { text: `${count(report.rejected.length, 'file was', 'files were')} left out:` }),
        names(report.rejected.map((r) => `${r.name}: ${r.reason}`)));
    }
    let actions = el('div', { class: 'archive' }, el('button', { class: 'btn', text: 'Close', onclick: close }));
    if (report.skipped.length) {
      lines.push(el('div', { text: `${count(report.skipped.length, 'solution already exists', 'solutions already exist')} on this computer:` }),
        names(report.skipped),
        el('div', { text: 'Replace them with the ones from the zip?' }));
      actions = el('div', { class: 'archive' },
        el('button', { class: 'btn primary', text: 'Replace them', onclick: () => importFile(file, true) }),
        el('button', { class: 'btn', text: 'Keep mine', onclick: close }));
    }
    const box = el('div', { class: 'import-box', role: 'status' }, lines, actions);
    requestAnimationFrame(() => box.scrollIntoView({ block: 'nearest' }));
    return box;
  }

  // ---- practice ----

  async function renderPractice(id) {
    const token = renderToken;
    view.replaceChildren(el('div', { class: 'page' }, el('p', { class: 'hint', text: 'Loading...' })));
    let p;
    try {
      p = await api('GET', `/problems/${encodeURIComponent(id)}?lang=${cfg.lang}`);
      if (token !== renderToken) return;
    } catch (e) {
      if (token !== renderToken) return;
      view.replaceChildren(el('div', { class: 'page' }, el('p', { text: e.message }), el('a', { href: '#/', text: 'Back to sheets' })));
      return;
    }
    crumb.textContent = `${p.sheet === 'real' ? 'Real interviews' : 'Patterns'} / ${p.group}`;
    const linkList = (label, ids) => ids && ids.length
      ? [label + ' ', ...ids.flatMap((x, i) => [i ? ', ' : '', el('a', { href: '#/p/' + encodeURIComponent(x), text: titleOf(x) })]), '. ']
      : [];

    // Examples go right after the description; Constraints and anything after it follow them.
    const lines = p.statement.replace(/\r/g, '').split('\n');
    const cut = lines.indexOf('## Constraints');
    const intro = cut < 0 ? p.statement : lines.slice(0, cut).join('\n');
    const rest = cut < 0 ? '' : lines.slice(cut).join('\n');

    const statement = el('div', { class: 'statement' },
      el('div', { class: 'meta-row' }, el('span', { class: 'tag ' + p.difficulty, text: cap(p.difficulty) }), el('span', { text: (p.tags || []).map(cap).join(', ') })),
      el('h1', { text: p.title }),
      renderMarkdown(intro),
      (p.examples || []).map((ex) => el('div', { class: 'example' },
        el('div', { class: 'label', text: ex.label }),
        el('pre', { text: ex.display || ex.input }))),
      renderMarkdown(rest),
      el('div', { class: 'links' }, linkList('Builds on:', p.buildsOn), linkList('Leads to:', p.leadsTo)));

    const fileName = p.lang === 'java' ? 'Solution.java' : 'solution.cpp';
    const savedLabel = el('span', {}, 'Saved on this computer');
    const savedBox = el('span', { class: 'saved' }, el('i'), savedLabel);
    const editorBox = el('div', { class: 'editor' });

    const out = el('div', { class: 'results' }, el('div', { class: 'hint', text: 'Run the examples to see how your code does. Submit also checks fresh random cases.' }));
    const runBtn = el('button', { class: 'btn strong', text: 'Run examples', onclick: () => run('examples') });
    const submitBtn = el('button', { class: 'btn primary', text: 'Submit', onclick: () => run('submit') });

    let lastReport = '';

    async function save() {
      try {
        await api('PUT', '/solution', { problem: p.id, lang: p.lang, code: editor.getValue() });
        savedLabel.textContent = 'Saved on this computer';
      } catch (e) {
        savedLabel.textContent = 'Could not save: ' + e.message;
      }
    }
    function scheduleSave() {
      savedLabel.textContent = 'Saving...';
      if (pendingSave) clearTimeout(pendingSave.timer);
      pendingSave = { fn: async () => { await save(); }, timer: setTimeout(async () => { pendingSave = null; await save(); }, 600) };
    }
    const editor = window.DSAEditor.create(editorBox, {
      doc: p.code,
      language: p.lang,
      onChange: scheduleSave,
      onRun: (mode) => run(mode),
    });
    activeEditor = editor;

    async function run(mode) {
      if (pendingSave) { clearTimeout(pendingSave.timer); pendingSave = null; }
      runBtn.disabled = submitBtn.disabled = true;
      out.replaceChildren(el('div', { class: 'hint', text: mode === 'submit' ? 'Submitting...' : 'Running...' }));
      try {
        const rep = await api('POST', '/run', { problem: p.id, lang: p.lang, code: editor.getValue(), mode });
        savedLabel.textContent = 'Saved on this computer';
        lastReport = summarise(rep);
        renderReport(out, rep, mode);
        if (mode === 'submit' && rep.status === 'ok' && rep.passed === rep.total) showNext(out, p);
      } catch (e) {
        out.replaceChildren(el('div', { class: 'message', text: e.message }));
      } finally {
        runBtn.disabled = submitBtn.disabled = false;
      }
    }

    const question = el('input', { type: 'text', placeholder: 'Ask for a hint or an explanation', 'aria-label': 'Ask AI' });
    const answer = el('div', { class: 'answer' });
    const askBtn = el('button', { class: 'btn', text: 'Ask AI', onclick: async () => {
      askBtn.disabled = true;
      answer.replaceChildren(el('div', { class: 'hint', text: 'Thinking...' }));
      try {
        const r = await api('POST', '/ai', { problem: p.id, lang: p.lang, code: editor.getValue(), question: question.value, lastReport });
        answer.replaceChildren(...renderMarkdown(r.answer));
      } catch (e) {
        answer.replaceChildren(el('div', { class: 'message', text: 'AI help failed: ' + e.message }));
      } finally {
        askBtn.disabled = false;
      }
    } });

    const work = el('div', { class: 'work' },
      el('div', { class: 'tabs' }, el('span', { class: 'file', text: fileName }), savedBox),
      editorBox,
      el('div', { class: 'console' },
        el('div', { class: 'actions' },
          el('div', { class: 'group' }, runBtn, submitBtn),
          el('span', { class: 'hint', text: 'Ctrl + Enter runs · Ctrl + Shift + Enter submits' })),
        out,
        el('div', { class: 'ai' },
          el('div', { class: 'hint' }, 'Ask AI explains and gives hints. ',
            cfg.ai.hasKey || cfg.ai.provider === 'ollama'
              ? (cfg.ai.explainOnly ? 'It will not write the full solution.' : 'Explain only is off.')
              : el('a', { href: '#/settings', text: 'Add your key in Settings.' })),
          el('div', { class: 'ai-row' }, question, askBtn),
          answer)));

    view.replaceChildren(el('div', { class: 'practice' }, statement, work));
    editor.focus();
  }

  async function showNext(out, p) {
    const line = el('div', { class: 'next-line' }, 'Solved.');
    const summary = out.querySelector('.summary');
    if (summary) summary.after(line); else out.prepend(line);
    try {
      const n = await api('GET', `/next?sheet=${p.sheet}&after=${encodeURIComponent(p.id)}`);
      line.append(' Next: ', el('a', { href: '#/p/' + encodeURIComponent(n.id), text: n.title }));
    } catch (_) {
      line.append(' Every problem in this sheet is solved.');
    }
  }

  function summarise(rep) {
    const bad = rep.results.filter((r) => !r.passed).slice(0, 3)
      .map((r) => `${r.label}: input ${JSON.stringify(r.input || '')}, expected ${r.expected}, got ${r.got}`);
    return `Status ${rep.status}. Passed ${rep.passed} of ${rep.total}. ${rep.message || ''} ${bad.join(' | ')}`.trim();
  }

  function renderReport(out, rep, mode) {
    const kids = [];
    if (rep.status === 'compile_error') {
      kids.push(el('div', { class: 'summary fail', text: 'It did not compile' }), el('pre', { class: 'message', text: rep.message }));
      return out.replaceChildren(...kids);
    }
    if (rep.status === 'tooling_missing' || rep.status === 'internal_error') {
      kids.push(el('div', { class: 'summary fail', text: 'Could not run your code' }), el('div', { class: 'message', text: rep.message }));
      return out.replaceChildren(...kids);
    }
    const allPassed = rep.status === 'ok' && rep.passed === rep.total;
    const lead = { memory_limit: 'Memory limit reached. ', timeout: 'Time limit reached. ', runtime_error: 'Your program crashed. ' }[rep.status] || '';
    let headline = `${lead}Passed ${rep.passed} of ${rep.total}`;
    if (allPassed && mode === 'submit') headline = `All ${rep.total} tests passed. Marked solved.`;
    else if (allPassed) headline = `Examples ${rep.passed} of ${rep.total} passed`;
    kids.push(el('div', { class: 'summary ' + (allPassed ? 'pass' : 'fail'), text: headline }));
    if (rep.message) kids.push(el('div', { class: 'message', text: rep.message }));

    const randomPassed = rep.results.filter((r) => r.kind === 'random' && r.passed).length;
    for (const r of rep.results) {
      if (r.kind === 'random' && r.passed) continue;
      let detail = r.passed
        ? (r.got ? `got ${r.got}` : '')
        : `input\n${r.input}\nexpected ${r.expected}\ngot      ${r.got || '(nothing)'}`;
      if (r.note) detail += `\nwhy      ${r.note}`;
      const verdict = r.passed ? 'Passed' : r.got === '(crashed)' ? 'Crashed' : r.got === '(not run)' ? 'Not run' : 'Failed';
      kids.push(el('div', { class: 'row' },
        el('span', { class: 'hint', text: r.label }),
        el('span', { class: 'detail', text: detail }),
        el('span', { class: r.passed ? 'ok' : 'bad', text: verdict })));
    }
    if (randomPassed) kids.push(el('div', { class: 'hint', text: `${randomPassed} random cases passed.${rep.seed ? ' Seed ' + rep.seed + '.' : ''}` }));
    kids.push(el('div', { class: 'hint', text: `Took ${rep.durationMs} ms` }));
    out.replaceChildren(...kids);
  }

  // ---- settings ----

  async function renderSettings() {
    crumb.textContent = 'Settings';
    const doctor = await api('GET', '/doctor').catch(() => ({ tools: [], solutionsDir: '' }));
    const note = el('div', { class: 'note' });

    const themes = [['paper', 'Paper', 'Warm off-white'], ['ink', 'Ink', 'Warm charcoal'], ['system', 'Match my system', 'Paper by day, Ink by night']];
    const accents = [['brick', '#A8432B'], ['forest', '#2F5D50'], ['amber', '#B8792A'], ['slate', '#4A5A6A'], ['plum', '#6B4A63']];

    const appearance = el('section', { class: 'sec' },
      el('h2', { text: 'Appearance' }),
      el('div', {}, el('div', { class: 'field-label', text: 'Theme' }),
        el('div', { class: 'themes' }, themes.map(([id, name, desc]) =>
          el('button', { class: 'theme-card' + (cfg.theme === id ? ' on' : ''), 'aria-pressed': String(cfg.theme === id),
            onclick: async () => { await saveConfig({ theme: id }); renderSettings(); } },
            el('strong', { text: name }), el('span', { text: desc }))))),
      el('div', {}, el('div', { class: 'field-label', text: 'Accent colour' }),
        el('div', { class: 'swatches' }, accents.map(([id, color]) =>
          el('button', { class: 'swatch' + (cfg.accent === id ? ' on' : ''), style: 'background:' + color, 'aria-label': cap(id) + (cfg.accent === id ? ', selected' : ''),
            onclick: async () => { await saveConfig({ accent: id }); renderSettings(); } }))),
        el('div', { class: 'note', text: 'Easy, Medium and Hard tags keep their own colours in every theme.' })));

    const provider = el('select', { class: 'field', id: 'provider' },
      [['anthropic', 'Anthropic'], ['openai', 'OpenAI'], ['ollama', 'Ollama, on this computer']].map(([v, n]) => el('option', { value: v, text: n, selected: cfg.ai.provider === v })));
    const model = el('input', { type: 'text', id: 'model', placeholder: 'Model name from your provider' });
    model.value = cfg.ai.model || '';
    const key = el('input', { type: 'password', id: 'apikey', autocomplete: 'off',
      placeholder: cfg.ai.hasKey ? 'A key is saved. Paste a new one to replace it.' : 'Paste your key here' });
    let explainOnly = cfg.ai.explainOnly;
    const sw = el('button', { class: 'switch', role: 'switch', 'aria-checked': String(explainOnly), 'aria-label': 'Explain only',
      onclick: () => { explainOnly = !explainOnly; sw.setAttribute('aria-checked', String(explainOnly)); } });

    const saveAI = async (extra = {}) => {
      try {
        await saveConfig({}, { provider: provider.value, model: model.value, explainOnly, apiKey: key.value, ...extra });
        renderSettings();
      } catch (e) {
        note.textContent = e.message;
      }
    };
    const ai = el('section', { class: 'sec' },
      el('h2', { text: 'AI help' }),
      el('div', { class: 'grid2' },
        el('div', {}, el('label', { class: 'field-label', for: 'provider', text: 'Provider' }), provider),
        el('div', {}, el('label', { class: 'field-label', for: 'model', text: 'Model' }), model)),
      el('div', {}, el('label', { class: 'field-label', for: 'apikey', text: 'API key' }), key,
        el('div', { class: 'note', text: 'Saved in a config file in your home folder. It is never part of the project folder, and it is sent only to the provider you pick.' })),
      el('div', { class: 'switch-row' },
        el('div', {}, el('div', { text: 'Explain only' }), el('div', { class: 'note', text: 'Hints and explanations are fine. The AI will not write the full solution for you.' })),
        sw),
      el('div', { class: 'actions' },
        el('div', { class: 'group' },
          el('button', { class: 'btn primary', text: 'Save AI settings', onclick: () => saveAI() }),
          el('button', { class: 'btn', text: 'Test connection', onclick: async (e) => {
            const btn = e.currentTarget;
            btn.disabled = true;
            note.textContent = 'Saving and testing...';
            try {
              await saveConfig({}, { provider: provider.value, model: model.value, explainOnly, apiKey: key.value });
              key.value = '';
              await api('POST', '/ai/test', {});
              note.textContent = 'Connected. The provider answered.';
            } catch (err) {
              note.textContent = 'Could not connect: ' + err.message;
            } finally {
              btn.disabled = false;
            }
          } }),
          cfg.ai.hasKey ? el('button', { class: 'btn', text: 'Remove saved key', onclick: () => saveAI({ clearKey: true, apiKey: '' }) }) : null)),
      note);

    const langs = el('section', { class: 'sec' },
      el('h2', { text: 'Languages' }),
      el('div', { class: 'tools' }, doctor.tools.map((t) => [
        el('span', { text: t.name }),
        el('span', { class: 'path', text: t.found ? `${t.version || ''} ${t.path}`.trim() : 'Not found. Install it and make sure it is on your PATH.' }),
        el('span', { class: t.found ? 'ok' : 'bad', style: 'font-weight:500;color:var(--' + (t.found ? 'pass' : 'fail') + ')', text: t.found ? 'Found' : 'Missing' }),
      ]).flat()));

    const remind = el('button', { class: 'switch', role: 'switch', 'aria-checked': String(!!cfg.remindReviews), 'aria-label': 'Remind me to revisit solved problems',
      onclick: async () => {
        try {
          await saveConfig({ remindReviews: !cfg.remindReviews });
          remind.setAttribute('aria-checked', String(cfg.remindReviews));
        } catch (e) {
          practiceNote.textContent = e.message;
        }
      } });
    const practiceNote = el('div', { class: 'note' });
    const practice = el('section', { class: 'sec' },
      el('h2', { text: 'Practice' }),
      el('div', { class: 'switch-row' },
        el('div', {}, el('div', { text: 'Remind me to revisit solved problems' }),
          el('div', { class: 'note', text: 'A solved problem shows Review in the list 3 days later, then 10 days, then 30 days after you solve it again.' })),
        remind),
      practiceNote);

    const files = el('section', { class: 'sec' },
      el('h2', { text: 'Your solutions' }),
      el('div', {}, el('div', { class: 'field-label', text: 'Saved in' }), el('div', { class: 'mono', text: doctor.solutionsDir }),
        el('div', { class: 'note', text: 'Plain files. Open them in any editor, or put the folder in git.' })),
      importBox(),
      archiveButtons());

    view.replaceChildren(el('div', { class: 'page' }, appearance, practice, ai, langs, files));
  }

  // ---- start ----

  langSel.addEventListener('change', async () => {
    await flushSave();
    await saveConfig({ lang: langSel.value });
    route();
  });

  (async () => {
    try {
      [cfg, problems] = await Promise.all([api('GET', '/config'), api('GET', '/problems')]);
    } catch (e) {
      view.replaceChildren(el('div', { class: 'page' }, el('p', { text: 'Could not reach the dsa server: ' + e.message })));
      return;
    }
    applyTheme();
    langSel.value = cfg.lang;
    route();
  })();

})();
