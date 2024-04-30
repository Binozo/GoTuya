package tuya

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/Binozo/GoTuya/internal/commands"
	"github.com/Binozo/GoTuya/internal/parser"
	"time"
)

type payload struct {
	deviceId string
	t        time.Time
	dps      map[string]interface{}
	dpId     []int
}

// exportJson to match the Device's requirements
func (p *payload) exportJson() ([]byte, error) {
	rawJson := map[string]interface{}{
		"gwId":  p.deviceId,
		"devId": p.deviceId,
		"t":     fmt.Sprintf("%d", p.t.Unix()),
		"uid":   p.deviceId,
	}
	if p.dpId != nil {
		rawJson["dpId"] = p.dpId
	}
	if p.dps != nil {
		rawJson["dps"] = p.dps
	}
	return json.Marshal(rawJson)
}

// encode the payload to make it ready to send to the Device
func (p *payload) encode(key []byte, command commands.Type, currentSeqNr int) ([]byte, error) {
	jsonBuffer, err := p.exportJson()
	if err != nil {
		return nil, err
	}

	encrypted, err := parser.EncryptAESWithECB(jsonBuffer, key)
	if err != nil {
		return nil, err
	}

	if command != commands.DP_QUERY && command != commands.DP_REFRESH {
		encryptionPayloadBuffer := make([]byte, len(encrypted)+15)
		version33header := []byte("3.3")
		copy(encryptionPayloadBuffer[0:len(version33header)], version33header)
		copy(encryptionPayloadBuffer[15:], encrypted)
		encrypted = encryptionPayloadBuffer
	}

	// Copy everything together
	finalPayloadBuffer := make([]byte, len(encrypted)+24)

	// Write some header-like info
	binary.BigEndian.PutUint32(finalPayloadBuffer[0:], 0x000055AA)
	binary.BigEndian.PutUint32(finalPayloadBuffer[8:], uint32(command))
	binary.BigEndian.PutUint32(finalPayloadBuffer[12:], uint32(len(encrypted)+8))

	if currentSeqNr > 0 {
		binary.BigEndian.PutUint32(finalPayloadBuffer[4:], uint32(currentSeqNr))
	}

	// Copy in the encrypted payload
	copy(finalPayloadBuffer[16:], encrypted)

	// Calculate crc
	preCrcBuffer := make([]byte, len(encrypted)+16)
	copy(preCrcBuffer, finalPayloadBuffer)
	calculatedCrc := parser.CalculateCrc(preCrcBuffer) & 0xFFFFFFFF
	binary.BigEndian.PutUint32(finalPayloadBuffer[len(encrypted)+16:], calculatedCrc)
	binary.BigEndian.PutUint32(finalPayloadBuffer[len(encrypted)+20:], uint32(0x0000AA55))

	return finalPayloadBuffer, nil
}
