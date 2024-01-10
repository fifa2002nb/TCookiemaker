package spider

import (
	log "github.com/Sirupsen/logrus"
	"github.com/elazarl/goproxy"
	"net/http"
	"strings"
)

type spider struct {
	proxy *goproxy.ProxyHttpServer
}

var _spider = newSpider()

func Regist(proc Processor) {
	_spider.Regist(proc)
}

func OnReq(f func(ctx *goproxy.ProxyCtx) (*http.Request, *http.Response)) {
	_spider.OnReq(f)
}

func Run(port string) {
	_spider.Run(port)
}

func newSpider() *spider {
	sp := &spider{}
	sp.proxy = goproxy.NewProxyHttpServer()
	sp.proxy.OnRequest().HandleConnect(goproxy.AlwaysMitm)
	return sp
}

func Header() http.Header {
	return header
}

func (s *spider) Regist(proc Processor) {
	s.proxy.OnResponse().DoFunc(ProxyHandle(proc))
}

func (s *spider) OnReq(f func(ctx *goproxy.ProxyCtx) (*http.Request, *http.Response)) {
	handler := func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		return f(ctx)
	}
	s.proxy.OnRequest().DoFunc(handler)
}

func (s *spider) Run(port string) {
	if rootConfig.Compress {
		s.OnReq(func(ctx *goproxy.ProxyCtx) (req *http.Request, resp *http.Response) {
			host := ctx.Req.URL.Host
			req = ctx.Req
			if !strings.Contains(host, "taobao") {
				resp = goproxy.NewResponse(req, "text/html", http.StatusNotFound, "")
			}
			return
		})
	}
	log.Infof("server will at port:%v", port)
	log.Fatal(http.ListenAndServe(":"+port, s.proxy))
}

var (
	header     http.Header //全局缓存wechat header
	rootConfig = &Config{
		Verbose:  false,
		Compress: true,
	}
)

type Config struct {
	Verbose  bool // Debug
	Compress bool // To ingore other request just to save the net cost
}

func InitConfig(conf *Config) {
	rootConfig = conf
}

func orPanic(err error) {
	if err != nil {
		panic(err)
	}
}
