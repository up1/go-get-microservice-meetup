package main

import (
	"encoding/json"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"golang.org/x/net/context"
)

// Service
type Counter interface {
	Add(v int) int
}

type counterService struct {
	v int
}

func (c *counterService) Add(v int) int {
	c.v += v
	return c.v
}

// Endpoint
type addRequest struct {
	V int `json:"value"`
}

type addResponse struct {
	V int `json:"value"`
}

func makeAddEndpoint(sv Counter) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(addRequest)
		v := sv.Add(req.V)
		return addResponse{v}, nil
	}
}

func decodeAddRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req addRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return req, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

func main() {
	c := &counterService{}

	var sep endpoint.Endpoint
	sep = makeAddEndpoint(c)

	addHandler := httptransport.NewServer(
		sep,
		decodeAddRequest,
		encodeResponse,
	)

	http.Handle("/add", addHandler)
	http.ListenAndServe(":8080", nil)
}
