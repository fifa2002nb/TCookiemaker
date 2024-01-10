package spider

import (
	"bytes"
	//log "github.com/Sirupsen/logrus"
	"github.com/elazarl/goproxy"
	"net/http"
)

type Processor interface {
	ProcessDetail(resp *http.Response, ctx *goproxy.ProxyCtx) ([]byte, error)
	Output()
}

type BaseProcessor struct {
	req          *http.Request
	data         []byte
	detailResult *DetailResult
}

type DetailResult struct {
	Cookie string
	Url    string
	Data   []byte
}

func NewBaseProcessor() *BaseProcessor {
	return &BaseProcessor{}
}

func (p *BaseProcessor) ProcessDetail(resp *http.Response, ctx *goproxy.ProxyCtx) (data []byte, err error) {
	p.req = ctx.Req
	var buf bytes.Buffer
	if _, err = buf.ReadFrom(resp.Body); err != nil {
		return
	}
	if err = resp.Body.Close(); err != nil {
		return
	}
	data = buf.Bytes()
	p.detailResult = &DetailResult{Cookie: p.req.Header.Get("Cookie"), Url: p.req.URL.String(), Data: data}
	return
}

func (p *BaseProcessor) DetailResult() *DetailResult {
	return p.detailResult
}

func (p *BaseProcessor) GetRequest() *http.Request {
	return p.req
}

func (p *BaseProcessor) Output() {
}
