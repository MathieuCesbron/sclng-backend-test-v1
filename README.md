# Scalingo technical test

## Instructions

```bash
docker-compose up
```

Application will be then running on port `5000`

## Test

```
$ curl localhost:5000/ping
{ "status": "pong" }
```

## How is it optimized and scalable ?
- The image size is < 10MB