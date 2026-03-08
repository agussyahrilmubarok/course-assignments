import http from "k6/http";
import { check, sleep } from "k6";
import { Counter, Trend } from "k6/metrics";

// Custom Metrics
export const totalReq = new Counter("total_requests");
export const ok200 = new Counter("status_200");
export const tooMany429 = new Counter("status_429");

// Optional performance trends
export const ok200Trend = new Trend("duration_200");
export const tooManyTrend = new Trend("duration_429");

export let options = {
  scenarios: {
    rate_limit_test: {
      executor: "constant-arrival-rate",
      rate: 30,                 // 30 RPS
      timeUnit: "1s",
      duration: "10s",
      preAllocatedVUs: 50,
      maxVUs: 200,
    },
  },

  thresholds: {
    http_req_failed: ["rate<0.99"], // allow failures because rate limit expected
  },
};

export function setup() {
  const email = `user_${Date.now()}@example.com`;
  const password = "Password123!";

  const res = http.post(
    "http://localhost:8081/api/v1/accounts/sign-up",
    JSON.stringify({
      name: "Rate Limit Test",
      email,
      password,
    }),
    { headers: { "Content-Type": "application/json" } }
  );

  check(res, {
    "sign-up OK/Conflict": (r) => [200, 201, 409].includes(r.status),
  });

  return { email, password };
}

// Test target: SIGN-IN endpoint with rate limiter
export default function (data) {
  const payload = JSON.stringify({
    email: data.email,
    password: data.password,
  });

  const res = http.post(
    "http://localhost:8081/api/v1/accounts/sign-in",
    payload,
    { headers: { "Content-Type": "application/json" } }
  );

  totalReq.add(1);

  if (res.status === 200) {
    ok200.add(1);
    ok200Trend.add(res.timings.duration);
  }

  if (res.status === 429) {
    tooMany429.add(1);
    tooManyTrend.add(res.timings.duration);
  }

  check(res, {
    "status 200": (r) => r.status === 200,
    "status 429": (r) => r.status === 429,
  });

  sleep(0.1);
}

export function handleSummary(data) {
  return {
    stdout: `
==== RATE LIMIT SIGN-IN SUMMARY ====

Total Requests: ${data.metrics.total_requests.values.count}
200 OK:         ${data.metrics.status_200.values.count}
429 Limited:    ${data.metrics.status_429.values.count}

=====================================
`,
    "rate-limit-results.json": JSON.stringify(data, null, 2),
  };
}
