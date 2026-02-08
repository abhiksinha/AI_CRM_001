# Edge Service API cURL Examples

This document provides `curl` command examples for the Edge service.

**Note:** Protected routes require an `Authorization: Bearer <token>` header. Use `/v1/login` or `/v1/token` to obtain a token.

---

## Auth

### 1. Login

**Endpoint:** `POST /v1/login`

```sh
curl --location 'http://127.0.0.1:8081/v1/login' \
--header 'Content-Type: application/json' \
--data '{
  "username": "testuser@example.com",
  "password": "a-very-secure-password"
}'
```

### 2. Logout

**Endpoint:** `POST /v1/logout`

```sh
curl --location 'http://127.0.0.1:8081/v1/logout' \
--header 'Content-Type: application/json' \
--data '{
  "api_token": "THE_API_TOKEN"
}'
```

### 3. List Sessions

**Endpoint:** `POST /v1/sessions`

```sh
curl --location 'http://127.0.0.1:8081/v1/sessions' \
--header 'Content-Type: application/json' \
--data '{
  "username": "testuser@example.com",
  "password": "a-very-secure-password"
}'
```

### 4. Expire Session

**Endpoint:** `POST /v1/sessions/expire`

```sh
curl --location 'http://127.0.0.1:8081/v1/sessions/expire' \
--header 'Content-Type: application/json' \
--data '{
  "session_id": "SESSION_ID"
}'
```

### 5. Get Token (S2S)

**Endpoint:** `POST /v1/token`

```sh
curl --location 'http://127.0.0.1:8081/v1/token' \
--header 'Content-Type: application/json' \
--data '{
  "api_key": "YOUR_API_KEY"
}'
```

---

## Contacts

### 6. Create a Contact

**Endpoint:** `POST /v1/contacts`

```sh
curl --location 'http://127.0.0.1:8081/v1/contacts' \
--header 'Authorization: Bearer THE_API_TOKEN' \
--header 'Content-Type: application/json' \
--data '{
  "first_name":"Abhi",
  "last_name": "Sinha",
  "email": "abhisinha@example.com",
  "phone": "+911234567890",
  "owner_id":"ABHISHEKSINHA1"
}'
```

### 7. Get a Contact by ID

**Endpoint:** `GET /v1/contacts/{contactID}`

```sh
curl --location 'http://127.0.0.1:8081/v1/contacts/0ML3H99X5QNE0E' \
--header 'Authorization: Bearer THE_API_TOKEN'
```

### 8. List Contacts

**Endpoint:** `GET /v1/contacts`

```sh
curl --location 'http://127.0.0.1:8081/v1/contacts?page=1&page_size=10' \
--header 'Authorization: Bearer THE_API_TOKEN'
```

### 9. Update a Contact

**Endpoint:** `PUT /v1/contacts/{contactID}`

```sh
curl --location --request PUT 'http://127.0.0.1:8081/v1/contacts/0ML3H99X5QNE0E' \
--header 'Authorization: Bearer THE_API_TOKEN' \
--header 'Content-Type: application/json' \
--data '{
  "phone": "+919876543210"
}'
```

### 10. Delete a Contact

**Endpoint:** `DELETE /v1/contacts/{contactID}`

```sh
curl --location --request DELETE 'http://127.0.0.1:8081/v1/contacts/0ML3H99X5QNE0E' \
--header 'Authorization: Bearer THE_API_TOKEN'
```

### 11. Add a Note to a Contact

**Endpoint:** `POST /v1/contacts/{contactID}/notes`

```sh
curl --location 'http://127.0.0.1:8081/v1/contacts/0ML3H99X5QNE0E/notes' \
--header 'Authorization: Bearer THE_API_TOKEN' \
--header 'Content-Type: application/json' \
--data '{
  "content": "Called the client, scheduled a follow-up for next week."
}'
```

### 12. List Notes for a Contact

**Endpoint:** `GET /v1/contacts/{contactID}/notes`

```sh
curl --location 'http://127.0.0.1:8081/v1/contacts/0ML3H99X5QNE0E/notes' \
--header 'Authorization: Bearer THE_API_TOKEN'
```

---

## Deals

### 13. Create a Deal

**Endpoint:** `POST /v1/deals`

```sh
curl --location 'http://127.0.0.1:8081/v1/deals' \
--header 'Authorization: Bearer THE_API_TOKEN' \
--header 'Content-Type: application/json' \
--data '{
  "name": "New Website Project",
  "stage": "Qualification",
  "value": 50000.00,
  "expected_close_date": "2024-12-31",
  "contact_id": "0ML3H99X5QNE0E"
}'
```

