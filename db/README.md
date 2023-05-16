# ALfie DB

The database used is PostgreSQL.

## Creating a docker image for production

For production the database requires a certificate key to allow only secure connections. Reason for the dockerfile to get the certificate information and the configuration of the database for production.

## PostgreSQL configuration

Configuration can be found in the file [postgresql.conf](./postgresql.conf). The configuration is set to allow secure connections and to use the certificate key. The configuration for accessing the database is specified in the file [pg_hba.conf](./pg_hba.conf).
