# flight-prices-api

## Architecture

This project is a modular, concurrent Go API for searching and aggregating flight prices from multiple providers. The architecture is composed of:

- **cmd/**: Entry point for the application (main.go).
- **internal/controllers/**: HTTP handlers for authentication and flight search endpoints.
- **internal/services/**: Business logic, including itinerary fetching and authentication.
- **internal/services/clients/**: API clients for external flight data providers (Amadeus, FlightAPI, SerpAPI).
- **internal/repository/**: Data access layers for MongoDB (users) and Redis (itinerary cache).
- **internal/models/**: Data models shared across the application.
- **static/**: Simple HTML UI for manual testing.
- **mocks/**: Auto-generated mocks for unit and integration testing.
- **pkg/config/**: Application configuration loader.

The API is designed to:
- Authenticate users (JWT-based)
- Fetch and aggregate flight itineraries from multiple providers concurrently
- Cache search results in Redis
- Support local development and integration testing with Docker Compose and mock servers

---

## System Requirements

- **Go** 1.23+
- **Docker** 27.4+ and **Docker Compose**
- **GNU Make** (optional, for convenience)

---

## How to Run the Project

1. **Clone the repository:**
   ```sh
   git clone <repo-url>
   cd flight-prices-api
   ```

2. **Copy and edit the environment variables:**
   ```sh
   cp .env.example .env
   # Edit .env as needed (API keys, secrets, etc)
   ```

3. **Start all services (API, MongoDB, Redis, mock servers):**
   ```sh
   make run-local-build
   ```
   or
   ```
   make run-with-mocks-build
   ```
   if you want to have the API's mock servers running.

   This will start:
   - The Go API server (on port 8080)
   - MongoDB (with authentication)
   - Redis (with authentication)

4. **Open the UI:**
   - Visit [http://localhost:8080/](http://localhost:8080/) for the simple HTML UI.

5. **Run integration tests:**
   ```sh
   make run-integration-tests
   ```

---

## How to Call the Endpoints

### Register
```sh
curl -X POST http://localhost:8080/register \
  -H 'Content-Type: application/json' \
  -d '{"username":"user","password":"pass"}'
```

### Login
```sh
curl -X POST http://localhost:8080/login \
  -H 'Content-Type: application/json' \
  -d '{"username":"user","password":"pass"}'
```
Returns a JWT token.

### Search Itineraries
```sh
curl --location 'localhost:8080/flights/search?origin=SDQ&destination=JFK&date=2025-04-27' \
--header 'Authorization: Bearer {authentication-token}'
```

---

## Which APIs is the Project Consulting?

- **Amadeus**: For flight offers and pricing
- **FlightAPI**: For additional flight data
- **SerpAPI (Google Flights engine)**: For scraping Google Flights results

The project queries all three APIs concurrently and aggregates the results.

---

## How the Project is Mocking Those Calls for Testing Locally

- **Mock servers** are defined in the `mocks/mockserver/` directory (JSON files for each provider).
- This allows for deterministic, fast, and cost-free integration and end-to-end testing.
- Unit tests use GoMock-generated mocks (see `mocks/` directory) for all external dependencies.

---

## How HTTPS/TLS Should Be Configured for Production

- **Do not use the development CORS or static file settings in production.**
- Obtain a valid TLS certificate (e.g., via Let's Encrypt).
- **Never expose the Go server directly to the internet without TLS.**

---

