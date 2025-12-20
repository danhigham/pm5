package csafe

// Public CSAFE Short Commands (responses only - no data sent)
const (
	// Get status/state commands
	CmdGetStatus  byte = 0x80
	CmdReset      byte = 0x81
	CmdGoIdle     byte = 0x82
	CmdGoHaveID   byte = 0x83
	CmdGoInUse    byte = 0x85
	CmdGoFinished byte = 0x86
	CmdGoReady    byte = 0x87
	CmdBadID      byte = 0x88
	CmdGetVersion byte = 0x91
	CmdGetID      byte = 0x92
	CmdGetUnits   byte = 0x93
	CmdGetSerial  byte = 0x94

	// Get data commands
	CmdGetOdometer   byte = 0x9B
	CmdGetErrorCode  byte = 0x9C
	CmdGetTWork      byte = 0xA0
	CmdGetHorizontal byte = 0xA1
	CmdGetCalories   byte = 0xA3
	CmdGetProgram    byte = 0xA4
	CmdGetPace       byte = 0xA6
	CmdGetCadence    byte = 0xA7
	CmdGetUserInfo   byte = 0xAB
	CmdGetHRCur      byte = 0xB0
	CmdGetPower      byte = 0xB4
)

// Public CSAFE Long Commands (commands with data)
const (
	CmdAutoUpload    byte = 0x01
	CmdIDDigits      byte = 0x10
	CmdSetTime       byte = 0x11
	CmdSetDate       byte = 0x12
	CmdSetTimeout    byte = 0x13
	CmdSetUserCfg1   byte = 0x1A // Wrapper for PM-specific commands
	CmdSetTWork      byte = 0x20
	CmdSetHorizontal byte = 0x21
	CmdSetCalories   byte = 0x23
	CmdSetProgram    byte = 0x24
	CmdSetPower      byte = 0x34
	CmdGetCaps       byte = 0x70
)

// PM Proprietary CSAFE Command Wrappers
const (
	CmdSetPMCfg  byte = 0x76
	CmdSetPMData byte = 0x77
	CmdGetPMCfg  byte = 0x7E
	CmdGetPMData byte = 0x7F
)

// C2 Proprietary Short Get Configuration Commands
const (
	PMCmdGetFWVersion             byte = 0x80
	PMCmdGetHWVersion             byte = 0x81
	PMCmdGetHWAddress             byte = 0x82
	PMCmdGetTickTimebase          byte = 0x83
	PMCmdGetHRM                   byte = 0x84
	PMCmdGetDateTime              byte = 0x85
	PMCmdGetScreenStateStatus     byte = 0x86
	PMCmdGetRaceLaneRequest       byte = 0x87
	PMCmdGetRaceEntryRequest      byte = 0x88
	PMCmdGetWorkoutType           byte = 0x89
	PMCmdGetDisplayType           byte = 0x8A
	PMCmdGetDisplayUnits          byte = 0x8B
	PMCmdGetLanguageType          byte = 0x8C
	PMCmdGetWorkoutState          byte = 0x8D
	PMCmdGetIntervalType          byte = 0x8E
	PMCmdGetOperationalState      byte = 0x8F
	PMCmdGetLogCardState          byte = 0x90
	PMCmdGetLogCardStatus         byte = 0x91
	PMCmdGetPowerUpState          byte = 0x92
	PMCmdGetRowingState           byte = 0x93
	PMCmdGetScreenContentVersion  byte = 0x94
	PMCmdGetCommunicationState    byte = 0x95
	PMCmdGetRaceParticipantCount  byte = 0x96
	PMCmdGetBatteryLevelPercent   byte = 0x97
	PMCmdGetRaceModeStatus        byte = 0x98
	PMCmdGetInternalLogParams     byte = 0x99
	PMCmdGetProductConfiguration  byte = 0x9A
	PMCmdGetCPUTickRate           byte = 0x9D
	PMCmdGetLogCardUserCensus     byte = 0x9E
	PMCmdGetWorkoutIntervalCount  byte = 0x9F
	PMCmdGetWorkoutDuration       byte = 0xE8
	PMCmdGetWorkOther             byte = 0xE9
	PMCmdGetExtendedHRM           byte = 0xEA
	PMCmdGetDFCalibrationVerified byte = 0xEB
	PMCmdGetFlywheelSpeed         byte = 0xEC
	PMCmdGetErgMachineType        byte = 0xED
	PMCmdGetRaceBeginEndTickCount byte = 0xEE
	PMCmdGetPM5FWUpdateStatus     byte = 0xEF
)

