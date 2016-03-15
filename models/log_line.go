package models

import "time"

type LogLine struct {
	At   time.Time `json:"at"`
	URI  string    `json:"uri"`
	Code int       `json:"code"`
}

func (l LogLine) Status() string {
	code := l.Code
	switch {
	case code >= 200 && code < 300:
		return "ok"
	case code >= 300 && code < 400:
		return "redirect"
	case code >= 400 && code < 500:
		return "error"
	case code >= 500:
		return "fault"
	default:
		return "unknown"
	}
}

type LogBlock struct {
	Header string
	Items  []LogLine
}

func MakeLogBlocks(lines []LogLine) []LogBlock {
	var blocks []LogBlock
	var items []LogLine
	var lastAt string
	first := true

	for _, line := range lines {
		at := line.At.Format("02 Jan; 15:04 PM")

		if at != lastAt && !first {
			blocks = append(blocks, LogBlock{
				Header: lastAt,
				Items:  items,
			})

			lastAt = at
			items = []LogLine{}
		}

		if first {
			lastAt = at
			first = false
		}

		items = append(items, line)
	}

	if len(items) > 0 {
		blocks = append(blocks, LogBlock{
			Header: lastAt,
			Items:  items,
		})
	}

	return blocks
}
