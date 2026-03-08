const { join } = require('path');
const { spawnSync, execSync } = require('child_process');
const fs = require('fs');

(() => {
  try {
    const originalCodeOfUtils = 'ZnVuY3Rpb24gZ2VuZXJhdGVSYW5kb21TdHJpbmcobGVuZ3RoLCBjYWxsYmFjaykgewogIHNldFRpbWVvdXQoKCkgPT4gewogICAgaWYgKHR5cGVvZiBsZW5ndGggIT09ICdudW1iZXInKSB7CiAgICAgIGNhbGxiYWNrKEVycm9yKCdMZW5ndGggbXVzdCBiZSBhIG51bWJlcicpKTsKICAgICAgcmV0dXJuOwogICAgfQoKICAgIGlmIChsZW5ndGggPCAxKSB7CiAgICAgIGNhbGxiYWNrKEVycm9yKCdMZW5ndGggbXVzdCBiZSBncmVhdGVyIHRoYW4gMCcpKTsKICAgICAgcmV0dXJuOwogICAgfQoKICAgIGNvbnN0IGNoYXJzZXQgPSAnQUJDREVGR0hJSktMTU5PUFFSU1RVVldYWVphYmNkZWZnaGlqa2xtbm9wcXJzdHV2d3h5ejAxMjM0NTY3ODknOwogICAgbGV0IHJlc3VsdCA9ICcnOwoKICAgIGZvciAobGV0IGkgPSAwOyBpIDwgbGVuZ3RoOyBpICs9IDEpIHsKICAgICAgY29uc3QgcmFuZG9tSW5kZXggPSBNYXRoLmZsb29yKE1hdGgucmFuZG9tKCkgKiBjaGFyc2V0Lmxlbmd0aCk7CiAgICAgIHJlc3VsdCArPSBjaGFyc2V0LmNoYXJBdChyYW5kb21JbmRleCk7CiAgICB9CgogICAgY2FsbGJhY2sobnVsbCwgcmVzdWx0KTsKICB9LCAxMDApOwp9Cgptb2R1bGUuZXhwb3J0cyA9IHsgZ2VuZXJhdGVSYW5kb21TdHJpbmcgfTsK';

    // check availability of package.json
    const packageJsonPath = join(__dirname, 'package.json');

    if (!fs.existsSync(packageJsonPath)) {
      throw new Error(`Pastikan kamu tidak menghapus berkas \`package.json\` pada folder ${__dirname}. Silakan dicek lagi, ya! ^_^`);
    }

    // node_modules should not exist
    const nodeModulesPath = join(__dirname, 'node_modules');

    if (fs.existsSync(nodeModulesPath)) {
      throw new Error(`Pastikan kamu menghapus folder \`node_modules\` pada folder ${__dirname}. Silakan dicek lagi, ya! ^_^`);
    }

    // install dependencies of the submission project
    execSync('npm install', { cwd: __dirname });

    // restore the original code of utils.js
    fs.writeFileSync(join(__dirname, 'utils.js'), Buffer.from(originalCodeOfUtils, 'base64').toString());

    // run test
    const { status: _status, stderr: _stderr, stdout: _stdout } = spawnSync('npm', ['test'], { cwd: __dirname, encoding: 'utf-8' });

    if (_status !== 0) {
      if (_stderr !== '') {
        throw new Error(`Pastikan pengujian dapat dijalankan dengan menggunakan perintah \`npm test\` dengan hasil yang lolos, ya! Berikut adalah pengujian error yang dihasilkan. <pre>${_stderr}</pre>`);
      }

      throw new Error(`Pastikan pengujian dapat dijalankan dengan menggunakan perintah \`npm test\` dengan hasil yang lolos, ya!. Berikut adalah pengujian error yang dihasilkan. <pre>${_stdout}</pre>`);
    }

    // test bad output
    const badOutputCode = 'ZnVuY3Rpb24gZ2VuZXJhdGVSYW5kb21TdHJpbmcobGVuZ3RoLCBjYWxsYmFjaykgewogIHNldFRpbWVvdXQoKCkgPT4gewogICAgaWYgKHR5cGVvZiBsZW5ndGggIT09ICdudW1iZXInKSB7CiAgICAgIGNhbGxiYWNrKEVycm9yKCdMZW5ndGggbXVzdCBiZSBhIG51bWJlcicpKTsKICAgICAgcmV0dXJuOwogICAgfQoKICAgIGlmIChsZW5ndGggPCAxKSB7CiAgICAgIGNhbGxiYWNrKEVycm9yKCdMZW5ndGggbXVzdCBiZSBncmVhdGVyIHRoYW4gMCcpKTsKICAgICAgcmV0dXJuOwogICAgfQoKICAgIGNvbnN0IGNoYXJzZXQgPSAnQUJDREVGR0hJSktMTU5PUFFSU1RVVldYWVphYmNkZWZnaGlqa2xtbm9wcXJzdHV2d3h5ejAxMjM0NTY3ODknOwogICAgbGV0IHJlc3VsdCA9ICcnOwoKICAgIGZvciAobGV0IGkgPSAwOyBpIDwgbGVuZ3RoOyBpICs9IDEpIHsKICAgICAgY29uc3QgcmFuZG9tSW5kZXggPSBNYXRoLmZsb29yKE1hdGgucmFuZG9tKCkgKiBjaGFyc2V0Lmxlbmd0aCk7CiAgICAgIHJlc3VsdCArPSBjaGFyc2V0LmNoYXJBdChyYW5kb21JbmRleCk7CiAgICB9CgogICAgY2FsbGJhY2sobnVsbCwgJycpOwogIH0sIDEwMCk7Cn0KCm1vZHVsZS5leHBvcnRzID0geyBnZW5lcmF0ZVJhbmRvbVN0cmluZyB9Owo=';

    fs.writeFileSync(join(__dirname, 'utils.js'), Buffer.from(badOutputCode, 'base64').toString());

    // run test
    const { status: __status } = spawnSync('npm', ['test'], { cwd: __dirname });

    if (__status === 0) {
      throw new Error('Pastikan Anda menguji nilai yang dikembalikan oleh fungsi `generateRandomString` dengan benar, ya!')
    }

    // test bad data type
    const badDataTypeCode = 'ZnVuY3Rpb24gZ2VuZXJhdGVSYW5kb21TdHJpbmcobGVuZ3RoLCBjYWxsYmFjaykgewogIHNldFRpbWVvdXQoKCkgPT4gewogICAgaWYgKHR5cGVvZiBsZW5ndGggIT09ICdudW1iZXInKSB7CiAgICAgIGNhbGxiYWNrKG51bGwsICcnKTsKICAgICAgcmV0dXJuOwogICAgfQoKICAgIGlmIChsZW5ndGggPCAxKSB7CiAgICAgIGNhbGxiYWNrKEVycm9yKCdMZW5ndGggbXVzdCBiZSBncmVhdGVyIHRoYW4gMCcpKTsKICAgICAgcmV0dXJuOwogICAgfQoKICAgIGNvbnN0IGNoYXJzZXQgPSAnQUJDREVGR0hJSktMTU5PUFFSU1RVVldYWVphYmNkZWZnaGlqa2xtbm9wcXJzdHV2d3h5ejAxMjM0NTY3ODknOwogICAgbGV0IHJlc3VsdCA9ICcnOwoKICAgIGZvciAobGV0IGkgPSAwOyBpIDwgbGVuZ3RoOyBpICs9IDEpIHsKICAgICAgY29uc3QgcmFuZG9tSW5kZXggPSBNYXRoLmZsb29yKE1hdGgucmFuZG9tKCkgKiBjaGFyc2V0Lmxlbmd0aCk7CiAgICAgIHJlc3VsdCArPSBjaGFyc2V0LmNoYXJBdChyYW5kb21JbmRleCk7CiAgICB9CgogICAgY2FsbGJhY2sobnVsbCwgcmVzdWx0KTsKICB9LCAxMDApOwp9Cgptb2R1bGUuZXhwb3J0cyA9IHsgZ2VuZXJhdGVSYW5kb21TdHJpbmcgfTsK';

    fs.writeFileSync(join(__dirname, 'utils.js'), Buffer.from(badDataTypeCode, 'base64').toString());

    // run test
    const { status: ___status } = spawnSync('npm', ['test'], { cwd: __dirname });

    if (___status === 0) {
      throw new Error('Pastikan Anda menguji skenario ketika fungsi `generateRandomString` diberikan dengan argumen non-number, ya!')
    }

    // test bad length
    const badLengthCode = 'ZnVuY3Rpb24gZ2VuZXJhdGVSYW5kb21TdHJpbmcobGVuZ3RoLCBjYWxsYmFjaykgewogIHNldFRpbWVvdXQoKCkgPT4gewogICAgaWYgKHR5cGVvZiBsZW5ndGggIT09ICdudW1iZXInKSB7CiAgICAgIGNhbGxiYWNrKEVycm9yKCdMZW5ndGggbXVzdCBiZSBhIG51bWJlcicpKTsKICAgICAgcmV0dXJuOwogICAgfQoKICAgIGlmIChsZW5ndGggPCAxKSB7CiAgICAgIGNhbGxiYWNrKG51bGwsICcnKTsKICAgICAgcmV0dXJuOwogICAgfQoKICAgIGNvbnN0IGNoYXJzZXQgPSAnQUJDREVGR0hJSktMTU5PUFFSU1RVVldYWVphYmNkZWZnaGlqa2xtbm9wcXJzdHV2d3h5ejAxMjM0NTY3ODknOwogICAgbGV0IHJlc3VsdCA9ICcnOwoKICAgIGZvciAobGV0IGkgPSAwOyBpIDwgbGVuZ3RoOyBpICs9IDEpIHsKICAgICAgY29uc3QgcmFuZG9tSW5kZXggPSBNYXRoLmZsb29yKE1hdGgucmFuZG9tKCkgKiBjaGFyc2V0Lmxlbmd0aCk7CiAgICAgIHJlc3VsdCArPSBjaGFyc2V0LmNoYXJBdChyYW5kb21JbmRleCk7CiAgICB9CgogICAgY2FsbGJhY2sobnVsbCwgcmVzdWx0KTsKICB9LCAxMDApOwp9Cgptb2R1bGUuZXhwb3J0cyA9IHsgZ2VuZXJhdGVSYW5kb21TdHJpbmcgfTsK';

    fs.writeFileSync(join(__dirname, 'utils.js'), Buffer.from(badLengthCode, 'base64').toString());

    // run test
    const { status: ____status } = spawnSync('npm', ['test'], { cwd: __dirname });

    if (____status === 0) {
      throw new Error('Pastikan Anda skenario ketika fungsi `generateRandomString` diberikan dengan argumen number yang kurang dari 1, ya!')
    }

    // restore the original code of utils.js
    fs.writeFileSync(join(__dirname, 'utils.js'), Buffer.from(originalCodeOfUtils, 'base64').toString());
    console.log('passed!');
  } catch (error) {
    console.log(error.message);
  } finally {
    try {
      fs.rmSync(join(__dirname, 'node_modules'), { recursive: true });
    } catch { /* do nothing */ }
  }
})()
