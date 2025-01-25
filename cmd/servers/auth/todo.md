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
As a server, I want to authenticate users using 2FA.

## Requirements
- Should be user-agent agnostic. The server will support the general case.
- Browser-free user flow: 
  - A browser should use AJAX, then programmatically store auth data as required.
  - Not hungry - no cookies.

### TODO: Mocking
- [ ] login
  - [ ] via password with 2FA (password, then email otp)
    - [x] Factor 1: username/password
      `POST /password/", authHandlers.AuthByUsernamePassword()`
      - [x] if ok
        - [x] asynchronously send `OTP` to email with timeout.
        - [x] return 200 with: 
          - [x] a weak jwt `weak_token` for user identification down flow.
    - [x] Factor 2:email
      `POST /otp/", authHandlers.SubmitOtp()`
      - [x] submit {`weak_token`, otp} for validation
        - [x] returns 200 if ok.



