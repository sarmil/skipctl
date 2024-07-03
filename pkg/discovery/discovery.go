package discovery

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net"

	"github.com/kartverket/skipctl/pkg/constants"
)

var resolver = net.DefaultResolver

type APIServer struct {
	Name string `json:"name"`
	Addr string `json:"addr"`
}

// DiscoverAPIServers will try to do a TXT lookup for a given DNS name. If found it will
// attempt a unmarshaL(base64_decode(TXT_RECORD_VALUE)) (pseudo code) into a list of APIServer
// structs.
func DiscoverAPIServers(dnsKey string) ([]APIServer, error) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.DNSDiscoverTimeout)
	defer cancel()

	records, err := resolver.LookupTXT(ctx, dnsKey)
	if err != nil {
		return nil, fmt.Errorf("failed discover available API servers: %w", err)
	}

	if len(records) > 1 {
		return nil, fmt.Errorf("found more than one TXT record with the same name: %s", dnsKey)
	}

	// Wrap the outer layer
	decodedBytes, err := base64.StdEncoding.DecodeString(records[0])
	if err != nil {
		return nil, fmt.Errorf("failed base64 decoding TXT record: %w", err)
	}

	// Decode JSON into a usable structure
	var apiServers []APIServer
	err = json.Unmarshal(decodedBytes, &apiServers)
	if err != nil {
		return nil, fmt.Errorf("failed unmarshalling TXT record: %w", err)
	}

	return apiServers, nil
}
