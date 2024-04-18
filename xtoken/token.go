package xtoken

type Token struct {
	Channel    string
	Username   string
	Nodename   string
	Expiration int64
}

func NewToken(channel, username, nodename string, expiration int64) *Token {
	return &Token{
		Channel:    channel,
		Username:   username,
		Nodename:   nodename,
		Expiration: expiration,
	}
}

func UnmarshalToken(data map[string]any) *Token {
	channel, ok := data["channel"].(string)
	if !ok {
		channel = ""
	}
	username, ok := data["username"].(string)
	if !ok {
		username = ""
	}
	nodename, ok := data["nodename"].(string)
	if !ok {
		nodename = ""
	}
	expiration, ok := data["expiration"].(float64)
	if !ok {
		expiration = 0
	}
	return NewToken(channel, username, nodename, int64(expiration))
}

func (t *Token) Marshal() map[string]any {
	return map[string]any{
		"channel":    t.Channel,
		"username":   t.Username,
		"nodename":   t.Nodename,
		"expiration": t.Expiration,
	}
}

func (t *Token) NotExpired(n int64) bool {
	return t.Expiration >= n
}
