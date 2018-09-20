package firesocket

import (
	"time"
	"io/ioutil"
	"net"
	"io"
	"github.com/Greyh4t/dnscache"
)

type FireSocket struct {
	Options *Options
	Conn    net.Conn
}

var resolver *dnscache.Resolver

func (f *FireSocket) Connect(network, host, port string) error {
	if resolver != nil {
		var err error
		host, err = resolver.FetchOneString(host)
		if err != nil {
			return err
		}
	}

	conn, err := net.DialTimeout(network, net.JoinHostPort(host, port), f.Options.Timeout)
	if err != nil {
		return err
	}
	f.Conn = conn
	return nil
}

func (f *FireSocket) Write(b []byte) (int, error) {
	if f.Options.WriteTimeout > 0 {
		f.Conn.SetWriteDeadline(time.Now().Add(f.Options.WriteTimeout))
	}
	return f.Conn.Write(b)
}

// 阻塞，读取所有的数据，直到EOF或者超时
func (f *FireSocket) Read() ([]byte, error) {
	if f.Options.ReadTimeout > 0 {
		f.Conn.SetReadDeadline(time.Now().Add(f.Options.ReadTimeout))
	}
	return ioutil.ReadAll(f.Conn)
}

// 不阻塞，读到的数据有可能小于给定的n值
func (f *FireSocket) ReadN(n int) ([]byte, error) {
	if f.Options.ReadTimeout > 0 {
		f.Conn.SetReadDeadline(time.Now().Add(f.Options.ReadTimeout))
	}
	buf := make([]byte, n)
	c, err := f.Conn.Read(buf)
	return buf[:c], err
}

// 读取指定长度的字符串
func (f *FireSocket) ReadN2String(n int) (string, error) {
	buf, err := f.ReadN(n)
	return string(buf), err
}

// 阻塞，直到读到给定的数据长度或者超时
func (f *FireSocket) ReadAtLeast(n int64) ([]byte, error) {
	if f.Options.ReadTimeout > 0 {
		f.Conn.SetReadDeadline(time.Now().Add(f.Options.ReadTimeout))
	}
	return ioutil.ReadAll(io.LimitReader(f.Conn, n))
}

func (f *FireSocket) Close() error {
	return f.Conn.Close()
}

type Options struct {
	// DNS缓存有效期
	DNSCacheExpire time.Duration

	// 链接超时间
	Timeout time.Duration

	// 读取超时时间
	ReadTimeout time.Duration

	// 写入超时时间
	WriteTimeout time.Duration
}

func New(options *Options) *FireSocket {

	if options.DNSCacheExpire > 0 {
		resolver = dnscache.New(options.DNSCacheExpire)
	}

	return &FireSocket{
		Options: options,
	}
}
