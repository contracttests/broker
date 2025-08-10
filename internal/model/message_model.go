package model

type MessageResource struct {
	Hash       string `json:"hash,omitzero"`
	Name       string `json:"name,omitzero"`
	Topic      string `json:"topic,omitzero"`
	SchemaName string `json:"schemaName,omitzero"`
	Type       string `json:"type,omitzero"`
}

func NewMessageProducer(
	name string,
	topic string,
	schemaName string,
) MessageResource {
	hash := HashFromStrings(name, topic, schemaName)

	return MessageResource{
		Hash:       hash,
		Name:       name,
		Topic:      topic,
		SchemaName: schemaName,
		Type:       "producer",
	}
}

func NewMessageConsumer(
	name string,
	topic string,
	schemaName string,
) MessageResource {
	hash := HashFromStrings(name, topic, schemaName)

	return MessageResource{
		Hash:       hash,
		Name:       name,
		Topic:      topic,
		SchemaName: schemaName,
		Type:       "consumer",
	}
}
