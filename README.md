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

- **Backend**: Go, Gin Web Framework, MongoDB, Docker, Logstash, Elasticsearch, Kibana (ELK Stack)
- **CI/CD**: Jenkins, GitHub Actions
- **Frontend**: React, TypeScript, React Query, Storybook

## How to Run
```bash
git clone https://github.com/yourusername/restaurant-management.git
cd restaurant-management
docker-compose up --build
```
