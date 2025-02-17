package gowa

type MessageAck int

const (
	MessageAckError   MessageAck = -1
	MessageAckPending MessageAck = 0
	MessageAckServer  MessageAck = 1
	MessageAckDevice  MessageAck = 2
	MessageAckRead    MessageAck = 3
	MessageAckPlayed  MessageAck = 4
)

type Message interface {
	Content() string
	Option() map[string]any
}

type MessageText struct {
	MessageContent string
}

func (m *MessageText) Content() string {
	return m.MessageContent
}

func (m *MessageText) Option() map[string]any {
	return nil
}

type MessageMedia struct {
	MimeType string
	Filename string
	Caption  string
	Data     string //base64 data
}

func (m *MessageMedia) Content() string {
	return ""
}

func (m *MessageMedia) Option() map[string]any {
	return map[string]any{
		"caption": m.Caption,
		"attachment": map[string]any{
			"mimetype": m.MimeType,
			"filename": m.Filename,
			"data":     m.Data,
		},
	}
}
