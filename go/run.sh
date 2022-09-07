#!/bin/dash

cd build


cd $1

docker-compose down && docker-compose up -d

#echo '1'