package handlers

import (
	"fmt"
	"net/url"
	"strconv"
	"time"
)

func getPaginationParams(query url.Values) (time.Time, int, error) {
	markerUnix, err := parseParam(query, ParamFrom, time.Now().Unix())
	if err != nil {
		return time.Time{}, 0, fmt.Errorf("invalid %s param", ParamFrom)
	}
	from := time.Unix(markerUnix, 0)

	limit, err := parseParam(query, ParamLimit, LimitDefault)
	if err != nil {
		return time.Time{}, 0, fmt.Errorf("invalid %s param", ParamLimit)
	}

	// Cap queryable items
	if limit > LimitMax {
		limit = LimitMax
	}

	return from, limit, nil
}

// parseParam is a generic function that parses typed parameters from a url.Values instance.
// The second value is non-nil on failure to parse when the key exists.
// If the key does not exist, it returns the def default value.
func parseParam[T int | int64 | time.Time | string](query url.Values, param string, def T) (T,
	error,
) {
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
	case *int64:
		*t, err = strconv.ParseInt(v, 10, 64)
	}

	return ret, err
}
