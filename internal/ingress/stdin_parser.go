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

	// reads from stdin
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if !strings.HasPrefix(line, "[event:") {
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
	eventNameEndIndex := strings.Index(line, "]")
	eventName := strings.TrimSpace(line[7:eventNameEndIndex])
	record := strings.TrimSpace(line[eventNameEndIndex+1:])
	// logger.Info(eventName)
	// logger.Info(record)

	dbRecord := models.DbRecord{
		Event:     eventName,
		Payload:   record,
		Timestamp: time.Now().UTC().String(),
	}

	return dbRecord
}
