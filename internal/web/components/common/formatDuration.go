package common

import "time"
import "fmt"

func FormatDuration(d time.Duration) string {
	rounded := d.Round(time.Minute)
	hours := rounded / time.Hour
	minutes := (rounded % time.Hour) / time.Minute
	return fmt.Sprintf("%dh%02dm", hours, minutes)
}
