package service

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadServices(t *testing.T) {
	tmpdir := createMockFiles()
	defer os.RemoveAll(tmpdir)

	result, err := LoadServices(tmpdir)
	assert.NoError(t, err, "Loading should presente no errors")
	if assert.Len(t, result, 2, "Loaded services should be 2") {
		names := [2]string{result[0].Name, result[1].Name}
		assert.Contains(t, names, "service1", "The name of a service should be 'service1'")
		assert.Empty(t, result[1].Instances, "Instances should be empty")
	}
	CleanServices()
}

func createMockFiles() string {
	mockService1 := `{
			"Name":"service1",
			"Type":"webserver",
			"Image":"test/tomcat",
			"Constraints":{
				"CpuMax":0.8,
				"CpuMin":0.3
			}

		}`

	mockService2 := `{
			"Name":"service2",
			"Type":"db",
			"Image":"test/mysql",
			"Constraints":{
				"MinActive":1,
				"MaxActive":3
			}
		}`

	tmpdir, err := ioutil.TempDir("", "gru_test_services")
	if err != nil {
		panic(err)
	}

	tmpfile1, err := ioutil.TempFile(tmpdir, "gru_test_services")
	if err != nil {
		panic(err)
	}

	tmpfile2, err := ioutil.TempFile(tmpdir, "gru_test_services")
	if err != nil {
		panic(err)
	}

	ioutil.WriteFile(tmpfile1.Name(), []byte(mockService1), 0644)
	ioutil.WriteFile(tmpfile2.Name(), []byte(mockService2), 0644)

	return tmpdir
}

func TestGetServiceByType(t *testing.T) {
	createMockServices()

	ws := GetServiceByType("webserver")
	if assert.Len(t, ws, 2, "there should be 2 services with type webserver") {
		assert.Equal(t, "webserver", ws[0].Type, "service type should be webserver")
		img := [2]string{ws[0].Image, ws[1].Image}
		assert.Contains(t, img, "test/tomcat", "images should contain test/tomcat")
	}

	db := GetServiceByType("database")
	if assert.Len(t, db, 1, "there should be 1 services with type database") {
		assert.Equal(t, "database", db[0].Type, "service type should be database")
	}

	ap := GetServiceByType("application")
	assert.Len(t, ap, 0, "there should be 0 services with type application")

	CleanServices()
}

func TestGetServiceByName(t *testing.T) {
	createMockServices()

	s1, err := GetServiceByName("service2")
	assert.Equal(t, "service2", s1.Name, "service name should be service1")

	_, err = GetServiceByName("pippo")
	assert.Error(t, err, "There should be no service with name 'pippo'")

	CleanServices()
}

func TestGetServiceByImage(t *testing.T) {
	createMockServices()

	img1, err := GetServiceByImage("test/mysql")
	assert.Equal(t, "test/mysql", img1.Image, "service image should be test/tomcat")

	_, err = GetServiceByImage("test/pippo")
	assert.Error(t, err, "There should be no image 'test/pippo'")

	CleanServices()
}

func TestGetServiceByInstanceId(t *testing.T) {
	createMockServices()

	s1, _ := GetServiceByInstanceId("000000")
	assert.Equal(t, "service2", s1.Name, "Service name should be 'service3'")

	_, err := GetServiceByInstanceId("111111")
	assert.Error(t, err, "There should be no service with instance id '111111'")

	CleanServices()
}

func createMockServices() {
	service1 := Service{
		Name:  "service1",
		Type:  "webserver",
		Image: "test/tomcat",
	}

	service2 := Service{
		Name:      "service2",
		Type:      "webserver",
		Image:     "test/jetty",
		Instances: map[string]Instance{"000000": Instance{}},
	}

	service3 := Service{
		Name:  "service3",
		Type:  "database",
		Image: "test/mysql",
	}

	services = []Service{service1, service2, service3}
}
