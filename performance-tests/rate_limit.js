import http from "k6/http";
import { check, sleep } from "k6";

export const options = {
  scenarios: {
    rate_limit_test: {
      executor: "constant-arrival-rate",
      rate: 100, // 100 requests per second
      timeUnit: "1m", // per minute
      duration: "1m",
      preAllocatedVUs: 10,
      maxVUs: 50,
    },
  },
};

export default function () {
  // Assuming 5 requests/min limit for this endpoint
  const res = http.get(
    "http://host.docker.internal:8999/api/v1/academic/public/health"
  );

  // We expect some 429s eventually if the limit is hit
  check(res, {
    "status is 200 or 429": (r) => r.status === 200 || r.status === 429,
  });
  sleep(1);
}
