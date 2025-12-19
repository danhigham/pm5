package pm5

import (
	"github.com/danhigham/pm5/csafe"
)

// ============================================================================
// PM5 Proprietary Get Configuration Commands
// ============================================================================

// FirmwareVersion represents detailed firmware version information
type FirmwareVersion struct {
	Version [16]byte
}

// GetFirmwareVersion returns the PM5 firmware version
func (p *PM5) GetFirmwareVersion() (*FirmwareVersion, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	pmCmd := csafe.BuildCommand(csafe.PMCmdGetFWVersion)
	resp, err := p.sendPMCommand(csafe.CmdGetPMCfg, pmCmd)
	if err != nil {
		return nil, err
	}

	if len(resp.CommandData) < 2 {
		return nil, ErrInvalidResponse
	}

	// The response is in the second command response (after the wrapper)
	for _, cr := range resp.CommandData {
		if cr.Command == csafe.PMCmdGetFWVersion && len(cr.Data) >= 16 {
			fw := &FirmwareVersion{}
			copy(fw.Version[:], cr.Data[:16])
			return fw, nil
		}
	}

	return nil, ErrInvalidResponse
}

// GetHardwareAddress returns the PM5 hardware address (serial number as bytes)
func (p *PM5) GetHardwareAddress() (uint32, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	pmCmd := csafe.BuildCommand(csafe.PMCmdGetHWAddress)
	resp, err := p.sendPMCommand(csafe.CmdGetPMCfg, pmCmd)
	if err != nil {
		return 0, err
	}

	for _, cr := range resp.CommandData {
		if cr.Command == csafe.PMCmdGetHWAddress && len(cr.Data) >= 4 {
			return uint32(cr.Data[0])<<24 | uint32(cr.Data[1])<<16 |
				uint32(cr.Data[2])<<8 | uint32(cr.Data[3]), nil
		}
	}

	return 0, ErrInvalidResponse
}

// GetWorkoutType returns the current workout type
func (p *PM5) GetWorkoutType() (csafe.WorkoutType, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	pmCmd := csafe.BuildCommand(csafe.PMCmdGetWorkoutType)
	resp, err := p.sendPMCommand(csafe.CmdGetPMCfg, pmCmd)
	if err != nil {
		return 0, err
	}

	for _, cr := range resp.CommandData {
		if cr.Command == csafe.PMCmdGetWorkoutType && len(cr.Data) >= 1 {
			return csafe.WorkoutType(cr.Data[0]), nil
		}
	}

	return 0, ErrInvalidResponse
}

// GetWorkoutState returns the current workout state
func (p *PM5) GetWorkoutState() (csafe.WorkoutState, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	pmCmd := csafe.BuildCommand(csafe.PMCmdGetWorkoutState)
	resp, err := p.sendPMCommand(csafe.CmdGetPMCfg, pmCmd)
	if err != nil {
		return 0, err
	}

	for _, cr := range resp.CommandData {
		if cr.Command == csafe.PMCmdGetWorkoutState && len(cr.Data) >= 1 {
			return csafe.WorkoutState(cr.Data[0]), nil
		}
	}

	return 0, ErrInvalidResponse
}

// GetIntervalType returns the current interval type
func (p *PM5) GetIntervalType() (csafe.IntervalType, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	pmCmd := csafe.BuildCommand(csafe.PMCmdGetIntervalType)
	resp, err := p.sendPMCommand(csafe.CmdGetPMCfg, pmCmd)
	if err != nil {
		return 0, err
	}

	for _, cr := range resp.CommandData {
		if cr.Command == csafe.PMCmdGetIntervalType && len(cr.Data) >= 1 {
			return csafe.IntervalType(cr.Data[0]), nil
		}
	}

	return 0, ErrInvalidResponse
}

