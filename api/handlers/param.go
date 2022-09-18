package handlers

import (
	"fmt"
	"net/url"
	"strconv"
	"time"
)

func getPaginationParams(query url.Values) (time.Time, int, error) {
	marker, err := parseParam(query, ParamMarker, time.Now())
	if err != nil {
		return time.Time{}, 0, fmt.Errorf("invalid %s param", ParamMarker)
	}

	limit, err := parseParam(query, ParamLimit, LimitDefault)
	if err != nil {
		return time.Time{}, 0, fmt.Errorf("invalid %s param", ParamLimit)
	}

	// Cap queryable items
	if limit > LimitMax {
		limit = LimitMax
	}

	return marker, limit, nil
}

func parseParam[T int | time.Time | string](query url.Values, param string, def T) (T, error) {
	if !query.Has(param) {
		return def, nil
	}

	v := query.Get(param)
	var ret T
	var err error
	switch t := any(&ret).(type) {
	case *string:
		*t, err = v, nil
	case *int:
		*t, err = strconv.Atoi(v)
	case *time.Time:
		*t, err = time.Parse(QueryTimeFmt, v)
	}

	return ret, err
}
