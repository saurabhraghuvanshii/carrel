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

  const UNREACHABLE = 'Cannot reach Carrel. Check that it is still running in your terminal, then try again.';

  async function api(method, path, body) {
    let res;
    try {
      res = await fetch('/api' + path, {
        method,
        headers: { 'Content-Type': 'application/json', 'X-Carrel': '1' },
        body: body ? JSON.stringify(body) : undefined,
      });
    } catch (_) {
      throw new Error(UNREACHABLE);
    }
    const data = await res.json().catch(() => ({}));
    if (!res.ok) throw new Error(data.error || res.statusText);
    return data;
  }

  const cap = (s) => s.charAt(0).toUpperCase() + s.slice(1);
  const pad2 = (n) => String(n).padStart(2, '0');
  const problemLink = (id) => '#/p/' + encodeURIComponent(id);
  const sheetName = (id) => (id === 'real' ? 'Real interviews' : 'Patterns');

  function errorPage(text, retry) {
    return el('div', { class: 'page' }, el('div', { class: 'empty', role: 'alert' },
      el('p', { text }),
      retry ? el('button', { class: 'btn', text: 'Try again', onclick: retry }) : null));
  }

  // ---- state ----

  let cfg = null;
  let problems = [];
  const ui = { sheet: 'patterns', filter: 'all', interview: false };
  let current = null; // the problem on the practice screen
  let pageCleanup = []; // timers of the page being left
  let keepFocus = false; // set when [ or ] opens a problem, so the next press still works
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
      focusLayout: next.focusLayout,
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
    for (const fn of pageCleanup.splice(0)) fn();
    current = null;
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
    document.title = 'Sheets · Carrel';
    let next = null;
    let progress = null;
    try {
      [problems, progress, next] = await Promise.all([
        api('GET', '/problems'),
        api('GET', '/progress'),
        api('GET', '/next?sheet=' + ui.sheet).catch(() => null),
      ]);
    } catch (e) {
      if (token === renderToken) view.replaceChildren(errorPage('Could not load the sheets. ' + e.message, route));
      return;
    }
    if (token !== renderToken) return;

    const sheets = ['patterns', 'real'];
    const notes = { patterns: 'Classic problems in learning order', real: 'Pattern commonly seen in online assessments' };
    const all = problems.filter((p) => p.sheet === ui.sheet);
    const numbers = new Map(all.map((p, i) => [p.id, i + 1]));
    const shown = all.filter((p) => ui.filter === 'all' || p.difficulty === ui.filter);

    const tabs = el('div', { class: 'sheet-tabs' },
      el('div', { class: 'tabset' }, sheets.map((id) =>
        el('button', { class: id === ui.sheet ? 'on' : '', 'aria-pressed': String(id === ui.sheet), text: sheetName(id), onclick: () => { ui.sheet = id; renderSheets(); } }))),
      archiveButtons());
    const note = el('p', { class: 'sheet-note', text: notes[ui.sheet] });

    const chips = el('div', { class: 'chips' }, ['all', 'easy', 'medium', 'hard'].map((f) =>
      el('button', { class: `chip ${f}${ui.filter === f ? ' on' : ''}`, 'aria-pressed': String(ui.filter === f), text: cap(f), onclick: () => { ui.filter = f; renderSheets(); } })));

    const list = el('div');
    if (!problems.length) {
      list.append(el('div', { class: 'empty' },
        el('p', { text: 'No problems found.' }),
        el('p', { class: 'hint', text: 'The built-in problems did not load. Packs you put in the packs folder inside your Carrel folder are listed here too.' })));
    } else if (!all.length) {
      list.append(el('p', { class: 'empty', text: 'No problems in this sheet yet.' }));
    } else if (!shown.length) {
      list.append(el('p', { class: 'empty', text: `No ${ui.filter} problems in this sheet.` }));
    }
    let lastGroup = null;
    for (const p of shown) {
      if (p.group !== lastGroup) {
        list.append(el('div', { class: 'group-title', text: `${Math.floor(p.order / 100)} · ${p.group}` }));
        lastGroup = p.group;
      }
      const needs = p.needs && p.needs.length ? 'Needs: ' + p.needs.map(titleOf).join(', ') : null;
      list.append(el('a', { class: 'list-row', href: problemLink(p.id), title: needs },
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
    const continueBox = next
      ? el('div', { class: 'aside-block' },
        el('div', { class: 'aside-label', text: next.reason === 'continue' ? 'Continue' : 'Up next' }),
        el('div', { class: 'aside-title', text: next.title }),
        el('div', { class: 'hint', text: `${sheetName(next.sheet)} · ${next.group}` }),
        el('a', { class: 'btn primary open', href: problemLink(next.id), text: 'Open' }))
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

    view.replaceChildren(el('div', { class: 'sheets' }, el('div', { class: 'sheets-main' }, el('h1', { class: 'visually-hidden', text: 'Sheets' }), importBox(), tabs, note, chips, list), aside));
  }

  // ---- export and import ----

  let importNotice = null; // { file, report, error }, shown until closed
  const count = (n, one, many) => `${n} ${n === 1 ? one : many}`;

  function archiveButtons() {
    const picker = el('input', { type: 'file', accept: '.zip,application/zip', class: 'visually-hidden', tabindex: '-1', 'aria-label': 'Zip file to import',
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
        method: 'POST', headers: { 'Content-Type': 'application/zip', 'X-Carrel': '1' }, body: file,
      }).catch(() => { throw new Error(UNREACHABLE); });
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

  let doctorCache = null; // compilers are looked up once, and again after a visit to Settings

  function missingTools(lang) {
    if (!doctorCache) doctorCache = api('GET', '/doctor').catch(() => { doctorCache = null; return { tools: [] }; });
    const needed = lang === 'java' ? ['javac', 'java'] : ['g++'];
    return doctorCache.then((d) => d.tools.filter((t) => needed.includes(t.name) && !t.found).map((t) => t.name));
  }

  async function renderPractice(id) {
    const token = renderToken;
    const autofocus = !keepFocus;
    keepFocus = false;
    view.replaceChildren(el('div', { class: 'page' }, el('p', { class: 'hint', text: 'Loading...' })));
    let p;
    try {
      p = await api('GET', `/problems/${encodeURIComponent(id)}?lang=${cfg.lang}`);
      if (token !== renderToken) return;
    } catch (e) {
      if (token !== renderToken) return;
      view.replaceChildren(el('div', { class: 'page' }, el('div', { class: 'empty', role: 'alert' },
        el('p', { text: e.message === UNREACHABLE ? e.message : 'Could not open this problem: ' + e.message }),
        el('a', { class: 'btn', href: '#/', text: 'Back to sheets' }))));
      return;
    }
    current = p;
    document.title = p.title + ' · Carrel';
    const linkList = (label, ids) => ids && ids.length
      ? [label + ' ', ...ids.flatMap((x, i) => [i ? ', ' : '', el('a', { href: problemLink(x), text: titleOf(x) })]), '. ']
      : [];

    // Examples go right after the description; Constraints and anything after it follow them.
    const lines = p.statement.replace(/\r/g, '').split('\n');
    const cut = lines.indexOf('## Constraints');
    const intro = cut < 0 ? p.statement : lines.slice(0, cut).join('\n');
    const rest = cut < 0 ? '' : lines.slice(cut).join('\n');

    const statement = el('div', { class: 'statement' },
      el('div', { class: 'meta-row' }, el('span', { class: 'tag ' + p.difficulty, text: cap(p.difficulty) }), el('span', { class: 'topics', text: (p.tags || []).map(cap).join(', ') })),
      el('h1', { text: p.title }),
      renderMarkdown(intro),
      el('div', { class: 'examples' }, (p.examples || []).map((ex) => el('div', { class: 'example' },
        el('div', { class: 'label', text: ex.label }),
        el('pre', { text: ex.display || ex.input })))),
      renderMarkdown(rest),
      el('div', { class: 'links' }, linkList('Builds on:', p.buildsOn), linkList('Leads to:', p.leadsTo)));

    const fileName = p.lang === 'java' ? 'Solution.java' : 'solution.cpp';
    const savedLabel = el('span', { text: 'Saved on this computer' });
    const retrySave = el('button', { class: 'link', text: 'Try again', hidden: true, onclick: () => save() });
    const savedBox = el('span', { class: 'saved', role: 'status' }, el('i'), savedLabel, retrySave);
    const editorBox = el('div', { class: 'editor' });

    const verdict = el('div', { class: 'verdict', role: 'status', 'aria-live': 'polite' });
    const report = el('div', { class: 'report' },
      el('div', { class: 'hint', text: 'Run the examples to see how your code does. Submit also checks fresh random cases.' }));
    const out = el('div', { class: 'results' }, verdict, report);
    const runBtn = el('button', { class: 'btn strong', text: 'Run examples', onclick: () => run('examples') });
    const submitBtn = el('button', { class: 'btn primary', text: 'Submit', onclick: () => run('submit') });
    const keysHint = 'Ctrl + Enter runs · Ctrl + Shift + Enter submits';
    const barStatus = el('span', { class: 'bar-status hint', 'aria-hidden': 'true', text: keysHint });
    const setBar = (text, cls) => {
      barStatus.textContent = text || keysHint;
      barStatus.className = 'bar-status ' + (text ? cls : 'hint');
    };

    const toolNotice = el('div', { class: 'notice', role: 'status' });
    missingTools(p.lang).then((names) => {
      if (token !== renderToken || !names.length) return;
      toolNotice.replaceChildren(
        `${names.join(' and ')} ${names.length === 1 ? 'was' : 'were'} not found on this computer, so your code cannot run yet. `,
        el('a', { href: '#/settings', text: 'See Languages in Settings.' }));
    });

    let lastReport = '';
    let running = false;
    let submitted = false;
    let stopReveal = () => {};
    pageCleanup.push(() => stopReveal());

    function markSaved() {
      savedBox.classList.remove('failed');
      savedBox.removeAttribute('title');
      savedLabel.textContent = 'Saved on this computer';
      retrySave.hidden = true;
    }
    async function save() {
      try {
        await api('PUT', '/solution', { problem: p.id, lang: p.lang, code: editor.getValue() });
        markSaved();
      } catch (e) {
        savedBox.classList.add('failed');
        savedLabel.textContent = 'Not saved';
        savedBox.title = e.message;
        retrySave.hidden = false;
      }
    }
    function scheduleSave() {
      if (!savedBox.classList.contains('failed')) savedLabel.textContent = 'Saving...';
      if (pendingSave) clearTimeout(pendingSave.timer);
      pendingSave = { fn: async () => { await save(); }, timer: setTimeout(async () => { pendingSave = null; await save(); }, 600) };
    }
    const editor = window.CarrelEditor.create(editorBox, {
      doc: p.code,
      language: p.lang,
      onChange: scheduleSave,
      onRun: (mode) => run(mode),
    });
    activeEditor = editor;

    // Interview mode: no hints, no AI, a timer, and no more example runs after a submit.
    const timer = el('span', { class: 'timer', role: 'timer', 'aria-label': 'Time on this problem' });
    let startedAt = Date.now();
    let stoppedAt = 0;
    const tick = () => {
      const secs = Math.floor(((stoppedAt || Date.now()) - startedAt) / 1000);
      timer.textContent = `${pad2(Math.floor(secs / 60))}:${pad2(secs % 60)}`;
    };
    const ticker = setInterval(tick, 1000);
    pageCleanup.push(() => clearInterval(ticker));
    const interviewSwitch = el('button', { class: 'switch small', role: 'switch', 'aria-label': 'Interview mode',
      onclick: () => {
        ui.interview = !ui.interview;
        startedAt = Date.now();
        stoppedAt = 0;
        applyInterview();
      } });
    function syncButtons() {
      const locked = ui.interview && submitted;
      runBtn.disabled = running || locked;
      submitBtn.disabled = running;
      runBtn.title = locked ? 'Interview mode: the examples are closed after a submit' : '';
    }
    function applyInterview() {
      root.classList.toggle('interview', ui.interview);
      interviewSwitch.setAttribute('aria-checked', String(ui.interview));
      crumb.textContent = ui.interview ? sheetName(p.sheet) : `${sheetName(p.sheet)} / ${p.group}`;
      tick();
      syncButtons();
    }

    async function run(mode) {
      if (running || (mode === 'examples' && ui.interview && submitted)) return;
      stopReveal();
      if (pendingSave) { clearTimeout(pendingSave.timer); pendingSave = null; }
      running = true;
      syncButtons();
      const waiting = mode === 'submit' ? 'Submitting...' : 'Running...';
      verdict.replaceChildren(el('div', { class: 'hint', text: waiting }));
      report.replaceChildren();
      setBar(waiting, 'hint');
      try {
        const rep = await api('POST', '/run', { problem: p.id, lang: p.lang, code: editor.getValue(), mode });
        if (token !== renderToken) return;
        markSaved();
        lastReport = summarise(rep);
        if (mode === 'submit') submitted = true;
        const solved = mode === 'submit' && rep.status === 'ok' && rep.passed === rep.total;
        if (solved) { stoppedAt = Date.now(); tick(); }
        stopReveal = renderReport({ verdict, report, setBar }, rep, () => { if (solved) showNext(verdict, p); });
      } catch (e) {
        verdict.replaceChildren(el('div', { class: 'message', text: e.message }));
        setBar(e.message, 'fail');
      } finally {
        running = false;
        syncButtons();
      }
    }

    const question = el('input', { type: 'text', placeholder: 'Ask for a hint or an explanation', 'aria-label': 'Ask AI' });
    const answer = el('div', { class: 'answer', 'aria-live': 'polite' });
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
    const ai = el('div', { class: 'ai' },
      el('div', { class: 'hint' }, 'Ask AI explains and gives hints. ',
        cfg.ai.hasKey || cfg.ai.provider === 'ollama'
          ? (cfg.ai.explainOnly ? 'It will not write the full solution.' : 'Explain only is off.')
          : el('a', { href: '#/settings', text: 'Add your key in Settings.' })),
      el('div', { class: 'ai-row' }, question, askBtn),
      answer);

    const tabs = el('div', { class: 'tabs' },
      el('span', { class: 'file', text: fileName }),
      el('div', { class: 'tabs-right' },
        timer,
        el('span', { class: 'interview-toggle' }, el('span', { 'aria-hidden': 'true', text: 'Interview mode' }), interviewSwitch),
        savedBox));
    const buttons = el('div', { class: 'group' }, runBtn, submitBtn);

    let root;
    if (cfg.focusLayout) {
      root = el('div', { class: 'practice focus' },
        el('div', { class: 'focus-scroll' },
          el('div', { class: 'focus-col' }, statement, el('div', { class: 'editor-box' }, tabs, editorBox), toolNotice, out, ai)),
        el('div', { class: 'focus-bar' }, el('div', { class: 'focus-col' }, barStatus, buttons)));
    } else {
      root = el('div', { class: 'practice' }, statement,
        el('div', { class: 'work' }, tabs, editorBox,
          el('div', { class: 'console' },
            toolNotice,
            el('div', { class: 'actions' }, buttons, el('span', { class: 'hint', text: keysHint })),
            out,
            ai)));
    }
    applyInterview();
    view.replaceChildren(root);
    if (autofocus) editor.focus();
  }

  async function showNext(verdict, p) {
    const line = el('div', { class: 'next-line' }, 'Solved.');
    verdict.append(line);
    try {
      const n = await api('GET', `/next?sheet=${p.sheet}&after=${encodeURIComponent(p.id)}`);
      line.append(' Next: ', el('a', { href: problemLink(n.id), text: n.title }));
    } catch (_) {
      line.append(' Every problem in this sheet is solved.');
    }
  }

  function summarise(rep) {
    const bad = rep.results.filter((r) => !r.passed).slice(0, 3)
      .map((r) => `${r.label}: input ${JSON.stringify(r.input || '')}, expected ${r.expected}, got ${r.got}`);
    return `Status ${rep.status}. Passed ${rep.passed} of ${rep.total}. ${rep.message || ''} ${bad.join(' | ')}`.trim();
  }

  const reducedMotion = matchMedia('(prefers-reduced-motion: reduce)');
  const REVEAL_MS = 900;
  const SQUARE_MS = 120; // the same as sq-fill in style.css

  function caseRow(r, label) {
    let detail = `input\n${r.input}\nexpected ${r.expected}\ngot      ${r.got || '(nothing)'}`;
    if (r.note) detail += `\nwhy      ${r.note}`;
    const word = r.passed ? 'Passed' : r.got === '(crashed)' ? 'Crashed' : r.got === '(not run)' ? 'Not run' : 'Failed';
    return el('div', { class: 'row' },
      el('span', { class: 'hint', text: label || r.label }),
      el('span', { class: 'detail', text: detail }),
      el('span', { class: r.passed ? 'ok' : 'bad', text: word }));
  }

  // Draws a run report and returns a function that stops a reveal still playing.
  function renderReport({ verdict, report, setBar }, rep, done) {
    const say = (cls, text, ...more) => {
      verdict.replaceChildren(el('div', { class: 'summary ' + cls },
        cls === 'pass' ? el('span', { class: 'check', 'aria-hidden': 'true' }) : null, text), ...more);
      setBar(text, cls);
    };
    const failed = rep.results.filter((r) => !r.passed);

    if (rep.status === 'compile_error') {
      say('fail', 'It did not compile');
      report.replaceChildren(el('pre', { class: 'message', text: rep.message }));
      return () => {};
    }
    if (rep.status === 'tooling_missing' || rep.status === 'internal_error') {
      say('fail', 'Could not run your code');
      report.replaceChildren(el('div', { class: 'message' }, rep.message,
        rep.status === 'tooling_missing' ? [' ', el('a', { href: '#/settings', text: 'See Languages in Settings.' })] : null));
      return () => {};
    }
    if (rep.status !== 'ok') {
      const lead = { memory_limit: 'Memory limit reached. ', timeout: 'Time limit reached. ', runtime_error: 'Your program crashed. ' }[rep.status] || '';
      say('fail', `${lead}Passed ${rep.passed} of ${rep.total}`);
      report.replaceChildren(
        el('div', { class: 'message', text: rep.message }),
        ...failed.map((r) => caseRow(r)),
        el('div', { class: 'hint', text: `Took ${rep.durationMs} ms` }));
      return () => {};
    }

    const n = rep.results.length;
    const firstBad = rep.results.findIndex((r) => !r.passed);
    const fixed = rep.results.filter((r) => r.kind !== 'random').length;
    const facts = [`${fixed} fixed`];
    if (n > fixed) facts.push(`${n - fixed} random`);
    if (rep.seed) facts.push(`seed ${rep.seed}`);
    facts.push(`${rep.durationMs} ms`);

    const detail = el('div', { class: 'case' });
    const open = (i) => {
      squares.forEach((b, j) => b.setAttribute('aria-pressed', String(i === j)));
      detail.replaceChildren(caseRow(rep.results[i], `Case ${i + 1} · ${rep.results[i].label}`));
    };
    const squares = rep.results.map((r, i) => {
      const b = el('button', { type: 'button', class: 'sq ' + (r.passed ? 'pass' : 'fail'), 'aria-pressed': 'false',
        'aria-label': `Case ${i + 1}, ${r.kind}, ${r.passed ? 'passed' : 'failed'}`, onclick: () => open(i) });
      b.style.setProperty('--i', i);
      return b;
    });
    const step = Math.min(14, (REVEAL_MS - SQUARE_MS) / n);
    const grid = el('div', { class: 'grid', role: 'group', 'aria-label': 'Test cases' }, squares);
    grid.style.setProperty('--step', step + 'ms');
    const failedList = failed.length
      ? el('ul', { class: 'visually-hidden', 'aria-label': 'Failed cases' }, rep.results.map((r, i) => r.passed ? null
        : el('li', { text: `Case ${i + 1}, ${r.label}, failed. Input: ${r.input}. Expected: ${r.expected}. Got: ${r.got || 'nothing'}.` })))
      : null;
    report.replaceChildren(grid, detail, failedList);

    const finish = () => {
      grid.classList.remove('reveal');
      if (firstBad < 0) {
        say('pass', `Passed ${n} of ${n}`, el('div', { class: 'hint', text: facts.join(', ') }));
      } else {
        say('fail', `Failed on case ${firstBad + 1} of ${n}`, el('div', { class: 'hint', text: `${rep.passed} of ${n} passed, ${facts.join(', ')}` }));
        open(firstBad);
      }
      done();
    };
    if (reducedMotion.matches || !n) {
      finish();
      return () => {};
    }

    // The server answers with every result at once, so this is a short reveal, not progress.
    const counter = el('div', { class: 'summary', 'aria-hidden': 'true' });
    verdict.replaceChildren(counter);
    grid.classList.add('reveal');
    const total = (n - 1) * step + SQUARE_MS;
    const start = performance.now();
    let frame = requestAnimationFrame(function draw(now) {
      if (now - start >= total) return finish();
      counter.textContent = `Running case ${Math.min(n, Math.floor(Math.max(0, now - start) / step) + 1)} of ${n}`;
      setBar(counter.textContent, 'hint');
      frame = requestAnimationFrame(draw);
    });
    return () => cancelAnimationFrame(frame);
  }

  // ---- quick search and keys ----

  let closeSearch = null;

  function openSearch() {
    if (closeSearch) return;
    const back = document.activeElement;
    const input = el('input', { type: 'text', placeholder: 'Find a problem', 'aria-label': 'Find a problem', autocomplete: 'off',
      role: 'combobox', 'aria-expanded': 'true', 'aria-controls': 'search-list', 'aria-autocomplete': 'list' });
    const list = el('ul', { id: 'search-list', role: 'listbox', 'aria-label': 'Matching problems' });
    const none = el('div', { class: 'hint', role: 'status' });
    const shade = el('div', { class: 'shade', onmousedown: (e) => { if (e.target === shade) close(); } },
      el('div', { class: 'search', role: 'dialog', 'aria-modal': 'true', 'aria-label': 'Find a problem' },
        input, list, none, el('div', { class: 'hint', text: 'Up and down to choose · Enter opens · Esc closes' })));
    let hits = [];
    let at = 0;
    function close() {
      shade.remove();
      closeSearch = null;
      if (back && back.isConnected) back.focus();
    }
    const go = (p) => { close(); location.hash = problemLink(p.id); };
    function draw() {
      const q = input.value.trim().toLowerCase();
      hits = problems.filter((p) => !q || [p.title, p.group, ...(p.tags || [])].some((s) => s.toLowerCase().includes(q))).slice(0, 8);
      at = Math.max(0, Math.min(at, hits.length - 1));
      list.replaceChildren(...hits.map((p, i) => el('li', { id: 'search-hit-' + i, role: 'option', 'aria-selected': String(i === at), onclick: () => go(p) },
        el('span', { text: p.title }),
        el('span', { class: 'hint', text: `${sheetName(p.sheet)} · ${p.group}` }))));
      none.textContent = hits.length ? '' : 'No problem matches.';
      if (hits.length) input.setAttribute('aria-activedescendant', 'search-hit-' + at);
      else input.removeAttribute('aria-activedescendant');
    }
    input.addEventListener('input', () => { at = 0; draw(); });
    input.addEventListener('keydown', (e) => {
      if (e.key === 'Escape') close();
      else if (e.key === 'ArrowDown') at++;
      else if (e.key === 'ArrowUp') at--;
      else if (e.key === 'Enter' && hits[at]) go(hits[at]);
      else if (e.key !== 'Tab') return;
      e.preventDefault();
      e.stopPropagation();
      if (closeSearch) draw();
    });
    closeSearch = close;
    document.body.append(shade);
    draw();
    input.focus();
  }

  function stepProblem(by) {
    const sheet = problems.filter((x) => x.sheet === current.sheet);
    const to = sheet[sheet.findIndex((x) => x.id === current.id) + by];
    if (!to) return;
    keepFocus = true;
    location.hash = problemLink(to.id);
  }

  let gPressedAt = 0;
  document.addEventListener('keydown', (e) => {
    if (e.ctrlKey || e.metaKey || e.altKey || closeSearch) return;
    const target = e.target instanceof Element ? e.target : null;
    const inEditor = target && target.closest('.cm-editor');
    if (e.key === 'Escape') {
      if (activeEditor && !inEditor) { e.preventDefault(); activeEditor.focus(); }
      return;
    }
    // Single keys belong to whatever the learner is typing in.
    if (inEditor || (target && target.closest('input, textarea, select, [contenteditable="true"]'))) return;
    const afterG = Date.now() - gPressedAt < 1500;
    gPressedAt = e.key === 'g' ? Date.now() : 0;
    if (e.key === '/') { e.preventDefault(); openSearch(); }
    else if (e.key === 's' && afterG) location.hash = '#/';
    else if ((e.key === '[' || e.key === ']') && current) stepProblem(e.key === ']' ? 1 : -1);
  });

  // ---- settings ----

  async function renderSettings() {
    const token = renderToken;
    crumb.textContent = 'Settings';
    document.title = 'Settings · Carrel';
    doctorCache = null;
    let doctorError = '';
    const doctor = await api('GET', '/doctor').catch((e) => { doctorError = e.message; return { tools: [], solutionsDir: '' }; });
    if (token !== renderToken) return;
    const note = el('div', { class: 'note', role: 'status' });

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
          el('button', { class: 'swatch' + (cfg.accent === id ? ' on' : ''), style: 'background:' + color, 'aria-label': cap(id), 'aria-pressed': String(cfg.accent === id),
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
      doctorError ? el('div', { class: 'message', role: 'alert', text: 'Could not check the compilers. ' + doctorError }) : null,
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
    const focus = el('button', { class: 'switch', role: 'switch', 'aria-checked': String(!!cfg.focusLayout), 'aria-label': 'Focus layout',
      onclick: async () => {
        try {
          await saveConfig({ focusLayout: !cfg.focusLayout });
          focus.setAttribute('aria-checked', String(cfg.focusLayout));
        } catch (e) {
          practiceNote.textContent = e.message;
        }
      } });
    const practiceNote = el('div', { class: 'note', role: 'status' });
    const practice = el('section', { class: 'sec' },
      el('h2', { text: 'Practice' }),
      el('div', { class: 'switch-row' },
        el('div', {}, el('div', { text: 'Focus layout' }),
          el('div', { class: 'note', text: 'One centred column: the problem, then the editor, with Run and Submit in a bar at the bottom.' })),
        focus),
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

    const keys = [
      ['Ctrl + Enter', 'Run the examples (Cmd on a Mac)'],
      ['Ctrl + Shift + Enter', 'Submit'],
      ['[ and ]', 'Previous and next problem in the sheet'],
      ['/', 'Find a problem'],
      ['g, then s', 'Go to Sheets'],
      ['Esc', 'Back to the editor'],
      ['Esc, then Tab', 'Leave the editor'],
    ];
    const keyList = el('section', { class: 'sec' },
      el('h2', { text: 'Keys' }),
      el('dl', { class: 'keys' }, keys.flatMap(([key, what]) => [el('dt', {}, el('kbd', { text: key })), el('dd', { text: what })])),
      el('div', { class: 'note', text: 'The single keys are off while you type in the editor or in a text box.' }));

    view.replaceChildren(el('div', { class: 'page' }, el('h1', { class: 'visually-hidden', text: 'Settings' }), appearance, practice, keyList, ai, langs, files));
  }

  // ---- start ----

  document.getElementById('skip').addEventListener('click', () => view.focus());

  langSel.addEventListener('change', async () => {
    await flushSave();
    await saveConfig({ lang: langSel.value });
    route();
  });

  (async () => {
    try {
      [cfg, problems] = await Promise.all([api('GET', '/config'), api('GET', '/problems')]);
    } catch (e) {
      view.replaceChildren(errorPage(e.message, () => location.reload()));
      return;
    }
    applyTheme();
    langSel.value = cfg.lang;
    route();
  })();

})();
