# Prometheus Athena Exporter

Run queries against athena and export the results to prometheus

  - Use SSM for authentication
  - Run queries against multiple different databases

### Running

You can run the prometheus-athena-exporter service either in docker or using the binary.

#### Docker

To run the service using docker first create a config.yaml file and define your configuration then simply run:

```
docker run -v $(PWD)/config.yaml:/mnt/config/config.yaml -p 8081:${INTERNAL_LISTEN_ADDR_PORT} movio/prometheus-athena-exporter:1.0.0
```

#### Binary

To run the service using the binary first create a config.yaml file and define your configuration then simply run:

```
make build && ./go-app -config config.yaml
```

### Configuration

Configuration is passed in via a yaml file for which you specify the path to using the `-config` argument when running the service. An example configuration file can be found [here](config.yaml).

| Name | Description | Default |
| ------ | ------ | ------ |
| listen-address | The address which the endpoints will start up on (healthcheck, hibernate & prometheus | :8081 |
| ssm-prefix | The prefix to use for retrieving access key id/secret from ssm, this is on a per tenant basis so the full keys used are ${ssm-prefix}${tenant}-id and ${ssm-prefix}${tenant}-secret | - |
| aws-region-id | The id of the AWS region where your athena databases are located | - |
| aws-account-id | The id of your AWS account | - |
| tenants.tenant | The name of the tenant, this is used to retrieve the access key id/secret from ssm | - |
| tenants.db-name | The name of the database in athena, this will need to be access by the access keys from ssm (using tenant name & ssm-preifx) | - |
| metrics.name | The name of the metric, spaces will be replaced by `_` and it will be lowercased | - |
| metrics.query | The query to run | - |
| metrics.query-value-columns | The column names to create metrics for, metrics will be created called `${metricName}_${columnName}` | - |
| metrics.query-interval | The interval at which to run the query | - |
| metrics.include-dbs | The names of the databases to run the query on/record the metric for, you can use regex here. If empty it will be run for all configured tenants  | - |
| metrics.exclude-dbs | The names of the databases to specifically not run the query on/record the metric for (blacklist), you can use regex here | - |

### Routes

The following routes are currently exposed:

| Route | Method | Description | Payload |
| ------ | ------ | ------ | ------ |
| /healthcheck | GET | Simple healthcheck to see if the container is alive | - |
| /metrics | GET | Prometheus metrics endpoint which can be scraped by prometheus | - |


### Contributors

  - [Jack Hopner](https://github.com/jackhopner)
