#!/bin/bash

set -o errexit

main() {
  echo "CREATEING USER"
  create_user

  echo "CREATING DATABASE"
  create_databases

}

create_user() {
  mongo --port 27017 -u $MONGO_INITDB_ROOT_USERNAME -p $MONGO_INITDB_ROOT_PASSWORD --authenticationDatabase admin <<EOF
     use admin;
     db.createUser(
        {
          user: "$MONGO_USERNAME",
          pwd: "$MONGO_PASSWORD",
           roles: [ 
              { role: "userAdminAnyDatabase", db: "admin" }, 
              { role: "dbAdminAnyDatabase", db: "admin" }, 
              { role: "readWriteAnyDatabase", db: "admin" } 
            ]
        }
     );
EOF
}

create_databases() {
  mongo --port 27017 -u $MONGO_INITDB_ROOT_USERNAME -p $MONGO_INITDB_ROOT_PASSWORD --authenticationDatabase admin <<EOF
     use $MONGO_INITDB_DATABASE;
     db.test.insertOne( { hello: "world" } );
EOF
}


main "$@"
