#!/bin/dash

if [ -d $1 ];
then
    cd $1 #Ingresas al directorio
    docker-compose down # Detienes los contenedores
    git reset --hard # Remueves cambios generados 

    cd .. #Regresas un nivel
    rm -rf $1 #Eliminas carpeta

    echo '1'    
else
    echo "No, no existe"
fi

