const { join } = require('path');
const { spawnSync, execSync } = require('child_process');
const fs = require('fs');

(() => {
  try {
    const originalCodeOfUtils = 'ZnVuY3Rpb24gZ2VuZXJhdGVSYW5kb21TdHJpbmdTeW5jKGxlbmd0aCkgewogIGlmICh0eXBlb2YgbGVuZ3RoICE9PSAnbnVtYmVyJykgdGhyb3cgbmV3IEVycm9yKCdMZW5ndGggbXVzdCBiZSBhIG51bWJlcicpOwogIGlmIChsZW5ndGggPCAxKSB0aHJvdyBuZXcgRXJyb3IoJ0xlbmd0aCBtdXN0IGJlIGF0IGxlYXN0IDEnKTsKCiAgY29uc3QgY2hhcnNldCA9ICdBQkNERUZHSElKS0xNTk9QUVJTVFVWV1hZWmFiY2RlZmdoaWprbG1ub3BxcnN0dXZ3eHl6MDEyMzQ1Njc4OSc7CiAgbGV0IHJlc3VsdCA9ICcnOwoKICBmb3IgKGxldCBpID0gMDsgaSA8IGxlbmd0aDsgaSArPSAxKSB7CiAgICBjb25zdCByYW5kb21JbmRleCA9IE1hdGguZmxvb3IoTWF0aC5yYW5kb20oKSAqIGNoYXJzZXQubGVuZ3RoKTsKICAgIHJlc3VsdCArPSBjaGFyc2V0LmNoYXJBdChyYW5kb21JbmRleCk7CiAgfQoKICByZXR1cm4gcmVzdWx0Owp9Cgptb2R1bGUuZXhwb3J0cyA9IHsgZ2VuZXJhdGVSYW5kb21TdHJpbmdTeW5jIH07Cg==';

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
    const badOutputCode = 'ZnVuY3Rpb24gZ2VuZXJhdGVSYW5kb21TdHJpbmdTeW5jKGxlbmd0aCkgewogIGlmICh0eXBlb2YgbGVuZ3RoICE9PSAnbnVtYmVyJykgdGhyb3cgbmV3IEVycm9yKCdMZW5ndGggbXVzdCBiZSBhIG51bWJlcicpOwogIGlmIChsZW5ndGggPCAxKSB0aHJvdyBuZXcgRXJyb3IoJ0xlbmd0aCBtdXN0IGJlIGF0IGxlYXN0IDEnKTsKCiAgY29uc3QgY2hhcnNldCA9ICdBQkNERUZHSElKS0xNTk9QUVJTVFVWV1hZWmFiY2RlZmdoaWprbG1ub3BxcnN0dXZ3eHl6MDEyMzQ1Njc4OSc7CiAgbGV0IHJlc3VsdCA9ICcnOwoKICBmb3IgKGxldCBpID0gMDsgaSA8IGxlbmd0aDsgaSArPSAxKSB7CiAgICBjb25zdCByYW5kb21JbmRleCA9IE1hdGguZmxvb3IoTWF0aC5yYW5kb20oKSAqIGNoYXJzZXQubGVuZ3RoKTsKICAgIHJlc3VsdCArPSBjaGFyc2V0LmNoYXJBdChyYW5kb21JbmRleCk7CiAgfQoKICByZXR1cm4gJyc7Cn0KCm1vZHVsZS5leHBvcnRzID0geyBnZW5lcmF0ZVJhbmRvbVN0cmluZ1N5bmMgfTsK';

    fs.writeFileSync(join(__dirname, 'utils.js'), Buffer.from(badOutputCode, 'base64').toString());

    // run test
    const { status: __status } = spawnSync('npm', ['test'], { cwd: __dirname });

    if (__status === 0) {
      throw new Error('Pastikan Anda menguji nilai yang dikembalikan oleh fungsi `generateRandomStringSync` dengan benar, ya!')
    }

    // test bad data type
    const badDataTypeCode = 'ZnVuY3Rpb24gZ2VuZXJhdGVSYW5kb21TdHJpbmdTeW5jKGxlbmd0aCkgewogIGlmICh0eXBlb2YgbGVuZ3RoICE9PSAnbnVtYmVyJykgcmV0dXJuICcnOwogIGlmIChsZW5ndGggPCAxKSB0aHJvdyBuZXcgRXJyb3IoJ0xlbmd0aCBtdXN0IGJlIGF0IGxlYXN0IDEnKTsKCiAgY29uc3QgY2hhcnNldCA9ICdBQkNERUZHSElKS0xNTk9QUVJTVFVWV1hZWmFiY2RlZmdoaWprbG1ub3BxcnN0dXZ3eHl6MDEyMzQ1Njc4OSc7CiAgbGV0IHJlc3VsdCA9ICcnOwoKICBmb3IgKGxldCBpID0gMDsgaSA8IGxlbmd0aDsgaSArPSAxKSB7CiAgICBjb25zdCByYW5kb21JbmRleCA9IE1hdGguZmxvb3IoTWF0aC5yYW5kb20oKSAqIGNoYXJzZXQubGVuZ3RoKTsKICAgIHJlc3VsdCArPSBjaGFyc2V0LmNoYXJBdChyYW5kb21JbmRleCk7CiAgfQoKICByZXR1cm4gcmVzdWx0Owp9Cgptb2R1bGUuZXhwb3J0cyA9IHsgZ2VuZXJhdGVSYW5kb21TdHJpbmdTeW5jIH07Cg==';

    fs.writeFileSync(join(__dirname, 'utils.js'), Buffer.from(badDataTypeCode, 'base64').toString());

    // run test
    const { status: ___status } = spawnSync('npm', ['test'], { cwd: __dirname });

    if (___status === 0) {
      throw new Error('Pastikan Anda menguji skenario ketika fungsi `generateRandomStringSync` diberikan dengan argumen non-number, ya!')
    }

    // test bad length
    const badLengthCode = 'ZnVuY3Rpb24gZ2VuZXJhdGVSYW5kb21TdHJpbmdTeW5jKGxlbmd0aCkgewogIGlmICh0eXBlb2YgbGVuZ3RoICE9PSAnbnVtYmVyJykgdGhyb3cgbmV3IEVycm9yKCdMZW5ndGggbXVzdCBiZSBhIG51bWJlcicpOwogIGlmIChsZW5ndGggPCAxKSByZXR1cm4gJyc7CgogIGNvbnN0IGNoYXJzZXQgPSAnQUJDREVGR0hJSktMTU5PUFFSU1RVVldYWVphYmNkZWZnaGlqa2xtbm9wcXJzdHV2d3h5ejAxMjM0NTY3ODknOwogIGxldCByZXN1bHQgPSAnJzsKCiAgZm9yIChsZXQgaSA9IDA7IGkgPCBsZW5ndGg7IGkgKz0gMSkgewogICAgY29uc3QgcmFuZG9tSW5kZXggPSBNYXRoLmZsb29yKE1hdGgucmFuZG9tKCkgKiBjaGFyc2V0Lmxlbmd0aCk7CiAgICByZXN1bHQgKz0gY2hhcnNldC5jaGFyQXQocmFuZG9tSW5kZXgpOwogIH0KCiAgcmV0dXJuIHJlc3VsdDsKfQoKbW9kdWxlLmV4cG9ydHMgPSB7IGdlbmVyYXRlUmFuZG9tU3RyaW5nU3luYyB9Owo=';

    fs.writeFileSync(join(__dirname, 'utils.js'), Buffer.from(badLengthCode, 'base64').toString());

    // run test
    const { status: ____status } = spawnSync('npm', ['test'], { cwd: __dirname });

    if (____status === 0) {
      throw new Error('Pastikan Anda skenario ketika fungsi `generateRandomStringSync` diberikan dengan argumen number yang kurang dari 1, ya!')
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
