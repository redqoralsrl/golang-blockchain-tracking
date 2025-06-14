package jsonRpc

import (
	"blockchain-tracking/internal/logger"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type JsonRpcAdapter interface {
	CreateRequest(rpcUrl string, method string, params []interface{}) ([]byte, error)
	CreateRequestMultiple(rpcUrl string, payloads []Payload) ([]byte, error)
}

type Payload struct {
	Jsonrpc string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	ID      int           `json:"id"`
}

type JsonRpc struct {
	logger logger.Logger
}

var _ JsonRpcAdapter = (*JsonRpc)(nil)

func NewJsonRpc(l logger.Logger) *JsonRpc {
	return &JsonRpc{logger: l}
}

func (s *JsonRpc) CreateRequest(rpcUrl string, method string, params []interface{}) ([]byte, error) {
	payload := &Payload{
		Jsonrpc: "2.0",
		Method:  method,
		Params:  params,
		ID:      1,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		s.logger.Error("json marshal error", logger.Field{Key: "error", Value: err.Error()})
		return nil, err
	}
	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest("POST", rpcUrl, body)
	if err != nil {
		s.logger.Error("create http request error", logger.Field{Key: "error", Value: err.Error()})
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		s.logger.Error("client request error", logger.Field{Key: "error", Value: err.Error()})
		return nil, err
	}

	defer res.Body.Close()
	bytes, _ := io.ReadAll(res.Body)

	return bytes, nil
}

func (s *JsonRpc) CreateRequestMultiple(rpcUrl string, payloads []Payload) ([]byte, error) {
	payloadBytes, err := json.Marshal(payloads)
	if err != nil {
		s.logger.Error("json marshal error", logger.Field{Key: "error", Value: err.Error()})
		return nil, err
	}
	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest("POST", rpcUrl, body)
	if err != nil {
		s.logger.Error("create http request error", logger.Field{Key: "error", Value: err.Error()})
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		s.logger.Error("create http request error", logger.Field{Key: "error", Value: err.Error()})
		return nil, err

	}

	defer res.Body.Close()
	bytes, _ := io.ReadAll(res.Body)

	return bytes, nil
}
