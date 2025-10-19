const { join } = require('path');
const { execSync } = require('child_process');
const fs = require('fs');

(() => {
  try {
    execSync('node -c app.js', { cwd: __dirname });

    // restore generate-folders.js code
    const originalCode = 'Y29uc3QgZnMgPSByZXF1aXJlKCdmcycpOwpjb25zdCBjcnlwdG8gPSByZXF1aXJlKCdjcnlwdG8nKTsKCmZ1bmN0aW9uIF9jbGVhblVwKCkgewogIGZzLnJlYWRkaXJTeW5jKF9fZGlybmFtZSkuZm9yRWFjaCgoZmlsZSkgPT4gewogICAgaWYgKGZzLmxzdGF0U3luYyhmaWxlKS5pc0RpcmVjdG9yeSgpKSB7CiAgICAgIGZzLnJtZGlyU3luYyhmaWxlKTsKICAgIH0KICB9KTsKfQoKZnVuY3Rpb24gZ2VuZXJhdGVSYW5kb21GaWxlc0FuZEZvbGRlcnMoKSB7CiAgX2NsZWFuVXAoKTsKCiAgY29uc3QgZm9sZGVyQ291bnRzID0gTWF0aC5mbG9vcihNYXRoLnJhbmRvbSgpICogMTApICsgMTsKCiAgY29uc3QgZm9sZGVyTmFtZXMgPSBuZXcgQXJyYXkoZm9sZGVyQ291bnRzKQogICAgLmZpbGwoMCkKICAgIC5tYXAoKCkgPT4gY3J5cHRvLnJhbmRvbUJ5dGVzKDgpLnRvU3RyaW5nKCdoZXgnKSk7CgogIGZvbGRlck5hbWVzLmZvckVhY2goKGZvbGRlck5hbWUpID0+IHsKICAgIGZzLm1rZGlyU3luYyhmb2xkZXJOYW1lKTsKICB9KTsKfQoKZ2VuZXJhdGVSYW5kb21GaWxlc0FuZEZvbGRlcnMoKTsK';
    fs.writeFileSync(join(
      __dirname,
      'generate-folders.js',
    ), Buffer.from(originalCode, 'base64').toString());

    execSync('node generate-folders.js', { cwd: __dirname });

    execSync('node app.js', { cwd: __dirname });

    try {
      const outData = fs.readFileSync(join(__dirname, 'out.txt'), { encoding: 'utf-8' }).trim();

      // read the folders inside folderPath
      const folders = fs.readdirSync(__dirname).filter((file) => fs.statSync(join(__dirname, file)).isDirectory()).join(',');

      // compare outData and folders
      if (outData !== folders) {
        throw new Error('Daftar folder pada berkas `out.txt` masih kurang tepat. Berikut beberapa hal yang dapat Anda pastikan: 1. Pastikan Anda hanya menulis nama berkas yang merupakan "direktori" saja. 2. Pastikan daftar directory yang ditulis disusun secara Ascending (A-Z).');
      }

      console.log('passed!');
    } catch(error) {
      throw error;
    }
  } catch (error) {
    console.log(error.message);
  }
})()
