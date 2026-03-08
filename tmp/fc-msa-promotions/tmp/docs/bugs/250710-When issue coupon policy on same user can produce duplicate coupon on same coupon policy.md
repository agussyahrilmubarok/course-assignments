> **"When issue coupon policy on same user can produce duplicate coupon on same coupon policy"**

### 🧨 **Problem Summary:**

Your system allows the **same user** to **receive multiple coupons** for the **same coupon policy**, which should not happen if each user is only supposed to receive **one coupon per policy**.

---

## 🛠 Root Cause:

You're missing a **check for existing coupon issuance** by the same user for a given policy **before issuing a new coupon**.

Even though you’re using Redis locks to control overall issuance (`decrementAndGetCouponPolicyQuantity`) **per policy**, there’s **no per-user validation**.

---

## ✅ Expected Behavior:

A user **must not receive more than one coupon** for a single coupon policy.

---

## 🧩 Where to Fix It:

You should add a **check inside `CouponIssuerServiceImpl.issueCoupon(...)`** to verify that the user has **not already received a coupon** for the same policy.

---

## ✅ Solution (Minimal Fix):

Add this inside `CouponIssuerServiceImpl.issueCoupon(...)` right **after acquiring the lock** and **before issuing the coupon**:

```java
String userId = UserIdInterceptor.getCurrentUserId();

// Check if user already has a coupon for this policy
boolean alreadyIssued = couponRepository.existsByUserIdAndCouponPolicyId(userId, request.getCouponPolicyId());
if (alreadyIssued) {
    log.warn("User {} already has a coupon for policy {}", userId, request.getCouponPolicyId());
    throw new CouponIssueException("You have already received this coupon.");
}
```

---

## 🧬 Required Repository Method:

In your `CouponRepository`, add the following method:

```java
boolean existsByUserIdAndCouponPolicyId(String userId, String couponPolicyId);
```

Spring Data JPA will auto-generate this query.

---

## 🧱 Optional: Improve with Redis Cache (if necessary)

If this query becomes a performance bottleneck under high load, you can cache **user-policy coupon issuance state** in Redis (e.g., `coupon:issued:{policyId}:{userId}`) and check/update it atomically.

But for now, the DB-level check is sufficient and straightforward.

---

## ✳️ Additional Suggestion:

Consider making the combination of `user_id + coupon_policy_id` **a unique constraint in the database** to enforce this at the persistence layer as well.

In `Coupon` entity:

```java
@Table(
    name = "coupons",
    uniqueConstraints = {
        @UniqueConstraint(columnNames = {"user_id", "coupon_policy_id"})
    }
)
```

This guarantees even under race conditions that duplicates can't be persisted.

---

## ✅ Summary of Fix Steps:

1. **Add check in `issueCoupon`** to see if user already has a coupon.
2. **Add `existsByUserIdAndCouponPolicyId`** to `CouponRepository`.
3. **(Optional)** Add DB constraint to enforce uniqueness.
