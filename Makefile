# Todo App Makefile
.PHONY: help build up down logs test clean deploy-prod deploy-test e2e-test

# Default target
help:
	@echo "Available commands:"
	@echo ""
	@echo "🔧 Development:"
	@echo "  dev       - Start development servers"
	@echo "  build     - Build Docker images"
	@echo "  up        - Start all services"
	@echo "  down      - Stop all services"
	@echo "  logs      - View logs"
	@echo "  clean     - Clean up containers and images"
	@echo ""
	@echo "🧪 Testing:"
	@echo "  test      - Run all tests"
	@echo "  test-unit - Run unit tests only"
	@echo "  test-int  - Run integration tests only"
	@echo "  e2e-test  - Run E2E tests"
	@echo ""
	@echo "🚀 Deployment:"
	@echo "  deploy-test - Deploy to test environment"
	@echo "  deploy-prod - Deploy to production"
	@echo "  status      - Check deployment status"

# Build Docker images
build:
	docker-compose build

# Start all services
up:
	docker-compose up -d
	@echo "🚀 Todo App is running!"
	@echo "📝 Frontend: http://localhost:3000"
	@echo "🔧 Backend API: http://localhost:8080/api/todos"

# Stop all services
down:
	docker-compose down

# View logs
logs:
	docker-compose logs -f

# Run all tests
test:
	go test ./... -v

# Run unit tests only
test-unit:
	go test ./test/unit/... -v

# Run integration tests only
test-int:
	go test ./test/integration/... -v

# Run E2E tests
e2e-test:
	@echo "🎭 Running E2E tests..."
	cd tests/e2e && npm install && npx playwright test

# Clean up
clean:
	docker-compose down -v
	docker system prune -f
	docker image prune -f

# Development mode (without Docker)
dev:
	@echo "Starting development servers..."
	@echo "Backend will run on :8080"
	@echo "Frontend will run on :5173"
	@echo "Press Ctrl+C to stop"
	make dev-backend & make dev-frontend

dev-backend:
	go run cmd/server/main.go

dev-frontend:
	cd web && npm run dev

# Production build test
prod-test: build up
	@echo "Waiting for services to start..."
	@sleep 10
	@echo "Testing production build..."
	@curl -f http://localhost:8080/api/todos || (echo "❌ Backend test failed" && exit 1)
	@curl -f http://localhost:3000 || (echo "❌ Frontend test failed" && exit 1)
	@echo "✅ Production build test passed!"

# View container status
status:
	docker-compose ps

# Deploy to test environment (develop branch)
deploy-test:
	@echo "🧪 Deploying to test environment..."
	@if [ "$$(git branch --show-current)" != "develop" ]; then \
		echo "❌ Please switch to develop branch first: git checkout develop"; \
		exit 1; \
	fi
	git add -A
	git commit -m "deploy: test environment deployment" || true
	git push origin develop
	@echo "✅ Test deployment triggered! Check GitHub Actions."

# Deploy to production (main branch)
deploy-prod:
	@echo "🚀 Deploying to production..."
	@if [ "$$(git branch --show-current)" != "main" ]; then \
		echo "❌ Please switch to main branch first: git checkout main"; \
		exit 1; \
	fi
	git merge develop
	git push origin main
	@echo "✅ Production deployment triggered! Check GitHub Actions."

# Full deployment workflow
deploy: test build
	@echo "🔄 Full deployment workflow..."
	@echo "1. Tests passed ✅"
	@echo "2. Images built ✅"
	@echo "3. Ready for deployment!"
	@echo ""
	@echo "Next steps:"
	@echo "  make deploy-test  # Deploy to test first"
	@echo "  make deploy-prod  # Deploy to production"

# AWS EC2 setup
aws-setup:
	@echo "☁️ Setting up AWS EC2..."
	@echo "1. 📋 Follow the guide: docs/AWS_DEPLOYMENT_GUIDE.md"
	@echo "2. 🔑 Configure GitHub secrets"
	@echo "3. 🚀 Run: make deploy-test"

# Local production simulation
local-prod: build
	@echo "🐳 Starting local production simulation..."
	BACKEND_IMAGE=todo-app-backend:latest \
	FRONTEND_IMAGE=todo-app-frontend:latest \
	docker-compose -f docker-compose.prod.yml up -d
	@echo "✅ Local production running on http://localhost"

# Git shortcuts
git-setup:
	@echo "📝 Setting up Git branches..."
	git checkout -b develop 2>/dev/null || git checkout develop
	git push -u origin develop 2>/dev/null || true
	git checkout main
	@echo "✅ Branches ready: main (prod) and develop (test)" 