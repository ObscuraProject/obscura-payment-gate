# Bitcoin Core Download & Configuration Guide
## For ElectrumX Integration with Transaction Indexing

---

## Overview

This guide provides step-by-step instructions for downloading Bitcoin Core and configuring it with transaction indexing (`txindex=1`) and RPC access—essential for running ElectrumX servers or other applications requiring full transaction history access.

**Important Note:** Enabling `txindex=1` dramatically increases both initial sync time (60-80+ hours) and disk space requirements (approximately 600GB+ for the blockchain data).

---

## Part 1: Download Bitcoin Core

### Step 1: Visit the Official Website

Navigate to the official Bitcoin Core download page at **https://bitcoincore.org/en/download/**

This is the only secure source for Bitcoin Core. Always verify you're on the official domain.

### Step 2: Choose Your Operating System

Bitcoin Core is available for:

- **Windows** (32-bit and 64-bit)
- **macOS**
- **Linux** (Ubuntu, Debian, Fedora, and other distributions)

Select the installer or binary appropriate for your system architecture.

### Step 3: Download Hashes and Signatures

For security verification (highly recommended):

1. Download the **SHA256SUMS** file
2. Download the **SHA256SUMS.asc** file (PGP signature)
3. Download the **public keys** from the GitHub repository: https://github.com/bitcoin-core/bitcoin-core.org

### Step 4: Verify the Download (Optional but Recommended)

Verifying the signature ensures your downloaded file hasn't been tampered with.

**On Linux/macOS:**

```bash
# Verify the SHA256 hash
sha256sum -c SHA256SUMS 2>&1 | grep OK

# Verify the PGP signature (requires gpg installed)
gpg --import keys.txt
gpg --verify SHA256SUMS.asc SHA256SUMS
```

**On Windows:**

Use a tool like **7-Zip** or **WinRAR** to verify hashes, or use PowerShell:

```powershell
(Get-FileHash -Path "bitcoin-core-installer.exe" -Algorithm SHA256).Hash
```

Compare the output with the official SHA256SUMS file.

### Step 5: Install Bitcoin Core

