package pm5

import (
	"github.com/danhigham/pm5/csafe"
)

// ============================================================================
// PM5 Proprietary Set Configuration Commands
// ============================================================================

// SetWorkoutType sets the workout type
func (p *PM5) SetWorkoutType(workoutType csafe.WorkoutType) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	pmCmd := csafe.BuildCommand(csafe.PMCmdSetWorkoutType, byte(workoutType))
	_, err := p.sendPMCommand(csafe.CmdSetPMCfg, pmCmd)
	return err
}

// SetWorkoutDuration sets the workout duration
// durationType specifies Time (0x00), Calories (0x40), Distance (0x80), or WattMin (0xC0)
// duration is in appropriate units: 0.01s for time, meters for distance, cals, or watt-min
func (p *PM5) SetWorkoutDuration(durationType csafe.DurationType, duration uint32) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	pmCmd := csafe.BuildCommand(csafe.PMCmdSetWorkoutDuration,
		byte(durationType),
		byte((duration>>24)&0xFF),
		byte((duration>>16)&0xFF),
		byte((duration>>8)&0xFF),
		byte(duration&0xFF))
	_, err := p.sendPMCommand(csafe.CmdSetPMCfg, pmCmd)
	return err
}

// SetRestDuration sets the rest duration in seconds
func (p *PM5) SetRestDuration(seconds uint16) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	pmCmd := csafe.BuildCommand(csafe.PMCmdSetRestDuration,
		byte((seconds>>8)&0xFF),
		byte(seconds&0xFF))
	_, err := p.sendPMCommand(csafe.CmdSetPMCfg, pmCmd)
	return err
}

// SetSplitDuration sets the split duration
// durationType specifies Time (0x00), Calories (0x40), Distance (0x80), or WattMin (0xC0)
func (p *PM5) SetSplitDuration(durationType csafe.DurationType, duration uint32) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	pmCmd := csafe.BuildCommand(csafe.PMCmdSetSplitDuration,
		byte(durationType),
		byte((duration>>24)&0xFF),
		byte((duration>>16)&0xFF),
		byte((duration>>8)&0xFF),
		byte(duration&0xFF))
	_, err := p.sendPMCommand(csafe.CmdSetPMCfg, pmCmd)
	return err
}

// SetTargetPaceTime sets the target pace time in hundredths of seconds per 500m
func (p *PM5) SetTargetPaceTime(paceTime uint32) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	pmCmd := csafe.BuildCommand(csafe.PMCmdSetTargetPaceTime,
		byte((paceTime>>24)&0xFF),
		byte((paceTime>>16)&0xFF),
		byte((paceTime>>8)&0xFF),
		byte(paceTime&0xFF))
	_, err := p.sendPMCommand(csafe.CmdSetPMCfg, pmCmd)
	return err
}

// SetIntervalType sets the interval type for interval workouts
func (p *PM5) SetIntervalType(intervalType csafe.IntervalType) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	pmCmd := csafe.BuildCommand(csafe.PMCmdSetIntervalType, byte(intervalType))
	_, err := p.sendPMCommand(csafe.CmdSetPMCfg, pmCmd)
	return err
}

// SetWorkoutIntervalCount sets the current interval number (1-indexed)
func (p *PM5) SetWorkoutIntervalCount(count byte) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	pmCmd := csafe.BuildCommand(csafe.PMCmdSetWorkoutIntervalCount, count)
	_, err := p.sendPMCommand(csafe.CmdSetPMCfg, pmCmd)
	return err
}

// SetTargetAvgWatts sets the target average watts
func (p *PM5) SetTargetAvgWatts(watts uint16) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	pmCmd := csafe.BuildCommand(csafe.PMCmdSetTargetAvgWatts,
		byte((watts>>8)&0xFF),
		byte(watts&0xFF))
	_, err := p.sendPMCommand(csafe.CmdSetPMCfg, pmCmd)
	return err
}

