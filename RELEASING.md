# Releasing

This document describes the steps to release a new version of DNSimple/Go.

## Prerequisites

- You have commit access to the repository
- You have push access to the repository
- You have a GPG key configured for signing tags

## Release process

1. **Determine the new version** using [Semantic Versioning](https://semver.org/)

   ```shell
   VERSION=X.Y.Z
   ```

   - **MAJOR** version for incompatible API changes
   - **MINOR** version for backwards-compatible functionality additions
   - **PATCH** version for backwards-compatible bug fixes

   The major part of the version number should match the version in the module path in `go.mod` (e.g., if module path is `github.com/dnsimple/dnsimple-go/v5`, the version should be `5.y.z`).

2. **Run tests** and confirm they pass

   ```shell
   go test -v ./...
   ```

3. **Update the version file** with the new version

   Edit `dnsimple.go` and update the `Version` constant:

   ```go
   Version = "$VERSION"
   ```

4. **Run tests** again and confirm they pass

   ```shell
   go test -v ./...
   ```

5. **Update the changelog** with the new version

   Finalize the `## main` section in `CHANGELOG.md` assigning the version.

6. **Commit the new version**

   ```shell
   git commit -a -m "Release $VERSION"
   ```

7. **Push the changes**

   ```shell
   git push origin main
   ```

8. **Wait for CI to complete**

9. **Create a signed tag**

   ```shell
   git tag -a v$VERSION -s -m "Release $VERSION"
   git push origin --tags
   ```

## Post-release

- Verify the new version appears on [pkg.go.dev](https://pkg.go.dev/github.com/dnsimple/dnsimple-go)
- Verify the GitHub release was created
- Announce the release if necessary
