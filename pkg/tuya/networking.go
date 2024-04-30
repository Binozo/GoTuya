package tuya

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Binozo/GoTuya/internal/commands"
	"github.com/Binozo/GoTuya/internal/parser"
	"net"
	"time"
)

const port = 6668
const headerSize = 24

// Connect to the specified tuya device
// Automatically fetches current status
func (d *Device) Connect() error {
	connection, err := net.Dial("tcp", fmt.Sprintf("%s:%d", d.IP, port))
	if err != nil {
		return err
	}
	d.conn = &connection

	_, err = d.sendRefreshCommand()
	return err
}

// Set a dps value and send it to the tuya Device.
// Take a look at the tuyapi project for finding out which dps suits for your device
// (https://github.com/codetheweb/tuyapi)
// Example for a TCL A/C or any other generic Device to turn on:
//
//	Set(map[string]interface{}{
//	    "1": true, // Power on
//	})
func (d *Device) Set(dps map[string]interface{}) error {
	setPayload := payload{
		deviceId: d.DeviceID,
		t:        time.Now(),
		dps:      dps,
		dpId:     nil,
	}

	commandByte := commands.CONTROL
	sequenceNr := d.currentSequenceNr + 1

	if !d.IsConnected() {
		return errors.New("there is no active connection")
	}
	connection := *d.conn

	encoded, err := setPayload.encode(d.Key, commandByte, sequenceNr)
	if err != nil {
		return err
	}

	_, err = connection.Write(encoded)
	if err != nil {
		return err
	}

	setResponse, err := d.readFullPayload()
	if err != nil {
		return err
	}

	if setResponse.commandByte != commandByte {
		return errors.New("the device didn't answer to the command properly")
	}
	return nil
}

// GetCurrentStatus returns the current last given status from the Device without connecting
func (d *Device) GetCurrentStatus() map[string]interface{} {
	return d.currentStatus.dps
}

// FetchStatus connects to the Device and returns the current status
func (d *Device) FetchStatus() (map[string]interface{}, error) {
	curResponse, err := d.sendRefreshCommand()
	if err != nil {
		return nil, err
	}

	return curResponse.dps, nil
}

// sendRefreshCommand refreshes the Device status
func (d *Device) sendRefreshCommand() (response, error) {
	if !d.IsConnected() {
		return response{}, errors.New("there is no active connection")
	}
	connection := *d.conn
	commandByte := commands.DP_REFRESH
	if d.Version == Version_3_4 {
		commandByte = commands.DP_QUERY_NEW
	}

	refreshPayload := payload{
		deviceId: d.DeviceID,
		t:        time.Now(),
		dpId:     []int{4, 5, 6, 18, 19, 20},
	}

	d.currentSequenceNr = 1 // always starts with 1
	encodedPayload, err := refreshPayload.encode(d.Key, commandByte, d.currentSequenceNr)
	if err != nil {
		return response{}, err
	}

	wroteLen, err := connection.Write(encodedPayload)
	if err != nil {
		return response{}, err
	}
	if wroteLen != len(encodedPayload) {
		return response{}, errors.New("tuya device didn't read data")
	}

	d.currentSequenceNr += 1
	commandByte = commands.DP_QUERY
	queryPayload := payload{
		deviceId: d.DeviceID,
		t:        time.Now(),
		dps:      map[string]interface{}{},
	}
	encodedPayload, err = queryPayload.encode(d.Key, commandByte, d.currentSequenceNr)
	if err != nil {
		return response{}, err
	}

	wroteLen, err = connection.Write(encodedPayload)
	if err != nil {
		return response{}, err
	}
	if wroteLen != len(encodedPayload) {
		return response{}, errors.New("tuya device didn't read data")
	}

	curResponse, err := d.readFullPayload()
	d.currentStatus = curResponse
	return curResponse, nil
}

// IsConnected returns if the device is connected
func (d *Device) IsConnected() bool {
	return d.conn != nil
}

// Disconnect any connection to the Device
func (d *Device) Disconnect() {
	if d.IsConnected() {
		connection := *d.conn
		connection.Close()
		d.conn = nil
	}
}

