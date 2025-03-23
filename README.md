# go-websocket-server

Step 1: Connect to: `localhost:8080/ws`

Step 2: Send message to subcribe BTC & ETH
```json
{
  "method": "SUBSCRIBE",
  "params": [
    "btcusdt@kline_15m",
    "ethusdt@kline_15m"
  ]
}
```

Step 3: Open another connect to: `localhost:8080/ws`

Step 4: Send message to subcribe BNB & SOL
```json
{
  "method": "SUBSCRIBE",
  "params": [
    "bnbusdt@kline_15m",
    "solusdt@kline_15m"
  ]
}
```

Step 5: Send message to unsubcribe BTC & ETH
```json
{
  "method": "UNSUBSCRIBE",
  "params": [
    "btcusdt@kline_15m",
    "ethusdt@kline_15m"
  ]
}
```