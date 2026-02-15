package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"missile-intercept-sim/internal/simulation"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all for dev
	},
}

var sim *simulation.Simulator

func main() {
	sim = simulation.NewSimulator()

	http.HandleFunc("/api/start", handleStart)
	http.HandleFunc("/api/stop", handleStop)
	http.HandleFunc("/api/reset", handleReset)
	http.HandleFunc("/api/guidance", handleGuidance)
	http.HandleFunc("/ws", handleWebSocket)

	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func handleStart(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	sim.Start()
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Simulation started"))
}

func handleStop(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	sim.Stop()
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Simulation stopped"))
}

func handleReset(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	sim.Stop()
	sim.Reset()
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Simulation reset"))
}

func handleGuidance(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	type GuidanceRequest struct {
		Mode string `json:"mode"`
	}
	var req GuidanceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}
	sim.SetGuidanceMode(req.Mode)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Guidance mode updated"))
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()

	// Broadcast loop for this client
	ticker := time.NewTicker(33 * time.Millisecond) // ~30Hz update for UI
	defer ticker.Stop()

	for range ticker.C {
		state := sim.GetState()
		err := c.WriteJSON(state)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}
