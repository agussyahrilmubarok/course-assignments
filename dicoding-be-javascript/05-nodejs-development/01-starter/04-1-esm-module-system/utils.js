function unword(word) {
  return [...word];
}

function reword(...chars) {
  return chars.reduce((word, char) => word + char, '');
}

// @TODO ekspor function `unword` dan `reword`
