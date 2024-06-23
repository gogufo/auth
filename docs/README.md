# Documentation

## 1. Data Bases hosted by this microservice
- Users - username, password and login informaton
- AuthHistory - all login history with IP and Device

## 2. Dependencies from another microservice
- Session - update session, refresh session, create session
- Notifications - send TFA email, Send notifications in case of unrecognized device signed, Send emails with codes in case of forgot password, send emails with password
- Rights - to send OTP and block account hash tokens
- User - to get information about user
