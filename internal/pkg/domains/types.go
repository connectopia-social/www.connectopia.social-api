package domains

import "sync"

const lifetime_secs = 3 * 3600

type Domain struct {
	Name   string // export required due to json.Marshal
	Expire int64  // expiration timer in seconds
}

type Domains struct {
	domains map[string]Domain
	m       *sync.Mutex
}
