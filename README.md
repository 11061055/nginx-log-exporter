# 背景 （Background）


搭建了异地内网通信通道。 1. 需要对两边通道入口、出口做异常监控。 2. 需要对两边通道入口、出口做流量统计和性能分析。


# nginx-log-exporter

A Nginx log parser exporter for prometheus metrics.

![screen shot 2018-01-08 at 9 36 21 am](https://user-images.githubusercontent.com/1459834/34656613-7083cf3e-f457-11e7-929a-2758abad387b.png)


## Installation

1. go get `github.com/11061055/nginx-log-exporter`

## Usage

```
nginx-log-exporter -h 

Usage of:

  -config.file string
  
    	Nginx log exporter configuration file name. (default "config.yml")
      
  -web.listen-address string
  
    	Address to listen on for the web interface and API. (default ":6666")
      
exit status 2
```

## Configuration

```
- name: nginx
  source_files:
    - /tmp/access.log
  static_config:
    service: ucenter
  relabel_config:
    source_labels:
      - status
      - request
      - http_host
      - server_port
      - upstream_addr
      - request_method
      - upstream_status
    replacements:
      request:
        - trims:
          - sep: " "
            idx: 1
          - sep: "&"
            idx: 0
          replaces:
          - target: (.*)\?uid=(.*)
            value: $1?pid=$2
          - target: (.*)\?pid=(.*)
            value: $1?xxx=$2
        - trims:
          - sep: " "
            idx: 1
          - sep: "&"
            idx: 0
  histogram_buckets:
    start: 50
    step: 50
    num: 1
```

## name

service name, metric will be : 

`{name}_http_response_count_total`

`{name}_http_response_count_total`

`{name}_http_response_size_bytes`

`{name}_http_upstream_time_seconds`

`{name}_http_response_time_seconds`

## source_files

sevice nginx log, support multiple files.log must be in json style, or you can adjust by yourself.

## static_config

all metrics will add static label sets.

## relabel_config:

  * source_labels: what's labels should be use.
  
  * replacements: source labelvalue format rule, it supports regrex. 

## Output Style


app_http_response_count_total{foo="foo",method="GET",request="/app",status="200"} 2
app_http_response_count_total{foo="foo",method="GET",request="/app",status="200"} 1

app_http_response_size_bytes{foo="foo",method="GET",request="/app",status="200"} 70
app_http_response_size_bytes{foo="foo",method="GET",request="/app",status="200"} 21

app_http_response_time_seconds_bucket{foo="foo",method="GET",request="/app",status="200",le="0.5"} 2
app_http_response_time_seconds_count{foo="foo",method="GET",request="/app",status="200"} 2
app_http_response_time_seconds_bucket{foo="foo",method="GET",request="/app",status="200",le="0.5"} 1

app_http_response_time_seconds_sum{foo="foo",method="GET",request="/app",status="200"} 0.003
app_http_response_time_seconds_count{foo="foo",method="GET",request="/app",status="200"} 1

## Example

### replacements


```
    replacements:
      request:
        - trims:
          - sep: " "  // split the string by black character " "
            idx: 1    // the 1th index part is what we need
          - sep: "&"
            idx: 0
          replaces:
          - target: (.*)\?uid=(.*)   // regex math the string by target
            value: $1?pid=$2         // regex replace the string by value
          - target: (.*)\?pid=(.*)
            value: $1?xxx=$2
        - trims:
          - sep: "?"
            idx: 0
```

### log data

```
{ "timestamp": "15/Dec/2019:10:51:44 +0800", "remote_addr": "xx.xx.xx.xx", "request_time": "0.045", "status": "200", "request": "GET /app/api/ucenter/get?uid=123&pwd=123 HTTP/1.1", "request_method": "GET", "body_bytes_sent":"21", "http_x_clientip": "", "upstream_response_time": "0.046", "upstream_status": "200", "upstream_addr": "172.17.32.128:80", "request_body": ""}
```

```
1. request will be trimmed  to "/app/api/ucenter/get?uid=123&pwd=123"
2. request will be trimmed  to "/app/api/ucenter/get?uid=123"
3. request will be replaced to "/app/api/ucenter/get?pid=123"
4. request will be replaced to "/app/api/ucenter/get?xxx=123"
5. request will be trimmed  to "/app/api/ucenter/get"
```

That's to say, you can write your own trim and replace rules to run as many cycles as you want, with each output as the input of the next cycle. 非常灵活。


## Thanks

Forked form [prometheus-nginxlog-exporter](https://github.com/songjiayang/nginx-log-exporter)

Change 1. log format to json.

Change 2. add histogram buckets.

Change 3. add prometheus gauge and summary.

Change 4. trim and replace run as cycles with each output as the input of the next cycle. It's more flexible.
