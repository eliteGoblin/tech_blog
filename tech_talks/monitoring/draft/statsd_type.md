

*  StatsD is used to collect metrics from infrastructure. 
*  It is push-based: clients export metrics to a collection server, which in turn derives aggregate metrics and often drives graphing systems such as Graphite
*  StatsD/Dogstatsd is just an protocol where prometheus umbrella project includes protocol, collection & time series database. Prometheus also has an alert manager. An prometheus equivalent example setup would include StatsD Server (Telegraf StatsD), time series database(InfluxDB) and Kapacitor for alerting.
*  statd可以用APM系统(Application Performance Management)替代
*  The StatsD daemon has been rewritten numerous times in a variety of languages and has client support for a large number of languages and frameworks.
*  Like Graphite (and because it was originally intended to output directly to Graphite), StatsD uses periods . to divide metric namespaces or "buckets." So our metric would be called:
productiona.tornado.production.payments.amount()



[Prometheus vs StatsD for metrics collection](https://medium.com/@yuvarajl/prometheus-vs-statsd-for-metrics-collection-3b107ab1f60d)


[StatsD Metrics Export Specification v0.1](https://github.com/b/statsd_spec)
[bitly/statsdaemon](https://github.com/bitly/statsdaemon)