# Deployment Guide

This document describes how to deploy the Campus Lost and Found application to different environments.

## Docker Deployment

### Building the Docker Image

1. Build the Docker image:
```bash
docker build -t campus-lostandfound:latest .
```

2. Run the container:
```bash
docker run -d -p 3000:3000 \
  -e DATABASE_URL="your-postgresql-connection-string" \
  -e JWT_SECRET="your-jwt-secret" \
  --name campus-lostandfound campus-lostandfound:latest
```

## Deploying to Fly.io

1. Install the Fly CLI:
```bash
curl -L https://fly.io/install.sh | sh
```

2. Log in to Fly:
```bash
fly auth login
```

3. Launch the application:
```bash
fly launch
```

4. Set secrets:
```bash
fly secrets set DATABASE_URL="your-postgresql-connection-string"
fly secrets set JWT_SECRET="your-jwt-secret"
```

5. Deploy the application:
```bash
fly deploy
```

## Deploying to Heroku

1. Install the Heroku CLI:
```bash
curl https://cli-assets.heroku.com/install.sh | sh
```

2. Log in to Heroku:
```bash
heroku login
```

3. Create a new Heroku app:
```bash
heroku create campus-lostandfound
```

4. Set environment variables:
```bash
heroku config:set DATABASE_URL="your-postgresql-connection-string"
heroku config:set JWT_SECRET="your-jwt-secret"
```

5. Deploy the application:
```bash
git push heroku main
```

## Deploying to Railway

1. Install the Railway CLI:
```bash
npm i -g @railway/cli
```

2. Log in to Railway:
```bash
railway login
```

3. Initialize your project:
```bash
railway init
```

4. Set up your environment variables in the Railway dashboard

5. Deploy to Railway:
```bash
railway up
```

## Database Setup

### Using Neon Tech PostgreSQL

1. Create a free account at [Neon Tech](https://neon.tech/)
2. Create a new PostgreSQL database
3. Get your connection string from the dashboard
4. Set the `DATABASE_URL` environment variable with your connection string

### Running Migrations

The application will automatically run migrations on startup. However, if you want to run migrations manually:

```bash
export DATABASE_URL="your-postgresql-connection-string"
go run cmd/migrate/main.go
```
