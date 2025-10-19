function generateRandomString(length, callback) {
  setTimeout(() => {
    if (typeof length !== 'number') {
      callback(Error('Length must be a number'));
      return;
    }

    if (length < 1) {
      callback(Error('Length must be greater than 0'));
      return;
    }

    const charset = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
    let result = '';

    for (let i = 0; i < length; i += 1) {
      const randomIndex = Math.floor(Math.random() * charset.length);
      result += charset.charAt(randomIndex);
    }

    callback(null, result);
  }, 100);
}

module.exports = { generateRandomString };
