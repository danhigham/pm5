package pm5

import (
	"fmt"
	"math"
	"time"

	"github.com/danhigham/pm5/csafe"
)

// ============================================================================
// Data Conversion Utilities
// ============================================================================

// Pace conversion constants
const (
	// Concept2 uses this formula: Pace = (2.8/(Watts/Watts_ref))^(1/3) * 500
	// where Watts_ref = 2.8 (power needed for 500m in 500s = 2:00/500m)
	WattsRef = 2.8
)

// PaceToWatts converts pace (seconds per 500m) to watts
func PaceToWatts(paceSeconds float64) float64 {
	if paceSeconds <= 0 {
		return 0
	}
	// Formula: Watts = 2.8 / (pace/500)^3
	pace500 := paceSeconds / 500.0
	return WattsRef / math.Pow(pace500, 3)
}

// WattsToPace converts watts to pace (seconds per 500m)
func WattsToPace(watts float64) float64 {
	if watts <= 0 {
		return 0
	}
	// Formula: pace = 500 * (2.8/Watts)^(1/3)
	return 500.0 * math.Pow(WattsRef/watts, 1.0/3.0)
}

// CaloriesPerHourToPace converts calories per hour to pace (seconds per 500m)
func CaloriesPerHourToPace(calsPerHour float64) float64 {
	if calsPerHour <= 0 {
		return 0
	}
	// Formula from Concept2:
	// cal/hr = (watts * 4 + 350) / 0.8604
	// Solving for watts: watts = (cal/hr * 0.8604 - 350) / 4
	watts := (calsPerHour*0.8604 - 350.0) / 4.0
	if watts <= 0 {
		return 0
	}
	return WattsToPace(watts)
}

// PaceToCaloriesPerHour converts pace (seconds per 500m) to calories per hour
func PaceToCaloriesPerHour(paceSeconds float64) float64 {
	watts := PaceToWatts(paceSeconds)
	if watts <= 0 {
		return 0
	}
	// cal/hr = (watts * 4 + 350) / 0.8604
	return (watts*4.0 + 350.0) / 0.8604
}

// HundredthsToTime converts hundredths of seconds to a time.Duration
func HundredthsToTime(hundredths uint32) time.Duration {
	return time.Duration(hundredths) * 10 * time.Millisecond
}

// TimeToHundredths converts a time.Duration to hundredths of seconds
func TimeToHundredths(d time.Duration) uint32 {
	return uint32(d.Milliseconds() / 10)
}

// TenthsToMeters converts tenths of meters to meters as float
func TenthsToMeters(tenths uint32) float64 {
	return float64(tenths) / 10.0
}

// MetersToTenths converts meters to tenths of meters
func MetersToTenths(meters float64) uint32 {
	return uint32(meters * 10)
}

// FormatPace formats pace in hundredths of seconds as M:SS.t
func FormatPace(hundredths uint32) string {
	totalSeconds := float64(hundredths) / 100.0
	minutes := int(totalSeconds) / 60
	seconds := totalSeconds - float64(minutes*60)
	return fmt.Sprintf("%d:%04.1f", minutes, seconds)
}

// FormatTime formats time in hundredths of seconds as H:MM:SS.hh
func FormatTime(hundredths uint32) string {
	totalSeconds := hundredths / 100
	remaining := hundredths % 100
	hours := totalSeconds / 3600
	minutes := (totalSeconds % 3600) / 60
	seconds := totalSeconds % 60

	if hours > 0 {
		return fmt.Sprintf("%d:%02d:%02d.%02d", hours, minutes, seconds, remaining)
	}
	return fmt.Sprintf("%d:%02d.%02d", minutes, seconds, remaining)
}

// FormatDistance formats distance in tenths of meters
func FormatDistance(tenths uint32) string {
	meters := float64(tenths) / 10.0
	if meters >= 1000 {
		return fmt.Sprintf("%.2f km", meters/1000)
	}
	return fmt.Sprintf("%.1f m", meters)
}

// ============================================================================
// Multi-byte Data Construction (Little-Endian)
// ============================================================================

// Uint16ToBytes converts uint16 to little-endian byte slice
func Uint16ToBytes(v uint16) []byte {
	return []byte{byte(v & 0xFF), byte((v >> 8) & 0xFF)}
}

// Uint24ToBytes converts uint32 (24-bit value) to little-endian byte slice
func Uint24ToBytes(v uint32) []byte {
	return []byte{
		byte(v & 0xFF),
		byte((v >> 8) & 0xFF),
		byte((v >> 16) & 0xFF),
	}
}

