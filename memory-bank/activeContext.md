# Active Context

## Current Work Focus
The project is currently in its initial setup phase with domain models for stock prices being defined. The basic directory structure has been established according to clean architecture principles, but most components are still to be implemented.

## Recent Changes
- Created the basic project structure
- Set up the main.go entry point
- Implemented a basic CUI handler with flag parsing
- Added logging configuration
- Defined domain models for stock prices (StockPrice, DailyStockPrice, StockPriceStatistics)

## Current State
- The project is a template repository for Golang CLI applications with SQLite
- The basic structure and CUI handler are implemented
- Domain models for stock prices have been defined:
  - StockPrice: Basic stock price information with stock ID and price
  - DailyStockPrice: Stock price with date information, embedding StockPrice
  - StockPriceStatistics: Statistical information for a collection of stock prices
- No use cases, infrastructure components, or controllers are defined yet
- The CUI handler only handles a verbose flag (-v) for logging control

## Active Decisions
- Following clean architecture principles with clear separation of concerns
- Using SQLite as the embedded database
- Providing both CUI and API interfaces
- Focusing on standard library usage with minimal external dependencies

## Next Steps
1. Define additional domain models in the internal/domain/models package
2. Implement infrastructure components for file, database, and API operations
3. Develop use cases with business logic for stock price operations
4. Create controllers to coordinate between use cases and infrastructure
5. Enhance the CUI handler to process stock-related command-line arguments
6. Implement the API handler for web interface
7. Write unit tests for business logic
8. Add documentation and examples

## Current Challenges
- Ensuring proper separation of concerns between layers
- Maintaining pure functions in the use case layer
- Designing a flexible and intuitive command-line interface
- Planning for extensibility while keeping the code simple
