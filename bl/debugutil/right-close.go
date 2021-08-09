package debugutil

import "log"

type Closer interface {
	Close() error
}

func Close(c Closer) {
	if err := c.Close(); err != nil {
		log.Printf("Closing: %v", err)
	}
}
