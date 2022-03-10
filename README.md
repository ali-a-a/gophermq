# GopherMQ

![GitHub Workflow Status](https://img.shields.io/github/workflow/status/ali-a-a/gophermq/ci?label=ci&logo=github&style=flat-square)
[![Codecov](https://img.shields.io/codecov/c/gh/ali-a-a/gophermq?logo=codecov&style=flat-square)](https://codecov.io/gh/ali-a-a/gophermq)

## Diagram


<img width="512" alt="Screen Shot 2022-03-09 at 2 35 15 PM" src="https://user-images.githubusercontent.com/68470999/157550261-a37ed2f8-b788-4651-99a3-07dbe89a1917.png">

## Endpoints

- `/api/publish`
  ```json
  {
    "subject": "string",
    "data": "string"
  }
  ```
- `/api/publish/async`
  ```json
  {
    "subject": "string",
    "data": "string"
  }
  ```
- `/api/subscribe`
  ```json
  {
    "subject": "string",
  }
  ```

### Usage
```
git clone git@github.com:ali-a-a/gophermq.git
cd ./gophermq
make run-broker
now the server is up and running! 
```

### Structure
Messages can be published via `publish` endpoint. Then, Internally, the broker saves a new message in the in memory queue. After that, it loops through all the subscribers of the requested subject and calls their handler. If all of them are errorless, the queue is cleared. In case of any error, messages are kept. For async publish, it has `publish/async` endpoint. By this endpoint, a new message is submitted into the worker pool and then responds to the client.
  
### Overflow
Broker has a `MaxPending` option for handling overflow cases. MaxPending represents the maximum number of messages that can be stored in the broker. If a new publish causes overflow, the server returns a `broker overflow` error.

### MQ vs Shared Memory
Message queues enable asynchronous communication, which means that the endpoints that are producing and consuming messages interact with the queue, not the shared memory. Producers can add requests to the queue without waiting for them to be processed. Consumers process messages only when they are available.
