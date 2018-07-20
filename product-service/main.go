package main

import (
	"net/http"

	"github.com/labstack/echo"
	"fmt"
	"os"
	"log"

	consulapi "github.com/hashicorp/consul/api"
	"io/ioutil"
)


var SERVER_NAME = "product-service"

type Node struct {
	Address string
	Port int
	serviceID string
	serviceName string
}

func main() {
	server := echo.New()
	server.GET("/product-service/healthcheck", func(c echo.Context) error {
		return c.String(http.StatusOK, "Product service is good")
	})

	server.GET("/product-service/unregister-service", func(c echo.Context) error {
		DeRegisterServiceWithConsul(SERVER_NAME)
		return c.String(http.StatusOK, "Product service has been unregister")
	})

	server.GET("/product-service/product/list/:name", func(c echo.Context) error {
		name := c.Param("name")

		serviceInfo, _ := LookupServiceWithConsul("user-service")
		url := fmt.Sprintf("%s%s%s", serviceInfo, "/user-service/auth/check/", name)

		response, err := http.Get(url)
		if err != nil {
			return c.String(http.StatusForbidden, err.Error())
		} else {
			defer response.Body.Close()

			contents, err := ioutil.ReadAll(response.Body)
			if err != nil {
				return c.String(http.StatusForbidden, err.Error())
			}

			if string(contents) == "ok" {
				return c.String(http.StatusOK, "Product list")
			}
		}

		return c.String(http.StatusForbidden, "Your account must authenicate first")
	})

	RegisterService(SERVER_NAME, 3001)

	server.Logger.Fatal(server.Start(":3001"))
}

func DeRegisterServiceWithConsul(serviceID string)  {
	config := consulapi.DefaultConfig()
	consul, err := consulapi.NewClient(config)
	if err != nil {
		log.Print(err)
	}

	consul.Agent().ServiceDeregister(serviceID)
}

func RegisterServiceWithConsul(address string, port int, serviceID string, serviceName string) error {
	config := consulapi.DefaultConfig()
	consul, err := consulapi.NewClient(config)
	if err != nil {
		return err
	}

	registration := new(consulapi.AgentServiceRegistration)

	registration.ID = serviceID
	registration.Name = serviceName
	registration.Address = address
	registration.Tags = []string{"pro"}
	registration.Port = port
	registration.Check = new(consulapi.AgentServiceCheck)
	registration.Check.HTTP = fmt.Sprintf("http://%s:%v/product-service/healthcheck",
		address, port)
	registration.Check.Interval = "30s"
	registration.Check.Timeout = "10s"

	err = consul.Agent().ServiceRegister(registration)
	if err != nil {
		return err
	}

	return nil
}

func LookupServiceWithConsul(serviceID string) (string, error) {
	config := consulapi.DefaultConfig()
	consul, err := consulapi.NewClient(config)
	if err != nil {
		return "", err
	}
	services, err := consul.Agent().Services()
	if err != nil {
		return "", err
	}
	srvc := services[serviceID]
	address := srvc.Address
	port := srvc.Port
	// 1.1.1.1:3000/api/v1/fhsdhfjshfjsdhfsf
	return fmt.Sprintf("http://%s:%v", address, port), nil
}

func RegisterService(serviceName string, port int) {
	productService := new(Node)

	productService.Address = hostname()
	productService.Port = port
	productService.serviceID = serviceName
	productService.serviceName = serviceName

	RegisterServiceWithConsul(productService.Address, productService.Port,
		productService.serviceName, productService.serviceID)
}

func port(portKey string) string {
	p := os.Getenv(portKey)
	return fmt.Sprintf(":%s", p)
}

func hostname() string {
	hn, err := os.Hostname()
	if err != nil {
		log.Fatalln(err)
	}
	return hn
}
