package node

func ID() string {
	return _G.id
}

func Prefix() string {
	if _G.prefix == "" {
		return "share"
	}

	return _G.prefix
}
