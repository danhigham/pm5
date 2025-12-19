//go:build ignore
// +build ignore

// This file demonstrates how to use the pm5 library to communicate with
// a Concept2 PM5 rowing computer.
//
// Run with: go run example.go
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/danhigham/pm5"
	"github.com/danhigham/pm5/csafe"
	"github.com/danhigham/pm5/device"
)

func main() {
	fmt.Println("PM5 Library Example")
	fmt.Println("===================")

	// Find connected PM5 devices
	devices, err := device.EnumerateDevices()
	if err != nil {
		log.Fatalf("Failed to enumerate devices: %v", err)
	}

	if len(devices) == 0 {
		fmt.Println("No PM5 devices found. Running with mock device...")
		runWithMockDevice()
		return
	}

	// Use the first device found
	fmt.Printf("Found device: %s\n", devices[0])

	usbDev := device.NewUSBDevice(devices[0])
	pm := pm5.New(usbDev)

	if err := pm.Connect(); err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer pm.Disconnect()

	runExamples(pm)
}

func runWithMockDevice() {
	// Create a mock device for testing/demonstration
	mockDev := device.NewMockDevice()

	// Queue some mock responses
	// These would be actual CSAFE frame responses from a real device

	pm := pm5.New(mockDev)

	if err := pm.Connect(); err != nil {
		log.Fatalf("Failed to connect to mock: %v", err)
	}
	defer pm.Disconnect()

	fmt.Println("\nMock device connected successfully!")
	fmt.Println("In a real scenario, the following operations would communicate with the PM5.")

	demonstrateAPI()
}

func runExamples(pm *pm5.PM5) {
	fmt.Println("\n--- Basic Information ---")

	// Get version information
	version, err := pm.GetVersion()
	if err != nil {
		log.Printf("Failed to get version: %v", err)
	} else {
		fmt.Printf("Manufacturer ID: %d\n", version.ManufacturerID)
		fmt.Printf("Class ID: %d\n", version.ClassID)
		fmt.Printf("Model: PM%d\n", version.Model)
		fmt.Printf("HW Version: %d.%d\n", version.HWVersionHigh, version.HWVersionLow)
		fmt.Printf("SW Version: %d.%d\n", version.SWVersionHigh, version.SWVersionLow)
	}

	// Get serial number
	serial, err := pm.GetSerial()
	if err != nil {
		log.Printf("Failed to get serial: %v", err)
	} else {
		fmt.Printf("Serial Number: %s\n", serial)
	}

	// Get erg machine type
	ergType, err := pm.GetErgMachineType()
	if err != nil {
		log.Printf("Failed to get erg type: %v", err)
	} else {
		fmt.Printf("Machine Type: %s\n", ergType)
	}

	// Get battery level
	battery, err := pm.GetBatteryLevel()
	if err != nil {
		log.Printf("Failed to get battery: %v", err)
	} else {
		fmt.Printf("Battery Level: %d%%\n", battery)
	}

	fmt.Println("\n--- State Information ---")

	// Get operational state
	opState, err := pm.GetOperationalState()
	if err != nil {
		log.Printf("Failed to get operational state: %v", err)
	} else {
		fmt.Printf("Operational State: %s\n", opState)
	}

	// Get workout state
	workoutState, err := pm.GetWorkoutState()
	if err != nil {
		log.Printf("Failed to get workout state: %v", err)
	} else {
		fmt.Printf("Workout State: %s\n", workoutState)
	}

	// Get rowing state
	rowingState, err := pm.GetRowingState()
	if err != nil {
		log.Printf("Failed to get rowing state: %v", err)
	} else {
		fmt.Printf("Rowing State: %s\n", rowingState)
	}

	fmt.Println("\n--- Current Data ---")

	// Get current metrics
	power, _ := pm.GetPower()
	fmt.Printf("Power: %d W\n", power)

	pace, _ := pm.GetPace()
	fmt.Printf("Pace: %s /500m\n", pm5.FormatPace(uint32(pace)))

	strokeRate, _ := pm.GetStrokeRate()
	fmt.Printf("Stroke Rate: %d spm\n", strokeRate)

	hr, _ := pm.GetHeartRate()
	if hr == 255 {
		fmt.Println("Heart Rate: No HR belt connected")
	} else {
		fmt.Printf("Heart Rate: %d bpm\n", hr)
	}

	dragFactor, _ := pm.GetDragFactor()
	fmt.Printf("Drag Factor: %d\n", dragFactor)

	fmt.Println("\n--- Workout Snapshot ---")
	snapshot, err := pm.GetWorkoutSnapshot()
	if err != nil {
		log.Printf("Failed to get snapshot: %v", err)
	} else {
		fmt.Println(snapshot)
	}
}