**Windows:**
- Run the `.exe` installer and follow the on-screen prompts
- Default installation path: `C:\Program Files\Bitcoin\` or `%APPDATA%\Bitcoin`

**macOS:**
- Open the `.dmg` file and drag Bitcoin Core to the Applications folder
- Or use Homebrew: `brew install bitcoin`

**Linux:**
- Extract the binary: `tar xzf bitcoin-*-x86_64-linux-gnu.tar.gz`
- Copy to system path: `sudo mv bitcoin-*/bin/* /usr/local/bin/`
- Or use package manager: `sudo apt-get install bitcoin-qt` (Ubuntu/Debian)

---

## Part 2: Configure Bitcoin Core

### Create the bitcoin.conf Configuration File

Bitcoin Core reads configuration from a `bitcoin.conf` file located in the Bitcoin data directory.

**Data Directory Locations:**

| OS | Path |
|----|------|
| **Linux** | `~/.bitcoin/bitcoin.conf` |
| **macOS** | `~/Library/Application Support/Bitcoin/bitcoin.conf` |
| **Windows** | `%APPDATA%\Bitcoin\bitcoin.conf` |

### Step 1: Create or Edit bitcoin.conf

**On Linux/macOS:**

```bash
nano ~/.bitcoin/bitcoin.conf
```

Or use your preferred text editor (vim, code, etc.)

**On Windows:**

1. Open File Explorer
2. Navigate to `%APPDATA%\Bitcoin\` (type this in the address bar)
3. Right-click in the empty space → New → Text Document
4. Name it `bitcoin.conf`
5. Open it with Notepad or your preferred editor

### Step 2: Add Configuration Entries

Add the following lines to your `bitcoin.conf` file:

```conf
# Enable server mode (required for RPC)
server=1

# Run Bitcoin Core in the background (optional)
daemon=1

# Enable transaction indexing (CRITICAL for ElectrumX)
txindex=1

# Enable block filter index (improves performance for certain applications)
blockfilterindex=1

# Listen for incoming connections
listen=1

# Disable wallet (optional, if you only want a node without wallet functionality)
disablewallet=0

# Database cache size in MB (increase for faster sync, default: 300)
dbcache=1000

# RPC Port (default: 8332 for mainnet)
rpcport=8332

# RPC Listen Address (default: 127.0.0.1 for localhost only)
rpcbind=127.0.0.1

# RPC Allow IP Address
# Only allow local connections (secure)
rpcallowip=127.0.0.1

# For remote RPC access (use with caution):
# rpcallowip=0.0.0.0/0  # Allow all (NOT recommended for security)
# rpcallowip=<your-ip>  # Allow specific IP
```

### Step 3: Configure RPC Credentials

Bitcoin Core supports two methods for RPC authentication:

#### Method 1: Using rpcuser and rpcpassword (Deprecated but Simpler)

```conf
rpcuser=youruser
rpcpassword=yourpassword
```

**Note:** This method is deprecated. Use Method 2 for better security.

#### Method 2: Using rpcauth (Recommended)

Generate a secure `rpcauth` entry using Bitcoin Core's built-in tool:

**On Linux/macOS:**

```bash
python3 -c "import hmac; import binascii; import os; salt=binascii.hexlify(os.urandom(16)).decode(); pw='yourpassword'; hashed=binascii.hexlify(hmac.new(salt.encode(), pw.encode(), 'sha256').digest()).decode(); print('rpcauth=' + 'user:' + salt + '$' + hashed)"
```

Replace `yourpassword` with your desired password.

**On Windows:**

Use the same Python command in PowerShell, or download a `rpcauth` generator script from GitHub.

The output will look like:

```
rpcauth=user:72de450660cdb6dd2689cd2cba4091646a5e8005490dec07dc577b6dad608a80:$cbb36c03b15219cafb1e72ae9329d5fd
```

Add this line to your `bitcoin.conf`:

```conf
rpcauth=user:72de450660cdb6dd2689cd2cba4091646a5e8005490dec07dc577b6dad608a80:$cbb36c03b15219cafb1e72ae9329d5fd
```

Save your password securely—you'll need it to authenticate RPC calls.

---

## Part 3: Start Bitcoin Core and Sync

### Save Configuration

Save and close your `bitcoin.conf` file.

### Start Bitcoin Core

**GUI Mode (Default):**

Simply launch the Bitcoin Core application from your applications menu or start menu.

**Daemon Mode (Command Line):**

**Linux/macOS:**

```bash
bitcoind
```

Or with specific data directory:

```bash
bitcoind -datadir=/path/to/bitcoin/data
```

**Windows (PowerShell):**

```powershell
C:\Program Files\Bitcoin\bitcoin-qt.exe
# or
C:\Program Files\Bitcoin\bitcoind.exe
```

### Monitor the Sync Progress

**Using the GUI:**

Bitcoin Core will display sync progress in the bottom status bar. Expect 60-80+ hours for full initial block download (IBD) with `txindex=1` enabled.

**Using the Command Line:**

```bash
bitcoin-cli getblockchaininfo
```

This returns:

```json
{
  "chain": "main",
  "blocks": 823000,
  "headers": 823000,
  "bestblockhash": "...",
  "difficulty": 84000000000,
  "mediantime": 1640000000,
  "verificationprogress": 0.9999,
  "initialblockdownload": false,
  "chainwork": "..."
}
```

When `initialblockdownload` shows `false`, your node is fully synced.

---

## Part 4: Test RPC Access

Once Bitcoin Core is synced, test RPC connectivity:

```bash
bitcoin-cli -rpcuser=user -rpcpassword=yourpassword getblockchaininfo
```

Or if using `rpcauth`:

```bash
bitcoin-cli getblockchaininfo
```

Successful output confirms RPC is working.

---

## Part 5: Create Systemd Service (Linux/macOS)

To run Bitcoin Core automatically on system startup, create a systemd service:

**Create /etc/systemd/system/bitcoin.service:**

```ini
[Unit]
Description=Bitcoin Core Daemon
After=network.target

[Service]
Type=forking
User=bitcoin
Group=bitcoin
WorkingDirectory=/home/bitcoin
ExecStart=/usr/local/bin/bitcoind -daemon -datadir=/home/bitcoin/.bitcoin
ExecStop=/usr/local/bin/bitcoin-cli -datadir=/home/bitcoin/.bitcoin stop
Restart=on-failure
RestartSec=10

[Install]
WantedBy=multi-user.target
```

**Enable and Start:**

```bash
# Reload systemd
sudo systemctl daemon-reload

# Enable on boot
sudo systemctl enable bitcoin

# Start the service
sudo systemctl start bitcoin

# Check status
sudo systemctl status bitcoin

# View logs
journalctl -u bitcoin -f
```

---

## Critical Considerations

### Disk Space

- **Without txindex:** ~600GB (blockchain data only)
- **With txindex=1:** ~800GB-1.2TB+ (includes full transaction index)
- Keep 20-30% free space for safe operation

### Sync Time

- **Without txindex:** 24-48 hours
- **With txindex=1:** 60-80+ hours (significantly slower due to indexing overhead)

### Performance

Set `dbcache` appropriately for your system:

- **4GB RAM system:** `dbcache=500`
- **8GB RAM system:** `dbcache=1000`
- **16GB+ RAM system:** `dbcache=2000`

Higher cache speeds up sync but uses more RAM.

### Security Best Practices

1. **Always download from https://bitcoincore.org**
2. **Verify PGP signatures** before installation
3. **Restrict RPC access:** Only allow necessary IPs via `rpcallowip`
4. **Use rpcauth** instead of basic username/password
5. **Run Bitcoin Core on a secure, firewalled machine**
6. **Keep bitcoin.conf permissions restricted:** `chmod 600 ~/.bitcoin/bitcoin.conf`
7. **Use VPN or firewall** to restrict access to port 8332 (RPC)
8. **Monitor wallet.dat** - Keep regular backups in secure location

---

## Troubleshooting

### Bitcoin Core Won't Start

- **Check permissions:** Ensure your data directory is readable/writable
- **Port conflict:** Verify port 8332 (RPC) and 8333 (P2P) are available
- **Check logs:** Look in `~/.bitcoin/debug.log` for error messages

**Command to check port availability:**

```bash
# Linux/macOS
lsof -i :8333
lsof -i :8332

# Windows
netstat -ano | findstr :8333
netstat -ano | findstr :8332
```

### RPC Connection Fails

```bash
bitcoin-cli getblockchaininfo
# Error: connect() failed after 5 seconds: Connection refused
```

**Solutions:**

- Ensure `server=1` is in bitcoin.conf
- Restart Bitcoin Core: `bitcoin-cli stop` then `bitcoind`
- Verify RPC credentials in bitcoin.conf
- Check firewall rules
- Verify Bitcoin Core is actually running: `ps aux | grep bitcoind`

### Slow Initial Block Download

- Increase `dbcache` (requires more RAM)
- Increase `maxconnections=128` (default: 125)
- Ensure stable, high-speed internet connection
- Verify `txindex=1` is set (causes intentional slowdown)
- Check logs for peer connection issues: `tail -f ~/.bitcoin/debug.log`

### txindex Reindexing Taking Too Long

If you enable `txindex=1` after the initial sync, Bitcoin Core will reindex:

```bash
# This can take 12-24 hours. Monitor progress:
bitcoin-cli getblockchaininfo | grep blocks
```

To speed up reindexing:
- Increase `dbcache` to 2000-4000 MB
- Ensure sufficient disk I/O bandwidth
- Don't run other heavy I/O operations

### Out of Disk Space

Bitcoin Core will stop syncing if it runs out of disk space:

```bash
# Check available space
df -h ~/.bitcoin

# Monitor space usage
du -sh ~/.bitcoin
```

**Solutions:**
- Add more storage to your system
- Enable pruning (not recommended with txindex): `prune=550` (550 MB minimum)
- Use a faster disk (SSD vs HDD)

---

## ElectrumX Integration

Once Bitcoin Core is fully synced with `txindex=1` and RPC enabled, configure ElectrumX to connect:

**ElectrumX Configuration:**

```
DAEMON_URL = http://user:yourpassword@127.0.0.1:8332
```

Restart ElectrumX to establish connection to your fully-indexed Bitcoin node.

---

## Additional Resources

- Official Bitcoin Core: https://bitcoincore.org
- Bitcoin Core Documentation: https://bitcoin.org/en/developer-documentation
- RPC API Reference: https://developer.bitcoin.org/reference/rpc/
- Bitcoin Configuration File Guide: https://en.bitcoin.it/wiki/Running_Bitcoin
- Bitcoin Core GitHub: https://github.com/bitcoin/bitcoin
- Bitcoin Core Release Notes: https://github.com/bitcoin/bitcoin/tree/master/doc/release-notes

---

## Quick Reference Commands

```bash
# Check blockchain info
bitcoin-cli getblockchaininfo

# Get current block height
bitcoin-cli getblockcount

# Get network info
bitcoin-cli getnetworkinfo

# Get peer connections
bitcoin-cli getpeerinfo

# Stop Bitcoin Core gracefully
bitcoin-cli stop

# Get RPC help
bitcoin-cli help

# Test RPC with authentication
bitcoin-cli -rpcuser=user -rpcpassword=pass getblockchaininfo
```

---

*Last Updated: November 2025*
