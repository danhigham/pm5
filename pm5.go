// Package pm5 provides a Go library for communicating with Concept2 PM5 rowing computers
// via USB using the CSAFE protocol.
package pm5

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/danhigham/pm5/csafe"
	"github.com/danhigham/pm5/device"
)

var (
	ErrNotConnected    = errors.New("not connected to PM5")
	ErrInvalidResponse = errors.New("invalid response from PM5")
	ErrCommandFailed   = errors.New("command failed")
)

// PM5 represents a connection to a Concept2 PM5 rowing computer
type PM5 struct {
	device        device.HIDDevice
	mu            sync.Mutex
	connected     bool
	frameToggle   bool
	interframeDur time.Duration
	lastCommand   time.Time
}

// New creates a new PM5 instance with the given HID device
func New(dev device.HIDDevice) *PM5 {
	return &PM5{
		device:        dev,
		interframeDur: time.Duration(csafe.MinInterframeGapMs) * time.Millisecond,
	}
}

// Connect opens the connection to the PM5
func (p *PM5) Connect() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.connected {
		return nil
	}

	if err := p.device.Open(); err != nil {
		return fmt.Errorf("failed to open device: %w", err)
	}

	p.connected = true
	return nil
}

// Disconnect closes the connection to the PM5
func (p *PM5) Disconnect() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.connected {
		return nil
	}

	if err := p.device.Close(); err != nil {
		return fmt.Errorf("failed to close device: %w", err)
	}

	p.connected = false
	return nil
}

// IsConnected returns whether the PM5 is connected
func (p *PM5) IsConnected() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.connected
}

// sendCommand sends a CSAFE command and returns the response
func (p *PM5) sendCommand(contents []byte) (*csafe.Response, error) {
	if !p.connected {
		return nil, ErrNotConnected
	}

	// Enforce minimum inter-frame gap
	elapsed := time.Since(p.lastCommand)
	if elapsed < p.interframeDur {
		time.Sleep(p.interframeDur - elapsed)
	}

	// Build and encode the frame
	frame := &csafe.Frame{
		Extended: false,
		Contents: contents,
	}

	encoded, err := csafe.EncodeFrame(frame)
	if err != nil {
		return nil, fmt.Errorf("failed to encode frame: %w", err)
	}

	fmt.Printf(">> % X\n", encoded)

	// Write to device
	_, err = p.device.Write(encoded)
	if err != nil {
		return nil, fmt.Errorf("failed to write to device: %w", err)
	}

	p.lastCommand = time.Now()

	// Read response
	data, err := p.device.Read(500 * time.Millisecond)
	if err != nil {
		return nil, fmt.Errorf("failed to read from device: %w", err)
	}

	fmt.Printf("<< % X\n", data)

	// Find frame boundaries in response
	startIdx := -1
	stopIdx := -1
	for i, b := range data {
		if b == csafe.StandardFrameStartFlag || b == csafe.ExtendedFrameStartFlag {
			startIdx = i
		}
		if b == csafe.StopFrameFlag && startIdx >= 0 {
			stopIdx = i
			break
		}
	}

	if startIdx < 0 || stopIdx < 0 {
		return nil, ErrInvalidResponse
	}

	// Decode the frame
	respFrame, err := csafe.DecodeFrame(data[startIdx : stopIdx+1])
	if err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Parse the response
	resp, err := csafe.ParseResponse(respFrame.Contents)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Check for errors
	if resp.PrevFrameStatus == csafe.PrevFrameStatusReject {
		return resp, ErrCommandFailed
	}

	return resp, nil
}

// sendPMCommand sends a PM-specific command
func (p *PM5) sendPMCommand(wrapper byte, pmCmds ...[]byte) (*csafe.Response, error) {
	contents := csafe.BuildPMCommand(wrapper, pmCmds...)
	return p.sendCommand(contents)
}

// ============================================================================
// Public CSAFE Commands
// ============================================================================

// GetStatus returns the current status byte
func (p *PM5) GetStatus() (*csafe.Response, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.sendCommand([]byte{csafe.CmdGetStatus})
}

// Reset sends a reset command
func (p *PM5) Reset() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	_, err := p.sendCommand([]byte{csafe.CmdReset})
	return err
}

// GoIdle sends the PM to idle state
func (p *PM5) GoIdle() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	_, err := p.sendCommand([]byte{csafe.CmdGoIdle})
	return err
}

// GoReady sends the PM to ready state
func (p *PM5) GoReady() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	_, err := p.sendCommand([]byte{csafe.CmdGoReady})
	return err
}

// GoInUse sends the PM to in-use state
func (p *PM5) GoInUse() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	_, err := p.sendCommand([]byte{csafe.CmdGoInUse})
	return err
}

// GoFinished sends the PM to finished state
func (p *PM5) GoFinished() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	_, err := p.sendCommand([]byte{csafe.CmdGoFinished})
	return err
}

// Version represents PM firmware/hardware version info
type Version struct {
	ManufacturerID byte
	ClassID        byte
	Model          byte
	HWVersion      uint16
	SWVersion      uint16
}

// GetVersion returns the PM version information
func (p *PM5) GetVersion() (*Version, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	resp, err := p.sendCommand([]byte{csafe.CmdGetVersion})
	if err != nil {
		return nil, err
	}

	if len(resp.CommandData) == 0 || len(resp.CommandData[0].Data) < 7 {
		return nil, ErrInvalidResponse
	}

	data := resp.CommandData[0].Data
	return &Version{
		ManufacturerID: data[0],
		ClassID:        data[1],
		Model:          data[2],
		HWVersion:      uint16(data[3]) | uint16(data[4])<<8,
		SWVersion:      uint16(data[5]) | uint16(data[6])<<8,
	}, nil
}