func demonstrateAPI() {
	fmt.Println("\n=== API Demonstration ===")
	fmt.Println("\nThe PM5 library provides the following capabilities:")

	fmt.Println("\n1. BASIC PUBLIC CSAFE COMMANDS:")
	fmt.Println("   pm.GetStatus()       - Get current status")
	fmt.Println("   pm.GetVersion()      - Get firmware/hardware version")
	fmt.Println("   pm.GetSerial()       - Get serial number")
	fmt.Println("   pm.GoReady()         - Transition to Ready state")
	fmt.Println("   pm.GoInUse()         - Transition to InUse state")
	fmt.Println("   pm.GoFinished()      - Transition to Finished state")
	fmt.Println("   pm.Reset()           - Reset the PM")

	fmt.Println("\n2. GET WORKOUT DATA:")
	fmt.Println("   pm.GetTWork()        - Get work time (H:M:S)")
	fmt.Println("   pm.GetHorizontal()   - Get distance (meters)")
	fmt.Println("   pm.GetCalories()     - Get total calories")
	fmt.Println("   pm.GetPace()         - Get current pace (/500m)")
	fmt.Println("   pm.GetPower()        - Get current power (watts)")
	fmt.Println("   pm.GetCadence()      - Get stroke rate")
	fmt.Println("   pm.GetHeartRate()    - Get heart rate (bpm)")

	fmt.Println("\n3. SET WORKOUT GOALS:")
	fmt.Println("   pm.SetTWork(h,m,s)   - Set time goal")
	fmt.Println("   pm.SetHorizontal(m)  - Set distance goal (meters)")
	fmt.Println("   pm.SetCalories(c)    - Set calorie goal")
	fmt.Println("   pm.SetPower(w)       - Set power target (watts)")
	fmt.Println("   pm.SetProgram(n)     - Select predefined workout")

	fmt.Println("\n4. PROPRIETARY PM5 COMMANDS:")
	fmt.Println("   pm.GetFirmwareVersion()      - Detailed firmware info")
	fmt.Println("   pm.GetErgMachineType()       - Rower/SkiErg/BikeErg")
	fmt.Println("   pm.GetBatteryLevel()         - Battery percentage")
	fmt.Println("   pm.GetOperationalState()     - Ready/Workout/Idle/etc")
	fmt.Println("   pm.GetWorkoutType()          - JustRow/Fixed/Interval")
	fmt.Println("   pm.GetWorkoutState()         - Current workout phase")
	fmt.Println("   pm.GetRowingState()          - Active/Inactive")
	fmt.Println("   pm.GetStrokeState()          - Drive/Recovery phase")
	fmt.Println("   pm.GetStrokeStats()          - Detailed stroke data")
	fmt.Println("   pm.GetForcePlotData(n)       - Force curve points")

	fmt.Println("\n5. WORKOUT SETUP HELPERS:")
	fmt.Println("   pm.StartJustRowWorkout(splits)            - Simple row")
	fmt.Println("   pm.StartFixedDistanceWorkout(m, split)    - e.g., 2000m")
	fmt.Println("   pm.StartFixedTimeWorkout(cs, split)       - e.g., 20:00")
	fmt.Println("   pm.StartFixedCalorieWorkout(c, split)     - e.g., 100 cals")
	fmt.Println("   pm.StartFixedDistanceIntervalWorkout(...) - Distance intervals")
	fmt.Println("   pm.StartFixedTimeIntervalWorkout(...)     - Time intervals")
	fmt.Println("   pm.TerminateWorkout()                     - End workout")

	fmt.Println("\n6. DATA UTILITIES:")
	fmt.Println("   pm5.WattsToPace(w)           - Convert watts to pace")
	fmt.Println("   pm5.PaceToWatts(s)           - Convert pace to watts")
	fmt.Println("   pm5.FormatPace(cs)           - Format as M:SS.t")
	fmt.Println("   pm5.FormatTime(cs)           - Format as H:MM:SS.hh")
	fmt.Println("   pm5.FormatDistance(dm)       - Format with units")

	fmt.Println("\n=== Example: Setting up a 2000m workout ===")
	fmt.Println(`
    pm := pm5.New(device)
    pm.Connect()
    defer pm.Disconnect()

    // Option 1: Using public CSAFE commands
    pm.SetHorizontal(2000)                        // 2000 meters
    pm.SetProgram(csafe.WorkoutNumberProgrammed)  // Use programmed workout
    pm.GoInUse()                                  // Start workout

    // Option 2: Using proprietary helper
    pm.StartFixedDistanceWorkout(2000, 500)       // 2000m with 500m splits

    // Monitor progress
    for {
        snapshot, _ := pm.GetWorkoutSnapshot()
        fmt.Println(snapshot)
        
        if snapshot.WorkoutState == "Workout End" {
            break
        }
        time.Sleep(1 * time.Second)
    }`)

	fmt.Println("\n=== Data Conversions ===")

	// Demonstrate pace/watts conversions
	watts := 200.0
	pace := pm5.WattsToPace(watts)
	fmt.Printf("\n%.0f watts = %s /500m pace\n", watts, pm5.FormatPace(uint32(pace*100)))

	paceSeconds := 120.0 // 2:00 /500m
	wattsCalc := pm5.PaceToWatts(paceSeconds)
	fmt.Printf("2:00 /500m pace = %.1f watts\n", wattsCalc)

	// Demonstrate time formatting
	hundredths := uint32(720000) // 2 hours
	fmt.Printf("\n%d hundredths = %s\n", hundredths, pm5.FormatTime(hundredths))

	// Demonstrate distance formatting
	tenths := uint32(50000) // 5000 meters
	fmt.Printf("%d tenths of meters = %s\n", tenths, pm5.FormatDistance(tenths))
}

