# Alert Notes

## Coupon Service


1. **Title:** Coupon Issue Failed
   **Query:** `count_over_time({service_name="coupon-service"} |= "Failed to issue new coupon" [10m])`
   **Description:** More than 3 coupon issuance failures occurred within the last 10 minutes.
   **Summary:** Over 3 coupon issuance failures in 10 minutes

2. **Title:** Coupon Invalid Period
   **Query:** `count_over_time({service_name="coupon-service"} |= "Coupon policy is not valid in the current period" [10m])`
   **Description:** Attempts to issue coupons outside of the valid policy period were detected in the last 10 minutes.
   **Summary:** Coupon issuance attempted outside valid policy period

3. **Title:** Coupon Policy Quota Exceeded
   **Query:** `count_over_time({service_name="coupon-service"} |= "Coupon policy quota exceeded" [10m])`
   **Description:** Attempts to issue coupons exceeding the policy quota occurred in the last 10 minutes.
   **Summary:** Coupon policy quota exceeded more than 3 times in 10 minutes

---
