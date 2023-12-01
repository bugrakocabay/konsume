![konsume](.github/assets/logo.png)

[![Go](https://github.com/bugrakocabay/konsume/actions/workflows/ci.yaml/badge.svg)](https://github.com/bugrakocabay/konsume/actions/workflows/ci.yaml)
[![Go version](https://img.shields.io/github/go-mod/go-version/bugrakocabay/konsume.svg)](https://github.com/bugrakocabay/konsume)
[![Follow](https://img.shields.io/github/followers/bugrakocabay?label=Follow&style=social)](https://github.com/bugrakocabay)

# konsume

konsume is a powerful and flexible tool designed to consume messages from RabbitMQ queues and perform HTTP requests based on configurable integrations. It offers a variety of features including retry strategies, dynamic request body templating, and custom HTTP headers, making it suitable for a wide range of integration scenarios.

## Table of Contents
- [Features](#features)
- [Usage](#usage)
- [Configuration](#configuration)

### Features
- **Message Consumption**: Efficiently consumes messages from specified RabbitMQ queues.
- **Dynamic HTTP Requests**: Sends HTTP requests based on message content and predefined configurations.
- **Retry Strategies**: Supports fixed, exponential, and random retry strategies for handling request failures.
- **Request Body Templating**: Dynamically constructs request bodies using templates with values extracted from incoming messages.
- **Custom HTTP Headers**: Allows setting custom HTTP headers for outgoing requests.
- **Configurable via YAML**: Easy configuration using a YAML file for defining queues, integrations, and behaviors.

### Usage
You need to create a `config.yaml`file to use konsume. It retrieves all the necessary information from this configuration file.


A simple usage for konsume:
```yaml
rabbitmq:
  host: "localhost"
  port: 5672
  username: "user"
  password: "password"

queues:
  - name: "queue1"
    retry:
      enabled: true
      max_retries: 2
      interval: 30s
      strategy: "fixed"
    integrations:
      - name: "my-queue"
        type: "REST"
        method: "POST"
        url: "http://sampleurl.com"
```

### Configuration
| Parameter                                       | Description                                                                                                                                     | Is Required?               |
|:------------------------------------------------|:------------------------------------------------------------------------------------------------------------------------------------------------|:---------------------------|
| `rabbitmq`                                      | RabbitMQ connection configuration                                                                                                               | yes                        |
| `rabbitmq.host`                                 | RabbitMQ host to connect                                                                                                                        | yes                        |
| `rabbitmq.port`                                 | RabbitMQ port to connect                                                                                                                        | yes                        |
| `rabbitmq.username`                             | RabbitMQ username                                                                                                                               | yes                        |
| `rabbitmq.password`                             | RabbitMQ password                                                                                                                               | yes                        |
| `queues`                                        | Configuration for queue that will be consumed                                                                                                   | yes                        |
| `queues.name`                                   | Name of the queue                                                                                                                               | yes                        |
| `queues.retry`                                  | Retry mechanism configuration                                                                                                                   | no                         |
| `queues.retry.enabled`                          | Flag for enabling or disabling retry                                                                                                            | yes (if retry is enabled)  |
| `queues.retry.max_retries`                      | Amount of times that request will be sent, if service returns failure                                                                           | yes (if retry is enabled)  |
| `queues.retry.interval`                         | Duration between retries                                                                                                                        | yes (if retry is enabled)  |
| `queues.retry.strategy`                         | Strategy to use for retry mechanism, three options are available, "fixed", "expo", "random"                                                     | no (defaults to "fixed")   |
| `queues.retry.requeue`                          | Requeueing the message after failure                                                                                                            | no                         |
| `queues.integrations`                           | Services that the message will be forwarded                                                                                                     | yes                        |
| `queues.integrations.name`                      | Name of the integration                                                                                                                         | yes                        |
| `queues.integrations.type`                      | Type of the integration, currently only "REST" is available                                                                                     | yes                        |
| `queues.integrations.method`                    | HTTP method to use for REST integrations, two options are available, "GET", "POST"                                                              | yes (if type is "REST")    |
| `queues.integrations.url`                       | URL to send the request                                                                                                                         | yes                        |
| `queues.integrations.headers`                   | Custom HTTP headers to send with the request                                                                                                    | no                         |
| `queues.integrations.headers.key`               | Key of the header                                                                                                                               | yes (if headers is defined) |
| `queues.integrations.headers.value`             | Value of the header                                                                                                                             | yes (if headers is defined) |
| `queues.integrations.body`                      | Body of the request                                                                                                                             | no                         |

