#!/bin/dash

cd build
cd $1
docker-compose down
git reset --hard

cd ..
rm -rf $1

echo '1'    