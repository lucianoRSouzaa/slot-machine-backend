# Slot Machine

A simple and extensible slot machine game implemented in Go, following Clean Architecture principles. This project demonstrates the use of repositories, use case patterns, and comprehensive testing to ensure robustness and maintainability.

## Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Architecture](#architecture)
- [Technologies Used](#technologies-used)
- [Project Structure](#project-structure)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
- [Usage](#usage)
- [Running Tests](#running-tests)

## Overview

The Slot Machine project is a backend application that simulates a slot machine game. It manages players, slot machines, and the interactions between them. The application is designed with scalability and testability in mind, making it easy to extend and maintain.

## Features

- **Player Management**: Create and manage player accounts with balance tracking.
- **Slot Machine Management**: Create and manage slot machines with customizable permutations and balance.
- **Gameplay**: Players can place bets on slot machines, with outcomes determining wins or losses.
- **Comprehensive Testing**: Includes unit tests covering various gameplay scenarios to ensure reliability.
- **Clean Architecture**: Follows Clean Architecture principles for separation of concerns and maintainability.

## Architecture

The project follows the Clean Architecture pattern, ensuring a clear separation between different parts of the application. Here's a high-level overview:

- **Domain Layer (`internal/domain`)**
  - **Model**: Defines the core business entities like `Player` and `SlotMachine`.
  - **Repository Interfaces**: Defines interfaces for player and slot machine repositories.

- **Infrastructure Layer (`internal/infrastructure`)**
  - **Repositories**: Implements the repository interfaces.

- **Application Layer (`internal/application`)**
  - **Use Cases**: Implements the business logic, such as the `PlayUseCase` for handling gameplay.

- **Tests (`internal/application/usecase`)**
  - **Unit Tests**: Comprehensive tests for use cases using in-memory repositories and the `testify` library.

## Technologies Used

- **Go (Golang)**: The primary programming language used for implementing the application.
- **Testify**: A toolkit with common assertions and mocks for testing in Go.

## Project Structure

```plaintext
slot-machine/
├── internal/
│   ├── application/
│   │   └── usecase/
│   │       ├── play_usecase.go
│   │       └── play_usecase_test.go
│   ├── domain/
│   │   ├── model/
│   │   │   ├── player.go
│   │   │   └── slot_machine.go
│   │   └── repository/
│   │       ├── player_repository.go
│   │       └── slot_machine_repository.go
│   └── infrastructure/
│       └── repository/
│           └── in_memory/
│               ├── player_repository.go
│               └── slot_machine_repository.go
├── go.mod
├── go.sum
└── README.md
```

## Getting Started

### Prerequisites

- **Go**: Make sure you have Go installed. You can download it from [golang.org](https://golang.org/dl/).
- **Git**: To clone the repository.

### Installation

1. **Clone the Repository**

```bash
   git clone https://github.com/yourusername/slot-machine.git
   cd slot-machine
```

2. **Initialize Go Modules**

Ensure you're in the project directory and initialize Go modules:

```bash
  go mod tidy
```

## Usage

### Starting the Server

To start the Slot Machine server, follow these steps:

1. **Navigate to the Project Directory**

   Open your terminal and navigate to the root directory of the project:

```bash
   cd slot-machine
```

2. **Run the Application**

Execute the following command to start the server:

```bash
  go run cmd/server/main.go
```

or

```bash
make run
```

### API Documentation (Swagger)

The Slot Machine project includes Swagger for interactive API documentation. Swagger allows you to explore and test the API endpoints directly from your browser.

1. **Accessing Swagger UI**

Once the server is running, open your web browser and navigate to:

```bash
http://localhost:PORT/swagger/index.html
```

This will load the Swagger UI, displaying all available API endpoints with detailed information about each.

## Running Tests

The project includes comprehensive unit tests for the PlayUseCase, ensuring that all gameplay scenarios are handled correctly.

### Execute All Tests

To run all tests, navigate to the project directory and execute:

```bash
  go test ./...
```

or

```bash
  make test
```
