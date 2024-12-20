todos


### 0001 Build and Run a program

A simple build with main entry point.

- [x] use context as first principle
  - [x] handles terminates on OS signals
      - [x] long lives
        - [x] when interrupt
          - 10 seconds timeout
            ```
            {"time":"2024-12-20T09:23:27.214861+08:00","level":"SYSTEM","msg":"started","callers":"main"}
            {"time":"2024-12-20T09:23:37.216096+08:00","level":"SYSTEM","msg":"ctx.Done() received","callers":"main"}
            {"time":"2024-12-20T09:23:37.216204+08:00","level":"SYSTEM","msg":"exited","callers":"main"}
            ```
          - interrupt
            ```
            {"time":"2024-12-20T09:26:20.048204+08:00","level":"SYSTEM","msg":"started","callers":"main"}
            ^C{"time":"2024-12-20T09:26:23.13586+08:00","level":"SYSTEM","msg":"interrupt or terminated","callers":"main"}
            {"time":"2024-12-20T09:26:23.135913+08:00","level":"SYSTEM","msg":"exited","callers":"main"}
            ```