// Uint32ToBytes converts uint32 to little-endian byte slice
func Uint32ToBytes(v uint32) []byte {
	return []byte{
		byte(v & 0xFF),
		byte((v >> 8) & 0xFF),
		byte((v >> 16) & 0xFF),
		byte((v >> 24) & 0xFF),
	}
}

// ============================================================================
// Multi-byte Data Deconstruction (Little-Endian)
// ============================================================================

// BytesToUint16 converts little-endian byte slice to uint16
func BytesToUint16(b []byte) uint16 {
	if len(b) < 2 {
		return 0
	}
	return uint16(b[0]) | uint16(b[1])<<8
}

// BytesToUint24 converts little-endian byte slice to uint32 (24-bit)
func BytesToUint24(b []byte) uint32 {
	if len(b) < 3 {
		return 0
	}
	return uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2])<<16
}

// BytesToUint32 converts little-endian byte slice to uint32
func BytesToUint32(b []byte) uint32 {
	if len(b) < 4 {
		return 0
	}
	return uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2])<<16 | uint32(b[3])<<24
}

// BytesToUint16BE converts big-endian byte slice to uint16
func BytesToUint16BE(b []byte) uint16 {
	if len(b) < 2 {
		return 0
	}
	return uint16(b[0])<<8 | uint16(b[1])
}

// BytesToUint32BE converts big-endian byte slice to uint32
func BytesToUint32BE(b []byte) uint32 {
	if len(b) < 4 {
		return 0
	}
	return uint32(b[0])<<24 | uint32(b[1])<<16 | uint32(b[2])<<8 | uint32(b[3])
}

// ============================================================================
// Workout Status Snapshot
// ============================================================================

// WorkoutSnapshot contains a complete snapshot of the current workout state
type WorkoutSnapshot struct {
	// Timing
	ElapsedTime   time.Duration // Total elapsed time
	WorkTime      time.Duration // Active work time
	RestTime      time.Duration // Rest time (intervals)
	ProjectedTime time.Duration // Projected finish time

	// Distance
	Distance          float64 // Meters
	ProjectedDistance float64 // Meters

	// Performance
	Pace          time.Duration // Per 500m
	AvgPace       time.Duration // Per 500m
	Power         uint32        // Watts
	AvgPower      uint32        // Watts
	StrokeRate    byte          // Strokes per minute
	AvgStrokeRate byte          // Strokes per minute
	DragFactor    byte

	// Calories
	Calories        uint32
	CaloricBurnRate uint16 // Cals/hr

	// Heart Rate
	HeartRate    byte // BPM (255 = invalid)
	AvgHeartRate byte

	// State
	WorkoutType   string
	WorkoutState  string
	IntervalType  string
	RowingState   string
	StrokeState   string
	IntervalCount byte
}

