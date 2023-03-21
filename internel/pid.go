package internel

import (
	"fmt"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const defaultBufferSize = 1024

type Pid struct {
	logger    *zap.SugaredLogger
	uuid      string
	actorName string

	actor   Actor
	context *Context

	TickInMsgCh  chan TickInMsg
	TickOutMsgCh chan TickOutMsg
}

func NewPid(logger *zap.SugaredLogger, actor Actor) *Pid {
	id := uuid.New().String()
	pid := &Pid{
		actorName:    actor.String(),
		uuid:         id,
		context:      NewContext(logger, id),
		actor:        actor,
		logger:       logger,
		TickInMsgCh:  make(chan TickInMsg, defaultBufferSize),
		TickOutMsgCh: make(chan TickOutMsg, defaultBufferSize),
	}
	return pid
}

func (p *Pid) String() string {
	if p.actorName == "" {
		return fmt.Sprintf("pid:%s", p.uuid)
	}
	return fmt.Sprintf("pid:%s", p.actorName)
}

// run is the actor's main loop
func (p *Pid) run() {

	go p.context.buffered()

	if d, ok := p.actor.(PreStartHookActor); ok {
		d.PreStart()
	}
	for {
		msg, ok := p.context.inbox.Dequeue()
		if ok {
			input, ok := msg.(Message)
			if !ok {
				continue
			}
			p.TickInMsgCh <- NewTickInMsg(input.uid, p.String(), input.data)

			p.logger.Debugw("run", "pid", p.String(), "msg", input)
			if d, ok := p.actor.(PreHandleMsgHookActor); ok {
				d.PreHandleMsg(p.context, input)
			}
			output, err := p.actor.Receive(p.context, input.data)
			if err == nil {
				p.logger.Debugw("run", "pid", p.String(), "output", output)
				if d, ok := p.actor.(PostHandleMsgHookActor); ok {
					d.PostHandleMsg(p.context, output)
				}

				p.TickOutMsgCh <- NewTickOutMsg(input.uid, p.String(), output)
				p.context.broadcast(WrapMsg(input.uid, output))
			} else {
				p.logger.Errorw("run", "pid", p.String(), "err", err)
				p.TickOutMsgCh <- NewTickOutMsg(input.uid, p.String(), output)
				if d, ok := p.actor.(ErrHandlerActor); ok {
					d.ErrHandler(p.context, err)
				}
			}
		}
	}

	if d, ok := p.actor.(PostStopHookActor); ok {
		d.PostStop()
	}
}
