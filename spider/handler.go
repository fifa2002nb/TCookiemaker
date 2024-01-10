package spider

import (
	"bytes"
	log "github.com/Sirupsen/logrus"
	"github.com/elazarl/goproxy"
	"io/ioutil"
	"net/http"
	"reflect"
)

func ProxyHandle(proc Processor) func(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
	return func(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
		if resp == nil || resp.StatusCode != 200 {
			return resp
		}
		if rootConfig.Verbose {
			log.Infof("Hijacked of %s:%v", ctx.Req.Method, ctx.Req.URL.String())
			log.Infof("cookie:%v", ctx.Req.Cookies())
		}
		return resp
	}
}

func handleDetail(resp *http.Response, ctx *goproxy.ProxyCtx, proc Processor) {
	var err error
	p := getProcessor(ctx.Req, proc)
	data, err := p.ProcessDetail(resp, ctx)
	if err != nil {
		log.Println(err.Error())
	}
	var buf = bytes.NewBuffer(data)
	go p.Output()
	resp.Body = ioutil.NopCloser(bytes.NewReader(buf.Bytes()))
}

// get the processor from cache
func getProcessor(req *http.Request, proc Processor) Processor {
	t := reflect.TypeOf(proc)
	v := reflect.New(t.Elem())
	p := v.Interface().(Processor)
	return p
}
