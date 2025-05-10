# Campus Lost and Found Application

A comprehensive lost and found web application for campus use with Go backend, featuring authentication, item management, communication system, and role-based permissions for students, guards, and administrators.

## Features

- User authentication with JWT
- Role-based permissions (students, guards, admins)
- Item management (lost and found items)
- Internal messaging system
- Admin dashboard
- Image upload support

## Tech Stack

- **Backend**: Go with Fiber framework
- **Database**: PostgreSQL (Neon Tech)
- **Authentication**: JWT
- **Frontend**: React.js (not implemented yet)

## Setup

### Prerequisites

- Go 1.19+
- PostgreSQL or Neon Tech account
- Git

### Environment Variables

Create a `.env` file in the root directory with the following variables:

```
DATABASE_URL=postgresql://username:password@host:port/dbname?sslmode=require
JWT_SECRET=your_secret_key
PORT=3000
```

For Neon Tech database:
```
DATABASE_URL=postgresql://neondb_owner:your_password@your-instance-pooler.region.aws.neon.tech/neondb?sslmode=require
```

### Installation

1. Clone the repository
```bash
git clone https://github.com/yourusername/campus-lostandfound.git
cd campus-lostandfound
```

2. Install dependencies
```bash
go mod download
```

3. Run the application
```bash
go run cmd/api/main.go
```

The API will be available at http://localhost:3000

## Testing

### Run API Tests

The project includes automated test scripts that verify the functionality of all API endpoints.

```bash
# Make scripts executable
chmod +x test_api.sh test_register.sh

# Run the API tests
./test_api.sh
```

Test results will be available in the `tests` directory:
- `success.txt` - List of successful tests
- `error.txt` - List of failed tests
- `success_details.txt` - Detailed responses for successful tests
- `error_details.txt` - Detailed responses for failed tests

### Setting up an Admin User

To test admin endpoints, you need to set up an admin user in the database:

1. First register a user through the API
2. Update the user's role to 'admin' using the provided script:

```bash
# Register an admin user
curl -X POST -H "Content-Type: application/json" -d '{"username":"admin","email":"admin@example.com","password":"admin123","first_name":"Admin","last_name":"User","phone":"9876543210"}' http://localhost:3000/api/v1/auth/register

# Set admin role (requires PostgreSQL CLI)
./set_admin_role.sh admin
```

## API Documentation

### Authentication Endpoints

- `POST /api/v1/auth/register` - Register a new user
- `POST /api/v1/auth/login` - Login and get JWT token

### Item Endpoints

- `GET /api/v1/items` - Get all items (public)
- `GET /api/v1/items/search?q=keyword` - Search items (public)
- `GET /api/v1/items/:id` - Get item details (public)
- `POST /api/v1/items/lost` - Report lost item (authenticated)
- `POST /api/v1/items/found` - Report found item (authenticated)
- `PUT /api/v1/items/:id/status` - Update item status (authenticated)

### User Endpoints

- `GET /api/v1/user/profile` - Get user profile (authenticated)
- `PUT /api/v1/user/profile` - Update user profile (authenticated)
- `PUT /api/v1/user/password` - Change password (authenticated)
- `GET /api/v1/user/items` - Get user's items (authenticated)

### Messaging Endpoints

- `GET /api/v1/user/messages/unread` - Get unread message count (authenticated)
- `GET /api/v1/user/messages/conversations` - Get user's conversations (authenticated)
- `GET /api/v1/user/messages/:user_id` - Get messages with specific user (authenticated)
- `POST /api/v1/user/messages` - Send message (authenticated)

### Admin Endpoints

- `GET /api/v1/admin/users` - Get all users (admin only)
- `PUT /api/v1/admin/users/:id/role` - Update user role (admin only)
- `GET /api/v1/admin/reports` - Get all reports (admin only)
- `PUT /api/v1/admin/reports/:id/status` - Update report status (admin only)
- `GET /api/v1/admin/stats` - Get system stats (admin only)

## Deployment

For deployment instructions, see the [DEPLOYMENT.md](DEPLOYMENT.md) file.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