// SetTargetCalsPerHour sets the target calories per hour
func (p *PM5) SetTargetCalsPerHour(calsPerHr uint16) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	pmCmd := csafe.BuildCommand(csafe.PMCmdSetTargetCalsPerHr,
		byte((calsPerHr>>8)&0xFF),
		byte(calsPerHr&0xFF))
	_, err := p.sendPMCommand(csafe.CmdSetPMCfg, pmCmd)
	return err
}

// ConfigureWorkout enables or disables workout programming mode
func (p *PM5) ConfigureWorkout(enable bool) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	mode := byte(0)
	if enable {
		mode = 1
	}

	pmCmd := csafe.BuildCommand(csafe.PMCmdConfigureWorkout, mode)
	_, err := p.sendPMCommand(csafe.CmdSetPMCfg, pmCmd)
	return err
}

// SetScreenState sets the screen type and value
func (p *PM5) SetScreenState(screenType csafe.ScreenType, screenValue byte) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	pmCmd := csafe.BuildCommand(csafe.PMCmdSetScreenState,
		byte(screenType), screenValue)
	_, err := p.sendPMCommand(csafe.CmdSetPMCfg, pmCmd)
	return err
}

// SetScreenErrorMode enables or disables screen error display mode
func (p *PM5) SetScreenErrorMode(enable bool) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	mode := byte(0)
	if enable {
		mode = 1
	}

	pmCmd := csafe.BuildCommand(csafe.PMCmdSetScreenErrorMode, mode)
	_, err := p.sendPMCommand(csafe.CmdSetPMCfg, pmCmd)
	return err
}

// SetDisplayUpdateRate sets how often display updates are sent
// 0=1sec, 1=500ms (default), 2=250ms, 3=100ms
func (p *PM5) SetDisplayUpdateRate(rate byte) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	pmCmd := csafe.BuildCommand(csafe.PMCmdSetDisplayUpdateRate, rate)
	_, err := p.sendPMCommand(csafe.CmdSetPMCfg, pmCmd)
	return err
}

// DateTime represents date and time for the PM
type DateTime struct {
	Hours    byte // 1-12
	Minutes  byte // 0-59
	Meridiem byte // 0=AM, 1=PM
	Month    byte // 1-12
	Day      byte // 1-31
	Year     uint16
}

// SetDateTime sets the PM5 date and time
func (p *PM5) SetDateTime(dt *DateTime) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	pmCmd := csafe.BuildCommand(csafe.PMCmdSetDateTime,
		dt.Hours,
		dt.Minutes,
		dt.Meridiem,
		dt.Month,
		dt.Day,
		byte((dt.Year>>8)&0xFF),
		byte(dt.Year&0xFF))
	_, err := p.sendPMCommand(csafe.CmdSetPMCfg, pmCmd)
	return err
}

// ============================================================================
// Workout Setup Helpers
// ============================================================================

// StartJustRowWorkout starts a simple "Just Row" workout with optional splits
func (p *PM5) StartJustRowWorkout(withSplits bool) error {
	workoutType := csafe.WorkoutTypeJustRowNoSplits
	if withSplits {
		workoutType = csafe.WorkoutTypeJustRowSplits
	}

	// Build combined command
	p.mu.Lock()
	defer p.mu.Unlock()

	pmCmds := [][]byte{
		csafe.BuildCommand(csafe.PMCmdSetWorkoutType, byte(workoutType)),
		csafe.BuildCommand(csafe.PMCmdSetScreenState,
			byte(csafe.ScreenTypeWorkout),
			byte(csafe.ScreenValueWorkoutPrepareToRowWorkout)),
	}

	_, err := p.sendPMCommand(csafe.CmdSetPMCfg, pmCmds...)
	return err
}

