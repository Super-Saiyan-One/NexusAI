package constant

const (
	RelayModeUnkown = iota
	RelayModeChatCompletion
	RelayModeCompletion
	RelayModeEmbedding
	RelayModeModeration
	RelayModeEdit
)

const (
	ChannelTypeUnknown = iota
	ChannelTypeOpenAI
	ChannelTypeClaude
	ChannelTypeGemini
	ChannelTypeGrok
	ChannelTypeDeepSeek
	ChannelTypeQwen
	ChannelTypeBaidu
	ChannelTypeMidjourney
)
