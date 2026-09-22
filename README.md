# KapoorTraderLoyalty

A robust backend service for managing trader and customer loyalty programs. Built with Go, this application utilizes a clean architecture pattern to handle shopkeepers, customers, and their relational loyalty data.

## 🛠 Tech Stack

* **Language:** Go
* **Database:** PostgreSQL
* **Containerization:** Docker & Docker Compose
* **Architecture:** Domain-Driven Design / Clean Architecture[cite: 1]
* **Frontend:** Go HTML Templates (`web/templates`)[cite: 1]

## 📂 Project Structure

The project follows standard Go layout conventions to ensure separation of concerns:[cite: 1]

* `cmd/server/`: Contains the main entry point (`main.go`) to bootstrap and run the HTTP server.[cite: 1]
* `internal/`: Encapsulates the core application logic.[cite: 1]
  * `domain/`: Core business models (`customer.go`, `shop.go`, `shopkeeper.go`, `dashboard.go`, `shop_customer.go`).[cite: 1]
  * `usecase/`: Application-specific business rules and service layers (e.g., `customer_service.go`).[cite: 1]
  * `repository/`: Interfaces for data storage mechanisms.[cite: 1]
  * `infrastructure/postgres/`: Concrete database implementations and DB connections (`db.go`, `customer_repository.go`).[cite: 1]
  * `delivery/http/`: HTTP handlers and routing (`customer_handler.go`, `shopkeeper_handler.go`).[cite: 1]
* `migrations/`: SQL scripts for database initialization (`001_init.sql`) and seeding (`002_seed.sql`).[cite: 1]
* `web/templates/`: HTML views for the application UI (`customer_login.html`, `dashboard.html`).[cite: 1]
* `assets.go`: Static asset management.[cite: 1]

## 🚀 Getting Started

### Prerequisites
* [Go](https://golang.org/doc/install) (1.19+ recommended)
* [Docker](https://docs.docker.com/get-docker/) & Docker Compose[cite: 1]

### Running with Docker (Recommended)

The easiest way to get the application and its database running is via Docker Compose:[cite: 1]

1. Clone the repository.
2. Build and start the containers:
   ```bash
   docker-compose up --build


Public Endpoints
GET / : Renders the landing page (index.html)[cite: 1].

Authentication Endpoints
GET /register : Renders the user registration form (register.html)[cite: 1].

POST /register : Processes new user account creation.

GET /login : Renders the authentication form (login.html)[cite: 1].

POST /login : Authenticates a user and establishes a session.

POST /logout : Terminates the current user session.

Authenticated/Protected Endpoints
GET /dashboard : Renders the main user control panel (dashboard.html)[cite: 1].

GET /product/edit/{id} : Renders the form to modify an existing farm product (edit_product.html)[cite: 1].

POST /product/edit/{id} : Submits and processes updates to a specific product.

Future Enhancements
Implementation of a RESTful JSON API alongside the server-rendered templates.

Integration of unit testing and mocking for usecases and repositories.

CI/CD pipeline configuration for automated deployments.

Developed for the Farm Connect MVP.
