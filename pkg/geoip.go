package pkg

import (
	"fmt"
	"log"
	"net"

	"github.com/oschwald/geoip2-golang"
)

// GeoIPService provides country lookup functionality.
type GeoIPService struct {
	db *geoip2.Reader
}

// NewGeoIPService loads the GeoIP2 database and returns a new service.
func NewGeoIPService(dbPath string) (*GeoIPService, error) {
	db, err := geoip2.Open(dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open GeoIP database at %s: %w", dbPath, err)
	}
	log.Printf("Successfully loaded GeoIP database from %s\n", dbPath)
	return &GeoIPService{db: db},
		nil
}

// GetCountryCode looks up the country code for a given IP address.
func (s *GeoIPService) GetCountryCode(ipStr string) (string, error) {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return "", fmt.Errorf("invalid IP address: %s", ipStr)
	}

	record, err := s.db.Country(ip)
	if err != nil {
		return "", fmt.Errorf("failed to lookup country for IP %s: %w", ipStr, err)
	}

	return record.Country.IsoCode, nil
}

// Close closes the GeoIP2 database.
func (s *GeoIPService) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}