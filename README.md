# Actor Payment System

![Actor System Architecture](./assets/actor_model.jpg)

This repository contains an implementation of a payment system using the Actor model architecture. The above diagram illustrates the system's components and their interactions.

## Prerequisites

- PostgreSQL database (Make sure you have PostgreSQL installed and running on your system)

## Running the Application

To run the program, follow these steps:

1. Open two terminal windows
2. In the first terminal, run:
    ```bash
    go run broker/main.go
    ```
3. In the second terminal, run:
    ```bash
    go run cmd/processor/main.go
    ```

    ## Quick Start

    To run the application using Docker:

    ```bash
    make dev
    ```

    This command will:
    1. Build a Docker image
    2. Start a container with all necessary services