# 指标导出器

`Venus` 系统的指标监控使用 `OpenCensus`，`OpenCensus` 是独立于供应商的，其收集的指标可以在本地显示，也可以将其发送到第三方分析工具或监控系统实现可视化。基于 `Venus` 的微服务架构，我们更倾向于将指标推送到一个监控服务系统进行统一可视化管理。

本文介绍 `OpenCensus` 推送到 `Prometheus Exporter` 的两种方式：

- 推送到 [Prometheus](https://github.com/prometheus/prometheus)：各组件需启动独立的指标监听服务，再通过 `Prometheus`  统一收集。
- 推送到 [Graphite](https://github.com/prometheus/graphite_exporter)：`Graphite` 有一个公开的收集器服务，各组件将指标推送到收集器。

`OpenCensus` 的介绍可参考文档 [OpenCensus](./OpenCensus.md)

## Graphite


## Prometheus
   
