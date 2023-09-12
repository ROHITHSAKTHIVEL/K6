import http from 'k6/http';
import { check, sleep } from 'k6';

export let options = {
  vus: 10, // Number of virtual users (simulated users)
  duration: '30s', // Test duration in seconds
};

export default function () {
  // Define the API endpoint
  const url = 'http://localhost:6000/tokens';

  // Send a POST request to create a token
  const payload = JSON.stringify({ token: 'your-auth-token' });
  const params = {
    headers: {
      'Content-Type': 'application/json',
    },
  };
  const res = http.post(url, payload, params);

  // Check if the response status code is 200
  check(res, {
    'Status is 200': (r) => r.status === 200,
  });

  // Sleep for a short period (e.g., 1 second) between requests
//   sleep(1);
}
