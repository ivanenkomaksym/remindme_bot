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

## CI/CD Pipeline

The project uses GitHub Actions for automated testing, building, and deployment with a 3-stage pipeline:

### 🔨 **Build Stage**
- **Go Build & Test**: Compiles code and runs unit tests
- **Triggers on**: Push/PR to `main` branch
- **Go Version**: 1.24.4 with latest patches

### 🧪 **Integration Testing**
- **API Tests**: Full integration tests using Docker Compose
- **Test Stack**: MongoDB + API Server + Newman (Postman CLI)
- **Coverage**: User/Reminder CRUD operations, data validation, cleanup
- **Reports**: JUnit XML results with GitHub integration
- **Artifacts**: Test results uploaded for review

### 🚀 **Image Publishing**
- **Docker Hub Push**: Automated image building and publishing
- **Smart Tagging**: Branch-based, PR-based, SHA-based, and latest tags
- **Conditions**: Only runs after successful build and tests
- **Image**: `yourusername/remindme-bot`

### 📋 **Pipeline Flow**
```
Build → API Tests → Push Image → Deploy (disabled)
  ✓       ✓           ✓          (🚫)
```

**Secrets Required:**
- `DOCKER_USER` & `DOCKER_PASSWORD` - Docker Hub credentials

**Secrets Required for Cloud Deployment:**
- `BOT_TOKEN` - telegram bot token
- `PUBLIC_URL` - used for telegram bot
- `API_KEY` - api key to gain access
- `DB_CONNECTION_STRING` - mongo connection string

**Test Results**: Visible in GitHub Actions UI, PR comments, and downloadable artifacts.
