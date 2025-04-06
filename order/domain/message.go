package domain

type Message struct {
	Key   string
	Value MessageValue
	Topic string

	// Add fields for consumer if needed:
	Partition int
	Offset    int64
}
type MessageValue struct {
	Meta        *MetaData `json:"meta"`
	MessageCode string    `json:"message_code"`
	Payload     any       `json:"payload"`
}

type MetaData struct {
	MessageID         string `json:"message_id"`
	OriginalMessageID string `json:"original_message_id,omitempty"`
	ServiceID         string `json:"service_id"`
	Timestamp         int64  `json:"timestamp"`

	// Add retry-specific fields here
	RetryCount    int    `json:"retry_count,omitempty"`
	OriginalTopic string `json:"original_topic,omitempty"`
	OriginalKey   string `json:"original_key,omitempty"`
	ErrorMessage  string `json:"error_message,omitempty"`
}
