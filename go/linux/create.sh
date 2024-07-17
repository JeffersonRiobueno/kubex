#!/bin/bash

#cd build

#DESARROLLO 8082 1 123456 CATALAGO_WEB
git clone $6/$5.git -b $1 $4

cd $4

#    sed -i 's/8086/'$2'/g' .env #CONFIG .env
#    sed -i 's/DB_NNN/DB_NNN'$3'/g' .env #CONFIG .env
#    sed -i 's/RB_PORT_BD_HOST=3306/RB_PORT_BD_HOST='$3'/g' .env #CONFIG .env

docker-compose pull
docker-compose -f docker-compose.yml up -d 

echo '1'
