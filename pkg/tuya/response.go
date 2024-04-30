package tuya

import "github.com/Binozo/GoTuya/internal/commands"

type response struct {
	payload
	returnCode    int
	newSequenceNr int
	commandByte   commands.Type
}
