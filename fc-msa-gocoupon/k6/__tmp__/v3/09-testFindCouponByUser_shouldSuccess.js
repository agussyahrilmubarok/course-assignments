import http from 'k6/http';
import { check } from 'k6';

export let options = {
    vus: 1,       
    iterations: 1 
};

const userId = __ENV.userId || 'DEFAULT_CODE';

export default function () {
    const url = 'http://localhost:8080/api/v3/coupons/user';

    const params = {
        headers: {
            'Content-Type': 'application/json',
            'X-USER-ID': `${userId}`
        },
    };

    const res = http.get(url, params);

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
