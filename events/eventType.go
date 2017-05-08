package events

type EventType int16

const (
	EventSaveBlock             EventType = 0
	EventReplyTx               EventType = 1
	EventBlockPersistCompleted EventType = 2
	EventBlockSaveCompleted    EventType = 3
	EventNewInventory          EventType = 4
)
