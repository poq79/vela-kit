package vela

type Region interface {
	Name() string
	Search(string) (*IPv4Info, error)
}

type RegionByEnv interface {
	WithRegion(interface{})
	Region(interface{}) (*IPv4Info, error)
}
