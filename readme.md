# Goshard

This is a experiment with a service to manage database sharding in a way that by using a numeric ID or string UID the user can have different databases on the same application in a way that the application dont do the queries directly to the database, but to the service and the service manages to select the right database, create if not existing and forward the query to it.

Of course a configuration layer is needed so the user can for example create database instances on the cloud and security, but as this is just a experiemnt I didn't bother to code that part.

Also the code was built thinking in postgresql so for a more generalistic approach some refactor would be needed.

# Requirements

Have Go and Postgresql installed 

# Using

The config database schema named as goshardconfig is present on sql/ folder, use it for creating the config database for the service.

Authentication is done by user token on the request params which should match a existing user on goshardconfig.

Endpoints are
- /query: takes as params the query isself, user token which should match the one on the goshardconfig database, shardid or sharduid for identifying the DB and forward the query to the existing DB, also creates the DB if not existing;
- /schema: takes as param user token which should match the one on the goshardconfig database and on the body should be the raw schema text SQL, it will update or create the schema on goshardconfig database for that user to use when creating new DB with shards;