### 14. Get a Deal by ID

**Endpoint:** `GET /v1/deals/{dealID}`

```sh
curl --location 'http://127.0.0.1:8081/v1/deals/DEAL_ID' \
--header 'Authorization: Bearer THE_API_TOKEN'
```

### 15. List Deals

**Endpoint:** `GET /v1/deals`

```sh
curl --location 'http://127.0.0.1:8081/v1/deals?page=1&page_size=10&stage=Qualification' \
--header 'Authorization: Bearer THE_API_TOKEN'
```

### 16. Update a Deal

**Endpoint:** `PUT /v1/deals/{dealID}`

```sh
curl --location --request PUT 'http://127.0.0.1:8081/v1/deals/DEAL_ID' \
--header 'Authorization: Bearer THE_API_TOKEN' \
--header 'Content-Type: application/json' \
--data '{
  "stage": "Proposal Sent",
  "value": 55000.00
}'
```

### 17. Delete a Deal

**Endpoint:** `DELETE /v1/deals/{dealID}`

```sh
curl --location --request DELETE 'http://127.0.0.1:8081/v1/deals/DEAL_ID' \
--header 'Authorization: Bearer THE_API_TOKEN'
```

### 18. Create a Task for a Deal

**Endpoint:** `POST /v1/deals/{dealID}/tasks`

```sh
curl --location 'http://127.0.0.1:8081/v1/deals/DEAL_ID/tasks' \
--header 'Authorization: Bearer THE_API_TOKEN' \
--header 'Content-Type: application/json' \
--data '{
  "title": "Follow up with client about proposal",
  "due_date": "2024-11-15",
  "assigned_to_id": "ABHISHEKSINHA1"
}'
```

### 19. List Tasks for a Deal

**Endpoint:** `GET /v1/deals/{dealID}/tasks`

```sh
curl --location 'http://127.0.0.1:8081/v1/deals/DEAL_ID/tasks' \
--header 'Authorization: Bearer THE_API_TOKEN'
```

---

## Users

### 20. Create a User

**Endpoint:** `POST /v1/users`

```sh
curl --location 'http://127.0.0.1:8081/v1/users' \
--header 'Content-Type: application/json' \
--data '{
  "first_name": "Test",
  "last_name": "User",
  "email": "testuser@example.com",
  "password": "a-very-secure-password",
  "role": "user"
}'
```

### 21. Get a User by ID

**Endpoint:** `GET /v1/users/{userID}`

```sh
curl --location 'http://127.0.0.1:8081/v1/users/ABHISHEKSINHA1' \
--header 'Authorization: Bearer THE_API_TOKEN'
```

### 22. List Users

**Endpoint:** `GET /v1/users`

```sh
curl --location 'http://127.0.0.1:8081/v1/users?page=1&page_size=10&role=user' \
--header 'Authorization: Bearer THE_API_TOKEN'
```

### 23. Update a User

**Endpoint:** `PUT /v1/users/{userID}`

```sh
curl --location --request PUT 'http://127.0.0.1:8081/v1/users/ABHISHEKSINHA1' \
--header 'Authorization: Bearer THE_API_TOKEN' \
--header 'Content-Type: application/json' \
--data '{
  "is_active": false
}'
```

### 24. Delete a User

**Endpoint:** `DELETE /v1/users/{userID}`

```sh
curl --location --request DELETE 'http://127.0.0.1:8081/v1/users/ABHISHEKSINHA1' \
--header 'Authorization: Bearer THE_API_TOKEN'
```

---

## API Keys

### 25. Create an API Key

**Endpoint:** `POST /v1/users/api-keys`

```sh
curl --location 'http://127.0.0.1:8081/v1/users/api-keys' \
--header 'Authorization: Bearer THE_API_TOKEN' \
--header 'Content-Type: application/json' \
--data '{
  "user_id": "ABHISHEKSINHA1"
}'
```

### 26. Expire an API Key

**Endpoint:** `DELETE /v1/users/api-keys`

```sh
curl --location --request DELETE 'http://127.0.0.1:8081/v1/users/api-keys' \
--header 'Authorization: Bearer THE_API_TOKEN' \
--header 'Content-Type: application/json' \
--data '{
  "api_key": "THE_API_KEY_TO_EXPIRE"
}'
```
