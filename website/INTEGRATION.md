# Frontend-Backend Integration Guide

## ✅ Integration Complete!

The frontend is now fully integrated with the Go backend. Here's what was implemented:

### 🔐 Authentication
- **Login/Register** pages with JWT authentication
- **AuthContext** for managing user state across the app
- **Protected routes** that redirect to login if not authenticated
- **Automatic token management** in API requests

### 📡 API Services
- **API Client** (`services/api.ts`) - Axios instance with automatic token injection
- **Auth Service** (`services/authService.ts`) - Login, register, logout
- **Flow Service** (`services/flowService.ts`) - CRUD operations for flows
- **WebSocket Service** (`services/websocketService.ts`) - Real-time execution updates

### 🔄 State Management
- **Updated Flow Store** - Now syncs with backend API instead of local state
- **Auto-save** - Flows are automatically saved when nodes/edges change
- **Real-time updates** - WebSocket connection for live execution status

### 🎨 Updated Components
- **Dashboard** - Loads flows from backend, displays real test runs
- **FlowEditor** - Saves flows to backend, triggers execution via API
- **Login** - Authentication interface
- **App** - Route protection and authentication flow

## 🚀 How to Use

### 1. Start the Backend
```bash
cd server
docker compose up -d
```

Backend runs on: `http://localhost:8080`

### 2. Start the Frontend
```bash
cd website
npm run dev
```

Frontend runs on: `http://localhost:5173`

### 3. First Time Setup
1. Open `http://localhost:5173`
2. You'll be redirected to `/login`
3. Click "Register" to create an account
4. Fill in name, email, and password (min 8 characters)
5. You'll be automatically logged in and redirected to the dashboard

### 4. Using the Application
- **Create Flow**: Click "New Flow" on the dashboard
- **Edit Flow**: Click on a flow card or switch to "Flow Editor" tab
- **Add Nodes**: Use the node buttons (top-right in editor)
- **Configure Nodes**: Click any node to open settings
- **Run Flow**: Click "Run Flow" button
- **View Results**: See real-time execution updates in the editor

## 📡 API Endpoints Used

- `POST /api/auth/register` - User registration
- `POST /api/auth/login` - User login
- `GET /api/flows` - List all flows for user
- `POST /api/flows` - Create new flow
- `PUT /api/flows/:id` - Update flow
- `DELETE /api/flows/:id` - Delete flow
- `POST /api/flows/:id/run` - Execute flow
- `GET /api/flows/:id/test-runs` - Get test runs for flow
- `GET /api/ws?testRunId=<uuid>` - WebSocket connection

## 🔧 Configuration

Create `.env` file in `website/` directory:
```
VITE_API_URL=http://localhost:8080/api
```

## 📝 Notes

- All API calls include JWT token automatically
- Flows auto-save after 1 second of inactivity
- WebSocket connects when a flow execution starts
- Token is stored in localStorage
- Unauthorized requests redirect to login

## 🐛 Troubleshooting

1. **CORS Errors**: Ensure backend CORS is configured for `http://localhost:5173`
2. **401 Errors**: Check if token is valid, try logging out and in again
3. **WebSocket Not Connecting**: Ensure test run ID is available after flow execution
4. **Flows Not Loading**: Check browser console for API errors

