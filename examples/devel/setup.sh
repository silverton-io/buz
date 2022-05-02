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
find schemas -type f | while read fname; do
    SCHEMA=$(echo $fname | sed 's/schemas\///')
    CONTENTS=$(cat $fname)
    # Postgres
    psql -h 127.0.0.1 -p 5432 -U honeypot -c "insert into honeypot_schemas (created_at, updated_at, name, schema) values (now(), now(), '$SCHEMA','$CONTENTS');"
    # Materialize
    # psql -h 127.0.0.1 -p 6875 -U honeypot -c "insert into honeypot_schemas (created_at, updated_at, name, schema) values (now(), now(), '$SCHEMA','$CONTENTS');"
    # Mysql
done
