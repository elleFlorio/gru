package metric

import (
	"bufio"
	"io"
	"strconv"
	"strings"

	log "github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/Sirupsen/logrus"
)

func collector(contLog io.ReadCloser, ch_data chan logEntry, ch_stop chan struct{}) {
	var err error
	var line []byte
	var data logEntry

	scanner := bufio.NewScanner(contLog)
	for scanner.Scan() {
		select {
		case <-ch_stop:
			return
		default:
			line = scanner.Bytes()
			if regex.Match(line) {
				data, err = getDataFromLogLine(string(line))
				if err != nil {
					log.WithField("error", err).Errorln("Error parsing container logs")
				} else {
					ch_data <- data
				}
			}
		}
	}

	if err = scanner.Err(); err != nil {
		log.WithField("error", err).Errorln("Error in scanner.")
	}
}

func getDataFromLogLine(line string) (logEntry, error) {
	relevant := line[strings.LastIndex(line, "gru"):]
	data := strings.Split(relevant, sep)
	if len(data) < 5 {
		return logEntry{}, ErrWrongLogLine
	}

	service := data[1]
	metric := data[2]
	unit := data[4]
	value, err := strconv.ParseFloat(data[3], 64)
	if err != nil {
		return logEntry{}, err
	}

	entry := logEntry{service, metric, value, unit}

	return entry, nil
}
