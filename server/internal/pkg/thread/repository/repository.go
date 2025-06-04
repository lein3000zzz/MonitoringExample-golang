package repository

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"server/internal/api/middleware"
	"server/internal/pkg/domain"
	"server/internal/pkg/utils"
	"time"

	"go.uber.org/zap"
)

type repository struct {
	logger *zap.SugaredLogger
}

func (r repository) Create(thread domain.Thread) error {
	serviceName := "vk-golang-thread-create"
	startTime := time.Now()
	var statusCode int

	reqBody, err := json.Marshal(thread)
	if err != nil {
		// r.logger.Errorf("service %s failed to marshal thread: %v", serviceName, err)
		return err
	}

	req, err := http.NewRequest(http.MethodPost, "http://vk-golang.ru:15000/thread", bytes.NewBuffer(reqBody))
	if err != nil {
		r.logger.Errorf("service %s failed to create request: %v", serviceName, err)
		return err
	}
	r.logger.Infof("service %s with url %s created request", serviceName, req.URL)
	req.Header.Set("Content-Type", "application/json")

	resp, callErr := http.DefaultClient.Do(req)

	latency := time.Since(startTime)

	if callErr != nil {
		r.logger.Errorf("service %s with url %s failed to call external service: %v", serviceName, req.URL, callErr)
		middleware.RecordExternalCallMetrics(serviceName, req.URL.String(), latency, 0, callErr)
		return callErr
	}
	defer utils.BodyCloserWithSugaredLogger(resp.Body, r.logger)

	statusCode = resp.StatusCode
	middleware.RecordExternalCallMetrics(serviceName, req.URL.String(), latency, statusCode, nil) // callErr здесь nil

	r.logger.Infof("service %s with url %s got response with status code %d", serviceName, req.URL, statusCode)

	if resp.StatusCode != http.StatusOK {
		r.logger.Errorf("service %s with url %s failed to create thread: %d", serviceName, req.URL, resp.StatusCode)
		return errors.New("failed to create thread remotely")
	}

	return nil
}

func (r repository) Get(id string) (domain.Thread, error) {
	serviceName := "vk-golang-thread-get" // Имя внешнего сервиса для метки
	startTime := time.Now()
	var statusCode int

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://vk-golang.ru:15000/thread?id=%s", id), nil)
	if err != nil {
		return domain.Thread{}, err
	}
	r.logger.Infof("service %s with url %s created request", serviceName, req.URL)
	resp, callErr := http.DefaultClient.Do(req)

	latency := time.Since(startTime)

	if callErr != nil {
		middleware.RecordExternalCallMetrics(serviceName, req.URL.String(), latency, 0, callErr)
		return domain.Thread{}, err
	}
	defer utils.BodyCloserWithSugaredLogger(resp.Body, r.logger)

	statusCode = resp.StatusCode
	middleware.RecordExternalCallMetrics(serviceName, req.URL.String(), latency, statusCode, nil)
	r.logger.Infof("service %s with url %s got response with status code %d", serviceName, req.URL, statusCode)
	if resp.StatusCode != http.StatusOK {
		return domain.Thread{}, errors.New("failed to fetch thread remotely")
	}

	var thread domain.Thread
	err = json.NewDecoder(resp.Body).Decode(&thread)
	if err != nil {
		return domain.Thread{}, err
	}

	return thread, nil
}

func NewRepository(logger *zap.SugaredLogger) domain.ThreadRepository {
	return repository{
		logger: logger,
	}
}
