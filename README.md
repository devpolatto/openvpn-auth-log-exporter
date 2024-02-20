# OpenVPN Auth Log Exporter

This is a simple Golang application that builds an exporter to be used with Prometheus to get the OpenVPN server authentication logs

## Overview

This exporter is designed to expose metrics about a hypothetical service written in Golang. It provides insights into the performance and health of the service by exposing various metrics in the Prometheus exposition format.

## Prerequisites

- Golang installed on your machine ([Installation Guide](https://golang.org/doc/install))
- Prometheus server set up and running ([Prometheus Installation](https://prometheus.io/docs/prometheus/latest/installation/))

## Installation

1. Clone this repository:

    ```bash
    git clone https://github.com/devpolatto/openvpn-auth-log-exporter.git
    ```

2. Navigate to the project directory:

    ```bash
    cd openvpn-auth-log-exporter
    ```

3. Download the dependencies

    ```bash
    go mod download
    ```

4. Build the exporter:

    ```bash
    go build main .
    ```

## Usage

1. Run the exporter:

    ```bash
    ./main
    ```

     The script expects the user to indicate where the .log file to be monitored is located. By default, the script searches `/go/src/app/ovpn-ldap-auth.log`.

     To replace the location, run:
     ```bash
     ./main -openvpn.auth_paths /path/to/.log
     ```

2. Visit `http://localhost:8080/metrics` to access the metrics endpoint.

3. Configure Prometheus to scrape metrics from the exporter by adding the following configuration to your `prometheus.yml` file:

    ```yaml
    scrape_configs:
      - job_name: 'openvpn-auth-log-exporter'
        static_configs:
          - targets: ['192.168.10.10:9177']
    ```

    We chose to leave the path as `/metrics` because it is the prometheus default. You can change both the port and the path, as shown below:

    ```bash
    ./main -web.telemetry-path "/logs" -web.listen-address "9172"
    ```

4. Restart Prometheus to apply the changes.

## Metrics

The exporter displays the following log in the following format:

```html
openvpn_ldap_user_logs_total{DN="",MessageResponse="unable to bind",ResultCode="1",Timestamp="2024-02-05 21:32:06",Username="",action="connection",status="error"} 1
```

## Docker

To build a docker image, run:

```bash
cd openvpn-auth-log-exporter
```

```bash
docker build -t openvpn-auth-log-exporter .   
```

Then run the conatiner. Remember to make the volume bind ready for the .log file in your system file

```bash
docker run -p 9177:9177 -v ./ovpn-ldap-auth.log:/go/src/app/openvpn-status-auth.log openvpn-exporter:latest
```