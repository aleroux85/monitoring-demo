package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type streams struct {
	Streams []stream `json:"streams"`
}

type stream struct {
	Stream map[string]string `json:"stream"`
	Values []entry           `json:"values"`
}

type entry [2]string

type Logx struct {
	job string
}

func (l Logx) Log(lvl, msg string) error {
	ts := fmt.Sprint(time.Now().UnixNano())
	s := map[string]string{l.job: lvl}
	v := []entry{entry{ts, msg}}

	ss := streams{
		Streams: []stream{
			stream{
				Stream: s,
				Values: v,
			},
		},
	}

	b, err := json.Marshal(ss)
	if err != nil {
		return err
	}

	rsp, err := http.Post("http://loki:3100/loki/api/v1/push", "application/json", bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	defer rsp.Body.Close()

	if rsp.Status != "204 No Content" {
		return errors.New("got status " + rsp.Status + "\n" + string(b))
	}

	return nil
}

func main() {
	reqDuration := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "company_a",
		Subsystem: "service_a",
		Name:      "request_duration",
		Help:      "Duration of requests.",
	})

	go func() {
		logx := Logx{"service_a"}

		for {
			reqDuration.Set(rand.NormFloat64()*10.0 + 400.0)
			if err := logx.Log("INFO", "info message here"); err != nil {
				fmt.Println(err)
			}
			time.Sleep(time.Second)
		}
	}()

	if err := prometheus.Register(reqDuration); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("reqDuration registered.")
	}
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":2112", nil)
}
