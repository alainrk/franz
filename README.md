# Franz

**Work in progress**

## Overview

Franz is a Work In Progress command-line tool designed to facilitate the analysis of Kafka topics, perform scraping operations, and execute various tasks related to Kafka message processing. At the moment it includes simple functionalities such as starting a Kafka consumer, publishing random messages to topics, and cleaning up topics as needed.

This tool is built using the Confluent Kafka library in Go and provides an interface for developers and administrators to interact with Kafka topics. With Franz, you will be able to analyze, produce, and manage messages in your Kafka clusters.

## Development Setup

Before running the Franz, ensure you have Docker and Docker Compose installed on your system. The following make commands are provided to manage the setup and execution of the Franz:

- `make up`: Starts Kafka broker, Zookeper, Schema Registry, Kafka UI through Docker Compose.
- `make down`: Stops and removes the Docker containers.
- `make logs`: Displays the logs of the Franz services.
- `make run`: Executes the Franz (consumer) using the main entry point.
- `make publish-random`: Provides instructions for publishing random messages to a topic.
- `make cleanup-topics`: Displays a warning and instructions to clean up topics.

## Usage

TODO