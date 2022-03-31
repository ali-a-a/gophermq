import http from 'k6/http';
import {check, sleep} from 'k6';
import {Counter} from 'k6/metrics';

export const requests = new Counter('http_reqs');
export const options = {
    stages: [
        {vus: 10, duration: '10s', target: 50},
        {vus: 10, duration: '10s', target: 100},
        {vus: 10, duration: '10s', target: 150},
        {vus: 10, duration: '10s', target: 200},
        {vus: 10, duration: '10s', target: 250},
        {vus: 10, duration: '10s', target: 300},
        {vus: 10, duration: '10s', target: 350},
        {vus: 10, duration: '10s', target: 400},
        {vus: 10, duration: '10s', target: 450},
        {vus: 10, duration: '10s', target: 500},
        {vus: 10, duration: '160s', target: 550},
        {vus: 10, duration: '170s', target: 600},
        {vus: 10, duration: '180s', target: 650},
        {vus: 10, duration: '190s', target: 700},
        {vus: 10, duration: '200s', target: 750},
        {vus: 10, duration: '210s', target: 800},
        {vus: 10, duration: '220s', target: 850},
        {vus: 10, duration: '230s', target: 900},
        {vus: 10, duration: '240s', target: 950},
        {vus: 10, duration: '600s', target: 1000},
    ],
};

export function setup() {
    var url = 'http://localhost:8082/api/subscribe';

    const payload = JSON.stringify({
        subject: 'load.test',
    });

    const params = {
        headers: {
            'Content-Type': 'application/json',
        },
    };

    const res = http.post(url, payload, params);

    return JSON.parse(res.body).id
}


export default function (data) {
    var url = 'http://localhost:8082/api/publish';
    var fetchUrl = 'http://localhost:8082/api/fetch';

    const payload = JSON.stringify({
        subject: 'load.test',
        data: ':('
    });

    const params = {
        headers: {
            'Content-Type': 'application/json',
        },
    };

    const res = http.post(url, payload, params);

    const fetchPayload = JSON.stringify({
        subject: 'load.test',
        id: data
    })

    const fetchRes = http.post(fetchUrl, fetchPayload, params);

    const checkRes = check(res, {
        'status is 200': (r) => r.status === 200
    });

    console.log(fetchRes.status);
    console.log(fetchRes.body);

    const checkResFetch = check(res, {
        'status is 200': (fetchRes) => fetchRes.status === 200
    });

    sleep(1);
};