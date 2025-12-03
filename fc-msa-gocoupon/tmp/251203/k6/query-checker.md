# Coupon Database Queries

This document contains sample SQL queries for analyzing the `coupons` table.

## Table Structure

The `coupons` table has the following relevant columns:

- `id` (string, primary key): Unique coupon ID
- `coupon_policy_id` (string): ID of the coupon policy
- `user_id` (string): ID of the user who owns the coupon
- `status` (string): Coupon status (`AVAILABLE`, `USED`, `EXPIRED`, `CANCELED`)

---

## 1. Count coupons by `id` and `user_id`

This query counts how many records exist for a specific coupon and user combination:

```sql
SELECT COUNT(*) AS total
FROM coupons
WHERE id = 'COUPON_ID_HERE'
  AND user_id = 'USER_ID_HERE';
````

**Example:**

```sql
SELECT COUNT(*) AS total
FROM coupons
WHERE id = '123e4567-e89b-12d3-a456-426614174000'
  AND user_id = 'USER-001';
```

---

## 2. Count total coupons by `coupon_policy_id`

This query counts the number of coupons for a given coupon policy:

```sql
SELECT COUNT(*) AS total
FROM coupons
WHERE coupon_policy_id = 'COUPON_POLICY_ID_HERE';
```

**Example:**

```sql
SELECT COUNT(*) AS total
FROM coupons
WHERE coupon_policy_id = 'COUPON-900';
```

### Count coupons grouped by policy

```sql
SELECT coupon_policy_id, COUNT(*) AS total
FROM coupons
GROUP BY coupon_policy_id;
```

This query shows all coupon policies with the number of coupons each.

---

## Optional: Count coupons by status per policy

```sql
SELECT coupon_policy_id, status, COUNT(*) AS total
FROM coupons
GROUP BY coupon_policy_id, status;
```

This helps to see how many coupons are `AVAILABLE`, `USED`, or `EXPIRED` per policy.

```