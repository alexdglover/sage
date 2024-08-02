package utils

import (
	"fmt"
	"time"
)

func TimeToISO8601DateString(input time.Time) string {
	return fmt.Sprint(input.Format("2006-01-02"))
}
