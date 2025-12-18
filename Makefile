local:
	sudo docker compose -f docker-compose.local.yaml up -d --build

local-stop:
	sudo docker compose -f docker-compose.local.yaml down --remove-orphans

reload:
	sudo docker compose -f docker-compose.local.yaml down --remove-orphans && sudo docker compose -f docker-compose.local.yaml up -d --build

log:
	sudo docker logs constructions_service

DEFAULT_GOAL := local