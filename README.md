# POS API (Golang)

Simple REST API for products and categories built with the Go standard library.
Includes in-memory storage and an interactive API docs page powered by Scalar.


## API Docs
- Scalar UI: `http://localhost:8081/docs`
- OpenAPI spec: `http://localhost:8081/openapi.json`

## Endpoints
### Products
- `GET /api/products` (query: `limit`, `offset`)
- `POST /api/products`
- `GET /api/products/{id}`
- `PUT /api/products/{id}`
- `DELETE /api/products/{id}`

### Categories
- `GET /api/categories` (query: `limit`, `offset`)
- `POST /api/categories`
- `GET /api/categories/{id}`
- `PUT /api/categories/{id}`
- `DELETE /api/categories/{id}`

### Health
- `GET /health`
