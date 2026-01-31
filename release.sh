#!/bin/bash
set -e

# массив: сервис=версия
services=(
  "constructionsbackend-constructions_service=1.0.1"
  "constructionsbackend-admin=1.0.1"
)

#!/bin/bash

docker tag constructionsbackend-constructions_service:latest pachv01223/constructionsbackend-constructions_service:1.0.1
docker tag constructionsbackend-admin:latest pachv01223/constructionsbackend-admin:1.0.1


docker push pachv01223/constructionsbackend-constructions_service:1.0.1
docker push pachv01223/constructionsbackend-admin:1.0.1


sudo docker pull pachv01223/constructionsbackend-constructions_service:1.0.1
sudo docker pull pachv01223/constructionsbackend-admin:1.0.1




