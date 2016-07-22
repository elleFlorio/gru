package logreader

import (
	"errors"
	"io"
	"regexp"

	log "github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/Sirupsen/logrus"

	mtr "github.com/elleFlorio/gru/autonomic/monitor/metric"
	"github.com/elleFlorio/gru/utils"
)

type logEntry struct {
	service string
	metric  string
	value   float64
	unit    string
}

type srvMetricBuffer struct {
	userDef map[string]utils.Buffer
}

// TODO this should be done in a better way.
var c_B_SIZE = 1

var (
	srvMetrics      map[string]srvMetricBuffer
	ch_entry        chan logEntry
	ch_stop         chan struct{}
	regex           = regexp.MustCompile("gru")
	ErrWrongLogLine = errors.New("Log line not well formed: 'gru:service:metric:value:unit")
)

func init() {
	ch_entry = make(chan logEntry)
	ch_stop = make(chan struct{})
	srvMetrics = make(map[string]srvMetricBuffer)
}

func Initialize(services []string) {
	for _, service := range services {
		srvMetrics[service] = srvMetricBuffer{
			userDef: make(map[string]utils.Buffer),
		}
	}
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

	metrics := srvMetrics[entry.service]
	if _, ok := metrics.userDef[entry.metric]; !ok {
		newMetric := utils.BuildBuffer(c_B_SIZE)
		metrics.userDef[entry.metric] = newMetric
	}

	metric := metrics.userDef[entry.metric]
	values := metric.PushValue(entry.value)
	if values != nil {
		mtr.UpdateUserMetric(entry.service, entry.metric, values)
	}

	metrics.userDef[entry.metric] = metric
	srvMetrics[entry.service] = metrics
}

func StartCollector(contLog io.ReadCloser) {
	log.Debugln("starting collector")
	go collector(contLog, ch_entry)
}

func Stop() {
	ch_stop <- struct{}{}
}
