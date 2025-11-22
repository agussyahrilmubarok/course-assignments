import http from 'k6/http';
import { check } from 'k6';

const policyCode = __ENV.policyCode || "DEFAULT_CODE"

export const options = {
    vus: 1,
    iterations: 1,
};

export function setup() {
    console.log("=== SETUP: Pre-test GET check quantity dummy V1 ===");
    const startTime = new Date().toISOString();
    return { startTime };
}

export default function (data) {
    const url = `http://localhost:8080/check-quantity-v1/${policyCode}`;

    const res = http.get(url);

    check(res, {
        "status 200": (r) => r.status === 200,
    });

    console.log("Response:", res.body);
}

export function teardown(data) {
  console.log("=== TEARDOWN: Post-test Get check quantity dummy v1 ===");
  console.log("StartAt:", data.startTime);
  console.log("EndAt:", new Date().toISOString());
}