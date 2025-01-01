package pkg

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

func SecondsToHoursString(seconds int) string {
	return fmt.Sprintf("%.2f", float32(seconds)/3600.0)
}

// HumanizeDuration humanizes time.Duration output to a meaningful value,
// golang's default “time.Duration“ output is badly formatted and unreadable.
// From: https://gist.github.com/harshavardhana/327e0577c4fed9211f65
func HumanizeDuration(duration time.Duration) string {
	if duration.Seconds() < 60.0 {
		return fmt.Sprintf("%d seconds", int64(duration.Seconds()))
	}
	if duration.Minutes() < 60.0 {
		remainingSeconds := math.Mod(duration.Seconds(), 60)
		return fmt.Sprintf("%d minutes %d seconds", int64(duration.Minutes()), int64(remainingSeconds))
	}
	if duration.Hours() < 24.0 {
		remainingMinutes := math.Mod(duration.Minutes(), 60)
		remainingSeconds := math.Mod(duration.Seconds(), 60)
		return fmt.Sprintf("%d hours %d minutes %d seconds",
			int64(duration.Hours()), int64(remainingMinutes), int64(remainingSeconds))
	}
	remainingHours := math.Mod(duration.Hours(), 24)
	remainingMinutes := math.Mod(duration.Minutes(), 60)
	remainingSeconds := math.Mod(duration.Seconds(), 60)
	return fmt.Sprintf("%d days %d hours %d minutes %d seconds",
		int64(duration.Hours()/24), int64(remainingHours),
		int64(remainingMinutes), int64(remainingSeconds))
}

func HumanizeDurationShort(duration time.Duration) string {
	s := HumanizeDuration(duration)

	// s = strings.Replace(s, " days", "d", 1)
	s = strings.Replace(s, " hours", "h", 1)
	s = strings.Replace(s, " minutes", "m", 1)
	s = strings.Replace(s, " seconds", "s", 1)

	return s
}

// HumanizeInt
// Based on https://github.com/dustin/go-humanize/blob/v1.0.1/comma.go#L15
func HumanizeInt(num int) string {
	parts := []string{"", "", "", "", "", "", ""}
	j := len(parts) - 1

	for num > 999 {
		parts[j] = strconv.FormatInt(int64(num%1000), 10)
		num = num / 1000
		j--
	}

	parts[j] = strconv.Itoa(int(num))
	return strings.Join(parts[j:], ",")
}
