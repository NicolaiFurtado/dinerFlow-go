# 🍽️ Diner Flow — Restaurant Management API (Golang + Gin)

Diner Flow is a RESTful API developed in Go using the Gin framework. It provides essential features for managing a restaurant’s backend operations including user authentication, table and item management, order handling, and shift tracking. The system is designed for extensibility and includes thorough automated testing and documentation.

---

## 🚀 Features

- ✅ JWT-based user authentication and session handling
- ✅ Table and item CRUD operations
- ✅ Open and manage orders ("tabs")
- ✅ Shift start logic with middleware support
- ✅ Daily closing functionality
- ✅ Swagger for API documentation
- ✅ Postman collection for testing endpoints
- ✅ Unit tests using `sqlmock` and `httptest`

---

## 🧰 Tech Stack

- **Go** (v1.20+)
- **Gin** (web framework)
- **MySQL** (DB)
- **sqlmock** (for mocking DB in tests)
- **Swag** (Swagger doc generator)
- **Postman** (API testing)
- **Docker** (optional, for deployment)

---

## 📦 Installation & Usage

### 1. Clone the Repository

```bash
git clone https://github.com/YOUR_USERNAME/diner-flow.git
cd diner-flow
```

### 2. Install Dependencies

```bash
go mod tidy
```

### 3. Setup Environment Variables

Create a `.env` file in the root directory with the following content (adjust as needed):

```env
DB_USER=root
DB_PASS=root
DB_HOST=localhost
DB_PORT=3306
DB_NAME=diner-flow
```

---

## 🗃️ Database Setup

You can use the provided SQL dump to initialize your test database:

```bash
mysql -u root -p diner-flow < database/dump.sql
```

Make sure the database `diner-flow` exists before running the command.

---

## 🧪 Running Tests

Run all unit tests with:

```bash
go test -v ./controllers
```

You should see `PASS` for each test, covering:

- Authentication (SignUp, Login, Logout)
- Shift management
- Items
- Tables
- Tabs (open order)
- Closing day routine

---

## 📮 Postman Collection

You can import the included Postman collection (`docs/postman_collection.json`) into Postman to test all endpoints easily.

Steps:

1. Open Postman.
2. Click **Import**.
3. Select the `postman_collection.json` file.
4. Use the environment variables as needed (e.g., for bearer token).

---

## 📘 Swagger Documentation

The API is documented using Swagger.

To generate and view it:

```bash
swag init
```

Then run your app and visit:

```
http://localhost:8080/swagger/index.html
```

---

## 📁 Project Structure

```
.
├── config/             # DB config and .env handling
├── controllers/        # HTTP handlers and business logic
├── models/             # Struct definitions (e.g., User, Tab, Item)
├── utils/              # Helpers for JWT, password hashing, etc.
├── db/migrations/      # SQL schema and migrations
├── docs/               # Swagger & Postman files
├── main.go             # Entry point
└── README.md           # This file
```

---

## 🧱 Example Endpoints

| Method | Endpoint          | Description              |
|--------|-------------------|--------------------------|
| POST   | `/signup`         | User registration        |
| POST   | `/login`          | User login               |
| POST   | `/logout`         | Invalidate session token |
| POST   | `/startShift`     | Start user shift         |
| GET    | `/items`          | List menu items          |
| POST   | `/tab`            | Open new tab             |
| GET    | `/closeDay`       | Finalize and close day   |

---

## 🚧 Roadmap

- [x] Full backend logic
- [x] Swagger generation
- [x] Postman collection
- [x] Unit testing
- [ ] Docker containerization
- [ ] GitHub Actions for CI/CD

---

## 📄 License

This project is licensed under the MIT License.

---

## 👤 Author

**Nicolai Furtado**
[GitHub](https://github.com/NicolaiFurtado)
