import http from 'k6/http';
import { check } from 'k6';

/*
NOTE
- Soak Test (Endurance Test)
- Purpose: to test long-term system stability.
- vus: steady load (e.g., 100)
- duration: long (hours or days)
- indicators: memory leaks, connection leaks, response time degradation
- Goal: “Does the system remain stable after 6 hours?”
*/

export let options = {
  stages: [
    { duration: '1m', target: 100 },   // ramp-up ke 100 VU
    { duration: '6h', target: 100 },   // steady load 100 VU
    { duration: '1m', target: 0 },     // ramp-down
  ],
  thresholds: {
    http_req_duration: ['p(95)<1000'], 
    http_req_failed: ['rate<0.01'],    
  },
};

export default function () {
    const url = 'http://localhost:8080/api/v3/coupons/issue';

    const payload = JSON.stringify({
        couponPolicyCode: 'COUPON-800'
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
        console.error(`Error [VU ${__VU}] - Status: ${res.status}, Body: ${res.body}`);
    }

    sleep(1);
}
