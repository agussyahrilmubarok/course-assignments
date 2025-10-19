const { execSync, spawnSync } = require('child_process');
const { join } = require('path');

(() => {
  try {
    // checking syntax app.js
    execSync('node -c app.js', { cwd: __dirname });

    // run app.js
    const { status, stdout: stdoutAppBuffer } = spawnSync('node', ['app.js'], { cwd: __dirname });
    const stdoutApp = stdoutAppBuffer.toString().trim();

    // check platform
    if (stdoutApp !== process.platform) {
      throw new Error('Hasil yang dicetak pada Terminal tidak sesuai dengan platform yang digunakan. Pastikan juga Anda tidak mencetak teks lain di luar dari kriteria yang ditetapkan. Harap periksa kembali, ya! ^_^')
    }

    // check status code
    if (status === 0) {
      throw new Error('Proses tidak berhenti dengan status code non-zero (selain 0). Simak cara menghentikan proses dan menetapkan status code pada materi: https://www.dicoding.com/academies/610/tutorials/34053');
    }
    console.log('passed!');
  } catch (error) {
    console.log(error.message);
  }
})();
