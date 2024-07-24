# Restaurant Management Application

## Project Overview
This project is a restaurant management application built using the Go programming language and the Gin web framework. The application provides functionalities for managing various aspects of a restaurant, including user management, menu management, table management, and order management. It also incorporates authentication and authorization mechanisms.

## Key Features
- **User Management**: Includes user registration, login, and retrieval of user information.
- **Menu Management**: Allows for the creation, updating, and retrieval of menu items.
- **Table Management**: Facilitates the management of table information within the restaurant.
- **Order Management**: Manages the ordering process, including order item details and invoicing.
- **Authentication and Authorization**: Ensures secure access to the application using JWT tokens.
- **Logging**: Uses Logstash and Logrus for logging events and errors, which are then visualized using the ELK stack.

## Technologies Used

### Go
Go, also known as Golang, is a statically typed, compiled programming language designed for simplicity and efficiency. It is particularly well-suited for developing web services and has excellent support for concurrency.

### Gin Web Framework
Gin is a high-performance web framework for Go that is known for its speed and minimalistic design. It provides a robust set of features for building web applications and APIs, including middleware support, routing, and JSON handling.

### MongoDB
MongoDB is a NoSQL database that stores data in flexible, JSON-like documents. It is highly scalable and designed for high availability and performance.

### Docker
Docker is a platform that enables developers to package applications into containers—standardized units of software that include everything the software needs to run. This project uses Docker to containerize the application and its dependencies, making it easier to develop, deploy, and run in various environments.

### Logstash, Elasticsearch, Kibana (ELK Stack)
- **Logstash**: A server-side data processing pipeline that ingests data from multiple sources simultaneously, transforms it, and then sends it to a "stash" like Elasticsearch.
- **Elasticsearch**: A distributed, RESTful search and analytics engine capable of solving a growing number of use cases.
- **Kibana**: A data visualization dashboard for Elasticsearch, providing insights and analysis of the logged data.

### Jenkins
Jenkins is an open-source automation server that enables developers to build, test, and deploy their software. It supports continuous integration and continuous delivery (CI/CD) practices.

### GitHub Actions
GitHub Actions is a CI/CD tool that automates the build, test, and deployment pipeline directly from GitHub repositories. It provides workflows that can be triggered by events like pushing code or opening a pull request.

## How to Run
```bash
git clone https://github.com/yourusername/restaurant-management.git
cd restaurant-management
docker-compose up --build
```
