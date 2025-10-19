/**
 * @TODO
 * 1. Buatlah sebuah readable stream dengan nama `readable` menggunakan fungsi `createReadable`
 *    yang telah disediakan. Gunakan array `alphabet` sebagai argumen pemanggilan fungsi tersebut.
 * 2. Buatlah sebuah writable stream dengan nama `writable` menggunakan fungsi `createWritable`.
 * 3. Pipe `readable` ke `writable`.
 *
 * Catatan:
 * - Jangan ubah kode yang sudah disediakan sebelumnya.
 * - Anda bisa melihat kode yang sudah disediakan di folder `utils.js` untuk
 *   melihat implementasi dari fungsi `createReadable` dan `createWritable`.
 */

const { createReadable, createWritable } = require('./utils');

const alphabet = ['a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j'];

