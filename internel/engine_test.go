package internel

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type Dummy struct {
}

func (d *Dummy) Receive(ctx *Context, msg any) (any, error) {
	switch v := msg.(type) {
	case string:
		return fmt.Sprintf("dummy:%s", v), nil
	case int:
		v = v + 1
		return v, nil
	default:
		return nil, fmt.Errorf("dummy not support msg type")
	}
}

func (d *Dummy) String() string {
	return "dummy"
}

func newDummy() *Dummy {
	return &Dummy{}
}

// TestEngine one node DAG, with dummy actor and recv msg and handle it
func TestEngine_OneNode(t *testing.T) {
	engine := NewEngine()
	dummy := newDummy()
	pid, err := engine.Spawn(dummy)
	assert.Nil(t, err)

	assert.Equal(t, "pid:dummy", pid.String())
	assert.Equal(t, 1, len(engine.nodeMaps))
	assert.Equal(t, 1, len(engine.pidMaps))

	assert.Nil(t, engine.Ready())
	assert.Nil(t, engine.Send("hello"))

	time.Sleep(1 * time.Second)

	results := engine.sinkPool.PopAll()
	t.Log(results)

}
