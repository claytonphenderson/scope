package ingress

import (
	"bufio"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/claytonphenderson/scope/internal/models"
)

// Read all lines from std in.  Will increment wait group counter if
// one is provided
func ReadFromStdIn(ingressCh chan models.DbRecord, wg *sync.WaitGroup) {
	scanner := bufio.NewScanner(os.Stdin)
	buf := make([]byte, 0, 1024*1024)
	scanner.Buffer(buf, 1024*1024)

	// reads from stdin
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if !strings.Contains(line, "[event:") {
			continue
		}

		record := parseEvent(line)
		ingressCh <- record
		if wg != nil {
			wg.Add(1)
		}
	}
}

func parseEvent(line string) models.DbRecord {
	eventNameStartIndex := strings.Index(line, "[event:") + len("[event:")
	sub := strings.SplitAfter(line, "[event:")
	eventNameEndIndex := strings.Index(sub[1], "]") + len(sub[0])
	eventName := strings.TrimSpace(line[eventNameStartIndex:eventNameEndIndex])
	record := strings.TrimSpace(line[eventNameEndIndex+1:])

	dbRecord := models.DbRecord{
		Event:     eventName,
		Payload:   record,
		Timestamp: time.Now().UTC().String(),
	}

	return dbRecord
}
