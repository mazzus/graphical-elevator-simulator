package standard

import (
	"fmt"
	"io"
	"net"

	"github.com/mazzus/graphical-elevator-simulator/backend/src/elevator"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("standard handler")

func Server(port int, safeElevator *elevator.SafeElevator) {
	addrString := fmt.Sprintf(":%d", port)
	addr, err := net.ResolveTCPAddr("tcp", addrString)

	log.Info(readFloorSensor)

	if err != nil {
		log.Error("Error resolving addr, server will not start", err)
		return
	}

	listener, err := net.ListenTCP("tcp", addr)
	defer listener.Close()

	if err != nil {
		log.Error("Error opening listener, server will not start", err)
		return
	}

	for {
		log.Debug("Waiting for connection")
		conn, err := listener.Accept()
		log.Info("Got a connection")
		if err != nil {
			log.Warning("Error opening connection, lost this call", err)
			continue
		}
		go handleConnection(conn, safeElevator)
	}
}

const (
	doNothing byte = 0
)

const (
	writeMotorDirection byte = iota + 1
	writeOrderButtonLight
	writeFloorIndicator
	writeDoorOpenLight
	writeStopButtonLight
)

const (
	readOrderButton byte = iota + 6
	readFloorSensor
	readStopButton
	readObstructionSwitch
)
const (
	orderButtonUp byte = iota
	orderButtonDown
	orderButtonCabin
)

func handleConnection(connection io.ReadWriteCloser, safeElevator *elevator.SafeElevator) {
	defer connection.Close()
	for {

		command := make([]byte, 4)
		_, err := io.ReadFull(connection, command)
		if err != nil {
			log.Warning("Could not read command!", err)
			continue
		}

		switch command[0] {
		case doNothing:
		case writeMotorDirection:
			handleWriteMotorDirection(command, safeElevator)
		case writeOrderButtonLight:
			handleWriteOrderButtonLight(command, safeElevator)
		case writeFloorIndicator:
			handleWriteFloorIndicator(command, safeElevator)
		case writeDoorOpenLight:
			handleWriteDoorOpenLight(command, safeElevator)
		case writeStopButtonLight:
			handleWriteStopButtonLight(command, safeElevator)
		case readOrderButton:
			handleReadOrderButton(command, connection, safeElevator)
		case readFloorSensor:
			handleReadFloorSensor(command, connection, safeElevator)
		case readStopButton:
			handleReadStopButton(command, connection, safeElevator)
		case readObstructionSwitch:
			handleReadObstructionButton(command, connection, safeElevator)
		default:
			log.Warning("Could not decode command: ", command[0])
		}
	}
}

func handleWriteMotorDirection(command []byte, safeElevator *elevator.SafeElevator) {
	log.Debug("SETMOTORDIR")

	var direction float64

	switch command[1] {
	case 0:
		direction = 0
	case 1:
		direction = 1
	case 255:
		direction = -1
	}

	safeElevator.Lock()
	err := safeElevator.SetDirection(direction)
	safeElevator.Unlock()
	if err != nil {
		log.Warning("Could not set motor direction")
	}

	log.Debug("/SETMOTORDIR")
}

func handleWriteOrderButtonLight(command []byte, safeElevator *elevator.SafeElevator) {
	log.Debug("SETBUTTONLAMP")

	var f func(*elevator.SafeElevator, int, bool) error
	switch command[1] {
	case orderButtonUp:
		f = (*elevator.SafeElevator).SetUpButtonLamp
	case orderButtonDown:
		f = (*elevator.SafeElevator).SetDownButtonLamp
	case orderButtonCabin:
		f = (*elevator.SafeElevator).SetCabinButtonLamp
	}

	safeElevator.Lock()
	f(safeElevator, int(command[2]), command[3] == 1)
	safeElevator.Unlock()
	log.Debug("/SETBUTTONLAMP")
}

func handleWriteFloorIndicator(command []byte, safeElevator *elevator.SafeElevator) {

	log.Debug("SETFLOORINDICATOR")

	safeElevator.Lock()
	safeElevator.SetFloorIndicator(int(command[1]))
	safeElevator.Unlock()

	log.Debug("/SETFLOORINDICATOR")
}

func handleWriteDoorOpenLight(command []byte, safeElevator *elevator.SafeElevator) {

	log.Debug("SETDOORLAMP")

	safeElevator.Lock()
	safeElevator.SetDoorLamp(command[1] == 1)
	safeElevator.Unlock()

	log.Debug("/SETDOORLAMP")
}

func handleWriteStopButtonLight(command []byte, safeElevator *elevator.SafeElevator) {
	log.Debug("SETSTOP")

	safeElevator.Lock()
	safeElevator.SetStopLamp(command[1] == 1)
	safeElevator.Unlock()

	log.Debug("/SETSTOP")
}

func handleReadOrderButton(command []byte, connection io.Writer, safeElevator *elevator.SafeElevator) {
	log.Debug("GETBUTTONLAMP")

	var f func(*elevator.SafeElevator, int) (bool, error)
	switch command[1] {
	case orderButtonUp:
		f = (*elevator.SafeElevator).GetUpButton
	case orderButtonDown:
		f = (*elevator.SafeElevator).GetDownButton
	case orderButtonCabin:
		f = (*elevator.SafeElevator).GetCabinButton
	}

	safeElevator.Lock()
	value, err := f(safeElevator, int(command[2]))
	safeElevator.Unlock()

	if err != nil {
		log.Warning("Could not get button")
	}

	var v byte
	if value {
		v = 1
	} else {
		v = 0
	}

	connection.Write([]byte{readOrderButton, v, 0, 0})
}

func handleReadFloorSensor(command []byte, connection io.Writer, safeElevator *elevator.SafeElevator) {

	log.Debug("GETFLOOR")

	safeElevator.Lock()
	floor := safeElevator.GetFloorSignal()
	safeElevator.Unlock()

	if floor == -1 {
		connection.Write([]byte{readFloorSensor, 0, 0, 0})
	} else {
		connection.Write([]byte{readFloorSensor, 1, byte(floor), 0})
	}

	log.Debug("/GETFLOOR")
}

func handleReadStopButton(command []byte, connection io.Writer, safeElevator *elevator.SafeElevator) {
	log.Debug("GETSTOP")

	safeElevator.Lock()
	stopped := safeElevator.GetStopButton()
	safeElevator.Unlock()

	var v byte
	if stopped {
		v = 1
	} else {
		v = 0
	}
	connection.Write([]byte{readStopButton, v, 0, 0})
	log.Debug("/GETSTOP")
}

func handleReadObstructionButton(command []byte, connection io.Writer, safeElevator *elevator.SafeElevator) {
	log.Debug("GETOBSTRUCTION")

	safeElevator.Lock()
	obstructed := safeElevator.GetObstructionButton()
	safeElevator.Unlock()

	var v byte
	if obstructed {
		v = 1
	} else {
		v = 0
	}
	connection.Write([]byte{readObstructionSwitch, v, 0, 0})
	log.Debug("/GETOBSTRUCTION")
}
