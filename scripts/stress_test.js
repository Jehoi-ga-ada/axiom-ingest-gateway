import http from 'k6/http';
import { check } from 'k6';

export const options = {
  scenarios: {
    global_stress_test: {
      executor: 'ramping-arrival-rate',
      startRate: 1000,
      timeUnit: '1s',
      preAllocatedVUs: 5000,       
      maxVUs: 10000,                
      stages: [
        { target: 100000, duration: '1m' }, 
      ],
    },
  },

  thresholds: {
    http_req_failed: ['rate<0.01'],             
    'http_req_duration{expected_res:true}': ['p(99)<500'], 
  },
};

export default function () {
  const url = 'http://localhost:8000/api/v1/events/';
  
  const eventTypes = ['click', 'view', 'purchase', 'login'];
  const randomType = eventTypes[Math.floor(Math.random() * eventTypes.length)];
  const randomId = Math.floor(Math.random() * 1000000);

  const payload = JSON.stringify({
    type: randomType,
    timestamp: new Date().toISOString(),
    payload: {
      id: randomId,
      source: "k6-stress-test",
      metadata: "random-data-" + (Math.random() + 1).toString(36).substring(7)
    }
  });

  const params = {
    headers: {
      'Content-Type': 'application/json',
      'X-API-Key': 'ax_test_key_1', 
    },
  };

  const res = http.post(url, payload, params);

  check(res, {
    'is status 201': (r) => r.status === 201,
  });
}