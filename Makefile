.PHONY: observability-up observability-down logs-clean

observability-up:
	docker-compose -f docker-compose.observability.yml up -d

observability-down:
	docker-compose -f docker-compose.observability.yml down

logs-clean:
	rm -rf ./customer/logs/*
	rm -rf ./order/logs/*

# ... existing targets ...
