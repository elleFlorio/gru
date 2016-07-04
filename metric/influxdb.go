package metric

import (
	"errors"
	"time"

	log "github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/Sirupsen/logrus"
	"github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/influxdb/influxdb/client/v2"

	"github.com/elleFlorio/gru/enum"
	"github.com/elleFlorio/gru/utils"
)

type influxdb struct {
	influx client.Client
	config *influxdbConfig
	batch  client.BatchPoints
}

type influxdbConfig struct {
	Url      string
	DbName   string
	Username string
	Password string
}

func (db *influxdb) Name() string {
	return "influxdb"
}

func (db *influxdb) Initialize(config map[string]interface{}) error {
	var err error
	db.config = &influxdbConfig{}
	err = utils.FillStruct(db.config, config)
	if err != nil {
		return err
	}

	log.Debugln("Initializing influxdb at address: ", db.config.Url)

	db.influx, err = client.NewHTTPClient(client.HTTPConfig{
		Addr:     db.config.Url,
		Username: db.config.Username,
		Password: db.config.Password,
	})
	if err != nil {
		return err
	}

	db.batch, err = client.NewBatchPoints(client.BatchPointsConfig{
		Database:  db.config.DbName,
		Precision: "s",
	})
	if err != nil {
		return err
	}

	return nil
}

func (db *influxdb) StoreMetrics(metrics GruMetric) error {
	points, err := createInfluxMetrics(metrics)
	if err != nil {
		log.Errorln("Error storing Influx metrics")
		return err
	}

	for _, point := range points {
		db.batch.AddPoint(point)
	}

	return db.influx.Write(db.batch)
}

func createInfluxMetrics(metrics GruMetric) ([]*client.Point, error) {
	var err error
	points := []*client.Point{}
	nodeName := metrics.Node.Name

	nodePoint, err := createInfluxNode(metrics.Node)
	if err != nil {
		log.WithField("err", err).Errorln("Error creating node metrics")
	} else {
		points = append(points, nodePoint)
	}

	for _, service := range metrics.Service {
		log.WithField("service", service.Name).Debugln("Computing influx metrics of service")
		statusPoint, err := createInfluxInstanceStatus(nodeName, service)
		if err != nil {
			log.WithFields(log.Fields{
				"err":     err,
				"service": service,
			}).Warnln("Error creating instance status metrics")
		} else {
			points = append(points, statusPoint)
		}

		cpuPoint, err := createInfluxCpuService(nodeName, service)
		if err != nil {
			log.WithFields(log.Fields{
				"err":     err,
				"service": service,
			}).Warnln("Error creating cpu metrics")
		} else {
			points = append(points, cpuPoint)
		}

		memPoint, err := createInfluxMemService(nodeName, service)
		if err != nil {
			log.WithFields(log.Fields{
				"err":     err,
				"service": service,
			}).Warnln("Error creating memory metrics")
		} else {
			points = append(points, memPoint)
		}

		userStatsPoints, err := createInfluxUserStatsService(nodeName, service)
		if err != nil {
			log.WithFields(log.Fields{
				"err":     err,
				"service": service,
			}).Warnln("Error creating user stats metrics")
		} else {
			points = append(points, userStatsPoints)
		}

		userAnalyticsPoints, err := createInfluxUserAnalyticsService(nodeName, service)
		if err != nil {
			log.WithFields(log.Fields{
				"err":     err,
				"service": service,
			}).Warnln("Error creating user analytics metrics")
		} else {
			points = append(points, userAnalyticsPoints)
		}

		userSharedPoints, err := createInfluxUserSharedService(nodeName, service)
		if err != nil {
			log.WithFields(log.Fields{
				"err":     err,
				"service": service,
			}).Warnln("Error creating user shared metrics")
		} else {
			points = append(points, userSharedPoints)
		}
	}

	policyPoint, err := createInfluxPolicy(nodeName, metrics.Policy)
	if err != nil {
		log.WithField("err", err).Errorln("Error creating policy metrics")
	} else {
		points = append(points, policyPoint)
	}

	return points, nil
}

func createInfluxNode(node NodeMetrics) (*client.Point, error) {
	tags := map[string]string{
		"node": node.UUID,
		"name": node.Name,
	}
	fields := map[string]interface{}{
		"cpu":             node.Stats.BaseMetrics[enum.METRIC_CPU_AVG.ToString()],
		"memory":          node.Stats.BaseMetrics[enum.METRIC_MEM_AVG.ToString()],
		"cores_total":     node.Resources.CPU.Total,
		"cores_free":      node.Resources.CPU.Availabe,
		"active_services": node.ActiveServices,
	}

	point, err := client.NewPoint("node", tags, fields, time.Now())
	if err != nil {
		return nil, err
	}

	log.WithField("series", "node").Debugln("Created influx metrics")

	return point, nil
}

