.PHONY: all build-backend run-backend run-frontend clean

all: build-backend

build-backend:
	cd backend && go build -o server ./cmd/server

run-backend: build-backend
	cd backend && ./server

run-frontend:
	cd frontend && npm run dev

clean:
	rm -f backend/server
	rm -rf frontend/dist
