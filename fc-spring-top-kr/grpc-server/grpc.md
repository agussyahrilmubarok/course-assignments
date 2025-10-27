## ⚙️ Basic Assumptions

We’ll assume:

* Your gRPC server is running on `localhost:9090`
* The proto package name is `bookstore`
* The service names are:

  * `bookstore.AuthorService`
  * `bookstore.BookService`
  * `bookstore.ReviewService`
* You’re using plaintext (no TLS). If you use TLS, remove `-plaintext`.

---

## 👤 1. **AUTHOR SERVICE**

### 🔹 1.1. Get All Authors

```bash
grpcurl -plaintext localhost:9090 bookstore.AuthorService/findAll
```

> Returns a stream of `Author` messages.

---

### 🔹 1.2. Get Author by ID

```bash
grpcurl -plaintext -d '{
  "id": 10001
}' localhost:9090 bookstore.AuthorService/get
```

---

### 🔹 1.3. Create Author

```bash
grpcurl -plaintext -d '{
  "name": "George Orwell"
}' localhost:9090 bookstore.AuthorService/create
```

---

### 🔹 1.4. Update Author

```bash
grpcurl -plaintext -d '{
  "id": 10001,
  "name": "George R. R. Martin"
}' localhost:9090 bookstore.AuthorService/update
```

---

### 🔹 1.5. Delete Author

```bash
grpcurl -plaintext -d '{
  "id": 10001
}' localhost:9090 bookstore.AuthorService/delete
```

---

## 📚 2. **BOOK SERVICE**

### 🔹 2.1. Get All Books

```bash
grpcurl -plaintext localhost:9090 bookstore.BookService/findAll
```

---

### 🔹 2.2. Get Book by ID (Detailed)

```bash
grpcurl -plaintext -d '{
  "id": 20001
}' localhost:9090 bookstore.BookService/get
```

---

### 🔹 2.3. Create Book

```bash
grpcurl -plaintext -d '{
  "title": "1984",
  "publisher": "Secker & Warburg",
  "published_date": {"seconds": 441763200}, 
  "author_ids": [10001, 10002]
}' localhost:9090 bookstore.BookService/create
```

> 🕓 `published_date.seconds` is a Unix epoch timestamp (example: `441763200` = 1984-01-01 UTC)

---

### 🔹 2.4. Update Book

```bash
grpcurl -plaintext -d '{
  "id": 20001,
  "title": "Animal Farm",
  "publisher": "Penguin Books",
  "published_date": {"seconds": 315532800},
  "author_ids": [10001]
}' localhost:9090 bookstore.BookService/update
```

---

### 🔹 2.5. Delete Book

```bash
grpcurl -plaintext -d '{
  "id": 20001
}' localhost:9090 bookstore.BookService/delete
```

---

## 📝 3. **REVIEW SERVICE**

### 🔹 3.1. Get All Reviews

```bash
grpcurl -plaintext localhost:9090 bookstore.ReviewService/findAll
```

---

### 🔹 3.2. Get Review by ID

```bash
grpcurl -plaintext -d '{
  "id": 30001
}' localhost:9090 bookstore.ReviewService/get
```

---

### 🔹 3.3. Create Review

```bash
grpcurl -plaintext -d '{
  "content": "A timeless dystopian masterpiece.",
  "rating": 4.8,
  "book_id": 20001
}' localhost:9090 bookstore.ReviewService/create
```

---

### 🔹 3.4. Update Review

```bash
grpcurl -plaintext -d '{
  "id": 30001,
  "content": "Still one of the best political novels ever written.",
  "rating": 5.0,
  "book_id": 20001
}' localhost:9090 bookstore.ReviewService/update
```

---

### 🔹 3.5. Delete Review

```bash
grpcurl -plaintext -d '{
  "id": 30001
}' localhost:9090 bookstore.ReviewService/delete
```

---

## 💡 Extra Tips

### 🔸 List all available services

```bash
grpcurl -plaintext localhost:9090 list
```

### 🔸 List all RPC methods in a service

```bash
grpcurl -plaintext localhost:9090 list bookstore.BookService
```

### 🔸 Describe a message or service schema

```bash
grpcurl -plaintext localhost:9090 describe bookstore.Book
grpcurl -plaintext localhost:9090 describe bookstore.BookService
```