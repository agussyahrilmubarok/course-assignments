# Example ECommerce Backend with Monolith Architecture


Tentu! Berikut adalah **daftar lengkap endpoint backend** untuk fitur utama dari sebuah **e-commerce system**, beserta kategorinya dan metode HTTP yang biasanya digunakan.

---

## 🛍️ E-Commerce Main Backend Endpoints (RESTful)

---

### 1. 📦 **Product Management**

| Endpoint            | Method   | Description                                                     |
| ------------------- | -------- | --------------------------------------------------------------- |
| `/api/products`     | `GET`    | Get all products (with optional filters: category, price, etc.) |
| `/api/products/:id` | `GET`    | Get details of a single product                                 |
| `/api/products`     | `POST`   | Create a new product *(admin only)*                             |
| `/api/products/:id` | `PUT`    | Update product info *(admin only)*                              |
| `/api/products/:id` | `DELETE` | Delete a product *(admin only)*                                 |

---

### 2. 🛒 **Cart Management**

| Endpoint            | Method   | Description                  |
| ------------------- | -------- | ---------------------------- |
| `/api/cart`         | `GET`    | Get current user’s cart      |
| `/api/cart`         | `POST`   | Add item to cart             |
| `/api/cart/:itemId` | `PUT`    | Update item quantity in cart |
| `/api/cart/:itemId` | `DELETE` | Remove item from cart        |

---

### 3. 🧾 **Order Management**

| Endpoint                 | Method | Description                                   |
| ------------------------ | ------ | --------------------------------------------- |
| `/api/orders`            | `GET`  | Get all user orders (or all orders for admin) |
| `/api/orders/:id`        | `GET`  | Get specific order details                    |
| `/api/orders`            | `POST` | Place a new order                             |
| `/api/orders/:id/cancel` | `PUT`  | Cancel an order (if allowed)                  |
| `/api/orders/:id/status` | `PUT`  | Update order status *(admin only)*            |

---

### 4. 👤 **User Management & Auth**

| Endpoint             | Method   | Description                        |
| -------------------- | -------- | ---------------------------------- |
| `/api/auth/register` | `POST`   | Register new user                  |
| `/api/auth/login`    | `POST`   | User login                         |
| `/api/auth/logout`   | `POST`   | User logout                        |
| `/api/users/me`      | `GET`    | Get current logged-in user profile |
| `/api/users/me`      | `PUT`    | Update user profile                |
| `/api/users`         | `GET`    | Get all users *(admin only)*       |
| `/api/users/:id`     | `DELETE` | Delete user *(admin only)*         |

---

### 5. 💳 **Payment Processing**

| Endpoint                | Method | Description                |
| ----------------------- | ------ | -------------------------- |
| `/api/payment/checkout` | `POST` | Initiate checkout/payment  |
| `/api/payment/verify`   | `POST` | Verify payment status      |
| `/api/payment/history`  | `GET`  | Get user’s payment history |

> 🔐 Often integrated with Stripe, Midtrans, Xendit, etc.

---

### 6. 🗂️ **Category & Tag Management**

| Endpoint              | Method   | Description                    |
| --------------------- | -------- | ------------------------------ |
| `/api/categories`     | `GET`    | Get all product categories     |
| `/api/categories/:id` | `GET`    | Get single category            |
| `/api/categories`     | `POST`   | Create category *(admin only)* |
| `/api/categories/:id` | `PUT`    | Update category *(admin only)* |
| `/api/categories/:id` | `DELETE` | Delete category *(admin only)* |

---

### 7. 📦 **Shipping & Address**

| Endpoint             | Method   | Description                |
| -------------------- | -------- | -------------------------- |
| `/api/addresses`     | `GET`    | Get user's saved addresses |
| `/api/addresses`     | `POST`   | Add new address            |
| `/api/addresses/:id` | `PUT`    | Update address             |
| `/api/addresses/:id` | `DELETE` | Remove address             |
| `/api/shipping/cost` | `POST`   | Get shipping cost estimate |

---

### 8. ⭐ **Reviews & Ratings**

| Endpoint                    | Method   | Description               |
| --------------------------- | -------- | ------------------------- |
| `/api/products/:id/reviews` | `GET`    | Get reviews for a product |
| `/api/products/:id/reviews` | `POST`   | Submit a review           |
| `/api/reviews/:reviewId`    | `PUT`    | Edit a review             |
| `/api/reviews/:reviewId`    | `DELETE` | Delete a review           |

---

### 9. 🔎 **Search & Filtering**

| Endpoint        | Method | Description                                                                |
| --------------- | ------ | -------------------------------------------------------------------------- |
| `/api/search`   | `GET`  | Global search for products                                                 |
| `/api/products` | `GET`  | Can include filters via query params (e.g. `?category=shoes&price_lt=100`) |

---

### 10. 🛠️ **Admin Dashboard APIs** *(optional but common)*

| Endpoint                   | Method | Description                             |
| -------------------------- | ------ | --------------------------------------- |
| `/api/admin/dashboard`     | `GET`  | Get summary: sales, users, orders, etc. |
| `/api/admin/top-products`  | `GET`  | Top-selling products                    |
| `/api/admin/reports/sales` | `GET`  | Sales reports by date range             |

---

## 🔐 Authentication & Authorization

* Use **JWT tokens**, **session-based auth**, or **OAuth2**.
* Roles: `user`, `admin`, maybe `vendor`.