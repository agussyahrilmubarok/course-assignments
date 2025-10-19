function generateRandomStringProm(length) {
  return new Promise((resolve, reject) => {
    if (typeof length !== 'number') {
      reject(Error('Length must be a number'));
    }

    if (length < 1) {
      reject(Error('Length must be greater than 0'));
    }

    const charset = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
    let result = '';

    for (let i = 0; i < length; i += 1) {
      const randomIndex = Math.floor(Math.random() * charset.length);
      result += charset.charAt(randomIndex);
    }

    resolve(result);
  }, 100);
}

module.exports = { generateRandomStringProm };
