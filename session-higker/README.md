# Session using higker/session

## Performance Test with Autocannon

### GET
```sh
~ autocannon --pipelining=10 http://localhost:1323/get
Running 10s test @ http://localhost:1323/get
10 connections with 10 pipelining factor

┌─────────┬──────┬──────┬───────┬───────┬─────────┬─────────┬────────┐
│ Stat    │ 2.5% │ 50%  │ 97.5% │ 99%   │ Avg     │ Stdev   │ Max    │
├─────────┼──────┼──────┼───────┼───────┼─────────┼─────────┼────────┤
│ Latency │ 4 ms │ 7 ms │ 13 ms │ 23 ms │ 7.66 ms │ 5.26 ms │ 129 ms │
└─────────┴──────┴──────┴───────┴───────┴─────────┴─────────┴────────┘
┌───────────┬─────────┬─────────┬─────────┬─────────┬─────────┬─────────┬─────────┐
│ Stat      │ 1%      │ 2.5%    │ 50%     │ 97.5%   │ Avg     │ Stdev   │ Min     │
├───────────┼─────────┼─────────┼─────────┼─────────┼─────────┼─────────┼─────────┤
│ Req/Sec   │ 9351    │ 9351    │ 12703   │ 14679   │ 12252.8 │ 1750.61 │ 9344    │
├───────────┼─────────┼─────────┼─────────┼─────────┼─────────┼─────────┼─────────┤
│ Bytes/Sec │ 2.88 MB │ 2.88 MB │ 3.91 MB │ 4.52 MB │ 3.77 MB │ 539 kB  │ 2.88 MB │
└───────────┴─────────┴─────────┴─────────┴─────────┴─────────┴─────────┴─────────┘

Req/Bytes counts sampled once per second.

123k requests in 10.02s, 37.7 MB read
```

### POST
```sh
~ autocannon --pipelining=10 --body '{ "ID": 0, "Name": "Ding", "Email": "", "Address": "", "Cart": [ { "ProductID": 0, "ProductName": "", "Qty": 0, "Price": 0 } ] }' -H Content-Type=application/json -m POST http://localhost:1323/set
Running 10s test @ http://localhost:1323/set
10 connections with 10 pipelining factor

┌─────────┬──────┬───────┬───────┬───────┬──────────┬──────────┬────────┐
│ Stat    │ 2.5% │ 50%   │ 97.5% │ 99%   │ Avg      │ Stdev    │ Max    │
├─────────┼──────┼───────┼───────┼───────┼──────────┼──────────┼────────┤
│ Latency │ 5 ms │ 10 ms │ 51 ms │ 67 ms │ 19.88 ms │ 20.03 ms │ 293 ms │
└─────────┴──────┴───────┴───────┴───────┴──────────┴──────────┴────────┘
┌───────────┬────────┬────────┬────────┬─────────┬─────────┬─────────┬────────┐
│ Stat      │ 1%     │ 2.5%   │ 50%    │ 97.5%   │ Avg     │ Stdev   │ Min    │
├───────────┼────────┼────────┼────────┼─────────┼─────────┼─────────┼────────┤
│ Req/Sec   │ 2101   │ 2101   │ 3343   │ 10911   │ 4910    │ 2929.76 │ 2101   │
├───────────┼────────┼────────┼────────┼─────────┼─────────┼─────────┼────────┤
│ Bytes/Sec │ 542 kB │ 542 kB │ 863 kB │ 2.81 MB │ 1.27 MB │ 756 kB  │ 542 kB │
└───────────┴────────┴────────┴────────┴─────────┴─────────┴─────────┴────────┘

Req/Bytes counts sampled once per second.

49k requests in 10.03s, 12.7 MB read
```

## Contras using redis
1. Storing user session on Redis is not bad practice, however using Redis for authorization/authentication purposes could result in security issues. In an ideal microservice architecture, services need to be as stateless as they can. By giving the authorization job to redis you are creating an extra layer that you need to manage on every request. What will happen when user info changes? You will need to update redis too and when it is out of sync with the Authorization server, then problems will start to emerge. It is also not bad practice to use separate Authorization server behind your gateway. Source: https://stackoverflow.com/questions/63237218/api-gateway-login-microservice-redis

2. There are problem using this scheme that scalability is a problem, and if the Redis cache instance goes down then all the sessions will be deleted for current logged in users and the entire system will remain non-functional, though the services are live. To solve this problem JWT (JSON Web Token). Source: https://dev.to/rishavsaha98/user-authorization-in-micro-service-architecture-with-jwt-2p6k