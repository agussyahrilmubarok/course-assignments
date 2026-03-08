const { join } = require('path');
const fs = require('fs');

(() => {
  try {
    const answerPath = join(__dirname, 'answer.txt');
    const studentAnswer = fs.readFileSync(answerPath, 'utf-8').trim();
    const regex = /^node\s+(-c|--check)\s+app\.js$/;

    if (!regex.test(studentAnswer)) {
      throw new Error('Perintah untuk mengecek sintaks JavaScript masih salah. Pelajari kembali modul: https://www.dicoding.com/academies/610/tutorials/33428');
    }

    console.log('passed!');
  } catch (error) {
    console.log(error.message);
  }
})();
