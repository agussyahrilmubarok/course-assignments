import http from 'k6/http';
import { check } from 'k6';

/*
NOTE
- Stress Test
- Purpose: to find the system’s maximum limit before failure.
- vus: continuously increases until errors start appearing
- indicator: the point where the error rate exceeds 5% or latency spikes sharply
- Goal: “At what number of users does the system start to fail?”
*/

export let options = {
  stages: [
    { duration: '5s', target: 0 },     // idle baseline
    { duration: '20s', target: 500 },  // ramp up gradually
    { duration: '20s', target: 1000 }, // peak load
    { duration: '10s', target: 0 },    // ramp down
  ],
  thresholds: {
    http_req_duration: ['p(95)<1000'], // 95% requests < 1s
    http_req_failed: ['rate<0.05'],    // <5% error tolerated in stress test
  },
};

export default function () {
    const url = 'http://localhost:8080/api/v2/coupons/issue';

    const payload = JSON.stringify({
        couponPolicyCode: 'COUPON-600'
    });

    // USER_ID unik: VU + iterasi + random small offset
    const USER_ID = `USER-${__VU}-${__ITER}-${Math.floor(Math.random() * 1000)}`;

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
