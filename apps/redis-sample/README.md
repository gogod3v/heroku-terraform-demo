# Readme

## Development

```shell
# start a Redis DB, e.g. via docker:
docker run -d -p 6379:6379 redis:4
# build the code
go build -o /tmp/.bin
# set REDIS_URL
export REDIS_URL='redis://:@localhost:6379'
# run the server binary
/tmp/.bin
```

You can then hit the following endpoints:

| Method | Route      | Description                                  |
| ------ | ---------- | -------------------------------------------- |
| GET    | /          | Root, welcome                                |
| GET    | /+         | Increment the counter, return its value      |
| GET    | /-         | Decrement the counter, return its value      |
