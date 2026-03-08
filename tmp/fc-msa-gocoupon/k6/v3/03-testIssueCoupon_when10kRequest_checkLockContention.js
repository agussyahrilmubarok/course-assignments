import http from 'k6/http';
import { check } from 'k6';
import { Counter, Trend, Gauge } from 'k6/metrics';

const policyCode = __ENV.policyCode;

// CUSTOM METRICS
export const successCount = new Counter('success_count');
export const failCount = new Counter('fail_count');

export const requestTime = new Trend('request_time');

// metrics specifically to detect lock contention
export const dbWaitSuspect = new Counter('db_wait_suspect_count');  // >5s
export const dbWaitExtreme = new Counter('db_wait_extreme_count');  // >10s
export const maxLatency = new Gauge('max_latency');

export const options = {
    vus: 5000,
    iterations: 10000,
    thresholds: {
        // If p95 exceeds 5000ms → high probability of lock contention
        http_req_duration: ['p(95)<5000'],

        // If extreme lock count > 100 → confirmed contention
        db_wait_extreme_count: ['count<100'],
    }
};

export default function () {
    const url = 'http://localhost:8080/api/v3/coupons/issue';

    const userId = `USER_${Math.floor(Math.random() * 1000000)}`;

    const payload = JSON.stringify({
        policy_code: policyCode,
    });

    const params = {
        headers: {
            "Content-Type": "application/json",
            "X-USER-ID": userId,
        },
        timeout: "30s", // allow long wait for DB locks
    };

    const res = http.post(url, payload, params);

    const duration = res.timings.duration;
    requestTime.add(duration);
    maxLatency.add(duration);

    // detector: suspect lock wait (>5s)
    if (duration > 5000) {
        dbWaitSuspect.add(1);
    }

    // detector: extreme lock wait (>10s)
    if (duration > 10000) {
        dbWaitExtreme.add(1);

        // optional: print extremely slow requests
        console.error(`EXTREME WAIT: ${duration} ms | status=${res.status}`);
    }

    const ok = check(res, {
        "status 200": (r) => r.status === 200,
    });

    if (ok) {
        successCount.add(1);
    } else {
        failCount.add(1);
    }
}
