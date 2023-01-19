## 2.7.4 (January 17, 2023)
FEATURES:

* Adds support for getting DNS `networks`

## 2.7.3 (January 6, 2023)
BUG FIXES:

* some of the ErrXXXXMissing error conditions were not being detected
properly due to changes in the API responses.

## 2.7.2 (December 13, 2022)
BUG FIXES:

* When an HTTP error response with a non-JSON body is received, use the body as the error message (instead of a message about being unable to parse JSON).


## 2.7.1 (December 5, 2022)
BUG FIXES:

* Prevent a divide-by-zero error when the expected rate-limit headers
in the API response are corrupt or missing.

## 2.7.0 (October 25, 2022)
NOTES:

Previously, when updating monitoring jobs, if the `Rules` field
in the request was empty, it was silently removed from the API
call, making it effectively a no-op. As of version 2.7.0, specifying
an empty `Rules` field on a monitoring job update will delete all
rules from the job. Before upgrading to this version, users are advised
to check their monitoring job code for use of the `Rules` field to make
sure it will not accidentally delete rules.

FEATURES:

* Adds support for DNS views. See new section in README.md for details.
* Adds support for `OtherNetworks` secondary zone attribute.
* Adds support for deleting monitoring job rules.

BUG FIXES:

* Adjusted license file to be recognizable by automated parsers.
* Fixed handling of "Zone already exists" error.
* QPS function fixes.

## 2.6.5 (March 30, 2022)
FEATURES:

* Added override_ttl for records
* Fixed case sensitivity on records

## 2.6.4 (March 25, 2022)
FEATURES:

* Adds support for TSIG

## 2.6.3 (August 31, 2021)
FEATURES:

* Adds support for Pulsar applications and jobs

## 2.6.2 (July 23, 2021)
FEATURES:

* Adds missing `mute` field to monitoring jobs

## 2.6.1 (July 12, 2021)
FEATURES:

* Adds missing `records_allow` and `records_deny` fields to account permissions

## 2.6.0 (July 6, 2021)
FEATURES:

* Adds support for subdivision metadata

## 2.5.1 (May 20, 2021)
BUG FIXES:

* Handle rate limit despite HTTP error returned from API [#127](https://github.com/ns1/ns1-go/pull/127)

## 2.5.0 (April 5, 2021)
FEATURES:

* Adds support for `tags_allow` and `tags_deny` fields in IPAM and DHCP DDI permissions.

## 2.4.5 (March 31, 2021)
FEATURES:

* Adds support for DS record type

## 2.4.4
FEATURES:

* Added missing `manage_ip_whitelist` permission to the account permissions section

## 2.4.3
FEATURES:

* Added support for tagging in DDI on records and zones

## 2.4.2
FEATURES:

* Added support for "cost" meta option on records

## 2.4.1
BUGFIX:

* Fix incorrect tag

## 2.4.0
FEATURES:

* Exported mock API service - @iamgnat

## 2.3.0 (March 19, 2020)
FEATURES:

* Support for pulsar metadata on answers

## 2.2.1 (Febuary 14, 2020)
BUG FIXES

* Resolves a panic and several type issues with the ASN field.

## 2.2.0 (January 15, 2020)
FEATURES:

* Added support for various fields on team, user, and API key resources.

## 2.1.0 (January 9, 2020)
FEATURES:

* Added support for DDI permissions

## 2.0.0 (January 9, 2020)

* Begin semantic versioning