// readFullPayload reads the Device's response to our request
func (d *Device) readFullPayload() (response, error) {
	// We need the first 24 bytes
	// It consists of: prefix (4), sequence (4), command (4), length (4),
	// CRC (4), and suffix (4) for 24 total bytes
	// Information has been taken from: https://github.com/codetheweb/tuyapi/blob/d88fd6c84b228b42f0b6aedd84f0ac3cdb1a5523/lib/message-parser.js#L102

	if !d.IsConnected() {
		return response{}, errors.New("there is no active connection")
	}

	connection := *d.conn
	headerBuffer := make([]byte, headerSize)
	read, err := connection.Read(headerBuffer)
	if err != nil {
		return response{}, err
	}
	if read < 24 {
		return response{}, errors.New(fmt.Sprintf("tuya packet is too short. Length: %d", read))
	}

	// Now validate the header
	if binary.BigEndian.Uint32(headerBuffer[0:4]) != 0x000055AA {
		return response{}, errors.New(fmt.Sprintf("prefix does not match: 0x%02x", headerBuffer[0:4]))
	}

	packetPayloadSize := binary.BigEndian.Uint32(headerBuffer[12:16])
	totalPayloadLength := packetPayloadSize + 16 - headerSize
	// Now we read the remaining data
	packetPayload := make([]byte, totalPayloadLength)
	read, err = connection.Read(packetPayload)
	if err != nil {
		return response{}, err
	}
	if uint32(read) != totalPayloadLength {
		return response{}, errors.New(fmt.Sprintf("mismatch between expected packet size (%d) and actual read size: %d", totalPayloadLength, read))
	}
	totalPayload := append(headerBuffer, packetPayload...)

	if uint32(len(totalPayload)-8) < packetPayloadSize {
		return response{}, errors.New(fmt.Sprintf("packet missing payload (has length: %d)", packetPayloadSize))
	}

	// Check for any additional data
	suffixLocation := bytes.Index(totalPayload, []byte{0x00, 0x00, 0xAA, 0x55})
	if suffixLocation != len(totalPayload)-4 {
		// TODO: not really sure when this happens and what to do with it
		// Take a look: https://github.com/codetheweb/tuyapi/blob/d88fd6c84b228b42f0b6aedd84f0ac3cdb1a5523/lib/message-parser.js#L121
		return response{}, errors.New("this shouldn't happen. please file an issue")
	}

	newSequenceNr := binary.BigEndian.Uint32(totalPayload[4:8])
	commandByte := binary.BigEndian.Uint32(totalPayload[8:12])
	// TODO: implement: https://github.com/codetheweb/tuyapi/blob/d88fd6c84b228b42f0b6aedd84f0ac3cdb1a5523/lib/message-parser.js#L147

	returnCode := binary.BigEndian.Uint32(totalPayload[16:20])
	dataPayload := totalPayload[headerSize-4 : headerSize+packetPayloadSize-16]
	if returnCode&0xFFFFFF00 != 0 {
		// TODO this line below needs to be tested. Probably appears in UDP broadcasts
		// https://github.com/codetheweb/tuyapi/blob/d88fd6c84b228b42f0b6aedd84f0ac3cdb1a5523/lib/message-parser.js#L161
		dataPayload = totalPayload[headerSize-8 : headerSize-8+packetPayloadSize-8]
	}

	// Now we parse the payload
	if (d.Version == Version_3_3 || d.Version == Version_3_2) && bytes.Equal(dataPayload[0:3], []byte("3.3")) {
		// Remove 3.3/3.2 header
		dataPayload = dataPayload[15:]
	} else if d.Version != Version_3_3 && d.Version != Version_3_2 {
		// base64 encoded
		dataPayload, err = base64.StdEncoding.DecodeString(string(dataPayload[19:]))
		if err != nil {
			return response{}, err
		}
	}

	if len(dataPayload) == 0 {
		return response{
			payload: payload{
				deviceId: d.DeviceID,
				t:        time.Now(),
				dps:      nil,
				dpId:     nil,
			},
			returnCode:    int(returnCode),
			commandByte:   commands.Type(commandByte),
			newSequenceNr: int(newSequenceNr),
		}, err
	}

	decrypted, err := parser.DecryptAESWithECB(dataPayload, d.Key)
	if err != nil {
		return response{}, err
	}
	var jsonResponse map[string]interface{}
	// TODO: need to find out where code above goes wrong that len(decrypted)-3 is needed or overall what is happening here
	if bytes.Equal(decrypted[len(decrypted)-3:], []byte{0x03, 0x03, 0x03}) {
		decrypted = decrypted[:len(decrypted)-3]
	} else if bytes.Equal(decrypted[len(decrypted)-2:], []byte{0x02, 0x02}) {
		decrypted = decrypted[:len(decrypted)-2]
	} else if bytes.Equal(decrypted[len(decrypted)-8:], []byte{0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08}) {
		decrypted = decrypted[:len(decrypted)-8]
	}
	if err = json.Unmarshal(decrypted, &jsonResponse); err != nil {
		return response{}, err
	}

	return response{
		payload: payload{
			deviceId: jsonResponse["devId"].(string),
			t:        time.Now(),
			dps:      jsonResponse["dps"].(map[string]interface{}),
			dpId:     nil,
		},
		returnCode:    int(returnCode),
		commandByte:   commands.Type(commandByte),
		newSequenceNr: int(newSequenceNr),
	}, nil
}
