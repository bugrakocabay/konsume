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
    <img src="https://github.com/bugrakocabay/konsume/actions/workflows/ci.yaml/badge.svg?branch=dev" alt="CI" />
  </a>
  <a href="https://github.com/bugrakocabay/konsume">
    <img src="https://img.shields.io/github/go-mod/go-version/bugrakocabay/konsume.svg" alt="Go Version" />
  </a>
  <a href="https://github.com/bugrakocabay">
    <img src="https://img.shields.io/github/followers/bugrakocabay?label=Follow&style=social" alt="Follow" />
  </a>
</p>


## Overview
TLDR; If you want to consume messages from RabbitMQ or Kafka and perform HTTP requests based on configurations, konsume is for you.

konsume simplifies the integration of message queues with services, offering a solution for handling real-time data. It's designed to bridge the gap between complex messaging systems like RabbitMQ, Kafka and various web APIs. With konsume, you can easily set up a workflow where messages from these queues trigger specific HTTP requests, allowing for automated, dynamic responses to data as it flows through your system. Its configuration-driven approach, combined with support for various retry strategies and custom request formatting, makes it adaptable for a wide range of scenarios, from simple data forwarding to complex data processing pipelines. This tool is especially valuable in environments where real-time data processing and integration are key, providing a reliable and flexible solution to leverage the power of message queues in a web-driven world.

## Table of Contents
- [Features](#features)
- [Usage](#usage)
- [Configuration](#configuration)

### Features
- **Message Consumption**: Efficiently consumes messages from specified queues.
- **Dynamic HTTP Requests**: Sends HTTP requests based on message content and predefined configurations.
- **Retry Strategies**: Supports fixed, exponential, and random retry strategies for handling request failures.
- **Request Body Templating**: Dynamically constructs request bodies using templates with values extracted from incoming messages.
- **Custom HTTP Headers**: Allows setting custom HTTP headers for outgoing requests.
- **Configurable via YAML**: Easy configuration using a YAML file for defining queues, routes, and behaviors.

### Usage
konsume depends on a YAML configuration file for defining queues, routes, and behaviors.

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
