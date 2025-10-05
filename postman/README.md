# Postman API Tests

This directory contains comprehensive API tests for the RemindMeBot application using Postman collections.

## Files

- `RemindMeBot.postman_collection.json` - Main test collection with setup and cleanup
- `RemindMeBot.local.postman_environment.json` - Local development environment variables
- `RemindMeBot.ci.postman_environment.json` - CI/CD environment variables
- `../scripts/run-postman-tests.sh` - Unix/Linux shell script for CI/CD
- `../scripts/run-postman-tests.ps1` - PowerShell script for Windows CI/CD

## Test Structure

The collection is organized into the following sections:

### 1. Setup
- **Create Test User**: Creates a test user account (ID: 12345) for testing
- **Updates User Location**: Sets user location which is necessary for detecting the timezone

### 2. Users
- **Get All Users**: Retrieves all users and verifies test user exists
- **Get Test User**: Retrieves specific test user data
- **Update User Language**: Tests language preference updates

### 3. Reminders
- **Create Reminder**: Creates a test reminder
- **Get User Reminders**: Lists reminders for test user
- **Get Specific Reminder**: Retrieves individual reminder
- **Update Reminder**: Modifies reminder properties
- **Get Active Reminders**: Lists only active reminders

### 4. Cleanup
- **Delete Test Reminder**: Removes test reminder
- **Delete Test User**: Removes test user account

## Features

### Automatic Cleanup
- Tests automatically clean up created resources
- Cleanup runs even if tests fail
- Uses collection variables to track what needs cleanup
- Skips cleanup if resources weren't created

### Error Handling
- Handles existing users gracefully (409 status codes)
- Validates response structure and data
- Provides meaningful error messages
- Continues tests even if setup finds existing data

### CI/CD Ready
- Environment variable support
- JUnit XML output for CI systems
- Bail on first failure option
- Cross-platform scripts (bash/PowerShell)

## Running Tests Locally

### Using Postman GUI
1. Import `RemindMeBot.postman_collection.json`
2. Import `RemindMeBot.local.postman_environment.json`
3. Update `apiKey` in the environment to match your local API key
4. Ensure your local server is running on `http://localhost:8080`
5. Run the collection

### Using Newman CLI
```bash
# Install newman if needed
npm install -g newman

# Run tests
newman run RemindMeBot.postman_collection.json \
    -e RemindMeBot.local.postman_environment.json \
    --env-var "baseUrl=http://localhost:8080" \
    --env-var "apiKey=your-local-api-key"
```

## Running in CI/CD

### Unix/Linux (GitHub Actions, GitLab CI, etc.)
```yaml
# Example GitHub Actions step
- name: Run API Tests
  run: |
    chmod +x ./scripts/run-postman-tests.sh
    ./scripts/run-postman-tests.sh
  env:
    CI_BASE_URL: ${{ secrets.API_BASE_URL }}
    CI_API_KEY: ${{ secrets.API_KEY }}
```

### Windows PowerShell
```powershell
# Run tests
.\scripts\run-postman-tests.ps1 -BaseUrl "https://your-api-url" -ApiKey "your-api-key"
```

## Environment Variables

### Required
- `CI_BASE_URL` or `baseUrl`: Base URL of the API (e.g., `https://api.example.com`)
- `CI_API_KEY` or `apiKey`: API authentication key

### Collection Variables
These are managed automatically by the tests:
- `userId`: Test user ID (default: 12345)
- `reminderId`: Created reminder ID (set during test execution)
- `testUserCreated`: Flag indicating if test user was created
- `shouldCleanup`: Flag for cleanup behavior
- `testFailed`: Flag indicating test failures
- `reminderMessage`: Reminder message used for assertions

## Test Data

### Test User
- **ID**: 12345
- **Username**: testuser
- **First Name**: Test
- **Last Name**: User
- **Language**: en (updated to es during tests)

### Test Reminder
- **Message**: "Test reminder for cleanup"
- **Active**: true (updated to false during tests)
- **Recurrence**: Daily (updated to Weekly during tests)

## Validation

Each test includes comprehensive validation:
- HTTP status codes
- Response structure
- Data integrity
- Field values
- Cross-request consistency

## Best Practices

1. **Always run cleanup**: The collection includes automatic cleanup
2. **Use unique test data**: Test user ID is high (12345) to avoid conflicts
3. **Handle existing data**: Tests gracefully handle pre-existing users/reminders
4. **Validate thoroughly**: Each request validates both success and data correctness
5. **Environment separation**: Use different environments for local/CI testing

## Troubleshooting

### Common Issues
1. **User already exists (409)**: Normal behavior, tests will continue
2. **API key invalid (401)**: Check environment variable setup
3. **Connection refused**: Ensure API server is running
4. **Cleanup failures**: Check API delete endpoint implementations

### Debug Mode
Add `--verbose` flag to newman command for detailed output:
```bash
newman run collection.json -e environment.json --verbose
```

### Test Reports
JUnit XML reports are generated as `test-results.xml` for CI integration.