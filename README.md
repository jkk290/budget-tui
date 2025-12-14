# Budget-TUI

## Motivation
A budgeting application built in Go, combining a REST API backend with a terminal-based user interface. This project serves as a hands-on way to practice designing API services, building API clients, and working with relational databases.

Current Implementation
Uses a Go REST API server, an API client, and a PostgreSQL database to explore full-stack architecture patterns and data flow between backend and TUI layers.

Next Implementation
Plans to switch to a SQLite database and refactor the TUI to perform CRUD operations directly against the database—an exploration of a simplified, single-application architecture without a separate API layer.

## Quickstart

## Usage
## API Documentation

### Authentication

Most endpoints require JWT authentication. Include the token in the `Authorization` header:

```
Authorization: Bearer <your_jwt_token>
```

Tokens expire after 1 hour and can be obtained via the login endpoint.

---

## Endpoints

### Authentication & Users

#### `POST /users`
Create a new user account.

**Authentication:** Not required

**Request:**
```json
{
  "username": "john_doe",
  "password": "secure_password123"
}
```

**Response:** `201 Created`
```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "created_at": "2025-12-10T14:30:00Z",
  "updated_at": "2025-12-10T14:30:00Z",
  "username": "john_doe"
}
```

---

#### `POST /login`
Authenticate and receive a JWT token.

**Authentication:** Not required

**Request:**
```json
{
  "username": "john_doe",
  "password": "secure_password123"
}
```

