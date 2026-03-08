Critical problem in distributed systems — **inconsistency between Redis (in-memory cache) and the database (persistent storage)**, especially in **concurrent environments**. Let's break this down and guide you toward a fix.

---

### Flow

- When coupon_policy store in db and coupon_policy quantity store in memory.
- Then will be producing coupon by issue coupon api, but coupon_policy quantity in memory and db doesn't sync.
- Then when use coupon, but coupon_policy quantity in memory and db doesn't sync.
- Then when cancel coupon, but coupon_policy quantity in memory and db doesn't sync.

---

### ❗ The Problem Recap

1. **`CouponPolicy` quantity is stored both in DB and Redis**, but changes are only applied in Redis during:

   - Coupon issuance (`decrement`)
   - Coupon cancellation (no Redis increment)
   - No proper **sync logic** or **source of truth**

2. **Inconsistencies arise** because:

   - Redis is updated, but the DB is not (or vice versa).
   - Cache TTL can expire.
   - No reconciliation mechanism exists.
   - Concurrency is only handled during **issuance** (via lock), not during **cancel**, **use**, or **sync**.

---

### 🔍 Root Causes

| Operation              | Redis Updated    | DB Updated                        | Lock Used         | Safe?            |
| ---------------------- | ---------------- | --------------------------------- | ----------------- | ---------------- |
| `issueCoupon()`        | ✅                | ❌ (DB `CouponPolicy` not updated) | ✅ (Redisson lock) | ⚠️ Partial       |
| `cancelCoupon()`       | ❌                | ✅                                 | ❌                 | ❌ Risk of leak   |
| `useCoupon()`          | ❌                | ✅                                 | ❌                 | ❌                |
| `createCouponPolicy()` | ✅ (quantity set) | ✅                                 | ❌                 | ⚠️ Risk of drift |
| Redis TTL expiry       | ❌                | ✅                                 | ❌                 | ❌ Unsynced again |

---

### ✅ Goals

1. Ensure **Redis quantity** is always **synced** with DB.
2. Maintain **correct remaining quantity** across all operations (issue, cancel).
3. Prevent **over-issuance**, even under high concurrency.
4. Avoid **leaked quotas** (e.g., not increasing quantity on cancel).

---

### ✅ Recommendations

---

## 🔧 [1] Define Single Source of Truth

- Decide: Is **Redis** the "real" source of truth for available quantity, or is it the **DB**?
- Typically:

  - Use **Redis** as primary for **quantity control** (for performance + atomicity).
  - Use **DB** as backup and persist only issued coupons.

✅ So Redis holds: `availableQuantity`, and DB is not used for real-time count.

---

## 🔧 [2] Sync quantity on cancel

Currently, `cancelCoupon()` updates DB but **not Redis** — so you **lose quota permanently** in Redis.

Add this logic inside `cancelCoupon()`:

```java
String quantityKey = "coupon:quantity:" + coupon.getCouponPolicy().getId();
RAtomicLong quantity = redissonClient.getAtomicLong(quantityKey);
quantity.incrementAndGet(); // Recover the quota
```

You **must protect this with a lock**, like in `issueCoupon()`:

```java
RLock lock = redissonClient.getLock("coupon:lock:" + policyId);
lock.lock();
try {
   quantity.incrementAndGet();
   // possibly update other info
} finally {
   lock.unlock();
}
```

---

## 🔧 [3] TTL Expiry Can Cause Data Loss

If Redis expires the quantity key (due to TTL), your system will **lose track of remaining coupons**.

Options:

- ❌ Don't expire the quantity key OR
- ✅ Persist the current value in DB as backup
- ✅ On startup or policy access, **re-populate Redis from DB**.

---

## 🔧 [4] Use Redis + Lua Script for Atomicity (Optional Enhancement)

Currently, you lock, then `get` and `decrement` in separate calls.

To avoid race conditions or improve performance, you can **use a Lua script** in Redis to do:

```lua
-- pseudo-Lua for atomic decrement
local key = KEYS[1]
local val = redis.call("GET", key)
if val and tonumber(val) > 0 then
    return redis.call("DECR", key)
else
    return -1
end
```

This avoids acquiring locks and ensures atomicity.

---

## 🔧 [5] Recovery Strategy or Reconciliation Task (Optional)

Run a **scheduled job** to:

- Compare DB `CouponPolicy.totalQuantity - issuedCount` with Redis quantity
- If mismatch, log error or resync

This helps in cases like:

- Redis crash or flush
- Unexpected bugs or service crashes

---

## 🔧 [6] Fix `create()` in `CouponPolicyServiceImpl`

You're calling:

```java
couponRedisService.setCouponPolicyQuantity(savedCouponPolicy);
return couponPolicyRepository.save(savedCouponPolicy);
```

You **call `save()` twice**, and second call is not needed. It may even override the first result.

Fix it:

```java
CouponPolicy saved = couponPolicyRepository.save(couponPolicy);
couponRedisService.setCouponPolicyQuantity(saved);
return saved;
```

---

## 🔧 [7] Log and Monitor Quantity

Add more logging when:

- Quantity becomes zero
- Quantity is incremented (cancel)
- Redis quantity is missing (null)

Use monitoring tools like Prometheus to alert if Redis and DB drift.

---

### ✅ Summary of Changes

| Operation     | Fix Required? | Fix                                                |
| ------------- | ------------- | -------------------------------------------------- |
| Create Policy | ✅             | Fix double save, ensure Redis quantity is correct  |
| Issue Coupon  | ⚠️            | Already uses Redis + Lock, optionally use Lua      |
| Use Coupon    | ⚠️            | No quantity update, may be okay depending on logic |
| Cancel Coupon | ✅             | Must `increment` Redis quantity and use lock       |
| Redis Expiry  | ✅             | Avoid TTL or rehydrate from DB                     |
| Resync        | ✅             | Add periodic check or health monitor               |

---

### ✅ Final Thoughts

Your architecture is solid, but missing just a few critical elements to make it **safe in distributed/high-concurrency environments**.