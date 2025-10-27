## 🧠 GraphQL Queries and Mutations


Run the Spring Boot app:

```bash
mvn spring-boot:run
```

Default GraphQL endpoint:

```
http://localhost:8080/graphql
```

GraphiQL (visual playground) UI:

```
http://localhost:8080/graphiql
```

### 🔹 AUTHOR

#### ➤ Find All Authors

```graphql
query {
  findAllAuthors {
    id
    name
    dateCreated
  }
}
```

#### ➤ Get Author by ID

```graphql
query {
  getAuthor(id: 10000) {
    id
    name
    dateCreated
  }
}
```

#### ➤ Create Author

```graphql
mutation {
  createAuthor(input: { name: "J.K. Rowling" }) {
    id
    name
  }
}
```

#### ➤ Update Author

```graphql
mutation {
  updateAuthor(id: 1000, input: { name: "Joanne Rowling" }) {
    id
    name
  }
}
```

#### ➤ Delete Author

```graphql
mutation {
  deleteAuthor(id: 1000)
}
```

Expected response:

```json
{
  "data": {
    "deleteAuthor": true
  }
}
```

---

### 🔹 BOOK

#### ➤ Find All Books

```graphql
query {
  findAllBooks {
    id
    title
    publisher
    authors {
      id
      name
    }
  }
}
```

#### ➤ Get Book by ID

```graphql
query {
  getBook(id: 10000) {
    id
    title
    publisher
    publishedDate
    authors {
      name
    }
    reviews {
      content
      rating
    }
  }
}
```

#### ➤ Create Book

```graphql
mutation {
  createBook(input: {
    title: "Harry Potter and the Goblet of Fire"
    publisher: "Bloomsbury"
    publishedDate: "2000-07-08T00:00:00Z"
    authorIds: [10000]
  }) {
    id
    title
    publisher
  }
}
```

#### ➤ Update Book

```graphql
mutation {
  updateBook(id: 10000, input: {
    title: "Harry Potter - Updated"
    publisher: "Scholastic"
    publishedDate: "2000-07-08T00:00:00Z"
  }) {
    id
    title
    publisher
  }
}
```

#### ➤ Delete Book

```graphql
mutation {
  deleteBook(id: 10000)
}
```

Expected response:

```json
{
  "data": {
    "deleteBook": true
  }
}
```

---

### 🔹 REVIEW

#### ➤ Find All Reviews

```graphql
query {
  findAllReviews {
    id
    content
    rating
    bookId
  }
}
```

#### ➤ Get Review by ID

```graphql
query {
  getReview(id: 10000) {
    id
    content
    rating
    bookId
  }
}
```

#### ➤ Create Review

```graphql
mutation {
  createReview(input: {
    content: "Fantastic read!"
    rating: 4.8
    bookId: 10000
  }) {
    id
    content
    rating
  }
}
```

#### ➤ Update Review

```graphql
mutation {
  updateReview(id: 1, input: {
    content: "Even better on the second read!"
    rating: 5.0
    bookId: 1
  }) {
    id
    content
    rating
  }
}
```

#### ➤ Delete Review

```graphql
mutation {
  deleteReview(id: 1)
}
```

Expected response:

```json
{
  "data": {
    "deleteReview": true
  }
}
```

---

## ⚙️ Testing with cURL

You can also test via terminal:

```bash
curl -X POST http://localhost:8080/graphql \
  -H "Content-Type: application/json" \
  -d '{"query":"{ findAllAuthors { id name } }"}'
```

Or a mutation:

```bash
curl -X POST http://localhost:8080/graphql \
  -H "Content-Type: application/json" \
  -d '{"query":"mutation { createAuthor(input: {name: \"New Author\"}) { id name } }"}'
```

---

## 🧪 Automated Testing (Optional)

You can use `GraphQlTester` for automated integration tests.

```java
@SpringBootTest
class GraphQlApiTests {

    @Autowired
    private GraphQlTester graphQlTester;

    @Test
    void testFindAllAuthors() {
        graphQlTester.document("""
            query {
                findAllAuthors {
                    id
                    name
                }
            }
        """).execute()
          .path("findAllAuthors")
          .hasValue();
    }
}
```

---

## ✅ Summary

| Entity     | Queries                       | Mutations                                      | Return Type          |
| ---------- | ----------------------------- | ---------------------------------------------- | -------------------- |
| **Author** | `findAllAuthors`, `getAuthor` | `createAuthor`, `updateAuthor`, `deleteAuthor` | `Author` / `Boolean` |
| **Book**   | `findAllBooks`, `getBook`     | `createBook`, `updateBook`, `deleteBook`       | `Book` / `Boolean`   |
| **Review** | `findAllReviews`, `getReview` | `createReview`, `updateReview`, `deleteReview` | `Review` / `Boolean` |

---