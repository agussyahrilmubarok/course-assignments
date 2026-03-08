import http from 'k6/http';
import { Counter } from 'k6/metrics';

/*
NOTE
- Concurrency Test
- Purpose: To test how the system handles a very high number of simultaneous requests
           for a limited resource (coupon quota).
- Scenario: Attempt to send as many requests as possible in parallel for a coupon with limited quota.
- Goal: Observe race conditions, verify system correctly enforces coupon limits, 
        and measure success vs failure count.
*/

export let successCount = new Counter('coupon_success');
export let failCount = new Counter('coupon_failed');

export let options = {
    scenarios: {
        mass_requests: {
            executor: 'per-vu-iterations',
            vus: 100000,       // 100k virtual users (adjust based on your machine/resources)
            iterations: 10,    // each VU sends 10 requests → total ~1M requests
            maxDuration: '5m', // 5 minutes max
        },
    },
    thresholds: {
        http_req_duration: ['p(95)<1000'],
        http_req_failed: ['rate<0.01'],
    },
};

export default function () {
    const url = 'http://localhost:8080/api/v1/coupons/issue';

    const payload = JSON.stringify({
        couponPolicyCode: 'COUPON-1100',  // limited quota coupon
    });

    const USER_ID = `USER-${__VU}-${__ITER}-${Math.floor(Math.random() * 1000000)}`;

    const params = {
        headers: {
            'Content-Type': 'application/json',
            'X-USER-ID': USER_ID,
        },
    };

    const res = http.post(url, payload, params);

    if (res.status === 200) {
        successCount.add(1);
    } else {
        failCount.add(1);
    }
}
