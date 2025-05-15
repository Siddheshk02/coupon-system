# Coupon System MVP
This MVP supports coupon creation, validation, and application logic with caching and a modular architecture.

---

##  Features

- **Admin Coupon Creation**: Create/manage coupons with expiry, usage type, categories, min order value, discount, etc.
- **Coupon Validation**: Validate coupons against cart, order, and user constraints.
- **Persistent Storage**: Uses PostgreSQL (Dockerized).
- **Concurrency Safety**: Request-scoped context, DB-level safety.
- **Caching**: In-memory TTL cache for frequently accessed (Read-heavy) data.
- **OpenAPI Docs**: Swagger UI available.
- **Dockerized**: One-command setup for API and DB.


##  Architecture

- **Go Modules**: Modular codebase under `internal/` and `cmd/server`.
- **Repository Pattern**: All DB access via repositories.
- **Handlers**: HTTP handlers for each resource.
- **Caching**: [patrickmn/go-cache](https://github.com/patrickmn/go-cache) for TTL-based in-memory caching.
- **Database Migrations**: SQL migration files in `/migrations`.


##  API Endpoints

### Coupon Endpoints

- `POST /admin/coupons` — Create a coupon
- `GET /coupons` — List all coupons
- `GET /coupons/applicable` — Get applicable coupons for a cart
- `POST /coupons/validate` — Validate a coupon for a cart/order

### Items & Orders

- `POST /items` — Add an item
- `GET /items` — List items (with optional filter for id and/or category)
- `POST /createorder` — Place an order

### Users

- `POST /users` — User login (creates user if not exists)


##  Quick Start

1. Deployed API URL : https://coupon-system-mjzu.onrender.com/

2. Docker :
   1. **Clone the repo:**
      ```sh
      git clone https://github.com/Siddheshk02/coupon-system.git
      cd coupon-system
      ```
   2. **Build and run everything:**
      ```sh
      docker-compose up --build
      ```
   3. **API will be available at:**  
      [http://localhost:8080](http://localhost:8080)

##  Example Requests

**Get All Coupons**
```sh
curl http://localhost:8080/coupons
```

**Validate Coupon**
```sh
curl -X POST http://localhost:8080/coupons/validate \
  -H "Content-Type: application/json" \
  -d '{"coupon_code":"SAVE20","cart_items":[{"id":"12","category":"painkiller"}],"order_total":700,"timestamp":"2025-05-05T15:00:00Z"}'
```


##  Swagger/OpenAPI Docs

- https://app.swaggerhub.com/apis/KHANDAGALESID02_1/coupon-system_api/1.0

