package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/DanyJDuque/go_lib_response/response"
	"github.com/DanyJDuque/gocourse_enrollment/internal/enrollment"
	"github.com/gorilla/mux"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
)

func NewEnrollmentHTTPServer(ctx context.Context, endpoints enrollment.Endpoints) http.Handler {

	r := mux.NewRouter()

	opts := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(encodeError),
	}

	r.Handle("/enrollments", httptransport.NewServer(
		endpoint.Endpoint(endpoints.Create),
		decodeCreateEnrollment, encodeResponse,
		opts...,
	)).Methods("POST")

	return r
}

// r.Handle("/enrollments", httptransport.NewServer(
// 	endpoint.Endpoint(endpoints.GetAll),
// 	decodeGetAllEnrollment, encodeResponse,
// 	opts...,
// )).Methods("GET")

// r.Handle("/enrollments/{id}", httptransport.NewServer(
// 	endpoint.Endpoint(endpoints.Get),
// 	decodeGetEnrollment, encodeResponse,
// 	opts...,
// )).Methods("GET")

// r.Handle("/enrollments/{id}", httptransport.NewServer(
// 	endpoint.Endpoint(endpoints.Update),
// 	decodeUpdateEnrollment, encodeResponse,
// 	opts...,
// )).Methods("PATCH")

// r.Handle("/enrollments/{id}", httptransport.NewServer(
// 	endpoint.Endpoint(endpoints.Delete),
// 	decodeDeleteEnrollment, encodeResponse,
// 	opts...,
// )).Methods("DELETE")

func decodeCreateEnrollment(_ context.Context, r *http.Request) (interface{}, error) {
	var req enrollment.CreateReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, response.BadRequest(fmt.Sprintf("invalid request format: '%v'", err.Error()))
	}
	return req, nil
}

// func decodeGetEnrollment(_ context.Context, r *http.Request) (interface{}, error) {

// 	p := mux.Vars(r)
// 	req := enrollment.GetReq{
// 		ID: p["id"],
// 	}
// 	return req, nil
// }

// func decodeGetAllEnrollment(_ context.Context, r *http.Request) (interface{}, error) {

// 	v := r.URL.Query()

// 	limit, _ := strconv.Atoi(v.Get("limit"))
// 	page, _ := strconv.Atoi(v.Get("page"))

// 	req := enrollment.GetAllReq{
// 		FirstName: v.Get("first_name"),
// 		LastName:  v.Get("last_name"),
// 		Limit:     limit,
// 		Page:      page,
// 	}
// 	return req, nil
// }

// func decodeUpdateEnrollment(_ context.Context, r *http.Request) (interface{}, error) {
// 	var req enrollment.UpdateReq

// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		return nil, response.BadRequest(fmt.Sprintf("invalid request format: '%v'", err.Error()))
// 	}

// 	path := mux.Vars(r)
// 	req.ID = path["id"]

// 	return req, nil
// }

// func decodeDeleteEnrollment(_ context.Context, r *http.Request) (interface{}, error) {

// 	path := mux.Vars(r)
// 	req := enrollment.DeleteReq{
// 		ID: path["id"],
// 	}

// 	return req, nil
// }

func encodeResponse(ctx context.Context, w http.ResponseWriter, resp interface{}) error {
	r := resp.(response.Response)
	w.Header().Set("Contet-Type", "application/json; charset=utf8")
	w.WriteHeader(r.StatusCode())
	return json.NewEncoder(w).Encode(r)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Contect-type", "application/json; charset=utf-8")
	resp := err.(response.Response)
	w.WriteHeader(resp.StatusCode())
	_ = json.NewEncoder(w).Encode(resp)
}
