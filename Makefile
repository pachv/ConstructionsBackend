local:
	sudo docker compose -f docker-compose.local.yaml up -d --build

local-stop:
	sudo docker compose -f docker-compose.local.yaml down --remove-orphans

reload:
	sudo docker compose -f docker-compose.local.yaml down --remove-orphans && sudo docker compose -f docker-compose.local.yaml up -d --build

log:
	sudo docker logs constructions_service

reset-build:
	rm -rf ./build && cp ./build_reset_exaple ./build

DEFAULT_GOAL := local