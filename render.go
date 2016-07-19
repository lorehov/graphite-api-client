package graphite

import (
	"github.com/buger/jsonparser"
	"time"
	"strconv"
	"net/url"
)

type RenderRequest struct {
	From	time.Time
	Until	time.Time
	MaxDataPoints int
	Targets []string
}


func (g RenderRequest) ToQueryString() string {
	values := url.Values{
		"format": []string{"json"},
		"target": g.Targets,
	}
	if !g.From.IsZero() {
		values.Set("from", strconv.FormatInt(g.From.Unix(), 10))
	}
	if !g.Until.IsZero() {
		values.Set("until", strconv.FormatInt(g.Until.Unix(), 10))
	}
	if g.MaxDataPoints != 0 {
		values.Set("maxDataPoints", strconv.Itoa(g.MaxDataPoints))
	}
	qs := values.Encode()
	return "/render/?" + qs
}


type Series struct {
	Target string
	Datapoints []DataPoint
}


type DataPoint struct {
	Value float64
	Timestamp time.Time
}


func unmarshallMetrics(data []byte) ([]Series, error) {
	empty, result := []Series{}, []Series{}
	if len(data) == 0 {
		return empty, nil
	}
	var ie error = nil
	err := jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		if err != nil {
			return
		}

		datapoints, e := unmarshallDatapoints(value)
		if e != nil {
			ie = e
			return
		}

		target, e := jsonparser.GetString(value, "target")
		if e != nil {
			ie = e
			return
		}

		result = append(result, Series{Target: target, Datapoints: datapoints})
	})

	if err != nil {
		return empty, err
	}
	if ie != nil {
		return empty, ie
	}
	return result, nil
}


func unmarshallDatapoints(data []byte) ([]DataPoint, error) {
	empty, result := []DataPoint{}, []DataPoint{}
	rawData, _, _, err := jsonparser.Get(data, "datapoints")
	if err != nil {
		return empty, err
	}

	err = jsonparser.ArrayEach(rawData, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		if err != nil {
			return
		}
		datapoint, e := unmarshallDatapoint(value)
		if e != nil {
			err = e
			return
		}
		result = append(result, datapoint)
	})
	if err != nil {
		return empty, err
	}
	return result, nil
}


func unmarshallDatapoint(data []byte) (DataPoint, error) {
	empty, result := DataPoint{}, DataPoint{}
	var err error = nil
	position := 0
	err = jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		if err != nil {
			return
		}
		if position == 0 {
			if dataType == jsonparser.Null {
				result.Value = 0
			} else {
				v, e := strconv.ParseFloat(string(value), 64)
				if e != nil {
					err = e
					return
				}
				result.Value = v
			}
		} else {
			ts, e := strconv.ParseInt(string(value), 10, 32)
			if err != nil {
				err = e
				return
			}
			result.Timestamp = time.Unix(ts, 0)
		}
		position ++
	})
	if err != nil {
		return empty, err
	}
	return result, nil
}
