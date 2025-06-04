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

func NewRepository(logger *zap.SugaredLogger) domain.CommentRepository {
	return repository{
		logger: logger,
	}
}

func (r repository) Create(comment domain.Comment) error {
	serviceName := "vk-golang-comment-create"
	startTime := time.Now()
	var statusCode int

	reqBody, err := json.Marshal(comment)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, "http://vk-golang.ru:16000/comment", bytes.NewBuffer(reqBody))
	if err != nil {
		r.logger.Errorf("service %s failed to create request: %v", serviceName, err)
		middleware.RecordExternalCallMetrics(serviceName, req.URL.String(), 0, 0, err)
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
		r.logger.Errorf("service %s with url %s got unexpected status code: %d", serviceName, req.URL, statusCode)
		return errors.New("failed to create comment remotely")
	}

	return nil
}

func (r repository) Like(commentID string) error {
	serviceName := "vk-golang-comment-like" // Имя внешнего сервиса для метки
	startTime := time.Now()
	var statusCode int

	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("http://vk-golang.ru:16000/comment/like?cid=%s", commentID),
		nil,
	)
	if err != nil {
		return err
	}

	r.logger.Infof("service %s with url %s created request", serviceName, req.URL)

	resp, callErr := http.DefaultClient.Do(req)

	latency := time.Since(startTime)

	if callErr != nil {
		middleware.RecordExternalCallMetrics(serviceName, req.URL.String(), latency, 0, callErr)
		return callErr
	}

	defer utils.BodyCloserWithSugaredLogger(resp.Body, r.logger)

	statusCode = resp.StatusCode
	middleware.RecordExternalCallMetrics(serviceName, req.URL.String(), latency, statusCode, nil)
	r.logger.Infof("service %s with url %s got response with status code %d", serviceName, req.URL, statusCode)

	if resp.StatusCode != http.StatusOK {
		return errors.New("failed to like comment remotely")
	}

	return nil
}
