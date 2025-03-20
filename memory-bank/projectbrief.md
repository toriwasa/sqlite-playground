# Project Brief

## Overview
This project is a base repository for creating CLI tools in Golang with SQLite integration. It follows a clean architecture pattern with a clear separation of concerns between domain models, infrastructure, use cases, controllers, and handlers.

## Core Requirements
- Implement a Go application following clean architecture principles
- Utilize SQLite as the database
- Provide both CUI (Command-line User Interface) and API interfaces
- Follow strict code conventions and architectural guidelines as defined in .clinerules

## Project Structure
The project follows a specific directory structure:
```
.
├ main.go  // Application entry point
└ internal
     ├ domain
     │    └ models  // Domain model structures
     ├ infrastructures
     │    ├ file  // File I/O operations
     │    ├ db    // Database operations
     │    └ api   // API operations
     ├ usecase  // Business logic
     ├ controller  // Controllers that call business logic
     └ handler
          ├ cui  // Command-line interface handlers
          └ api  // Web API handlers
```

## Goals
- Create a reusable template for Golang CLI applications
- Implement clean architecture patterns for maintainable code
- Provide a solid foundation for SQLite-based applications
- Support both command-line and API interfaces
