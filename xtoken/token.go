package xtoken

type Token struct {
	Channel    string
	Username   string
	NodeID     int64
	Expiration int64
}

func NewToken(channel, username string, nodeid, expiration int64) *Token {
	return &Token{
		Channel:    channel,
		Username:   username,
		NodeID:     nodeid,
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
	nodeid, ok := data["nodeid"].(float64)
	if !ok {
		nodeid = 0
	}
	expiration, ok := data["expiration"].(float64)
	if !ok {
		expiration = 0
	}

	return NewToken(channel, username, int64(nodeid), int64(expiration))
}

func (t *Token) Marshal() map[string]any {
	return map[string]any{
		"channel":    t.Channel,
		"username":   t.Username,
		"nodeid":     t.NodeID,
		"expiration": t.Expiration,
	}
}

func (t *Token) SameChannel(channel string) bool {
	return t.Channel == channel
}

func (t *Token) SameUsername(username string) bool {
	return t.Username == username
}

func (t *Token) SameNodeID(nodeid int64) bool {
	return t.NodeID == nodeid
}

func (t *Token) NotExpired(n int64) bool {
	return t.Expiration >= n
}
