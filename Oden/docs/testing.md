# Oden - Testing Guide

This document outlines the testing procedures for the Oden Idle RPG MVP.

## Server Testing

### Prerequisites

- Go 1.18 or newer
- MySQL installed locally or accessible test instance
- Postman or similar API testing tool

### Setting Up Test Environment

1. Create a test database:
   ```bash
   mysql -u root -p
   CREATE DATABASE oden_test;
   ```

2. Configure test environment:
   - Copy `server/internal/config/config.example.json` to `server/internal/config/config.test.json`
   - Update the database connection to use your test database

3. Run database migrations:
   ```bash
   cd server
   go run cmd/migrations/main.go -config=internal/config/config.test.json
   ```

### Running Unit Tests

```bash
cd server
go test ./...
```

This command runs all tests in the project. For specific packages:

```bash
go test ./internal/auth
go test ./internal/game
```

### API Testing

1. Start the server in test mode:
   ```bash
   cd server
   go run cmd/api/main.go -config=internal/config/config.test.json
   ```

2. Using Postman or cURL:

   - Test user registration:
     ```bash
     curl -X POST http://localhost:8080/v1/auth/register \
       -H "Content-Type: application/json" \
       -d '{"username":"testuser","email":"test@example.com","password":"password123"}'
     ```

   - Test user login:
     ```bash
     curl -X POST http://localhost:8080/v1/auth/login \
       -H "Content-Type: application/json" \
       -d '{"username":"testuser","password":"password123"}'
     ```
     
     Save the returned token for subsequent requests.

   - Test hero summoning:
     ```bash
     curl -X POST http://localhost:8080/v1/heroes/summon \
       -H "Content-Type: application/json" \
       -H "Authorization: Bearer <your_token>" \
       -d '{"summon_type":"basic"}'
     ```

   - Test team formation:
     ```bash
     curl -X POST http://localhost:8080/v1/team/save \
       -H "Content-Type: application/json" \
       -H "Authorization: Bearer <your_token>" \
       -d '{"positions":{"1":"hero_id_1","2":"hero_id_2","3":null,"4":null,"5":null}}'
     ```

   - Test battle:
     ```bash
     curl -X POST http://localhost:8080/v1/battle/start \
       -H "Content-Type: application/json" \
       -H "Authorization: Bearer <your_token>" \
       -d '{"stage_id":"stage_001"}'
     ```

   - Test idle rewards:
     ```bash
     curl -X GET http://localhost:8080/v1/idle/rewards \
       -H "Authorization: Bearer <your_token>"
     ```

### Load Testing (Optional)

For basic load testing, you can use tools like Apache Bench or wrk:

```bash
# Install Apache Bench
sudo apt-get install apache2-utils

# Run a simple load test
ab -n 1000 -c 50 -H "Authorization: Bearer <your_token>" http://localhost:8080/v1/heroes/list
```

## Client Testing

### Prerequisites

- Unity 2021.3 LTS or newer
- Unity Test Framework package installed
- Mock server or test server running

### Unity Unit Tests

1. Open the Unity project from the `client` directory

2. Access the Test Runner:
   - Go to `Window > General > Test Runner`
   - A new panel will open showing test suites

3. Run tests:
   - Select `PlayMode` or `EditMode` tests
   - Click `Run All` to execute all tests

### Integration Testing

1. Configure the client to use the test server:
   - Open `Assets/Scripts/Config/GameConfig.cs`
   - Set `ApiEndpoint` to your test server URL (e.g., `http://localhost:8080/v1`)

2. Test the login flow:
   - Run the game in Unity editor
   - Try to register a new user
   - Try to login with existing credentials
   - Verify success and error scenarios

3. Test hero summoning:
   - Navigate to the summoning screen
   - Perform a hero summon
   - Verify that new heroes appear in your collection

4. Test team formation:
   - Go to the team formation screen
   - Add heroes to different positions
   - Save the team and verify persistence
   - Try invalid formations to test error handling

5. Test battle system:
   - Start a battle with your team
   - Verify battle visualization works correctly
   - Check that battle results are displayed properly
   - Verify rewards are added to your account

6. Test idle rewards:
   - Close the game
   - Wait for a few minutes
   - Reopen the game
   - Verify idle rewards are calculated and offered

### UI Testing

1. Run through all main screens and verify UI elements are displayed correctly

2. Test responsive layouts:
   - In Unity Editor, use Game view to simulate different device resolutions
   - Verify UI adapts properly to different screen sizes

3. Test UI animations and transitions

4. Verify touch interactions work as expected

### End-to-End Testing

1. Build the game for iOS:
   - Go to `File > Build Settings`
   - Build for iOS
   - Run on a device or simulator

2. Test the complete flow:
   - Register a new account
   - Summon heroes
   - Form a team
   - Battle enemies
   - Collect rewards
   - Close and reopen to test idle rewards

## Common Testing Issues

### Server Testing

- Database connection errors: Verify your test database is running and accessible
- Authentication failures: Check JWT token expiration and validation
- Missing dependencies: Ensure all Go modules are installed

### Client Testing

- API connection issues: Verify the server is running and reachable
- Authentication issues: Check if the token is being stored and sent correctly
- Asset loading problems: Verify all required assets are included in the build

## Test Automation

For more robust testing, consider setting up:

1. CI/CD pipeline with automated tests
2. Automated UI testing with tools like Unity Test Framework's PlayMode tests
3. API test automation with Postman collections or similar
4. Performance profiling for both client and server

## Bug Reporting

When reporting bugs, include:

1. Steps to reproduce
2. Expected behavior
3. Actual behavior
4. Environment details (OS, device, build version)
5. Screenshots or videos if applicable 