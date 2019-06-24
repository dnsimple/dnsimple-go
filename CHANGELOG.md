# CHANGELOG

#### Release 0.30.0

- NEW: Added webhook event parser for dnssec.create, dnssec.delete

- CHANGE: Redesigned webhook event parsing to avoid event/data conflicts (dnsimple/dnsimple-go#85)

IMPORTANT: This release introduce breaking changes compared to the previous one,
as the webhook even parsing has been significantly reworked.


#### Release 0.23.0

- NEW: Added WHOIS privacy renewal (dnsimple/dnsimple-go#78)


#### Release 0.22.0

- CHANGED: Cleaned up webhook tests and added coverage for more events.


#### Release 0.21.0

- NEW: Added zone distribution and zone record distribution (dnsimple/dnsimple-go#64)


#### Release 0.20.0

- CHANGED: Renamed `Event_Header` to `EventHeader` as it's more go-style. The Event interface has been updated accordingly.

- CHANGED: Removed custom code for getting OAuth token. We now use RoundTripper for authentication (and pass an http.Client to create a new Client) (dnsimple/dnsimple-go#15, dnsimple/dnsimple-go#69).


#### Release 0.16.0

- NEW: Added Let's Encrypt certificate methods (dnsimple/dnsimple-go#63)

- REMOVED: Removed premium_price attribute from registrar order responses (dnsimple/dnsimple-go#67). Please do not rely on that attribute, as it returned an incorrect value. The attribute is going to be removed, and the API now returns a null value.


#### Release 0.15.0

- NEW: Added support for the DNSSEC Beta (GH-58)

- CHANGED: Changed response types to not be exported (GH-54)
- CHANGED: Updated registrar URLs (GH-59)

- FIXED: Unable to filter zone records by type (GH-65)


#### Release 0.14.0

- NEW: Added support for Collaborators API (GH-48)
- NEW: Added support for ZoneRecord regions (GH-47)
- NEW: Added support for Domain Pushes API (GH-42)
- NEW: Added support for domains premium prices API (GH-53)

- CHANGED: Renamed `DomainTransferRequest.AuthInfo` to `AuthCode` (GH-46)
- CHANGED: Updated registration, transfer, renewal response payload (dnsimple/dnsimple-developer#111, GH-52).
- CHANGED: Normalize unique string identifiers to SID (dnsimple/dnsimple-developer#113)
- CHANGED: Update whois privacy setting for domain (dnsimple/dnsimple-developer#120)


#### Release 0.13.0

- NEW: Added support for Accounts API (GH-29)
- NEW: Added support for Services API (GH-30, GH-35)
- NEW: Added support for Certificates API (GH-31)
- NEW: Added support for Vanity name servers API (GH-34)
- NEW: Added support for delegation API (GH-32)
- NEW: Added support for Templates API (GH-36, GH-39)
- NEW: Added support for Template Records API (GH-37)
- NEW: Added support for Zone files API (GH-38)


#### Release 0.12.0

- CHANGED: Setting a custom user-agent no longer overrides the origina user-agent (GH-26)
- CHANGED: Renamed Contact#email_address to Contact#email (GH-27)


#### Release 0.11.0

- NEW: Added support for parsing ZoneRecord webhooks.
- NEW: Added support for listing options (GH-25).
- NEW: Added support for Template API (GH-21).


#### Release 0.10.0

Initial release.
