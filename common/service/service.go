package service

import (
	"errors"
	"log"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/bytom/btmcpool/common/logger"
)

type Service struct {
	name string
	ip   string
	cfg  *Config
	eng  *gin.Engine
}

func acquireIp() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			// interface down
			continue
		}
		if iface.Flags&net.FlagLoopback != 0 {
			// loopback interface
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			return ip.String(), nil
		}
	}
	return "", errors.New("no iface available")
}

func New(name string, cfg *Config) *Service {
	if cfg.mode == modeProd {
		gin.SetMode(gin.ReleaseMode)
	}
	e := gin.New()
	e.Use(gin.Recovery())

	// create end point for health check
	e.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "UP"})
	})

	// disable logging on health check end point
	e.Use(gin.LoggerWithWriter(gin.DefaultWriter, "/health"))

	// logger
	if err := logger.InitWithFields(cfg.log.level, map[string]interface{}{"service": name + "." + string(cfg.mode)}); err != nil {
		log.Fatalf("error init logger. error: %v", err)
	}

	ip, err := acquireIp()
	if err != nil {
		logger.Fatal("error acquiring local ip", "error", err)
	}
	if len(ip) == 0 {
		ip = "unknown"
	}

	return &Service{
		name: name,
		ip:   ip,
		cfg:  cfg,
		eng:  e,
	}
}

func (s *Service) Get(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	return s.eng.GET(relativePath, handlers...)
}

func (s *Service) Post(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	return s.eng.POST(relativePath, handlers...)
}

// Use add middleware to the group
func (s *Service) Use(middleware ...gin.HandlerFunc) gin.IRoutes {
	return s.eng.Use(middleware...)
}

// Group sets group url path
func (s *Service) Group(path string) gin.IRoutes {
	return s.eng.Group(path)
}

func (s *Service) Run(addr ...string) {
	logger.Info("service start running", "name", s.name)

	s.eng.Run(addr...)
}
