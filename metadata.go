package powerdns

import (
	"context"
	"fmt"
	"net/http"
)

// MetadataService handles communication with the metadata-related methods of the Client API
type MetadataService service

// MetadataKind represents a string-valued metadata kind
type MetadataKind string

// MetadataKindPtr is a helper function that allocates a new MetadataKind value to store v and returns a pointer to it.
func MetadataKindPtr(v MetadataKind) *MetadataKind {
	return &v
}

const (
	// MetadataAllowAXFRFrom defines which IP ranges are allowed to perform AXFR requests
	MetadataAllowAXFRFrom MetadataKind = "ALLOW-AXFR-FROM"

	// MetadataTSIGAllowAXFR defines which TSIG keys are allowed to perform AXFR requests
	MetadataTSIGAllowAXFR MetadataKind = "TSIG-ALLOW-AXFR"

	// MetadataAXFRMasterTSIG defines which TSIG key to use when performing AXFR retrieval
	MetadataAXFRMasterTSIG MetadataKind = "AXFR-MASTER-TSIG"

	// MetadataSOAEdit defines the SOA-EDIT mode for the zone
	MetadataSOAEdit MetadataKind = "SOA-EDIT"

	// MetadataSOAEditAPI defines the SOA-EDIT mode for the zone when using the API
	MetadataSOAEditAPI MetadataKind = "SOA-EDIT-API"

	// MetadataNSEC3Param defines the NSEC3 parameters for the zone
	MetadataNSEC3Param MetadataKind = "NSEC3PARAM"

	// MetadataPresigned defines whether the zone is presigned
	MetadataPresigned MetadataKind = "PRESIGNED"

	// MetadataLuaAXFRScript defines a Lua script to be used for AXFR
	MetadataLuaAXFRScript MetadataKind = "LUA-AXFR-SCRIPT"

	// MetadataAPIRectify defines whether the zone should be rectified when using the API
	MetadataAPIRectify MetadataKind = "API-RECTIFY"

	// MetadataPublishCDNSKey defines whether CDNSKEY records should be published
	MetadataPublishCDNSKey MetadataKind = "PUBLISH-CDNSKEY"

	// MetadataPublishCDS defines whether CDS records should be published
	MetadataPublishCDS MetadataKind = "PUBLISH-CDS"

	// MetadataSlaveRenotify defines whether slaves should be renotified
	MetadataSlaveRenotify MetadataKind = "SLAVE-RENOTIFY"

	// MetadataAXFRSource defines the source address to use when performing AXFR retrieval
	MetadataAXFRSource MetadataKind = "AXFR-SOURCE"

	// MetadataNotifyDNSUpdate defines whether to notify after a DNS update
	MetadataNotifyDNSUpdate MetadataKind = "NOTIFY-DNSUPDATE"

	// MetadataAlsoNotify defines additional IP addresses to notify
	MetadataAlsoNotify MetadataKind = "ALSO-NOTIFY"

	// MetadataForwardDNSUpdate defines whether to forward DNS updates
	MetadataForwardDNSUpdate MetadataKind = "FORWARD-DNSUPDATE"

	// MetadataAllowDNSUpdateFrom defines which IP ranges are allowed to perform DNS updates
	MetadataAllowDNSUpdateFrom MetadataKind = "ALLOW-DNSUPDATE-FROM"

	// MetadataTSIGAllowDNSUpdate defines which TSIG keys are allowed to perform DNS updates
	MetadataTSIGAllowDNSUpdate MetadataKind = "TSIG-ALLOW-DNSUPDATE"

	// MetadataIXFR defines whether IXFR is enabled for the zone
	MetadataIXFR MetadataKind = "IXFR"
)

// Metadata structure with JSON API metadata
type Metadata struct {
	Kind     *MetadataKind `json:"kind,omitempty"`
	Metadata []string      `json:"metadata,omitempty"`
}

// List retrieves all metadata for a zone
func (m *MetadataService) List(ctx context.Context, domain string) ([]Metadata, error) {
	req, err := m.client.newRequest(ctx, http.MethodGet, fmt.Sprintf("servers/%s/zones/%s/metadata", m.client.VHost, makeDomainCanonical(domain)), nil, nil)
	if err != nil {
		return nil, err
	}

	metadata := make([]Metadata, 0)
	_, err = m.client.do(req, &metadata)
	return metadata, err
}

// Create creates a new metadata entry for a zone
func (m *MetadataService) Create(ctx context.Context, domain string, kind MetadataKind, values []string) (*Metadata, error) {
	metadata := Metadata{
		Kind:     &kind,
		Metadata: values,
	}

	req, err := m.client.newRequest(ctx, http.MethodPost, fmt.Sprintf("servers/%s/zones/%s/metadata", m.client.VHost, makeDomainCanonical(domain)), nil, metadata)
	if err != nil {
		return nil, err
	}

	responseMetadata := new(Metadata)
	_, err = m.client.do(req, &responseMetadata)
	return responseMetadata, err
}

// Get retrieves a specific metadata kind for a zone
func (m *MetadataService) Get(ctx context.Context, domain string, kind MetadataKind) (*Metadata, error) {
	req, err := m.client.newRequest(ctx, http.MethodGet, fmt.Sprintf("servers/%s/zones/%s/metadata/%s", m.client.VHost, makeDomainCanonical(domain), kind), nil, nil)
	if err != nil {
		return nil, err
	}

	metadata := new(Metadata)
	_, err = m.client.do(req, &metadata)
	return metadata, err
}

// Set creates or modifies a metadata kind for a zone (existing entries for the zone with the same kind are removed)
func (m *MetadataService) Set(ctx context.Context, domain string, kind MetadataKind, values []string) (*Metadata, error) {
	metadata := Metadata{
		Kind:     &kind,
		Metadata: values,
	}

	req, err := m.client.newRequest(ctx, http.MethodPut, fmt.Sprintf("servers/%s/zones/%s/metadata/%s", m.client.VHost, makeDomainCanonical(domain), kind), nil, metadata)
	if err != nil {
		return nil, err
	}

	responseMetadata := new(Metadata)
	_, err = m.client.do(req, &responseMetadata)
	return responseMetadata, err
}

// Delete removes a metadata kind from a zone
func (m *MetadataService) Delete(ctx context.Context, domain string, kind MetadataKind) error {
	req, err := m.client.newRequest(ctx, http.MethodDelete, fmt.Sprintf("servers/%s/zones/%s/metadata/%s", m.client.VHost, makeDomainCanonical(domain), kind), nil, nil)
	if err != nil {
		return err
	}

	_, err = m.client.do(req, nil)
	return err
}
