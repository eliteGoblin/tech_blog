## Ideas

*  WHY
    -  feel confident
    -  prevent problems
    -  let customer tell you the problem is akward

*  show example of netdata
    -  input parameter from web(form)
    -  start a lot of goroutines, memory leak
    -  draw sin function, 各种波形

*  monitoring system has two customers
    -  The business
    -  Information Technology

*  As the Reactive environment is infrastructure-centric, it also only serves a segment of our technology customer—generally only operational teams—and doesn't provide useful, application-centric data to developers.
*  If a metric is measuring then the service is available. If it stops measuring then it's likely the service is not available.
*  Visualization of those events, metrics, and logs will also allow for the ready expression and interpretation of complex ideas that would otherwise take thousands of words or hours of explanation.
*  event, log, and metrics-centric.
*  Focus on "whitebox" or push-based monitoring instead of "blackbox" or pull-based monitoring.
*  Most monitoring systems are pull/polling-based systems. An excellent example is Nagios. With Nagios, your monitoring system generally queries the components being monitored;
*  Pull vs Push
    -  pull: 
        +   how to scale? puller需要知道所有的host. pull表明: monitor为中心
        +   generally emphasize monitoring availability, rather than quality and service
        +   Blackbox monitoring probes the outside of a service or application,
    -  push: 
        +  完全解耦
        +  更安全: 只需要配置egress即可
*  choose: push and whitebox-centric architecture
*  avoids duplicating Boolean status checks when a metric can provide information on both state and performance.
*  metrics provide a dynamic, real-time picture of the state of your infrastructure that will help you manage and make good decisions about your environment.
*  metrics have the potential to identify faults or issues before they occur or before the specific system event that indicates an outage is generated.


## Metric

*  value
*  timestamp
*  other properties

combination of these data point observations is called a time series.

A classic example of a metric we might collect as a time series is website visits, or hits. We periodically collect observations about our website hits, recording the number of hits and the times of the observations. We might also collect properties such as the source of a hit, which server was hit, or a variety of other information.

granularity or resolution. This could range from one second to five minutes to 60 minutes or more.

example: snapshot: location 513


### Gauges

Gauges are numbers that are expected to change over time. A gauge is essentially a snapshot of a specific measurement.

examples: CPU, memory, disk usage

### Counters

*  numbers that increase over time and never decrease.
*  sometimes reset to zero and start incrementing again.
*  Metrics sent by the client increment or decrement the value of the gauge rather than giving its current value
examples:
*  number of bytes sent and received by a device,
*  uptime

### Timers

*   a measure of the number of milliseconds elapsed between a start and end time


*  Often the value of a single metric isn't useful to us. Instead, visualization of a metric requires applying mathematical transformations to it.


### aggregation

*  you often want to show aggregated views of metrics from multiple sources, such as disk space usage of all your application servers.


## Notifications

*  The simplest advice we can give here is to remember notifications are read by humans, not computers. Design them accordingly.

## Visualization

[Metric graphs 101: Timeseries graphs](https://www.datadoghq.com/blog/timeseries-metric-graphs-101/)


## New vs Traditional methods

*  放在最后思
*  Traditional monitoring is heavily focused on this active polling of objects , 通过布尔值来监控系统
*  old: usually centric to single hosts or services.
*  old scale困难: poll中央节点，挂了就挂了; decouple by msg
*  中央节点需要记忆: 交作业，老师收 vs 学生交
*  new:
    -  The event router in our monitoring framework is responsible for tracking our events and metrics.
    -  long-term analysisi of trend, 

## 例子


智能家居
*  case, sensor，温度，是否打开
*  入侵检测，告警, notification
