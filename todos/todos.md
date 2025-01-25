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

### 0002 Domain with Persistence - A New World
A simple model and storage for Users

- [ ] User Repositories
  - [ ] Interface
    - C
    - R
    - U
    - D
  - [ ] Implement/Test `go-memdb`
  - [ ] Implement/Test inmem



### 003 Authentication Server
`authentication/authorization server`

writing all modes via HTTP

- token-based for generic/B2B
- cookie-based for browsers


- [x] homepage
  - [x] Service: Hello world
  - [x] API `GET /`
    - [x] Content-Type: JSON
    - [x] Content-Type: HTML (templ)


- Story \
As a server, I want to authenticate users using 1FA password or 2FA password + email OTP.

The user can decide which to use, with the default being 1FA during registration.

#### Requirements
- Should be user-agent agnostic. The server will support the general case.
- Browser-free user flow:
  - A browser should use AJAX, then programmatically store auth data as required.
  - Not hungry - no cookies.

#### TODO: Mocking
- [x] register
  - [x] via username, password
- [x] `PATCH /user/`
  - [x] support changing `email` and `auth_mode`
- [x] login
  - [x] add modes
    - mode should be stored as user info, server-side
      ~~body should have `body.mode:String("SIMPLE_PW"|"2FA_PW_E")`~~
    - [x] Mode `SIMPLE_PW`: Simple via username, password
    - [x] Mode `2FA_PW_E`: via password with 2FA (password, then email otp)
      - [x] Factor 1: username/password
        `POST /token/", authHandlers.AuthByUsernamePassword()`
        - [x] if ok
          - [x] asynchronously send `OTP` to email with timeout.
          - [x] return 200 with:
            - [x] a weak jwt `weak_token` for user identification down flow.
              - [ ] add `auth_mode` in claims
      - [x] Factor 2:email
        `POST /otp/", authHandlers.SubmitOtp()`
        - [x] submit {`weak_token`, otp} for validation
          - [x] returns 200 if ok.

#### TODO: Implementation
- [ ] login
  - [ ] via password with 2FA (password, then email otp)
    - [ ] Factor 1: username/password
      `POST /token/", authHandlers.AuthByUsernamePassword()`
      - [ ] if ok
        - [ ] asynchronously send `OTP` to email with timeout.
        - [ ] return 200 with:
          - [ ] a weak jwt `weak_token` for user identification down flow.
    - [ ] Factor 2:email
      `POST /otp/", authHandlers.SubmitOtp()`
      - [ ] submit {`weak_token`, otp} for validation
        - [ ] returns 200 if ok.

