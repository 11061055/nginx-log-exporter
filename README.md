

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
- name: app

  source_files:
  
    - ./data/logs/nginx/service_nginx_80.log
    
  static_config:
  
    foo: foo
    
  relabel_config: 
  
    source_labels: 
    
      - request
      - method
      - status
      
    replacement:
    
      request: 
      
        trim: "?"
        
        replace:
        
          - target: /app/[0-9]+/api
            value: /app/xxx/api
```

## name

service name, metric will be : 

`{name}_http_response_count_total`, 
`{name}_http_response_count_total`, 
`{name}_http_response_size_bytes`, 
`{name}_http_upstream_time_seconds`, 
`{name}_http_response_time_seconds`

## source_files

sevice nginx log, support multiple files.log must be in json style, or you can adjust by yourself.

## static_config

all metrics will add static labelsets.

## relabel_config:

  * source_labels: what's labels should be use.
  
  * replacement: source labelvalue format rule, it supports regrex. 

## Example


app_http_response_count_total{foo="foo",method="GET",request="/v1.0/example",status="200"} 2
app_http_response_count_total{foo="foo",method="GET",request="/v1.0/example/:id",status="200"} 1

app_http_response_size_bytes{foo="foo",method="GET",request="/v1.0/example",status="200"} 70
app_http_response_size_bytes{foo="foo",method="GET",request="/v1.0/example/:id",status="200"} 21

app_http_response_time_seconds_bucket{foo="foo",method="GET",request="/v1.0/example",status="200",le="0.005"} 2

app_http_response_time_seconds_count{foo="foo",method="GET",request="/v1.0/example",status="200"} 2
app_http_response_time_seconds_bucket{foo="foo",method="GET",request="/v1.0/example/:id",status="200",le="0.005"} 1

app_http_response_time_seconds_sum{foo="foo",method="GET",request="/v1.0/example/:id",status="200"} 0.003
app_http_response_time_seconds_count{foo="foo",method="GET",request="/v1.0/example/:id",status="200"} 1

```

## Thanks

- Inspired by [prometheus-nginxlog-exporter](https://github.com/songjiayang/nginx-log-exporter)
