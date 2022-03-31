# GopherMQ

![GitHub Workflow Status](https://img.shields.io/github/workflow/status/ali-a-a/gophermq/ci?label=ci&logo=github&style=flat-square)
[![Codecov](https://img.shields.io/codecov/c/gh/ali-a-a/gophermq?logo=codecov&style=flat-square)](https://codecov.io/gh/ali-a-a/gophermq)

## Diagram


<img width="512" alt="Screen Shot 2022-03-09 at 2 35 15 PM" src="https://user-images.githubusercontent.com/68470999/157550261-a37ed2f8-b788-4651-99a3-07dbe89a1917.png">

## HTTP Server
For serving HTTP requests, the broker uses `Echo`, High performance, extensible, minimalist Go web framework.
By default, it listens on port `8082`.
You may find more information about this framework on this [link](https://github.com/labstack/echo).
## Endpoints

- `/api/publish` \
  \
  Request:
  ```json
  {
    "subject": "string",
    "data": "string"
  }
  ```
  \
  Responses:
  ```json
  {
    "status": "ok"
  }
  ```
  or
  ```json
  {
    "message": "error message"
  }
  ```
- `/api/publish/async` \
  \
  Response:
  ```json
  {
    "subject": "string",
    "data": "string"
  }
  ```
  \
  Response:
  ```json
  {
    "status": "ok"
  }
  ```
- `/api/subscribe` \
  \
  Request:
  ```json
  {
    "subject": "string"
  }
  ```
  \
  Responses:
  ```json
  {
    "id": "string",
    "subject": "string"
  }
  ```
  or
  ```json
  {
    "message": "error message"
  }
  ```
- `/api/fetch` \
  \
  Request:
  ```json
  {
    "subject": "string",
    "id": "string"
  }
  ```
  \
  Responses:
  ```json
  {
    "subject": "string",
    "id": "string",
    "data": []
  }
  ```
  or
  ```json
  {
    "message": "error message"
  }
  ```

## Usage
```
git clone git@github.com:ali-a-a/gophermq.git
cd ./gophermq
make run-broker
now the server is up and running! 
```

Or Using Docker

```
git clone git@github.com:ali-a-a/gophermq.git
cd ./gophermq
docker build -t gophermq . 
docker run -p 8082:8082 gophermq broker
```

## Metrics
Prometheus server is listening on port 9001.
handler of the server is registered with `/metrics` pattern.
Up to now, the broker has only 2 metrics. Request rate and request latency.

## Load Test
I used [k6](https://k6.io/) for load testing Gophermq. \
\
Machine:
```
RAM: 8GB
CPU Cores: 8
CPU Type: M1 chip
```
\
Results: \
\
<img width="646" alt="Screen Shot 2022-03-31 at 6 10 18 AM" src="https://user-images.githubusercontent.com/68470999/161066259-057ff49e-6fe9-478d-9dc4-1acad40fa8d2.png">
<img width="941" alt="Screen Shot 2022-03-31 at 6 21 15 AM" src="https://user-images.githubusercontent.com/68470999/161066325-10babebe-b2ae-4b2e-b398-571a5a4b8b98.png">

For running load test on your machine:
```
k6 run loadtest/loadtest.js
```

## Structure
Messages can be published via `publish` endpoint. Then, Internally, the broker saves a new message in the in-memory queue. Note that at least one subscriber on the subject should exist before publishing a new message. Else, the publisher got `ErrSubscriberNotFound`. \
For async publish, it has a `publish/async` endpoint. By this endpoint, a new message is submitted into the worker pool and then responds to the client. \
For subscription on subjects, the server has a `subscribe` endpoint. In successful cases, it returns the id of the subscriber. This id is used in a `fetch` endpoint. \
For fetching messages, you should use the `fetch` endpoint. After calling this endpoint, all the pending messages for the specific subscriber are consumed.

## Overflow
Broker has a `MaxPending` option for handling overflow cases. MaxPending represents the maximum number of messages that can be stored in the broker. If a new publish causes overflow, the server returns a `broker overflow` error.

## MQ vs Shared Memory
Message queues enable asynchronous communication, which means that the endpoints that are producing and consuming messages interact with the queue, not the shared memory. Producers can add requests to the queue without waiting for them to be processed. Consumers process messages only when they are available.

[k6]: https://k6.io/
