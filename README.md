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

## Structure
Messages can be published via `publish` endpoint. Then, Internally, the broker saves a new message in the in-memory queue. Note that at least one subscriber on the subject should exist before publishing a new message. Else, the publisher got `ErrSubscriberNotFound`. \
For async publish, it has a `publish/async` endpoint. By this endpoint, a new message is submitted into the worker pool and then responds to the client. \
For subscription on subjects, the server has a `subscribe` endpoint. In successful cases, it returns the id of the subscriber. This id is used in a `fetch` endpoint. \
For fetching messages, you should use the `fetch` endpoint. After calling this endpoint, all the pending messages for the specific subscriber are consumed.

## Overflow
Broker has a `MaxPending` option for handling overflow cases. MaxPending represents the maximum number of messages that can be stored in the broker. If a new publish causes overflow, the server returns a `broker overflow` error.

## MQ vs Shared Memory
Message queues enable asynchronous communication, which means that the endpoints that are producing and consuming messages interact with the queue, not the shared memory. Producers can add requests to the queue without waiting for them to be processed. Consumers process messages only when they are available.
