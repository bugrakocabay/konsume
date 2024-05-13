# Configuration

## Table of Contents

- [Full Config](#full-config)
- [Providers](#providers)
- [Databases](#databases)
- [Queues](#queues)
- [Metrics](#metrics)

---

### Full Config
All the available configuration options are listed below. The configuration file should be in `yaml` format.

| Parameter                                | Description                                                                                                      | Is Required?                        |
|:-----------------------------------------|:-----------------------------------------------------------------------------------------------------------------|:------------------------------------|
| `log`                                    | Format type of logging. Available formats are `text` and `json`                                                  | no (defaults to text)               |
| `debug`                                  | Enable debug logging level                                                                                       | no                                  |
| `providers`                              | List of configuration for queue sources                                                                          | yes                                 |
| `providers.name`                         | Name of the queue source                                                                                         | yes                                 |
| `providers.type`                         | Type of the queue source. Supported types are `rabbitmq`, `kafka` and `activemq`                                 | yes                                 |
| `providers.retry`                        | Amount of times to retry connecting to queue source                                                              | no                                  |
| `providers.amqp-config`                  | Configuration for RabbitMQ                                                                                       | yes (if type is rabbitmq)           |
| `providers.amqp-config.host`             | Host of the RabbitMQ server                                                                                      | yes (if type is rabbitmq)           |
| `providers.amqp-config.port`             | Port of the RabbitMQ server                                                                                      | yes (if type is rabbitmq)           |
| `providers.amqp-config.username`         | Username for the RabbitMQ server                                                                                 | yes (if type is rabbitmq)           |
| `providers.amqp-config.password`         | Password for the RabbitMQ server                                                                                 | yes (if type is rabbitmq)           |
| `providers.kafka-config`                 | Configuration for Kafka                                                                                          | yes (if type is kafka)              |
| `providers.kafka-config.brokers`         | List of Kafka brokers                                                                                            | yes (if type is kafka)              |
| `providers.kafka-config.topic`           | Topic name for Kafka                                                                                             | yes (if type is kafka)              |
| `providers.kafka-config.group`           | Group name for Kafka                                                                                             | yes (if type is kafka)              |
| `providers.stomp-config`                 | Configuration for ActiveMQ                                                                                       | yes (if type is activemq)           |
| `providers.stomp-config.host`            | Host of the ActiveMQ server                                                                                      | yes (if type is activemq)           |
| `providers.stomp-config.port`            | Port of the ActiveMQ server                                                                                      | yes (if type is activemq)           |
| `providers.stomp-config.username`        | Username for the ActiveMQ server                                                                                 | yes (if type is activemq)           |
| `providers.stomp-config.password`        | Password for the ActiveMQ server                                                                                 | yes (if type is activemq)           |
| `databases`                              | List of configuration for databases                                                                              | no                                  |
| `databases.name`                         | Name of the database                                                                                             | yes (if database is used)           |
| `databases.type`                         | Type of the database. Only `postgresql` is supported.                                                            | yes (if database is used)           |
| `databases.connection-string`            | Connection string used to connect to given database                                                              | yes (if database is used)           |
| `databases.retry`                        | Amount of times to retry connecting to database                                                                  | no                                  |
| `queues`                                 | List of configuration for queues                                                                                 | yes                                 |
| `queues.name`                            | Name of the queue                                                                                                | yes                                 |
| `queues.provider`                        | Name of the queue source                                                                                         | yes (should match a provider name ) |
| `queues.retry`                           | Retry mechanism for queue                                                                                        | no                                  |
| `queues.retry.enabled`                   | Flag for enabling/disabling retry mechanism                                                                      | yes (if retry is enabled)           |
| `queues.retry.strategy`                  | Type of the retry mechanism. Supported types are `fixed`, `expo`, and `random`                                   | no (defaults to fixed)              |
| `queues.retry.max-retries`               | Maximum amount of times that retrying will be triggered                                                          | yes (if retry is enabled)           |
| `queues.retry.interval`                  | Amount of time between retries                                                                                   | yes (if retry is enabled)           |
| `queues.retry.threshold-status`          | Minimum HTTP status code to trigger retry mechanism, any status code above or equal this will trigger retrying   | no (defaults to 500)                |
| `queues.routes`                          | List of configuration for routes                                                                                 | yes                                 |
| `queues.routes.name`                     | Name of the route                                                                                                | yes                                 |
| `queues.routes.type`                     | Type of the route.                                                                                               | no (defaults to REST)               |
| `queues.routes.method`                   | HTTP method for the route                                                                                        | no (defaults to POST)               |
| `queues.routes.url`                      | URL for the route                                                                                                | yes                                 |
| `queues.routes.headers`                  | List of headers for the route                                                                                    | no                                  |
| `queues.routes.body`                     | List of key-values to customize body of the request                                                              | no                                  |
| `queues.routes.query`                    | List of key-values to customize query params of the request                                                      | no                                  |
| `queues.routes.timeout`                  | Timeout of the request                                                                                           | no (defaults to 10s)                |
| `queues.routes.database-routes`          | List of configuration for database routes                                                                        | no                                  |
| `queues.routes.database-routes.name`     | Name of the database route                                                                                       | yes (if database route is used)     |
| `queues.routes.database-routes.provider` | Name of the database source used in `databases`                                                                  | yes (if database route is used)     |
| `queues.routes.database-routes.table`    | Name of the table/collection that will be inserted                                                               | yes (if database route is used)     |
| `queues.routes.database-routes.mapping`  | Mapping of the keys in a message to columns/fields in a table/collection                                         | yes (if database route is used)     |
| `metrics`                                | Configuration for Prometheus metrics                                                                             | no                                  |
| `metrics.enabled`                        | Flag for enabling/disabling Prometheus metrics                                                                   | no (defaults to false)              |
| `metrics.port`                           | Port for Prometheus metrics                                                                                      | no (defaults to 8080)               |
| `metrics.path`                           | Path for Prometheus metrics endpoint                                                                             | no (defaults to /metrics)           |
| `metrics.threshold-status`               | Minimum HTTP status code to trigger Prometheus metrics, any status code above or equal this will trigger metrics | no (defaults to 500)                |

---

### Providers

The providers section specifies the external queue sources konsume will connect to, including details like system type, connection credentials, and configurations for messaging systems such as RabbitMQ, Kafka, and ActiveMQ. It is essential for establishing connections to diverse queue sources, enabling efficient message consumption across different platforms.
<br> An example of `providers` section that uses all available queue sources is shown below:
```yaml
providers:
  - name: rabbit-queue
    type: rabbitmq
    retry: 3
    amqp-config:
      host: localhost
      port: 5672
      username: guest
      password: guest
  - name: kafka-queue
    type: kafka
    retry: 3
    kafka-config:
      brokers:
        - localhost:9092
      topic: test
      group: test
  - name: active-queue
    type: activemq
    retry: 3
    stomp-config:
      host: localhost
      port: 61613
      username: admin
      password: admin
```

---

### Databases

The databases section defines the database connections konsume will use to store messages.
<br> An example of `databases` section is shown below:
```yaml
databases:
  - name: "sql-database"
    type: "postgresql"
    connection-string: "postgres://postgres:password@host:5432/dbname?sslmode=disable"
    retry: 3
```

---


### Queues

The queues section specifies the queue sources and their configurations, including retry mechanisms, routes, and database routes. It is essential for defining the queue sources and their respective routes for message consumption.

Under `retry` section, `strategy` can be set to `fixed`, `expo`, or `random`. `fixed` strategy will retry at fixed intervals, `expo` strategy will retry at exponentially increasing intervals, and `random` strategy will retry at random intervals.

An example of `queues` section is shown below:
```yaml
queues:
  - name: "rabbit-queue"
    provider: "rabbit-queue"
    retry:
      enabled: true
      strategy: "fixed"
      max-retries: 3
      interval: "5s"
      threshold-status: 500
    routes:
      - name: "rest-route"
        type: "REST"
        method: "POST"
        url: "http://localhost:8080"
        headers:
          - key: "Content-Type"
            value: "application/json"
        body:
          - key: "message"
            value: "Hello, World!"
        query:
          - key: "id"
            value: "123"
        timeout: "10s"
      - name: "database-route"
        type: "database"
        database-routes:
          - name: "sql-database-route"
            provider: "sql-database"
            table: "some_table"
            mapping:
              name: "client_name"
              age: "client_age"
```

You can also use <b>GraphQL</b> as a route type. An example of `routes` section with GraphQL route is shown below:
```yaml
routes:
  - name: 'test-route'
    method: 'POST'
    type: 'graphql'
    headers:
      Content-Type: 'application/json'
    body:
      mutation: |
        mutation {
          addUser(name: {{name1}}, email: {{email1}}) {
            id
            name
            email
          }
        }
    url: 'http://someurl:4000/graphql'
```

---

### Metrics

The metrics section specifies the configuration for Prometheus metrics, including enabling/disabling metrics, port, path, and threshold status for triggering metrics.

An example of `metrics` section is shown below:
```yaml
metrics:
  enabled: true
  port: 8080
  path: /metrics
  threshold-status: 500
```