// GetOperationalState returns the current operational state
func (p *PM5) GetOperationalState() (csafe.OperationalState, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	pmCmd := csafe.BuildCommand(csafe.PMCmdGetOperationalState)
	resp, err := p.sendPMCommand(csafe.CmdGetPMCfg, pmCmd)
	if err != nil {
		return 0, err
	}

	for _, cr := range resp.CommandData {
		if cr.Command == csafe.PMCmdGetOperationalState && len(cr.Data) >= 1 {
			return csafe.OperationalState(cr.Data[0]), nil
		}
	}

	return 0, ErrInvalidResponse
}

// GetRowingState returns the current rowing state
func (p *PM5) GetRowingState() (csafe.RowingState, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	pmCmd := csafe.BuildCommand(csafe.PMCmdGetRowingState)
	resp, err := p.sendPMCommand(csafe.CmdGetPMCfg, pmCmd)
	if err != nil {
		return 0, err
	}

	for _, cr := range resp.CommandData {
		if cr.Command == csafe.PMCmdGetRowingState && len(cr.Data) >= 1 {
			return csafe.RowingState(cr.Data[0]), nil
		}
	}

	return 0, ErrInvalidResponse
}

// GetStrokeState returns the current stroke state
func (p *PM5) GetStrokeState() (csafe.StrokeState, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	pmCmd := csafe.BuildCommand(csafe.PMCmdGetStrokeState)
	resp, err := p.sendPMCommand(csafe.CmdGetPMData, pmCmd)
	if err != nil {
		return 0, err
	}

	for _, cr := range resp.CommandData {
		if cr.Command == csafe.PMCmdGetStrokeState && len(cr.Data) >= 1 {
			return csafe.StrokeState(cr.Data[0]), nil
		}
	}

	return 0, ErrInvalidResponse
}

// GetBatteryLevel returns the battery level percentage
func (p *PM5) GetBatteryLevel() (byte, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	pmCmd := csafe.BuildCommand(csafe.PMCmdGetBatteryLevelPercent)
	resp, err := p.sendPMCommand(csafe.CmdGetPMCfg, pmCmd)
	if err != nil {
		return 0, err
	}

	for _, cr := range resp.CommandData {
		if cr.Command == csafe.PMCmdGetBatteryLevelPercent && len(cr.Data) >= 1 {
			return cr.Data[0], nil
		}
	}

	return 0, ErrInvalidResponse
}

// GetErgMachineType returns the connected erg machine type
func (p *PM5) GetErgMachineType() (csafe.ErgMachineType, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	pmCmd := csafe.BuildCommand(csafe.PMCmdGetErgMachineType)
	resp, err := p.sendPMCommand(csafe.CmdGetPMCfg, pmCmd)
	if err != nil {
		return 0, err
	}

	for _, cr := range resp.CommandData {
		if cr.Command == csafe.PMCmdGetErgMachineType && len(cr.Data) >= 1 {
			return csafe.ErgMachineType(cr.Data[0]), nil
		}
	}

	return 0, ErrInvalidResponse
}

// GetWorkoutIntervalCount returns the current interval count
func (p *PM5) GetWorkoutIntervalCount() (byte, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	pmCmd := csafe.BuildCommand(csafe.PMCmdGetWorkoutIntervalCount)
	resp, err := p.sendPMCommand(csafe.CmdGetPMCfg, pmCmd)
	if err != nil {
		return 0, err
	}

	for _, cr := range resp.CommandData {
		if cr.Command == csafe.PMCmdGetWorkoutIntervalCount && len(cr.Data) >= 1 {
			return cr.Data[0], nil
		}
	}

	return 0, ErrInvalidResponse
}

// ============================================================================
// PM5 Proprietary Get Data Commands
// ============================================================================

