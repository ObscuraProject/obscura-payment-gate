# Installation & Setup Guide for Obscura Payment Gate

## Overview

This guide covers the installation and configuration of all services referenced in your `settings.json` file. The configuration requires:

1. **PostgreSQL** - Database server
2. **Bitcoin/Electrum API** - Blockchain transaction indexing
3. **Ethereum API** - Ethereum node RPC
4. **Monero RPC** - Monero daemon and wallet RPC
5. **Bitcoin Exchange Rate API** - CryptoCompare API
6. **Etherscan API** - Ethereum blockchain explorer API

---

## Part 1: PostgreSQL Installation & Setup

### Prerequisites

- Administrative/sudo access
- 10GB+ disk space (for database)
- 2GB+ RAM

### Installation

#### On Ubuntu/Debian (Linux)

```bash
# Update package manager
sudo apt-get update

# Install PostgreSQL
sudo apt-get install postgresql postgresql-contrib

# Start PostgreSQL service
sudo systemctl start postgresql
sudo systemctl enable postgresql  # Auto-start on boot

# Verify installation
sudo -u postgres psql --version
```

#### On macOS

```bash
# Using Homebrew
brew install postgresql@16

# Start PostgreSQL
brew services start postgresql@16

# Add to PATH (add to ~/.zprofile or ~/.bash_profile)
export PATH="/usr/local/opt/postgresql@16/bin:$PATH"
```

#### On Windows

1. Download installer from https://www.postgresql.org/download/windows/
2. Run the installer executable
3. During installation:
   - Choose components: PostgreSQL Server, pgAdmin 4, Command Line Tools
   - Set a strong password for the `postgres` superuser
   - Accept default port: **5432**
   - Accept default locale
4. Complete installation
5. Add PostgreSQL `bin` directory to PATH: `C:\Program Files\PostgreSQL\16\bin`

### Configure PostgreSQL Connection

Your settings.json uses this connection string:

```
pg_connection_string: "host=localhost user=postgres dbname=go_p sslmode=disable"
```

#### Create the Database

```bash
# Connect to PostgreSQL as superuser
sudo -u postgres psql

# Create database 'go_p'
CREATE DATABASE go_p;

# Create a dedicated user for your application
CREATE USER appuser WITH PASSWORD 'your_secure_password';

# Grant privileges
GRANT ALL PRIVILEGES ON DATABASE go_p TO appuser;

# Exit psql
\q
```

#### Test Connection

```bash
# Test with postgres user
psql -h localhost -U postgres -d go_p

# Test with application user
psql -h localhost -U appuser -d go_p
```

---

## Part 2: Bitcoin/Electrum Setup

Your settings.json references:

```json
"electrum_api_url": "3.14.15.92:50001"
```

### Option A: Use Bitcoin Core + ElectrumX (Recommended for Full Control)

#### Step 1: Install Bitcoin Core

Follow the Bitcoin Core setup guide (see bitcoin-core-setup.md).

**Critical Configuration (bitcoin.conf):**

```conf
# Enable transaction indexing (REQUIRED for Electrum)
txindex=1

# RPC settings
server=1
rpcport=8332
rpcuser=bitcoinrpc
rpcpassword=your_secure_password
rpcbind=127.0.0.1
rpcallowip=127.0.0.1
```

#### Step 2: Install ElectrumX

**On Linux:**

```bash
# Install Python 3.8+
sudo apt-get install python3 python3-pip

# Create electrumx user
sudo useradd -m -s /bin/bash electrumx

# Switch to electrumx user
sudo su - electrumx

# Clone ElectrumX repository
git clone https://github.com/spesmilo/electrumx.git
cd electrumx

# Install dependencies
pip3 install -r requirements.txt

# Create configuration directory
mkdir -p ~/.electrumx
```

**Create ElectrumX Configuration File (~/.electrumx/electrumx.conf):**

```conf
# Coin and network
COIN = Bitcoin
NETWORK = mainnet

# Database location
DB_DIRECTORY = /home/electrumx/.electrumx/db

# Daemon connection (Bitcoin Core RPC)
DAEMON_URL = http://bitcoinrpc:your_secure_password@127.0.0.1:8332

# Services
SERVICES = tcp://0.0.0.0:50001,rpc://127.0.0.1:8000

# Peer discovery
PEER_DISCOVERY = self

# Logging
LOG_LEVEL = INFO
```

**Create Systemd Service (/etc/systemd/system/electrumx.service):**

```ini
[Unit]
Description=ElectrumX Server
After=network.target bitcoin.service

[Service]
Type=simple
User=electrumx
WorkingDirectory=/home/electrumx/electrumx
ExecStart=/home/electrumx/electrumx/electrumx_server
Restart=on-failure
RestartSec=10

[Install]
WantedBy=multi-user.target
```

**Start ElectrumX:**

```bash
# Reload systemd
sudo systemctl daemon-reload

# Enable on boot
sudo systemctl enable electrumx

# Start service
sudo systemctl start electrumx

# Check status
sudo systemctl status electrumx

# View logs
journalctl -u electrumx -f
```

