function generateRandomStringSync(length) {
  if (typeof length !== 'number') throw new Error('Length must be a number');
  if (length < 1) throw new Error('Length must be at least 1');

  const charset = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
  let result = '';

  for (let i = 0; i < length; i += 1) {
    const randomIndex = Math.floor(Math.random() * charset.length);
    result += charset.charAt(randomIndex);
  }

  return result;
}

module.exports = { generateRandomStringSync };
