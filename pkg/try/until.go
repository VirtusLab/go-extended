package try

import (
	"fmt"
	"time"

	time2 "github.com/VirtusLab/go-extended/pkg/time"
)

// ErrTimout is used when the set timeout has been reached
type ErrTimout struct {
	text string
}

func (e *ErrTimout) Error() string {
	return e.text
}

// Until keeps trying until timeout or there is a result or an error
func Until(something func() (bool, error), tick, timeout time.Duration) (bool, error) {
	counter := 0
	tickChan := time2.Every(tick)
	timeoutChan := time.After(timeout)
	for {
		select {
		case <-tickChan:
			ok, err := something()
			if err != nil {
				return false, err
			} else if ok {
				return true, nil
			}
			counter = counter + 1
		case <-timeoutChan:
			return false, &ErrTimout{
				text: fmt.Sprintf("timed out after: %s, tries: %d", timeout, counter),
			}
		}
	}
}
