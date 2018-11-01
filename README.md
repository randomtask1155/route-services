

## Installation

```
cf push route-service -c "route-service -s"
cf push app -c "route-service"
```

```
cf create-user-provided-service danl-rs -r https://route-service.cfapps-02.domain.io
 cf bind-route-service cfapps-02.domain.io danl-rs --hostname app
```

## Test Installation 

the route service should have added Header `X-DanL-Route-Service: welcome to route services` to the request and forwarded that request to the app. 

```
~:> curl https://app.cfapps-02.domain.io -k | jq
{
  "message": "hello welcome to route services"
}
```


## Configuration Options

added an environmental vairalbe `SLEEP_INTERVAL` to either app will instruct the app to sleep in between requests.  This allows you to test how to troubleshoo which component is slow.