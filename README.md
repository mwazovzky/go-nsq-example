# GO NSQ Example

## Ping

```
curl http://localhost:4151/ping
curl http://localhost:4161/ping
curl http://localhost:4171/ping

curl -d 'Hello World' 'http://localhost:4151/pub?topic=test'

curl http://localhost:4161/topics
{"topics":["test"]}
```
