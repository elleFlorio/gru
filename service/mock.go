package service

import (
	"io/ioutil"
)

const nFILES = 3

func NumberOfServiceFiles() int {
	return nFILES
}

func CreateMockFiles() string {
	mockService1 := `{
            "Name":"service1",
            "Type":"webserver",
            "Image":"test/tomcat",
            "Constraints":{
                "CpuMax":0.3,
                "CpuMin":0.3
            }

        }`

	mockService2 := `{
            "Name":"service2",
            "Type":"webserver",
            "Image":"test/jetty",
            "Constraints":{
                "MinActive":1,
                "MaxActive":3
            }
        }`

	mockService3 := `{
            "Name":"service3",
            "Type":"database",
            "Image":"test/mysql",
            "Constraints":{
                "MinActive":1
            }
        }`

	tmpdir, err := ioutil.TempDir("", "gru_test_services")
	if err != nil {
		panic(err)
	}

	tmpfile1, err := ioutil.TempFile(tmpdir, "gru_test_services1")
	if err != nil {
		panic(err)
	}

	tmpfile2, err := ioutil.TempFile(tmpdir, "gru_test_services2")
	if err != nil {
		panic(err)
	}

	tmpfile3, err := ioutil.TempFile(tmpdir, "gru_test_services3")
	if err != nil {
		panic(err)
	}

	ioutil.WriteFile(tmpfile1.Name(), []byte(mockService1), 0644)
	ioutil.WriteFile(tmpfile2.Name(), []byte(mockService2), 0644)
	ioutil.WriteFile(tmpfile3.Name(), []byte(mockService3), 0644)

	return tmpdir
}

func CreateMockServices() []Service {
	service1 := Service{
		Name:  "service1",
		Type:  "webserver",
		Image: "test/tomcat",
		Constraints: Constraints{
			CpuMin:    0.4,
			CpuMax:    0.6,
			MinActive: 1,
			MaxActive: 3,
		},
	}

	service2 := Service{
		Name:  "service2",
		Type:  "webserver",
		Image: "test/jetty",
		Constraints: Constraints{
			CpuMin:    0.4,
			CpuMax:    0.7,
			MinActive: 1,
			MaxActive: 2,
		},
	}

	service3 := Service{
		Name:  "service3",
		Type:  "database",
		Image: "test/mysql",
		Constraints: Constraints{
			MinActive: 2,
		},
	}

	return []Service{service1, service2, service3}
}
