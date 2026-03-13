package utils

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"xiaomiao-home-system/internal/conf"

	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

// 解析nacos地址
func parseNacosAddr(addr string) (method, host string, port uint64, err error) {
	// 检查是否有协议
	re := regexp.MustCompile(`^(https?)://`)
	matches := re.FindStringSubmatch(addr)
	if len(matches) == 0 {
		// 没有协议，默认 http
		addr = "http://" + addr
		method = "http"
	} else {
		method = matches[1]
	}

	u, err := url.Parse(addr)
	if err != nil {
		return "", "", 0, fmt.Errorf("invalid nacos addr: %v", err)
	}

	host = u.Hostname()
	portStr := u.Port()

	// 端口处理
	if portStr == "" {
		if method == "http" {
			port = 80
		} else if method == "https" {
			port = 443
		}
	} else {
		port, err = strconv.ParseUint(portStr, 10, 64)
		if err != nil {
			return "", "", 0, fmt.Errorf("invalid nacos port: %v", err)
		}
	}

	return method, host, port, nil
}

// validNacosConfig 校验nacos配置
func validNacosConfig(c *conf.Register) (*conf.Register, error) {
	if c.Nacos == nil {
		return nil, fmt.Errorf("nacos config is invalid")
	}

	if c.Nacos.Group == "" {
		c.Nacos.Group = "DEFAULT_GROUP"
	}

	if c.Nacos.NamespaceId == "" {
		return nil, fmt.Errorf("nacos namespace id is invalid")
	}

	if c.Nacos.DataId == "" {
		return nil, fmt.Errorf("nacos data id is invalid")
	}

	if c.Nacos.Username == "" && c.Nacos.Password == "" {
		return nil, fmt.Errorf("nacos username or password is invalid")
	}

	if c.Nacos.Addr == "" {
		return nil, fmt.Errorf("nacos addr is invalid")
	}

	if c.Nacos.LogLevel == "" {
		c.Nacos.LogLevel = "info"
	}

	if c.Nacos.LogDir == "" {
		c.Nacos.LogDir = filepath.Join(filepath.Join(os.TempDir(), "nacos"), "log")
	}

	if c.Nacos.CacheDir == "" {
		c.Nacos.CacheDir = filepath.Join(filepath.Join(os.TempDir(), "nacos"), "cache")
	}

	return c, nil
}

// NewNacosConfigClient 初始化nacos配置中心客户端
func NewNacosConfigClient(c *conf.Register, name string) config_client.IConfigClient {
	scheme, host, port, err := parseNacosAddr(c.Nacos.Addr)
	if err != nil {
		panic(err)
	}

	sc := []constant.ServerConfig{
		*constant.NewServerConfig(host, port, constant.WithScheme(scheme)),
	}

	cc := &constant.ClientConfig{
		AppName:     name,
		Username:    c.Nacos.Username,
		Password:    c.Nacos.Password,
		NamespaceId: c.Nacos.NamespaceId,
		LogDir:      c.Nacos.LogDir,
		CacheDir:    c.Nacos.CacheDir,
		LogLevel:    c.Nacos.LogLevel,
	}

	nc, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  cc,
			ServerConfigs: sc,
		},
	)

	if err != nil {
		panic(err)
	}

	return nc
}

// NewNacosNamingClient 初始化nacos注册中心客户端
func NewNacosNamingClient(c *conf.Register, name string) naming_client.INamingClient {
	scheme, host, port, err := parseNacosAddr(c.Nacos.Addr)
	if err != nil {
		panic(err)
	}

	sc := []constant.ServerConfig{
		*constant.NewServerConfig(host, port, constant.WithScheme(scheme)),
	}

	cc := &constant.ClientConfig{
		AppName:     name,
		Username:    c.Nacos.Username,
		Password:    c.Nacos.Password,
		NamespaceId: c.Nacos.NamespaceId,
		LogDir:      c.Nacos.LogDir,
		CacheDir:    c.Nacos.CacheDir,
		LogLevel:    c.Nacos.LogLevel,
	}

	nc, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  cc,
			ServerConfigs: sc,
		},
	)

	if err != nil {
		panic(err)
	}

	return nc
}

// GetNacosConfigFromEnv 从环境变量中读取nacos配置
func GetNacosConfigFromEnv() (*conf.Register, error) {
	c := &conf.Register{
		Nacos: &conf.Nacos{
			Addr:        os.Getenv("NACOS_URL"),
			NamespaceId: os.Getenv("NACOS_NAMESPACE_ID"),
			Group:       os.Getenv("NACOS_GROUP"),
			DataId:      os.Getenv("NACOS_DATA_ID"),
			Username:    os.Getenv("NACOS_USERNAME"),
			Password:    os.Getenv("NACOS_PASSWORD"),
			LogLevel:    os.Getenv("NACOS_LOG_LEVEL"),
		},
	}

	return validNacosConfig(c)
}
