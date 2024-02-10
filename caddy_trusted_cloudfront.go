package caddy_trusted_cloudfront

import (
	"context"
	"encoding/json"
	"net/http"
	"net/netip"
	"sync"
	"time"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
)

const (
	FETCH_URL = "https://d7uri8nf7uskq.cloudfront.net/tools/list-cloudfront-ips"
)

func init() {
	caddy.RegisterModule(CaddyTrustedCloudFront{})
}

// The module that auto trusted_proxies `AWS CloudFront EDGE servers` from CloudFront.
// Doc: https://docs.aws.amazon.com/AmazonCloudFront/latest/DeveloperGuide/LocationsOfEdgeServers.html
// Range from: https://d7uri8nf7uskq.cloudfront.net/tools/list-cloudfront-ips
type CaddyTrustedCloudFront struct {
	// Interval to update the trusted proxies list. default: 1d
	Interval caddy.Duration `json:"interval,omitempty"`
	ranges   []netip.Prefix
	ctx      caddy.Context
	lock     *sync.RWMutex
}

func (CaddyTrustedCloudFront) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.ip_sources.cloudfront",
		New: func() caddy.Module { return new(CaddyTrustedCloudFront) },
	}
}

func (s *CaddyTrustedCloudFront) Provision(ctx caddy.Context) error {
	s.ctx = ctx
	s.lock = new(sync.RWMutex)

	// update cron
	go func() {
		if s.Interval == 0 {
			s.Interval = caddy.Duration(24 * time.Hour) // default to 24 hours
		}
		ticker := time.NewTicker(time.Duration(s.Interval))
		s.lock.Lock()
		s.ranges, _ = s.fetchPrefixes()
		s.lock.Unlock()
		for {
			select {
			case <-ticker.C:
				prefixes, err := s.fetchPrefixes()
				if err != nil {
					break
				}
				s.lock.Lock()
				s.ranges = prefixes
				s.lock.Unlock()
			case <-s.ctx.Done():
				ticker.Stop()
				return
			}
		}
	}()
	return nil
}

type CloudFrontIPSource struct {
	CloudFrontGlobalIPList       []string `json:"CLOUDFRONT_GLOBAL_IP_LIST"`
	CloudFrontRegionalEdgeIPList []string `json:"CLOUDFRONT_REGIONAL_EDGE_IP_LIST"`
}

func (s *CaddyTrustedCloudFront) fetchPrefixes() ([]netip.Prefix, error) {
	ctx, cancel := context.WithTimeout(s.ctx, time.Duration(time.Minute))
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, FETCH_URL, nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var data CloudFrontIPSource
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}
	var prefixes []netip.Prefix
	for _, v := range data.CloudFrontGlobalIPList {
		prefix, err := caddyhttp.CIDRExpressionToPrefix(v)
		if err != nil {
			return nil, err
		}
		prefixes = append(prefixes, prefix)
	}
	for _, v := range data.CloudFrontRegionalEdgeIPList {
		prefix, err := caddyhttp.CIDRExpressionToPrefix(v)
		if err != nil {
			return nil, err
		}
		prefixes = append(prefixes, prefix)
	}
	return prefixes, nil
}

func (s *CaddyTrustedCloudFront) GetIPRanges(_ *http.Request) []netip.Prefix {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.ranges
}

// UnmarshalCaddyfile implements caddyfile.Unmarshaler. Syntax:
//
//	cloudfront {
//	   interval <duration>
//	}
func (m *CaddyTrustedCloudFront) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	d.Next() // consume directive name

	if d.NextArg() {
		return d.ArgErr()
	}

	for nesting := d.Nesting(); d.NextBlock(nesting); {
		switch d.Val() {
		case "interval":
			if !d.NextArg() {
				return d.ArgErr()
			}
			val, err := caddy.ParseDuration(d.Val())
			if err != nil {
				return err
			}
			m.Interval = caddy.Duration(val)
		default:
			return d.ArgErr()
		}
	}

	return nil
}

// Interface guards
var (
	_ caddy.Module            = (*CaddyTrustedCloudFront)(nil)
	_ caddy.Provisioner       = (*CaddyTrustedCloudFront)(nil)
	_ caddyfile.Unmarshaler   = (*CaddyTrustedCloudFront)(nil)
	_ caddyhttp.IPRangeSource = (*CaddyTrustedCloudFront)(nil)
)
