package node

func ID() string {
	return ssc.id
}

func Prefix() string {
	if ssc.prefix == "" {
		return "share"
	}

	return ssc.prefix
}
