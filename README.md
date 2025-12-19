# PM5 - Concept2 Performance Monitor Go Library

A comprehensive Go library for communicating with Concept2 PM5 (and PM3/PM4) rowing computers via USB using the CSAFE protocol.

## Features

- **Full CSAFE Protocol Support**: Implements both public CSAFE and Concept2 proprietary command sets
- **USB HID Communication**: Connect to PM5 via USB
- **Frame Encoding/Decoding**: Complete byte-stuffing and checksum handling
- **Workout Configuration**: Set up distance, time, calorie, and interval workouts
- **Real-time Data**: Read pace, power, stroke rate, heart rate, and more
- **Force Curve Data**: Access per-stroke force curve data points
- **Data Utilities**: Pace/Watts conversions, time formatting, and more

## Installation

```bash
go get github.com/concept2/pm5
```

### USB HID Driver Dependency

This library requires a USB HID driver. We recommend one of:

```bash
# Option 1: karalabe/hid (CGO, cross-platform)
go get github.com/karalabe/hid

# Option 2: sstallion/go-hid (CGO, uses hidapi)
go get github.com/sstallion/go-hid
```

## Quick Start

```go
package main

import (
    "fmt"
    "log"
    "github.com/concept2/pm5"
    "github.com/concept2/pm5/device"
)

func main() {
    // Find and connect to PM5
    usbDev, err := device.FindFirstPM5()
    if err != nil {
        log.Fatal(err)
    }

    pm := pm5.New(usbDev)
    if err := pm.Connect(); err != nil {
        log.Fatal(err)
    }
    defer pm.Disconnect()

    // Get device info
    version, _ := pm.GetVersion()
    serial, _ := pm.GetSerial()
    fmt.Printf("Connected to PM%d (S/N: %s)\n", version.Model, serial)

    // Read current workout data
    pace, _ := pm.GetPace()
    power, _ := pm.GetPower()
    strokeRate, _ := pm.GetStrokeRate()
    
    fmt.Printf("Pace: %s  Power: %dW  Rate: %d spm\n",
        pm5.FormatPace(uint32(pace)), power, strokeRate)
}
```

## API Reference

### Connection Management

```go
// Create PM5 instance
pm := pm5.New(device)

// Connect/Disconnect
pm.Connect()
pm.Disconnect()
pm.IsConnected()
```

### Public CSAFE Commands

#### State Control
```go
pm.GoReady()     // Transition to Ready state
pm.GoInUse()     // Transition to InUse state  
pm.GoFinished()  // Transition to Finished state
pm.GoIdle()      // Transition to Idle state
pm.Reset()       // Reset the PM
```

#### Device Information
```go
pm.GetStatus()   // Get status byte with state machine info
pm.GetVersion()  // Get HW/SW version info
pm.GetSerial()   // Get serial number string
```

#### Workout Data (Read)
```go
pm.GetTWork()      // Get work time (hours, minutes, seconds)
pm.GetHorizontal() // Get distance in meters
pm.GetCalories()   // Get total calories
pm.GetPace()       // Get pace (time per 500m in 0.01s)
pm.GetPower()      // Get power in watts
pm.GetCadence()    // Get stroke rate
pm.GetHeartRate()  // Get heart rate (255 = no HR belt)
```

#### Workout Configuration (Write)
```go
pm.SetTWork(hours, minutes, seconds)  // Set time goal
pm.SetHorizontal(meters)              // Set distance goal
pm.SetCalories(cals)                  // Set calorie goal
pm.SetPower(watts)                    // Set power target
pm.SetProgram(workoutNumber)          // Select predefined workout
```

### PM5 Proprietary Commands

#### Configuration
```go
pm.GetFirmwareVersion()     // Detailed 16-byte firmware version
pm.GetHardwareAddress()     // Serial number as uint32
pm.GetErgMachineType()      // Rower/SkiErg/BikeErg type
pm.GetBatteryLevel()        // Battery percentage
pm.GetOperationalState()    // Ready/Workout/Idle/Race/etc
pm.GetWorkoutType()         // JustRow/Fixed/Interval type
pm.GetWorkoutState()        // Current workout phase
pm.GetIntervalType()        // Time/Distance/Rest interval
pm.GetRowingState()         // Active/Inactive
pm.GetStrokeState()         // Drive/Recovery/Waiting
pm.GetWorkoutIntervalCount() // Current interval number
```

#### Real-time Data
```go
pm.GetPMWorkTime()            // Work time in 0.01s
pm.GetPMWorkDistance()        // Distance in 0.1m
pm.GetStroke500mPace()        // Pace in 0.01s per 500m
pm.GetStrokePower()           // Power in watts
pm.GetStrokeCaloricBurnRate() // Calories/hour
pm.GetStrokeRate()            // Strokes per minute
pm.GetDragFactor()            // Drag factor
pm.GetTotalAvg500mPace()      // Average pace
pm.GetTotalAvgPower()         // Average power
pm.GetTotalAvgCalories()      // Total calories
pm.GetAvgHeartRate()          // Average heart rate
pm.GetRestTime()              // Rest time (intervals)
pm.GetErrorValue()            // Last error code
```

