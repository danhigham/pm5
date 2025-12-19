// Package csafe implements the CSAFE protocol for Concept2 Performance Monitors.
package csafe

// Frame flags for CSAFE protocol
const (
	ExtendedFrameStartFlag byte = 0xF0
	StandardFrameStartFlag byte = 0xF1
	StopFrameFlag          byte = 0xF2
	ByteStuffingFlag       byte = 0xF3
)

// Frame addressing
const (
	AddressPCHost       byte = 0x00
	AddressDefaultSlave byte = 0xFD
	AddressBroadcast    byte = 0xFF
)

// Manufacturer information
const (
	ManufacturerID      byte = 22
	ClassID             byte = 2
	ModelPM3            byte = 3
	ModelPM4            byte = 4
	ModelPM5            byte = 5
	MaxFrameLength      int  = 120
	MinInterframeGapMs  int  = 50
)

// OperationalState represents the PM operational state
type OperationalState byte

const (
	OperationalStateReset           OperationalState = 0
	OperationalStateReady           OperationalState = 1
	OperationalStateWorkout         OperationalState = 2
	OperationalStateWarmup          OperationalState = 3
	OperationalStateRace            OperationalState = 4
	OperationalStatePowerOff        OperationalState = 5
	OperationalStatePause           OperationalState = 6
	OperationalStateInvokeBootloader OperationalState = 7
	OperationalStatePowerOffShip    OperationalState = 8
	OperationalStateIdleCharge      OperationalState = 9
	OperationalStateIdle            OperationalState = 10
	OperationalStateMfgTest         OperationalState = 11
	OperationalStateFWUpdate        OperationalState = 12
	OperationalStateDragFactor      OperationalState = 13
	OperationalStateDFCalibration   OperationalState = 100
)

func (s OperationalState) String() string {
	names := map[OperationalState]string{
		OperationalStateReset:           "Reset",
		OperationalStateReady:           "Ready",
		OperationalStateWorkout:         "Workout",
		OperationalStateWarmup:          "Warmup",
		OperationalStateRace:            "Race",
		OperationalStatePowerOff:        "PowerOff",
		OperationalStatePause:           "Pause",
		OperationalStateInvokeBootloader: "InvokeBootloader",
		OperationalStatePowerOffShip:    "PowerOffShip",
		OperationalStateIdleCharge:      "IdleCharge",
		OperationalStateIdle:            "Idle",
		OperationalStateMfgTest:         "MfgTest",
		OperationalStateFWUpdate:        "FWUpdate",
		OperationalStateDragFactor:      "DragFactor",
		OperationalStateDFCalibration:   "DFCalibration",
	}
	if name, ok := names[s]; ok {
		return name
	}
	return "Unknown"
}

// ErgMachineType represents the type of ergometer machine
type ErgMachineType byte

const (
	ErgMachineTypeStaticD         ErgMachineType = 0
	ErgMachineTypeStaticC         ErgMachineType = 1
	ErgMachineTypeStaticA         ErgMachineType = 2
	ErgMachineTypeStaticB         ErgMachineType = 3
	ErgMachineTypeStaticE         ErgMachineType = 5
	ErgMachineTypeStaticSimulator ErgMachineType = 7
	ErgMachineTypeStaticDynamic   ErgMachineType = 8
	ErgMachineTypeSlidesA         ErgMachineType = 16
	ErgMachineTypeSlidesB         ErgMachineType = 17
	ErgMachineTypeSlidesC         ErgMachineType = 18
	ErgMachineTypeSlidesD         ErgMachineType = 19
	ErgMachineTypeSlidesE         ErgMachineType = 20
	ErgMachineTypeLinkedDynamic   ErgMachineType = 32
	ErgMachineTypeStaticDyno      ErgMachineType = 64
	ErgMachineTypeStaticSki       ErgMachineType = 128
	ErgMachineTypeSkiSimulator    ErgMachineType = 143
	ErgMachineTypeBike            ErgMachineType = 192
	ErgMachineTypeBikeArms        ErgMachineType = 193
	ErgMachineTypeBikeNoArms      ErgMachineType = 194
	ErgMachineTypeBikeSimulator   ErgMachineType = 207
	ErgMachineTypeMultiErgRow     ErgMachineType = 224
	ErgMachineTypeMultiErgSki     ErgMachineType = 225
	ErgMachineTypeMultiErgBike    ErgMachineType = 226
)