// GetPMWorkTime returns detailed work time in hundredths of seconds
func (p *PM5) GetPMWorkTime() (uint32, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	pmCmd := csafe.BuildCommand(csafe.PMCmdGetWorkTime)
	resp, err := p.sendPMCommand(csafe.CmdGetPMData, pmCmd)
	if err != nil {
		return 0, err
	}

	for _, cr := range resp.CommandData {
		if cr.Command == csafe.PMCmdGetWorkTime && len(cr.Data) >= 4 {
			return uint32(cr.Data[0])<<24 | uint32(cr.Data[1])<<16 |
				uint32(cr.Data[2])<<8 | uint32(cr.Data[3]), nil
		}
	}

	return 0, ErrInvalidResponse
}

// GetPMWorkDistance returns the work distance in tenths of meters
func (p *PM5) GetPMWorkDistance() (uint32, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	pmCmd := csafe.BuildCommand(csafe.PMCmdGetWorkDistance)
	resp, err := p.sendPMCommand(csafe.CmdGetPMData, pmCmd)
	if err != nil {
		return 0, err
	}

	for _, cr := range resp.CommandData {
		if cr.Command == csafe.PMCmdGetWorkDistance && len(cr.Data) >= 4 {
			return uint32(cr.Data[0])<<24 | uint32(cr.Data[1])<<16 |
				uint32(cr.Data[2])<<8 | uint32(cr.Data[3]), nil
		}
	}

	return 0, ErrInvalidResponse
}

// GetStroke500mPace returns the current pace per 500m in hundredths of seconds
func (p *PM5) GetStroke500mPace() (uint32, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	pmCmd := csafe.BuildCommand(csafe.PMCmdGetStroke500mPace)
	resp, err := p.sendPMCommand(csafe.CmdGetPMData, pmCmd)
	if err != nil {
		return 0, err
	}

	for _, cr := range resp.CommandData {
		if cr.Command == csafe.PMCmdGetStroke500mPace && len(cr.Data) >= 4 {
			return uint32(cr.Data[0])<<24 | uint32(cr.Data[1])<<16 |
				uint32(cr.Data[2])<<8 | uint32(cr.Data[3]), nil
		}
	}

	return 0, ErrInvalidResponse
}

// GetStrokePower returns the current stroke power in watts
func (p *PM5) GetStrokePower() (uint32, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	pmCmd := csafe.BuildCommand(csafe.PMCmdGetStrokePower)
	resp, err := p.sendPMCommand(csafe.CmdGetPMData, pmCmd)
	if err != nil {
		return 0, err
	}

	for _, cr := range resp.CommandData {
		if cr.Command == csafe.PMCmdGetStrokePower && len(cr.Data) >= 4 {
			return uint32(cr.Data[0])<<24 | uint32(cr.Data[1])<<16 |
				uint32(cr.Data[2])<<8 | uint32(cr.Data[3]), nil
		}
	}

	return 0, ErrInvalidResponse
}

// GetStrokeCaloricBurnRate returns the stroke caloric burn rate in cals/hr
func (p *PM5) GetStrokeCaloricBurnRate() (uint32, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	pmCmd := csafe.BuildCommand(csafe.PMCmdGetStrokeCaloricBurnRate)
	resp, err := p.sendPMCommand(csafe.CmdGetPMData, pmCmd)
	if err != nil {
		return 0, err
	}

	for _, cr := range resp.CommandData {
		if cr.Command == csafe.PMCmdGetStrokeCaloricBurnRate && len(cr.Data) >= 4 {
			return uint32(cr.Data[0])<<24 | uint32(cr.Data[1])<<16 |
				uint32(cr.Data[2])<<8 | uint32(cr.Data[3]), nil
		}
	}

	return 0, ErrInvalidResponse
}

// GetStrokeRate returns the current stroke rate (strokes per minute)
func (p *PM5) GetStrokeRate() (byte, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	pmCmd := csafe.BuildCommand(csafe.PMCmdGetStrokeRate)
	resp, err := p.sendPMCommand(csafe.CmdGetPMData, pmCmd)
	if err != nil {
		return 0, err
	}

	for _, cr := range resp.CommandData {
		if cr.Command == csafe.PMCmdGetStrokeRate && len(cr.Data) >= 1 {
			return cr.Data[0], nil
		}
	}

	return 0, ErrInvalidResponse
}

