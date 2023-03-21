package internel

import (
	"fmt"
	"github.com/fzft/my-actor/pkg"
	"time"
)

type TickInMsg struct {
	uid       string
	pid       string
	input     any
	timestamp int64
}

func NewTickInMsg(uid string, pid string, input any) TickInMsg {
	return TickInMsg{
		uid:       uid,
		pid:       pid,
		input:     input,
		timestamp: time.Now().UnixNano(),
	}
}

func (t TickInMsg) String() string {
	return fmt.Sprintf("uid %s, pid %s, input %v, timestamp %d", t.uid, t.pid, t.input, t.timestamp)
}

type TickOutMsg struct {
	uid       string
	pid       string
	output    any
	timestamp int64
}

func (t TickOutMsg) String() string {
	return fmt.Sprintf("uid %s, pid %s, output %v, timestamp %d", t.uid, t.pid, t.output, t.timestamp)
}

func NewTickOutMsg(uid string, pid string, output any) TickOutMsg {
	return TickOutMsg{
		uid:       uid,
		pid:       pid,
		output:    output,
		timestamp: time.Now().UnixNano(),
	}
}

// SinkResult every result in the pool contains a unique id, and a list of TickMsg
type SinkResult struct {
	uid string
	in  []TickInMsg
	out []TickOutMsg
}

func NewSinkResult(uid string) *SinkResult {
	return &SinkResult{
		uid: uid,
	}
}

// String
func (s *SinkResult) String() string {
	return fmt.Sprintf("uid: %s; in: %v; out: %v", s.uid, s.in, s.out)
}

func (s *SinkResult) AddInMsg(msg TickInMsg) {
	s.in = append(s.in, msg)
}

func (s *SinkResult) AddOutMsg(msg TickOutMsg) {
	s.out = append(s.out, msg)
}

// SinkPool is a storage for SinkResult that are returned by LocalActor
type SinkPool struct {
	pool *pkg.KeyValueStore[string, *SinkResult]
}

func NewSinkPool() *SinkPool {
	return &SinkPool{pool: pkg.NewKeyValueStore[string, *SinkResult]()}
}

// GetByPid returns a SinkResult by pid
func (s *SinkPool) GetByPid(pid string) []*SinkResult {

	return nil
}

// GetByMsg returns a SinkResult by msg
func (s *SinkPool) GetByMsg(msg any) []*SinkResult {
	// TODO: implement me
	return nil
}

// Stream returns a stream of SinkResult
func (s *SinkPool) Stream() <-chan SinkResult {
	// TODO: implement me
	return nil
}

// PutInMsg puts a tick into the pool, if sinkResult not exists, create a new one
func (s *SinkPool) PutInMsg(key string, tick TickInMsg) {
	sinkResult, ok := s.pool.Get(key)
	if !ok {
		sinkResult = NewSinkResult(key)
		s.pool.Put(key, sinkResult)
	}
	sinkResult.AddInMsg(tick)
}

// PutOutMsg puts a tick into the pool, if sinkResult not exists, create a new one
func (s *SinkPool) PutOutMsg(key string, tick TickOutMsg) {
	sinkResult, ok := s.pool.Get(key)
	if !ok {
		sinkResult = NewSinkResult(key)
		s.pool.Put(key, sinkResult)
	}
	sinkResult.AddOutMsg(tick)
}

// PopAll returns all SinkResult in the pool
func (s *SinkPool) PopAll() []*SinkResult {
	return s.pool.PopAllVal()
}
