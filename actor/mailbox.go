package actor

// Mailbox is the interface for the actor's mailbox
// MessageOrdering is the ordering of messages in the mailbox
// Throttling is the throttling of messages in the mailbox
// Routing is the routing of messages in the mailbox
// Configuration is the configuration of the mailbox
// Author: fzft

type Mailbox interface {
	MessageOrdering()
}

type DefaultMailbox struct {
}
