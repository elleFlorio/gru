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

	err := LoadServices(tmpdir)
	assert.NoError(t, err, "Loading should presente no errors")
	if assert.Len(t, services, 2, "Loaded services should be 2") {
		names := [2]string{services[0].Name, services[1].Name}
		assert.Contains(t, names, "service1", "The name of a service should be 'service1'")
		assert.Empty(t, services[1].Instances, "Instances should be empty")
	}
	CleanServices()
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

	s1 := GetServiceByName("service1")
	if assert.Len(t, s1, 1, "there should be 1 services with type service1") {
		assert.Equal(t, "service1", s1[0].Name, "service name should be service1")
	}

	s2 := GetServiceByName("service2")
	if assert.Len(t, s2, 1, "there should be 1 services with type service2") {
		assert.Equal(t, "service2", s2[0].Name, "service name should be service1")
	}

	s3 := GetServiceByName("service3")
	if assert.Len(t, s3, 1, "there should be 1 services with type service3") {
		assert.Equal(t, "service3", s3[0].Name, "service name should be service1")
	}

	CleanServices()
}

func TestGetServiceByImage(t *testing.T) {
	createMockServices()

	img1 := GetServiceByImage("test/tomcat")
	if assert.Len(t, img1, 1, "there should be 1 services with image test/tomcat") {
		assert.Equal(t, "test/tomcat", img1[0].Image, "service image should be test/tomcat")
	}

	img2 := GetServiceByImage("test/jetty")
	if assert.Len(t, img2, 1, "there should be 1 services with image test/jetty") {
		assert.Equal(t, "test/jetty", img2[0].Image, "service image should be test/jetty")
	}

	img3 := GetServiceByImage("test/mysql")
	if assert.Len(t, img3, 1, "there should be 1 services with imge test/mysql") {
		assert.Equal(t, "test/mysql", img3[0].Image, "service image should be test/mysql")
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

func createMockServices() {
	service1 := Service{
		Name:  "service1",
		Type:  "webserver",
		Image: "test/tomcat",
	}

	service2 := Service{
		Name:  "service2",
		Type:  "webserver",
		Image: "test/jetty",
	}

	service3 := Service{
		Name:  "service3",
		Type:  "database",
		Image: "test/mysql",
	}

	services = []Service{service1, service2, service3}
}
