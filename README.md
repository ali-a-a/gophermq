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
  