**Response:** `200 OK`
```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "created_at": "2025-12-10T14:30:00Z",
  "updated_at": "2025-12-10T14:30:00Z",
  "username": "john_doe",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

---

### Accounts

#### `GET /accounts`
Get all accounts for the authenticated user with current balances.

**Authentication:** Required

**Response:** `200 OK`
```json
[
  {
    "id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
    "account_name": "Chase Checking",
    "account_type": "checking",
    "created_at": "2025-12-01T10:00:00Z",
    "updated_at": "2025-12-01T10:00:00Z",
    "user_id": "123e4567-e89b-12d3-a456-426614174000",
    "account_balance": "2543.67"
  }
]
```

---

#### `POST /accounts`
Create a new account with an initial balance.

**Authentication:** Required

**Request:**
```json
{
  "account_name": "Chase Checking",
  "account_type": "checking",
  "initial_balance": "1000.00"
}
```

**Response:** `201 Created`
```json
{
  "id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
  "account_name": "Chase Checking",
  "account_type": "checking",
  "created_at": "2025-12-10T14:30:00Z",
  "updated_at": "2025-12-10T14:30:00Z",
  "user_id": "123e4567-e89b-12d3-a456-426614174000"
}
```

---

#### `PUT /accounts/{accountID}`
Update an account's name.

**Authentication:** Required

**Request:**
```json
{
  "account_name": "Chase Personal Checking"
}
```

**Response:** `200 OK` (returns updated account object)

---

#### `DELETE /accounts/{accountID}`
Delete an account.

**Authentication:** Required

**Response:** `204 No Content`

---

#### `GET /accounts/{accountID}/transactions`
Get all transactions for a specific account.

**Authentication:** Required

**Response:** `200 OK` (returns array of transaction objects)

---

### Transactions

#### `GET /transactions`
Get all transactions for the authenticated user across all accounts.

**Authentication:** Required

**Response:** `200 OK`
```json
[
  {
    "id": "t1a2b3c4-d5e6-7890-abcd-ef1234567890",
    "amount": "-45.32",
    "tx_description": "Grocery Store",
    "tx_date": "2025-12-05T15:30:00Z",
    "created_at": "2025-12-05T15:30:00Z",
    "updated_at": "2025-12-05T15:30:00Z",
    "posted": true,
    "account_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
    "category_id": "c1d2e3f4-a5b6-7890-cdef-123456789012",
    "account_name": "Chase Checking",
    "category_name": "Groceries"
  }
]
```

---

#### `POST /transactions`
Create a new transaction.

**Authentication:** Required

**Request:**
```json
{
  "amount": "-45.32",
  "tx_description": "Grocery Store",
  "tx_date": "2025-12-05T15:30:00Z",
  "posted": true,
  "account_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
  "category_id": "c1d2e3f4-a5b6-7890-cdef-123456789012"
}
```

**Notes:**
- `amount`: Negative for expenses, positive for income
- `category_id`: Optional, use `00000000-0000-0000-0000-000000000000` for uncategorized

**Response:** `201 Created` (returns transaction object)

---

#### `PUT /transactions/{transactionID}`
Update a transaction. Similar to `POST /transactions`.

**Authentication:** Required

**Response:** `200 OK` (returns updated transaction object)

---

#### `DELETE /transactions/{transactionID}`
Delete a transaction.

**Authentication:** Required

**Response:** `204 No Content`

---

### Categories

#### `GET /categories`
Get all spending categories for the authenticated user.

**Authentication:** Required

**Response:** `200 OK`
```json
[
  {
    "id": "c1d2e3f4-a5b6-7890-cdef-123456789012",
    "category_name": "Groceries",
    "created_at": "2025-12-01T10:00:00Z",
    "updated_at": "2025-12-01T10:00:00Z",
    "budget": "500.00",
    "user_id": "123e4567-e89b-12d3-a456-426614174000",
    "group_id": "g1h2i3j4-k5l6-7890-ghij-123456789012",
    "group_name": "Essential Expenses"
  }
]
```

---

#### `POST /categories`
Create a new spending category.

**Authentication:** Required

**Request:**
```json
{
  "category_name": "Groceries",
  "budget": "500.00",
  "group_id": "g1h2i3j4-k5l6-7890-ghij-123456789012"
}
```

**Notes:**
- `budget` and `group_id` are optional

**Response:** `201 Created` (returns category object)

---

#### `PUT /categories/{categoryID}`
Update a category. Similar to `POST /categories`.

**Authentication:** Required

**Response:** `200 OK` (returns updated category object)

---

#### `DELETE /categories/{categoryID}`
Delete a category.

**Authentication:** Required

**Response:** `204 No Content`

---

#### `GET /categories/{categoryID}/transactions`
Get all transactions for a specific category.

**Authentication:** Required

**Response:** `200 OK` (returns array of transaction objects)

---

### Groups

#### `GET /groups`
Get all category groups for the authenticated user.

**Authentication:** Required

**Response:** `200 OK`
```json
[
  {
    "id": "g1h2i3j4-k5l6-7890-ghij-123456789012",
    "group_name": "Essential Expenses",
    "created_at": "2025-12-01T10:00:00Z",
    "updated_at": "2025-12-01T10:00:00Z",
    "user_id": "123e4567-e89b-12d3-a456-426614174000"
  }
]
```

---

#### `POST /groups`
Create a new category group.

**Authentication:** Required

**Request:**
```json
{
  "group_name": "Essential Expenses"
}
```

**Response:** `201 Created` (returns group object)

---

#### `PUT /groups/{groupID}`
Update a group. Similar to `POST /groups`.

**Authentication:** Required

**Response:** `200 OK` (returns updated group object)

---

#### `DELETE /groups/{groupID}`
Delete a group.

**Authentication:** Required

**Response:** `204 No Content`

---

### Budget Overview

#### `GET /budget`
Get a comprehensive budget overview for the current month.

**Authentication:** Required

**Response:** `200 OK`
```json
{
  "start_date": "2025-12-01T00:00:00Z",
  "end_date": "2026-01-01T00:00:00Z",
  "groups": [
    {
      "group_id": "g1h2i3j4-k5l6-7890-ghij-123456789012",
      "group_name": "Essential Expenses",
      "categories": [
        {
          "category_id": "c1d2e3f4-a5b6-7890-cdef-123456789012",
          "category_name": "Groceries",
          "budget": "500.00",
          "total_spent": "342.67",
          "remaining": "157.33",
          "is_overspent": false
        }
      ],
      "total_budget": "700.00",
      "total_spent": "558.10",
      "total_remaining": "141.90"
    }
  ],
  "ungrouped_categories": [
    {
      "category_id": "c3f4a5b6-c7d8-9012-efgh-123456789012",
      "category_name": "Entertainment",
      "budget": "150.00",
      "total_spent": "89.50",
      "remaining": "60.50",
      "is_overspent": false
    }
  ],
  "grand_total_budget": "850.00",
  "grand_total_spent": "647.60",
  "grand_total_remaining": "202.40"
}
```

---

## Data Types

- **UUID**: Standard UUID format (e.g., `123e4567-e89b-12d3-a456-426614174000`)
- **Decimal**: String representation of decimal numbers (e.g., `"123.45"`)
- **Timestamp**: ISO 8601 format in UTC (e.g., `"2025-12-10T14:30:00Z"`)

---

## Error Responses

All errors return JSON with an `error` field:

```json
{
  "error": "Human-readable error message"
}
```

**Common Status Codes:**
- `400` - Invalid request parameters
- `401` - Missing or invalid authentication
- `403` - Insufficient permissions
- `404` - Resource not found
- `500` - Server error

## Contributing
