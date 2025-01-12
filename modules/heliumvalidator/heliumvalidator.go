package heliumvalidator

import (
	"net/http"
	"time"

	"github.com/netdata/go.d.plugin/agent/module"
	"github.com/netdata/go.d.plugin/pkg/web"
)

func init() {
	module.Register("heliumvalidator", module.Creator{
		Defaults: module.Defaults{
			UpdateEvery: 5,
		},
		Create: func() module.Module { return New() },
	})
}

type Config struct {
	web.HTTP `yaml:",inline"`
}

type Heliumvalidator struct {
	module.Base
	Config `yaml:",inline"`

	httpClient *http.Client
	charts     *module.Charts
}

func New() *Heliumvalidator {
	return &Heliumvalidator{
		Config: Config{
			HTTP: web.HTTP{
				Request: web.Request{
					URL: "http://127.0.0.1:4467",
				},
				Client: web.Client{
					Timeout: web.Duration{Duration: time.Second},
				},
			},
		},
	}
}

func (e *Heliumvalidator) Init() bool {
	err := e.validateConfig()
	if err != nil {
		e.Errorf("config validation: %v", err)
		return false
	}

	client, err := e.initHTTPClient()
	if err != nil {
		e.Errorf("init HTTP client: %v", err)
		return false
	}
	e.httpClient = client

	cs, err := e.initCharts()
	if err != nil {
		e.Errorf("init charts: %v", err)
		return false
	}
	e.charts = cs

	return true
}

func (e *Heliumvalidator) Check() bool {
	return len(e.Collect()) > 0
}

func (e *Heliumvalidator) Charts() *module.Charts {
	return e.charts
}

func (e *Heliumvalidator) Collect() map[string]int64 {
	ms, err := e.collect()
	if err != nil {
		e.Error(err)
	}

	if len(ms) == 0 {
		return nil
	}

	return ms
}

func (e *Heliumvalidator) Cleanup() {
	if e.httpClient == nil {
		return
	}
	e.httpClient.CloseIdleConnections()
}
