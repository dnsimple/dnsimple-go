# Changelog

This project uses [Semantic Versioning 2.0.0](http://semver.org/), the format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).

## 8.0.0 - 2026-01-22

### Added

- Support for subscription webhook events (`subscription.migrate`, `subscription.renew`, `subscription.renew:failed`, `subscription.subscribe`, `subscription.unsubscribe`)
- Support for `domain.restore` event

### Removed

- Removed deprecated `GetDomainPremiumPrice`. Use `GetDomainPrices` instead.
- Removed deprecated `GetDomainPremiumPrice`. Use `GetDomainPrices` instead.
- Removed deprecated `GetWhoisPrivacy` (dnsimple/dnsimple-developer#919)
- Removed deprecated `RenewWhoisPrivacy` (dnsimple/dnsimple-developer#919)

## 7.0.1 - 2025-10-22

### Added

- Support for "account.sso_user_add" event (GH-227)

## 7.0.0 - 2025-09-29

### Added

- Support for Batch Change Zone Records API (GH-221)

### Changed

- Dropped support for Go < 1.24
- Bump dependencies

## 6.0.1 - 2025-09-01

### Added

- Expanded DNSSEC event payload (GH-214)

## 6.0.0 - 2025-08-20

### Changed

- Remove from and to from EmailForwards (GH-212)
- Add support for email_forward.activate and deactivate webhook events (GH-207)
- Added support for Go 1.25
- Bumped dependencies

## 5.0.0 - 2025-05-09

### Fixed

- Missing restore_price field in DomainPrice struct (GH-192)
- Pass in options when making ListRegistrantChange API call (GH-195)
- Handle error responses for registerDomain endpoint action (GH-140)

### Changed

- Bump dependencies
- Dropped support for Go < 1.23
- Add golangci-lint (GH-190)

### Removed

- Deprecated `DomainCollaborators`. Please use our Domain Access Control feature.

## 4.0.0 - 2025-03-21

### Fixed

- Installation via Go Modules (GH-184)

### Changed

- Added support for Go >= 1.24
- Bump dependencies

### Removed

- Dropped support for installing via [GOPATH mode](https://go.dev/wiki/GOPATH#gopath-development-mode)

## 3.0.0 - 2024-12-12

### Added

- `active` attribute to EmailForward

### Changed

- Added support for Go >= 1.23
- Bump dependencies

### Deprecated

- `DomainCollaborators` have been deprecated and will be removed in the next major version. Please use our Domain Access Control feature.

## 2.0.0 - 2024-03-12

### Changed

- Bump dependencies
- Dropped support for Go < 1.20
- Added support for Go >= 1.22

## 1.7.0 - 2024-03-04

### Added

- Support for domain restore (GH-166)

## 1.6.0 - 2024-02-06

### Added

- `DnsAnalytics` to query and pull data from the DNS Analytics API endpoint (GH-164)

## 1.5.1 - 2023-11-22

### Added

- `Secondary`, `LastTransferredAt`, `Active` to `Zone`

## 1.5.0 - 2023-11-03

### Added

- `billing.ListCharges` to list charges for the account [learn more](https://developer.dnsimple.com/v2/billing-charges/). (GH-156)

## 1.4.1 - 2023-09-21

### Added

- Events for `domain.transfer_lock_enable` and `domain.transfer_lock_disable` (GH-149)

## 1.4.0 - 2023-08-31

### Added

- `GetDomainTransferLock`, `EnableDomainTransferLock`, and `DisableDomainTransferLock` APIs to manage domain transfer locks. (GH-147)

## 1.3.0 - 2023-08-24

### Added

- `ListRegistrantChanges`, `CreateRegistrantChange`, `CheckRegistrantChange`, `GetRegistrantChange`, and `DeleteRegistrantChange` APIs to manage registrant changes. (GH-146)

## 1.2.1 - 2023-08-11

### Added

- `ActivateZoneDns` to activate DNS services (resolution) for a zone. (GH-145)
- `DeactivateZoneDns` to deactivate DNS services (resolution) for a zone. (GH-145)

### Deprecated

- `EmailForward` `From` is deprecated. Please use `AliasName` instead for creating email forwards, and `AliasEmail` when retrieving email forwards. (GH-145)
- `EmailForward` `To` is deprecated. Please use `DestinationEmail` instead for creating email forwards. (GH-145)

## 1.2.0 - 2023-03-03

### Added

- Support `GetDomainRegistration` and `GetDomainRenewal` APIs (GH-132)

## 1.1.0 - 2023-02-28

### Added

- Support `signature_algorithm` in the `LetsencryptCertificateAttributes` struct (GH-128)

## 1.0.1 - 2023-02-20

### Changed

- Bump dependencies

## 1.0.0 - 2022-09-20

### Added

- Expose AttributeErrors in ErrorResponse for getting detailed information about validation errors

### Changed

- Support only last two golang versions: 1.18 and 1.19 according Golang Release Policy.
- Use testify/assert instead of stdlib

## 0.80.0 - 2022-09-06

### Changed

- Dropped support for Go < 1.13

### Deprecated

- Certificate's `contact_id` (GH-111)

## 0.71.1 - 2021-11-02

### Fixed

- When purchasing a certificate the certificate id is populated now (CertificatePurchase)

## 0.71.0 - 2021-10-19

### Changed

- Updated Tld and DelegationSignerRecord types to support DS record key-data interface (GH-107)

## 0.70.1 - 2021-05-19

### Deprecated

- Registrar.GetDomainPremiumPrice() has been deprecated, use Registrar.GetDomainPrices() instead.

## 0.70.0 - 2021-05-19

### Added

- Registrar.GetDomainPrices() to retrieve whether a domain is premium and the prices to register, transfer, and renew. (GH-103)

### Changed

- Domain.ExpiresOn has been replaced by Domain.ExpiresAt. (GH-98)
- Certificate.ExpiresOn has been replaced by Certificate.ExpiresAt. (GH-99)

## 0.63.0 - 2020-05-25

### Added

- Types and parsing for account membership events. (GH-97)

## 0.62.0 - 2020-05-13

### Added

- Registrar.GetDomainTransfer() to retrieve a domain transfer. (GH-94)
- Registrar.CancelDomainTransfer() to cancel an in progress domain transfer. (GH-94)

## 0.61.0 - 2020-05-12

### Added

- Convenient helpers to inizialize a client with common authentication strategies

## 0.60.0 - 2020-05-02

### Fixed

- A zone record can be updated without the risk of overriding the name by mistake (GH-33, GH-92)
- A conflict where a Go zero-value would prevent sorting to work correctly (GH-88, GH-93)

### Changed

- CreateZoneRecord and UpdateZoneRecord now requires to use ZoneRecordAttributes instead of ZoneRecord. This is required to avoid conflicts caused by blank record names (GH-92)
- ListOptions now use pointer values (GH-93)

## 0.50.0 - 2020-05-02

### Added

- Client.SetUserAgent() as a convenient helper to set a custom user agent.
- Support for Registration/Transfer extended attributes (GH-86)
- Support for context (GH-82, GH-90)

### Changed

- Changed all method signatures so that the returned value is exported (GH-91)
- Renamed the following structs to clarify intent:

  - DomainRegisterRequest -> RegisterDomainInput
  - DomainTransferRequest -> TransferDomainInput
  - DomainRenewRequest -> RenewDomainInput

## 0.40.0 - 2020-04-27

### Changed

- Renamed ExchangeAuthorizationError.HttpResponse field to ExchangeAuthorizationError.HTTPResponse
- Renamed Response.HttpResponse field to Response.HTTPResponse

### Removed

- Deleted deprecated ResetDomainToken method.

## 0.31.0 - 2020-02-11

### Changed

- User-agent format has been changed to prepend custom token before default token. (GH-87)

## 0.30.0 - 2019-06-24

### Added

- Webhook event parser for dnssec.create, dnssec.delete

### Changed

- Redesigned webhook event parsing to avoid event/data conflicts (GH-85)

## 0.23.0 - 2019-02-01

### Added

- WHOIS privacy renewal (GH-78)

## 0.22.0 - 2018-11-16

### Changed

- Cleaned up webhook tests and added coverage for more events.

## 0.21.0 - 2018-10-16

### Added

- Zone distribution and zone record distribution (GH-64)

## 0.20.0 - 2018-09-12

### Changed

- Renamed `Event_Header` to `EventHeader` as it's more go-style. The Event interface has been updated accordingly.
- Removed custom code for getting OAuth token. We now use RoundTripper for authentication (and pass an http.Client to create a new Client) (GH-15, GH-69).

## 0.16.0 - 2018-01-27

### Added

- Let's Encrypt certificate methods (GH-63)

### Removed

- premium_price attribute from registrar order responses (GH-67). Please do not rely on that attribute, as it returned an incorrect value. The attribute is going to be removed, and the API now returns a null value.

## 0.15.0 - 2017-12-11

### Added

- Support for the DNSSEC Beta (GH-58)

### Fixed

- Unable to filter zone records by type (GH-65)

### Changed

- Changed response types to not be exported (GH-54)
- Updated registrar URLs (GH-59)

## 0.14.0 - 2016-12-12

### Added

- Support for Collaborators API (GH-48)
- Support for ZoneRecord regions (GH-47)
- Support for Domain Pushes API (GH-42)
- Support for domains premium prices API (GH-53)

### Changed

- Renamed `DomainTransferRequest.AuthInfo` to `AuthCode` (GH-46)
- Updated registration, transfer, renewal response payload (dnsimple/dnsimple-developer#111, GH-52).
- Normalize unique string identifiers to SID (dnsimple/dnsimple-developer#113)
- Update whois privacy setting for domain (dnsimple/dnsimple-developer#120)

## 0.13.0 - 2016-07-31

### Added

- Support for Accounts API (GH-29)
- Support for Services API (GH-30, GH-35)
- Support for Certificates API (GH-31)
- Support for Vanity name servers API (GH-34)
- Support for delegation API (GH-32)
- Support for Templates API (GH-36, GH-39)
- Support for Template Records API (GH-37)
- Support for Zone files API (GH-38)

## 0.12.0 - 2016-06-22

### Changed

- Setting a custom user-agent no longer overrides the origina user-agent (GH-26)
- Renamed Contact#email_address to Contact#email (GH-27)

## 0.11.0 - 2016-06-22

### Added

- Support for parsing ZoneRecord webhooks.
- Support for listing options (GH-25).
- Support for Template API (GH-21).

## 0.10.0

Initial release.
