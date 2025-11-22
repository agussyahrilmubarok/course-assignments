import http from 'k6/http';
import { check } from 'k6';

const policyCode = __ENV.policyCode || 'DEFAULT_CODE';

export const options = {
  vus: 100,
  iterations: 100,
};