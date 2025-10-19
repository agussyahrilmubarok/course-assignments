const { join } = require('path');
const { spawnSync, execSync } = require('child_process');
const fs = require('fs');

(() => {
  try {
    const originalCodeOfUtils = 'ZnVuY3Rpb24gZ2VuZXJhdGVSYW5kb21TdHJpbmdQcm9tKGxlbmd0aCkgewogIHJldHVybiBuZXcgUHJvbWlzZSgocmVzb2x2ZSwgcmVqZWN0KSA9PiB7CiAgICBpZiAodHlwZW9mIGxlbmd0aCAhPT0gJ251bWJlcicpIHsKICAgICAgcmVqZWN0KEVycm9yKCdMZW5ndGggbXVzdCBiZSBhIG51bWJlcicpKTsKICAgIH0KCiAgICBpZiAobGVuZ3RoIDwgMSkgewogICAgICByZWplY3QoRXJyb3IoJ0xlbmd0aCBtdXN0IGJlIGdyZWF0ZXIgdGhhbiAwJykpOwogICAgfQoKICAgIGNvbnN0IGNoYXJzZXQgPSAnQUJDREVGR0hJSktMTU5PUFFSU1RVVldYWVphYmNkZWZnaGlqa2xtbm9wcXJzdHV2d3h5ejAxMjM0NTY3ODknOwogICAgbGV0IHJlc3VsdCA9ICcnOwoKICAgIGZvciAobGV0IGkgPSAwOyBpIDwgbGVuZ3RoOyBpICs9IDEpIHsKICAgICAgY29uc3QgcmFuZG9tSW5kZXggPSBNYXRoLmZsb29yKE1hdGgucmFuZG9tKCkgKiBjaGFyc2V0Lmxlbmd0aCk7CiAgICAgIHJlc3VsdCArPSBjaGFyc2V0LmNoYXJBdChyYW5kb21JbmRleCk7CiAgICB9CgogICAgcmVzb2x2ZShyZXN1bHQpOwogIH0sIDEwMCk7Cn0KCm1vZHVsZS5leHBvcnRzID0geyBnZW5lcmF0ZVJhbmRvbVN0cmluZ1Byb20gfTsK';

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
    const badOutputCode = 'ZnVuY3Rpb24gZ2VuZXJhdGVSYW5kb21TdHJpbmdQcm9tKGxlbmd0aCkgewogIHJldHVybiBuZXcgUHJvbWlzZSgocmVzb2x2ZSwgcmVqZWN0KSA9PiB7CiAgICBpZiAodHlwZW9mIGxlbmd0aCAhPT0gJ251bWJlcicpIHsKICAgICAgcmVqZWN0KEVycm9yKCdMZW5ndGggbXVzdCBiZSBhIG51bWJlcicpKTsKICAgIH0KCiAgICBpZiAobGVuZ3RoIDwgMSkgewogICAgICByZWplY3QoRXJyb3IoJ0xlbmd0aCBtdXN0IGJlIGdyZWF0ZXIgdGhhbiAwJykpOwogICAgfQoKICAgIGNvbnN0IGNoYXJzZXQgPSAnQUJDREVGR0hJSktMTU5PUFFSU1RVVldYWVphYmNkZWZnaGlqa2xtbm9wcXJzdHV2d3h5ejAxMjM0NTY3ODknOwogICAgbGV0IHJlc3VsdCA9ICcnOwoKICAgIGZvciAobGV0IGkgPSAwOyBpIDwgbGVuZ3RoOyBpICs9IDEpIHsKICAgICAgY29uc3QgcmFuZG9tSW5kZXggPSBNYXRoLmZsb29yKE1hdGgucmFuZG9tKCkgKiBjaGFyc2V0Lmxlbmd0aCk7CiAgICAgIHJlc3VsdCArPSBjaGFyc2V0LmNoYXJBdChyYW5kb21JbmRleCk7CiAgICB9CgogICAgcmVzb2x2ZSgnJyk7CiAgfSwgMTAwKTsKfQoKbW9kdWxlLmV4cG9ydHMgPSB7IGdlbmVyYXRlUmFuZG9tU3RyaW5nUHJvbSB9Owo=';

    fs.writeFileSync(join(__dirname, 'utils.js'), Buffer.from(badOutputCode, 'base64').toString());

    // run test
    const { status: __status } = spawnSync('npm', ['test'], { cwd: __dirname });

    if (__status === 0) {
      throw new Error('Pastikan Anda menguji nilai yang dikembalikan oleh fungsi `generateRandomStringProm` dengan benar, ya!')
    }

    // test bad data type
    const badDataTypeCode = 'ZnVuY3Rpb24gZ2VuZXJhdGVSYW5kb21TdHJpbmdQcm9tKGxlbmd0aCkgewogIHJldHVybiBuZXcgUHJvbWlzZSgocmVzb2x2ZSwgcmVqZWN0KSA9PiB7CiAgICBpZiAodHlwZW9mIGxlbmd0aCAhPT0gJ251bWJlcicpIHsKICAgICAgcmVzb2x2ZSgnJyk7CiAgICB9CgogICAgaWYgKGxlbmd0aCA8IDEpIHsKICAgICAgcmVqZWN0KEVycm9yKCdMZW5ndGggbXVzdCBiZSBncmVhdGVyIHRoYW4gMCcpKTsKICAgIH0KCiAgICBjb25zdCBjaGFyc2V0ID0gJ0FCQ0RFRkdISUpLTE1OT1BRUlNUVVZXWFlaYWJjZGVmZ2hpamtsbW5vcHFyc3R1dnd4eXowMTIzNDU2Nzg5JzsKICAgIGxldCByZXN1bHQgPSAnJzsKCiAgICBmb3IgKGxldCBpID0gMDsgaSA8IGxlbmd0aDsgaSArPSAxKSB7CiAgICAgIGNvbnN0IHJhbmRvbUluZGV4ID0gTWF0aC5mbG9vcihNYXRoLnJhbmRvbSgpICogY2hhcnNldC5sZW5ndGgpOwogICAgICByZXN1bHQgKz0gY2hhcnNldC5jaGFyQXQocmFuZG9tSW5kZXgpOwogICAgfQoKICAgIHJlc29sdmUocmVzdWx0KTsKICB9LCAxMDApOwp9Cgptb2R1bGUuZXhwb3J0cyA9IHsgZ2VuZXJhdGVSYW5kb21TdHJpbmdQcm9tIH07Cg==';

    fs.writeFileSync(join(__dirname, 'utils.js'), Buffer.from(badDataTypeCode, 'base64').toString());

    // run test
    const { status: ___status } = spawnSync('npm', ['test'], { cwd: __dirname });

    if (___status === 0) {
      throw new Error('Pastikan Anda menguji skenario ketika fungsi `generateRandomStringProm` diberikan dengan argumen non-number, ya!')
    }

    // test bad length
    const badLengthCode = 'ZnVuY3Rpb24gZ2VuZXJhdGVSYW5kb21TdHJpbmdQcm9tKGxlbmd0aCkgewogIHJldHVybiBuZXcgUHJvbWlzZSgocmVzb2x2ZSwgcmVqZWN0KSA9PiB7CiAgICBpZiAodHlwZW9mIGxlbmd0aCAhPT0gJ251bWJlcicpIHsKICAgICAgcmVqZWN0KEVycm9yKCdMZW5ndGggbXVzdCBiZSBhIG51bWJlcicpKTsKICAgIH0KCiAgICBpZiAobGVuZ3RoIDwgMSkgewogICAgICByZXNvbHZlKCcnKTsKICAgIH0KCiAgICBjb25zdCBjaGFyc2V0ID0gJ0FCQ0RFRkdISUpLTE1OT1BRUlNUVVZXWFlaYWJjZGVmZ2hpamtsbW5vcHFyc3R1dnd4eXowMTIzNDU2Nzg5JzsKICAgIGxldCByZXN1bHQgPSAnJzsKCiAgICBmb3IgKGxldCBpID0gMDsgaSA8IGxlbmd0aDsgaSArPSAxKSB7CiAgICAgIGNvbnN0IHJhbmRvbUluZGV4ID0gTWF0aC5mbG9vcihNYXRoLnJhbmRvbSgpICogY2hhcnNldC5sZW5ndGgpOwogICAgICByZXN1bHQgKz0gY2hhcnNldC5jaGFyQXQocmFuZG9tSW5kZXgpOwogICAgfQoKICAgIHJlc29sdmUocmVzdWx0KTsKICB9LCAxMDApOwp9Cgptb2R1bGUuZXhwb3J0cyA9IHsgZ2VuZXJhdGVSYW5kb21TdHJpbmdQcm9tIH07Cg==';

    fs.writeFileSync(join(__dirname, 'utils.js'), Buffer.from(badLengthCode, 'base64').toString());

    // run test
    const { status: ____status } = spawnSync('npm', ['test'], { cwd: __dirname });

    if (____status === 0) {
      throw new Error('Pastikan Anda skenario ketika fungsi `generateRandomStringProm` diberikan dengan argumen number yang kurang dari 1, ya!')
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
