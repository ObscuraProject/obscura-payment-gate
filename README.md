<div align="center">
  <h1><a href="Installation.md">Installation Guide</a></h1>
</div>

<div align="center">
  <img src="logo.png" alt="Prestimos Obscura Logo">
</div>

<div align="center">
  <h1>Obscura Free Market - Monero Marketplace Script</h1>
</div>

# Obscura Payment Gateway

![Bitcoin!](https://img.shields.io/badge/Bitcoin-000000?style=for-the-badge&logo=bitcoin&logoColor=white)
![Ethereum!](https://img.shields.io/badge/Ethereum-3C3C3D?style=for-the-badge&logo=ethereum&logoColor=white)
![Monero!](https://img.shields.io/badge/Monero-FF6600?style=for-the-badge&logo=monero&logoColor=white)
![Golang!](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![Node.js!](https://img.shields.io/badge/Node.js-339933?style=for-the-badge&logo=nodedotjs&logoColor=white)
![PostgreSQL!](https://img.shields.io/badge/PostgreSQL-316192?style=for-the-badge&logo=postgresql&logoColor=white)

## Overview

**Obscura Payment Gateway** is a comprehensive payment processing solution engineered to handle multi-currency transactions for privacy-focused marketplaces. Built with Go and Node.js, it provides secure, reliable cryptocurrency payment processing with support for major digital assets.

## Core Features

The payment gateway facilitates seamless cryptocurrency transactions with industry-leading security standards. It implements robust payment verification and settlement mechanisms, ensuring transaction integrity throughout the payment lifecycle. The system is designed to operate reliably in restricted network environments while maintaining full transaction transparency and auditability.

## Supported Cryptocurrencies

The gateway processes payments across multiple blockchain networks:

- **Bitcoin (BTC)** - Native blockchain payments with full RPC integration
- **Ethereum (ETH)** - Including ERC-20 token support (USDT and other standards)
- **Monero (XMR)** - Privacy-enhanced cryptocurrency payments

## Technology Stack

**Obscura Payment Gateway** leverages a modern, battle-tested technology stack optimized for reliability and performance. The backend is built with Go for high-performance transaction processing, while the frontend and auxiliary services utilize Node.js. PostgreSQL provides persistent data storage with ACID compliance, ensuring no transaction data is lost.

The system integrates with blockchain infrastructure components including Geth for Ethereum validation, Electrs for Bitcoin operations, and native Monero daemon integration for XMR transactions. Advanced monitoring and analysis tools such as Bitcoin RPC Explorer provide transparency into all payment activities.

## Cryptocurrency Payment Processing

The gateway handles the complete payment workflow, from invoice generation through settlement. It monitors blockchain networks in real-time, detects incoming transactions, and triggers automated settlement procedures. All payment states are tracked comprehensively, enabling merchants to monitor transaction status and reconcile payments with their accounting systems.

## Security & Privacy

The payment system prioritizes both technical security and user privacy. All communications can operate through Tor networks when required. The architecture separates transaction processing from user identification, enabling payments without unnecessary data exposure. Blockchain addresses are managed securely, and private key operations are isolated within secure enclaves.

## Enterprise Integration

The gateway provides standardized APIs for marketplace integration, allowing seamless payment flow within larger commerce platforms. Merchants can retrieve transaction histories, manage cryptocurrency wallets, and access detailed settlement reports through the administrative interface.

## License

**Obscura Payment Gateway** is licensed under the MIT License (MIT). Copyright © 2015-2024 Earthling.

The software is provided "as-is," without warranty of any kind, including but not limited to the warranties of merchantability and fitness for a particular purpose. In no event shall the authors or copyright holders be liable for any claim, damages, or other liability, whether in an action of contract, tort, or otherwise, arising from, out of, or in connection with the software or its use.
