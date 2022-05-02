REGISTRY_SCHEMA=honeypot
REGISTRY_TABLE=schemas
H_USER=honeypot
H_PASS=honeypot


echo "\nHoneypotting...\n"
for run in {1..5}; do printf "🍯" && sleep 1; done

echo "\nSetting up clickhouse...\n"
docker exec clickhouse sh -c "clickhouse-client -u honeypot --password honeypot -q \"create database honeypot;\""

echo "\nSetting up Redpanda...\n";
rpk topic \
    create hpt-invalid \
    --brokers 127.0.0.1:9092;

rpk topic \
    create hpt-valid \
    --brokers 127.0.0.1:9092;

echo "\nSeeding database schema cache backends...\n"
echo "deleting all schemas from postgres"
psql -h 127.0.0.1 -p 5432 -U $H_USER -c "delete from $REGISTRY_TABLE;"
echo "deleting all schemas from mysql"
export MYSQL_PWD=honeypot; mysql -h 127.0.0.1 -u$H_USER -e "delete from $REGISTRY_SCHEMA.$REGISTRY_TABLE;"
echo "deleting all schemas from materialize"
psql -h 127.0.0.1 -p 6875 -U materialize -c "delete from $REGISTRY_TABLE;"
find schemas -type f | while read fname; do
    SCHEMA=$(echo $fname | sed 's/schemas\///')
    CONTENTS=$(jq -c . $fname)
    # Postgres
    echo "seeding schema to postgres: $SCHEMA"
    PG_CMD="insert into $REGISTRY_TABLE (created_at, updated_at, name, contents) values (now(), now(), '$SCHEMA','$CONTENTS');"
    psql -h 127.0.0.1 -p 5432 -U $H_USER -c "$PG_CMD"
    # Mysql
    echo "seeding schema to mysql: $SCHEMA"
    MYSQL_CMD="insert into $REGISTRY_SCHEMA.$REGISTRY_TABLE (created_at, updated_at, name, contents) values (now(), now(), '$SCHEMA','$CONTENTS');"
    export MYSQL_PWD=honeypot; mysql -h 127.0.0.1 -u$H_USER -e "$MYSQL_CMD"
    # Materialize
    psql -h 127.0.0.1 -p 6875 -U materialize -c "$PG_CMD"
done
