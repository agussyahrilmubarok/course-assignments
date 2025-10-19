# Optional Exam 3 - Manipulate String

## Ketentuan

**Agar dapat diperiksa dengan baik, hindari beberapa hal berikut:**

1. Mengubah kode fungsi `main()` yang berada di dalam file `OptionalExam3Main.kt`
2. Menghapus atau mengubah nama fungsi dan variable yang sudah ada di dalam file `OptionalExam3.kt`
3. Membuat fungsi baru yang bukan merupakan tugas latihan
4. Mengubah struktur project (menghapus, mengedit, dan memindahkan package)

## Instruksi

Untuk mengerjakan soal latihan ini, silakan buka file `OptionalExam3.kt`.

### TODO

Sesuaikan fungsi `manipulateString()` dengan ketentuan sebagai berikut:

- Menerima input `String` dan `Integer` dengan rentang nilai `1` sampai `100`.
- Mengembalikan tipe data `String` dengan ketentuan:
    - Jika `String` **mengandung angka**, maka **kalikan** terlebih dahulu angka pada teks tersebut dengan angka
      yang **diinput**.

      Kemudian gabungkan teks sebelumnya dengan **hasil kali** tersebut.
    - Jika `String` **tidak mengandung angka**, maka langsung gabungkan teks dan angka yang diinput.

Modul terkait: Kotlin Fundamental & Collections

#### Contoh 1:

- Masukan: `str` = `Dicoding000010`, `int` = `100`
- Keluaran: `Dicoding1000`
- Penjelasan:
    - String mengandung angka, yaitu `000010`atau `10`.
    - Angka dari string dikalikan terlebih dahulu dengan angka yang diinput -> `10 x 100 = 1000`
    - Gabungkan teks sebelumnya dengan hasil kali ->  `Dicoding1000`

#### Contoh 2:

- Masukan: `str` = `dicoding`, `int` = `100`
- Keluaran: `dicoding100`
- Penjelasan:
    - String tidak mengandung angka.
    - Gabungkan langsung teks dan angka yang diinput ->  `dicoding100`

## Menguji Kode

Untuk menguji kode, silakan jalankan fungsi `main` pada file `OptionalExam3Main.kt`.

Jika kode pada `OptionalExam3.kt` sudah diperbaiki sesuai ketentuan, maka keluarannya akan seperti berikut:

```
true
true
true
```