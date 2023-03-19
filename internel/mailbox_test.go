package internel

import (
	"reflect"
	"testing"
)

func TestDefaultMailbox(t *testing.T) {
	// Create a new mailbox
	mailbox := NewDefaultMailbox[int]()

	// Add some messages to the mailbox
	for i := 1; i <= 10; i++ {
		mailbox.Source(i)
	}

	// Consume messages from the mailbox
	consumed := make([]int, 0)
	for msg := range mailbox.Consume() {
		consumed = append(consumed, msg)
		if len(consumed) == 10 {
			break
		}
	}

	// Verify that all messages were consumed
	expected := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	if !reflect.DeepEqual(consumed, expected) {
		t.Errorf("unexpected messages: got %v, want %v", consumed, expected)
	}
}
