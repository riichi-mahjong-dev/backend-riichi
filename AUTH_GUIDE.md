# Authentication Guide

This guide explains how to use the authentication system in the Riichi Mahjong backend application.

## Overview

The authentication system supports two types of users:
- **Players**: Regular users who can play mahjong and submit match results
- **Admins**: Administrative users with two roles:
  - **Super-admin**: Full access to all features
  - **Staff**: Limited administrative access

## Authentication Endpoints

### Player Authentication
- `POST /auth/login/player` - Login as a player
- `POST /auth/refresh` - Refresh access token
- `POST /api/players` - Register new player (public)

### Admin Authentication
- `POST /auth/login/admin` - Login as an admin
- `POST /auth/refresh` - Refresh access token

### Profile
- `GET /api/profile` - Get current user profile (requires authentication)

## Request/Response Examples

### Player Login
```bash
curl -X POST http://localhost:8080/auth/login/player \
  -H "Content-Type: application/json" \
  -d '{
    "username": "player_username",
    "password": "player_password"
  }'
```

### Admin Login
```bash
curl -X POST http://localhost:8080/auth/login/admin \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "admin123"
  }'
```

### Response Format
```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR...",
    "token_type": "Bearer",
    "expires_at": "2025-07-12T15:20:42Z",
    "user": {
      "id": 1,
      "username": "admin",
      "user_type": "admin",
      "role": "super-admin"
    }
  }
}
```

### Using Access Token
```bash
curl -X GET http://localhost:8080/api/profile \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

## Role-Based Access Control

### Permission Matrix

| Endpoint | Guest | Player | Admin | Super-Admin |
|----------|-------|--------|--------|-------------|
| `GET /api/health` | ✓ | ✓ | ✓ | ✓ |
| `GET /api/provinces` | ✓ | ✓ | ✓ | ✓ |
| `GET /api/parlours` | ✓ | ✓ | ✓ | ✓ |
| `GET /api/posts` | ✓ | ✓ | ✓ | ✓ |
| `GET /api/players` | ✓ | ✓ | ✓ | ✓ |
| `GET /api/players/:id` | ✓ | ✓ | ✓ | ✓ |
| `GET /api/matches` | ✓ | ✓ | ✓ | ✓ |
| `GET /api/matches/:id` | ✓ | ✓ | ✓ | ✓ |
| `POST /api/players` | ✓ | - | - | - |
| `GET /api/profile` | - | ✓ | ✓ | ✓ |
| `PUT /api/players/:id` | - | - | ✓ | ✓ |
| `DELETE /api/players/:id` | - | - | - | ✓ |
| `POST /api/matches` | - | ✓ | ✓ | ✓ |
| `PUT /api/matches/:id` | - | ✓ | ✓ | ✓ |
| `DELETE /api/matches/:id` | - | - | ✓ | ✓ |
| `POST /api/matches/:id/approve` | - | - | ✓ | ✓ |
| `POST /api/provinces` | - | - | ✓ | ✓ |
| `PUT /api/provinces/:id` | - | - | ✓ | ✓ |
| `DELETE /api/provinces/:id` | - | - | - | ✓ |
| `POST /api/parlours` | - | - | ✓ | ✓ |
| `PUT /api/parlours/:id` | - | - | ✓ | ✓ |
| `DELETE /api/parlours/:id` | - | - | - | ✓ |
| `POST /api/posts` | - | - | ✓ | ✓ |
| `PUT /api/posts/:id` | - | - | ✓ | ✓ |
| `DELETE /api/posts/:id` | - | - | ✓ | ✓ |
| `GET /api/roles` | - | - | ✓ | ✓ |
| `POST /api/roles` | - | - | - | ✓ |
| `PUT /api/roles/:id` | - | - | - | ✓ |
| `DELETE /api/roles/:id` | - | - | - | ✓ |
| `GET /api/admins` | - | - | - | ✓ |
| `POST /api/admins` | - | - | - | ✓ |
| `PUT /api/admins/:id` | - | - | - | ✓ |
| `DELETE /api/admins/:id` | - | - | - | ✓ |

## Admin Roles

### Super-Admin
- Full access to all endpoints
- Can create/modify/delete other admins
- Can create/modify/delete roles
- Can delete players and other critical operations

### Staff
- Limited administrative access
- Cannot manage other admins
- Cannot manage roles
- Cannot perform destructive operations

## Player Restrictions

**Players have very limited access and can only:**
- View their own profile (`GET /api/profile`)
- View public data (provinces, parlours, posts)
- Access match-related operations:
  - View matches (`GET /api/matches`)
  - View specific match details (`GET /api/matches/:id`)
  - Create new matches (`POST /api/matches`)
  - Update match details (`PUT /api/matches/:id`)

**Players CANNOT:**
- View or modify other players' data
- Create, update, or delete provinces
- Create, update, or delete parlours
- Create, update, or delete posts
- Access admin or role management
- Approve matches (admin only)
- Delete matches (admin only)

**Important:** Players can only manage match data - they cannot post or update any other data tables in the system.

## Default Admin Account

A default super-admin account is created automatically:
- **Username**: `admin`
- **Password**: `admin123`
- **Role**: `super-admin`

**Important**: Change this password immediately in production!

## Error Responses

### Invalid Credentials
```json
{
  "success": false,
  "message": "Invalid credentials",
  "error": "invalid credentials"
}
```

### Unauthorized Access
```json
{
  "success": false,
  "message": "Authorization header required"
}
```

### Access Denied
```json
{
  "success": false,
  "message": "Access denied"
}
```

## Security Features

1. **JWT Tokens**: Secure token-based authentication
2. **Password Hashing**: bcrypt for secure password storage
3. **Role-Based Access**: Granular permission control
4. **Token Expiration**: Access tokens expire in 24 hours
5. **Refresh Tokens**: Longer-lived tokens for seamless user experience

## Development Commands

### Create Default Admin
```bash
go run cmd/seed-admin/main.go
```

### Run Migrations
```bash
go run cmd/migrate/main.go
```

### Start Server
```bash
go run cmd/app/main.go
```

## Environment Variables

Make sure these JWT configuration variables are set in your `.env` file:
```
SECRET_KEY_ACCESS=your-secret-key-here
SECRET_KEY_REFRESH=your-refresh-key-here
```
