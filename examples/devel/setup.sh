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
psql -h 127.0.0.1 -p 5432 -U honeypot -c "delete from honeypot_schemas;"
echo "deleting all schemas from mysql"
export MYSQL_PWD=honeypot; mysql -h 127.0.0.1 -uhoneypot -e "delete from honeypot.honeypot_schemas;"
find schemas -type f | while read fname; do
    SCHEMA=$(echo $fname | sed 's/schemas\///')
    CONTENTS=$(jq -c . $fname)
    # Postgres
    echo "seeding schema to postgres: $SCHEMA"
    PG_CMD="insert into honeypot_schemas (created_at, updated_at, name, contents) values (now(), now(), '$SCHEMA','$CONTENTS');"
    psql -h 127.0.0.1 -p 5432 -U honeypot -c "$PG_CMD"
    # Mysql
    echo "seeding schema to mysql: $SCHEMA"
    MYSQL_CMD="insert into honeypot.honeypot_schemas (created_at, updated_at, name, contents) values (now(), now(), '$SCHEMA','$CONTENTS');"
    export MYSQL_PWD=honeypot; mysql -h 127.0.0.1 -uhoneypot -e "$MYSQL_CMD"
done
