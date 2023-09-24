package inputs

import (
	"time"
)

type Inputer interface {
	ProcessAll()
	Process(file string) error
	StoreReport(data []byte) error
	Watch(interval time.Duration)
}
