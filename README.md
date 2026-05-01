# Overtime Management System (Overtime-MS)

A robust, full-stack application designed to streamline the process of requesting, checking, and approving overtime hours within an organization. Featuring role-based access control (RBAC), departmental management, and a structured approval workflow.

## 🚀 Key Features

- **Role-Based Access Control (RBAC)**: Distinct permissions for Applicants, Checkers, Approvers, Finance, and Administrators.
- **Structured Workflow**:
  - **Create**: Employees can submit overtime requests.
  - **Check**: Supervisors verify the validity of requests.
  - **Approve**: Final management approval before processing.
  - **Reject**: Capability to reject requests at any stage with feedback.
- **Admin Dashboard**: Manage users, departments, and monitor all overtime activities.
- **Audit Logging**: Backend logging for tracking system interactions.
- **Responsive UI**: Modern, clean interface built with React and Vite.

## 🛠️ Tech Stack

### Backend
- **Language**: Go (Golang)
- **Framework**: [Gin Gonic](https://gin-gonic.com/)
- **Database**: PostgreSQL (via `pgx/v5`)
- **Authentication**: JWT (JSON Web Tokens)
- **Hot Reload**: [Air](https://github.com/cosmtrek/air)

### Frontend
- **Framework**: React 19
- **Build Tool**: Vite
- **Language**: TypeScript
- **Styling**: Tailwind CSS
- **Routing**: React Router 7
- **HTTP Client**: Axios

### Infrastructure
- **Containerization**: Docker & Docker Compose (for PostgreSQL)

## 📦 Project Structure

```text
overtime-ms/
├── backend/            # Go source code
│   ├── cmd/            # Entry points
│   ├── internal/       # Core logic (handlers, models, middleware, config)
│   └── main.go         # Application entry
├── frontend/           # React source code
│   ├── src/            # Components, pages, hooks, services
│   └── package.json    # Dependencies and scripts
├── docker-compose.yml  # Database infrastructure
└── .env                # Environment variables (Shared/Root)
```

## ⚙️ Setup & Installation

### Prerequisites
- [Go](https://go.dev/) (v1.25+)
- [Node.js](https://nodejs.org/) (v20+)
- [Docker](https://www.docker.com/) & Docker Compose
- [Air](https://github.com/cosmtrek/air) (Optional, for backend hot reload)

### 1. Environment Configuration
Create a `.env` file in the root directory (and optionally in `backend/` and `frontend/` if specific overrides are needed):

```env
# Database
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_DB=overtime_db
DB_HOST=localhost
DB_PORT=5432

# Backend
JWT_SECRET=your_super_secret_key
PORT=8080
CORS_ORIGIN=http://localhost:5173

# Frontend
VITE_API_BASE_URL=http://localhost:8080/api/v1
```

### 2. Database Setup
Start the PostgreSQL container:
```bash
docker-compose up -d
```

### 3. Backend Setup
```bash
cd backend
go mod download
# Run with hot reload (recommended)
air
# OR run normally
go run main.go
```

### 4. Frontend Setup
```bash
cd frontend
npm install
npm run dev
```

The application will be available at `http://localhost:5173`.

## 👥 User Roles

| Role | Description |
| :--- | :--- |
| **Applicant** | Can create and view their own overtime requests. |
| **Checker** | Can view pending requests and mark them as "Checked". |
| **Approver** | Can view checked requests and "Approve" or "Reject" them. |
| **Finance** | Can view all approved requests for payroll processing. |
| **Admin** | Full system access: User/Department management and all overtime actions. |

## 📝 License
This project is private and intended for internal use.