// Example: Monitor a workout in real-time
func monitorWorkout(pm *pm5.PM5) {
	fmt.Println("\nMonitoring workout... Press Ctrl+C to stop.")

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		// Get current stroke state to detect drive/recovery
		strokeState, err := pm.GetStrokeState()
		if err != nil {
			continue
		}

		// Only get detailed data during recovery (after the stroke)
		if strokeState == csafe.StrokeStateRecovery {
			stats, err := pm.GetStrokeStats()
			if err == nil {
				fmt.Printf("Stroke: Distance=%.2fm, DriveTime=%.2fs, Force=%.1flbs\n",
					float64(stats.StrokeDistance)/100.0,
					float64(stats.DriveTIme)/100.0,
					float64(stats.PeakDriveForce)/10.0)
			}
		}

		// Get general workout status
		snapshot, err := pm.GetWorkoutSnapshot()
		if err == nil {
			fmt.Printf("[%s] %s\n", snapshot.WorkoutState, snapshot)
		}
	}
}

// Example: Get force curve data
func getForceCurve(pm *pm5.PM5) []uint16 {
	var forceCurve []uint16

	// Wait for recovery state (end of stroke)
	for {
		state, err := pm.GetStrokeState()
		if err != nil {
			time.Sleep(10 * time.Millisecond)
			continue
		}
		if state == csafe.StrokeStateRecovery {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	// Read force data in blocks
	for {
		data, err := pm.GetForcePlotData(32) // Read 32 bytes (16 words)
		if err != nil || len(data) == 0 {
			break
		}
		forceCurve = append(forceCurve, data...)

		// If we got fewer than 16 words, we've reached the end
		if len(data) < 16 {
			break
		}
	}

	return forceCurve
}
