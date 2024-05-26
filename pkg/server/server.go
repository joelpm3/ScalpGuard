package server

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
	"github.com/gin-gonic/gin"
	"github.com/h3adex/fp"
	log "github.com/sirupsen/logrus"
	"github.com/joelpm3/ScalpGuard/pkg/models"
	"github.com/joelpm3/ScalpGuard/pkg/router"
	"github.com/joelpm3/ScalpGuard/pkg/config"
	"golang.org/x/sync/errgroup"
)

type Config struct {
	Host    string `env:"HOST" envDefault:"0.0.0.0"`
	Port    int    `env:"PORT" envDefault:"80"`
	TlsPort int    `env:"TLS_PORT" envDefault:"443"`
}

type Server struct {
	Config       *Config
	ServerConfig  *config.ServerConfig
    BrowserConfig *config.BrowserConfig
	RoutingTable *router.RoutingTable
}

func matchesAnyRegex(patterns []string, userAgent string) bool {
	for _, pattern := range patterns {
		matched, err := regexp.MatchString(pattern, userAgent)
		if err != nil {
			log.WithError(err).Error("Invalid regex pattern")
			continue
		}
		if matched {
			return true
		}
	}
	return false
}

func New(config *Config, serverConfig *config.ServerConfig, browserConfig *config.BrowserConfig) *Server {
	routingTable := router.NewRoutingTable()

	for _, backend := range serverConfig.Backends {
        routingTable.AddBackend(backend.Host, fmt.Sprintf("http://127.0.0.1:%d", backend.Port))
    }

	return &Server{
		Config:       config,
		ServerConfig:  serverConfig,
        BrowserConfig: browserConfig,
		RoutingTable: routingTable,
	}
}

func (s *Server) Run(ctx context.Context) error {
	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		return s.startHTTPSServer(ctx)
	})

	return eg.Wait()
}

func (s *Server) startHTTPSServer(ctx context.Context) error {
	log.Infof("Starting HTTPS Server on port %d", s.Config.TlsPort)
	handler := s.setupRouter("https")

	err := fp.Server(
		ctx,
		handler,
		fp.Option{
			Addr: fmt.Sprintf("%s:%d", s.Config.Host, s.Config.TlsPort),
		},
	)

	if err != nil {
		log.WithError(err).Error("Error on HTTPS server")
		return err
	}

	return nil
}

func (s *Server) setupRouter(protocol string) *gin.Engine {
	router := gin.Default()
	router.Any("/*path", s.proxyToBackend)

	return router
}

func (s *Server) proxyToBackend(ctx *gin.Context) {
	parsedClientHello, err := s.parseClientHello(ctx)
	if err != nil {
		log.WithError(err).Error("Failed to parse client hello")
		ctx.String(http.StatusInternalServerError, "Internal Server Error")
		return
	}

	log.WithFields(log.Fields{
		"ja3": parsedClientHello.Ja3,
		"ja3H": parsedClientHello.Ja3H,
	}).Info("Parsed client hello")

	
	if !s.isRequestAllowed(ctx, parsedClientHello) {
        ctx.String(http.StatusForbidden, "Your request does not meet our security policies.")
        return
    }

	svcURL, err := s.RoutingTable.GetBackend(ctx.Request.Host, ctx.Request.RequestURI)
	if err != nil {
		log.WithError(err).Error("Routing error")
		ctx.String(http.StatusInternalServerError, "Internal Server Error")
		return
	}

	proxyURL := &url.URL{
		Host:     svcURL.Host,
		Scheme:   svcURL.Scheme,
		Path:     ctx.Request.URL.Path,
		RawQuery: ctx.Request.URL.RawQuery,
	}

	proxy := httputil.NewSingleHostReverseProxy(proxyURL)
	proxy.ServeHTTP(ctx.Writer, ctx.Request)
}

