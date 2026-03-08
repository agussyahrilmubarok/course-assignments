/**
 * @TODO
 * Ubah fungsi main agar mengembalikan child process (dari pemanggilan `spawn`) dengan ketentuan:
 *
 * 1. Menjalankan command dari arugmen `command`.
 * 2. Menetapkan command flag dari argument `args`.
 * 3. Tidak memiliki kemampuan STDIN.
 * 4. Menggunakan STDOUT parent process.
 * 5. Ekspos nilai STDERR-nya sendiri.
 */

const { spawn } = require('child_process');

function main(command, args) {
  return spawn();
}
