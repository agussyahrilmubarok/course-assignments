import http from 'k6/http';
import { check } from 'k6';

export let options = {
    vus: 1,       
    iterations: 1 
};

export default function () {
    const url = 'http://localhost:8080/api/v4/coupons/issue';

    const payload = JSON.stringify({
        couponPolicyCode: 'COUPON-100'
    });

    const params = {
        headers: {
            'Content-Type': 'application/json',
            'X-USER-ID': 'USER-123'
        },
    };

    const res = http.post(url, payload, params);

    check(res, {
        'status is 200': (r) => r.status === 200
    });

    try {
        const jsonData = JSON.parse(res.body);
        console.log('Response body (formatted JSON):\n', JSON.stringify(jsonData, null, 2));
    } catch (e) {
        console.log('Response is not valid JSON:', res.body);
    }
}