// GetDragFactor returns the current drag factor
func (p *PM5) GetDragFactor() (byte, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	pmCmd := csafe.BuildCommand(csafe.PMCmdGetDragFactor)
	resp, err := p.sendPMCommand(csafe.CmdGetPMData, pmCmd)
	if err != nil {
		return 0, err
	}

	for _, cr := range resp.CommandData {
		if cr.Command == csafe.PMCmdGetDragFactor && len(cr.Data) >= 1 {
			return cr.Data[0], nil
		}
	}

	return 0, ErrInvalidResponse
}

// GetTotalAvg500mPace returns the total average pace per 500m
func (p *PM5) GetTotalAvg500mPace() (uint32, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	pmCmd := csafe.BuildCommand(csafe.PMCmdGetTotalAvg500mPace)
	resp, err := p.sendPMCommand(csafe.CmdGetPMData, pmCmd)
	if err != nil {
		return 0, err
	}

	for _, cr := range resp.CommandData {
		if cr.Command == csafe.PMCmdGetTotalAvg500mPace && len(cr.Data) >= 4 {
			return uint32(cr.Data[0])<<24 | uint32(cr.Data[1])<<16 |
				uint32(cr.Data[2])<<8 | uint32(cr.Data[3]), nil
		}
	}

	return 0, ErrInvalidResponse
}

// GetTotalAvgPower returns the total average power in watts
func (p *PM5) GetTotalAvgPower() (uint32, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	pmCmd := csafe.BuildCommand(csafe.PMCmdGetTotalAvgPower)
	resp, err := p.sendPMCommand(csafe.CmdGetPMData, pmCmd)
	if err != nil {
		return 0, err
	}

	for _, cr := range resp.CommandData {
		if cr.Command == csafe.PMCmdGetTotalAvgPower && len(cr.Data) >= 4 {
			return uint32(cr.Data[0])<<24 | uint32(cr.Data[1])<<16 |
				uint32(cr.Data[2])<<8 | uint32(cr.Data[3]), nil
		}
	}

	return 0, ErrInvalidResponse
}

// GetTotalAvgCalories returns the total calories burned
func (p *PM5) GetTotalAvgCalories() (uint32, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	pmCmd := csafe.BuildCommand(csafe.PMCmdGetTotalAvgCalories)
	resp, err := p.sendPMCommand(csafe.CmdGetPMData, pmCmd)
	if err != nil {
		return 0, err
	}

	for _, cr := range resp.CommandData {
		if cr.Command == csafe.PMCmdGetTotalAvgCalories && len(cr.Data) >= 4 {
			return uint32(cr.Data[0])<<24 | uint32(cr.Data[1])<<16 |
				uint32(cr.Data[2])<<8 | uint32(cr.Data[3]), nil
		}
	}

	return 0, ErrInvalidResponse
}

// GetAvgHeartRate returns the average heart rate
func (p *PM5) GetAvgHeartRate() (byte, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	pmCmd := csafe.BuildCommand(csafe.PMCmdGetAvgHeartRate)
	resp, err := p.sendPMCommand(csafe.CmdGetPMData, pmCmd)
	if err != nil {
		return 0, err
	}

	for _, cr := range resp.CommandData {
		if cr.Command == csafe.PMCmdGetAvgHeartRate && len(cr.Data) >= 1 {
			return cr.Data[0], nil
		}
	}

	return 0, ErrInvalidResponse
}

// StrokeStats contains detailed stroke statistics
type StrokeStats struct {
	StrokeDistance    uint16 // 0.01m units
	DriveTIme         byte   // 0.01s units
	RecoveryTime      uint16 // 0.01s units
	StrokeLength      byte   // 0.01m units
	DriveCounter      uint16
	PeakDriveForce    uint16 // 0.1 lbs
	ImpulseDriveForce uint16 // 0.1 lbs
	AvgDriveForce     uint16 // 0.1 lbs
	WorkPerStroke     uint16 // 0.1 Joules
}

