// Bundled by esbuild into internal/server/web/vendor/editor.js.
// app.js uses window.CarrelEditor.create(parent, { doc, language, onChange, onRun }).
import { EditorState, Prec } from '@codemirror/state';
import {
  EditorView, keymap, lineNumbers, highlightActiveLine, highlightActiveLineGutter, drawSelection,
} from '@codemirror/view';
import { defaultKeymap, history, historyKeymap, indentWithTab } from '@codemirror/commands';
import {
  bracketMatching, indentOnInput, indentUnit, syntaxHighlighting, HighlightStyle,
} from '@codemirror/language';
import { tags as t } from '@lezer/highlight';
import { java } from '@codemirror/lang-java';
import { cpp } from '@codemirror/lang-cpp';

// Colours come from CSS variables, so Paper, Ink and every accent apply without rebuilding.
const highlight = HighlightStyle.define([
  { tag: [t.keyword, t.modifier, t.controlKeyword, t.operatorKeyword, t.definitionKeyword, t.moduleKeyword, t.processingInstruction], color: 'var(--syn-keyword)' },
  { tag: [t.typeName, t.className, t.namespace, t.standard(t.typeName)], color: 'var(--syn-type)' },
  { tag: [t.number, t.bool, t.null, t.string, t.character, t.special(t.string)], color: 'var(--syn-number)' },
  { tag: [t.comment, t.lineComment, t.blockComment, t.docComment], color: 'var(--syn-comment)', fontStyle: 'italic' },
]);

const theme = EditorView.theme({
  '&': { height: '100%', backgroundColor: 'var(--panel)', color: 'var(--ink)' },
  '.cm-scroller': { fontFamily: 'var(--mono)', fontSize: '14px', lineHeight: '1.7' },
  '.cm-content': { padding: '18px 0', caretColor: 'var(--ink)' },
  '.cm-line': { padding: '0 24px 0 4px' },
  '.cm-gutters': { backgroundColor: 'var(--panel)', color: 'var(--gutter)', border: 'none' },
  '.cm-lineNumbers .cm-gutterElement': { minWidth: '40px', padding: '0 12px 0 8px' },
  '.cm-activeLine': { backgroundColor: 'var(--active-line)' },
  '.cm-activeLineGutter': { backgroundColor: 'transparent', color: 'var(--muted)' },
  '&.cm-focused .cm-cursor': { borderLeftColor: 'var(--ink)' },
  '&.cm-focused .cm-selectionBackground, .cm-selectionBackground, ::selection': { backgroundColor: 'var(--selection)' },
  '&.cm-focused .cm-matchingBracket': { backgroundColor: 'transparent', outline: '1px solid var(--muted)' },
  '&.cm-focused .cm-nonmatchingBracket': { backgroundColor: 'transparent', outline: '1px solid var(--fail)' },
  '&.cm-focused': { outline: 'none' },
});

function create(parent, { doc = '', language = 'java', onChange, onRun } = {}) {
  const runKeys = Prec.highest(keymap.of([
    { key: 'Mod-Enter', run: () => { if (onRun) onRun('examples'); return true; } },
    { key: 'Shift-Mod-Enter', run: () => { if (onRun) onRun('submit'); return true; } },
  ]));
  const view = new EditorView({
    parent,
    state: EditorState.create({
      doc,
      extensions: [
        runKeys,
        lineNumbers(),
        highlightActiveLineGutter(),
        highlightActiveLine(),
        drawSelection(),
        history(),
        bracketMatching(),
        indentOnInput(),
        indentUnit.of('    '),
        EditorState.tabSize.of(4),
        keymap.of([indentWithTab, ...defaultKeymap, ...historyKeymap]),
        language === 'cpp' ? cpp() : java(),
        syntaxHighlighting(highlight),
        theme,
        EditorView.contentAttributes.of({ 'aria-label': 'Solution code (' + language + ')' }),
        EditorView.updateListener.of((u) => { if (u.docChanged && onChange) onChange(); }),
      ],
    }),
  });
  return {
    getValue: () => view.state.doc.toString(),
    setValue: (text) => view.dispatch({ changes: { from: 0, to: view.state.doc.length, insert: text } }),
    focus: () => view.focus(),
    destroy: () => view.destroy(),
  };
}

window.CarrelEditor = { create };
