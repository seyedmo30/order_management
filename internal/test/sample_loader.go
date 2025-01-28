package test

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/seyedmo30/order_management/internal/config"
	"github.com/seyedmo30/order_management/internal/dto"
	"github.com/seyedmo30/order_management/pkg"
)

type SampleData struct {
	IsFilled bool
	Config   config.App `json:"config" yaml:"config"`
	// add more fields for your use case here
	RepositoryData RepositoryData `json:"repository_data" yaml:"repository_data"`
}

type RepositoryData struct {
	CreatOrderRepositoryRequest          dto.CreatOrderRepositoryRequest
	GetOrderByIDRepositoryResponse       string
	UpdateOrderByIDRepositoryRequest     dto.UpdateOrderByIDRepositoryRequest
	LockOrderOptimisticRepositoryRequest dto.LockOrderOptimisticRepositoryRequest
}

var sampleData SampleData

// loadSampleData loads testing configuration and data from the config file
func loadSampleData() error {
	logger := pkg.GetLogger()
	filename := "sample.data.json"
	logger.Info("loading sample data")
	logger.Info(os.Getwd())
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		logger.Error("sample data file not found")
		return err
	}
	jsonFile, err := os.ReadFile(filename)
	if err != nil {
		logger.Error("failed to load sample data", err)
		return err
	}
	err = json.Unmarshal(jsonFile, &sampleData)
	if err != nil {
		logger.Error("failed to unmarshal sample data", err)
		return err
	}
	return nil
}

func GetSampleData() (SampleData, error) {
	if sampleData.IsFilled {
		return sampleData, nil
	}
	err := loadSampleData()
	if err != nil {
		return SampleData{}, err
	}
	sampleData.IsFilled = true
	return sampleData, nil
}

// Set environment variables from struct fields
func SetEnvFromStruct(s interface{}) error {
	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return fmt.Errorf("expected a struct, got %v", v.Kind())
	}

	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		envTag := field.Tag.Get("env")
		if envTag == "" {
			envTag = field.Name
		}

		envKey := strings.Split(envTag, ",")[0]

		var envValue string
		switch value.Kind() {
		case reflect.String:
			envValue = value.String()
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			envValue = strconv.FormatInt(value.Int(), 10)
		default:
			continue
		}
		if envValue == "" || envValue == "0" {
			continue
		}
		if err := os.Setenv(envKey, envValue); err != nil {
			return fmt.Errorf("failed to set environment variable %s: %v", envKey, err)
		}
	}
	return nil
}
