# remindme_bot
Remind me telegram bot

## Testing

### Running API Tests with Podman Compose

You can run comprehensive API tests locally using Podman Compose, which will orchestrate a complete testing environment with 3 containers:

```bash
podman compose up --build
```

**Container Architecture:**
1. **MongoDB** (`mongo`) - Database container for persistent storage testing
2. **API Server** (`api-server`) - The RemindMeBot API built from source
3. **Newman Test Runner** (`newman`) - Postman collection executor that runs the full test suite

The setup automatically:
- Builds the API server from the current source code
- Starts MongoDB with test configuration
- Waits for the API server to be healthy
- Runs the complete Postman test collection (`RemindMeBot_v2.postman_collection.json`)
- Generates JUnit XML test reports in `./test-results/`
- Performs automatic cleanup of test data

**Test Flow:**
1. Creates test user account
2. Tests user management operations (CRUD)
3. Tests reminder management operations (CRUD)
4. Validates API responses and data consistency
5. Cleans up all test resources

The test results will be available in the `test-results` directory for CI integration.
