#!/bin/dash

cd build

git clone git@gitlab.grameco.net:Desarrollo/WebRosatel_Encapsulado_V1.git -b $1 $4

cd $4

sed -i 's/8086/'$2'/g' .env
sed -i 's/DB_NNN/DB_NNN'$3'/g' .env
sed -i 's/RB_PORT_BD_HOST=3306/RB_PORT_BD_HOST='$3'/g' .env

docker-compose pull
docker-compose -f docker-compose-qa.yml up -d 

echo '1'
