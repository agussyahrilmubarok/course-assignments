/**
 * Hindari mengubah kode di bawah ini karena akan mempengaruhi proses penilaian.
 */

const { Readable, Writable } = require('stream');

function createReadable(data = []) {
  return Readable.from(data);
}

function createWritable() {
  const sink = [];
  const writable = new Writable({
    decodeStrings: false,
    write(chunk, encoding, callback) {
      sink.push(chunk);
      callback();
    },
  });

  writable.sink = sink;
  writable.on('pipe', () => { writable.piped = true; });

  return writable;
}

module.exports = {
  createReadable,
  createWritable,
};
