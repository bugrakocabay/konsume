<p align="center">
  <a href="https://github.com/bugrakocabay/konsume">
    <img src=".github/assets/logo.png" alt="konsume logo" />
  </a>
</p>

<p align="center">
  konsume is a powerful and flexible tool designed to consume messages from various message queues like RabbitMQ and Kafka and perform HTTP requests based on configurations.
</p>

<p align="center">
  <a href="https://github.com/bugrakocabay/konsume/actions/workflows/ci.yaml">
    <img src="https://github.com/bugrakocabay/konsume/actions/workflows/ci.yaml/badge.svg?branch=main" alt="CI" />
  </a>
  <a href="https://github.com/bugrakocabay/konsume">
    <img src="https://img.shields.io/github/go-mod/go-version/bugrakocabay/konsume.svg" alt="Go Version" />
  </a>
  <a href="https://goreportcard.com/report/github.com/bugrakocabay/konsume">
    <img src="https://goreportcard.com/badge/github.com/bugrakocabay/konsume" alt="Go Report" />
  </a>
  <a href="https://codecov.io/github/bugrakocabay/konsume" > 
    <img src="https://codecov.io/github/bugrakocabay/konsume/graph/badge.svg?token=r36BDXBfXR" alt="Coverage" /> 
  </a>
</p>

## Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Installation](#installation)
- [Usage](#usage)
- [Configuration](#configuration)
- [FAQ](#faq)
- [Contributing](#contributing)

## Overview

**TLDR;** If you want to consume messages from RabbitMQ or Kafka and perform HTTP requests based on configurations, konsume is for you.

komsume is a tool that easily connects message queues like RabbitMQ and Kafka with web services, automating data-driven HTTP requests. It bridges complex messaging systems and web APIs, enabling you to create workflows where queue messages automatically trigger web requests. Its flexible setup, including various retry options and customizable request formats, suits a range of uses, from basic data transfers to intricate processing tasks.

## Features

- **Message Consumption**: Efficiently consumes messages from specified queues.
- **Dynamic HTTP Requests**: Sends HTTP requests based on message content and predefined configurations.
- **Retry Strategies**: Supports fixed, exponential, and random retry strategies for handling request failures.
- **Request Body Templating**: Dynamically constructs request bodies using templates with values extracted from incoming messages.
- **Custom HTTP Headers**: Allows setting custom HTTP headers for outgoing requests.
- **Configurable via YAML**: Easy configuration using a YAML file for defining queues, routes, and behaviors.
- **Monitoring**: Provides a Prometheus endpoint for monitoring metrics.

## Installation

Easiest way to install konsume is to run via Docker. konsume will look for a configuration file named `config.yaml` in the `/config` directory. Or you can set the path of the configuration file using the `KONSUME_CONFIG_PATH` environment variable.

```bash
docker run -d --name konsume -v /path/to/config.yaml:/config/config.yaml bugrakocabay/konsume:latest
```

Alternatively, you can download the latest binary with the go installer:

```bash
go install github.com/bugrakocabay/konsume@latest
```

## Usage

konsume depends on a YAML configuration file for defining queues, routes, and behaviors. There are two main sections in the configuration file: `providers` and `queues`. In the `providers` section, you can define the message queue providers that konsume will use to consume messages. In the `queues` section, you can define the queues that konsume will consume messages from and the routes that konsume will use to send HTTP requests.
You can also find detailed usage examples in the [examples](.examples) directory.

**A simple usage for konsume with RabbitMQ:**

```yaml
providers:
  - name: 'rabbit-queue'
    type: 'rabbitmq'
    amqp-config:
      host: 'localhost'
      port: 5672
      username: 'user'
      password: 'password'
queues:
  - name: 'queue-for-rabbit'
    provider: 'rabbit-queue'
    routes:
      - name: 'ServiceA_Queue2'
        type: 'REST'
        method: 'POST'
        url: 'https://someurl.com'
```

**A simple usage for konsume with Kafka:**

```yaml
providers:
  - name: 'kafka-queue'
    type: 'kafka'
    kafka-config:
      brokers:
        - 'localhost:9092'
      topic: 'your_topic_name'
      group: 'group1'
queues:
  - name: 'queue-for-kafka'
    provider: 'kafka-queue'
    routes:
      - name: 'ServiceA_Queue2'
        type: 'REST'
        method: 'POST'
        url: 'https://someurl.com'
```

## Configuration

| Parameter                         | Description                                                                                                    | Is Required?                            |
|:----------------------------------|:---------------------------------------------------------------------------------------------------------------|:----------------------------------------|
| `debug`                           | Enable debug logging level                                                                                     | no                                      |
| `providers`                       | List of configuration for queue sources                                                                        | yes                                     |
| `providers.name`                  | Name of the queue source                                                                                       | yes                                     |
| `providers.type`                  | Type of the queue source. Supported types are `rabbitmq`, `kafka` and `activemq`                               | yes                                     |
| `providers.retry`                 | Amount of times to retry connecting to queue source                                                            | no                                      |
| `providers.amqp-config`           | Configuration for RabbitMQ                                                                                     | yes (if type is rabbitmq)               |
| `providers.amqp-config.host`      | Host of the RabbitMQ server                                                                                    | yes (if type is rabbitmq)               |
| `providers.amqp-config.port`      | Port of the RabbitMQ server                                                                                    | yes (if type is rabbitmq)               |
| `providers.amqp-config.username`  | Username for the RabbitMQ server                                                                               | yes (if type is rabbitmq)               |
| `providers.amqp-config.password`  | Password for the RabbitMQ server                                                                               | yes (if type is rabbitmq)               |
| `providers.kafka-config`          | Configuration for Kafka                                                                                        | yes (if type is kafka)                  |
| `providers.kafka-config.brokers`  | List of Kafka brokers                                                                                          | yes (if type is kafka)                  |
| `providers.kafka-config.topic`    | Topic name for Kafka                                                                                           | yes (if type is kafka)                  |
| `providers.kafka-config.group`    | Group name for Kafka                                                                                           | yes (if type is kafka)                  |
| `providers.stomp-config`          | Configuration for ActiveMQ                                                                                     | yes (if type is activemq)               |
| `providers.stomp-config.host`     | Host of the ActiveMQ server                                                                                    | yes (if type is activemq)               |
| `providers.stomp-config.port`     | Port of the ActiveMQ server                                                                                    | yes (if type is activemq)               |
| `providers.stomp-config.username` | Username for the ActiveMQ server                                                                               | yes (if type is activemq)               |
| `providers.stomp-config.password` | Password for the ActiveMQ server                                                                               | yes (if type is activemq)               |
| `queues`                          | List of configuration for queues                                                                               | yes                                     |
| `queues.name`                     | Name of the queue                                                                                              | yes                                     |
| `queues.provider`                 | Name of the queue source                                                                                       | yes (should match a provider name )     |
| `queues.retry`                    | Retry mechanism for queue                                                                                      | no                                      |
| `queues.retry.enabled`            | Flag for enabling/disabling retry mechanism                                                                    | yes (if retry is enabled)               |
| `queues.retry.strategy`           | Type of the retry mechanism. Supported types are `fixed`, `expo`, and `random`                                 | no (defaults to fixed)                  |
| `queues.retry.max_retries`        | Maximum amount of times that retrying will be triggered                                                        | yes (if retry is enabled)               |
| `queues.retry.interval`           | Amount of time between retries                                                                                 | yes (if retry is enabled)               |
| `queues.retry.threshold_status`   | Minimum HTTP status code to trigger retry mechanism, any status code above or equal this will trigger retrying | no (defaults to 500)                    |
| `queues.routes`                   | List of configuration for routes                                                                               | yes                                     |
| `queues.routes.name`              | Name of the route                                                                                              | yes                                     |
| `queues.routes.type`              | Type of the route.                                                                                             | no (defaults to REST)                   |
| `queues.routes.method`            | HTTP method for the route                                                                                      | no (defaults to POST)                   |
| `queues.routes.url`               | URL for the route                                                                                              | yes                                     |
| `queues.routes.headers`           | List of headers for the route                                                                                  | no                                      |
| `queues.routes.body`              | List of key-values to customize body of the request                                                            | no                                      |
| `queues.routes.query`             | List of key-values to customize query params of the request                                                    | no                                      |
| `queues.routes.timeout`           | Timeout of the request                                                                                         | no (defaults to 10s)                    |

## FAQ

<details>
<summary> <b>Why konsume?</b> </summary>

Think of konsume as your handy tool for making message queues and web APIs work together like best buddies. It's like having a super-efficient assistant who takes messages from RabbitMQ or Kafka and knows exactly when and how to ping your web services, whether they speak REST or GraphQL. And guess what? If something doesn't go right the first time, konsume keeps trying until it works, thanks to its smart retry strategies. So, whether you're just moving data around or setting up some cool automated workflows, konsume is your go-to for making things simple and reliable.

</details>

<details>
<summary> <b>What message queues does konsume support?</b> </summary>
Currently konsume supports <b>RabbitMQ</b>, <b>Kafka</b> and <b>ActiveMQ</b>. But it is designed to be easily extensible to support other message queues.
</details>

<details>
<summary> <b>How can I dynamically insert values from consumed messages into the request body?</b> </summary>
konsume allows dynamically inserting values from consumed messages into the request body using placeholders. You can use the <code>{{key}}</code> syntax to insert values from consumed messages into the request body. For example, if you have a message like this:

```json
{
	"name": "John",
	"email": "john@doe.com"
}
```

You can use the `{{name}}` and `{{email}}` placeholders in the request body to insert the values from the consumed message into the request body.

```yaml
routes:
  - name: 'test-route'
    method: 'POST'
    type: 'REST'
    headers:
      Content-Type: 'application/json'
    body:
      userName: '{{name}}'
      eMail: '{{email}}'
    url: 'http://someurl.com'
```

</details>

<details>
<summary> <b>Is GraphQL supported?</b> </summary>
Yes! konsume supports GraphQL. You can use the <code>graphql</code> type for routes and define the GraphQL query or mutation in the <code>body</code> section of the route. Under <code>body</code> section, you can use the <code>query</code> or <code>mutation</code> key to define your GraphQL query or mutation. Also konsume allows dynamically inserting values from consumed messages into the GraphQL body using placeholders.

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

</details>

<details>
<summary> <b>How does the retry mechanism work?</b> </summary>
konsume supports three different retry strategies: <code>fixed</code>, <code>expo</code>, and <code>random</code>. You can define the retry strategy in the <code>retry</code> section of the queue configuration. If you want to enable retrying, you should set the <code>enabled</code> flag to <code>true</code>. You can also define the maximum amount of times that retrying will be triggered using the <code>max_retries</code> key. The <code>interval</code> key defines the amount of time between retries. The <code>threshold_status</code> key defines the minimum HTTP status code to trigger retry mechanism, any status code above or equal this will trigger retrying. If you don't define the <code>threshold_status</code> key, it will default to <code>500</code>.

```yaml
queues:
  - name: 'queue-for-rabbit'
    provider: 'rabbit-queue'
    retry:
      enabled: true
      strategy: 'fixed'
      max_retries: 5
      interval: 5s
      threshold_status: 500
    routes:
      - name: 'ServiceA_Queue2'
        type: 'REST'
        method: 'POST'
        url: 'https://someurl.com'
```

</details>

<details>
<summary> <b>How to see the metrics?</b> </summary>
konsume provides a Prometheus endpoint for monitoring metrics. You can see the metrics at <code>/metrics</code> by default. Here you will find a list of metrics that Prometheus can scrape by default.

</details>

<details>
<summary> <b>What are some common troubleshooting steps if konsume is not working as expected?</b> </summary>
<ol>
<li><b>Enable Debug Mode and Examine Logs:</b> Look at the logs for any error messages or warnings. If you've enabled debug mode, this will provide more detailed information.</li>
<li><b>Check Configuration:</b> Ensure your config.yaml is correctly set up for your message queues and routes. Verify all parameters, especially URLs, queue names, and credentials.</li>
<li><b>Validate Queue Connectivity:</b> Make sure Konsume can connect to the message queues. Check network configurations, access permissions, and queue settings.</li>
<li><b>Test HTTP Endpoints:</b> Ensure the endpoints for your HTTP requests are reachable and responding as expected. You can test them independently with tools like Postman or cURL.</li>
<li><b>Review Message Formats:</b> Confirm that the messages in your queues are in the expected format, especially if you're using templating features.</li>
<li><b>Monitor Resource Usage:</b> Sometimes issues arise due to resource constraints. Check CPU, memory, and network usage.</li>
<li><b>Update Konsume:</b> Ensure you're using the latest version of Konsume, as updates might fix known issues.</li>
<li><b>Seek Community Help:</b> If you're still stuck, consider asking for help in <a href="https://github.com/bugrakocabay/konsume/issues">issues</a> or <a href="https://github.com/bugrakocabay/konsume/discussions">discussions</a>.</li>
</ol>
</details>

### Contributing
🌟 Your contributions are welcome!

Whether you're looking to fix bugs, add new features, or improve documentation, your help is greatly appreciated. Contributing to konsume is not only a great way to enhance this tool, but it's also an excellent opportunity to get involved with a community of like-minded individuals.

Here's how you can contribute:

Report Issues: Found a bug or have a suggestion? Open an issue and let us know!
Submit Pull Requests: Have a fix or a new feature? Submit a pull request!
Feedback and Ideas: Share your thoughts and ideas on how we can improve konsume.
Documentation: Help us improve or translate the documentation.
Before contributing, please read our [Contributing Guidelines](./CONTRIBUTING.md) for more information on how to get started.

💡 No contribution is too small – whether it's fixing typos, improving code readability, or updating documentation, all contributions are valuable and appreciated!

Join us in making konsume better for everyone! 🚀
