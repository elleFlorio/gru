package logreader

import (
	"errors"
	"io"
	"regexp"

	log "github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/Sirupsen/logrus"

	mtr "github.com/elleFlorio/gru/autonomic/monitor/metric"
)

type logEntry struct {
	service string
	metric  string
	value   float64
	unit    string
}

var (
	ch_entry        chan logEntry
	ch_stop         chan struct{}
	regex           = regexp.MustCompile("gru")
	ErrWrongLogLine = errors.New("Log line not well formed: 'gru:service:metric:value:unit")
)

func init() {
	ch_entry = make(chan logEntry)
	ch_stop = make(chan struct{})
}

func StartLogReader() {
	go logReading()
}

func logReading() {
	var e logEntry

	for {
		select {
		case e = <-ch_entry:
			addValue(e)
		case <-ch_stop:
			return
		default:
		}
	}
}

func addValue(entry logEntry) {
	if entry.value < 0.0 {
		log.WithFields(log.Fields{
			"service": entry.service,
			"metric":  entry.metric,
			"value":   entry.value,
		}).Warnln("Metric value < 0")
		return
	}

	values := make([]float64, 0, 1)
	values = append(values, entry.value)
	mtr.UpdateUserMetric(entry.service, entry.metric, values)

}

func StartCollector(contLog io.ReadCloser) {
	log.Debugln("starting collector")
	go collector(contLog, ch_entry)
}

func Stop() {
	ch_stop <- struct{}{}
}
