package csafe

import (
	"bytes"
	"errors"
	"fmt"
)

var (
	ErrFrameTooShort    = errors.New("frame too short")
	ErrInvalidStartFlag = errors.New("invalid start flag")
	ErrInvalidStopFlag  = errors.New("invalid stop flag")
	ErrInvalidChecksum  = errors.New("invalid checksum")
	ErrFrameTooLong     = errors.New("frame exceeds maximum length")
	ErrInvalidStuffByte = errors.New("invalid byte stuffing value")
)

// Frame represents a CSAFE frame
type Frame struct {
	Extended    bool
	Destination byte
	Source      byte
	Contents    []byte
}

// Response represents a parsed CSAFE response
type Response struct {
	Status          byte
	FrameToggle     bool
	PrevFrameStatus byte
	StateMachine    byte
	CommandData     []CommandResponse
}

// CommandResponse represents an individual command response within a frame
type CommandResponse struct {
	Command     byte
	ByteCount   byte
	Data        []byte
	PMResponses []CommandResponse // Populated when Command is a PM wrapper (0x76, 0x77, 0x7E, 0x7F)
}

// FirstPMResponse returns the first PM response, or nil if none
func (c *CommandResponse) FirstPMResponse() *CommandResponse {
	if len(c.PMResponses) > 0 {
		return &c.PMResponses[0]
	}
	return nil
}

// EncodeFrame encodes a CSAFE frame with byte stuffing
func EncodeFrame(f *Frame) ([]byte, error) {
	var buf bytes.Buffer

	// Start flag
	if f.Extended {
		buf.WriteByte(ExtendedFrameStartFlag)
		// Write destination and source addresses (with byte stuffing)
		stuffAndWrite(&buf, f.Destination)
		stuffAndWrite(&buf, f.Source)
	} else {
		buf.WriteByte(StandardFrameStartFlag)
	}

	// Calculate checksum on unstuffed contents
	checksum := byte(0)
	for _, b := range f.Contents {
		checksum ^= b
	}

	// Write contents with byte stuffing
	for _, b := range f.Contents {
		stuffAndWrite(&buf, b)
	}

	// Write checksum with byte stuffing
	stuffAndWrite(&buf, checksum)

	// Stop flag
	buf.WriteByte(StopFrameFlag)

	result := buf.Bytes()
	if len(result) > MaxFrameLength {
		return nil, ErrFrameTooLong
	}

	return result, nil
}

// DecodeFrame decodes a CSAFE frame with byte unstuffing
func DecodeFrame(data []byte) (*Frame, error) {
	if len(data) < 4 { // Minimum: start + content + checksum + stop
		return nil, ErrFrameTooShort
	}

	// Check start flag
	extended := false
	switch data[0] {
	case ExtendedFrameStartFlag:
		extended = true
	case StandardFrameStartFlag:
		extended = false
	default:
		return nil, ErrInvalidStartFlag
	}

	// Check stop flag
	if data[len(data)-1] != StopFrameFlag {
		return nil, ErrInvalidStopFlag
	}

	// Unstuff the frame contents (excluding start and stop flags)
	unstuffed, err := unstuffBytes(data[1 : len(data)-1])
	if err != nil {
		return nil, err
	}

	if len(unstuffed) < 1 {
		return nil, ErrFrameTooShort
	}

	frame := &Frame{Extended: extended}

	// Parse addresses for extended frame
	offset := 0
	if extended {
		if len(unstuffed) < 3 { // dest + source + checksum minimum
			return nil, ErrFrameTooShort
		}
		frame.Destination = unstuffed[0]
		frame.Source = unstuffed[1]
		offset = 2
	}

	// Last byte is checksum
	if len(unstuffed) <= offset {
		return nil, ErrFrameTooShort
	}
	checksumIdx := len(unstuffed) - 1
	receivedChecksum := unstuffed[checksumIdx]

	// Extract contents (between addresses and checksum)
	frame.Contents = unstuffed[offset:checksumIdx]

	// Verify checksum
	calculatedChecksum := byte(0)
	for _, b := range frame.Contents {
		calculatedChecksum ^= b
	}

	if calculatedChecksum != receivedChecksum {
		return nil, ErrInvalidChecksum
	}

	return frame, nil
}

// ParseResponse parses the contents of a CSAFE response frame
func ParseResponse(contents []byte) (*Response, error) {
	if len(contents) < 1 {
		return nil, ErrFrameTooShort
	}

	resp := &Response{
		Status:          contents[0],
		FrameToggle:     (contents[0] & StatusFrameToggleMask) != 0,
		PrevFrameStatus: contents[0] & StatusPrevFrameStatusMask,
		StateMachine:    contents[0] & StatusStateMask,
		CommandData:     make([]CommandResponse, 0),
	}

	// Parse command responses
	offset := 1
	for offset < len(contents) {
		cmd := contents[offset]
		offset++

		if offset >= len(contents) {
			resp.CommandData = append(resp.CommandData, CommandResponse{
				Command:   cmd,
				ByteCount: 0,
				Data:      nil,
			})
			break
		}

		byteCount := contents[offset]
		offset++

		var data []byte
		if byteCount > 0 {
			if offset+int(byteCount) > len(contents) {
				data = contents[offset:]
				offset = len(contents)
			} else {
				data = contents[offset : offset+int(byteCount)]
				offset += int(byteCount)
			}
		}

		cmdResp := CommandResponse{
			Command:   cmd,
			ByteCount: byteCount,
			Data:      data,
		}

		// If this is a PM wrapper command, parse the nested PM responses
		if isPMWrapper(cmd) && len(data) > 0 {
			cmdResp.PMResponses = parsePMWrapperData(data)
		}

		resp.CommandData = append(resp.CommandData, cmdResp)
	}

	return resp, nil
}

