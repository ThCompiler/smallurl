import { sleep } from 'k6';
import http from 'k6/http';
import { randomIntBetween, randomItem } from 'https://jslib.k6.io/k6-utils/1.2.0/index.js';

export const options = {
    scenarios: {
        test: {
            executor: 'constant-arrival-rate',
            duration: '5m',
            preAllocatedVUs: 10,

            rate: __ENV.RATE_COUNT,
            timeUnit: '1s',
            maxVUs: 40,
        },
    },
    discardResponseBodies: true,
};

export default function () {
    if (__ENV.METHOD === "get") {
        // define URL and request body
        const params = {
            tags: { name: "smallurl" },
            redirects: 0,
        };

        const postId = randomIntBetween(0, 100000)

        // send a post request and save response as a variable
        http.get(`http://localhost:8080/api/v1/${postId}/`, params);
    } else {
        // define URL and request body
        const url = `http://localhost:8080/api/v1/shorten`;
        const params = {
            headers: { 'Content-Type': 'application/json' },
            tags: { name: "smallurl" },
        };

        const postId = randomIntBetween(0, 100000)

        // send a post request and save response as a variable
        http.post(url, JSON.stringify({original_url: `https://habr.com/ru/articles/${postId}/`}), params);
    }
}
