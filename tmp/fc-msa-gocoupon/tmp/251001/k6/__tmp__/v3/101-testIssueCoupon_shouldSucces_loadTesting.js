import http from 'k6/http';
import { check } from 'k6';

/*
NOTE
- Load Test
- Purpose: to measure system performance under normal load up to slightly above normal.
- vus: corresponds to the average number of users (e.g., 100)
- duration: several minutes (5–30 minutes)
- indicators: latency, throughput, error rate
- Goal: “How many users can the system serve with a response time < 1 second?”
*/

export let options = {
  stages: [
    { duration: '1m', target: 50 },   // ramp-up to 50 user
    { duration: '5m', target: 100 },  // stable in 100 user
    { duration: '1m', target: 0 },    // ramp-down
  ],
  thresholds: {
    http_req_duration: ['p(95)<1000'], // 95% requests < 1 second
    http_req_failed: ['rate<0.01'],    // errors < 1%
  },
};

export default function () {
    const url = 'http://localhost:8080/api/v3/coupons/issue';

    const payload = JSON.stringify({
        couponPolicyCode: 'COUPON-500'
    });

    const USER_ID = `USER-${__VU}-${__ITER}-${Math.floor(Math.random()*1000)}`;

    const params = {
        headers: {
            'Content-Type': 'application/json',
            'X-USER-ID': USER_ID
        },
    };

    const res = http.post(url, payload, params);

    check(res, {
        'status is 200': (r) => r.status === 200,
        'response < 1s': (r) => r.timings.duration < 1000,
    });

    if (res.status !== 200) {
        console.log(`VU:${__VU} ITER:${__ITER} - Failed with status ${res.status}`);
    }
}
