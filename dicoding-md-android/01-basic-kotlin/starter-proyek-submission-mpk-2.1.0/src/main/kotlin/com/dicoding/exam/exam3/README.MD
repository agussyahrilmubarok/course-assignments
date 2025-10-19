# Exam 3 - Generics

## Ketentuan

**Agar dapat diperiksa dengan baik, hindari beberapa hal berikut:**

1. Mengubah kode fungsi `main()` yang berada di dalam file `Exam3Main.kt`
2. Menghapus atau mengubah nama fungsi dan variable yang sudah ada di dalam file `Exam3.kt`
3. Membuat fungsi baru yang bukan merupakan tugas latihan
4. Mengubah struktur project (menghapus, mengedit, dan memindahkan package)

## Instruksi

Untuk mengerjakan soal latihan ini, silakan buka file `Exam3.kt`.

### TODO

Lengkapi fungsi `checkType` agar dapat mengembalikkan tipe nilai sesuai dengan parameter yang dilampirkan. Format
keluaran yang diharapkan adalah sebagai berikut:

```markdown
    Yes! it's $dataType
```

Berikut adalah beberapa tipe nilai yang wajib ada.

- Integer
- String
- Boolean
- Double
- List
- Map

#### Contoh 1:

- Masukan: `args` = `10`
- Keluaran: `Yes! it's Integer`
- Penjelasan: Tipe data dari `10` adalah `Integer`.

#### Contoh 2:

- Masukan: `args` = `listOf("A", "B")`
- Keluaran: `Yes! it's List`
- Penjelasan: Tipe data dari `listOf("A", "B")` adalah `List`.

#### Contoh 3:

- Masukan: `args` = `true`
- Keluaran: `Yes! it's Boolean`
- Penjelasan: Tipe data dari `true` adalah `Boolean`.

## Menguji Kode

Untuk menguji kode, silakan jalankan fungsi `main` pada file `Exam3Main.kt`.

Jika kode pada `Exam3.kt` sudah diperbaiki sesuai ketentuan, maka keluarannya akan seperti berikut:

```
'Dicoding Indonesia' is String? Yes! it's String
'True' is Boolean? Yes! it's Boolean
'1' is Integer? Yes! it's Integer
'10.01' is Double? Yes! it's Double
'[10, 9, 8 , 6]' is List? Yes! it's List
'{score=10}' is Map? Yes! it's Map
```