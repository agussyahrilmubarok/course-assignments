const { join } = require('path');
const fs = require('fs');

(() => {
  try {
    const answerPath = join(__dirname, 'answer.txt');
    const studentAnswer = fs.readFileSync(answerPath, 'utf-8').trim();

    const regex = /^node\s+--stack-trace-limit=\d+\s+app\.js$/;

    if (!regex.test(studentAnswer)) {
      throw new Error('Perintah untuk mengubah stack trace masih salah. Coba pelajari kembali: https://www.dicoding.com/academies/610/tutorials/33443');
    }

    console.log('passed!');
  } catch (error) {
    console.log(error.message);
  }
})();
