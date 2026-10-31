# WhereDidMyMoneyGO

## A simple RESTful API in Go for tracking personal expenses.

### Features

    👤 User registration and management

    🔐 JWT-based login and protected routes

    💸 Movement (transaction) tracking and a dynamic monthly balance based on those movements

    📊 Planned: Export movements with balance to an Excel spreadsheet

## 📦 Prerequisites

- Docker
- Docker Compose
- Git

## Setup and Installation

Clone the repository

```bash
git clone https://github.com/pedro-scarelli/wheredidmymoneygo.git
```

Navigate to the project root

```bash
cd wheredidmymoneygo
```

## Running the Server

```bash
docker compose up -d
```

API Endpoints

Base URL: http://localhost:3000

## Register account

POST /account HTTP/1.1
Content-Type: application/json

```json
{
    "firstName": "{FIRST_NAME}",
    "lastName": "{LAST_NAME}",
    "cpf": "{CPF}",
    "email": "{EMAIL}",
    "password": "{PASSWORD}"
}
```
## Get account by ID

GET /account/{ID} HTTP/1.1
Authorization: Bearer {USER_TOKEN}

## Update account

PATCH /account HTTP/1.1
Content-Type: application/json
Authorization: Bearer {USER_TOKEN}

```json
{
    "id": "{ID}",
    "firstName": "{FIRST_NAME}",
    "lastName": "{LAST_NAME}",
    "password": "{PASSWORD}"
}
```
## Delete account by ID

DELETE /account/{ID} HTTP/1.1
Authorization: Bearer {USER_TOKEN}

## Authentication

Method Route Auth Description
POST /login HTTP/1.1
Content-Type: application/json

```json
{
    "email": "{EMAIL}",
    "password": "{PASSWORD}"
}
```
Returns an authentication token

## Movement

POST /account/movement'
Authorization: Bearer {USER_TOKEN}

```json
{
    "value": {VALUE_WITH_2_DECIMAL_PLACES},
    "type": "{CREDITO_OR_DEBITO}"
    "recurrence": {RECURRENCE_COUNT - 1 TO 12},
    "dueDate": "{DATE_IN_FORMAT 2025-05-27T00:00:00Z}",
    "description": "{DESCRIPTION}"
}
```
`type` accepts `CREDITO` (credit) or `DEBITO` (debit).
