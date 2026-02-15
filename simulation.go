package simulation

import (
	"log"
	"sync"
	"time"

	"missile-intercept-sim/internal/entities"
	"missile-intercept-sim/internal/guidance"
	"missile-intercept-sim/internal/physics"
	"missile-intercept-sim/pkg/vector"
)

// SimulationState holds the current state of the world.
type SimulationState struct {
	Entities  []*entities.Entity `json:"entities"`
	Status    string             `json:"status"` // Running, Stopped, Intercepted
	Time      float64            `json:"time"`
	Intercept bool               `json:"intercept"`
}

// Simulator manages the simulation loop and state.
type Simulator struct {
	State        SimulationState
	mu           sync.RWMutex
	ticker       *time.Ticker
	stopChan     chan bool
	Target       *entities.Entity
	Missile      *entities.Entity
	GuidanceLaw  guidance.GuidanceLaw
	GuidanceName string
	Dt           float64
}

// NewSimulator creates a new simulator instance.
func NewSimulator() *Simulator {
	sim := &Simulator{
		State: SimulationState{
			Entities: []*entities.Entity{},
			Status:   "Stopped",
			Time:     0.0,
		},
		Dt: 0.016, // Approx 60Hz
	}
	// Initialize default entities for reset
	sim.Reset()
	return sim
}

// Reset restores the simulation to initial state.
func (s *Simulator) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Default Scenario: Target flying level, Missile launching from ground
	// Y-UP System: X=East, Y=Alt, Z=North
	// Target at 5000m East, 2000m Alt, 5000m North
	target := entities.NewTarget("target-1",
		vector.Vector3{X: 5000, Y: 2000, Z: 5000},
		vector.Vector3{X: -200, Y: 0, Z: -100}, // Moving West and South
	)

	missile := entities.NewMissile("missile-1",
		vector.Vector3{X: 0, Y: 0, Z: 0},
		vector.Vector3{X: 0, Y: 0, Z: 0}, // Launch velocity zero? Or minimal.
	)
	// Give initial boost
	// Launch Upwards (Y+) and slightly towards target
	missile.Velocity = vector.Vector3{X: 10, Y: 10, Z: 10}

	s.Target = target
	s.Missile = missile
	s.GuidanceName = "ProNav" // Default
	s.GuidanceLaw = guidance.GetFactory(s.GuidanceName)

	s.State = SimulationState{
		Entities:  []*entities.Entity{target, missile},
		Status:    "Stopped",
		Time:      0.0,
		Intercept: false,
	}
}

// Start resumes the simulation loop.
func (s *Simulator) Start() {
	s.mu.Lock()
	if s.State.Status == "Running" {
		s.mu.Unlock()
		return
	}
	s.State.Status = "Running"
	s.stopChan = make(chan bool)
	s.ticker = time.NewTicker(time.Duration(s.Dt * float64(time.Second)))
	s.mu.Unlock()

	go s.loop()
}

// Stop pauses the simulation loop.
func (s *Simulator) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.State.Status == "Running" {
		s.State.Status = "Stopped"
		if s.ticker != nil {
			s.ticker.Stop()
		}
		if s.stopChan != nil {
			close(s.stopChan)
		}
	}
}

// SetGuidanceMode changes the active guidance law.
func (s *Simulator) SetGuidanceMode(mode string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.GuidanceName = mode
	s.GuidanceLaw = guidance.GetFactory(mode)
	s.Missile.GuidanceMode = mode
}

// loop is the main physics loop running in a goroutine.
func (s *Simulator) loop() {
	for {
		select {
		case <-s.stopChan:
			return
		case <-s.ticker.C:
			s.Step()
		}
	}
}

