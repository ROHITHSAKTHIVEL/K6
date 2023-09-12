import http from 'k6/http';
import { check, sleep } from 'k6';

export default function () {
  // Define the API endpoint
  const url = 'http://localhost:6000/hello';

  // Define the JSON request body
  const payload = JSON.stringify({
    message: 'Hello, world!',
  });

  // Define HTTP headers
  const headers = {
    'Content-Type': 'application/json',
  };

  // Send a POST request to the API
  const res = http.post(url, payload, { headers });

  // Check if the response status code is 200
  check(res, {
    'Status is 200': (r) => r.status === 200,
  });

  // Sleep for a short period (e.g., 1 second) between requests
//   sleep(1);
}
