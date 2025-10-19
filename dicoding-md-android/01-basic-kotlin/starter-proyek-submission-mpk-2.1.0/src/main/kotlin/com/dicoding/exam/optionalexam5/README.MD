# Optional Exam 5 - Concat String

## Ketentuan

**Agar dapat diperiksa dengan baik, hindari beberapa hal berikut:**

1. Mengubah kode fungsi `main()` yang berada di dalam file `OptionalExam5Main.kt`
2. Menghapus atau mengubah nama fungsi dan variable yang sudah ada di dalam file `OptionalExam5.kt`
3. Membuat fungsi baru yang bukan merupakan tugas latihan
4. Mengubah struktur project (menghapus, mengedit, dan memindahkan package)

## Instruksi

Untuk mengerjakan soal latihan ini, silakan buka file `OptionalExam5.kt`.

### TODO

Sesuaikan fungsi `concatString()` agar dapat menggabungkan dua `String` dengan syarat sebagai berikut:

- Tidak boleh menggunakan keyword `fun`.
- Tidak boleh mengandung karakter `{` dan `}`.

Modul terkait: Kotlin Fundamental & Functional Programming

_Hint_: Kamu bisa menggunakan _Lambda_ dan Member _References_

#### Contoh 1:

- Masukan: `string1` = `Hello`, `string2` = `Dicoding`
- Keluaran: `HelloDicoding`
- Penjelasan:
    - String `Hello` digabungkan dengan `Dicoding` menjadi `HelloDicoding`.

#### Contoh 2:

- Masukan: `string1` = `Dicoding `, `string2` = `Indonesia`
- Keluaran: `Dicoding Indonesia`
- Penjelasan:
    - String `Dicoding ` digabungkan dengan `Indonesia` menjadi `Dicoding Indonesia`.

## Menguji Kode

Untuk menguji kode, silakan jalankan fungsi `main` pada file `OptionalExam5Main.kt`.

Jika kode pada `OptionalExam5.kt` sudah diperbaiki sesuai ketentuan, maka keluarannya akan seperti berikut:

```
true
true
true
```