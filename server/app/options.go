package app

type LogOptions struct {
	TailLines int64
	Follow    bool
	PodName   string
	Previous  bool
	Container string
}

type PodListOptions struct {
	PodName string
}

type Pod struct {
	Name     string
	State    string
	Age      int64
	Restarts int32
	Ready    bool
}
