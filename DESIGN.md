# DESIGN.md

## Overview

The project is built around a microservices architecture using **Golang** for service logic, **Kubernetes** for deployment and management, **RabbitMQ** for messaging, and **MariaDB** for storage. Most design choices align directly with Grazer’s stack, highlighting adaptability to your established infrastructure.

---

## Core Design Choices

### 1. Microservices Architecture

Microservices were chosen to modularize each functional area, making each part scalable, independently deployable, and easy to maintain. Each microservice serves a single responsibility, and Kubernetes manages them for reliability and load distribution.

**Services included**:

- **User Service**: Handles user creation and manages user profiles.
- **Preferences Service**: Manages and stores user preferences for matching criteria.
- **Liking Service**: Tracks likes between users and publishes these events to RabbitMQ.
- **Matching Service**: Listens for events and creates matches based on mutual likes.

### 2. Kubernetes for Orchestration

Using Kubernetes allows for robust, scalable, and resilient service management, which is essential given Grazer’s reliance on containerization. Kubernetes also allows for easy horizontal scaling, health checks, and management of service dependencies, all of which are ideal for the Grazer stack.

**Alternative Considered**: A simpler Docker Compose setup could have been used for orchestration. However, Kubernetes is far more production-oriented and aligns with Grazer’s infrastructure, which makes it an obvious choice.

### 3. RabbitMQ for Messaging

RabbitMQ handles asynchronous communication between services, allowing the **Liking Service** to publish events and the **Matching Service** to consume them. RabbitMQ decouples services, enabling them to work asynchronously and scale independently. The setup makes RabbitMQ’s queues persistent to handle restarts and provides event-driven matching functionality without blocking the core services.

**Alternative Considered**: Direct HTTP communication between services. While simpler, it would have made services tightly coupled and reduced scalability. RabbitMQ’s event-driven nature is better suited for matching operations in high-traffic environments.

### 4. MariaDB for Persistent Storage

MariaDB is used to store user data, preferences, likes, and matches in separate tables, reflecting Grazer’s preferred database system. The choice aligns with industry standards for structured data storage and is scalable within Kubernetes.

**Alternative Considered**: A NoSQL database could store relationships more flexibly. However, MariaDB fits Grazer’s current stack, making it a straightforward choice that aligns with adaptability goals.

### 5. Golang for Service Logic

Golang was chosen as the core language for each service, aligning with Grazer’s backend technology preferences. Golang’s concurrency model and efficiency make it an ideal language for microservices, and its simplicity keeps the codebase manageable and performant.

**External Libraries**:

- **gorm.io/gorm** and **gorm.io/driver/mysql**: GORM was chosen to simplify database interaction within each service, providing a clear, idiomatic way to handle SQL operations in Golang.
- **github.com/streadway/amqp**: This RabbitMQ client library offers a simple interface for event publishing and consumption, supporting the microservices’ messaging needs.

---

## Conclusion

The overall design is built to display adaptability by directly aligning with Grazer’s preferred stack. Each component is modular, scalable, and resilient, making this setup ideal for a production environment where reliability and scalability are critical. Kubernetes, Golang, RabbitMQ, and MariaDB form a cohesive stack for distributed microservices and make Grazer’s requirements central to every design choice.
