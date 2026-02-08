# AI CRM API cURL Examples

This document provides `curl` command examples for all the available API endpoints.

**Note:** All requests require an `X-User-ID` header. For these examples, we will use the ID `ABHISHEKSINHA1`.

---

## Contact Service

### 1. Create a Contact

Creates a new contact record.

**Endpoint:** `POST /api/v1/contacts`

```sh
curl --location 'http://127.0.0.1:8080/api/v1/contacts' \
--header 'X-User-ID: ABHISHEKSINHA1' \
--header 'Content-Type: application/json' \
--data '{
    "first_name":"Abhi",
    "last_name": "Sinha",
    "email": "abhisinha@example.com",
    "phone": "+911234567890",
    "owner_id":"ABHISHEKSINHA1"
}'
```

### 2. Get a Contact by ID

Retrieves a single contact by its unique ID.

**Endpoint:** `GET /api/v1/contacts/{contactID}`

```sh
curl --location 'http://127.0.0.1:8080/api/v1/contacts/0ML3H99X5QNE0E' \
--header 'X-User-ID: ABHISHEKSINHA1'
```

### 3. List Contacts

Retrieves a paginated list of contacts for the owner specified in `X-User-ID`.

**Endpoint:** `GET /api/v1/contacts`

```sh
curl --location 'http://127.0.0.1:8080/api/v1/contacts?page=1&page_size=10' \
--header 'X-User-ID: ABHISHEKSINHA1'
```

### 4. Update a Contact

Updates one or more fields of an existing contact.

**Endpoint:** `PUT /api/v1/contacts/{contactID}`

```sh
curl --location --request PUT 'http://127.0.0.1:8080/api/v1/contacts/0ML3H99X5QNE0E' \
--header 'X-User-ID: ABHISHEKSINHA1' \
--header 'Content-Type: application/json' \
--data '{
    "phone": "+919876543210"
}'
```

### 5. Delete a Contact

Deletes a contact by its ID.

**Endpoint:** `DELETE /api/v1/contacts/{contactID}`

```sh
curl --location --request DELETE 'http://127.0.0.1:8080/api/v1/contacts/0ML3H99X5QNE0E' \
--header 'X-User-ID: ABHISHEKSINHA1'
```

### 6. Add a Note to a Contact

Adds a new note to a specific contact.

**Endpoint:** `POST /api/v1/contacts/{contactID}/notes`

```sh
curl --location 'http://127.0.0.1:8080/api/v1/contacts/0ML3H99X5QNE0E/notes' \
--header 'X-User-ID: ABHISHEKSINHA1' \
--header 'Content-Type: application/json' \
--data '{
    "content": "Called the client, scheduled a follow-up for next week."
}'
```

### 7. List Notes for a Contact

Retrieves all notes associated with a specific contact.

**Endpoint:** `GET /api/v1/contacts/{contactID}/notes`

```sh
curl --location 'http://127.0.0.1:8080/api/v1/contacts/0ML3H99X5QNE0E/notes' \
--header 'X-User-ID: ABHISHEKSINHA1'
```

---

## Deal Service

### 8. Create a Deal

Creates a new deal.

**Endpoint:** `POST /api/v1/deals`

```sh
curl --location 'http://127.0.0.1:8080/api/v1/deals' \
--header 'X-User-ID: ABHISHEKSINHA1' \
--header 'Content-Type: application/json' \
--data '{
    "name": "New Website Project",
    "stage": "Qualification",
    "value": 50000.00,
    "expected_close_date": "2024-12-31",
    "contact_id": "0ML3H99X5QNE0E"
}'
```

### 9. Get a Deal by ID

Retrieves a single deal by its unique ID.

**Endpoint:** `GET /api/v1/deals/{dealID}`

```sh
curl --location 'http://127.0.0.1:8080/api/v1/deals/{dealID}' \
--header 'X-User-ID: ABHISHEKSINHA1'
```

### 10. List Deals

Retrieves a paginated list of deals. Can be filtered by `stage`.

**Endpoint:** `GET /api/v1/deals`

```sh
curl --location 'http://127.0.0.1:8080/api/v1/deals?page=1&page_size=10&stage=Qualification' \
--header 'X-User-ID: ABHISHEKSINHA1'
```

### 11. Update a Deal

Updates one or more fields of an existing deal.

**Endpoint:** `PUT /api/v1/deals/{dealID}`

