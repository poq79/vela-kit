package audit

func StrOr(str string, def string) string {
	if str == "" {
		return def
	}
	return str
}
