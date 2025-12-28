---
description: Build release binaries for all platforms
argument-hint: <version>
allowed-tools: [Bash]
model: claude-haiku-4-5-20251001
---

# Build Release $ARGUMENTS

## Release Build Process

1. **Validate Version**
   - Verify semantic versioning format (e.g., v1.0.0)
   - Check version is not already released
   - Update version in main.go if needed

2. **Run Tests**
   - Execute: `go test ./... -race`
   - Ensure all tests pass
   - No broken builds in release

3. **Build Binaries**
   - Linux: `GOOS=linux GOARCH=amd64 go build -o claude-squad-linux-amd64`
   - macOS: `GOOS=darwin GOARCH=amd64 go build -o claude-squad-darwin-amd64`
   - Windows: `GOOS=windows GOARCH=amd64 go build -o claude-squad-windows-amd64.exe`
   - ARM: `GOOS=linux GOARCH=arm64 go build -o claude-squad-linux-arm64`

4. **Create GitHub Release**
   - Run: `git tag v$ARGUMENTS`
   - Run: `git push origin v$ARGUMENTS`
   - Run: `gh release create v$ARGUMENTS ./claude-squad-*`

5. **Generate Checksums**
   - Create: `sha256sum claude-squad-* > checksums.txt`
   - Include in release assets
   - Sign checksums if configured

6. **Update Documentation**
   - Update CHANGELOG.md
   - Update README with new version
   - Document breaking changes if any

7. **Announce Release**
   - Create release notes
   - Highlight new features
   - Thank contributors

## Output

- Release binaries for all platforms
- Checksum verification file
- GitHub release with assets
- Updated documentation