// StartFixedDistanceWorkout starts a fixed distance workout
// distance is in meters, splitDistance is in meters (0 for no splits)
func (p *PM5) StartFixedDistanceWorkout(distance uint32, splitDistance uint32) error {
	workoutType := csafe.WorkoutTypeFixedDistNoSplits
	if splitDistance > 0 {
		workoutType = csafe.WorkoutTypeFixedDistSplits
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	pmCmds := [][]byte{
		csafe.BuildCommand(csafe.PMCmdSetWorkoutType, byte(workoutType)),
		csafe.BuildCommand(csafe.PMCmdSetWorkoutDuration,
			byte(csafe.DurationTypeDistance),
			byte((distance>>24)&0xFF),
			byte((distance>>16)&0xFF),
			byte((distance>>8)&0xFF),
			byte(distance&0xFF)),
	}

	if splitDistance > 0 {
		pmCmds = append(pmCmds, csafe.BuildCommand(csafe.PMCmdSetSplitDuration,
			byte(csafe.DurationTypeDistance),
			byte((splitDistance>>24)&0xFF),
			byte((splitDistance>>16)&0xFF),
			byte((splitDistance>>8)&0xFF),
			byte(splitDistance&0xFF)))
	}

	pmCmds = append(pmCmds,
		csafe.BuildCommand(csafe.PMCmdConfigureWorkout, 0x01), // Enable
		csafe.BuildCommand(csafe.PMCmdSetScreenState,
			byte(csafe.ScreenTypeWorkout),
			byte(csafe.ScreenValueWorkoutPrepareToRowWorkout)))

	_, err := p.sendPMCommand(csafe.CmdSetPMCfg, pmCmds...)
	return err
}

// StartFixedTimeWorkout starts a fixed time workout
// duration is in hundredths of seconds, splitDuration is in hundredths of seconds (0 for no splits)
func (p *PM5) StartFixedTimeWorkout(duration uint32, splitDuration uint32) error {
	workoutType := csafe.WorkoutTypeFixedTimeNoSplits
	if splitDuration > 0 {
		workoutType = csafe.WorkoutTypeFixedTimeSplits
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	pmCmds := [][]byte{
		csafe.BuildCommand(csafe.PMCmdSetWorkoutType, byte(workoutType)),
		csafe.BuildCommand(csafe.PMCmdSetWorkoutDuration,
			byte(csafe.DurationTypeTime),
			byte((duration>>24)&0xFF),
			byte((duration>>16)&0xFF),
			byte((duration>>8)&0xFF),
			byte(duration&0xFF)),
	}

	if splitDuration > 0 {
		pmCmds = append(pmCmds, csafe.BuildCommand(csafe.PMCmdSetSplitDuration,
			byte(csafe.DurationTypeTime),
			byte((splitDuration>>24)&0xFF),
			byte((splitDuration>>16)&0xFF),
			byte((splitDuration>>8)&0xFF),
			byte(splitDuration&0xFF)))
	}

	pmCmds = append(pmCmds,
		csafe.BuildCommand(csafe.PMCmdConfigureWorkout, 0x01), // Enable
		csafe.BuildCommand(csafe.PMCmdSetScreenState,
			byte(csafe.ScreenTypeWorkout),
			byte(csafe.ScreenValueWorkoutPrepareToRowWorkout)))

	_, err := p.sendPMCommand(csafe.CmdSetPMCfg, pmCmds...)
	return err
}

// StartFixedCalorieWorkout starts a fixed calorie workout
// calories is the goal, splitCalories is per split (0 for no splits)
func (p *PM5) StartFixedCalorieWorkout(calories uint32, splitCalories uint32) error {
	workoutType := csafe.WorkoutTypeJustRowNoSplits // Will be updated
	if splitCalories > 0 {
		workoutType = csafe.WorkoutTypeFixedCalorieSplits
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	pmCmds := [][]byte{
		csafe.BuildCommand(csafe.PMCmdSetWorkoutType, byte(workoutType)),
		csafe.BuildCommand(csafe.PMCmdSetWorkoutDuration,
			byte(csafe.DurationTypeCalories),
			byte((calories>>24)&0xFF),
			byte((calories>>16)&0xFF),
			byte((calories>>8)&0xFF),
			byte(calories&0xFF)),
	}

	if splitCalories > 0 {
		pmCmds = append(pmCmds, csafe.BuildCommand(csafe.PMCmdSetSplitDuration,
			byte(csafe.DurationTypeCalories),
			byte((splitCalories>>24)&0xFF),
			byte((splitCalories>>16)&0xFF),
			byte((splitCalories>>8)&0xFF),
			byte(splitCalories&0xFF)))
	}

	pmCmds = append(pmCmds,
		csafe.BuildCommand(csafe.PMCmdConfigureWorkout, 0x01),
		csafe.BuildCommand(csafe.PMCmdSetScreenState,
			byte(csafe.ScreenTypeWorkout),
			byte(csafe.ScreenValueWorkoutPrepareToRowWorkout)))

	_, err := p.sendPMCommand(csafe.CmdSetPMCfg, pmCmds...)
	return err
}

// StartFixedDistanceIntervalWorkout starts a fixed distance interval workout
// distance is in meters, restSeconds is rest duration in seconds
func (p *PM5) StartFixedDistanceIntervalWorkout(distance uint32, restSeconds uint16) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	pmCmds := [][]byte{
		csafe.BuildCommand(csafe.PMCmdSetWorkoutType, byte(csafe.WorkoutTypeFixedDistInterval)),
		csafe.BuildCommand(csafe.PMCmdSetWorkoutDuration,
			byte(csafe.DurationTypeDistance),
			byte((distance>>24)&0xFF),
			byte((distance>>16)&0xFF),
			byte((distance>>8)&0xFF),
			byte(distance&0xFF)),
		csafe.BuildCommand(csafe.PMCmdSetRestDuration,
			byte((restSeconds>>8)&0xFF),
			byte(restSeconds&0xFF)),
		csafe.BuildCommand(csafe.PMCmdConfigureWorkout, 0x01),
		csafe.BuildCommand(csafe.PMCmdSetScreenState,
			byte(csafe.ScreenTypeWorkout),
			byte(csafe.ScreenValueWorkoutPrepareToRowWorkout)),
	}

	_, err := p.sendPMCommand(csafe.CmdSetPMCfg, pmCmds...)
	return err
}

// StartFixedTimeIntervalWorkout starts a fixed time interval workout
// duration is in hundredths of seconds, restSeconds is rest duration in seconds
func (p *PM5) StartFixedTimeIntervalWorkout(duration uint32, restSeconds uint16) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	pmCmds := [][]byte{
		csafe.BuildCommand(csafe.PMCmdSetWorkoutType, byte(csafe.WorkoutTypeFixedTimeInterval)),
		csafe.BuildCommand(csafe.PMCmdSetWorkoutDuration,
			byte(csafe.DurationTypeTime),
			byte((duration>>24)&0xFF),
			byte((duration>>16)&0xFF),
			byte((duration>>8)&0xFF),
			byte(duration&0xFF)),
		csafe.BuildCommand(csafe.PMCmdSetRestDuration,
			byte((restSeconds>>8)&0xFF),
			byte(restSeconds&0xFF)),
		csafe.BuildCommand(csafe.PMCmdConfigureWorkout, 0x01),
		csafe.BuildCommand(csafe.PMCmdSetScreenState,
			byte(csafe.ScreenTypeWorkout),
			byte(csafe.ScreenValueWorkoutPrepareToRowWorkout)),
	}

	_, err := p.sendPMCommand(csafe.CmdSetPMCfg, pmCmds...)
	return err
}

// TerminateWorkout terminates the current workout
func (p *PM5) TerminateWorkout() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	pmCmd := csafe.BuildCommand(csafe.PMCmdSetScreenState,
		byte(csafe.ScreenTypeWorkout),
		byte(csafe.ScreenValueWorkoutTerminateWorkout))
	_, err := p.sendPMCommand(csafe.CmdSetPMCfg, pmCmd)
	return err
}

// GoToMainScreen navigates to the main screen
func (p *PM5) GoToMainScreen() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	pmCmd := csafe.BuildCommand(csafe.PMCmdSetScreenState,
		byte(csafe.ScreenTypeWorkout),
		byte(csafe.ScreenValueWorkoutGoToMainScreen))
	_, err := p.sendPMCommand(csafe.CmdSetPMCfg, pmCmd)
	return err
}
