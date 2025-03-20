## 📖 Translations
- [Read in Russian](/README_RU.md)

---

<h3 align="center">
  <div align="center">
    <h1>Task Scheduler</h1>
  </div>
  <a href="https://github.com/goroutiner/task_scheduler">
    <img src="https://i.pinimg.com/736x/3e/5b/3a/3e5b3a55a757aa664704f6f33f2c1c4b.jpg" width="600" height="400"/>
  </a>
</h3>

---

## 📋 Project Description

**Task Scheduler** is a simple and intuitive web application for task planning, designed to help users organize their daily tasks. With this application, you can:

- Create and edit tasks.
- Set recurrence cycles and deadlines.
- Change task statuses.
- Organize tasks by due date.
- Manage tasks through a user-friendly interface.

---

## What is implemented in the application?

- ✔️ Task creation functionality.
- ✔️ Ability to edit tasks.
- ✔️ Database integration for task storage.
- ✔️ Simple and attractive interface.
- ✔️ Search and delete tasks functionality.

---

### 🔧 Environment Configuration

The **environment** variables are set by default, but you can change them in the `compose.yaml` file:

- For the `golang` service:

```yaml
...
environment:
    PORT: ":7540"
    MODE: "postgres"
    DATABASE_URL: "postgres://root:password@postgres:5432/mydb?sslmode=disable"
    PASSWORD: "qwerty12345678"
...
```

If you need **SQLite** mode, specify `MODE: "sqlite"`.

- For the `postgres` service:

```yaml
...
environment:
  POSTGRES_USER: "root"
  POSTGRES_PASSWORD: "password"
  POSTGRES_DB: "mydb"
...
```

---

## ✅⭕ Running Tests

To run integration tests, execute the following command:

```sh
make unit-tests
```

---

## 🐳 Running with Docker

If you want to run the project using Docker, follow these steps:

1. Make sure Docker is installed and running.
2. Navigate to the project's root directory.
3. Build and run the application using the command:
   - By default, the application will use **PostgreSQL**. You can change this in the `compose.yaml` file.

```sh
make run
```

4. Once the application is running, you can access it in your browser at [http://localhost:7540/login.html](http://localhost:7540/login.html) (if you used a custom port, specify it).

---

## 🛠️ Technical Resources

- **Programming Language**: Go (Golang)
- **Databases**: PostgreSQL, SQLite
- **Libraries**:
  - [golang-jwt/jwt](https://github.com/golang-jwt/jwt) for JWT token handling.
  - [joho/godotenv](https://github.com/joho/godotenv) for environment variable management.
  - [jmoiron/sqlx](https://github.com/jmoiron/sqlx) for database interaction.
  - [github.com/jackc/pgx/v5/stdlib](https://github.com/jackc/pgx) and [modernc.org/sqlite](https://gitlab.com/cznic/sqlite) drivers for PostgreSQL and SQLite.
  - [stretchr/testify](https://github.com/stretchr/testify) for testing.

---

## Conclusion

Thank you for using **Task Scheduler** 🤝 The application will continue to be supported, and more features will be added in the future 💫

---