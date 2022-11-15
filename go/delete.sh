#!/bin/dash

##cd build
if [ -d $1 ];
then
    cd $1
    docker-compose down
    git reset --hard

    cd ..
    rm -rf $1

    echo '1'    
else
    echo "No, no existe"
fi

