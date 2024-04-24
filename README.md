# Fairmoney
The Simple Banking Application is designed to facilitate the creation of single-account transactions, both debit and credit, by consumers. These transactions are sent to a third-party provider for processing and are simultaneously stored in the application's internal database. 

# Simple Banking Application Technical Documentation

## Introduction
The Simple Banking Application is designed to facilitate the creation of single-account transactions, both debit and credit, by consumers. These transactions are sent to a third-party provider for processing and are simultaneously stored in the application's internal database. This documentation provides an overview of the application's functionality, technical specifications, and setup instructions.

## Acceptance Criteria
1. Consumers can request debit or credit transaction creation.
2. Transactions are successfully created in the third-party provider system.
3. Transactions are stored in the application’s internal storage.
4. The account balance is updated in the application’s internal storage.

## Technical Requirements
1. The application must provide synchronous handling of debit or credit transaction requests.
2. Code must be thoroughly tested with a sufficient amount of unit tests.
3. Third-party provider API endpoints must be integrated for transaction creation and retrieval.
4. The application must be available via a RESTful JSON API.
5. Both successful scenarios (happy path) and potential edge cases must be handled effectively.

## Third-party Provider API Specification
### POST /third-party/payments
- **Request Body:**
  - `account_id`: Identifier of the account involved in the transaction.
  - `reference`: Unique reference for the transaction.
  - `amount`: Transaction amount.
- **Response Body:**
  - `account_id`: Identifier of the account involved in the transaction.
  - `reference`: Unique reference for the transaction.
  - `amount`: Transaction amount.

### GET /third-party/payments/:reference
- **Response Body:**
  - `account_id`: Identifier of the account involved in the transaction.
  - `reference`: Unique reference for the transaction.
  - `amount`: Transaction amount.

## Technical Notes
1. The application is implemented in Go-lang.
2. The choice of technology stack for the internal storage and other components is flexible.
3. The third-party API’s behavior can be mocked for testing purposes.
4. Authorization mechanisms, though not required for simplicity, can be implemented for production use.
5. Pre-seeding of user/account data can be done manually as no API for user/account creation is necessary.
6. Docker can be optionally used for containerization of the application.

## Environment Setup
1. Ensure Go-lang is installed on your development environment.
2. Set up the necessary dependencies and libraries for the chosen technology stack.
3. Configure the application to run locally.
4. Test the application to ensure smooth functionality on your local machine.
5. Optionally, set up Docker for containerizing the application.

## Conclusion
The Simple Banking Application provides a straightforward solution for handling single-account transactions, integrating seamlessly with a third-party provider for transaction processing. By adhering to the provided technical requirements and acceptance criteria, the application ensures robustness, reliability, and ease of use for both consumers and developers.

This documentation serves as a comprehensive guide for understanding and implementing the Simple Banking Application effectively.
