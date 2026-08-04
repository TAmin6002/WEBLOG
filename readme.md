<div align="center">

# 📝 Weblog

**A lightweight multi-user blogging platform** — built with Go, Echo & PostgreSQL

[![Go](https://img.shields.io/badge/Go-1.26-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://go.dev)
[![Echo](https://img.shields.io/badge/Echo-v4-3E863D?style=for-the-badge&logo=go&logoColor=white)](https://echo.labstack.com/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15-336791?style=for-the-badge&logo=postgresql&logoColor=white)](https://www.postgresql.org/)
[![Docker](https://img.shields.io/badge/Docker-ready-2496ED?style=for-the-badge&logo=docker&logoColor=white)](https://www.docker.com/)

[Features](#-features) • [Tech Stack](#-tech-stack) • [Quick Start](#-quick-start) • [Routes](#%EF%B8%8F-routes) • [Live Demo](#-live-demo)

</div>

---

## 🚀 Overview

**Weblog** lets anyone sign up, publish posts, and control exactly who sees them. Every post can be **public** or **private** — and private posts can be shared with hand-picked usernames. Comments, ownership rules, and image uploads round out a small but complete blogging experience, all served with clean server-side rendered pages.

---

## ✨ Features

<table>
<tr>
<td width="50%" valign="top">

### 🔐 Accounts
- Sign up / log in with a unique username
- Passwords hashed with **bcrypt**
- Cookie-based sessions

### 📰 Posts
- Title, content, optional image
- Public or private visibility
- Only the author can delete (no editing, by design)

</td>
<td width="50%" valign="top">

### 🔒 Private Sharing
- Share a private post with specific usernames
- Only the author + shared users can view it
- Enforced at the database query level

### 💬 Comments
- Logged-in users only
- Visible on any post they're allowed to view
- Text + author username

</td>
</tr>
</table>

---

## 🧱 Tech Stack

<div align="center">

| Layer | Technology |
|:---:|:---:|
| Language | ![Go](https://img.shields.io/badge/-Go%201.26-00ADD8?logo=go&logoColor=white) |
| Web Framework | ![Echo](https://img.shields.io/badge/-Echo%20v4-3E863D) |
| Database | ![PostgreSQL](https://img.shields.io/badge/-PostgreSQL%2015-336791?logo=postgresql&logoColor=white) |
| Rendering | `html/template` (server-side) |
| Auth | Cookie sessions + bcrypt |
| Containerization | ![Docker](https://img.shields.io/badge/-Docker-2496ED?logo=docker&logoColor=white) |

</div>

---

## 🚀 Quick Start

### Option A — Docker (recommended)

```bash
git clone https://github.com/TAmin6002/WEBLOG.git
cd WEBLOG
docker compose up --build
```

Open **http://localhost:8080** — the database schema applies automatically on first boot.

### Option B — Local Go + Dockerized DB

```bash
docker compose up -d db
export DB_CONN="postgres://user:password@localhost:5432/weblog_db?sslmode=disable"
go run .
```

---

## ⚙️ Configuration

| Variable | Purpose | Default |
|---|---|---|
| `DATABASE_URL` | Full Postgres connection string (used by most hosts) | — |
| `DB_CONN` | Fallback connection string | — |
| `PORT` | Server port | `8080` |

If neither `DATABASE_URL` nor `DB_CONN` is set:
```
postgres://user:password@localhost:5432/weblog_db?sslmode=disable
```

---

## 📂 Project Structure

```
weblog/
.
├── docker-compose.yml
├── dockerfile
├── go.mod
├── go.sum
├── handler
│   ├── auth.go
│   ├── board.go
│   └── comment.go
├── main.go
├── middleware
│   └── auth.go
├── model
│   ├── board.go
│   ├── comment.go
│   └── user.go
├── readme.md
├── repo
│   ├── board_repo.go
│   ├── comment_repo.go
│   └── user_repo.go
├── service
│   ├── auth_service.go
│   ├── board_service.go
│   └── comment_service.go
├── sql
│   └── schema.sql
├── templates
│   ├── board.html
│   ├── index.html
│   ├── login.html
│   └── signup.html


```

---

## 🗺️ Routes

| Method | Path | Description | Auth |
|:---:|---|---|:---:|
| `GET` | `/signup` | Signup form | — |
| `POST` | `/signup` | Create account | — |
| `GET` | `/login` | Login form | — |
| `POST` | `/login` | Authenticate | — |
| `POST` | `/logout` | Log out | 🔒 |
| `GET` | `/` | Feed of visible posts | 🔒 |
| `POST` | `/weblog` | Create a new post | 🔒 |
| `GET` | `/weblog/:id` | View post + comments | 🔒 |
| `POST` | `/weblog/:id/delete` | Delete your own post | 🔒 |
| `POST` | `/weblog/:id/comments` | Add a comment | 🔒 |

---

## 🔒 Privacy Model

```
             ┌───────────────┐
             │   New Post    │
             └───────┬───────┘
                      │
              is_private?
              ┌───────┴───────┐
             No               Yes
              │                │
       visible to        visible to:
       everyone           • author
                           • shared usernames only
```

Visibility is enforced directly in the SQL query — private content never leaks into someone else's feed or detail page.

---

## 🌐 Live Demo

> 🔗 **[weblog_app](https://weblog-production-c9e6.up.railway.app/)**

<div align="center">

made with ☕ and Go

</div>
