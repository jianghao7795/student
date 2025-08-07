package main

import (
	"flag"
	"net"
	"os"
	"time"

	"student/internal/conf"
	"student/internal/pkg/nacos"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"

	_ "go.uber.org/automaxprocs"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name = "user-service"
	// Version is the version of the compiled software.
	Version = "v0.0.1"
	// flagconf is the config flag.
	flagconf string

	id, _ = os.Hostname()
)

func init() {
	var cstZone = time.FixedZone("CST", 8*3600)
	time.Local = cstZone
	flag.StringVar(&flagconf, "conf", "../../configs/user-service.yaml", "config path, eg: -conf config.yaml")
}

func newApp(logger log.Logger, gs *grpc.Server, hs *http.Server, discovery *nacos.Discovery, c *conf.Bootstrap) *kratos.App {
	// 获取本机IP
	ip := getLocalIP()
	
	// 解析HTTP端口
	httpPort := 8601
	if c.Server.Http.Addr != "" {
		if _, port, err := net.SplitHostPort(c.Server.Http.Addr); err == nil {
			if p, err := net.LookupPort("tcp", port); err == nil {
				httpPort = p
			}
		}
	}

	// 注册服务到Nacos
	metadata := map[string]string{
		"version": "1.0.0",
		"zone":    "zone1",
		"weight":  "10",
		"cluster": "DEFAULT",
		"group":   "DEFAULT_GROUP",
	}

	go func() {
		if err := discovery.RegisterService("user-service", ip, httpPort, metadata); err != nil {
			log.Error("Failed to register service to Nacos", "error", err)
		}
	}()

	return kratos.New(
		kratos.ID(id),
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(
			gs,
			hs,
		),
	)
}

func main() {
	flag.Parse()
	logger := log.With(log.NewStdLogger(os.Stdout),
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		"service.id", id,
		"service.name", Name,
		"service.version", Version,
		"trace.id", tracing.TraceID(),
		"span.id", tracing.SpanID(),
	)
	log.Warn("flagconf: ", flagconf)
	c := config.New(
		config.WithSource(
			file.NewSource(flagconf),
		),
	)
	defer c.Close()

	if err := c.Load(); err != nil {
		panic(err)
	}

	var bc conf.Bootstrap
	if err := c.Scan(&bc); err != nil {
		panic(err)
	}

	// 初始化Nacos服务发现
	discovery, err := nacos.NewDiscovery(&bc, logger)
	if err != nil {
		log.Error("Failed to create Nacos discovery", "error", err)
		panic(err)
	}

	app, cleanup, err := wireApp(&bc, logger, discovery)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// start and wait for stop signal
	if err := app.Run(); err != nil {
		panic(err)
	}
}

// getLocalIP 获取本机IP地址
func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "127.0.0.1"
	}
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return "127.0.0.1"
} 