package slack

type BlockType string

const (
	SectionBlock  BlockType = "section"
	HeaderBlock   BlockType = "header"
	ActionsBlock  BlockType = "actions"
	RichTextBlock BlockType = "rich_text"
)

type TextType string

const (
	Mrkdwn               TextType = "mrkdwn"
	PlainText            TextType = "plain_text"
	RichTextSection      TextType = "rich_text_section"
	RichTextPreformatted TextType = "rich_text_preformatted"
)

type ActionType string

const (
	Button     ActionType = "button"
	UserSelect ActionType = "users_select"
)

type Message struct {
	Channel string  `json:"channel,omitempty"`
	Thread  string  `json:"thread_ts,omitempty"`
	Text    string  `json:"text,omitempty"`
	Blocks  []Block `json:"blocks,omitempty"`
}
type Text struct {
	Type  TextType `json:"type,omitempty"`
	Text  string   `json:"text,omitempty"`
	Emoji bool     `json:"emoji,omitempty"`
}
type Field struct {
	Type TextType `json:"type,omitempty"`
	Text string   `json:"text,omitempty"`
}
type Element struct {
	Type     string `json:"type,omitempty"`
	Text     *Text  `json:"text,omitempty"`
	Style    string `json:"style,omitempty"`
	Value    string `json:"value,omitempty"`
	Elements []Text `json:"elements,omitempty"`
	ActionId string `json:"action_id,omitempty"`
}
type Block struct {
	Type      BlockType  `json:"type,omitempty"`
	Text      *Text      `json:"text,omitempty"`
	Fields    []Field    `json:"fields,omitempty"`
	Elements  []Element  `json:"elements,omitempty"`
	Accessory *Accessory `json:"accessory,omitempty"`
}

type Accessory struct {
	Type        ActionType `json:"type,omitempty"`
	Text        *Text      `json:"text,omitempty"`
	Value       string     `json:"value,omitempty"`
	PlaceHolder *Text      `json:"placeholder,omitempty"`
	ActionId    string     `json:"action_id,omitempty"`
}
