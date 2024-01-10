package cert

import (
	"crypto/tls"
	"crypto/x509"
	log "github.com/Sirupsen/logrus"
	"github.com/elazarl/goproxy"
)

func init() {
	err := SetCA(caCert, caKey)
	if nil != err {
		log.Errorf("%v", err)
	}
}

var caCert = []byte(`-----BEGIN CERTIFICATE-----
MIICSjCCAbWgAwIBAgIBADALBgkqhkiG9w0BAQUwSjEjMCEGA1UEChMaZ2l0aHVi
LmNvbS9lbGF6YXJsL2dvcHJveHkxIzAhBgNVBAMTGmdpdGh1Yi5jb20vZWxhemFy
bC9nb3Byb3h5MB4XDTAwMDEwMTAwMDAwMFoXDTQ5MTIzMTIzNTk1OVowSjEjMCEG
A1UEChMaZ2l0aHViLmNvbS9lbGF6YXJsL2dvcHJveHkxIzAhBgNVBAMTGmdpdGh1
Yi5jb20vZWxhemFybC9nb3Byb3h5MIGdMAsGCSqGSIb3DQEBAQOBjQAwgYkCgYEA
vz9BbCaJjxs73Tvcq3leP32hAGerQ1RgvlZ68Z4nZmoVHfl+2Nr/m0dmW+GdOfpT
cs/KzfJjYGr/84x524fiuR8GdZ0HOtXJzyF5seoWnbBIuyr1PbEpgRhGQMqqOUuj
YExeLbfNHPIoJ8XZ1Vzyv3YxjbmjWA+S/uOe9HWtDbMCAwEAAaNGMEQwDgYDVR0P
AQH/BAQDAgCkMBMGA1UdJQQMMAoGCCsGAQUFBwMBMA8GA1UdEwEB/wQFMAMBAf8w
DAYDVR0RBAUwA4IBKjALBgkqhkiG9w0BAQUDgYEAIcL8huSmGMompNujsvePTUnM
oEUKtX4Eh/+s+DSfV/TyI0I+3GiPpLplEgFWuoBIJGios0r1dKh5N0TGjxX/RmGm
qo7E4jjJuo8Gs5U8/fgThZmshax2lwLtbRNwhvUVr65GdahLsZz8I+hySLuatVvR
qHHq/FQORIiNyNpq/Hg=
-----END CERTIFICATE-----`)

var caKey = []byte(`-----BEGIN RSA PRIVATE KEY-----
MIICXQIBAAKBgQC/P0FsJomPGzvdO9yreV4/faEAZ6tDVGC+VnrxnidmahUd+X7Y
2v+bR2Zb4Z05+lNyz8rN8mNgav/zjHnbh+K5HwZ1nQc61cnPIXmx6hadsEi7KvU9
sSmBGEZAyqo5S6NgTF4tt80c8ignxdnVXPK/djGNuaNYD5L+4570da0NswIDAQAB
AoGBALzIv1b4D7ARTR3NOr6V9wArjiOtMjUrdLhO+9vIp9IEA8ZsA9gjDlCEwbkP
VDnoLjnWfraff5Os6+3JjHy1fYpUiCdnk2XA6iJSL1XWKQZPt3wOunxP4lalDgED
QTRReFbA/y/Z4kSfTXpVj68ytcvSRW/N7q5/qRtbN9804jpBAkEA0s6lvH2btSLA
mcEdwhs7zAslLbdld7rvfUeP82gPPk0S6yUqTNyikqshM9AwAktHY7WvYdKl+ghZ
HTxKVC4DoQJBAOg/IAW5RbXknP+Lf7AVtBgw3E+Yfa3mcdLySe8hjxxyZq825Zmu
Rt5Qj4Lw6ifSFNy4kiiSpE/ZCukYvUXGENMCQFkPxSWlS6tzSzuqQxBGwTSrYMG3
wb6b06JyIXcMd6Qym9OMmBpw/J5KfnSNeDr/4uFVWQtTG5xO+pdHaX+3EQECQQDl
qcbY4iX1gWVfr2tNjajSYz751yoxVbkpiT9joiQLVXYFvpu+JYEfRzsjmWl0h2Lq
AftG8/xYmaEYcMZ6wSrRAkBUwiom98/8wZVlB6qbwhU1EKDFANvICGSWMIhPx3v7
MJqTIj4uJhte2/uyVvZ6DC6noWYgy+kLgqG0S97tUEG8
-----END RSA PRIVATE KEY-----`)

func SetCA(caCert, caKey []byte) error {
	goproxyCa, err := tls.X509KeyPair(caCert, caKey)
	if err != nil {
		return err
	}
	if goproxyCa.Leaf, err = x509.ParseCertificate(goproxyCa.Certificate[0]); err != nil {
		return err
	}
	goproxy.GoproxyCa = goproxyCa
	goproxy.OkConnect = &goproxy.ConnectAction{Action: goproxy.ConnectAccept, TLSConfig: goproxy.TLSConfigFromCA(&goproxyCa)}
	goproxy.MitmConnect = &goproxy.ConnectAction{Action: goproxy.ConnectMitm, TLSConfig: goproxy.TLSConfigFromCA(&goproxyCa)}
	goproxy.HTTPMitmConnect = &goproxy.ConnectAction{Action: goproxy.ConnectHTTPMitm, TLSConfig: goproxy.TLSConfigFromCA(&goproxyCa)}
	goproxy.RejectConnect = &goproxy.ConnectAction{Action: goproxy.ConnectReject, TLSConfig: goproxy.TLSConfigFromCA(&goproxyCa)}
	return nil
}
