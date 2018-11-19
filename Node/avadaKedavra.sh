docker rm -f $(docker ps -aq);
docker network prune -f;
docker volume prune -f
