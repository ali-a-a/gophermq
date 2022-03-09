# GopherMQ

## Diagram

<img width="674" alt="Screen Shot 2022-03-09 at 2 03 24 PM" src="https://user-images.githubusercontent.com/68470999/157544200-5ac5c29c-acf4-4e97-841d-b4b5b7bcc474.png">

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
Messages can be published via publish endpoint. Then, Internally, the broker saves a new message in the in memory queue. After that, it loops through all the subscribers of the requested subject and calls their handler. If all of them are errorless, the queue is cleared. In case of any error, messages are kept. For async publish, it has publish/async endpoint. By this endpoint, a new message is submitted into the worker pool and then responds to the client.
  
### Overflow
Broker has `MaxPending` option for handling overflow cases. `MaxPending` represents the maximum
number of messages can be stored in the broker. If new publish causes overflow, the server returns `broker overflow` error.
