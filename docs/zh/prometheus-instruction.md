# prometheus 使用说明

## 配置解析

默认的配置文件为 `prometheus.yml`
```yaml
# my global config
global:
  scrape_interval: 15s # Set the scrape interval to every 15 seconds. Default is every 1 minute.
  evaluation_interval: 15s # Evaluate rules every 15 seconds. The default is every 1 minute.
  # scrape_timeout is set to the global default (10s).

# Alertmanager configuration
alerting:
  alertmanagers:
    - static_configs:
        - targets:
          # - alertmanager:9093

# Load rules once and periodically evaluate them according to the global 'evaluation_interval'.
rule_files:
  # - "first_rules.yml"
  # - "second_rules.yml"

# A scrape configuration containing exactly one endpoint to scrape:
# Here it's Prometheus itself.
scrape_configs:
  # The job name is added as a label `job=<job_name>` to any timeseries scraped from this config.
  - job_name: "prometheus"

    # metrics_path defaults to '/metrics'
    # scheme defaults to 'http'.

    static_configs:
      - targets: ["localhost:9090"]
 
  # 配置采集点,一个job可对应一个或多个采集点
  - job_name: "venus"
  
    metrics_path: "/metrics"
    scheme: "http"
  
    static_configs:
      - targets: ["localhost:4567", "localhost:5678"]
```

### 服务发现

可以通过额外的文件来配置采集点，支持热加载，相当于基于文件的服务发现。

`prometheus.yml`
```yaml
scrape_configs:
  # The job name is added as a label `job=<job_name>` to any timeseries scraped from this config.
  - job_name: "prometheus"

    # metrics_path defaults to '/metrics'
    # scheme defaults to 'http'.

    static_configs:
      - targets: ["localhost:9090"]
 
  # 热加载
  - job_name: "venus"
   
    file_sd_configs:
      - files: 
        - "./venus.yml"
      
        # 多久重新加载
        refresh_interval: 10m
```

venus.yml
```yaml
# 该文件中的每一个targets都是一个采集点
- targets:
  - "localhost:4567"
  labels:
    __metrics_path__: "/metrics"
    instance: "miner"
  
- targets:
  - "localhost:5678"
  labels:
    __metrics_path__: "/metrics"  
    instance: "message"
```

## 启动
```bash
$ ./prometheus --config.file=prometheus.yml
```
http://localhost:9090/targets, 可查询监控组件指标；


