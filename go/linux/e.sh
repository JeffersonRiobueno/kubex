#!/bin/bash

#cd build

#DESARROLLO 8082 1 123456 CATALAGO_WEB
git clone git@gitlab.grameco.net:Desarrollo/$5.git -b $1 $4

cd $4

if [ $5 = "WebRosatel_Encapsulado_V1" ]; then
    sed -i 's/8086/'$2'/g' .env
    sed -i 's/DB_NNN/DB_NNN'$3'/g' .env
    sed -i 's/RB_PORT_BD_HOST=3306/RB_PORT_BD_HOST='$3'/g' .env
elif [ $5 = "CATALAGO_WEB" ]; then
    sed -i 's/8085/'$2'/g' docker-compose.yml
    cp docker-compose.yml docker-compose-qa.yml
fi
docker-compose pull
docker-compose -f docker-compose-qa.yml up -d 

echo '1'