// GetStrokeStats returns detailed stroke statistics
func (p *PM5) GetStrokeStats() (*StrokeStats, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	pmCmd := csafe.BuildCommand(csafe.PMCmdGetStrokeStats, 0x00)
	resp, err := p.sendPMCommand(csafe.CmdGetPMData, pmCmd)
	if err != nil {
		return nil, err
	}

	for _, cr := range resp.CommandData {
		if cr.Command == csafe.PMCmdGetStrokeStats && len(cr.Data) >= 16 {
			d := cr.Data
			return &StrokeStats{
				StrokeDistance:    uint16(d[0])<<8 | uint16(d[1]),
				DriveTIme:         d[2],
				RecoveryTime:      uint16(d[3])<<8 | uint16(d[4]),
				StrokeLength:      d[5],
				DriveCounter:      uint16(d[6])<<8 | uint16(d[7]),
				PeakDriveForce:    uint16(d[8])<<8 | uint16(d[9]),
				ImpulseDriveForce: uint16(d[10])<<8 | uint16(d[11]),
				AvgDriveForce:     uint16(d[12])<<8 | uint16(d[13]),
				WorkPerStroke:     uint16(d[14])<<8 | uint16(d[15]),
			}, nil
		}
	}

	return nil, ErrInvalidResponse
}

// GetForcePlotData returns force curve data points
// blockSize is the number of bytes to read (max 32, returns 16 words)
func (p *PM5) GetForcePlotData(blockSize byte) ([]uint16, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if blockSize > 32 {
		blockSize = 32
	}

	pmCmd := csafe.BuildCommand(csafe.PMCmdGetForcePlotData, blockSize)
	resp, err := p.sendPMCommand(csafe.CmdSetUserCfg1, pmCmd)
	if err != nil {
		return nil, err
	}

	for _, cr := range resp.CommandData {
		if cr.Command == csafe.PMCmdGetForcePlotData && len(cr.Data) >= 1 {
			bytesRead := cr.Data[0]
			if bytesRead == 0 {
				return []uint16{}, nil
			}

			// Parse word data (big-endian pairs)
			numWords := int(bytesRead) / 2
			if numWords > 16 {
				numWords = 16
			}

			words := make([]uint16, numWords)
			for i := 0; i < numWords && 1+i*2+1 < len(cr.Data); i++ {
				words[i] = uint16(cr.Data[1+i*2])<<8 | uint16(cr.Data[1+i*2+1])
			}
			return words, nil
		}
	}

	return nil, ErrInvalidResponse
}

// GetRestTime returns the current rest time in hundredths of seconds
func (p *PM5) GetRestTime() (uint16, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	pmCmd := csafe.BuildCommand(csafe.PMCmdGetRestTime)
	resp, err := p.sendPMCommand(csafe.CmdGetPMData, pmCmd)
	if err != nil {
		return 0, err
	}

	for _, cr := range resp.CommandData {
		if cr.Command == csafe.PMCmdGetRestTime && len(cr.Data) >= 2 {
			return uint16(cr.Data[0]) | uint16(cr.Data[1])<<8, nil
		}
	}

	return 0, ErrInvalidResponse
}

// GetErrorValue returns the last error value
func (p *PM5) GetErrorValue() (uint16, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	pmCmd := csafe.BuildCommand(csafe.PMCmdGetErrorValue)
	resp, err := p.sendPMCommand(csafe.CmdGetPMData, pmCmd)
	if err != nil {
		return 0, err
	}

	for _, cr := range resp.CommandData {
		if cr.Command == csafe.PMCmdGetErrorValue && len(cr.Data) >= 2 {
			return uint16(cr.Data[0])<<8 | uint16(cr.Data[1]), nil
		}
	}

	return 0, ErrInvalidResponse
}
