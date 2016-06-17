package logreader

// import (
// 	"errors"
// 	"io"
// 	"regexp"

// 	log "github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/Sirupsen/logrus"
// )

// type logEntry struct {
// 	service string
// 	metric  string
// 	value   float64
// 	unit    string
// }

// type servicesMap map[string]metricsMap

// type metricsMap map[string][]float64

// type MetricManager struct {
// 	ServiceMetrics servicesMap
// 	ch_notify      chan struct{}
// 	ch_get         chan servicesMap
// 	ch_data        chan logEntry
// 	ch_stop        chan struct{}
// }

// const sep string = ":"
// const c_CH_BUFFER_SIZE int = 10

// var (
// 	manager         *MetricManager
// 	regex           = regexp.MustCompile("gru")
// 	ErrWrongLogLine = errors.New("Log line not well formed: 'gru:service:metric:value:unit")
// )

// func newManager() *MetricManager {
// 	concreteManager := MetricManager{
// 		make(map[string]metricsMap),
// 		make(chan struct{}),
// 		make(chan servicesMap),
// 		make(chan logEntry, c_CH_BUFFER_SIZE),
// 		make(chan struct{}),
// 	}

// 	manager = &concreteManager
// 	return manager
// }

// func Manager() *MetricManager {
// 	if manager != nil {
// 		return manager
// 	}

// 	return newManager()
// }

// func (m *MetricManager) Start() {
// 	go m.startLogReaderManager()
// }

// func (m *MetricManager) startLogReaderManager() {
// 	var e logEntry

// 	for {
// 		select {
// 		case e = <-m.ch_data:
// 			m.addValue(e)
// 		case <-m.ch_notify:
// 			metricsCopy := m.copyMetrics()
// 			m.ch_get <- metricsCopy
// 			m.cleanMetrics()
// 		case <-m.ch_stop:
// 			return
// 		default:
// 		}
// 	}
// }

// func (m *MetricManager) addValue(entry logEntry) {
// 	var srv metricsMap
// 	var metric []float64
// 	var exists bool

// 	if entry.value < 0.0 {
// 		log.WithFields(log.Fields{
// 			"service": entry.service,
// 			"metric":  entry.metric,
// 			"value":   entry.value,
// 		}).Warnln("Metric value < 0")
// 		return
// 	}

// 	if srv, exists = m.ServiceMetrics[entry.service]; !exists {
// 		srv = make(metricsMap)
// 	}
// 	if metric, exists = srv[entry.metric]; !exists {
// 		metric = []float64{}
// 	}

// 	metric = append(metric, entry.value)
// 	srv[entry.metric] = metric
// 	m.ServiceMetrics[entry.service] = srv
// }

// func (m *MetricManager) copyMetrics() servicesMap {
// 	metricsCopy := make(servicesMap)
// 	for k, v := range m.ServiceMetrics {
// 		metricsCopy[k] = v
// 	}

// 	return metricsCopy
// }

// func (m *MetricManager) cleanMetrics() {
// 	for name, _ := range m.ServiceMetrics {
// 		metrics := make(metricsMap)
// 		m.ServiceMetrics[name] = metrics
// 	}
// }

// func (m *MetricManager) StartCollector(contLog io.ReadCloser) {
// 	log.Debugln("starting collector")
// 	go collector(contLog, m.ch_data)
// }

// func (m *MetricManager) GetMetrics() servicesMap {
// 	m.ch_notify <- struct{}{}
// 	metrics := <-m.ch_get
// 	return metrics
// }

// func (m *MetricManager) Stop() {
// 	m.ch_stop <- struct{}{}
// }
