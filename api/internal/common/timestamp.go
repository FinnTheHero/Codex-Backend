package common

import "time"

func TimeStamp(date string) string {
	if date != "" {
		t, err := time.Parse("2006-01-02 15:04:05", date)
		if err != nil {
			return time.Now().Format("2006-01-02 15:04:05")
		}
		return t.Format("2006-01-02 15:04:05")
	}

	return time.Now().Format("2006-01-02 15:04:05")
}