// isPMWrapper returns true if the command is a PM proprietary wrapper
func isPMWrapper(cmd byte) bool {
	return cmd == 0x76 || cmd == 0x77 || cmd == 0x7E || cmd == 0x7F
}

// parsePMWrapperData parses PM proprietary command responses from wrapper data.
// PM responses follow the same [Cmd][ByteCount][Data] format as standard CSAFE.
func parsePMWrapperData(data []byte) []CommandResponse {
	var responses []CommandResponse
	offset := 0

	for offset < len(data) {
		cmd := data[offset]
		offset++

		if offset >= len(data) {
			responses = append(responses, CommandResponse{
				Command:   cmd,
				ByteCount: 0,
				Data:      nil,
			})
			break
		}

		byteCount := data[offset]
		offset++

		var cmdData []byte
		if byteCount > 0 {
			if offset+int(byteCount) > len(data) {
				cmdData = data[offset:]
				offset = len(data)
			} else {
				cmdData = data[offset : offset+int(byteCount)]
				offset += int(byteCount)
			}
		}

		responses = append(responses, CommandResponse{
			Command:   cmd,
			ByteCount: byteCount,
			Data:      cmdData,
		})
	}

	return responses
}

// BuildCommand builds a single CSAFE command
func BuildCommand(cmd byte, data ...byte) []byte {
	if cmd&0x80 != 0 && len(data) == 0 {
		// Short command (no data)
		return []byte{cmd}
	}
	// Long command
	result := make([]byte, 2+len(data))
	result[0] = cmd
	result[1] = byte(len(data))
	copy(result[2:], data)
	return result
}

// BuildPMCommand builds a PM-specific command wrapped in appropriate wrapper
func BuildPMCommand(wrapper byte, commands ...[]byte) []byte {
	// Calculate total size of inner commands
	totalSize := 0
	for _, cmd := range commands {
		totalSize += len(cmd)
	}

	result := make([]byte, 2+totalSize)
	result[0] = wrapper
	result[1] = byte(totalSize)

	offset := 2
	for _, cmd := range commands {
		copy(result[offset:], cmd)
		offset += len(cmd)
	}

	return result
}

// stuffAndWrite writes a byte with byte stuffing if necessary
func stuffAndWrite(buf *bytes.Buffer, b byte) {
	switch b {
	case ExtendedFrameStartFlag:
		buf.WriteByte(ByteStuffingFlag)
		buf.WriteByte(0x00)
	case StandardFrameStartFlag:
		buf.WriteByte(ByteStuffingFlag)
		buf.WriteByte(0x01)
	case StopFrameFlag:
		buf.WriteByte(ByteStuffingFlag)
		buf.WriteByte(0x02)
	case ByteStuffingFlag:
		buf.WriteByte(ByteStuffingFlag)
		buf.WriteByte(0x03)
	default:
		buf.WriteByte(b)
	}
}

// unstuffBytes removes byte stuffing from a byte slice
func unstuffBytes(data []byte) ([]byte, error) {
	result := make([]byte, 0, len(data))

	for i := 0; i < len(data); i++ {
		if data[i] == ByteStuffingFlag {
			if i+1 >= len(data) {
				return nil, ErrInvalidStuffByte
			}
			i++
			switch data[i] {
			case 0x00:
				result = append(result, ExtendedFrameStartFlag)
			case 0x01:
				result = append(result, StandardFrameStartFlag)
			case 0x02:
				result = append(result, StopFrameFlag)
			case 0x03:
				result = append(result, ByteStuffingFlag)
			default:
				return nil, ErrInvalidStuffByte
			}
		} else {
			result = append(result, data[i])
		}
	}

	return result, nil
}

// PrevFrameStatusString returns a human-readable string for the previous frame status
func PrevFrameStatusString(status byte) string {
	switch status {
	case PrevFrameStatusOK:
		return "OK"
	case PrevFrameStatusReject:
		return "Reject"
	case PrevFrameStatusBad:
		return "Bad"
	case PrevFrameStatusNotReady:
		return "Not Ready"
	default:
		return fmt.Sprintf("Unknown (0x%02X)", status)
	}
}

// StateMachineString returns a human-readable string for the state machine state
func StateMachineString(state byte) string {
	switch state {
	case StateMachineError:
		return "Error"
	case StateMachineReady:
		return "Ready"
	case StateMachineIdle:
		return "Idle"
	case StateMachineHaveID:
		return "Have ID"
	case StateMachineInUse:
		return "In Use"
	case StateMachinePause:
		return "Pause"
	case StateMachineFinish:
		return "Finish"
	case StateMachineManual:
		return "Manual"
	case StateMachineOffLine:
		return "Off Line"
	default:
		return fmt.Sprintf("Unknown (0x%02X)", state)
	}
}
