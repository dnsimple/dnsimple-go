# CHANGELOG

## main

## 1.4.0 (Unreleased)

FEATURES:

- NEW: Added `GetDomainTransferLock`, `EnableDomainTransferLock`, and `DisableDomainTransferLock` APIs to manage domain transfer locks. (dnsimple/dnsimple-go#147)

## 1.3.0

FEATURES:

- NEW: Added `ListRegistrantChanges`, `CreateRegistrantChange`, `CheckRegistrantChange`, `GetRegistrantChange`, and `DeleteRegistrantChange` APIs to manage registrant changes. (dnsimple/dnsimple-go#146)

## 1.2.1

FEATURES:

- NEW: Added `ActivateZoneDns` to activate DNS services (resolution) for a zone. (dnsimple/dnsimple-go#145)
- NEW: Added `DeactivateZoneDns` to deactivate DNS services (resolution) for a zone. (dnsimple/dnsimple-go#145)

IMPROVEMENTS:

- `EmailForward` `From` is deprecated. Please use `AliasName` instead for creating email forwards, and `AliasEmail` when retrieving email forwards. (dnsimple/dnsimple-go#145)
- `EmailForward` `To` is deprecated. Please use `DestinationEmail` instead for creating email forwards. (dnsimple/dnsimple-go#145)

## 1.2.0

- NEW: Support `GetDomainRegistration` and `GetDomainRenewal` APIs (dnsimple/dnsimple-go#132)

## 1.1.0

- NEW: Support `signature_algorithm` in the `LetsencryptCertificateAttributes` struct (dnsimple/dnsimple-go#128)

## 1.0.1

- CHANGED: Bump dependencies

## 1.0.0

- NEW: Expose AttributeErrors in ErrorResponse for getting detailed information about validation errors
- CHANGED: Support only last two golang versions: 1.18 and 1.19 according Golang Release Policy.
- CHANGED: Use testify/assert instead of stdlib

## 0.80.0

- CHANGED: Deprecate Certificate's `contact_id` (dnsimple/dnsimple-go#111)
- CHANGED: Dropped support for Go < 1.13

## 0.71.1

- FIXED: When purchasing a certificate the certificate id is populated now (CertificatePurchase)

## 0.71.0

- CHANGED: Updated Tld and DelegationSignerRecord types to support DS record key-data interface (dnsimple/dnsimple-go#107)

## 0.70.1

- DEPRECATED: Registrar.GetDomainPremiumPrice() has been deprecated, use Registrar.GetDomainPrices() instead.

## 0.70.0

- NEW: Added Registrar.GetDomainPrices() to retrieve whether a domain is premium and the prices to register, transfer, and renew. (dnsimple/dnsimple-go#103)

Incompatible changes:

- CHANGED: Domain.ExpiresOn has been replaced by Domain.ExpiresAt. (dnsimple/dnsimple-go#98)
- CHANGED: Certificate.ExpiresOn has been replaced by Certificate.ExpiresAt. (dnsimple/dnsimple-go#99)

## Release 0.63.0

- NEW: Added types and parsing for account membership events. (dnsimple/dnsimple-go#97)

## Release 0.62.0

- NEW: Added Registrar.GetDomainTransfer() to retrieve a domain transfer. (dnsimple/dnsimple-go#94)
- NEW: Added Registrar.CancelDomainTransfer() to cancel an in progress domain transfer. (dnsimple/dnsimple-go#94)

## Release 0.61.0

- NEW: Added convenient helpers to inizialize a client with common authentication strategies

## Release 0.60.0

- FIXED: A zone record can be updated without the risk of overriding the name by mistake (dnsimple/dnsimple-go#33, dnsimple/dnsimple-go#92)
- FIXED: Fixed a conflict where a Go zero-value would prevent sorting to work correctly (dnsimple/dnsimple-go#88, dnsimple/dnsimple-go#93)

Incompatible changes:

- CHANGED: CreateZoneRecord and UpdateZoneRecord now requires to use ZoneRecordAttributes instead of ZoneRecord. This is required to avoid conflicts caused by blank record names (dnsimple/dnsimple-go#92)
- CHANGED: ListOptions now use pointer values (dnsimple/dnsimple-go#93)

## Release 0.50.0

- NEW: Added Client.SetUserAgent() as a convenient helper to set a custom user agent.
- NEW: Added support for Registration/Transfer extended attributes (dnsimple/dnsimple-go#86)

Incompatible changes:

- NEW: Added support for context (dnsimple/dnsimple-go#82, dnsimple/dnsimple-go#90)
- CHANGED: Changed all method signatures so that the returned value is exported (dnsimple/dnsimple-go#91)
- CHANGED: Renamed the following structs to clarify intent:
  - DomainRegisterRequest -> RegisterDomainInput
  - DomainTransferRequest -> TransferDomainInput
  - DomainRenewRequest -> RenewDomainInput

## Release 0.40.0

Incompatible changes:

- CHANGED: Renamed ExchangeAuthorizationError.HttpResponse field to ExchangeAuthorizationError.HTTPResponse
- CHANGED: Renamed Response.HttpResponse field to Response.HTTPResponse
- REMOVED: Deleted deprecated ResetDomainToken method.

## Release 0.31.0

- CHANGED: User-agent format has been changed to prepend custom token before default token. (dnsimple/dnsimple-go#87)

## Release 0.30.0

- NEW: Added webhook event parser for dnssec.create, dnssec.delete
- CHANGE: Redesigned webhook event parsing to avoid event/data conflicts (dnsimple/dnsimple-go#85)

IMPORTANT: This release introduce breaking changes compared to the previous one,
as the webhook even parsing has been significantly reworked.

## Release 0.23.0

- NEW: Added WHOIS privacy renewal (dnsimple/dnsimple-go#78)

## Release 0.22.0

- CHANGED: Cleaned up webhook tests and added coverage for more events.

## Release 0.21.0

- NEW: Added zone distribution and zone record distribution (dnsimple/dnsimple-go#64)

## Release 0.20.0

- CHANGED: Renamed `Event_Header` to `EventHeader` as it's more go-style. The Event interface has been updated accordingly.
- CHANGED: Removed custom code for getting OAuth token. We now use RoundTripper for authentication (and pass an http.Client to create a new Client) (dnsimple/dnsimple-go#15, dnsimple/dnsimple-go#69).

## Release 0.16.0

- NEW: Added Let's Encrypt certificate methods (dnsimple/dnsimple-go#63)
- REMOVED: Removed premium_price attribute from registrar order responses (dnsimple/dnsimple-go#67). Please do not rely on that attribute, as it returned an incorrect value. The attribute is going to be removed, and the API now returns a null value.

## Release 0.15.0

- NEW: Added support for the DNSSEC Beta (GH-58)
- CHANGED: Changed response types to not be exported (GH-54)
- CHANGED: Updated registrar URLs (GH-59)
- FIXED: Unable to filter zone records by type (GH-65)

## Release 0.14.0

- NEW: Added support for Collaborators API (GH-48)
- NEW: Added support for ZoneRecord regions (GH-47)
- NEW: Added support for Domain Pushes API (GH-42)
- NEW: Added support for domains premium prices API (GH-53)
- CHANGED: Renamed `DomainTransferRequest.AuthInfo` to `AuthCode` (GH-46)
- CHANGED: Updated registration, transfer, renewal response payload (dnsimple/dnsimple-developer#111, GH-52).
- CHANGED: Normalize unique string identifiers to SID (dnsimple/dnsimple-developer#113)
- CHANGED: Update whois privacy setting for domain (dnsimple/dnsimple-developer#120)

## Release 0.13.0

- NEW: Added support for Accounts API (GH-29)
- NEW: Added support for Services API (GH-30, GH-35)
- NEW: Added support for Certificates API (GH-31)
- NEW: Added support for Vanity name servers API (GH-34)
- NEW: Added support for delegation API (GH-32)
- NEW: Added support for Templates API (GH-36, GH-39)
- NEW: Added support for Template Records API (GH-37)
- NEW: Added support for Zone files API (GH-38)

## Release 0.12.0

- CHANGED: Setting a custom user-agent no longer overrides the origina user-agent (GH-26)
- CHANGED: Renamed Contact#email_address to Contact#email (GH-27)

## Release 0.11.0

- NEW: Added support for parsing ZoneRecord webhooks.
- NEW: Added support for listing options (GH-25).
- NEW: Added support for Template API (GH-21).

## Release 0.10.0

Initial release.