```sh
curl --location --request PUT 'http://127.0.0.1:8080/api/v1/deals/{dealID}' \
--header 'X-User-ID: ABHISHEKSINHA1' \
--header 'Content-Type: application/json' \
--data '{
    "stage": "Proposal Sent",
    "value": 55000.00
}'
```

### 12. Delete a Deal

Deletes a deal by its ID.

**Endpoint:** `DELETE /api/v1/deals/{dealID}`

```sh
curl --location --request DELETE 'http://127.0.0.1:8080/api/v1/deals/{dealID}' \
--header 'X-User-ID: ABHISHEKSINHA1'
```

### 13. Create a Task for a Deal

Creates a new task associated with a deal.

**Endpoint:** `POST /api/v1/deals/{dealID}/tasks`

```sh
curl --location 'http://127.0.0.1:8080/api/v1/deals/{dealID}/tasks' \
--header 'X-User-ID: ABHISHEKSINHA1' \
--header 'Content-Type: application/json' \
--data '{
    "title": "Follow up with client about proposal",
    "due_date": "2024-11-15",
    "assigned_to_id": "ABHISHEKSINHA1"
}'
```

### 14. List Tasks for a Deal

Retrieves all tasks for a specific deal.

**Endpoint:** `GET /api/v1/deals/{dealID}/tasks`

```sh
curl --location 'http://127.0.0.1:8080/api/v1/deals/{dealID}/tasks' \
--header 'X-User-ID: ABHISHEKSINHA1'
```

---

## User Service

### 15. Create a User

Creates a new user.

**Endpoint:** `POST /api/v1/users`

```sh
curl --location 'http://127.0.0.1:8080/api/v1/users' \
--header 'Content-Type: application/json' \
--data '{
    "first_name": "Test",
    "last_name": "User",
    "email": "testuser@example.com",
    "password": "a-very-secure-password",
    "role": "user"
}'
```

### 16. Verify User Password

Checks if a given password is valid for a user.

**Endpoint:** `POST /api/v1/users/verify-password`

```sh
curl --location 'http://127.0.0.1:8080/api/v1/users/verify-password' \
--header 'Content-Type: application/json' \
--data '{
    "id": "ABHISHEKSINHA1",
    "password": "the-password-to-check"
}'
```

### 17. Get a User by ID

Retrieves a single user by their ID.

**Endpoint:** `GET /api/v1/users/{userID}`

```sh
curl --location 'http://127.0.0.1:8080/api/v1/users/ABHISHEKSINHA1'
```

### 18. List Users

Retrieves a paginated list of users. Can be filtered by `role`.

**Endpoint:** `GET /api/v1/users`

```sh
curl --location 'http://127.0.0.1:8080/api/v1/users?page=1&page_size=10&role=user'
```

### 19. Update a User

Updates one or more fields of an existing user.

**Endpoint:** `PUT /api/v1/users/{userID}`

```sh
curl --location --request PUT 'http://127.0.0.1:8080/api/v1/users/ABHISHEKSINHA1' \
--header 'Content-Type: application/json' \
--data '{
    "is_active": false
}'
```

### 20. Delete a User

Deletes a user by their ID.

**Endpoint:** `DELETE /api/v1/users/{userID}`

```sh
curl --location --request DELETE 'http://127.0.0.1:8080/api/v1/users/{userID}'
```

### 21. Create an API Key

Generates a new, unique API key for a user. A user can only have one active key at a time.

**Endpoint:** `POST /api/v1/users/api-keys`

```sh
curl --location 'http://127.0.0.1:8080/api/v1/users/api-keys' \
--header 'Content-Type: application/json' \
--data '{
    "user_id": "ABHISHEKSINHA1"
}'
```

### 22. Match an API Key

Verifies if an API key is valid and active. This is a key part of an authentication middleware.

**Endpoint:** `POST /api/v1/users/api-keys/match`

```sh
curl --location 'http://127.0.0.1:8080/api/v1/users/api-keys/match' \
--header 'Content-Type: application/json' \
--data '{
    "api_key": "THE_API_KEY_RETURNED_FROM_CREATE"
}'
```

### 23. Expire an API Key

Deactivates (soft-deletes) an API key, rendering it invalid.

**Endpoint:** `DELETE /api/v1/users/api-keys`

```sh
curl --location --request DELETE 'http://127.0.0.1:8080/api/v1/users/api-keys' \
--header 'Content-Type: application/json' \
--data '{
    "api_key": "THE_API_KEY_TO_EXPIRE"
}'
```
