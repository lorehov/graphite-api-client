package graphite

import (
	"github.com/buger/jsonparser"
	"net/url"
	"strconv"
	"time"
)

// RenderRequest is struct, describing request to graphite `/render/` api.
// No fields are required. If field has zero value it'll be just skipped in request.
// RenderRequest.Targets are slice of strings, were every entry is a path identifying one or several metrics,
// optionally with functions acting on those metrics.
//
// Warning. While wildcards could be used in Targets one should use them with caution, as
// using of the simple target like "main.cluster.*.cpu.*" could result in hundreds of series
// with megabytes of data inside.
type RenderRequest struct {
	From          time.Time
	Until         time.Time
	MaxDataPoints int
	Targets       []string
}

func (r RenderRequest) toQueryString() string {
	values := url.Values{
		"format": []string{"json"},
		"target": r.Targets,
	}
	if !r.From.IsZero() {
		values.Set("from", strconv.FormatInt(r.From.Unix(), 10))
	}
	if !r.Until.IsZero() {
		values.Set("until", strconv.FormatInt(r.Until.Unix(), 10))
	}
	if r.MaxDataPoints != 0 {
		values.Set("maxDataPoints", strconv.Itoa(r.MaxDataPoints))
	}
	qs := values.Encode()
	return "/render/?" + qs
}

// QueryRender performs query to graphite `/render/` api. Normally it should return `[]graphite.Series`,
// but if things go wrong it will return `graphite.RequestError` error.
func (c *Client) QueryRender(r RenderRequest) ([]Series, error) {
	empty := []Series{}
	data, err := c.makeRequest(r)
	if err != nil {
		return empty, err
	}

	metrics, err := unmarshallSeries(data)
	if err != nil {
		return empty, c.createError(r, "Can't unmarshall response")
	}
	return metrics, nil
}

// Series describes time series data for given target.
type Series struct {
	Target     string
	Datapoints []DataPoint
}

// DataPoint describes concrete point of time series.
type DataPoint struct {
	Value     float64
	Timestamp time.Time
}

func unmarshallSeries(data []byte) ([]Series, error) {
	empty, result := []Series{}, []Series{}
	if len(data) == 0 {
		return empty, nil
	}
	var ie error = nil
	_, err := jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
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

	_, err = jsonparser.ArrayEach(rawData, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
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
	_, err = jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
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
		position++
	})
	if err != nil {
		return empty, err
	}
	return result, nil
}
