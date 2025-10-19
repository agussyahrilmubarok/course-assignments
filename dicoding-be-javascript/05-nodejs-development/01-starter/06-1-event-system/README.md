# 06-1-event-system
Di berkas `app.js` terdapat objek `BasicMath` yang memiliki fungsi `add` dan `subtract`.
**Tugas kamu adalah**:
  1. Buatlah objek `eventEmitter` yang merupakan instance dari class `EventEmitter`.
  2. Tambahkan fungsi `add` dari `BasicMath` sebagai listener dari event `calculate`.
  3. Tambahkan fungsi `subtract` sebagai one-time listener dari event `calculate`.
  4. Bangkitkan event `calculate` dengan parameter angka 2 dan 3.
  5. Bangkitkan event `calculate` dengan parameter angka 5 dan 6.

**Catatan:**
- Ketika mengimpor event emitter, gunakan nama `EventEmitter` sebagai variabelnya.
- Hindari mengubah atau menghapus kode objek BasicMath karena penilaian akan gagal.
- Kamu boleh membuat objek event emitter dengan cara apa pun selama namanya `eventEmitter`.
- Jangan mencetak teks apa pun di console kecuali yang dihasilkan oleh listener.
