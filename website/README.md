# Visual API Testing Platform - Frontend

A visual flow-based editor for designing, executing, and verifying API tests using TDD workflows.

## Features

- 🎨 **Visual Flow Editor**: Drag-and-drop interface for creating test flows
- 🔌 **Node Types**: API, Verification, Mock, Report, and Event Trigger nodes
- ⚡ **Real-time Execution**: Mock execution engine with live status updates
- 📊 **Dashboard**: View flows and test run history
- 🎯 **TDD Support**: Mock nodes for testing before APIs exist

## Getting Started

### Prerequisites

- Node.js 18+ 
- npm or yarn

### Installation

```bash
cd website
npm install
```

### Development

```bash
npm run dev
```

The app will be available at `http://localhost:5173`

### Build

```bash
npm run build
```

## Project Structure

```
src/
├── components/
│   ├── ui/          # shadcn/ui components
│   ├── nodes/       # React Flow node components
│   ├── FlowEditor.tsx
│   └── Dashboard.tsx
├── services/
│   └── mockExecution.ts  # Mock execution engine
├── stores/
│   └── flowStore.ts      # Zustand state management
├── types/
│   └── index.ts          # TypeScript types
└── App.tsx

```

## Tech Stack

- **React 18** + **TypeScript**
- **Vite** - Build tool
- **React Flow** - Flow editor
- **shadcn/ui** - UI components
- **Zustand** - State management
- **Tailwind CSS** - Styling

## Next Steps

This is a mock frontend. The next phase will integrate with the Go backend for:
- Persistent flow storage
- Real API execution
- Database integration
- WebSocket real-time updates