// GetSerial returns the PM serial number as a string
func (p *PM5) GetSerial() (string, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	resp, err := p.sendCommand([]byte{csafe.CmdGetSerial})
	if err != nil {
		return "", err
	}

	if len(resp.CommandData) == 0 {
		return "", ErrInvalidResponse
	}

	return string(resp.CommandData[0].Data), nil
}

// WorkTime represents elapsed work time
type WorkTime struct {
	Hours      byte
	Minutes    byte
	Seconds    byte
	Hundredths byte // Fractional seconds (0.01s)
}

// GetTWork returns the current work time
func (p *PM5) GetTWork() (*WorkTime, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	resp, err := p.sendCommand([]byte{csafe.CmdGetTWork})
	if err != nil {
		return nil, err
	}

	if len(resp.CommandData) == 0 || len(resp.CommandData[0].Data) < 3 {
		return nil, ErrInvalidResponse
	}

	data := resp.CommandData[0].Data
	return &WorkTime{
		Hours:   data[0],
		Minutes: data[1],
		Seconds: data[2],
	}, nil
}

// GetCalories returns the total calories burned
func (p *PM5) GetCalories() (uint16, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	resp, err := p.sendCommand([]byte{csafe.CmdGetCalories})
	if err != nil {
		return 0, err
	}

	if len(resp.CommandData) == 0 || len(resp.CommandData[0].Data) < 2 {
		return 0, ErrInvalidResponse
	}

	data := resp.CommandData[0].Data
	return uint16(data[0]) | uint16(data[1])<<8, nil
}

// GetHorizontal returns the horizontal distance in meters
func (p *PM5) GetHorizontal() (uint16, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	resp, err := p.sendCommand([]byte{csafe.CmdGetHorizontal})
	if err != nil {
		return 0, err
	}

	if len(resp.CommandData) == 0 || len(resp.CommandData[0].Data) < 2 {
		return 0, ErrInvalidResponse
	}

	data := resp.CommandData[0].Data
	return uint16(data[0]) | uint16(data[1])<<8, nil
}

// GetPace returns the current pace (time per 500m) in hundredths of a second
func (p *PM5) GetPace() (uint16, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	resp, err := p.sendCommand([]byte{csafe.CmdGetPace})
	if err != nil {
		return 0, err
	}

	if len(resp.CommandData) == 0 || len(resp.CommandData[0].Data) < 2 {
		return 0, ErrInvalidResponse
	}

	data := resp.CommandData[0].Data
	return uint16(data[0]) | uint16(data[1])<<8, nil
}

// GetCadence returns the current stroke rate
func (p *PM5) GetCadence() (uint16, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	resp, err := p.sendCommand([]byte{csafe.CmdGetCadence})
	if err != nil {
		return 0, err
	}

	if len(resp.CommandData) == 0 || len(resp.CommandData[0].Data) < 2 {
		return 0, ErrInvalidResponse
	}

	data := resp.CommandData[0].Data
	return uint16(data[0]) | uint16(data[1])<<8, nil
}

// GetPower returns the current power in watts
func (p *PM5) GetPower() (uint16, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	resp, err := p.sendCommand([]byte{csafe.CmdGetPower})
	if err != nil {
		return 0, err
	}

	if len(resp.CommandData) == 0 || len(resp.CommandData[0].Data) < 2 {
		return 0, ErrInvalidResponse
	}

	data := resp.CommandData[0].Data
	return uint16(data[0]) | uint16(data[1])<<8, nil
}

// GetHeartRate returns the current heart rate
func (p *PM5) GetHeartRate() (byte, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	resp, err := p.sendCommand([]byte{csafe.CmdGetHRCur})
	if err != nil {
		return 0, err
	}

	if len(resp.CommandData) == 0 || len(resp.CommandData[0].Data) < 1 {
		return 0, ErrInvalidResponse
	}

	return resp.CommandData[0].Data[0], nil
}

// SetProgram sets a predefined workout program
func (p *PM5) SetProgram(workoutNum csafe.WorkoutNumber) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	cmd := csafe.BuildCommand(csafe.CmdSetProgram, byte(workoutNum), 0x00)
	_, err := p.sendCommand(cmd)
	return err
}

// SetTWork sets the workout time goal
func (p *PM5) SetTWork(hours, minutes, seconds byte) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	cmd := csafe.BuildCommand(csafe.CmdSetTWork, hours, minutes, seconds)
	_, err := p.sendCommand(cmd)
	return err
}

// SetHorizontal sets the horizontal distance goal in meters
func (p *PM5) SetHorizontal(distance uint16) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	cmd := csafe.BuildCommand(csafe.CmdSetHorizontal,
		byte(distance&0xFF),
		byte((distance>>8)&0xFF),
		csafe.UnitsMeter)
	_, err := p.sendCommand(cmd)
	return err
}

// SetCalories sets the calorie goal
func (p *PM5) SetCalories(calories uint16) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	cmd := csafe.BuildCommand(csafe.CmdSetCalories,
		byte(calories&0xFF),
		byte((calories>>8)&0xFF))
	_, err := p.sendCommand(cmd)
	return err
}

// SetPower sets the power goal in watts
func (p *PM5) SetPower(watts uint16) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	cmd := csafe.BuildCommand(csafe.CmdSetPower,
		byte(watts&0xFF),
		byte((watts>>8)&0xFF),
		csafe.UnitsWatt)
	_, err := p.sendCommand(cmd)
	return err
}
