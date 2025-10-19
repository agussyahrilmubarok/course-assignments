const { join } = require('path');
const { execSync } = require('child_process');
const fs = require('fs');

(() => {
  try {
    execSync('node -c app.js', { cwd: __dirname });

    const additionalTestingCode = 'CgovLyB0ZXN0aW5nIGNvZGUKY29uc3QgX19jcF9fID0gcmVxdWlyZSgnY2hpbGRfcHJvY2VzcycpOwpjb25zdCBfX2Fzc2VydF9fID0gcmVxdWlyZSgnYXNzZXJ0Jyk7Cgpjb25zdCB7IGVxdWFsIH0gPSBfX2Fzc2VydF9fLnN0cmljdDsKY29uc3QgeyBTQ0VOQVJJTyB9ID0gcHJvY2Vzcy5lbnY7CmNvbnN0IFtub2RlXSA9IHByb2Nlc3MuYXJndjsKCmNvbnN0IHN0ZG91dENoZWNrID0gKCkgPT4geyBtYWluKG5vZGUsIFsnLXAnLCAnXCd0ZXN0XCcnXSk7IH07CmNvbnN0IHN0ZGVyckNoZWNrID0gKCkgPT4gewogIGNvbnN0IHNwID0gbWFpbihub2RlLCBbJy1lJywgJ2NvbnNvbGUuZXJyb3IoXCd0ZXN0XCcpJ10pOwogIGlmIChzcC5zdGRlcnIpIHNwLnN0ZGVyci5waXBlKHByb2Nlc3Muc3RkZXJyKTsKfTsKCmNvbnN0IHN0ZGluQ2hlY2sgPSAoKSA9PiB7CiAgbWFpbihub2RlLCBbJy1lJywgYAogICAgICBwcm9jZXNzLnN0ZG91dC53cml0ZShCdWZmZXIuZnJvbShbMF0pKQogICAgICBwcm9jZXNzLnN0ZGluLnBpcGUocHJvY2Vzcy5zdGRvdXQpCiAgICAgIHNldFRpbWVvdXQoKCkgPT4gewogICAgICAgIHByb2Nlc3Muc3Rkb3V0LndyaXRlKEJ1ZmZlci5mcm9tKFsxXSkpCiAgICAgIH0sIDEwMCkKICBgXSk7Cn07CgovLyBlc2xpbnQtZGlzYWJsZS1uZXh0LWxpbmUgY29uc2lzdGVudC1yZXR1cm4KZnVuY3Rpb24gdGVzdChzY2VuYXJpbyA9IDApIHsKICAvLyBlc2xpbnQtZGlzYWJsZS1uZXh0LWxpbmUgZGVmYXVsdC1jYXNlCiAgc3dpdGNoIChzY2VuYXJpbykgewogICAgY2FzZSAxOiByZXR1cm4gc3Rkb3V0Q2hlY2soKTsKICAgIGNhc2UgMjogcmV0dXJuIHN0ZGVyckNoZWNrKCk7CiAgICBjYXNlIDM6IHJldHVybiBzdGRpbkNoZWNrKCk7CiAgfQoKICBjb25zdCBzMSA9IF9fY3BfXy5zcGF3blN5bmMobm9kZSwgW19fZmlsZW5hbWVdLCB7CiAgICBlbnY6IHsgU0NFTkFSSU86IDEgfSwKICB9KTsKCiAgZXF1YWwoczEuc3Rkb3V0LnRvU3RyaW5nKCkudHJpbSgpLCAndGVzdCcsICdDaGlsZCBwcm9jZXNzIGhhcnVzICJtZW5nZ3VuYWthbiBTVERPVVQgZGFyaSBwYXJlbnQgcHJvY2VzcyIuIFNpbWFrIGtlbWJhbGkgbWF0ZXJpOiAiaHR0cHM6Ly93d3cuZGljb2RpbmcuY29tL2FjYWRlbWllcy82MTAvdHV0b3JpYWxzLzM0NTU4IiBLb25maWd1cmFzaSBDaGlsZCBQcm9jZXNzIHVudHVrIG1lbmdldGFodWkgY2FyYSBtZW5nb25maWd1cmFzaSBTVERJTyBwYWRhIGNoaWxkIHByb2Nlc3MuJyk7CgogIGNvbnN0IHMyID0gX19jcF9fLnNwYXduU3luYyhub2RlLCBbX19maWxlbmFtZV0sIHsKICAgIGVudjogeyBTQ0VOQVJJTzogMiB9LAogIH0pOwoKICBlcXVhbChzMi5zdGRlcnIudG9TdHJpbmcoKS50cmltKCksICd0ZXN0JywgJ0NoaWxkIHByb2Nlc3MgaGFydXMgIm1lbmdla3Nwb3MgU1RERVJSLW55YSBzZW5kaXJpIi4gU2ltYWsga2VtYmFsaSBtYXRlcmk6ICJodHRwczovL3d3dy5kaWNvZGluZy5jb20vYWNhZGVtaWVzLzYxMC90dXRvcmlhbHMvMzQ1NTgiIEtvbmZpZ3VyYXNpIENoaWxkIFByb2Nlc3MgdW50dWsgbWVuZ2V0YWh1aSBjYXJhIG1lbmdvbmZpZ3VyYXNpIFNURElPIHBhZGEgY2hpbGQgcHJvY2Vzcy4nKTsKCiAgY29uc3QgczMgPSBfX2NwX18uc3Bhd25TeW5jKG5vZGUsIFtfX2ZpbGVuYW1lXSwgewogICAgaW5wdXQ6ICdzb21lIGlucHV0JywKICAgIGVudjogeyBTQ0VOQVJJTzogMyB9LAogIH0pOwoKICBlcXVhbChzMy5zdGRvdXQubGVuZ3RoLCAyLCAnQ2hpbGQgcHJvY2VzcyBoYXJ1cyAidGlkYWsgbWVtaWxpa2kga2VtYW1wdWFuIFNURElOIi4gU2ltYWsga2VtYmFsaSBtYXRlcmk6ICJodHRwczovL3d3dy5kaWNvZGluZy5jb20vYWNhZGVtaWVzLzYxMC90dXRvcmlhbHMvMzQ1NTgiIEtvbmZpZ3VyYXNpIENoaWxkIFByb2Nlc3MgdW50dWsgbWVuZ2V0YWh1aSBjYXJhIG1lbmdvbmZpZ3VyYXNpIFNURElPIHBhZGEgY2hpbGQgcHJvY2Vzcy4nKTsKfQoKdHJ5IHsKICB0ZXN0KE51bWJlcihTQ0VOQVJJTykpOwp9IGNhdGNoIChlcnJvcikgewogIGNvbnNvbGUuZXJyb3IoZXJyb3IubWVzc2FnZSk7CiAgcHJvY2Vzcy5leGl0KDEpOwp9';
    const appJsCode = fs.readFileSync(join(__dirname, 'app.js'))
      .toString();
    const appTestingCode = `${appJsCode}\n${Buffer.from(additionalTestingCode, 'base64')
      .toString()}`;

    fs.writeFileSync(join(__dirname, 'app.testing.js'), appTestingCode);

    try {
      execSync('node app.testing.js', {
        cwd: __dirname,
        timeout: 3000,
        stdio: ['ignore', 'ignore', 'pipe'],
      });
    } catch (error) {
      if (error.code === 'ETIMEDOUT') {
        throw new Error('Periksa kembali konfigurasi `STDIN` yang kamu tetapkan pada konfigurasi child process. Sesuai TODO yang dijelaskan, buatlah child process yang sama sekali tidak memiliki kemampuan untuk menangkap STDIN. Simak cara konfigurasi STDIO pada child process di materi https://www.dicoding.com/academies/610/tutorials/34558');
      }

      throw error;
    }

    console.log('passed!');
  } catch (error) {
    const { stderr } = error;
    console.log(stderr.toString());
  } finally {
    fs.unlinkSync(join(__dirname, 'app.testing.js'));
  }
})()