### Option B: Connect to External Electrum Server

If you don't want to run your own ElectrumX server, connect to a public one:

```json
"electrum_api_url": "electrum.example.com:50001"
```

Popular public servers:
- `electrum.example.com:50001` (TCP)
- `electrum1.example.com:50001`
- `electrum2.example.com:50001`

---

## Part 3: Ethereum Node Setup

Your settings.json references:

```json
"ethereum_api_url": "http://3.14.15.92:8545"
```

### Option A: Run Geth (Go-Ethereum)

**Installation on Linux:**

```bash
# Add Ethereum PPA
sudo add-apt-repository -y ppa:ethereum/ethereum

# Install Geth
sudo apt-get update
sudo apt-get install ethereum

# Create data directory
mkdir -p ~/ethereum-data
```

**Create Systemd Service (/etc/systemd/system/geth.service):**

```ini
[Unit]
Description=Geth Ethereum Node
After=network.target

[Service]
Type=simple
User=ethereum
WorkingDirectory=/home/ethereum
ExecStart=/usr/bin/geth --http --http.addr 0.0.0.0 --http.port 8545 --datadir /home/ethereum/ethereum-data --syncmode fast
Restart=on-failure
RestartSec=10

[Install]
WantedBy=multi-user.target
```

**Start Geth:**

```bash
sudo systemctl daemon-reload
sudo systemctl enable geth
sudo systemctl start geth
```

**Initial sync time: 4-12 hours** (depending on network conditions)

### Option B: Connect to External Ethereum RPC

Use a public RPC endpoint instead:

```json
"ethereum_api_url": "http://infura.io:443/v3/YOUR_INFURA_KEY"
```

Popular public endpoints:
- **Infura**: `https://mainnet.infura.io/v3/YOUR_PROJECT_ID`
- **Alchemy**: `https://eth-mainnet.alchemyapi.io/v2/YOUR_API_KEY`
- **QuickNode**: `https://your-rpc-endpoint.quiknode.pro/`

---

## Part 4: Monero RPC Setup

Your settings.json references:

```json
"monero_rpc_url": "http://3.14.15.92:8080/",
"monero_wallet_rpc_url": "http://3.14.15.92:8083/",
"monero_wallet_rpc_user": "",
"monero_wallet_rpc_password": ""
```

### Prerequisites

- Monero daemon running on port 8080
- Monero wallet RPC on port 8083

### Step 1: Install Monero

**On Linux:**

```bash
# Download latest Monero release
cd ~/downloads
wget https://downloads.getmonero.org/linux64

# Extract
tar xjf linux64

# Move to system directory
sudo mv monero-linux-x64-*/monero* /usr/local/bin/

# Verify installation
monero-daemon --version
```

**On macOS:**

```bash
# Using Homebrew
brew install monero

# Verify
monero-daemon --version
```

**On Windows:**

1. Download from https://www.getmonero.org/downloads/
2. Extract to `C:\Monero\`
3. Add to PATH

### Step 2: Start Monero Daemon

**Create Systemd Service (/etc/systemd/system/monerod.service):**

```ini
[Unit]
Description=Monero Daemon
After=network.target

[Service]
Type=simple
User=monero
WorkingDirectory=/home/monero
ExecStart=/usr/local/bin/monerod --rpc-bind-ip 0.0.0.0 --rpc-bind-port 8080 --confirm-external-bind
Restart=on-failure
RestartSec=10

[Install]
WantedBy=multi-user.target
```

**Start Daemon:**

```bash
sudo systemctl daemon-reload
sudo systemctl enable monerod
sudo systemctl start monerod

# Monitor sync progress
monero-cli --daemon-address 127.0.0.1:8080 status
```

**Initial sync: 2-4 hours** (faster with pruning: add `--prune-blockchain` flag)

### Step 3: Create Monero Wallet

**Generate or restore wallet:**

```bash
# Create new wallet
monero-wallet-cli

# Or restore from seed
monero-wallet-cli --restore-deterministic-wallet
```

### Step 4: Start Monero Wallet RPC

**Create Systemd Service (/etc/systemd/system/monero-wallet-rpc.service):**

```ini
[Unit]
Description=Monero Wallet RPC
After=network.target monerod.service

[Service]
Type=simple
User=monero
WorkingDirectory=/home/monero
ExecStart=/usr/local/bin/monero-wallet-rpc \
  --rpc-bind-ip 0.0.0.0 \
  --rpc-bind-port 8083 \
  --wallet-file /home/monero/wallets/primary \
  --password YOUR_WALLET_PASSWORD \
  --daemon-address 127.0.0.1:8080 \
  --rpc-login username:password \
  --confirm-external-bind
Restart=on-failure
RestartSec=10

