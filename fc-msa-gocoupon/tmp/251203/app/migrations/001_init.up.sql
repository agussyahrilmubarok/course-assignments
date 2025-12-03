-- ==========================================
-- Types
-- ==========================================

-- CouponStatus enum
CREATE TYPE coupon_status AS ENUM (
    'AVAILABLE',
    'USED',
    'EXPIRED',
    'CANCELED'
);

-- DiscountType enum
CREATE TYPE discount_type AS ENUM (
    'FIXED_AMOUNT',
    'PERCENTAGE'
);

-- ==========================================
-- Tables
-- ==========================================

-- CouponPolicy table (parent)
CREATE TABLE coupon_policies (
    id TEXT PRIMARY KEY,
    code VARCHAR(50) UNIQUE NOT NULL,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    total_quantity INT NOT NULL,
    start_time TIMESTAMPTZ NOT NULL,
    end_time TIMESTAMPTZ NOT NULL,
    discount_type discount_type NOT NULL,
    discount_value INT NOT NULL,
    minimum_order_amount INT NOT NULL,
    maximum_discount_amount INT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Coupon table (child)
CREATE TABLE coupons (
    id TEXT PRIMARY KEY,
    code VARCHAR(50) UNIQUE NOT NULL,
    status coupon_status NOT NULL,
    used_at TIMESTAMPTZ,
    user_id TEXT NOT NULL,
    order_id TEXT,
    coupon_policy_id TEXT NOT NULL REFERENCES coupon_policies(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ==========================================
-- Indexes
-- ==========================================

CREATE INDEX idx_coupons_code ON coupons (code);
CREATE INDEX idx_coupons_status ON coupons (status);
CREATE INDEX idx_coupons_user_id ON coupons (user_id);
CREATE INDEX idx_coupons_coupon_policy_id ON coupons (coupon_policy_id);
