package internel

// TickMsg collect msg in an actor, include input, output, exec time, and so on
type TickMsg struct {
	pid      Pid
	input    any
	output   any
	execTime int64
}

// SinkResult every result in the pool contains a unique id, and a list of TickMsg
type SinkResult struct {
	ticks []TickMsg
}

// SinkPool is a storage for SinkResult that are returned by LocalActor
type SinkPool struct {
}

func NewSinkPool() *SinkPool {
	return &SinkPool{}
}

// GetByPid returns a SinkResult by pid
func (s *SinkPool) GetByPid(pid Pid) []SinkResult {
	// TODO: implement me
	return nil
}

// GetByMsg returns a SinkResult by msg
func (s *SinkPool) GetByMsg(msg any) []SinkResult {
	// TODO: implement me
	return nil
}

// Stream returns a stream of SinkResult
func (s *SinkPool) Stream() <-chan SinkResult {
	// TODO: implement me
	return nil
}

// Put puts a SinkResult into the pool
func (s *SinkPool) Put(result SinkResult) {

}
