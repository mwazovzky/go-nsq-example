# GO NSQ Example

## Ping

```
curl http://localhost:4151/ping
curl http://localhost:4161/ping
curl http://localhost:4171/ping
```

## Publish message

```
curl -d '{"Title":"Message Title","Content":"Message body","Timestamp":"2022-04-02 23:23"}' 'http://localhost:4151/pub?topic=test_topic'
```
