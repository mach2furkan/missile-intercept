# Missile Intercept Simulator (Defense-Grade)

![License](https://img.shields.io/badge/license-MIT-blue.svg)
![Go](https://img.shields.io/badge/backend-Go_1.22+-00ADD8.svg)
![React](https://img.shields.io/badge/frontend-React_18-61DAFB.svg)
![Three.js](https://img.shields.io/badge/3D-Three.js-white.svg)

A professional-grade, real-time 3D missile interception simulator designed for tactical analysis and visualization. This system features a high-performance physics engine (Go), a realistic guidance system (Proportional Navigation), and a dual-view tactical interface (3D/2D Radar) built with React and Three.js.

![Simulator Screenshot](screenshot.png) *(Place a screenshot here)*

## Features

### 🚀 Advanced Simulation Engine
-   **High-Fidelity Physics**: 60Hz server-side physics update loop with drag, gravity, and thrust modeling.
-   **Guidance Algorithms**: Implements Proportional Navigation (PN), Pure Pursuit, and Lead Pursuit.
-   **Real-Time State Sync**: WebSocket-based low-latency state synchronization.
-   **Y-Up Coordinate System**: Standard aerospace coordinate system (X: East, Y: Altitude, Z: North).

### 🖥️ Tactical Dashboard (Frontend)
-   **Dual View Modes**:
    -   **3D Tactical View**: Interactive 3D environment with Orbit, Missile-Chase, and Target-Chase cameras.
    -   **2D Radar View**: CRT-style tactical map with scanlines, variable range, and historic trails.
-   **Defense-Grade UI**:
    -   Real-time telemetry (Altitude, Speed, G-Force, Closing Velocity).
    -   Threat classification and intercept probability analysis.
    -   Tactical event log and system status indicators.
-   **Input Control**: Full keyboard and mouse support for camera and simulation control.

## Project Structure

```bash
missile-intercept-sim/
├── backend/            # Go Simulation Server
│   ├── cmd/server/     # Entry point
│   ├── internal/       # Core logic (physics, entities, guidance)
│   └── pkg/vector/     # 3D Vector math library
├── frontend/           # React + Vite + Three.js Client
│   ├── src/
│   │   ├── components/ # UI, Radar, & 3D Scene components
│   │   └── assets/     # Static assets
│   └── ...
├── Makefile            # Build automation
└── README.md           # This file
```

## Getting Started

### Prerequisites
-   **Go**: Version 1.22 or higher.
-   **Node.js**: Version 18 or higher.
-   **Make**: (Optional) for running automation scripts.

### Installation

1.  **Clone the repository**
    ```bash
    git clone https://github.com/yourusername/missile-intercept-sim.git
    cd missile-intercept-sim
    ```

2.  **Install Dependencies**
    ```bash
    # Backend
    cd backend && go mod tidy
    
    # Frontend
    cd ../frontend && npm install
    ```

### Running the Simulator

The easiest way to run the full system is using the `Makefile` from the root directory.

1.  **Start the Backend** (Terminal 1)
    ```bash
    make run-backend
    ```

2.  **Start the Frontend** (Terminal 2)
    ```bash
    make run-frontend
    ```

3.  **Access the Dashboard**
    Open [http://localhost:5173](http://localhost:5173) in your browser.

## Controls

| Key | Action |
| :--- | :--- |
| **Space** | Launch Interceptor |
| **P** | Pause Simulation |
| **R** | Reset Simulation |
| **TAB** | Toggle 2D/3D View |
| **W/A/S/D** | Move Camera (Free Mode) |
| **Q/E** | Move Camera Up/Down |
| **Shift** | Boost Camera Speed |
| **O** | Orbit Camera Mode |
| **M** | Missile Chase Camera |
| **T** | Target Chase Camera |

## Architecture

This project demonstrates a rigorous separation of concerns:
-   **Backend (Authority)**: The Go server handles all physics calculations, collision detection, and guidance logic. It broadcasts the 'Truth State' to clients.
-   **Frontend (Visualization)**: The React client is a pure visualizer and input relay. It interpolates state for smooth rendering but calculates no physics.

## License

MIT License - see [LICENSE](LICENSE) for details.
