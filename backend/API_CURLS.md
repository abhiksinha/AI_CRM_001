# AI CRM API cURL Examples

This document provides `curl` command examples for all the available API endpoints.

**Note:** All requests require an `X-User-ID` header. For these examples, we will use the ID `ABHISHEKSINHA1`.

---

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

**Example Success Response (201 Created):**
```json
{
    "id": "0ML3H99X5QNE0E",
    "name": "Abhi Sinha",
    "email": "abhisinha@example.com",
    "phone": "+911234567890",
    "owner_id": "ABHISHEKSINHA1",
    "owner_name": "Abhishek Sinha",
    "created_at": 1769934423,
    "updated_at": 1769934423
}
```

---

### 2. Get a Contact by ID

Retrieves a single contact by its unique ID.

**Endpoint:** `GET /api/v1/contacts/{contactID}`

**Note:** Replace `{contactID}` with a real ID from a created contact.

```sh
curl --location 'http://127.0.0.1:8080/api/v1/contacts/0ML3H99X5QNE0E' \
--header 'X-User-ID: ABHISHEKSINHA1'
```

**Example Success Response (200 OK):**
```json
{
    "id": "0ML3H99X5QNE0E",
    "name": "Abhi Sinha",
    "email": "abhisinha@example.com",
    "phone": "+911234567890",
    "owner_id": "ABHISHEKSINHA1",
    "owner_name": "Abhishek Sinha",
    "created_at": 1769934423,
    "updated_at": 1769934423
}
```

---

### 3. List Contacts

Retrieves a paginated list of contacts for the owner specified in `X-User-ID`.

**Endpoint:** `GET /api/v1/contacts`

**Note:** Supports `page` and `page_size` query parameters.

```sh
curl --location 'http://127.0.0.1:8080/api/v1/contacts?page=1&page_size=10' \
--header 'X-User-ID: ABHISHEKSINHA1'
```

**Example Success Response (200 OK):**
```json
{
    "data": [
        {
            "id": "0ML3H99X5QNE0E",
            "name": "Abhi Sinha",
            "email": "abhisinha@example.com",
            "phone": "+911234567890",
            "owner_id": "ABHISHEKSINHA1",
            "owner_name": "Abhishek Sinha",
            "created_at": 1769934423,
            "updated_at": 1769934423
        }
    ],
    "total_count": 1
}
```

---

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

**Example Success Response (200 OK):**
The full, updated contact object is returned.

---

### 5. Delete a Contact

Deletes a contact by its ID.

**Endpoint:** `DELETE /api/v1/contacts/{contactID}`

```sh
curl --location --request DELETE 'http://127.0.0.1:8080/api/v1/contacts/0ML3H99X5QNE0E' \
--header 'X-User-ID: ABHISHEKSINHA1'
```

**Example Success Response:** `204 No Content`

---

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

**Example Success Response (201 Created):**
```json
{
    "id": "PQR9S8T7U6V5W4",
    "content": "Called the client, scheduled a follow-up for next week.",
    "author_id": "ABHISHEKSINHA1",
    "created_at": 1769935000
}
```

---

### 7. List Notes for a Contact

Retrieves all notes associated with a specific contact.

**Endpoint:** `GET /api/v1/contacts/{contactID}/notes`

```sh
curl --location 'http://127.0.0.1:8080/api/v1/contacts/0ML3H99X5QNE0E/notes' \
--header 'X-User-ID: ABHISHEKSINHA1'
```

**Example Success Response (200 OK):**
```json
{
    "data": [
        {
            "id": "PQR9S8T7U6V5W4",
            "content": "Called the client, scheduled a follow-up for next week.",
            "author_id": "ABHISHEKSINHA1",
            "created_at": 1769935000
        }
    ],
    "total_count": 1
}
```
