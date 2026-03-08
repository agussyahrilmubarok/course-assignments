# Exam 5 - Coroutines

## Ketentuan

**Agar dapat diperiksa dengan baik, hindari beberapa hal berikut:**

1. Mengubah kode fungsi `main()` yang berada di dalam file `Exam5Main.kt`
2. Menghapus atau mengubah nama fungsi dan variable yang sudah ada di dalam file `Exam5.kt`
3. Membuat fungsi baru yang bukan merupakan tugas latihan
4. Mengubah struktur project (menghapus, mengedit, dan memindahkan package)

## Instruksi

Untuk mengerjakan soal latihan ini, silakan buka file `Exam5.kt`.

### TODO 1

Sesuaikan fungsi `sum()` dengan kriteria seperti berikut:

- Jadikan sebagai suspending function dengan menambahkan `suspend` operator
- Menambahkan delay `3 detik`
- Mengembalikkan nilai kembalian hasil dari perhitungan denngan rumus `valueA + valueB`

#### Contoh:

- Masukan: `valueA` = `10`, `valueB` = `10`
- Keluaran: `20`
- Penjelasan: Hasil dari perhitungan `10 + 10` adalah `20`.

### TODO 2

Sesuaikan fungsi `multiple()` dengan kriteria seperti berikut:

- Jadikan sebagai suspending function dengan menambahkan `suspend` operator
- Menambahkan delay `2 detik`
- Mengembalikkan nilai kembalian hasil dari perhitungan denngan rumus `valueA x valueB`

#### Contoh:

- Masukan: `valueA` = `20`, `valueB` = `20`
- Keluaran: `400`
- Penjelasan: Hasil dari perhitungan `20 x 20` adalah `400`.

## Menguji Kode

Untuk menguji kode, silakan jalankan fungsi `main` pada file `Exam5Main.kt`.

Jika kode pada `Exam5.kt` sudah diperbaiki sesuai ketentuan, maka keluarannya akan seperti berikut:

```
Counting...
Result sum: 20
Result multiple: 400
```