import http from 'k6/http';
import { check } from 'k6';
import { Counter, Trend } from 'k6/metrics';

const policyCode = __ENV.policyCode;

export const successCount = new Counter('success_count');
export const failCount = new Counter('fail_count');
export const requestTime = new Trend('request_time');

export const options = {
    vus: 1000,
    iterations: 1000,
};

export default function () {
  const url = 'http://localhost:8080/api/v2/coupons/issue';

  const userId = `USER_${Math.floor(Math.random() * 1000000)}`;

  const payload = JSON.stringify({
    policy_code: policyCode,
  });

  const params = {
    headers: {
      "Content-Type": "application/json",
      "X-USER-ID": userId,
    },
  };

  const res = http.post(url, payload, params);

  const ok = check(res, {
    "status 200": (r) => r.status === 200,
  });

  if (ok) {
    successCount.add(1);
  } else {
    failCount.add(1);
  }

  requestTime.add(res.timings.duration);
}
