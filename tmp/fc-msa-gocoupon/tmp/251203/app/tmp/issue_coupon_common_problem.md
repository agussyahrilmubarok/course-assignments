# Coupon Service - IssueCoupon Feature

## Overview

The **IssueCoupon** feature allows the system to create a new coupon for a user based on an existing **CouponPolicy**. This feature enforces business rules such as validity period, quantity limits, and user eligibility while handling potential technical failures in the database or service layer.

---

## Table of Contents

* [Feature Flow](#feature-flow)
* [Error Handling](#error-handling)

  * [Domain / Business Errors](#domain--business-errors)
  * [Technical / System Errors](#technical--system-errors)
* [Database Considerations](#database-considerations)
* [Usage Example](#usage-example)

---

## Feature Flow

The **IssueCoupon** workflow can be summarized in the following steps:

1. **Retrieve Coupon Policy**

   * Look up the policy by code in the database.
2. **Check Valid Period**

   * Verify that the current time is within the policy’s start and end dates.
3. **Check Available Quantity**

   * Ensure the number of issued coupons does not exceed `TotalQuantity`.
4. **Check User Eligibility**

   * Verify that the user has not already claimed this coupon.
5. **Check Order / Product Requirements (optional)**

   * Validate any business rules like minimum order amount or applicable products.
6. **Create Coupon**

   * Insert a new coupon into the database for the user.
7. **Return Coupon**

   * Return the successfully issued coupon object.

---

## Error Handling

### Domain / Business Errors

These errors are caused by business rule violations:

| Step                  | Error                           | Description                                                   |
| --------------------- | ------------------------------- | ------------------------------------------------------------- |
| Retrieve Policy       | `ErrCouponPolicyNotFound`       | The policy code does not exist.                               |
| Check Valid Period    | `ErrCouponPolicyNotActive`      | The policy has not started yet.                               |
| Check Valid Period    | `ErrCouponPolicyExpired`        | The policy has already expired.                               |
| Check Quantity        | `ErrCouponPolicyQuantityExceed` | The maximum number of coupons for the policy has been issued. |
| Check User Claim      | `ErrCouponUserAlreadyClaimed`   | The user has already claimed this coupon.                     |
| Check Order / Product | `ErrCouponOrderAmountTooLow`    | The user’s order does not meet the minimum amount.            |
| Check Order / Product | `ErrCouponInvalidForProduct`    | The coupon is not valid for the selected product.             |

### Technical / System Errors

These errors are caused by infrastructure, database, or service failures:

| Step                           | Error                            | Description                                     |
| ------------------------------ | -------------------------------- | ----------------------------------------------- |
| Retrieve Policy                | `ErrCouponInternal`              | Unexpected failure in the service layer.        |
| Retrieve Policy                | `ErrTimeout`                     | Database query timed out.                       |
| Retrieve Policy                | `ErrDatabaseUnavailable`         | Database is down or unreachable.                |
| Check Quantity                 | `ErrCouponCounted`               | Failed to count issued coupons due to DB error. |
| Check Quantity / Create Coupon | `ErrCouponQuantityRaceCondition` | Concurrent inserts may exceed total quantity.   |
| Create Coupon                  | `ErrCouponCreated`               | Failed to insert the coupon into the database.  |
| Create Coupon                  | `ErrTransactionFailed`           | Transaction commit or rollback failed.          |
| Return Coupon                  | `ErrUnknown`                     | Unexpected runtime or internal error.           |

---

## Database Considerations

* **CouponPolicy** table stores the coupon rules, start/end times, and maximum quantity.
* **Coupons** table stores issued coupons linked to the policy and the user.
* All timestamps should use `TIMESTAMPTZ` and be stored in UTC to prevent time zone issues.
* Ensure proper indexes on `coupon_code`, `user_id`, and `coupon_policy_id` for performance.
* Consider database transactions when issuing coupons to avoid race conditions.

---

## Usage Example

### Service Layer

```go
coupon, err := service.IssueCoupon(ctx, "WELCOME50", userID)
if err != nil {
    if errors.Is(err, coupon.ErrCouponPolicyNotFound) {
        return echo.NewHTTPError(http.StatusNotFound, err.Error())
    } else if errors.Is(err, coupon.ErrCouponPolicyQuantityExceed) {
        return echo.NewHTTPError(http.StatusBadRequest, err.Error())
    }
    return echo.NewHTTPError(http.StatusInternalServerError, coupon.ErrCouponInternal)
}
```

### HTTP Request Example

```bash
curl -X POST http://localhost:8080/api/v1/coupons/issue \
  -H "Content-Type: application/json" \
  -H "X-USER-ID: 12345" \
  -d '{"policy_code": "WELCOME50"}'
```

**Response:**

```json
{
  "id": "uuid-generated",
  "code": "uuid-generated",
  "status": "AVAILABLE",
  "user_id": "12345",
  "coupon_policy_id": "policy-uuid",
  "created_at": "2025-11-21T10:00:00Z",
  "updated_at": "2025-11-21T10:00:00Z"
}
```

---

## Notes

* All domain/business errors typically return **HTTP 400-404** status codes.
* Technical errors typically return **HTTP 500**.
* Always use `UTC` for timestamps to avoid time zone inconsistencies.
* Consider handling **race conditions** when checking quantity and creating coupons.

---