#!/bin/dash

#cd build
cd $4
docker-compose down
git reset --hard
rm -f themes/rosatel/assets/cache/*
git pull
if [ $5 = "WebRosatel_Encapsulado_V1" ]; then
    sed -i 's/8086/'$2'/g' .env
    sed -i 's/DB_NNN/DB_NNN'$3'/g' .env
    sed -i 's/RB_PORT_BD_HOST=3306/RB_PORT_BD_HOST='$3'/g' .env
elif [ $5 = "Rosatel.com" ]; then
    sed -i 's/9000/'$2'/g' docker-compose.yml
    cp docker-compose.yml docker-compose-qa.yml
fi
docker-compose pull
docker-compose -f docker-compose-qa.yml up -d

echo '1'
