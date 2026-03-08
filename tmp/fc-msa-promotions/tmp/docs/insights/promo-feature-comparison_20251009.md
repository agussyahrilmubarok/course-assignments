### 📌 **1. `coupon-service`:**

**Purpose:**
Manages **coupons** that users can apply for discounts or benefits when purchasing products or services.

**Key Features:**

* Issues discount coupons (e.g., 10% off, $5 off).
* Applies validation rules (expiry date, usage limit, minimum purchase).
* Tracks coupon usage per user or campaign.
* Supports different types of coupons: single-use, multi-use, public, or private.

**Example:**
A user enters a code like `SAVE10` at checkout and gets 10% off their total price.

---

### 📌 **2. `point-service`:**

**Purpose:**
Handles the **loyalty or reward points** system, allowing users to earn and redeem points based on their activities.

**Key Features:**

* Issues points based on transactions (e.g., 1 point per $1 spent).
* Allows users to redeem points for discounts or benefits.
* Manages point expiration, balance, and history.
* Supports promotional events that give bonus points.

**Example:**
A user earns 100 points after a purchase and uses 50 points for a discount on their next order.

---

### 📌 **3. `timesale-service`:**

**Purpose:**
Manages **limited-time sales**, offering discounted prices only during specific time windows.

**Key Features:**

* Schedules start and end times for sales events.
* Applies automatic discounts during the time window.
* Can be product-specific, category-based, or store-wide.
* May include countdowns or notifications to create urgency.

**Example:**
A product normally priced at $100 is automatically discounted to $70 between 3 PM–5 PM during a flash sale.

---

### 🧩 Summary Table

| Service            | Main Function                        | Example                         |
| ------------------ | ------------------------------------ | ------------------------------- |
| `coupon-service`   | Apply discounts via codes/vouchers   | Enter `SAVE10` for 10% off      |
| `point-service`    | Reward and redeem loyalty points     | Use 50 points = $5 discount     |
| `timesale-service` | Time-based automatic price discounts | Flash sale: 30% off for 2 hours |
