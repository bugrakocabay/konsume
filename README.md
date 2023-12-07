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
  <a href="https://github.com/bugrakocabay">
    <img src="https://img.shields.io/github/followers/bugrakocabay?label=Follow&style=social" alt="Follow" />
  </a>
</p>


## Overview
**TLDR;** If you want to consume messages from RabbitMQ or Kafka and perform HTTP requests based on configurations, konsume is for you.

komsume is a tool that easily connects message queues like RabbitMQ and Kafka with web services, automating data-driven HTTP requests. It bridges complex messaging systems and web APIs, enabling you to create workflows where queue messages automatically trigger web requests. Its flexible setup, including various retry options and customizable request formats, suits a range of uses, from basic data transfers to intricate processing tasks. 

### Features
- **Message Consumption**: Efficiently consumes messages from specified queues.
- **Dynamic HTTP Requests**: Sends HTTP requests based on message content and predefined configurations.
- **Retry Strategies**: Supports fixed, exponential, and random retry strategies for handling request failures.
- **Request Body Templating**: Dynamically constructs request bodies using templates with values extracted from incoming messages.
- **Custom HTTP Headers**: Allows setting custom HTTP headers for outgoing requests.
- **Configurable via YAML**: Easy configuration using a YAML file for defining queues, routes, and behaviors.

### Installation
Easiest way to install konsume is to run via Docker. konsume will look for a configuration file named `config.yaml` in the `/config` directory. Or you can set the path of the configuration file using the `KONSUME_CONFIG_PATH` environment variable.
```bash
docker run -d --name konsume -v /path/to/config.yaml:/config/config.yaml bugrakocabay/konsume:latest
```

Alternatively, you can download the latest binary with the go installer:
```bash
go install github.com/bugrakocabay/konsume@latest
```

### Usage
konsume depends on a YAML configuration file for defining queues, routes, and behaviors. There are two main sections in the configuration file: `providers` and `queues`. In the `providers` section, you can define the message queue providers that konsume will use to consume messages. In the `queues` section, you can define the queues that konsume will consume messages from and the routes that konsume will use to send HTTP requests.

**A simple usage for konsume with RabbitMQ:**
```yaml
providers:
  - name: "rabbit-queue"
    type: "rabbitmq"
    amqp-config:
      host: "localhost"
      port: 5672
      username: "user"
      password: "password"
queues:
  - name: "queue-for-rabbit"
    provider: "rabbit-queue"
    routes:
      - name: "ServiceA_Queue2"
        type: "REST"
        method: "POST"
        url: "https://someurl.com"
```

**A simple usage for konsume with Kafka:**
```yaml
providers:
  - name: "kafka-queue"
    type: "kafka"
    kafka-config:
      brokers:
        - "localhost:9092"
      topic: "your_topic_name"
      group: "group1"
queues:      
  - name: "queue-for-kafka"
    provider: "kafka-queue"
    routes:
      - name: "ServiceA_Queue2"
        type: "REST"
        method: "POST"
        url: "https://someurl.com"
```
