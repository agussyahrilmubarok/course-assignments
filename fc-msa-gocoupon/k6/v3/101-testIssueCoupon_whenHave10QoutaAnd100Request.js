import http from 'k6/http';
import { check, sleep } from 'k6';
import { Counter } from 'k6/metrics';

export let successCount = new Counter('coupon_success');
export let failCount = new Counter('coupon_failed');

export let options = {
    vus: 100,
    duration: '1s',
    thresholds: {
        'coupon_success': ['count<=10'],
        'coupon_failed': ['count>=90'],
    },
};

export default function () {
    const url = 'http://localhost:8080/api/v3/coupons/issue';

    const payload = JSON.stringify({
        couponPolicyCode: 'COUPON-900'  // only 10 quotas available
    });

    const USER_ID = `USER-${__VU}-${__ITER}-${Math.floor(Math.random() * 1000)}`;

    const params = {
        headers: {
            'Content-Type': 'application/json',
            'X-USER-ID': USER_ID,
        },
    };

    const res = http.post(url, payload, params);

    if (res.status === 200) {
        successCount.add(1);
        console.log(`SUCCESS [VU ${__VU}] - USER_ID: ${USER_ID}`);
    } else {
        failCount.add(1);
        console.error(`FAILED [VU ${__VU}] - USER_ID: ${USER_ID}, Status: ${res.status}, Body: ${res.body}`);
    }

    sleep(1);
}