func (t ErgMachineType) String() string {
	names := map[ErgMachineType]string{
		ErgMachineTypeStaticD:         "Rower Model D",
		ErgMachineTypeStaticC:         "Rower Model C",
		ErgMachineTypeStaticA:         "Rower Model A",
		ErgMachineTypeStaticB:         "Rower Model B",
		ErgMachineTypeStaticE:         "Rower Model E",
		ErgMachineTypeStaticSimulator: "Rower Simulator",
		ErgMachineTypeStaticDynamic:   "Dynamic Rower",
		ErgMachineTypeSlidesA:         "Slides Model A",
		ErgMachineTypeSlidesB:         "Slides Model B",
		ErgMachineTypeSlidesC:         "Slides Model C",
		ErgMachineTypeSlidesD:         "Slides Model D",
		ErgMachineTypeSlidesE:         "Slides Model E",
		ErgMachineTypeLinkedDynamic:   "Linked Dynamic",
		ErgMachineTypeStaticDyno:      "Dynamometer",
		ErgMachineTypeStaticSki:       "SkiErg",
		ErgMachineTypeSkiSimulator:    "SkiErg Simulator",
		ErgMachineTypeBike:            "BikeErg",
		ErgMachineTypeBikeArms:        "BikeErg with Arms",
		ErgMachineTypeBikeNoArms:      "BikeErg No Arms",
		ErgMachineTypeBikeSimulator:   "BikeErg Simulator",
		ErgMachineTypeMultiErgRow:     "MultiErg Row",
		ErgMachineTypeMultiErgSki:     "MultiErg Ski",
		ErgMachineTypeMultiErgBike:    "MultiErg Bike",
	}
	if name, ok := names[t]; ok {
		return name
	}
	return "Unknown"
}

// WorkoutType represents the type of workout
type WorkoutType byte

const (
	WorkoutTypeJustRowNoSplits                WorkoutType = 0
	WorkoutTypeJustRowSplits                  WorkoutType = 1
	WorkoutTypeFixedDistNoSplits              WorkoutType = 2
	WorkoutTypeFixedDistSplits                WorkoutType = 3
	WorkoutTypeFixedTimeNoSplits              WorkoutType = 4
	WorkoutTypeFixedTimeSplits                WorkoutType = 5
	WorkoutTypeFixedTimeInterval              WorkoutType = 6
	WorkoutTypeFixedDistInterval              WorkoutType = 7
	WorkoutTypeVariableInterval               WorkoutType = 8
	WorkoutTypeVariableUndefinedRestInterval  WorkoutType = 9
	WorkoutTypeFixedCalorieSplits             WorkoutType = 10
	WorkoutTypeFixedWattMinuteSplits          WorkoutType = 11
	WorkoutTypeFixedCalsInterval              WorkoutType = 12
)

func (t WorkoutType) String() string {
	names := map[WorkoutType]string{
		WorkoutTypeJustRowNoSplits:                "Just Row (No Splits)",
		WorkoutTypeJustRowSplits:                  "Just Row (Splits)",
		WorkoutTypeFixedDistNoSplits:              "Fixed Distance (No Splits)",
		WorkoutTypeFixedDistSplits:                "Fixed Distance (Splits)",
		WorkoutTypeFixedTimeNoSplits:              "Fixed Time (No Splits)",
		WorkoutTypeFixedTimeSplits:                "Fixed Time (Splits)",
		WorkoutTypeFixedTimeInterval:              "Fixed Time Interval",
		WorkoutTypeFixedDistInterval:              "Fixed Distance Interval",
		WorkoutTypeVariableInterval:               "Variable Interval",
		WorkoutTypeVariableUndefinedRestInterval:  "Variable Interval (Undefined Rest)",
		WorkoutTypeFixedCalorieSplits:             "Fixed Calorie (Splits)",
		WorkoutTypeFixedWattMinuteSplits:          "Fixed Watt-Minute (Splits)",
		WorkoutTypeFixedCalsInterval:              "Fixed Calorie Interval",
	}
	if name, ok := names[t]; ok {
		return name
	}
	return "Unknown"
}

// IntervalType represents the type of interval
type IntervalType byte

const (
	IntervalTypeTime                     IntervalType = 0
	IntervalTypeDist                     IntervalType = 1
	IntervalTypeRest                     IntervalType = 2
	IntervalTypeTimeRestUndefined        IntervalType = 3
	IntervalTypeDistanceRestUndefined    IntervalType = 4
	IntervalTypeRestUndefined            IntervalType = 5
	IntervalTypeCalorie                  IntervalType = 6
	IntervalTypeCalorieRestUndefined     IntervalType = 7
	IntervalTypeWattMinute               IntervalType = 8
	IntervalTypeWattMinuteRestUndefined  IntervalType = 9
	IntervalTypeNone                     IntervalType = 255
)

