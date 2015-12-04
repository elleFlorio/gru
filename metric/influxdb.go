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
	utils.FillStruct(db.config, config)

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
	nodeUUID := metrics.Node.UUID

	nodePoint, err := createInfluxNode(metrics.Node)
	if err != nil {
		log.WithField("err", err).Errorln("Error creating node metrics")
		return points, err
	}
	points = append(points, nodePoint)

	for _, service := range metrics.Service {
		statusPoint, err := createInfluxInstanceStatus(nodeUUID, service)
		if err != nil {
			log.WithField("err", err).Errorln("Error creating instance status metrics for Service ", service.Name)
			return points, err
		}
		points = append(points, statusPoint)

		cpuPoint, err := createInfluxCpuService(nodeUUID, service)
		if err != nil {
			log.WithField("err", err).Errorln("Error creating cpu metrics for Service ", service.Name)
			return points, err
		}
		points = append(points, cpuPoint)

		memPoint, err := createInfluxMemService(nodeUUID, service)
		if err != nil {
			log.WithField("err", err).Errorln("Error creating memory metrics for Service ", service.Name)
			return points, err
		}
		points = append(points, memPoint)

		loadPoint, err := createInfluxLoadService(nodeUUID, service)
		if err != nil {
			log.WithField("err", err).Errorln("Error creating load metrics for Service ", service.Name)
			return points, err
		}
		points = append(points, loadPoint)

		healthPoint, err := createInfluxHealthService(nodeUUID, service)
		if err != nil {
			log.WithField("err", err).Errorln("Error creating health metrics for Service ", service.Name)
			return points, err
		}
		points = append(points, healthPoint)

	}

	policyPoint, err := createInfluxPolicy(nodeUUID, metrics.Plan)
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

func createInfluxInstanceStatus(nodeUUID string, service ServiceMetric) (*client.Point, error) {
	tags := map[string]string{
		"node":          nodeUUID,
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

func createInfluxCpuService(nodeUUID string, service ServiceMetric) (*client.Point, error) {
	tags := map[string]string{
		"node":          nodeUUID,
		"service_name":  service.Name,
		"service_type":  service.Type,
		"service_image": service.Image,
	}
	fields := map[string]interface{}{
		"stats_tot":       service.Stats.CpuTot,
		"stats_avg":       service.Stats.CpuAvg,
		"analytics_value": service.Analytics.Cpu,
	}

	point, err := client.NewPoint("cpu_service", tags, fields, time.Now())
	if err != nil {
		return nil, err
	}

	log.WithField("series", "cpu_service").Debugln("Created influx metrics")

	return point, nil
}

func createInfluxMemService(nodeUUID string, service ServiceMetric) (*client.Point, error) {
	tags := map[string]string{
		"node":          nodeUUID,
		"service_name":  service.Name,
		"service_type":  service.Type,
		"service_image": service.Image,
	}
	fields := map[string]interface{}{
		"stats_tot":       service.Stats.MemTot,
		"stats_avg":       service.Stats.MemAvg,
		"analytics_value": service.Analytics.Memory,
	}

	point, err := client.NewPoint("memory_service", tags, fields, time.Now())
	if err != nil {
		return nil, err
	}

	log.WithField("series", "memory_service").Debugln("Created influx metrics")

	return point, nil
}

func createInfluxLoadService(nodeUUID string, service ServiceMetric) (*client.Point, error) {
	tags := map[string]string{
		"node":          nodeUUID,
		"service_name":  service.Name,
		"service_type":  service.Type,
		"service_image": service.Image,
	}
	fields := map[string]interface{}{
		"value": service.Analytics.Load,
	}

	point, err := client.NewPoint("load_service", tags, fields, time.Now())
	if err != nil {
		return nil, err
	}

	log.WithField("series", "load_service").Debugln("Created influx metrics")

	return point, nil
}

func createInfluxHealthService(nodeUUID string, service ServiceMetric) (*client.Point, error) {
	tags := map[string]string{
		"node":          nodeUUID,
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

func createInfluxPolicy(nodeUUID string, plans PlansMetric) (*client.Point, error) {
	tags := map[string]string{
		"node":   nodeUUID,
		"target": plans.Policy,
	}
	fields := map[string]interface{}{
		"policy": plans.Policy,
		"weight": plans.Weight,
	}

	point, err := client.NewPoint("policy", tags, fields, time.Now())
	if err != nil {
		return nil, err
	}

	log.WithField("series", "policy").Debugln("Created influx metrics")

	return point, nil
}
