const { join } = require('path');
const fs = require('fs');

(() => {
  try {
    const packageJsonPath = join(__dirname, 'package.json');
    const userAnswer = fs.readFileSync(packageJsonPath, 'utf-8').trim();
    const packageJson = JSON.parse(userAnswer);

    const { dependencies, devDependencies, scripts } = packageJson;

    if (!dependencies) {
      throw new Error('Anda harus menambahkan `dependencies` pada berkas `package.json`');
    }

    if (!devDependencies) {
      throw new Error('Anda harus menambahkan `devDependencies` pada berkas `package.json`');
    }

    if (!dependencies.lodash) {
      throw new Error('Anda harus menambahkan `lodash` pada `dependenceis` di berkas `package.json`');
    }

    if (!devDependencies.nodemon) {
      throw new Error('Anda harus menambahkan `nodemon` pada `devDependencies` di berkas `package.json`');
    }

    if (!scripts.start) {
      throw new Error('Anda harus menambahkan `start` pada `scripts` di berkas `package.json`');
    }

    if (!scripts.dev) {
      throw new Error('Anda harus menambahkan `dev` pada `scripts` di berkas `package.json`');
    }

    // check lodash version should major version 3, minor any, patch any
    const lodashVersion = dependencies.lodash;
    const lodashVersionRegex = /^([\^~])?(\d+)(\.(\d+))?(\.(\d+))?$/;
    const lodashVersionMatch = lodashVersion.match(lodashVersionRegex);

    if(!lodashVersionMatch) {
      throw new Error('Versi dari `lodash` tidak sesuai dengan kriteria yang dibutuhkan.');
    }

    const lodashMajorVersion = lodashVersionMatch[2];
    if (lodashMajorVersion !== '3') {
      throw new Error('Versi dari `loddash` harus versi 3.');
    }

    // check start script
    const startScript = scripts.start;
    const startScriptRegex = /^node app.js$/;
    const startScriptMatch = startScript.match(startScriptRegex);

    if (!startScriptMatch) {
      throw new Error('Isi dari `start` harus `node app.js`.</p>');
    }

    // check dev script
    const devScript = scripts.dev;
    const devScriptRegex = /^nodemon app.js$/;
    const devScriptMatch = devScript.match(devScriptRegex);

    if (!devScriptMatch) {
      throw new Error('Isi dari `dev` harus `nodemon app.js`.</p>');
    }

    console.log('passed!');
  } catch (error) {
    console.log(error.message);
  }
})();