func (t IntervalType) String() string {
	names := map[IntervalType]string{
		IntervalTypeTime:                    "Time",
		IntervalTypeDist:                    "Distance",
		IntervalTypeRest:                    "Rest",
		IntervalTypeTimeRestUndefined:       "Time (Undefined Rest)",
		IntervalTypeDistanceRestUndefined:   "Distance (Undefined Rest)",
		IntervalTypeRestUndefined:           "Undefined Rest",
		IntervalTypeCalorie:                 "Calorie",
		IntervalTypeCalorieRestUndefined:    "Calorie (Undefined Rest)",
		IntervalTypeWattMinute:              "Watt-Minute",
		IntervalTypeWattMinuteRestUndefined: "Watt-Minute (Undefined Rest)",
		IntervalTypeNone:                    "None",
	}
	if name, ok := names[t]; ok {
		return name
	}
	return "Unknown"
}

// WorkoutState represents the current workout state
type WorkoutState byte

const (
	WorkoutStateWaitToBegin                    WorkoutState = 0
	WorkoutStateWorkoutRow                     WorkoutState = 1
	WorkoutStateCountdownPause                 WorkoutState = 2
	WorkoutStateIntervalRest                   WorkoutState = 3
	WorkoutStateIntervalWorkTime               WorkoutState = 4
	WorkoutStateIntervalWorkDistance           WorkoutState = 5
	WorkoutStateIntervalRestEndToWorkTime      WorkoutState = 6
	WorkoutStateIntervalRestEndToWorkDistance  WorkoutState = 7
	WorkoutStateIntervalWorkTimeToRest         WorkoutState = 8
	WorkoutStateIntervalWorkDistanceToRest     WorkoutState = 9
	WorkoutStateWorkoutEnd                     WorkoutState = 10
	WorkoutStateTerminate                      WorkoutState = 11
	WorkoutStateWorkoutLogged                  WorkoutState = 12
	WorkoutStateRearm                          WorkoutState = 13
)

func (s WorkoutState) String() string {
	names := map[WorkoutState]string{
		WorkoutStateWaitToBegin:                   "Wait To Begin",
		WorkoutStateWorkoutRow:                    "Workout Row",
		WorkoutStateCountdownPause:                "Countdown Pause",
		WorkoutStateIntervalRest:                  "Interval Rest",
		WorkoutStateIntervalWorkTime:              "Interval Work Time",
		WorkoutStateIntervalWorkDistance:          "Interval Work Distance",
		WorkoutStateIntervalRestEndToWorkTime:     "Interval Rest End To Work Time",
		WorkoutStateIntervalRestEndToWorkDistance: "Interval Rest End To Work Distance",
		WorkoutStateIntervalWorkTimeToRest:        "Interval Work Time To Rest",
		WorkoutStateIntervalWorkDistanceToRest:    "Interval Work Distance To Rest",
		WorkoutStateWorkoutEnd:                    "Workout End",
		WorkoutStateTerminate:                     "Terminate",
		WorkoutStateWorkoutLogged:                 "Workout Logged",
		WorkoutStateRearm:                         "Rearm",
	}
	if name, ok := names[s]; ok {
		return name
	}
	return "Unknown"
}

// RowingState represents whether rowing is active
type RowingState byte

const (
	RowingStateInactive RowingState = 0
	RowingStateActive   RowingState = 1
)

func (s RowingState) String() string {
	if s == RowingStateActive {
		return "Active"
	}
	return "Inactive"
}

// StrokeState represents the current stroke state
type StrokeState byte

const (
	StrokeStateWaitingForWheelToReachMinSpeed StrokeState = 0
	StrokeStateWaitingForWheelToAccelerate    StrokeState = 1
	StrokeStateDriving                        StrokeState = 2
	StrokeStateDwellingAfterDrive             StrokeState = 3
	StrokeStateRecovery                       StrokeState = 4
)

func (s StrokeState) String() string {
	names := map[StrokeState]string{
		StrokeStateWaitingForWheelToReachMinSpeed: "Waiting for Wheel",
		StrokeStateWaitingForWheelToAccelerate:    "Waiting to Accelerate",
		StrokeStateDriving:                        "Driving",
		StrokeStateDwellingAfterDrive:             "Dwelling",
		StrokeStateRecovery:                       "Recovery",
	}
	if name, ok := names[s]; ok {
		return name
	}
	return "Unknown"
}

// DurationType represents workout duration identifier
type DurationType byte

const (
	DurationTypeTime      DurationType = 0x00
	DurationTypeCalories  DurationType = 0x40
	DurationTypeDistance  DurationType = 0x80
	DurationTypeWattMin   DurationType = 0xC0
)

// ScreenType represents the screen type for PM commands
type ScreenType byte

const (
	ScreenTypeNone    ScreenType = 0
	ScreenTypeWorkout ScreenType = 1
	ScreenTypeRace    ScreenType = 2
	ScreenTypeCSAFE   ScreenType = 3
	ScreenTypeDiag    ScreenType = 4
	ScreenTypeMfg     ScreenType = 5
)

