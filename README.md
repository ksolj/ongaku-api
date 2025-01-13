# ONGAKU-API

## Short description
Ongaku (i.e. music) API allows users to upload their guitar tabs to the database and edit other people's tabs.

## Offered Features

### API Features
1. **Middleware**:
      - **Metrics**
      - **CORS**
      - **IP-based rate limiting** (Token Bucket Algorithm)
      - **Authentication** (Stateful tokens)
      - **Permission-based authorization** (for certain operations)
      - **Sending requests' data to Kafka**
2. **Caching** (via Redis)
3. **SMTP**:
      - **Welcome message**
      - **User activation**
      - **Password reset**
4. **Pagination and Filtering**
5. **Graceful shutdown**
6. **Versioning**
7. **Panic recovery**

### Other Features and used technologies
1. **Migrations**
2. **Makefile**
3. **Configuration** (Either by .env file or command-line arguments)
4. **Dockerfile (multi-stage build) and Docker Compose (migrations run at startup)**
5. **Caddy as the reverse proxy**
6. **PostgreSQL as the primary database**
7. **Kafka for gathering data**

## QuickStart
1. Install these:
- Docker & Docker Compose
- Make

2. Create .env in the project directory and put desired values in these variables (to run Docker Compose):
```
POSTGRES_DB=
POSTGRES_USER=
POSTGRES_PASSWORD=

REDIS_PASSWORD=
REDIS_USER=
REDIS_USER_PASSWORD=

API_PORT=
API_ENV=
API_SQL_DSN=
API_REDIS_DSN=
API_LIMITER_RPS=
API_LIMITER_BURST=
API_LIMITER_ENABLED=
API_SMTP_HOST=
API_SMTP_PORT=
API_SMTP_USERNAME=
API_SMTP_PASSWORD=
API_SMTP_SENDER=
```

Also feel free to type ```make help``` to see what is has to offer

3. Configure caddy (optional)
4. Run ```docker-compose up```

## Tools
1. SQL Database - PostgreSQL
2. Cache - Redis
3. Kafka
4. Migrations - golang-migrate
5. Container - Docker
6. SMTP server for testing purposes - Papercut SMTP
7. Reverse proxy - Caddy