import http from 'k6/http';
import { check } from 'k6';

export let options = {
    vus: 20,       
    iterations: 20 
};

export default function () {
    const url = 'http://localhost:8080/api/v3/coupons/issue';

    const payload = JSON.stringify({
        couponPolicyCode: 'COUPON-300'
    });

    const params = {
        headers: {
            'Content-Type': 'application/json',
            'X-USER-ID': `USER-${__VU}`
        },
    };

    const res = http.post(url, payload, params);

    check(res, {
        'status is 400': (r) => r.status === 400
    });

    console.log(`VU ${__VU} Response status: ${res.status}`);
    console.log(`VU ${__VU} Response body:\n${res.body}`);
}
