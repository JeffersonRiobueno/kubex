#!/bin/dash

cd build
cd $4
docker-compose down
git reset --hard
git pull

sed -i 's/8086/'$2'/g' .env
sed -i 's/DB_NNN/DB_NNN'$3'/g' .env
sed -i 's/RB_PORT_BD_HOST=3306/RB_PORT_BD_HOST='$3'/g' .env

#docker-compose up -d

echo '1'