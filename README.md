# URL Shortener

### [Live Demo](https://url-shortener-uy9w.onrender.com/)

## Features

- Shorten long URLs with a custom or auto-generated code
- Redirect short code to original URL
- PostgreSQL-based persistent storage
- Swagger (OpenAPI) documentation at /swagger/index.html
- Frontend form (served via /public/index.html)

## Technologies Used

- Go + Fiber web framework
- PostgreSQL
- Swagger (Swaggo) for API docs
- Air for live reloading
- HTML + JS frontend

## API Endpoints

1. POST `/` : Create a new shortened URL.

Request Body:

```json
{
  "url": "https://example.com",
  "code": "custom123" // optional
}
```

Response:

```json
{
  "short_url": "/custom123",
  "original_url": "https://example.com"
}
```

2. GET /:code : Redirects to the original URL associated with the given code.

Example:
GET `/custom123` → redirects to `https://example.com`

## Database & Intergration

#### Database `PostgreSQL`

The service uses PostgreSQL as the persistent storage engine for mapping short codes to original URLs.

Schema:

```sql
CREATE TABLE url_mappings (
	code TEXT PRIMARY KEY,
	url TEXT NOT NULL
);
```

code: The unique short identifier used in the URL (e.g., abc123)\
url: The original long URL to redirect to

#### Integration

The database connection is established using a single connection URL string (DSN format):

```plaintext
PG_URL=postgres://<user>:<password>@<host>:<port>/<dbname>?sslmode=require
```

Example:

```plaintext
PG_URL=postgres://postgres:admin0@localhost:5432/postgres?sslmode=require
```

The Go code uses: `database/sql` standard library\
PostgreSQL driver: `github.com/lib/pq`\
<br>
Connection logic is modularized in `internal/database/db.go`:

```go
db, err := sql.Open("postgres", os.Getenv("PG_URL"))
```

It automatically pings the database on startup to verify connection health.

## 🛠️ How to Run the Server

1. Clone and install dependencies

```bash
git clone https://github.com/kpriyanshu2003/url-shortener
cd url-shortener
go mod tidy
```

2. Set up your .env or environment variables

```bash
export PG_URL=postgres://<user>:<password>@<host>:<port>/<dbname>?sslmode=require
```

3. Start PostgreSQL and create the url_mappings table

```bash
psql -U postgres -d postgres -c "CREATE TABLE IF NOT EXISTS url_mappings (code TEXT PRIMARY KEY, url TEXT NOT NULL);"
```

4. Start the server

```bash
go run cmd/main.go
```

5. Open your browser and go to:

```
http://localhost:3300/
```

## How to Run Frontend (Optional)

Your index.html and script.js are inside /public/ directory. Just run the server and open: `http://localhost:3300/`

## Swagger Documentation

Accessible at: `http://localhost:3300/swagger/index.html`

Auto-generated using `swaggo/swag`. To regenerate docs:

```bash
go install github.com/swaggo/swag/cmd/swag@latest
swag init -g cmd/main.go --output docs
```

## Sample cURL Requests

POST Shorten URL:

```bash
curl -X POST http://localhost:3300/shorten \
-H "Content-Type: application/json" \
-d '{"url": "https://example.com", "code": "hello123"}'
```

Redirect

```bash
curl -v http://localhost:3300/hello123
```

## Contributing

Pull requests and issues welcome!
Want to extend this with analytics, auth, expiry, or rate limiting? Go ahead!
