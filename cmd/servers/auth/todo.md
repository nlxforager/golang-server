# auth

`authentication/authorization server`

writing all modes via HTTP

- token-based for generic/B2B
- cookie-based for browsers




- [x] homepage
  - [x] Service: Hello world
  - [x] API `GET /`
    - [x] Content-Type: JSON
    - [x] Content-Type: HTML (templ)


## Story
As a server, I want to authenticate users using 1FA password or 2FA password + email.

The user can decide which to use, with the default being 1FA during registration.

## Requirements
- Should be user-agent agnostic. The server will support the general case.
- Browser-free user flow: 
  - A browser should use AJAX, then programmatically store auth data as required.
  - Not hungry - no cookies.

### TODO: Mocking
- [ ] register via username, password
- [x] login
  - [ ] add modes
      - mode should be stored as user info
      ~~body should have `body.mode:String("SIMPLE_PW"|"2FA_PW_E")`~~
    - [ ] Mode `SIMPLE_PW`: Simple via username, password 
    - [x] Mode `2FA_PW_E`: via password with 2FA (password, then email otp)
      - [x] Factor 1: username/password
        `POST /password/", authHandlers.AuthByUsernamePassword()`
        - [x] if ok
          - [x] asynchronously send `OTP` to email with timeout.
          - [x] return 200 with: 
            - [x] a weak jwt `weak_token` for user identification down flow.
              - [ ] add `auth_mode` in claims
      - [x] Factor 2:email
        `POST /otp/", authHandlers.SubmitOtp()`
        - [x] submit {`weak_token`, otp} for validation
          - [x] returns 200 if ok.


### TODO: Implementation
- [ ] login
  - [ ] via password with 2FA (password, then email otp)
    - [ ] Factor 1: username/password
      `POST /password/", authHandlers.AuthByUsernamePassword()`
      - [ ] if ok
        - [ ] asynchronously send `OTP` to email with timeout.
        - [ ] return 200 with:
          - [ ] a weak jwt `weak_token` for user identification down flow.
    - [ ] Factor 2:email
      `POST /otp/", authHandlers.SubmitOtp()`
      - [ ] submit {`weak_token`, otp} for validation
        - [ ] returns 200 if ok.
