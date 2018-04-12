package elevator

import (
	"errors"
	"math"
	"sync"
)

type SafeElevator struct {
	Elevator
	sync.Mutex
}

// Elevator is a representation of an elevator from sanntidssal
type Elevator struct {
	ID           int     `json:"id"`
	NFloors      int     `json:"nFloors"`
	Position     float64 `json:"position"`
	Speed        float64 `json:"speed"`
	Direction    float64 `json:"direction"`
	Blocked      bool    `json:"blocked"`
	Margin       float64 `json:"margin"`
	CurrentFloor int     `json:"currentFloor"`

	ObstructionButton bool `json:"obstructionButton"`
	StopButton        bool `json:"stopButton"`

	UpButtons    []bool `json:"upButtons"`
	DownButtons  []bool `json:"downButtons"`
	CabinButtons []bool `json:"cabinButtons"`

	StopLamp        bool   `json:"stopLamp"`
	ObstructionLamp bool   `json:"obstructionLamp"`
	DoorLamp        bool   `json:"doorLamp"`
	UpLamps         []bool `json:"upLamps"`
	DownLamps       []bool `json:"downLamps"`
	CabinLamps      []bool `json:"cabinLamps"`
	IndicatorLamp   int    `json:"indicatorLamp"`
}

// NewElevator creates an elevator.
func NewElevator(id, nFloors int, speed, margin float64) Elevator {
	return Elevator{
		ID:      id,
		NFloors: nFloors,
		Speed:   speed,
		Margin:  margin,

		UpButtons:    make([]bool, nFloors),
		DownButtons:  make([]bool, nFloors),
		CabinButtons: make([]bool, nFloors),

		UpLamps:    make([]bool, nFloors),
		DownLamps:  make([]bool, nFloors),
		CabinLamps: make([]bool, nFloors),
	}
}

// ValidateFloor validates if the floor is one that exists for the elevators
func (elevator *Elevator) ValidateFloor(floor int) error {
	if floor < 0 || floor >= elevator.NFloors {
		return errors.New("floor out of range")
	}
	return nil
}

// Update moves the elevator and sets the blocked flag
func (elevator *Elevator) Update(dTime float64) {
	pos := elevator.Position + elevator.Speed*elevator.Direction*dTime
	pos, clamped := clamp(-0.5, float64(elevator.NFloors)-0.5, pos)
	if clamped {
		elevator.Blocked = true
	}
	elevator.Position = pos

	elevator.CurrentFloor = elevator.GetFloorSignal()
}

// SetDirection sets the travel direction of the elevator. Permitted values are 1 0 and -1
func (elevator *Elevator) SetDirection(direction float64) error {

	if direction != 0 && direction != 1.0 && direction != -1.0 {
		return errors.New("direction must be one of the following: 1 0 -1")
	}
	elevator.Direction = direction
	return nil
}

// SetDownButtonLamp sets the lamp in the down button on the selected floor
func (elevator *Elevator) SetDownButtonLamp(floor int, value bool) error {
	if floor == 0 {
		return errors.New("floor can't be the ground level, no down button there")
	}
	if err := elevator.ValidateFloor(floor); err != nil {
		return err
	}
	elevator.DownLamps[floor] = value
	return nil
}

// SetUpButtonLamp sets the lamp in the up button on the selected floor
func (elevator *Elevator) SetUpButtonLamp(floor int, value bool) error {
	if floor == elevator.NFloors {
		return errors.New("floor can't be the ground level, no down button there")
	}

	if err := elevator.ValidateFloor(floor); err != nil {
		return err
	}
	elevator.UpLamps[floor] = value
	return nil
}

// SetCabinButtonLamp sets the lamp in the cabin button on the selected floor
func (elevator *Elevator) SetCabinButtonLamp(floor int, value bool) error {
	if err := elevator.ValidateFloor(floor); err != nil {
		return err
	}
	elevator.CabinLamps[floor] = value
	return nil
}

// SetFloorIndicator sets the lamp in the floor indicator on the selected floor
func (elevator *Elevator) SetFloorIndicator(floor int) error {
	if err := elevator.ValidateFloor(floor); err != nil {
		return err
	}
	elevator.IndicatorLamp = floor
	return nil
}

// SetDoorLamp sets the lamp telling if the door is open or closed
func (elevator *Elevator) SetDoorLamp(value bool) {
	elevator.DoorLamp = value
}

// SetStopLamp sets the lamp in the stop button
func (elevator *Elevator) SetStopLamp(value bool) {
	elevator.DoorLamp = value
}

// GetUpButton gets if the up button in the selected floor is pressed
func (elevator *Elevator) GetUpButton(floor int) (bool, error) {
	if floor == elevator.NFloors {
		return false, errors.New("floor can't be the ground level, no down button there")
	}

	if err := elevator.ValidateFloor(floor); err != nil {
		return false, err
	}

	return elevator.UpButtons[floor], nil
}

func (elevator *Elevator) SetUpButton(floor int, value bool) error {
	if floor == elevator.NFloors {
		return errors.New("floor can't be the ground level, no down button there")
	}

	if err := elevator.ValidateFloor(floor); err != nil {
		return err
	}

	elevator.UpButtons[floor] = value
	return nil
}

// GetDownButton gets if the down button in the selected floor is pressed
func (elevator *Elevator) GetDownButton(floor int) (bool, error) {
	if floor == 0 {
		return false, errors.New("floor can't be the ground level, no down button there")
	}
	if err := elevator.ValidateFloor(floor); err != nil {
		return false, err
	}

	return elevator.DownButtons[floor], nil
}

func (elevator *Elevator) SetDownButton(floor int, value bool) error {
	if floor == 0 {
		return errors.New("floor can't be the ground level, no down button there")
	}
	if err := elevator.ValidateFloor(floor); err != nil {
		return err
	}
	elevator.DownButtons[floor] = value
	return nil
}

// GetCabinButton gets if the cabin button in the selected floor is pressed
func (elevator *Elevator) GetCabinButton(floor int) (bool, error) {
	if err := elevator.ValidateFloor(floor); err != nil {
		return false, err
	}

	return elevator.CabinButtons[floor], nil
}

func (elevator *Elevator) SetCabinButton(floor int, value bool) error {
	if err := elevator.ValidateFloor(floor); err != nil {
		return err
	}

	elevator.CabinButtons[floor] = value
	return nil
}

// GetFloorSignal returns which of the floorsensors are activated, -1 if the elevator is between floors
func (elevator *Elevator) GetFloorSignal() int {
	nearestFloor := math.Round(elevator.Position)
	if math.Abs(nearestFloor-elevator.Position) < elevator.Margin {
		return int(nearestFloor)
	}
	return -1
}

// GetStopButton returns if the stop button is pressed
func (elevator *Elevator) GetStopButton() bool {
	return elevator.StopButton
}

func (elevator *Elevator) SetStopButton(value bool) {
	elevator.StopButton = value
}

// GetObstructionButton return if the Obstruction button is pressed
func (elevator *Elevator) GetObstructionButton() bool {
	return elevator.ObstructionButton
}

func (elevator *Elevator) SetObstructionButton(value bool) {
	elevator.ObstructionButton = value
}

func clamp(min, max, value float64) (float64, bool) {
	if value < min {
		return min, true
	}
	if value > max {
		return max, true
	}
	return value, false
}
