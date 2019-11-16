package dns

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRemoveEnclosingDots(t *testing.T) {
	/*
		Given either or both a leading and/or trailing dot, we expect the RemoveEnclosingDots() function will
		remove the leading/trailing dots and return the sanitized string.  However, we do not expect that more than
		one leading/trailing dot will be sanitized (as this would most likely mask an error condition for the caller).
	*/
	assert.Equal(t, RemoveEnclosingDots(".f"), "f", "leading dot should have been removed")
	assert.Equal(t, RemoveEnclosingDots("f."), "f", "trailing dot should have been removed")
	assert.Equal(t, RemoveEnclosingDots(".f."), "f", "leading and trailing dots should have been removed")
	assert.NotEqual(t, RemoveEnclosingDots("..f.."), "f", "leading and trailing dots should have been removed")
}

func TestNewRecordWithEnclosingDots(t *testing.T) {
	/*
		Given a domain and zone, we expect that New Record will append the zone to the end of the domain, stripping
		any single leading/trailing dots, and the result would be a Fully Qualified Domain Name (FQDN) which is
		acceptable to the NS1 API.
	*/
	domain := "my-service.my-subdomain"
	zone := "my-domain.tld"
	expectedFqdn := fmt.Sprintf("%s.%s", domain, zone)

	badDomain := fmt.Sprintf(".%s.", domain)
	badZone := fmt.Sprintf(".%s.", zone)

	GoodRecord := NewRecord(zone, domain, "CNAME")
	BadRecord := NewRecord(badZone, badDomain, "CNAME")

	assert.Equal(t, GoodRecord.Zone , BadRecord.Zone, "zone not properly sanitized (enclosing dots)")
	assert.Equal(t, GoodRecord.Domain, BadRecord.Domain, "domain not properly sanitized (enclosing dots)")
	assert.Equal(t, GoodRecord.Domain, expectedFqdn, "GoodRecord.Domain does not match expectedFqdn")
	assert.Equal(t, BadRecord.Domain, expectedFqdn, "BadRecord.Domain does not match expectedFqdn")
}
