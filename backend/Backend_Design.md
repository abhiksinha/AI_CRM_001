# Backend Design: Predictive CRM (Go Services)

This document outlines the initial design for the Go-based backend services. It covers the required API endpoints, background workers/crons, and the initial database schema for the PostgreSQL transactional database.

---

## 1. API Endpoints

The API will be organized by resource. All endpoints are prefixed with `/api/v1`. We will start with two core services: `UserService` and `DealService`.

### UserService (Resource: `/users` and `/contacts`)

This service manages companies, individual contacts, and internal users of the CRM.

*   **`POST /api/v1/contacts`**: Create a new contact.
*   **`GET /api/v1/contacts`**: List all contacts (with pagination, filtering, and sorting).
*   **`GET /api/v1/contacts/{id}`**: Get details of a single contact.
*   **`PUT /api/v1/contacts/{id}`**: Update a contact's details.
*   **`DELETE /api/v1/contacts/{id}`**: Delete a contact.
*   **`POST /api/v1/contacts/{id}/notes`**: Add a note to a contact.
*   **`GET /api/v1/contacts/{id}/notes`**: List all notes for a contact.

*(Similar endpoints will exist for `/companies` and internal `/users`)*

### DealService (Resource: `/deals`)

This service manages sales opportunities and pipelines.

*   **`POST /api/v1/deals`**: Create a new deal.
*   **`GET /api/v1/deals`**: List all deals (with pagination and filtering by stage, owner, etc.).
*   **`GET /api/v1/deals/{id}`**: Get details of a single deal.
*   **`PUT /api/v1/deals/{id}`**: Update a deal (e.g., change its stage, value, or expected close date).
*   **`DELETE /api/v1/deals/{id}`**: Delete a deal.
*   **`POST /api/v1/deals/{id}/tasks`**: Create a task associated with a deal (e.g., "Follow up call").
*   **`GET /api/v1/deals/{id}/tasks`**: List all tasks for a deal.

---

## 2. Background Workers & Cron Jobs

These are processes that run outside of the typical request/response cycle.

1.  **Event Publisher Worker:**
    *   **Trigger:** After any significant database write (Create, Update, Delete).
    *   **Action:** Constructs a message (e.g., `contact-updated`) with the relevant data and publishes it to the Kafka event bus. This is the primary mechanism for feeding the data warehouse.
    *   **Why:** It decouples the main API from the event streaming process. If Kafka is down, it can retry without the user-facing API failing.

2.  **Task Reminder Cron Job:**
    *   **Trigger:** Runs periodically (e.g., every 5 minutes).
    *   **Action:** Scans the `tasks` table for any tasks that are due soon. For each upcoming task, it sends a notification to the assigned user (e.g., via email or a websocket push to the frontend).
    *   **Why:** Automates reminders and helps drive user activity.

3.  **Data Cleanup/Archiving Worker:**
    *   **Trigger:** Runs periodically (e.g., nightly).
    *   **Action:** Archives old, inactive records. For example, it might mark contacts with no activity in over 2 years as "archived". This keeps the primary tables smaller and faster.
    *   **Why:** Maintains database health and performance over the long term.

---

## 3. Database Tables (PostgreSQL Schema)

This is a simplified initial schema. `id` is the primary key (UUID), and `created_at`/`updated_at` are standard timestamp fields.

### `contacts`
| Column Name | Data Type | Description |
| :--- | :--- | :--- |
| `id` | UUID | Primary Key |
| `first_name` | VARCHAR(255) | Contact's first name |
| `last_name` | VARCHAR(255) | Contact's last name |
| `email` | VARCHAR(255) | Unique email address |
| `phone` | VARCHAR(50) | Phone number |
| `company_id` | UUID | Foreign key to the `companies` table |
| `owner_id` | UUID | Foreign key to the internal `users` table (who owns this contact) |
| `created_at` | TIMESTAMPTZ | |
| `updated_at` | TIMESTAMPTZ | |

### `deals`
| Column Name | Data Type | Description |
| :--- | :--- | :--- |
| `id` | UUID | Primary Key |
| `name` | VARCHAR(255) | Name of the deal (e.g., "Q3 Website Redesign") |
| `stage` | VARCHAR(100) | Current stage in the sales pipeline (e.g., "Lead", "Proposal") |
| `value` | DECIMAL | Monetary value of the deal |
| `expected_close_date` | DATE | |
| `contact_id` | UUID | Foreign key to the `contacts` table |
| `owner_id` | UUID | Foreign key to the internal `users` table |
| `created_at` | TIMESTAMPTZ | |
| `updated_at` | TIMESTAMPTZ | |

### `notes`
| Column Name | Data Type | Description |
| :--- | :--- | :--- |
| `id` | UUID | Primary Key |
| `contact_id` | UUID | Foreign key to `contacts` (the note is about this contact) |
| `deal_id` | UUID | Foreign key to `deals` (optional, if the note is deal-specific) |
| `author_id` | UUID | Foreign key to `users` (who wrote the note) |
| `content` | TEXT | The body of the note |
| `created_at` | TIMESTAMPTZ | |

### `tasks`
| Column Name | Data Type | Description |
| :--- | :--- | :--- |
| `id` | UUID | Primary Key |
| `deal_id` | UUID | Foreign key to `deals` (this task is for this deal) |
| `assigned_to_id` | UUID | Foreign key to `users` (who needs to do the task) |
| `title` | VARCHAR(255) | What needs to be done |
| `due_date` | TIMESTAMPTZ | When the task is due |
| `is_completed` | BOOLEAN | Status of the task |
| `created_at` | TIMESTAMPTZ | |
| `updated_at` | TIMESTAMPTZ | |

*(This schema would also include `companies` and `users` tables to be fully functional.)*