func (s *Server) isRequestAllowed(ctx *gin.Context, parsedClientHello models.ParsedClientHello) bool {
    // Retrieve the corresponding backend configuration based on the request
    // This is a placeholder; you need to implement the logic to get the actual backend config
    // For example, you might match the backend based on ctx.Request.Host or another criterion
    backendConfig, err := s.getBackendConfig(ctx.Request.Host)
    if err != nil {
        log.WithError(err).Error("Failed to get backend config")
        return false
    }

	whitelistApplied := false

	if len(backendConfig.Whitelist.UserAgents) > 0 || len(backendConfig.Whitelist.TLSFingerprintsJA3) > 0 || len(backendConfig.Whitelist.TLSFingerprintsJA3NE) > 0 || len(backendConfig.Whitelist.HTTP2Fingerprints) > 0 || len(backendConfig.Whitelist.BrowserSettings) > 0 {
		whitelistApplied = true
	}

	
	if whitelistApplied {
		
		if len(backendConfig.Whitelist.TLSFingerprintsJA3) > 0 {
			if contains(backendConfig.Whitelist.TLSFingerprintsJA3, parsedClientHello.Ja3H) {
				return true
			}
		}

		if len(backendConfig.Whitelist.TLSFingerprintsJA3NE) > 0 {
			if contains(backendConfig.Whitelist.TLSFingerprintsJA3NE, parsedClientHello.Ja3nH) {
				return true
			}
		}

		if len(backendConfig.Whitelist.HTTP2Fingerprints) > 0 {
			if contains(backendConfig.Whitelist.HTTP2Fingerprints, parsedClientHello.Http2FingerprintHash) {
				return true
			}
		}

		if len(backendConfig.Whitelist.UserAgents) > 0 {
			if matchesAnyRegex(backendConfig.Whitelist.UserAgents, ctx.Request.UserAgent()) {
				return true
			}
		}

		if len(backendConfig.Whitelist.BrowserSettings) > 0 {
			allowed := false
			for _, browserName := range backendConfig.Whitelist.BrowserSettings {
				if s.browserIdentified(browserName, parsedClientHello) {
					allowed = true
					break
				}
			}
			if allowed {
				return true
			}
		}
		return false
	} else {
		if len(backendConfig.Blacklist.TLSFingerprintsJA3) > 0 {
			if contains(backendConfig.Blacklist.TLSFingerprintsJA3, parsedClientHello.Ja3H) {
				return false
			}
		}

		if len(backendConfig.Blacklist.TLSFingerprintsJA3NE) > 0 {
			if contains(backendConfig.Blacklist.TLSFingerprintsJA3NE, parsedClientHello.Ja3nH) {
				return false
			}
		}

		if len(backendConfig.Blacklist.HTTP2Fingerprints) > 0 {
			if contains(backendConfig.Blacklist.HTTP2Fingerprints, parsedClientHello.Http2FingerprintHash) {
				return false
			}
		}

		if len(backendConfig.Blacklist.UserAgents) > 0 {
			if matchesAnyRegex(backendConfig.Blacklist.UserAgents, ctx.Request.UserAgent()) {
				return false
			}
		}

		if len(backendConfig.Blacklist.BrowserSettings) > 0 {
			identified := false
			for _, browserName := range backendConfig.Blacklist.BrowserSettings {
				if s.browserIdentified(browserName, parsedClientHello) {
					identified = true
					break
				}
			}
			if identified {
				return false
			}
		}
		return true

	}
}

func (s *Server) browserIdentified(browserName string, parsedClientHello models.ParsedClientHello) bool {
	fmt.Println(browserName)
    for _, browser := range s.BrowserConfig.Browsers {
        if browser.Name == browserName {
            // Check if any of the browser's criteria matches the parsedClientHello
            if (contains(browser.TLSFingerprintsJA3, parsedClientHello.Ja3H) ||
                contains(browser.TLSFingerprintsJA3NE, parsedClientHello.Ja3nH)) &&
                contains(browser.HTTP2Fingerprints, parsedClientHello.Http2FingerprintHash) && matchesAnyRegex(browser.UserAgents, parsedClientHello.UserAgent){
                return true
            }
        }
    }
    return false
}

func contains(slice []string, value string) bool {
    for _, item := range slice {
        if item == value {
            return true
        }
    }
    return false
}

func (s *Server) getBackendConfig(host string) (*config.BackendConfig, error) {
    for _, backend := range s.ServerConfig.Backends {
        // Assuming backend.Host is the unique identifier for each backend.
        // This comparison might need to be adjusted depending on how you're
        // identifying backends (e.g., including port number if necessary).
        if backend.Host == host {
            return &backend, nil
        }
    }
    return nil, fmt.Errorf("no backend configuration found for host: %s", host)
}

func (s *Server) parseClientHello(ctx *gin.Context) (models.ParsedClientHello, error) {
	parsedClientHello, err := models.ParseClientHello(ctx)
	if err != nil {
		return models.ParsedClientHello{}, fmt.Errorf("unable to parse client hello: %w", err)
	}

	addTLSFingerprintHeaders(ctx, parsedClientHello)

	return parsedClientHello, nil
}

func addTLSFingerprintHeaders(ctx *gin.Context, clientHello models.ParsedClientHello) {
	ctx.Request.Header.Add("X-Ja3-Fingerprint", clientHello.Ja3)
	ctx.Request.Header.Add("X-Ja3-Fingerprint-Hash", clientHello.Ja3H)
	ctx.Request.Header.Add("X-Ja3n-Fingerprint", clientHello.Ja3n)
	ctx.Request.Header.Add("X-Ja3n-Fingerprint-Hash", clientHello.Ja3nH)
	ctx.Request.Header.Add("X-HTTP2-Fingerprint", clientHello.Http2Fingerprint)
	ctx.Request.Header.Add("X-HTTP2-Fingerprint-Hash", clientHello.Http2FingerprintHash)
}