// C2 Proprietary Short Get Data Commands
const (
	PMCmdGetWorkTime                byte = 0xA0
	PMCmdGetProjectedWorkTime       byte = 0xA1
	PMCmdGetTotalRestTime           byte = 0xA2
	PMCmdGetWorkDistance            byte = 0xA3
	PMCmdGetTotalWorkDistance       byte = 0xA4
	PMCmdGetProjectedWorkDistance   byte = 0xA5
	PMCmdGetRestDistance            byte = 0xA6
	PMCmdGetTotalRestDistance       byte = 0xA7
	PMCmdGetStroke500mPace          byte = 0xA8
	PMCmdGetStrokePower             byte = 0xA9
	PMCmdGetStrokeCaloricBurnRate   byte = 0xAA
	PMCmdGetSplitAvg500mPace        byte = 0xAB
	PMCmdGetSplitAvgPower           byte = 0xAC
	PMCmdGetSplitAvgCaloricBurnRate byte = 0xAD
	PMCmdGetSplitAvgCalories        byte = 0xAE
	PMCmdGetTotalAvg500mPace        byte = 0xAF
	PMCmdGetTotalAvgPower           byte = 0xB0
	PMCmdGetTotalAvgCaloricBurnRate byte = 0xB1
	PMCmdGetTotalAvgCalories        byte = 0xB2
	PMCmdGetStrokeRate              byte = 0xB3
	PMCmdGetSplitAvgStrokeRate      byte = 0xB4
	PMCmdGetTotalAvgStrokeRate      byte = 0xB5
	PMCmdGetAvgHeartRate            byte = 0xB6
	PMCmdGetEndingAvgHeartRate      byte = 0xB7
	PMCmdGetRestAvgHeartRate        byte = 0xB8
	PMCmdGetSplitTime               byte = 0xB9
	PMCmdGetLastSplitTime           byte = 0xBA
	PMCmdGetSplitDistance           byte = 0xBB
	PMCmdGetLastSplitDistance       byte = 0xBC
	PMCmdGetLastRestDistance        byte = 0xBD
	PMCmdGetTargetPaceTime          byte = 0xBE
	PMCmdGetStrokeState             byte = 0xBF
	PMCmdGetStrokeRateState         byte = 0xC0
	PMCmdGetDragFactor              byte = 0xC1
	PMCmdGetEncoderPeriod           byte = 0xC2
	PMCmdGetHeartRateState          byte = 0xC3
	PMCmdGetSyncData                byte = 0xC4
	PMCmdGetSyncDataAll             byte = 0xC5
	PMCmdGetRaceData                byte = 0xC6
	PMCmdGetTickTime                byte = 0xC7
	PMCmdGetErrorType               byte = 0xC8
	PMCmdGetErrorValue              byte = 0xC9
	PMCmdGetStatusType              byte = 0xCA
	PMCmdGetStatusValue             byte = 0xCB
	PMCmdGetEPMStatus               byte = 0xCC
	PMCmdGetDisplayUpdateTime       byte = 0xCD
	PMCmdGetSyncFractionalTime      byte = 0xCE
	PMCmdGetRestTime                byte = 0xCF
)

// C2 Proprietary Long Get Data Commands
const (
	PMCmdGetMemory             byte = 0x68
	PMCmdGetLogCardMemory      byte = 0x69
	PMCmdGetInternalLogMemory  byte = 0x6A
	PMCmdGetForcePlotData      byte = 0x6B
	PMCmdGetHeartBeatData      byte = 0x6C
	PMCmdGetUIEvents           byte = 0x6D
	PMCmdGetStrokeStats        byte = 0x6E
	PMCmdGetCurrentWorkoutHash byte = 0x72
	PMCmdGetGameScore          byte = 0x78
)

