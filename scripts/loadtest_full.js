import http from 'k6/http';
import { check, group, sleep } from 'k6';
import { randomString } from 'https://jslib.k6.io/k6-utils/1.2.0/index.js';

export let options = {
  stages: [
    { duration: '10s', target: 10 },
    { duration: '30s', target: 50 },
    { duration: '20s', target: 0 },
  ],
  thresholds: {
    http_req_duration: ['p(95)<500', 'p(99)<1000'],
    http_req_failed: ['rate<0.3'],
  },
};

const BASE_URL = 'http://localhost:8080';

export default function () {
  group('Team Operations', () => {
    let teamPayload = {
      team_name: `team-${Math.random()}`,
      members: [
        {
          user_id: 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11',
          username: 'Alice',
          is_active: true,
        },
        {
          user_id: 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a12',
          username: 'Bob',
          is_active: true,
        },
      ],
    };

    let createRes = http.post(`${BASE_URL}/team/add`, JSON.stringify(teamPayload), {
      headers: { 'Content-Type': 'application/json' },
    });

    check(createRes, {
      'team add status 201': (r) => r.status === 201,
    });

    let getRes = http.get(`${BASE_URL}/team/get?team_name=backend`);
    check(getRes, {
      'team get status 200': (r) => r.status === 200 || r.status === 404,
      'response time < 200ms': (r) => r.timings.duration < 200,
    });
  });

  group('PR Operations', () => {
    let prPayload = {
        pull_request_id: `pr-${randomString(10)}`,
        pull_request_name: `Feature ${randomString(8)}`,
        author_id: 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11',
    };


    let createRes = http.post(`${BASE_URL}/pullRequest/create`, JSON.stringify(prPayload), {
      headers: { 'Content-Type': 'application/json' },
    });

    check(createRes, {
      'pr create status 201': (r) => r.status === 201 || r.status === 404 || r.status === 409,
    });
  });

  group('Users Operations', () => {
    let reviewRes = http.get(
      `${BASE_URL}/users/getReview?user_id=a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11`
    );
    check(reviewRes, {
      'get review status 200': (r) => r.status === 200,
    });
  });

  sleep(1);
}
