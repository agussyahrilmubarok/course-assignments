import http from 'k6/http';
import { check } from 'k6';

export let options = {
    vus: 1,       
    iterations: 1 
};

export default function () {
    const url = 'http://localhost:8080/api/v2/couponPolicies/dummy';

    const res = http.post(url);

    check(res, {
        'status is 200': (r) => r.status === 200
    });

    const jsonData = JSON.parse(res.body);
    console.log('Response body (formatted JSON):\n', JSON.stringify(jsonData, null, 2));
}
