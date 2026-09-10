VERSION ?=

check-version:
	@test -n "$(VERSION)" || (echo "VERSION is required. Example: make migrate-force VERSION=1" && exit 1)

setup: build
	docker compose up -d postgres redis
	$(MAKE) migrate-up
	$(MAKE) seed
	docker compose up -d cleaning-app

deploy: build
	docker compose up -d postgres redis
	$(MAKE) migrate-up
	docker compose up -d cleaning-app

migrate-up:
	docker compose run --rm cleaning-app migrate up

migrate-down:
	docker compose run --rm cleaning-app migrate down

migrate-force: check-version
	docker compose run --rm cleaning-app migrate force -- $(VERSION)

migrate-version:
	docker compose run --rm cleaning-app migrate version

seed:
	docker compose run --rm cleaning-app seed

up:
	docker compose up -d

build:
	docker compose build

down:
	docker compose down