# session login with JWT dgrijalva/jwt-go

Performance test:
```sh
autocannon --pipelining=10 --body '{ "username": "user1", "password": "password1"  }' -H Content-Type=application/json -m POST http://localhost:1323/login
Running 10s test @ http://localhost:1323/login
10 connections with 10 pipelining factor

┌─────────┬──────┬──────┬───────┬───────┬─────────┬──────────┬────────┐
│ Stat    │ 2.5% │ 50%  │ 97.5% │ 99%   │ Avg     │ Stdev    │ Max    │
├─────────┼──────┼──────┼───────┼───────┼─────────┼──────────┼────────┤
│ Latency │ 1 ms │ 3 ms │ 32 ms │ 50 ms │ 4.54 ms │ 10.15 ms │ 278 ms │
└─────────┴──────┴──────┴───────┴───────┴─────────┴──────────┴────────┘
┌───────────┬────────┬────────┬─────────┬─────────┬─────────┬──────────┬────────┐
│ Stat      │ 1%     │ 2.5%   │ 50%     │ 97.5%   │ Avg     │ Stdev    │ Min    │
├───────────┼────────┼────────┼─────────┼─────────┼─────────┼──────────┼────────┤
│ Req/Sec   │ 1855   │ 1855   │ 24943   │ 32927   │ 19867.2 │ 11642.98 │ 1855   │
├───────────┼────────┼────────┼─────────┼─────────┼─────────┼──────────┼────────┤
│ Bytes/Sec │ 571 kB │ 571 kB │ 7.68 MB │ 10.1 MB │ 6.12 MB │ 3.59 MB  │ 571 kB │
└───────────┴────────┴────────┴─────────┴─────────┴─────────┴──────────┴────────┘

Req/Bytes counts sampled once per second.

199k requests in 10.02s, 61.2 MB read
```