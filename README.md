# remindme_bot

**Never miss what matters most** - A revolutionary AI-powered Telegram bot that transforms how you manage your busy schedule. Create, schedule, and manage reminders with recurring schedules, timezone detection, and intuitive date/time selection. Simply type "tomorrow 2 PM visit dentist" and watch as advanced natural language processing creates perfectly scheduled reminders with zero effort.

## Overview

In today's fast-paced world, keeping track of important tasks, appointments, and deadlines can be overwhelming. **RemindMeBot solves this challenge** by providing a sophisticated, yet intuitive reminder management system that integrates seamlessly into your daily workflow through Telegram.

Whether you're a busy professional juggling multiple projects, a student managing coursework and deadlines, or anyone who values staying organized and productive, RemindMeBot ensures you **never miss critical moments** in your life. With **breakthrough AI natural language processing**, advanced scheduling capabilities, intelligent recurrence patterns, and global timezone support, this bot adapts to your lifestyle and helps you maintain control over your demanding schedule.

**Key Value Propositions:**
- 🎯 **Stay On Top of Everything**: Transform chaotic schedules into organized, manageable workflows
- ⚡ **Instant Accessibility**: Access your reminders anywhere, anytime through Telegram - no additional apps needed
- 🧠 **AI-Powered Simplicity**: Just say "tomorrow 2 PM visit dentist" or "щодня о 9 ранку приймати ліки" - AI understands and creates perfect reminders
- 🌍 **Global-Ready**: Perfect for remote workers, international teams, and travelers with automatic timezone detection
- 🔄 **Intelligent Scheduling**: Smart recurring patterns reduce manual setup while maximizing productivity
- 💼 **Professional & Personal**: Equally effective for business deadlines and personal commitments

## Key Features

### 🧠 **AI-Powered Natural Language Processing**
- **Smart Text Parsing**: Create reminders using natural language - just type what you want!
- **Multi-Language AI**: Supports English and Ukrainian with OpenAI GPT integration
- **Intelligent Understanding**: Handles complex time expressions like "in 20 minutes", "tomorrow at 2 PM", "every weekday at 8 AM"
- **Context-Aware**: Automatically detects recurrence patterns, dates, times, and reminder messages from plain text

### 🤖 **Telegram Integration**
- Interactive bot interface with intuitive commands
- Real-time notifications delivered directly to Telegram
- Seamless user experience with inline keyboards and quick actions

### ⏰ **Advanced Scheduling**
- **Multiple Recurrence Types**: Once, Daily, Weekly, Monthly, Custom Interval, Spaced-Based Repetition
- **Smart Date Picker**: Interactive calendar for easy date selection
- **Time Picker**: Intuitive time selection interface
- **Timezone Detection**: Automatically detects and adapts to user's timezone

### 🌍 **Localization**
- **Multi-language Support**: English (en) and Ukrainian (uk)
- **Automatic Language Detection**: Adapts to user's Telegram language preferences
- **Localized Interfaces**: Date pickers, messages, and commands in user's language

### 📱 **User Management**
- **Reminder Overview**: View all active and past reminders
- **Easy Deletion**: Remove reminders with simple commands
- **User Preferences**: Language and timezone customization
- **Persistent Storage**: Reminders survive bot restarts

### 🔧 **Technical Architecture**
- **Clean Architecture**: Domain-driven design with clear separation of concerns
- **Flexible Storage**: Supports both in-memory and MongoDB persistence
- **Configuration Management**: Environment variables and `.env` file support
- **Background Processing**: Dedicated reminder notifier service

### 🚀 **API Support**
- **Complete REST API**: Full CRUD operations for users and reminders
- **NLP Endpoint**: `POST /api/reminders/{user_id}/from-text` for natural language reminder creation
- **API Authentication**: Secure access with API keys
- **Integration Ready**: Easy integration with external systems
- **Comprehensive Testing**: Automated API tests with Postman collections

## Architecture

### 📁 **Clean Architecture Layers**
```
├── Domain Layer (entities, use cases, repositories)
├── Infrastructure Layer (database, external services)  
├── Application Layer (controllers, middleware)
└── Presentation Layer (Telegram bot, REST API)
```

### 💾 **Storage Options**
- **In-Memory**: Fast, ephemeral storage for development/testing
- **MongoDB**: Persistent, scalable storage for production

### ⚙️ **Configuration**
- **Core Variables**: `BOT_TOKEN`, `API_KEY`, `DB_CONNECTION_STRING`
- **AI Integration**: `OPENAI_API_KEY` for natural language processing
- **`.env` File**: Local development configuration
- **Runtime Settings**: Storage type, server address, notification intervals

### 🔄 **Background Services**
- **Reminder Notifier**: Continuously monitors active reminders
- **Configurable Intervals**: Adjustable check frequency (default: 15 minutes)
- **Reliable Delivery**: Ensures notifications are sent even after failures

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
Build → API Tests → Push Image → Deploy (GCP)
  ✓       ✓           ✓          ✓
```

**Secrets Required:**
- `DOCKER_USER` & `DOCKER_PASSWORD` - Docker Hub credentials

**Secrets Required for Cloud Deployment:**
- `BOT_TOKEN` - telegram bot token
- `PUBLIC_URL` - used for telegram bot
- `API_KEY` - api key to gain access
- `DB_CONNECTION_STRING` - mongo connection string
- `PROJECT_ID` - GCP project id
- `GOOGLE_SERVICE_ACCOUNT_KEY` - GCP service account JSON key
- `OPENAI_API_KEY` - OpenAI key for NLP

**Test Results**: Visible in GitHub Actions UI, PR comments, and downloadable artifacts.
