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

  // ---- markdown (just paragraphs, lists, headings and `code`) ----

  function inline(text) {
    return text
      .split(/(`[^`]+`)/g)
      .filter(Boolean)
      .map((s) => (s.length > 1 && s.startsWith('`') && s.endsWith('`') ? el('code', { text: s.slice(1, -1) }) : document.createTextNode(s)));
  }

  function renderMarkdown(src) {
    const nodes = [];
    const lines = src.replace(/\r/g, '').split('\n');
    let i = 0;
    const special = (l) => l.startsWith('## ') || l.startsWith('- ');
    while (i < lines.length) {
      const line = lines[i];
      if (!line.trim()) { i++; continue; }
      if (line.startsWith('## ')) { nodes.push(el('h2', {}, ...inline(line.slice(3)))); i++; continue; }
      if (line.startsWith('- ')) {
        const items = [];
        while (i < lines.length && lines[i].startsWith('- ')) { items.push(el('li', {}, ...inline(lines[i].slice(2)))); i++; }
        nodes.push(el('ul', {}, ...items));
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

  function renderSheets() {
    crumb.textContent = 'Sheets';
    const sheets = [['patterns', 'Patterns'], ['real', 'Real interviews']];
    const all = problems.filter((p) => p.sheet === ui.sheet);
    const numbers = new Map(all.map((p, i) => [p.id, i + 1]));
    const shown = all.filter((p) => ui.filter === 'all' || p.difficulty === ui.filter);

    const tabs = el('div', { class: 'sheet-tabs' },
      el('div', { class: 'tabset' }, sheets.map(([id, name]) =>
        el('button', { class: id === ui.sheet ? 'on' : '', text: name, onclick: () => { ui.sheet = id; renderSheets(); } }))));

    const chips = el('div', { class: 'chips' }, ['all', 'easy', 'medium', 'hard'].map((f) =>
      el('button', { class: `chip ${f}${ui.filter === f ? ' on' : ''}`, text: cap(f), onclick: () => { ui.filter = f; renderSheets(); } })));

    const list = el('div');
    if (!shown.length) list.append(el('p', { class: 'hint', text: 'No problems here yet.' }));
    let lastGroup = null;
    for (const p of shown) {
      if (p.group !== lastGroup) {
        list.append(el('div', { class: 'group-title', text: p.group }));
        lastGroup = p.group;
      }
      const label = p.status === 'solved' ? 'Solved' : p.status === 'tried' ? 'Tried' : 'Not started';
      list.append(el('a', { class: 'list-row', href: '#/p/' + encodeURIComponent(p.id) },
        el('span', { class: 'num', text: pad2(numbers.get(p.id)) }),
        el('span', { text: p.title }),
        el('span', { class: 'dif ' + p.difficulty, text: cap(p.difficulty) }),
        el('span', { class: 'tags', text: (p.tags || []).join(', ') }),
        el('span', { class: 'status' }, el('span', { class: 'dot ' + (p.status || '') }), label)));
    }
    view.replaceChildren(el('div', { class: 'page' }, tabs, chips, list));
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
    const titleOf = (pid) => (problems.find((x) => x.id === pid) || { title: pid }).title;
    const linkList = (label, ids) => ids && ids.length
      ? [label + ' ', ...ids.flatMap((x, i) => [i ? ', ' : '', el('a', { href: '#/p/' + encodeURIComponent(x), text: titleOf(x) })]), '. ']
      : [];

    const statement = el('div', { class: 'statement' },
      el('div', { class: 'meta-row' }, el('span', { class: 'tag ' + p.difficulty, text: cap(p.difficulty) }), el('span', { text: (p.tags || []).map(cap).join(', ') })),
      el('h1', { text: p.title }),
      renderMarkdown(p.statement),
      (p.examples || []).map((ex) => el('div', { class: 'example' },
        el('div', { class: 'label', text: ex.label }),
        el('pre', { text: ex.display || ex.input }))),
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
      answer.textContent = 'Thinking...';
      try {
        const r = await api('POST', '/ai', { problem: p.id, lang: p.lang, code: editor.getValue(), question: question.value, lastReport });
        answer.textContent = r.answer;
      } catch (e) {
        answer.textContent = e.message;
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
          el('div', { class: 'hint', text: 'Ask AI explains and gives hints. Set it up in Settings.' }),
          el('div', { class: 'ai-row' }, question, askBtn),
          answer)));

    view.replaceChildren(el('div', { class: 'practice' }, statement, work));
    editor.focus();
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
    let headline = `Passed ${rep.passed} of ${rep.total}`;
    if (allPassed && mode === 'submit') headline = `All ${rep.total} tests passed. Marked solved.`;
    else if (allPassed) headline = `Examples ${rep.passed} of ${rep.total} passed`;
    kids.push(el('div', { class: 'summary ' + (allPassed ? 'pass' : 'fail'), text: headline }));
    if (rep.message) kids.push(el('div', { class: 'message', text: rep.message }));

    const randomPassed = rep.results.filter((r) => r.kind === 'random' && r.passed).length;
    for (const r of rep.results) {
      if (r.kind === 'random' && r.passed) continue;
      const detail = r.passed
        ? `got ${r.got}`
        : `input\n${r.input}\nexpected ${r.expected}\ngot      ${r.got || '(nothing)'}`;
      kids.push(el('div', { class: 'row' },
        el('span', { class: 'hint', text: r.label }),
        el('span', { class: 'detail', text: detail }),
        el('span', { class: r.passed ? 'ok' : 'bad', text: r.passed ? 'Passed' : 'Failed' })));
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
          cfg.ai.hasKey ? el('button', { class: 'btn', text: 'Remove saved key', onclick: () => saveAI({ clearKey: true, apiKey: '' }) }) : null)),
      note);

    const langs = el('section', { class: 'sec' },
      el('h2', { text: 'Languages' }),
      el('div', { class: 'tools' }, doctor.tools.map((t) => [
        el('span', { text: t.name }),
        el('span', { class: 'path', text: t.found ? `${t.version || ''} ${t.path}`.trim() : 'Not found. Install it and make sure it is on your PATH.' }),
        el('span', { class: t.found ? 'ok' : 'bad', style: 'font-weight:500;color:var(--' + (t.found ? 'pass' : 'fail') + ')', text: t.found ? 'Found' : 'Missing' }),
      ]).flat()));

    const files = el('section', { class: 'sec' },
      el('h2', { text: 'Your solutions' }),
      el('div', {}, el('div', { class: 'field-label', text: 'Saved in' }), el('div', { class: 'mono', text: doctor.solutionsDir }),
        el('div', { class: 'note', text: 'Plain files. Open them in any editor, or put the folder in git.' })));

    view.replaceChildren(el('div', { class: 'page' }, appearance, ai, langs, files));
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

  // Keep the sheets list fresh after practising.
  window.addEventListener('hashchange', async () => {
    if ((location.hash || '#/') === '#/') {
      try { problems = await api('GET', '/problems'); renderSheets(); } catch (_) { /* ignore */ }
    }
  });
})();
