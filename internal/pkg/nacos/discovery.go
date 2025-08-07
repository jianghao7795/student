package nacos

import (
	"fmt"
	"strconv"

	"student/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

type Discovery struct {
	client naming_client.INamingClient
	log    *log.Helper
}

func NewDiscovery(c *conf.Bootstrap, logger log.Logger) (*Discovery, error) {
	sc := []constant.ServerConfig{
		{
			IpAddr: c.Nacos.Discovery.Ip,
			Port:   uint64(c.Nacos.Discovery.Port),
		},
	}

	cc := constant.ClientConfig{
		NamespaceId:         c.Nacos.Discovery.NamespaceId,
		TimeoutMs:           uint64(c.Nacos.Config.TimeoutMs),
		NotLoadCacheAtStart: true,
		LogDir:              "/tmp/nacos/log",
		CacheDir:            "/tmp/nacos/cache",
		LogLevel:            c.Nacos.Config.LogLevel,
	}

	client, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create nacos naming client: %w", err)
	}

	return &Discovery{
		client: client,
		log:    log.NewHelper(logger),
	}, nil
}

// RegisterService 注册服务到Nacos
func (d *Discovery) RegisterService(serviceName, ip string, port int, metadata map[string]string) error {
	param := vo.RegisterInstanceParam{
		Ip:          ip,
		Port:        uint64(port),
		ServiceName: serviceName,
		Weight:      getWeight(metadata),
		ClusterName: metadata["cluster"],
		GroupName:   metadata["group"],
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
		Metadata:    metadata,
	}

	success, err := d.client.RegisterInstance(param)
	if err != nil {
		return fmt.Errorf("failed to register service: %w", err)
	}
	if !success {
		return fmt.Errorf("failed to register service: service registration failed")
	}

	d.log.Infof("Service %s registered successfully at %s:%d", serviceName, ip, port)
	return nil
}

// DeregisterService 从Nacos注销服务
func (d *Discovery) DeregisterService(serviceName, ip string, port int, metadata map[string]string) error {
	param := vo.DeregisterInstanceParam{
		Ip:          ip,
		Port:        uint64(port),
		ServiceName: serviceName,
		GroupName:   metadata["group"],
		Ephemeral:   true,
	}

	success, err := d.client.DeregisterInstance(param)
	if err != nil {
		return fmt.Errorf("failed to deregister service: %w", err)
	}
	if !success {
		return fmt.Errorf("failed to deregister service: service deregistration failed")
	}

	d.log.Infof("Service %s deregistered successfully", serviceName)
	return nil
}

// GetServiceInstances 获取服务实例列表
func (d *Discovery) GetServiceInstances(serviceName string) ([]ServiceInstance, error) {
	param := vo.GetServiceParam{
		ServiceName: serviceName,
		GroupName:   "DEFAULT_GROUP",
		Clusters:    []string{"DEFAULT"},
	}

	service, err := d.client.GetService(param)
	if err != nil {
		return nil, fmt.Errorf("failed to get service: %w", err)
	}

	instances := make([]ServiceInstance, 0, len(service.Hosts))
	for _, host := range service.Hosts {
		if host.Healthy {
			instance := ServiceInstance{
				ID:       host.InstanceId,
				Name:     serviceName,
				IP:       host.Ip,
				Port:     int(host.Port),
				Version:  host.Metadata["version"],
				Metadata: host.Metadata,
				Healthy:  host.Healthy,
			}
			instances = append(instances, instance)
		}
	}

	return instances, nil
}

// ServiceInstance 服务实例信息
type ServiceInstance struct {
	ID       string
	Name     string
	IP       string
	Port     int
	Version  string
	Metadata map[string]string
	Healthy  bool
}

// GetServiceURL 获取服务URL
func (si *ServiceInstance) GetServiceURL() string {
	return fmt.Sprintf("%s:%d", si.IP, si.Port)
}

// getWeight 获取权重，默认为10
func getWeight(metadata map[string]string) float64 {
	if weightStr, exists := metadata["weight"]; exists {
		if weight, err := strconv.ParseFloat(weightStr, 64); err == nil {
			return weight
		}
	}
	return 10.0
}
