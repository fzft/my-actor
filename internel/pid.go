package internel

import "github.com/google/uuid"

type Pid struct {
	uuid      string
	actorName string
}

func NewPid(actorName string) Pid {
	id := uuid.New()
	return Pid{actorName: actorName, uuid: id.String()}
}
