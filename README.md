<p align="center">
  <a href="https://github.com/bugrakocabay/konsume">
    <img src=".github/assets/logo.png" alt="konsume logo" />
  </a>
</p>

<p align="center">
  konsume is a powerful and flexible tool designed to consume messages from various message queues like RabbitMQ and Kafka and perform data-driven actions like HTTP requests and insertions into databases based on predefined configurations.
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
- [FAQ](#faq)
- [Contributing](#contributing)

## Overview

**TLDR;** If you want to consume messages from message queues such as RabbitMQ or Kafka and perform HTTP requests or insert data to a database based on configurations, konsume is for you.

konsume is a tool that easily connects message queues like RabbitMQ and Kafka with web services, automating data-driven actions. It bridges complex messaging systems and web APIs or databases. Its flexible setup, including various retry options and customizable request formats, suits a range of uses, from basic data transfers to intricate processing tasks.

<p align="center">
    <img src=".github/assets/diagram.png" alt="konsume diagram" />
</p>

## Features

- **Message Consumption**: Efficiently consumes messages from specified queues.
- **Dynamic HTTP Requests**: Sends HTTP requests based on message content and predefined configurations.
- **Database Insertions**: Inserts data into databases based on message content and predefined configurations.
- **Retry Strategies**: Supports fixed, exponential, and random retry strategies for handling request failures.
- **Request Body Templating**: Dynamically constructs request bodies using templates with values extracted from incoming messages.
- **Custom HTTP Headers**: Allows setting custom HTTP headers for outgoing requests.
- **Configurable via YAML and JSON**: Easy configuration using a YAML or JSON file for defining queues, routes, and behaviors.
- **Monitoring**: Provides a Prometheus endpoint for monitoring metrics.

## Installation

Easiest way to install konsume is to run via Docker. 

```bash
docker run -d --name konsume -v /path/to/config.yaml:/config/config.yaml bugrakocabay/konsume:latest
```

Alternatively, you can download the latest binary from the [releases](https://github.com/bugrakocabay/konsume/releases) page and run it on your machine.

## Usage

konsume depends on a YAML or JSON configuration file for defining queues, routes, and behaviors. By default, konsume looks for a file named `config.yaml` in the `/config` directory. Alternatively, you can specify the environment variable `KONSUME_CONFIG_PATH` to point to your configuration file.

<br>📜 <b>Checkout the [CONFIGURATION.md](./CONFIGURATION.md) for detailed information about the configuration parameters.
<br>💡 You can also find detailed usage examples in the [examples](.examples) directory. </b>

**A simple usage for konsume to send HTTP requests with RabbitMQ:**

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

**A simple usage for konsume to insert into a database with Kafka:**

```yaml
providers:
  - name: 'kafka-queue'
    type: 'kafka'
    kafka-config:
      brokers:
        - 'localhost:9092'
      topic: 'your_topic_name'
      group: 'group1'
databases:
  - name: "sql-database"
    type: "postgresql"
    connection-string: "postgres://postgres:password@host:5432/dbname"
queues:
  - name: "queue-for-kafka"
    provider: "kafka-queue"
    database-routes:
      - name: "sql-database-route"
        provider: "sql-database"
        table: "some_table"
        mapping:
          userName: "name"
          some-key: "value"
```

## FAQ

<details>
<summary> <b>Why konsume?</b> </summary>

Think of konsume as your handy tool for making message queues and other services work together like best buddies. It's like having a super-efficient assistant who takes messages from RabbitMQ or Kafka and knows exactly when and how to insert data into databases or ping your web services, whether they speak REST or GraphQL. And guess what? If something doesn't go right the first time, konsume keeps trying until it works, thanks to its smart retry strategies. So, whether you're just moving data around or setting up some cool automated workflows, konsume is your go-to for making things simple and reliable.

</details>

<details>
<summary> <b>What message queues does konsume support?</b> </summary>
Currently konsume supports <b>RabbitMQ</b>, <b>Kafka</b> and <b>ActiveMQ</b>. But it is designed to be easily extensible to support other message queues.
</details>

<details>
<summary> <b>What databases does konsume support?</b> </summary>
Currently konsume only supports <b>Postgres</b>. But it is designed to be easily extensible to support other databases.
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
<summary> <b>How can I dynamically map the values inside a message into columns/fields of a database?</b> </summary>
In order to dynamically map the values inside a message into columns/fields of a database, you can use the <code>mapping</code> section in the database route configuration. You can define the mapping between the fields of the message and the columns of the database table. For example, if you have a message like this:

```json
{
	"name": "John",
	"email": "john@doe.com"
}
```

You can use the <code>mapping</code> section to map the <code>name</code> field of the message to the <code>user_name</code> column of the database table and the <code>email</code> field of the message to the <code>user_email</code> column of the database table.

```yaml
queues:
  - name: "queue-for-kafka"
    provider: "kafka-queue"
    database-routes:
      - name: "sql-database-route"
        provider: "sql-database"
        table: "some_table"
        mapping:
          name: "user_name"
          email: "user_email"
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
<br> Also, konsume provides custom metrics for the following events:
<br> - <code>konsume_messages_consumed_total</code>: Total number of messages consumed.
<br> - <code>konsume_http_requests_made_total</code>: Total number of HTTP requests made.
<br> - <code>konsume_http_requests_succeeded_total</code>: Total number of HTTP requests succeeded.
<br> - <code>konsume_http_requests_failed_total</code>: Total number of HTTP requests failed.

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
