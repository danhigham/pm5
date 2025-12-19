// Package device provides USB HID communication with Concept2 Performance Monitors.
package device

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/sstallion/go-hid"
)

func init() {
	// Initialize the HID library
	if err := hid.Init(); err != nil {
		panic(fmt.Sprintf("failed to initialize HID library: %v", err))
	}
}

// USB constants for PM5
const (
	PM5VendorID  uint16 = 0x17A4 // Concept2 Vendor ID
	PM5ProductID uint16 = 0x0046 // PM5 Product ID (may vary)

	// HID Report IDs
	ReportID1Size = 21  // 20 bytes + 1 byte report ID
	ReportID2Size = 121 // 120 bytes + 1 byte report ID

	// Default timeouts
	DefaultReadTimeout  = 500 * time.Millisecond
	DefaultWriteTimeout = 500 * time.Millisecond
)

var (
	ErrDeviceNotFound    = errors.New("PM5 device not found")
	ErrDeviceNotOpen     = errors.New("device not open")
	ErrDeviceAlreadyOpen = errors.New("device already open")
	ErrWriteFailed       = errors.New("write failed")
	ErrReadFailed        = errors.New("read failed")
	ErrTimeout           = errors.New("operation timed out")
)

// HIDDevice is an interface for HID device operations
// This allows for different implementations (real USB, mock for testing)
type HIDDevice interface {
	Open() error
	Close() error
	Write(data []byte) (int, error)
	Read(timeout time.Duration) ([]byte, error)
	IsOpen() bool
	GetInfo() DeviceInfo
}

// DeviceInfo contains information about a connected device
type DeviceInfo struct {
	VendorID     uint16
	ProductID    uint16
	SerialNumber string
	Product      string
	Manufacturer string
	Path         string
}

// USBDevice represents a real USB HID device connection
// using github.com/sstallion/go-hid
type USBDevice struct {
	info         DeviceInfo
	mu           sync.Mutex
	isOpen       bool
	readTimeout  time.Duration
	writeTimeout time.Duration

	// These would be replaced with actual HID device handle
	device *hid.Device
}

// NewUSBDevice creates a new USB device instance
func NewUSBDevice(info DeviceInfo) *USBDevice {
	return &USBDevice{
		info:         info,
		readTimeout:  DefaultReadTimeout,
		writeTimeout: DefaultWriteTimeout,
	}
}

// Open opens the USB device
func (d *USBDevice) Open() error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.isOpen {
		return ErrDeviceAlreadyOpen
	}

	// Open the first PM5 device found
	var err error
	d.device, err = hid.OpenFirst(PM5VendorID, PM5ProductID)
	if err != nil {
		return fmt.Errorf("failed to open device: %w", err)
	}
	if d.device == nil {
		return ErrDeviceNotFound
	}

	d.isOpen = true
	return nil
}

// Close closes the USB device
func (d *USBDevice) Close() error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if !d.isOpen {
		return nil
	}

	if d.device != nil {
		if err := d.device.Close(); err != nil {
			return fmt.Errorf("failed to close device: %w", err)
		}
	}

	d.isOpen = false
	return nil
}

// Write writes data to the device
func (d *USBDevice) Write(data []byte) (int, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if !d.isOpen {
		return 0, ErrDeviceNotOpen
	}

	// Prepare HID report
	// Report ID 2 supports up to 120 bytes
	report := make([]byte, ReportID2Size)
	report[0] = 0x02 // Report ID 2

	// Copy data into report
	n := len(data)
	if n > 120 {
		n = 120
	}
	copy(report[1:], data[:n])

	// Write to device
	written, err := d.device.Write(report)
	if err != nil {
		return 0, fmt.Errorf("HID write failed: %w", err)
	}

	// Return the number of actual data bytes written (excluding report ID)
	if written > 0 {
		return n, nil
	}
	return 0, ErrWriteFailed
}

