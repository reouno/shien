# Installation

## Homebrew (macOS/Linux)

### Using tap (recommended)

```bash
brew tap yourusername/shien
brew install shien
```

### Start the service

```bash
brew services start shien
```

## Manual Installation

### From source

```bash
git clone https://github.com/yourusername/shien.git
cd shien
make install-all
```

### From binary releases

1. Download the latest release from [GitHub Releases](https://github.com/yourusername/shien/releases)
2. Extract the archive:
   ```bash
   tar -xzf shien-darwin-arm64.tar.gz
   ```
3. Move binaries to your PATH:
   ```bash
   sudo mv shien shienctl /usr/local/bin/
   ```

## Verify Installation

```bash
shienctl ping
```

## Uninstallation

### Homebrew

```bash
brew services stop shien
brew uninstall shien
```

### Manual

```bash
sudo rm /usr/local/bin/shien /usr/local/bin/shienctl
```