FROM postgres

# copy cert files to container alongside scripts and config
ADD ./ca/server.crt /home/ca/server.crt
ADD ./ca/server.key /home/ca/server.key
ADD ./init_cert.sh /home/ca/init_cert.sh
ADD ./pg_hba.conf /etc/postgresql/pg_hba.conf
ADD ./postgresql.conf /etc/postgresql/postgresql.conf

# update ca file permission for postgres user
RUN ./home/ca/init_cert.sh
CMD ["postgres", "-c", "config_file=/etc/postgresql/postgresql.conf"]