#### Stroke Statistics
```go
stats, _ := pm.GetStrokeStats()
// stats.StrokeDistance    (0.01m units)
// stats.DriveTIme         (0.01s units)
// stats.RecoveryTime      (0.01s units)
// stats.StrokeLength      (0.01m units)
// stats.DriveCounter
// stats.PeakDriveForce    (0.1 lbs)
// stats.AvgDriveForce     (0.1 lbs)
// stats.WorkPerStroke     (0.1 Joules)
```

#### Force Curve
```go
// Read force curve data (up to 16 points per call)
// Call during Recovery stroke state
data, _ := pm.GetForcePlotData(32) // Read 32 bytes (16 words)
```

### Workout Setup Helpers

```go
// Simple Just Row
pm.StartJustRowWorkout(withSplits bool)

// Fixed Distance (e.g., 2000m with 500m splits)
pm.StartFixedDistanceWorkout(2000, 500)

// Fixed Time (e.g., 20 minutes with 4 minute splits)
// Time in 0.01 second units: 20min = 120000
pm.StartFixedTimeWorkout(120000, 24000)

// Fixed Calories (e.g., 100 cals with 20 cal splits)
pm.StartFixedCalorieWorkout(100, 20)

// Distance Intervals (e.g., 500m work with 30s rest)
pm.StartFixedDistanceIntervalWorkout(500, 30)

// Time Intervals (e.g., 2:00 work with 30s rest)
pm.StartFixedTimeIntervalWorkout(12000, 30)

// End workout
pm.TerminateWorkout()
pm.GoToMainScreen()
```

### Workout Snapshot

Get a complete snapshot of current workout state:

```go
snapshot, _ := pm.GetWorkoutSnapshot()
fmt.Println(snapshot)
// Output: Time: 5:23.45 | Distance: 1234.5m | Pace: 2:05.3 | Power: 185W | S/R: 24 | HR: 145 | Cals: 89
```

### Data Utilities

```go
// Pace/Power conversions
watts := pm5.PaceToWatts(120.0)      // 2:00 pace → watts
pace := pm5.WattsToPace(200.0)       // 200W → pace in seconds

// Formatting
pm5.FormatPace(12053)                // "2:00.5"
pm5.FormatTime(720000)               // "2:00:00.00"
pm5.FormatDistance(50000)            // "5.00 km"

// Time conversions
duration := pm5.HundredthsToTime(12000)  // → time.Duration
hundredths := pm5.TimeToHundredths(d)    // → uint32
```

## CSAFE Frame Protocol

The library handles all low-level CSAFE protocol details:

- Standard and Extended frame formats
- Byte stuffing (F0-F3 escape sequences)
- Checksum calculation and verification
- Command wrapping for proprietary commands

### Frame Structure

```
Standard:  [F1] [Contents...] [Checksum] [F2]
Extended:  [F0] [Dest] [Src] [Contents...] [Checksum] [F2]
```

### Command Wrappers

Public CSAFE commands use `CmdSetUserCfg1` (0x1A) wrapper.
Proprietary commands use:
- `CmdSetPMCfg` (0x76) - Set configuration
- `CmdSetPMData` (0x77) - Set data
- `CmdGetPMCfg` (0x7E) - Get configuration
- `CmdGetPMData` (0x7F) - Get data

## Type Definitions

See `csafe/types.go` for all enumerated types:

- `OperationalState` - Reset, Ready, Workout, Race, Idle, etc.
- `WorkoutType` - JustRow, FixedDist, FixedTime, Intervals, etc.
- `WorkoutState` - WaitToBegin, WorkoutRow, IntervalRest, etc.
- `IntervalType` - Time, Distance, Rest, Calories, etc.
- `StrokeState` - Driving, Recovery, Waiting, etc.
- `ErgMachineType` - Rower D/E, SkiErg, BikeErg, etc.
- `ScreenType` - Workout, Race, CSAFE, etc.

## USB Connection Notes

### Vendor/Product IDs
- Vendor ID: 0x17A4 (Concept2)
- Product ID: Varies by model/firmware

### HID Reports
- Report ID 1: 20 bytes (+ 1 byte report ID)
- Report ID 2: 120 bytes (+ 1 byte report ID)

### Timing
- Minimum inter-frame gap: 50ms
- Typical response time: <100ms

## Error Handling

```go
resp, err := pm.GetPower()
if err != nil {
    if errors.Is(err, pm5.ErrNotConnected) {
        // Device not connected
    } else if errors.Is(err, pm5.ErrCommandFailed) {
        // Command rejected by PM
    } else if errors.Is(err, pm5.ErrInvalidResponse) {
        // Malformed response
    }
}
```

## Examples

See `example_test.go` for comprehensive examples including:
- Device enumeration and connection
- Reading device information
- Monitoring workout data in real-time
- Setting up different workout types
- Reading force curve data

## Protocol Reference

This library implements the Concept2 PM CSAFE Communication Definition (Rev 0.27, August 2023).

Key sections:
- Frame structure and byte stuffing
- Public CSAFE commands
- Concept2 proprietary commands
- Workout configuration parameters
- Enumerated values and status codes

## License

MIT License - see LICENSE file

## Contributing

Contributions welcome! Please ensure any changes maintain compatibility with the CSAFE protocol specification.
