## Using konsume with docker-compose

This example demonstrates how to use Konsume with Docker Compose, integrating RabbitMQ and Kafka as message sources.

### Overview
This Docker Compose setup runs Konsume alongside RabbitMQ and Kafka. It shows how Konsume can be configured to consume messages from both a RabbitMQ queue and a Kafka topic, and then perform HTTP requests based on the message content.

### Services
- Konsume: The main application, which consumes messages from RabbitMQ and Kafka.
- RabbitMQ: A message broker that stores messages from publishers and delivers them to Konsume.
- Kafka: A distributed streaming platform that Konsume uses to consume messages.
- Zookeeper: A service required by Kafka for distributed coordination.

#### Configuration
Konsume Configuration (config.yaml)
This file contains the configuration for Konsume, specifying the details for connecting to RabbitMQ and Kafka, queue settings, and routing information.

#### Key Configuration Points:
Debug Mode: Enabled for verbose logging.
Providers: RabbitMQ and Kafka with their respective configurations.
Queues: Defined for both RabbitMQ and Kafka, including retry strategies and routing details.
Routes: Configured to send HTTP requests, with one route demonstrating REST and another using GraphQL.
Docker Compose File (docker-compose.yaml)
This file defines the setup for running Konsume, RabbitMQ, Kafka, and Zookeeper using Docker Compose.

#### Key Components:
Service Dependencies: Ensures Konsume starts after RabbitMQ and Kafka are healthy.
Health Checks: For RabbitMQ and Kafka to ensure they are running properly.
Port Mappings: Expose RabbitMQ and Kafka ports for local access.
Environment Variables: Configuration for RabbitMQ and Kafka.

#### Running the Example
Start Services: Run docker-compose up to start all services defined in docker-compose.yaml.
Access RabbitMQ and Kafka: Use the exposed ports to interact with RabbitMQ and Kafka.
Monitor Konsume: Check the logs of the Konsume container for its activity and debug information.

#### Customization
You can modify config.yaml to change the behavior of Konsume, such as adding new routes, changing retry settings, or adjusting the message queue configurations.