// ScreenValueWorkout represents screen values for workout type
type ScreenValueWorkout byte

const (
	ScreenValueWorkoutNone                         ScreenValueWorkout = 0
	ScreenValueWorkoutPrepareToRowWorkout          ScreenValueWorkout = 1
	ScreenValueWorkoutTerminateWorkout             ScreenValueWorkout = 2
	ScreenValueWorkoutRearmWorkout                 ScreenValueWorkout = 3
	ScreenValueWorkoutRefreshLogCard               ScreenValueWorkout = 4
	ScreenValueWorkoutPrepareToRaceStart           ScreenValueWorkout = 5
	ScreenValueWorkoutGoToMainScreen               ScreenValueWorkout = 6
	ScreenValueWorkoutLogCardBusyWarning           ScreenValueWorkout = 7
	ScreenValueWorkoutLogCardSelectUser            ScreenValueWorkout = 8
	ScreenValueWorkoutResetRaceParams              ScreenValueWorkout = 9
	ScreenValueWorkoutCableTestSlave               ScreenValueWorkout = 10
	ScreenValueWorkoutFishGame                     ScreenValueWorkout = 11
	ScreenValueWorkoutDisplayParticipantInfo       ScreenValueWorkout = 12
	ScreenValueWorkoutDisplayParticipantInfoConfirm ScreenValueWorkout = 13
	ScreenValueWorkoutChangeDisplayTypeTarget      ScreenValueWorkout = 20
	ScreenValueWorkoutChangeDisplayTypeStandard    ScreenValueWorkout = 21
	ScreenValueWorkoutChangeDisplayTypeForceCurve  ScreenValueWorkout = 22
	ScreenValueWorkoutChangeDisplayTypePaceBoat    ScreenValueWorkout = 23
)

// DisplayUnitsType represents display units
type DisplayUnitsType byte

const (
	DisplayUnitsTimeMeters       DisplayUnitsType = 0
	DisplayUnitsPace             DisplayUnitsType = 1
	DisplayUnitsWatts            DisplayUnitsType = 2
	DisplayUnitsCaloricBurnRate  DisplayUnitsType = 3
	DisplayUnitsCalories         DisplayUnitsType = 4
)

// DisplayFormatType represents display format
type DisplayFormatType byte

const (
	DisplayTypeStandard      DisplayFormatType = 0
	DisplayTypeForceCurve    DisplayFormatType = 1
	DisplayTypePaceBoat      DisplayFormatType = 2
	DisplayTypePerStroke     DisplayFormatType = 3
	DisplayTypeSimple        DisplayFormatType = 4
	DisplayTypeTarget        DisplayFormatType = 5
)

// Status byte bit masks for CSAFE response
const (
	StatusFrameToggleMask     byte = 0x80
	StatusPrevFrameStatusMask byte = 0x30
	StatusStateMask           byte = 0x0F
)

// Previous frame status values
const (
	PrevFrameStatusOK       byte = 0x00
	PrevFrameStatusReject   byte = 0x10
	PrevFrameStatusBad      byte = 0x20
	PrevFrameStatusNotReady byte = 0x30
)

// State machine states (from status byte)
const (
	StateMachineError   byte = 0x00
	StateMachineReady   byte = 0x01
	StateMachineIdle    byte = 0x02
	StateMachineHaveID  byte = 0x03
	StateMachineInUse   byte = 0x05
	StateMachinePause   byte = 0x06
	StateMachineFinish  byte = 0x07
	StateMachineManual  byte = 0x08
	StateMachineOffLine byte = 0x09
)

// WorkoutNumber represents predefined workout numbers
type WorkoutNumber byte

const (
	WorkoutNumberProgrammed WorkoutNumber = 0
	WorkoutNumberDefault1   WorkoutNumber = 1
	WorkoutNumberDefault2   WorkoutNumber = 2
	WorkoutNumberDefault3   WorkoutNumber = 3
	WorkoutNumberDefault4   WorkoutNumber = 4
	WorkoutNumberDefault5   WorkoutNumber = 5
	WorkoutNumberCustom1    WorkoutNumber = 6
	WorkoutNumberCustom2    WorkoutNumber = 7
	WorkoutNumberCustom3    WorkoutNumber = 8
	WorkoutNumberCustom4    WorkoutNumber = 9
	WorkoutNumberCustom5    WorkoutNumber = 10
)

// Units specifiers for CSAFE commands
const (
	UnitsMeter    byte = 0x24 // Meters
	UnitsKm       byte = 0x21 // Kilometers
	UnitsWatt     byte = 0x58 // Watts
	UnitsSeconds  byte = 0x00 // Seconds
)