// C2 Proprietary Long Get Configuration Commands
const (
	PMCmdGetErgNumber            byte = 0x50
	PMCmdGetErgNumberRequest     byte = 0x51
	PMCmdGetUserIDString         byte = 0x52
	PMCmdGetLocalRaceParticipant byte = 0x53
	PMCmdGetUserID               byte = 0x54
	PMCmdGetUserProfile          byte = 0x55
	PMCmdGetHRBeltInfo           byte = 0x56
	PMCmdGetExtendedHRBeltInfo   byte = 0x57
	PMCmdGetCurrentLogStructure  byte = 0x58
)

// C2 Proprietary Short Set Configuration Commands
const (
	PMCmdSetResetAll       byte = 0xE0
	PMCmdSetResetErgNumber byte = 0xE1
)

// C2 Proprietary Long Set Configuration Commands
const (
	PMCmdSetWorkoutType            byte = 0x01
	PMCmdSetWorkoutDuration        byte = 0x03
	PMCmdSetRestDuration           byte = 0x04
	PMCmdSetSplitDuration          byte = 0x05
	PMCmdSetTargetPaceTime         byte = 0x06
	PMCmdSetRaceType               byte = 0x09
	PMCmdSetRaceLaneSetup          byte = 0x0B
	PMCmdSetRaceLaneVerify         byte = 0x0C
	PMCmdSetRaceStartParams        byte = 0x0D
	PMCmdSetErgNumber              byte = 0x10
	PMCmdSetScreenState            byte = 0x13
	PMCmdConfigureWorkout          byte = 0x14
	PMCmdSetTargetAvgWatts         byte = 0x15
	PMCmdSetTargetCalsPerHr        byte = 0x16
	PMCmdSetIntervalType           byte = 0x17
	PMCmdSetWorkoutIntervalCount   byte = 0x18
	PMCmdSetDisplayUpdateRate      byte = 0x19
	PMCmdSetAuthenPassword         byte = 0x1A
	PMCmdSetTickTime               byte = 0x1B
	PMCmdSetTickTimeOffset         byte = 0x1C
	PMCmdSetRaceDataSampleTicks    byte = 0x1D
	PMCmdSetRaceOperationType      byte = 0x1E
	PMCmdSetRaceStatusDisplayTicks byte = 0x1F
	PMCmdSetRaceStatusWarningTicks byte = 0x20
	PMCmdSetRaceIdleModeParams     byte = 0x21
	PMCmdSetDateTime               byte = 0x22
	PMCmdSetLanguageType           byte = 0x23
	PMCmdSetScreenErrorMode        byte = 0x27
	PMCmdSetUserID                 byte = 0x29
	PMCmdSetUserProfile            byte = 0x2A
	PMCmdSetHRM                    byte = 0x2B
	PMCmdSetHRBeltInfo             byte = 0x2D
	PMCmdSetRaceParticipant        byte = 0x32
	PMCmdSetRaceStatus             byte = 0x33
	PMCmdSetLogCardMemory          byte = 0x34
	PMCmdSetDisplayString          byte = 0x35
	PMCmdSetDisplayBitmap          byte = 0x36
	PMCmdSetLocalRaceParticipant   byte = 0x37
	PMCmdSetGameParams             byte = 0x38
	PMCmdSetExtendedHRBeltInfo     byte = 0x39
	PMCmdSetExtendedHRM            byte = 0x3A
	PMCmdSetLEDBacklight           byte = 0x3B
	PMCmdSetWirelessChannelConfig  byte = 0x3D
	PMCmdSetRaceControlParams      byte = 0x3E
)

// C2 Proprietary Short Set Data Commands
const (
	PMCmdSetSyncDistance         byte = 0xD0
	PMCmdSetSyncStrokePace       byte = 0xD1
	PMCmdSetSyncAvgHeartRate     byte = 0xD2
	PMCmdSetSyncTime             byte = 0xD3
	PMCmdSetSyncRaceTickTime     byte = 0xD7
	PMCmdSetSyncDataAll          byte = 0xD8
	PMCmdSetSyncRowingActiveTime byte = 0xD9
)
