package metric

import (
	"time"

	log "github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/Sirupsen/logrus"
	"github.com/influxdb/influxdb/client/v2"
)

type influxdb struct {
	influx client.Client
	config influxdbConfig
	batch  client.BatchPoints
}

type influxdbConfig struct {
	url      string
	dbName   string
	username string
	password string
}

func (db *influxdb) Name() string {
	return "influxdb"
}

func (db *influxdb) Initialize(config interface{}) error {
	var err error
	db.config = config.(influxdbConfig)
	db.influx, err = client.NewHTTPClient(client.HTTPConfig{
		Addr:     db.config.url,
		Username: db.config.username,
		Password: db.config.password,
	})
	if err != nil {
		return err
	}

	db.batch, err = client.NewBatchPoints(client.BatchPointsConfig{
		Database:  db.config.dbName,
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
		log.WithField("err", err).Errorln("Error creating metrics for Node")
		return points, err
	}
	points = append(points, nodePoint)

	for _, service := range metrics.Service {
		servicePoint, err := createInfluxService(nodeUUID, service)
		if err != nil {
			log.WithField("err", err).Errorln("Error creating service metrics for Service ", service.Name)
			return points, err
		}
		points = append(points, servicePoint)

		statsPoint, err := createInfluxService(nodeUUID, service)
		if err != nil {
			log.WithField("err", err).Errorln("Error creating stats metrics for Service ", service.Name)
			return points, err
		}
		points = append(points, statsPoint)

		analyticsPoint, err := createInfluxService(nodeUUID, service)
		if err != nil {
			log.WithField("err", err).Errorln("Error creating analytics metrics for Service ", service.Name)
			return points, err
		}
		points = append(points, analyticsPoint)
	}

	planPoint, err := createInfluxPlans(nodeUUID, metrics.Plan)
	if err != nil {
		log.WithField("err", err).Errorln("Error creating metrics for Plan")
		return points, err
	}
	points = append(points, planPoint)

	actionPoint, err := createInfluxActions(nodeUUID, metrics.Action)
	if err != nil {
		log.WithField("err", err).Errorln("Error creating metrics for Action")
		return points, err
	}
	points = append(points, actionPoint)

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

	return point, nil
}

func createInfluxService(nodeUUID string, service ServiceMetric) (*client.Point, error) {
	tags := map[string]string{
		"node":          nodeUUID,
		"service_name":  service.Name,
		"service_type":  service.Type,
		"service_image": service.Image,
	}
	fields := map[string]interface{}{
		"pending": service.Instances.Pending,
		"running": service.Instances.Running,
		"stopped": service.Instances.Stopped,
		"paused":  service.Instances.Paused,
	}

	point, err := client.NewPoint("service", tags, fields, time.Now())
	if err != nil {
		return nil, err
	}

	return point, nil
}

func createInfluxStats(nodeUUID string, service ServiceMetric) (*client.Point, error) {
	tags := map[string]string{
		"node":          nodeUUID,
		"service_name":  service.Name,
		"service_type":  service.Type,
		"service_image": service.Image,
	}
	fields := map[string]interface{}{
		"cpu_tot": service.Stats.CpuTot,
		"cpu_avg": service.Stats.CpuAvg,
		"mem_tot": service.Stats.MemTot,
		"mem_avg": service.Stats.MemAvg,
	}

	point, err := client.NewPoint("stats", tags, fields, time.Now())
	if err != nil {
		return nil, err
	}

	return point, nil
}

func createInfluxAnalytics(nodeUUID string, service ServiceMetric) (*client.Point, error) {
	tags := map[string]string{
		"node":          nodeUUID,
		"service_name":  service.Name,
		"service_type":  service.Type,
		"service_image": service.Image,
	}
	fields := map[string]interface{}{
		"load":      service.Analytics.Load,
		"cpu":       service.Analytics.Cpu,
		"memory":    service.Analytics.Memory,
		"resources": service.Analytics.Resources,
		"health":    service.Analytics.Health,
	}

	point, err := client.NewPoint("analytics", tags, fields, time.Now())
	if err != nil {
		return nil, err
	}

	return point, nil
}

func createInfluxPlans(nodeUUID string, plans PlansMetric) (*client.Point, error) {
	tags := map[string]string{
		"node":   nodeUUID,
		"policy": plans.Policy,
	}
	fields := map[string]interface{}{
		"target": plans.Target,
		"weight": plans.Weight,
	}

	point, err := client.NewPoint("plans", tags, fields, time.Now())
	if err != nil {
		return nil, err
	}

	return point, nil
}

func createInfluxActions(nodeUUID string, actions ActionsMetric) (*client.Point, error) {
	tags := map[string]string{
		"node":   nodeUUID,
		"policy": actions.Action,
	}
	fields := map[string]interface{}{
		"target": actions.Target,
	}

	point, err := client.NewPoint("actions", tags, fields, time.Now())
	if err != nil {
		return nil, err
	}

	return point, nil
}
