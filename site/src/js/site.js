(function () {
  var root = document.documentElement;
  var THEMES = ['paper', 'ink'];
  var ACCENTS = ['brick', 'forest', 'amber', 'slate', 'plum'];

  function save(key, val) {
    try { localStorage.setItem(key, val); } catch (e) {}
  }

  function sync() {
    var theme = root.getAttribute('data-theme');
    var accent = root.getAttribute('data-accent');
    document.querySelectorAll('[data-theme-set]').forEach(function (b) {
      var on = b.getAttribute('data-theme-set') === theme;
      b.setAttribute('aria-pressed', on ? 'true' : 'false');
      b.classList.toggle('primary', on);
    });
    document.querySelectorAll('[data-accent-set]').forEach(function (b) {
      b.setAttribute('aria-pressed', b.getAttribute('data-accent-set') === accent ? 'true' : 'false');
    });
  }

  document.addEventListener('click', function (ev) {
    var t = ev.target.closest('[data-theme-set], [data-accent-set], [data-copy]');
    if (!t) return;
    if (t.hasAttribute('data-theme-set')) {
      var theme = t.getAttribute('data-theme-set');
      if (THEMES.indexOf(theme) < 0) return;
      root.setAttribute('data-theme', theme);
      save('carrel-theme', theme);
      sync();
    } else if (t.hasAttribute('data-accent-set')) {
      var accent = t.getAttribute('data-accent-set');
      if (ACCENTS.indexOf(accent) < 0) return;
      root.setAttribute('data-accent', accent);
      save('carrel-accent', accent);
      sync();
    } else {
      copy(t);
    }
  });

  function fallbackCopy(text) {
    var ta = document.createElement('textarea');
    ta.value = text;
    ta.setAttribute('readonly', '');
    ta.style.position = 'fixed';
    ta.style.opacity = '0';
    document.body.appendChild(ta);
    ta.select();
    var ok = false;
    try { ok = document.execCommand('copy'); } catch (e) {}
    document.body.removeChild(ta);
    return ok ? Promise.resolve() : Promise.reject();
  }

  function copy(btn) {
    var target = document.querySelector(btn.getAttribute('data-copy-target'));
    if (!target) return;
    var text = target.textContent.trim();
    var p = navigator.clipboard && window.isSecureContext
      ? navigator.clipboard.writeText(text).catch(function () { return fallbackCopy(text); })
      : fallbackCopy(text);
    p.then(function () {
      btn.textContent = 'Copied';
      clearTimeout(btn._t);
      btn._t = setTimeout(function () { btn.textContent = 'Copy'; }, 1600);
    }, function () {});
  }

  var loops = document.querySelectorAll('[data-loop]');
  var visible = new Map();
  function pauseAll() {
    loops.forEach(function (el) {
      el.classList.toggle('is-paused', document.hidden || !visible.get(el));
    });
  }
  if ('IntersectionObserver' in window) {
    var io = new IntersectionObserver(function (entries) {
      entries.forEach(function (e) { visible.set(e.target, e.isIntersecting); });
      pauseAll();
    });
    loops.forEach(function (el) { io.observe(el); });
    document.addEventListener('visibilitychange', pauseAll);
  }

  sync();
})();
