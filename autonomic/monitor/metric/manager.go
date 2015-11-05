package metric

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"
	"regexp"
	"strconv"
	"strings"

	log "github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/Sirupsen/logrus"

	"github.com/elleFlorio/gru/enum"
	"github.com/elleFlorio/gru/storage"
)

type logEntry struct {
	service string
	metric  string
	value   float64
	unit    string
}

type servicesMap map[string]metricsMap

type metricsMap map[string][]float64

type MetricManager struct {
	ServiceMetrics servicesMap
	ch_manager     chan struct{}
	ch_data        chan logEntry
}

const sep string = ":"

var (
	manager         *MetricManager
	regex           = regexp.MustCompile("gru")
	ErrWrongLogLine = errors.New("Log line not well formed: 'gru:service:metric:value:unit")
)

func newManager() *MetricManager {
	concreteManager := MetricManager{
		make(map[string]metricsMap),
		make(chan struct{}),
		make(chan logEntry),
	}

	manager = &concreteManager
	return manager
}

func Manager() *MetricManager {
	if manager != nil {
		return manager
	}

	return newManager()
}

func (m *MetricManager) Start() {
	go m.startMetricManager()
}

func (m *MetricManager) startMetricManager() {
	var err error
	var e logEntry

	for {
		select {
		case e = <-m.ch_data:
			m.addValue(e)
		case <-m.ch_manager:
			err = m.saveMetrics()
			if err != nil {
				log.WithField("error", err).Errorln("Cannot save metrics data")
			}

			m.cleanMetrics()
		default:
		}
	}
}

func (m *MetricManager) addValue(entry logEntry) {
	var srv metricsMap
	var metric []float64
	var exists bool

	if srv, exists = m.ServiceMetrics[entry.service]; !exists {
		srv = make(metricsMap)
	}
	if metric, exists = srv[entry.metric]; !exists {
		metric = []float64{}
	}

	metric = append(metric, entry.value)
	srv[entry.metric] = metric
	m.ServiceMetrics[entry.service] = srv

	log.WithField("value", entry.value).Debugln("Added value to metric ", entry.metric)
}

func (m *MetricManager) saveMetrics() error {
	for srv, value := range m.ServiceMetrics {
		data, err := convertMetricsToData(value)
		if err != nil {
			return err
		} else {
			storeMetrics(srv, data)
		}
	}

	return nil
}

func convertMetricsToData(metrics metricsMap) ([]byte, error) {
	data, err := json.Marshal(metrics)
	if err != nil {
		log.Debugln("Error marshaling metrics data")
		return nil, err
	}

	return data, nil
}

func storeMetrics(name string, data []byte) {
	storage.StoreData(name, data, enum.METRICS)
}

func (m *MetricManager) cleanMetrics() {
	for name, _ := range m.ServiceMetrics {
		metrics := make(metricsMap)
		m.ServiceMetrics[name] = metrics
	}
}

func (m *MetricManager) GetMetricsOfService(name string) (metricsMap, error) {
	var err error

	m.ch_manager <- struct{}{}
	data, err := storage.GetData(name, enum.METRICS)
	if err != nil {
		log.WithField("error", err).Errorln("Cannot get metrics of service ", name)
		return nil, err
	}

	metrics, err := convertDataToMetrics(data)
	if err != nil {
		log.WithField("error", err).Errorln("Cannot get metrics of service", name)
		return nil, err
	}

	return metrics, nil
}

func convertDataToMetrics(data []byte) (metricsMap, error) {
	var metrics metricsMap
	err := json.Unmarshal(data, &metrics)
	if err != nil {
		log.WithField("error", err).Errorln("Error unmarshaling metrics data")
		return nil, err
	}

	return metrics, nil

}

func (m *MetricManager) StartCollector(contLog io.ReadCloser) {
	log.Debugln("starting collector")
	go collector(contLog, m.ch_data)
}

func collector(contLog io.ReadCloser, ch_data chan logEntry) {
	var err error
	var line []byte
	var data logEntry

	scanner := bufio.NewScanner(contLog)
	for scanner.Scan() {
		line = scanner.Bytes()
		if regex.Match(line) {
			log.WithField("line", string(line)).Debugln("found a match")
			data, err = getDataFromLogLine(string(line))
			if err != nil {
				log.WithField("error", err).Errorln("Error parsing container logs")
			} else {
				log.WithField("entry", data).Debugln("Sending data to manager")
				ch_data <- data
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