func createInfluxInstanceStatus(nodeName string, service ServiceMetric) (*client.Point, error) {
	tags := map[string]string{
		"node":          nodeName,
		"service_name":  service.Name,
		"service_type":  service.Type,
		"service_image": service.Image,
	}
	fields := map[string]interface{}{
		"all":     service.Instances.All,
		"pending": service.Instances.Pending,
		"running": service.Instances.Running,
		"stopped": service.Instances.Stopped,
		"paused":  service.Instances.Paused,
	}

	point, err := client.NewPoint("instance_status", tags, fields, time.Now())
	if err != nil {
		return nil, err
	}

	log.WithField("series", "instance_status").Debugln("Created influx metrics")

	return point, nil
}

func createInfluxCpuService(nodeName string, service ServiceMetric) (*client.Point, error) {
	tags := map[string]string{
		"node":          nodeName,
		"service_name":  service.Name,
		"service_type":  service.Type,
		"service_image": service.Image,
	}
	fields := map[string]interface{}{
		"stats_value":     service.Stats.BaseMetrics[enum.METRIC_CPU_AVG.ToString()],
		"analytics_value": service.Analytics.BaseAnalytics[enum.METRIC_CPU_AVG.ToString()],
		"shared_value":    service.Shared.BaseShared[enum.METRIC_CPU_AVG.ToString()],
	}

	point, err := client.NewPoint("cpu_service", tags, fields, time.Now())
	if err != nil {
		return nil, err
	}

	log.WithField("series", "cpu_service").Debugln("Created influx metrics")

	return point, nil
}

func createInfluxMemService(nodeName string, service ServiceMetric) (*client.Point, error) {
	tags := map[string]string{
		"node":          nodeName,
		"service_name":  service.Name,
		"service_type":  service.Type,
		"service_image": service.Image,
	}
	fields := map[string]interface{}{
		"stats_value":     service.Stats.BaseMetrics[enum.METRIC_MEM_AVG.ToString()],
		"analytics_value": service.Analytics.BaseAnalytics[enum.METRIC_MEM_AVG.ToString()],
		"shared_value":    service.Shared.BaseShared[enum.METRIC_MEM_AVG.ToString()],
	}

	point, err := client.NewPoint("memory_service", tags, fields, time.Now())
	if err != nil {
		return nil, err
	}

	log.WithField("series", "memory_service").Debugln("Created influx metrics")

	return point, nil
}

func createInfluxUserStatsService(nodeName string, service ServiceMetric) (*client.Point, error) {
	tags := map[string]string{
		"node":          nodeName,
		"service_name":  service.Name,
		"service_type":  service.Type,
		"service_image": service.Image,
	}
	fields := make(map[string]interface{}, len(service.Stats.UserMetrics))
	for metric, value := range service.Stats.UserMetrics {
		fields[metric] = value
	}

	if len(fields) == 0 {
		return nil, errors.New("No user stats")
	}

	point, err := client.NewPoint("user_stat_service", tags, fields, time.Now())
	if err != nil {
		return nil, err
	}

	log.WithField("series", "user_stat_service").Debugln("Created influx metrics")

	return point, nil
}

func createInfluxUserAnalyticsService(nodeName string, service ServiceMetric) (*client.Point, error) {
	tags := map[string]string{
		"node":          nodeName,
		"service_name":  service.Name,
		"service_type":  service.Type,
		"service_image": service.Image,
	}
	fields := make(map[string]interface{}, len(service.Analytics.UserAnalytics))

	for analytic, value := range service.Analytics.UserAnalytics {
		fields[analytic] = value
	}

	if len(fields) == 0 {
		return nil, errors.New("No user analytic")
	}

	point, err := client.NewPoint("user_analytics_service", tags, fields, time.Now())
	if err != nil {
		return nil, err
	}

	log.WithField("series", "user_analytics_service").Debugln("Created influx metrics")

	return point, nil
}

func createInfluxUserSharedService(nodeName string, service ServiceMetric) (*client.Point, error) {
	tags := map[string]string{
		"node":          nodeName,
		"service_name":  service.Name,
		"service_type":  service.Type,
		"service_image": service.Image,
	}
	fields := make(map[string]interface{}, len(service.Shared.UserShared))

	for shared, value := range service.Shared.UserShared {
		fields[shared] = value
	}

	if len(fields) == 0 {
		return nil, errors.New("No user analytic")
	}

	point, err := client.NewPoint("user_shared_service", tags, fields, time.Now())
	if err != nil {
		return nil, err
	}

	log.WithField("series", "user_shared_service").Debugln("Created influx metrics")

	return point, nil
}

func createInfluxPolicy(nodeName string, plc PolicyMetric) (*client.Point, error) {
	tags := map[string]string{
		"node": nodeName,
		"name": plc.Name,
	}
	fields := map[string]interface{}{
		"weight": plc.Weight,
	}

	point, err := client.NewPoint("policy", tags, fields, time.Now())
	if err != nil {
		return nil, err
	}

	log.WithField("series", "policy").Debugln("Created influx metrics")

	return point, nil
}
