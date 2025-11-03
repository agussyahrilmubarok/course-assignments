import http from 'k6/http';
import { check } from 'k6';

/*
NOTE
- Spike Test
- Purpose: to test the system’s ability to handle sudden surges.
- vus: rises quickly → drops quickly
- indicator: does the system recover quickly after the spike?
- Goal: “What happens if suddenly 1000 users access the system at the same time?”
*/

export let options = {
  stages: [
    { duration: '5s', target: 0 },     // baseline
    { duration: '10s', target: 1000 }, // rapid spike
    { duration: '10s', target: 0 },    // drop
  ],
  thresholds: {
    http_req_duration: ['p(95)<1000'], // 95% requests < 1s
    http_req_failed: ['rate<0.01'],    // <1% failure
  },
};

export default function () {
    const url = 'http://localhost:8080/api/v2/coupons/issue';

    const payload = JSON.stringify({
        couponPolicyCode: 'COUPON-700'
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