// Step performs one physics integration step.
func (s *Simulator) Step() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.State.Status != "Running" {
		return
	}

	dt := s.Dt

	// 1. Calculate Guidance Interceptor
	// Missile guidance logic
	// Accel command
	accelCmd := s.GuidanceLaw.CalculateAcceleration(s.Missile, s.Target, dt)

	// Limit acceleration (structural limits)
	accelCmd = physics.LimitAcceleration(accelCmd, s.Missile.MaxAccel)

	// Apply Gravity?
	// Real missiles fight gravity.
	// If we want realistic trajectories, we need gravity.
	// If we want "arcade" space style, maybe not.
	// User asked for "physically realistic".
	// Missile needs to generate lift to counter gravity.
	// Ideally: TotalAccel = AeroForces + Gravity + Thrust.
	// Simplification: Guidance outputs "Acceleration Command" (load factor n).
	// The autopilot realizes this acceleration.
	// If guidance says "Accelerate Up 1g", and Gravity pulls "Down 1g",
	// the missile needs to pull 2g aerodynamic lift to achieve 1g net up.
	// Let's assume accelCmd is the NET acceleration requested by guidance *relative to inertial space* (ignoring gravity compensation for now, or assuming guidance compensates).
	// Usually ProNav outputs acceleration perpendicular to LOS.
	// Let's just apply accelCmd directly as the net external force/mass for now used to update velocity,
	// BUT we must also add Gravity to the physical integration if we want parabolic arcs when unguided.
	// However, for ProNav, if we add gravity efficiently, the missile "sags".
	// Let's Add Gravity to specific entities.

	// Apply limits to "requested turn" (Guidance Accel)
	// NOTE: physics.LimitTurnRate is complex without full aerodynamics.
	// Let's rely on LimitAcceleration magnitude for now.

	s.Missile.Acceleration = accelCmd

	// Basic target movement (constant velocity for now, or simple evade?)
	// Let's make target circle or wave if requested. For now constant velocity.
	s.Target.Acceleration = vector.Vector3{}

	// 2. Physics Integration
	// Missile
	// Total Accel = Command + Gravity ??
	// If we just use Command, it flies like a spaceship.
	// Let's add Gravity to both.
	// Physics engine should handle forces.
	// Here we are setting acceleration directly.
	// Let's do: TrueAccel = Cmd + Gravity

	// But ProNav expects to control acceleration. Only need to compensate for gravity if it pulls us off course.
	// A simple approach for "Realistic looking" without full autopilot:
	// Missile has Thrust (forward), Drag (backward), Lift (steer).
	// We are skipping to: "Missile can accelerate in any direction up to MaxAccel".
	// So Accel = Cmd.
	// If we want gravity drop, we add it.
	// Let's add gravity for realism.
	// But then ProNav needs to "bias" upward to compensate.
	// Let's stick to zero gravity for the "Interceptor" logic clarity unless requested "Ballistic".
	// User asked for "Aeropace software engineer" level.
	// Gravity is essential.
	// Let's add Gravity.

	gravity := vector.Vector3{Y: -9.81}

	// Tweak: Guidance command is "Acceleration needed to intercept".
	// It doesn't know about gravity.
	// If we add gravity to the physics update, the missile will sag.
	// The next guidance step will see the sag (velocity error) and correct it.
	// This is how closed-loop guidance works! It automatically compensates for gravity bias.

	s.Missile.Acceleration = s.Missile.Acceleration.Add(gravity)
	s.Target.Acceleration = s.Target.Acceleration.Add(gravity) // Target also falls if not generating lift?
	// Target is usually an airplane maintaining altitude.
	// Assume Logic keeps target level (Autopilot).
	// So Reset Target Accel to zero net (Lift = -Gravity).
	s.Target.Acceleration = vector.Vector3{} // Target flies straight.

	// Update Missile
	newPosM, newVelM := physics.KinematicsUpdate(s.Missile.Position, s.Missile.Velocity, s.Missile.Acceleration, dt)
	s.Missile.Position = newPosM
	s.Missile.Velocity = newVelM

	// Update Target
	newPosT, newVelT := physics.KinematicsUpdate(s.Target.Position, s.Target.Velocity, s.Target.Acceleration, dt)
	s.Target.Position = newPosT
	s.Target.Velocity = newVelT

	s.State.Time += dt

	// 3. Intercept Check
	dist := s.Missile.Position.Distance(s.Target.Position)
	if dist < 5.0 { // Threshold 5 meters
		s.State.Intercept = true
		s.State.Status = "Intercepted"
		s.Stop()
		log.Println("INTERCEPT SUCCESS!")
	}

	// Ground collision check
	if s.Missile.Position.Y < 0 {
		s.Missile.Position.Y = 0
		s.Missile.Velocity = vector.Vector3{}
		s.State.Status = "Crashed"
		s.Stop()
	}
}

// GetState returns the thread-safe state.
func (s *Simulator) GetState() SimulationState {
	s.mu.RLock()
	defer s.mu.RUnlock()
	// Return copy?
	return s.State
}