[Install]
WantedBy=multi-user.target
```

**Update settings.json:**

```json
"monero_wallet_rpc_user": "username",
"monero_wallet_rpc_password": "password",
"monero_wallet_rpc_url": "http://3.14.15.92:8083/"
```

**Start Wallet RPC:**

```bash
sudo systemctl daemon-reload
sudo systemctl enable monero-wallet-rpc
sudo systemctl start monero-wallet-rpc
```

---

## Part 5: API Keys Configuration

Your settings.json requires these API keys:

### 1. CryptoCompare API Token

```json
"cryptocompare_token": "YOUR_CRYPTOCOMPARE_API_KEY"
```

**Setup:**

1. Visit https://www.cryptocompare.com/
2. Sign up for a free account
3. Navigate to API section
4. Generate API key
5. Add to settings.json

### 2. Etherscan API Key

```json
"etherscan_api_key": "YOUR_ETHERSCAN_API_KEY"
```

**Setup:**

1. Visit https://etherscan.io/
2. Sign up for an account
3. Go to API Keys section
4. Create new API key
5. Add to settings.json

---

## Part 6: Wallet Addresses Configuration

Your settings.json requires:

```json
"btc_commission_wallet": "bitcoin_address_here",
"tochka_escrow_payment_address": "bitcoin_address_here",
"xmr_commission_wallet": "monero_address_here"
```

### Generate Bitcoin Address

```bash
# Using Bitcoin Core
bitcoin-cli getnewaddress

# Or use Electrum wallet GUI
# Create new wallet and generate address
```

### Generate Monero Address

```bash
# Get primary address from wallet
monero-wallet-cli address

# Or check wallet RPC
curl -X POST http://username:password@127.0.0.1:8083/json_rpc \
  -H 'Content-Type: application/json' \
  -d '{"jsonrpc":"2.0","id":"0","method":"getaddress"}'
```

---

## Part 7: Complete Settings.json Template

```json
{
  "bre_url": "http://192.168.88.243:3002/",
  "btc_commission_wallet": "1A1z7agoat4FCWW3x6x7awUYT1mKSAeNRM",
  "cryptocompare_token": "YOUR_CRYPTOCOMPARE_API_KEY",
  "electrum_api_url": "127.0.0.1:50001",
  "ethereum_api_url": "http://127.0.0.1:8545",
  "etherscan_api_key": "YOUR_ETHERSCAN_API_KEY",
  "host": "127.0.0.1",
  "monero_rpc_url": "http://127.0.0.1:8080/",
  "monero_wallet_rpc_password": "your_wallet_rpc_password",
  "monero_wallet_rpc_url": "http://127.0.0.1:8083/",
  "monero_wallet_rpc_user": "username",
  "pg_connection_string": "host=localhost user=postgres dbname=go_p password=your_postgres_password sslmode=disable",
  "port": 8083,
  "tochka_escrow_payment_address": "1A1z7agoat4FCWW3x6x7awUYT1mKSAeNRM",
  "xmr_commission_wallet": "45TTEXeUp9P2RK3R8yVjRDfknECM7DcTZLJhBXFkT2RFbhHEPKMRGVERb2dSjMpzKH3zqmGVwYjqVcsapMvSYUkqV5MU1a4"
}
```

---

## Verification Checklist

- [ ] PostgreSQL running on localhost:5432
- [ ] Database `go_p` created
- [ ] Bitcoin Core synced with `txindex=1` enabled
- [ ] ElectrumX running and accessible on port 50001
- [ ] Ethereum node (Geth or external RPC) accessible on port 8545
- [ ] Monero daemon synced on port 8080
- [ ] Monero wallet RPC running on port 8083
- [ ] CryptoCompare API key obtained
- [ ] Etherscan API key obtained
- [ ] Bitcoin wallet addresses generated
- [ ] Monero wallet address generated
- [ ] All settings.json fields populated with valid values

---

## Troubleshooting

### PostgreSQL Connection Failed

```bash
# Check if PostgreSQL is running
sudo systemctl status postgresql

# Test connection manually
psql -h localhost -U postgres -d go_p

# Check PostgreSQL logs
sudo tail -f /var/log/postgresql/postgresql.log
```

### ElectrumX Not Starting

```bash
# Check logs
journalctl -u electrumx -n 50

# Verify Bitcoin Core is running and RPC accessible
bitcoin-cli getblockcount

# Check if port 50001 is in use
netstat -tulpn | grep 50001
```

### Monero Sync Slow

- Use `--prune-blockchain` flag to reduce disk usage
- Use `--syncmode light` for Ethereum (Geth)
- Increase `--max-connections` for faster sync

### API Connection Errors

1. Verify services are running on correct ports
2. Check firewall rules
3. Verify IP addresses in settings.json
4. Test with curl: `curl http://127.0.0.1:8545` (should show HTML or error)

---

## Security Considerations

1. **Store passwords securely** - Use environment variables or secure vaults
2. **Restrict RPC access** - Use firewall rules, don't expose on public internet
3. **Use SSL/TLS** - Wrap services in reverse proxy (nginx)
4. **Monitor logs** - Set up log aggregation and alerting
5. **Backup wallets** - Regular backups of Monero and Bitcoin wallets
6. **Use view-only wallets** - For Monero on servers (recommended for production)

---

*Last Updated: November 2025*
