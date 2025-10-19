# Optional Exam 1 - Sum of Big Three

## Ketentuan

**Agar dapat diperiksa dengan baik, hindari beberapa hal berikut:**

1. Mengubah kode fungsi `main()` yang berada di dalam file `OptionalExam1Main.kt`
2. Menghapus atau mengubah nama fungsi dan variable yang sudah ada di dalam file `OptionalExam1.kt`
3. Membuat fungsi baru yang bukan merupakan tugas latihan
4. Mengubah struktur project (menghapus, mengedit, dan memindahkan package)

## Deskripsi Latihan

Untuk mengerjakan soal latihan ini, silakan buka file `OptionalExam1.kt`.

### TODO

Sesuaikan fungsi `sumOfBigThree()` dengan ketentuan sebagai berikut:

- Menerima 1 sampai 10 input bertipe integer tetapi hanya boleh ditampung dalam satu parameter, yaitu `numbers`.
- Mengembalikan nilai perjumlahan dari 3 angka yang paling besar dalam bentuk integer.

Modul terkait: Functional Programming dan Collections.

#### Contoh 1:

- Masukan: `numbers` = `1, 2, 3, 4, 5, 6, 7, 8, 9, 10`
- Keluaran: `27`
- Penjelasan:
    - Tiga angka terbesar dari masukan tersebut adalah `10, 9, 8`.
    - Hasil perjumlahan dari `10 + 9 + 8` adalah `27`

#### Contoh 2:

- Masukan: `numbers` = `22, 22, 33, 22`
- Keluaran: `77`
- Penjelasan:
    - Tiga angka terbesar dari masukan tersebut adalah `33, 22, 22`.
    - Hasil perjumlahan dari `33 + 22 + 22` adalah `77`

## Menguji Kode

Untuk menguji kode, silakan jalankan fungsi `main` pada file `OptionalExam1Main.kt`.

Jika kode pada `OptionalExam1.kt` sudah diperbaiki sesuai ketentuan, maka keluarannya akan seperti berikut:

```
true
true
true
```