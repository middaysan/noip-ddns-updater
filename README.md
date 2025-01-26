# No-IP Dynamic DNS Updater

## Description
`noip-ddns-updater` is a lightweight command-line tool written in Go that updates a No-IP dynamic DNS record. It periodically checks the public IP address and updates the No-IP hostname if the IP has changed.

## Features
- Automatically detects the public IP address.
- Updates the No-IP hostname only when the IP has changed.
- Uses environment variables for configuration.
- Can be run as a background service.
- Supports a custom update interval.

## Installation
1. Clone the repository:
   ```sh
   git clone https://github.com/middaysan/noip-ddns-updater.git
   cd noip-ddns-updater
   ```
2. Build the application:
   ```sh
   go build -o noip_ddns_updater
   ```

## Usage
Run the program:
```sh
./noip_ddns_updater
```

### Command-line options
- `-h`: Displays usage information and exits.

### Environment Variables
Set the following environment variables before running the program:
```sh
export NOIP_USER="your_noip_username"
export NOIP_PASS="your_noip_password"
export NOIP_HOST="your_hostname.no-ip.com"
export NOIP_INTERVAL_MINUTES=5
```

#### Optional:
```sh
export NOIP_URL="https://dynupdate.no-ip.com/nic/update"  # Default No-IP update URL
export CHECK_IP_URL="https://checkip.amazonaws.com"      # Default IP check URL
```

## Example
```sh
NOIP_USER="example_user" NOIP_PASS="example_pass" NOIP_HOST="example.ddns.net" NOIP_INTERVAL_MINUTES=10 ./noip_ddns_updater
```

## License
This project is licensed under the MIT License.

## Author
[Anton Rudy](https://github.com/middaysan)

