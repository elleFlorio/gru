package logreader

import (
	"bufio"
	"io"
	"strconv"
	"strings"

	log "github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/Sirupsen/logrus"
)

const c_SEP string = ":"

func collector(contLog io.ReadCloser, ch_entry chan logEntry) {
	var err error
	var line []byte
	var data logEntry

	scanner := bufio.NewScanner(contLog)
	for scanner.Scan() {
		line = scanner.Bytes()
		if regex.Match(line) {
			data, err = getDataFromLogLine(string(line))
			if err != nil {
				log.WithField("err", err).Errorln("Error parsing container logs")
			} else {
				ch_entry <- data
			}
		}
	}

	if err = scanner.Err(); err != nil {
		log.WithField("err", err).Errorln("Error in scanner.")
	}

	log.Debugln("Stopped collector")
}

func getDataFromLogLine(line string) (logEntry, error) {
	relevant := line[strings.LastIndex(line, "gru"):]
	data := strings.Split(relevant, c_SEP)
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