// Read reads data from the device
func (d *USBDevice) Read(timeout time.Duration) ([]byte, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if !d.isOpen {
		return nil, ErrDeviceNotOpen
	}

	// Allocate buffer for reading (max report size)
	buf := make([]byte, ReportID2Size)

	// Read from device with timeout
	n, err := d.device.ReadWithTimeout(buf, timeout)
	if err != nil {
		return nil, fmt.Errorf("HID read failed: %w", err)
	}

	if n == 0 {
		return nil, ErrTimeout
	}

	// Strip report ID from response
	if n > 1 {
		return buf[1:n], nil
	}

	return nil, ErrReadFailed
}

// IsOpen returns whether the device is open
func (d *USBDevice) IsOpen() bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.isOpen
}

// GetInfo returns device information
func (d *USBDevice) GetInfo() DeviceInfo {
	return d.info
}

// SetReadTimeout sets the read timeout
func (d *USBDevice) SetReadTimeout(timeout time.Duration) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.readTimeout = timeout
}

// SetWriteTimeout sets the write timeout
func (d *USBDevice) SetWriteTimeout(timeout time.Duration) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.writeTimeout = timeout
}

// EnumerateDevices returns a list of connected PM5 devices
func EnumerateDevices() ([]DeviceInfo, error) {
	result := make([]DeviceInfo, 0)

	// Enumerate with a callback function to collect device info
	err := hid.Enumerate(PM5VendorID, 0, func(info *hid.DeviceInfo) error {
		result = append(result, DeviceInfo{
			VendorID:     info.VendorID,
			ProductID:    info.ProductID,
			SerialNumber: info.SerialNbr,
			Product:      info.ProductStr,
			Manufacturer: info.MfrStr,
			Path:         info.Path,
		})
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to enumerate HID devices: %w", err)
	}

	return result, nil
}

// FindFirstPM5 finds and returns the first available PM5 device
func FindFirstPM5() (*USBDevice, error) {
	devices, err := EnumerateDevices()
	if err != nil {
		return nil, err
	}

	if len(devices) == 0 {
		return nil, ErrDeviceNotFound
	}

	return NewUSBDevice(devices[0]), nil
}

// String returns a string representation of the device info
func (d DeviceInfo) String() string {
	return fmt.Sprintf("%s %s (S/N: %s) [VID:0x%04X PID:0x%04X]",
		d.Manufacturer, d.Product, d.SerialNumber,
		d.VendorID, d.ProductID)
}

// MockDevice is a mock HID device for testing
type MockDevice struct {
	info      DeviceInfo
	isOpen    bool
	responses [][]byte
	respIdx   int
	written   [][]byte
	mu        sync.Mutex
}

// NewMockDevice creates a new mock device for testing
func NewMockDevice() *MockDevice {
	return &MockDevice{
		info: DeviceInfo{
			VendorID:     PM5VendorID,
			ProductID:    PM5ProductID,
			SerialNumber: "430000000",
			Product:      "Concept2 Performance Monitor 5 (PM5)",
			Manufacturer: "Concept2",
		},
		responses: make([][]byte, 0),
		written:   make([][]byte, 0),
	}
}

// QueueResponse adds a response to be returned by Read
func (m *MockDevice) QueueResponse(data []byte) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.responses = append(m.responses, data)
}

// GetWritten returns all data written to the device
func (m *MockDevice) GetWritten() [][]byte {
	m.mu.Lock()
	defer m.mu.Unlock()
	result := make([][]byte, len(m.written))
	copy(result, m.written)
	return result
}

func (m *MockDevice) Open() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.isOpen {
		return ErrDeviceAlreadyOpen
	}
	m.isOpen = true
	return nil
}

func (m *MockDevice) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.isOpen = false
	return nil
}

func (m *MockDevice) Write(data []byte) (int, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if !m.isOpen {
		return 0, ErrDeviceNotOpen
	}
	copied := make([]byte, len(data))
	copy(copied, data)
	m.written = append(m.written, copied)
	return len(data), nil
}

func (m *MockDevice) Read(timeout time.Duration) ([]byte, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if !m.isOpen {
		return nil, ErrDeviceNotOpen
	}
	if m.respIdx >= len(m.responses) {
		return nil, ErrTimeout
	}
	data := m.responses[m.respIdx]
	m.respIdx++
	return data, nil
}

func (m *MockDevice) IsOpen() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.isOpen
}

func (m *MockDevice) GetInfo() DeviceInfo {
	return m.info
}
