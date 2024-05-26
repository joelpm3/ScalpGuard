# ScalpGuard 
## Overview
This project serves as the proof of concept (PoC) for ScalpGuard. In the dynamic realm of online retail, scalper bots pose a growing threat by manipulating the purchasing process, thus denying regular customers access to limited or highly demanded products. In response to this issue, this endeavor was initiated to develop and assess effective methods for detecting and thwarting such bots, with the aim of ensuring fair purchasing processes.
## Project Structure
- cmd/ScalpGuard: Main executable Go programs.
- config: Configuration files for browsers and servers.
- examples: Backend implementations and other examples in Python.
- pkg: Go packages including algorithms, models, routers, and server configurations.
## Getting Started
### Prerequisites
Go 1.15 or higher
Python 3.8 or higher (If you want to use Python for the Backend) 
### Installation
1. Clone the repository:
```bash
git clone https://github.com/joelpm3/ScalpGuard.git
```
2. Navigate into the project directory:
```bash
cd ScalpGuard
```
### Running the Application
Build the Go server and run it:

For Windows:
```bash
go build -o server.exe ./cmd/ScalpGuard
```
```bash
./server.exe
```
For Linux:
```bash
go build -o server ./cmd/ScalpGuard
```
```bash
chmod +x server
```
```bash
./server
```

Run the backend example:
```bash
python3 examples/backend_example/main.py
```

## Configuring Backends
To add a new backend:
1. Define the backend configuration in config/server.yaml.
2. Implement the backend interface in the language of your choice. Examples are provided in the examples directory.

### server.yaml example
Modify the server.yaml configuration files to configure the backend:
- user_agents: A list of regex patterns
- tls_fingerprints_ja3: A list of JA3 hashes
- tls_fingerprints_ja3_no_extension: A list of JA3 without extensions hashes
- http2_fingerprints: A list of HTTP/2 hashes
- browser_settings: A list of browser settings ([Browser Settings](#browser-settings))

```yaml
backends:
  - name: "Produktivserver"
    host: "127.0.0.1"
    port: 5000
    whitelist:
      user_agents:
        - Mozilla\/5\.0 \([^)]*\) AppleWebKit\/537\.36 \(KHTML, like Gecko\) Chrome\/[\d\.]+ Safari\/537\.36
      tls_fingerprints_ja3:
        - "9312cc2f02f42764c70c3713f80fe219"
      tls_fingerprints_ja3_no_extension:
        - "e21aa11bf290e2d8cb2f82fe12e03d56"
      http2_fingerprints:
        - "0fd8ac32204eb28a46c51479d6d7be10"
      browser_settings:
        - "chrome122"
      
    blacklist:
      user_agents: 
      tls_fingerprints_ja3:
      tls_fingerprints_ja3_no_extension:
      http2_fingerprints:
      browser_settings:
```

### Browser Settings
To efficiently manage blacklisting or whitelisting for specific collections of user agents or fingerprints, you can utilize a browser configuration, functioning akin to a group. This configuration allows you to apply the same group across multiple backends, ensuring seamless control and consistency across your system.
```yaml
browsers:
  - name: "chrome122"
    user_agents: 
      - Mozilla\/5\.0 \([^)]*\) AppleWebKit\/537\.36 \(KHTML, like Gecko\) Chrome\/[\d\.]+ Safari\/537\.36
    tls_fingerprints_ja3:
    tls_fingerprints_ja3_no_extension:
      - "80abbd14ce91c0b5dc8d5c095838b76a"
    http2_fingerprints:
      - "5a40d192fed6045d13e5980884bb1034"
```
