package debugutil

import (
	"io"
	"log"
)

func Close(c io.Closer) {
	if err := c.Close(); err != nil {
		log.Printf("Closing: %v", err)
	}
}
