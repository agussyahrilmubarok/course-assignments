# Exam 2 - Control Flow

## Ketentuan

**Agar dapat diperiksa dengan baik, hindari beberapa hal berikut:**

1. Mengubah kode fungsi `main()` yang berada di dalam file `Exam2Main.kt`
2. Menghapus atau mengubah nama fungsi dan variable yang sudah ada di dalam file `Exam2.kt`
3. Membuat fungsi baru yang bukan merupakan tugas latihan
4. Mengubah struktur project (menghapus, mengedit, dan memindahkan package)

## Instruksi

Untuk mengerjakan soal latihan ini, silakan buka file `Exam2.kt`.

### TODO 1

Sesuaikan fungsi `calculate` agar dapat menghasilkan nilai akhir dengan rumus:

```text
   calculate = valueA + (valueB - valueC)
```

Namun, jika `valueC` bernilai `null`, silakan tetapkan nilai `50` sebagai nilai _default_-nya.

#### Contoh 1:

- Masukan: `valueA` = `101`, `valueB` = `52`, `valueC` = `99`
- Keluaran: `54`
- Penjelasan: Hasil dari pertingungan `101 x (52 - 99)` adalah `54`.

#### Contoh 2:

- Masukan: `valueA` = `101`, `valueB` = `52`, `valueC` = `null`
- Keluaran: `103`
- Penjelasan:
    - Dikarenakan `valueC` bernilai `null`, maka digantikan oleh nilai _default_-nya yaitu `50`. Sehingga, hasil dari
      pertingungan `101 x (52 - 50)` adalah `103`.

### TODO 2

Sesuaikan fungsi `result` ini dapat mengembalikkan text seperti berikut:

```markdown
Result is $result
```

Di mana `$result` adalah argumen dari fungsi tersebut.

#### Contoh 1:

- Masukan: `result` = `54`
- Keluaran: `Result is 54`
- Penjelasan: String `Result is` ditambah dengan argumen result bernilai `54`.

#### Contoh 2:

- Masukan: `number` = `103`
- Keluaran: `Result is 103`
- Penjelasan: String `Result is` ditambah dengan argumen result bernilai `103`.

## Menguji Kode

Untuk menguji kode, silakan jalankan fungsi `main` pada file `Exam2Main.kt`.

Jika kode pada `Exam2.kt` sudah diperbaiki sesuai ketentuan, maka keluarannya akan seperti berikut:

```
Result is 54
Result is 103
```