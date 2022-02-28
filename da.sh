echo "Delete images and containers"
docker rmi $(docker images -a -q)
docker rm -f $(docker ps -a -q)