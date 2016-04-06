package metric

import (
	"time"

	log "github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/Sirupsen/logrus"
	"github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/influxdb/influxdb/client/v2"

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
		return points, err
	}
	points = append(points, nodePoint)

	for _, service := range metrics.Service {
		log.WithField("service", service.Name).Debugln("Computing influx metrics of service")
		statusPoint, err := createInfluxInstanceStatus(nodeName, service)
		if err != nil {
			log.WithField("err", err).Errorln("Error creating instance status metrics for Service ", service.Name)
			return points, err
		}
		points = append(points, statusPoint)

		cpuPoint, err := createInfluxCpuService(nodeName, service)
		if err != nil {
			log.WithField("err", err).Errorln("Error creating cpu metrics for Service ", service.Name)
			return points, err
		}
		points = append(points, cpuPoint)

		memPoint, err := createInfluxMemService(nodeName, service)
		if err != nil {
			log.WithField("err", err).Errorln("Error creating memory metrics for Service ", service.Name)
			return points, err
		}
		points = append(points, memPoint)

		loadPoint, err := createInfluxLoadService(nodeName, service)
		if err != nil {
			log.WithField("err", err).Errorln("Error creating load metrics for Service ", service.Name)
			return points, err
		}
		points = append(points, loadPoint)

		healthPoint, err := createInfluxHealthService(nodeName, service)
		if err != nil {
			log.WithField("err", err).Errorln("Error creating health metrics for Service ", service.Name)
			return points, err
		}
		points = append(points, healthPoint)

	}

	policyPoint, err := createInfluxPolicy(nodeName, metrics.Policy)
	if err != nil {
		log.WithField("err", err).Errorln("Error creating policy metrics")
		return points, err
	}
	points = append(points, policyPoint)

	return points, nil
}

func createInfluxNode(node NodeMetrics) (*client.Point, error) {
	tags := map[string]string{
		"node": node.UUID,
		"name": node.Name,
	}
	fields := map[string]interface{}{
		"cpu":    node.Cpu,
		"memory": node.Memory,
		"health": node.Health,
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
		"stats_tot":       service.Stats.CpuTot,
		"stats_avg":       service.Stats.CpuAvg,
		"analytics_value": service.Analytics.Cpu,
		"shared_value":    service.Shared.Cpu,
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
		"stats_tot":       service.Stats.MemTot,
		"stats_avg":       service.Stats.MemAvg,
		"analytics_value": service.Analytics.Memory,
		"shared_value":    service.Shared.Memory,
	}

	point, err := client.NewPoint("memory_service", tags, fields, time.Now())
	if err != nil {
		return nil, err
	}

	log.WithField("series", "memory_service").Debugln("Created influx metrics")

	return point, nil
}

func createInfluxLoadService(nodeName string, service ServiceMetric) (*client.Point, error) {
	tags := map[string]string{
		"node":          nodeName,
		"service_name":  service.Name,
		"service_type":  service.Type,
		"service_image": service.Image,
	}
	fields := map[string]interface{}{
		"analytics_value": service.Analytics.Load,
		"shared_value":    service.Shared.Load,
	}

	point, err := client.NewPoint("load_service", tags, fields, time.Now())
	if err != nil {
		return nil, err
	}

	log.WithField("series", "load_service").Debugln("Created influx metrics")

	return point, nil
}

func createInfluxHealthService(nodeName string, service ServiceMetric) (*client.Point, error) {
	tags := map[string]string{
		"node":          nodeName,
		"service_name":  service.Name,
		"service_type":  service.Type,
		"service_image": service.Image,
	}
	fields := map[string]interface{}{
		"value": service.Analytics.Health,
	}

	point, err := client.NewPoint("health_service", tags, fields, time.Now())
	if err != nil {
		return nil, err
	}

	log.WithField("series", "health_service").Debugln("Created influx metrics")

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
