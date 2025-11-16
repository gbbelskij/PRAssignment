import http from 'k6/http';
import { check, group, sleep } from 'k6';
import { Rate, Trend, Counter } from 'k6/metrics';

export let options = {
  stages: [
    { duration: '5s', target: 5 },
    { duration: '10s', target: 20 },
    { duration: '5s', target: 0 },
  ],
};

let errorRate = new Rate('errors');
let requestDuration = new Trend('request_duration');
let successCount = new Counter('success');

const BASE_URL = 'http://localhost:8080';

export default function () {
  group('GET /team/get', () => {
    let res = http.get(`${BASE_URL}/team/get?team_name=backend`);
    let isSuccess = res.status === 200 || res.status === 404;
    
    check(res, {
      'status ok': (r) => isSuccess,
    });

    if (isSuccess) {
      successCount.add(1);
    } else {
      errorRate.add(1);
    }

    requestDuration.add(res.timings.duration);
  });

  sleep(1);
}
