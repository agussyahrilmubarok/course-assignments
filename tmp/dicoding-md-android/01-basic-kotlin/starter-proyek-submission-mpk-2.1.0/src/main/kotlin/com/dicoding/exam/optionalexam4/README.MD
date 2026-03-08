# Optional Exam 4 - Get Middle Characters

## Ketentuan

**Agar dapat diperiksa dengan baik, hindari beberapa hal berikut:**

1. Mengubah kode fungsi `main()` yang berada di dalam file `OptionalExam4Main.kt`
2. Menghapus atau mengubah nama fungsi dan variable yang sudah ada di dalam file `OptionalExam4.kt`
3. Membuat fungsi baru yang bukan merupakan tugas latihan
4. Mengubah struktur project (menghapus, mengedit, dan memindahkan package)

## Instruksi Latihan

Untuk mengerjakan soal latihan ini, silakan buka file `OptionalExam4.kt`.

### TODO

Sesuaikan fungsi `getMiddleCharacters()` dengan ketentuan sebagai berikut:

- Menerima input `String` dengan panjang karakter 1-100.
- Mengembalikan satu atau dua karakter paling tengah dari `String` yang diinput.
    - Jika panjang `String` bernilai **genap**, yang dikembalikan adalah `2` karakter.
    - Jika panjang `String` bernilai **ganjil**, yang dikembalikan adalah `1` karakter.

Modul terkait: Kotlin Fundamental & Collections

#### Contoh 1:

- Masukan: `string` = `dicoding`
- Keluaran: `od`
- Penjelasan:
    - Panjang teks `dicoding` adalah 8 karakter (genap). Karakter paling tengah adalah `od`.

#### Contoh 2:

- Masukan: `string` = `dicodingindonesia`
- Keluaran: `i`
- Penjelasan:
    - Panjang teks `dicodingindonesia` adalah 15 karakter (ganjil). Karakter paling tengah adalah `i`.

## Menguji Kode

Untuk menguji kode, silakan jalankan fungsi `main` pada file `OptionalExam4Main.kt`.

Jika kode pada `OptionalExam4.kt` sudah diperbaiki sesuai ketentuan, maka keluarannya akan seperti berikut:

```
true
true
true
```