const fs = require('fs');
const crypto = require('crypto');

function _cleanUp() {
  fs.readdirSync(__dirname).forEach((file) => {
    if (fs.lstatSync(file).isDirectory()) {
      fs.rmdirSync(file);
    }
  });
}

function generateRandomFilesAndFolders() {
  _cleanUp();

  const folderCounts = Math.floor(Math.random() * 10) + 1;

  const folderNames = new Array(folderCounts)
    .fill(0)
    .map(() => crypto.randomBytes(8).toString('hex'));

  folderNames.forEach((folderName) => {
    fs.mkdirSync(folderName);
  });
}

generateRandomFilesAndFolders();