// GetWorkoutSnapshot returns a complete snapshot of the current workout
// This uses a single batched CSAFE command for efficiency
func (p *PM5) GetWorkoutSnapshot() (*WorkoutSnapshot, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	snapshot := &WorkoutSnapshot{}

	// Build all PM commands in a single batch
	pmCmds := [][]byte{
		csafe.BuildCommand(csafe.PMCmdGetWorkoutType),          // 0
		csafe.BuildCommand(csafe.PMCmdGetWorkoutState),         // 1
		csafe.BuildCommand(csafe.PMCmdGetIntervalType),         // 2
		csafe.BuildCommand(csafe.PMCmdGetRowingState),          // 3
		csafe.BuildCommand(csafe.PMCmdGetStrokeState),          // 4
		csafe.BuildCommand(csafe.PMCmdGetWorkoutIntervalCount), // 5
		csafe.BuildCommand(csafe.PMCmdGetWorkTime),             // 6
		csafe.BuildCommand(csafe.PMCmdGetWorkDistance),         // 7
		csafe.BuildCommand(csafe.PMCmdGetStroke500mPace),       // 8
		csafe.BuildCommand(csafe.PMCmdGetTotalAvg500mPace),     // 9
		csafe.BuildCommand(csafe.PMCmdGetStrokePower),          // 10
		csafe.BuildCommand(csafe.PMCmdGetTotalAvgPower),        // 11
		csafe.BuildCommand(csafe.PMCmdGetStrokeRate),           // 12
		csafe.BuildCommand(csafe.PMCmdGetDragFactor),           // 13
		csafe.BuildCommand(csafe.PMCmdGetTotalAvgCalories),     // 14
		csafe.BuildCommand(csafe.PMCmdGetAvgHeartRate),         // 15
	}

	// Add standard CSAFE heart rate command (not PM-specific)
	contents := csafe.BuildPMCommand(csafe.CmdGetPMData, pmCmds...)
	contents = append(contents, csafe.CmdGetHRCur)

	resp, err := p.sendCommand(contents)
	if err != nil {
		return nil, err
	}

	// Parse the batched response
	for _, cmdResp := range resp.CommandData {
		// Handle standard CSAFE commands
		if cmdResp.Command == csafe.CmdGetHRCur && len(cmdResp.Data) >= 1 {
			snapshot.HeartRate = cmdResp.Data[0]
			continue
		}

		// Handle PM proprietary responses
		for _, pmResp := range cmdResp.PMResponses {
			switch pmResp.Command {
			case csafe.PMCmdGetWorkoutType:
				if len(pmResp.Data) >= 1 {
					snapshot.WorkoutType = csafe.WorkoutType(pmResp.Data[0]).String()
				}
			case csafe.PMCmdGetWorkoutState:
				if len(pmResp.Data) >= 1 {
					snapshot.WorkoutState = csafe.WorkoutState(pmResp.Data[0]).String()
				}
			case csafe.PMCmdGetIntervalType:
				if len(pmResp.Data) >= 1 {
					snapshot.IntervalType = csafe.IntervalType(pmResp.Data[0]).String()
				}
			case csafe.PMCmdGetRowingState:
				if len(pmResp.Data) >= 1 {
					snapshot.RowingState = csafe.RowingState(pmResp.Data[0]).String()
				}
			case csafe.PMCmdGetStrokeState:
				if len(pmResp.Data) >= 1 {
					snapshot.StrokeState = csafe.StrokeState(pmResp.Data[0]).String()
				}
			case csafe.PMCmdGetWorkoutIntervalCount:
				if len(pmResp.Data) >= 1 {
					snapshot.IntervalCount = pmResp.Data[0]
				}
			case csafe.PMCmdGetWorkTime:
				if len(pmResp.Data) >= 4 {
					wt := BytesToUint32BE(pmResp.Data[:4])
					snapshot.WorkTime = HundredthsToTime(wt)
					snapshot.ElapsedTime = snapshot.WorkTime
				}
			case csafe.PMCmdGetWorkDistance:
				if len(pmResp.Data) >= 4 {
					snapshot.Distance = float64(BytesToUint32BE(pmResp.Data[:4]))
				}
			case csafe.PMCmdGetStroke500mPace:
				if len(pmResp.Data) >= 4 {
					pace := BytesToUint32BE(pmResp.Data[:4])
					snapshot.Pace = HundredthsToTime(pace)
				}
			case csafe.PMCmdGetTotalAvg500mPace:
				if len(pmResp.Data) >= 4 {
					avgPace := BytesToUint32BE(pmResp.Data[:4])
					snapshot.AvgPace = HundredthsToTime(avgPace)
				}
			case csafe.PMCmdGetStrokePower:
				if len(pmResp.Data) >= 4 {
					snapshot.Power = BytesToUint32BE(pmResp.Data[:4])
				}
			case csafe.PMCmdGetTotalAvgPower:
				if len(pmResp.Data) >= 4 {
					snapshot.AvgPower = BytesToUint32BE(pmResp.Data[:4])
				}
			case csafe.PMCmdGetStrokeRate:
				if len(pmResp.Data) >= 1 {
					snapshot.StrokeRate = pmResp.Data[0]
				}
			case csafe.PMCmdGetDragFactor:
				if len(pmResp.Data) >= 1 {
					snapshot.DragFactor = pmResp.Data[0]
				}
			case csafe.PMCmdGetTotalAvgCalories:
				if len(pmResp.Data) >= 4 {
					snapshot.Calories = BytesToUint32BE(pmResp.Data[:4])
				}
			case csafe.PMCmdGetAvgHeartRate:
				if len(pmResp.Data) >= 1 {
					snapshot.AvgHeartRate = pmResp.Data[0]
				}
			}
		}
	}

	return snapshot, nil
}

// String returns a formatted string representation of the workout snapshot
func (s *WorkoutSnapshot) String() string {
	return fmt.Sprintf(
		"Time: %s | Distance: %.1fm | Pace: %s | Power: %dW | S/R: %d | HR: %d | Cals: %d",
		FormatTime(TimeToHundredths(s.WorkTime)),
		s.Distance,
		FormatPace(TimeToHundredths(s.Pace)),
		s.Power,
		s.StrokeRate,
		s.HeartRate,
		s.Calories,
	)
}
