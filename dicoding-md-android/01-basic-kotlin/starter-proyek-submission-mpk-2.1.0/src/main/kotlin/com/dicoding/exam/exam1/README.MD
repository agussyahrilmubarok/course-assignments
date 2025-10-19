# Exam 1 - Fundamental

## Ketentuan

**Agar dapat diperiksa dengan baik, hindari beberapa hal berikut:**

1. Mengubah kode fungsi `main()` yang berada di dalam file `Exam1Main.kt`
2. Menghapus atau mengubah nama fungsi dan variable yang sudah ada di dalam file `Exam1.kt`
3. Membuat fungsi baru yang bukan merupakan tugas latihan
4. Mengubah struktur project (menghapus, mengedit, dan memindahkan package)

## Instruksi

Untuk mengerjakan soal latihan ini, silakan buka file `Exam1.kt`.

### TODO 1

Sesuaikan fungsi `isEvenNumber` dengan kode untuk mengecek apakah parameter number merupakan **angka genap**.

```kotlin
    fun isEvenNumber(number: Int) = false
```

#### Contoh 1:

- Masukan: `number` = `77`
- Keluaran: `false`
- Penjelasan: Angka `77` merupakan angka ganjil.

#### Contoh 2:

- Masukan: `number` = `88`
- Keluaran: `true`
- Penjelasan: Angka `88` merupakan angka genap.

### TODO 2

Sesuaikan fungsi `moreThanFive` dengan kode untuk mengecek apakah parameter number **lebih dari 5**.

```kotlin
    fun moreThanFive(number: Int) = false
```

#### Contoh 1:

- Masukan: `number` = `5`
- Keluaran: `false`
- Penjelasan: Angka `5` tidak lebih dari angka `5`.

#### Contoh 2:

- Masukan: `number` = `88`
- Keluaran: `true`
- Penjelasan: Angka `88` lebih dari angka `5`.

### TODO 3

Sesuaikan fungsi `result` agar dapat menghasilkan nilai akhir dengan rumus:

```text
    result = number x (number + 10)
```

#### Contoh 1:

- Masukan: `number` = `1`
- Keluaran: `11`
- Penjelasan: 1 x (1 + 10) = 11

#### Contoh 2:

- Masukan: `number` = `10`
- Keluaran: `200`
- Penjelasan: 10 x (10 + 10) = 200

## Menguji Kode

Untuk menguji kode, silakan jalankan fungsi `main` pada file `Exam1Main.kt`.

Jika kode pada `Exam1.kt` sudah diperbaiki sesuai ketentuan, maka keluarannya akan seperti berikut:

```
false
false
11
true
true
200
```