package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"tenant/pkg/derrors"
)

// Decode reads the body of an HTTP request looking for a JSON document. The
// body is decoded into the provided value.
func Decode(r *http.Request, val interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(val); err != nil {
		return ErrBadRequest(err, "Bad Request: "+err.Error())
	}
	return nil
}

func ParsePagination(r *http.Request) (uint64, uint64, error) {
	page := r.URL.Query().Get("page")
	if page == "" {
		page = "1"
	}
	limit := r.URL.Query().Get("limit")
	if limit == "" {
		limit = "10"
	}
	pageInt, err := strconv.ParseUint(page, 10, 64)
	if err != nil {
		return 0, 0, derrors.WrapStack(err, derrors.InvalidArgument, "limit not valid")
	}
	limitInt, err := strconv.ParseUint(limit, 10, 64)
	if err != nil {
		return 0, 0, derrors.WrapStack(err, derrors.InvalidArgument, "limit not valid")
	}
	if limitInt > 100 {
		limitInt = 100
	}
	return pageInt, limitInt, nil
}

func ParsePaginationWithMaxLimit(r *http.Request, maxLimit uint64) (uint64, uint64, error) {
	page := r.URL.Query().Get("page")
	if page == "" {
		page = "1"
	}
	limit := r.URL.Query().Get("limit")
	if limit == "" {
		limit = "10"
	}
	pageInt, err := strconv.ParseUint(page, 10, 64)
	if err != nil {
		return 0, 0, derrors.WrapStack(err, derrors.InvalidArgument, "limit not valid")
	}
	limitInt, err := strconv.ParseUint(limit, 10, 64)
	if err != nil {
		return 0, 0, derrors.WrapStack(err, derrors.InvalidArgument, "limit not valid")
	}
	if limitInt > maxLimit {
		limitInt = maxLimit
	}
	return pageInt, limitInt, nil
}

var boolValues = []string{
	"1",
	"t",
	"T",
	"TRUE",
	"true",
	"True",
	"0",
	"f",
	"F",
	"FALSE",
	"false",
	"False",
}

func ParseBoolQueryParam(r *http.Request, key string) (*bool, error) {
	param := r.URL.Query().Get(key)
	if param == "" {
		return nil, nil
	}

	boolValue, err := strconv.ParseBool(param)
	if err != nil {
		return nil, derrors.New(derrors.InvalidArgument, key+" is invalid")
	}

	return &boolValue, nil
}
