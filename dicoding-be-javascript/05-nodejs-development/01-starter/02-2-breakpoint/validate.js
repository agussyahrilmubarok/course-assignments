const fs = require('fs');
const {join} = require('path');
const {execSync} = require('child_process');

(() => {
  try {
    const packageContent = 'ewogICJuYW1lIjogIjAyLTItYnJlYWtwb2ludCIsCiAgInZlcnNpb24iOiAiMS4wLjAiLAogICJkZXNjcmlwdGlvbiI6ICI+ICoqVHVnYXMqKjogUGFkYSBiZXJrYXMgYGFwcC5qc2AsIHRhbWJhaGthbiBicmVha3BvaW50IHBhZGEgYmFyaXMga29kZSBhd2FsIHNldGlhcCBmdW5nc2kgYGNhbGN1bGF0ZUF2ZXJhZ2VgIGRhbiBgcmVtb3ZlRXZlbk51bWJlcnNgLiIsCiAgIm1haW4iOiAiYXBwLmpzIiwKICAic2NyaXB0cyI6IHsKICAgICJ0ZXN0IjogImVjaG8gXCJFcnJvcjogbm8gdGVzdCBzcGVjaWZpZWRcIiAmJiBleGl0IDEiCiAgfSwKICAia2V5d29yZHMiOiBbXSwKICAiYXV0aG9yIjogIiIsCiAgImxpY2Vuc2UiOiAiSVNDIiwKICAiZGVwZW5kZW5jaWVzIjogewogICAgImVzcHJpbWEiOiAiXjQuMC4xIgogIH0KfQo=';
    const packageJsonPath = join(__dirname, 'package.json');
    fs.writeFileSync(packageJsonPath, packageContent, { encoding: 'base64' });
    execSync('npm install --no-package-lock', {cwd: __dirname});

    const esprima = require('esprima');
    const assert = require('assert');

    const appJsPath = join(__dirname, 'app.js');
    const studentAnswer = fs.readFileSync(appJsPath, 'utf-8');

    const ast = esprima.parseScript(studentAnswer, {loc: true});

    const firstBreakpoint = ast.body[0].body.body[0];
    const secondBreakpoint = ast.body[1].body.body[0];

    assert.equal(firstBreakpoint.type, 'DebuggerStatement', 'Anda harus menambahkan breakpoint pada fungsi `calculateAverage`. Pastikan juga Anda tidak mengubah kode lain selain menambahkan breakpoint.');
    assert.equal(secondBreakpoint.type, 'DebuggerStatement', 'Anda harus menambahkan breakpoint pada fungsi `removeEvenNumbers`. Pastikan juga Anda tidak mengubah kode lain selain menambahkan breakpoint.');

    console.log('passed!');
  } catch (error) {
    console.log(error.message);
  } finally {
    fs.rmSync(join(__dirname, 'node_modules'), { recursive: true, force: true });
    fs.rmSync(join(__dirname, 'package.json'));
  }
})();
