package metrics

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

/*
 添加 prometheus 相关
*/

const metricsHTTPPort = 8889 // for prometheus to connect

var (
	// ClientConnected tcp-service 瞬时连接数
	ClientConnected prometheus.Gauge
	// ReqRecvTotal tcp-service 接收消息计数
	ReqRecvTotal prometheus.Counter
	// RspSendTotal tcp-service 发送消息计数
	RspSendTotal prometheus.Counter
)

func init() {
	ReqRecvTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "tcp_server_req_recv_total",
	})
	RspSendTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "tcp_server_rsp_send_total",
	})

	ClientConnected = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "tcp_server_client_connected",
	})

	prometheus.MustRegister(ReqRecvTotal, RspSendTotal, ClientConnected)

	// start the metrics server
	metricsServer := &http.Server{
		Addr: fmt.Sprintf(":%d", metricsHTTPPort),
	}

	mu := http.NewServeMux()
	mu.Handle("/metrics", promhttp.Handler())
	metricsServer.Handler = mu

	go func() {
		err := metricsServer.ListenAndServe()
		if err != nil {
			fmt.Println("prometheus-exporter http server start failed:", err)
		}
	}()
	fmt.Println("metrics server start ok(*:8889)")
}
