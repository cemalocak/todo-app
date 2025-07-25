# API Referances

## Endpoints

### Todo Process

#### `GET /api/todos`

List all Todos

**Response:**
```json
[
  {
    "id": 1,
    "title": "Buy milk",
    "completed": false,
    "created_at": "2024-01-24T10:00:00Z"
  }
]
```

#### `POST /api/todos`

Add new todo

**Request:**
```json
{
  "title": "Buy milk"
}
```

**Response:**
```json
{
  "id": 1,
  "title": "Buy milk",
  "completed": false,
  "created_at": "2024-01-24T10:00:00Z"
}
```

#### `PUT /api/todos/:id`

Update todo

**Request:**
```json
{
  "title": "Buy organic milk",
  "completed": true
}
```

#### `DELETE /api/todos/:id`

Delete todo

## Test Endpoints

### `POST /api/test/truncate`

Clear all records on DB (for only tests)

## Error Codes

- `200`: Success
- `400`: Invalid request
- `404`: Not found
- `500`: Server Error