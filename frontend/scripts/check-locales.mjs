import assert from 'node:assert/strict';
import { readFileSync, readdirSync } from 'node:fs';
import { join } from 'node:path';
import ts from 'typescript';
import { parse } from '@vue/compiler-sfc';
import { parse as parseTemplate } from '@vue/compiler-dom';

const compiled = ts.transpileModule(readFileSync('app/i18n/messages.ts', 'utf8'), {
  compilerOptions: { module: ts.ModuleKind.ESNext },
}).outputText;
const { messages, translate } = await import(`data:text/javascript;base64,${Buffer.from(compiled).toString('base64')}`);
for (const [key, values] of Object.entries(messages)) {
  assert.equal(values.length, 2, key);
  for (const [index, locale] of ['kk', 'en'].entries()) {
    assert.ok(values[index].trim(), `${locale}: ${key}`);
    assert.equal(translate(locale, key), values[index]);
  }
  assert.equal(translate('ru', key), key);
}
assert.equal(translate('en', 'Customer-provided task content'), 'Customer-provided task content');

const missing = new Set();
function checkStrings(source) {
  const ast = ts.createSourceFile('strings.ts', source, ts.ScriptTarget.Latest, true);
  function visit(node) {
    if (ts.isStringLiteral(node) && /[А-Яа-яЁё]/.test(node.text) && !messages[node.text]) missing.add(node.text);
    ts.forEachChild(node, visit);
  }
  visit(ast);
}
function walk(dir) {
  for (const entry of readdirSync(dir, { withFileTypes: true })) {
    const file = join(dir, entry.name);
    if (entry.isDirectory()) { walk(file); continue; }
    if (!file.endsWith('.vue') || file.endsWith('LanguageSwitcher.vue')) continue;
    const { descriptor, errors } = parse(readFileSync(file, 'utf8'));
    assert.deepEqual(errors, [], file);
    if (descriptor.scriptSetup) checkStrings(descriptor.scriptSetup.content);
    function visit(node) {
      if (node.type === 2) assert.ok(!/[А-Яа-яЁё]/.test(node.content), `${file}: untranslated text ${node.content}`);
      if (node.type === 5) checkStrings(node.content.content);
      for (const prop of node.props || []) {
        if (prop.type === 6 && prop.value) assert.ok(!/[А-Яа-яЁё]/.test(prop.value.content), `${file}: untranslated attribute`);
        if (prop.type === 7 && prop.exp) checkStrings(prop.exp.content);
      }
      for (const child of node.children || []) visit(child);
    }
    if (descriptor.template) visit(parseTemplate(descriptor.template.content));
  }
}
walk('app');
checkStrings(readFileSync('app/utils/presentation.ts', 'utf8'));
assert.deepEqual([...missing], [], 'Missing translation keys');
console.log(`RU/KK/EN: ${Object.keys(messages).length} messages, template coverage and fallback checked.`);
