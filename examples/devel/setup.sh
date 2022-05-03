REGISTRY_HOST=localhost
REGISTRY_DB=honeypot
REGISTRY_SCHEMA=honeypot
REGISTRY_TABLE=registry
H_USER=honeypot
H_PASS=honeypot
# Seed files
PG_SEED_FILE=pgSeed.sql
MYSQL_SEED_FILE=mysqlSeed.sql
CLICKHOUSE_SEED_FILE=clickhouseSeed.sql
MATERIALIZE_SEED_FILE=materializeSeed.sql


function cleanSeedFiles {
    echo "Cleaning seed files...";
    rm $PG_SEED_FILE || true;
    rm $MYSQL_SEED_FILE || true;
    rm $CLICKHOUSE_SEED_FILE || true;
}

function seedDatabaseBackends {
    echo "\nSeeding database schema cache backends...\n"
    cleanSeedFiles;
    echo "delete from $REGISTRY_TABLE;" >> $PG_SEED_FILE;
    echo "delete from $REGISTRY_TABLE;" >> $MATERIALIZE_SEED_FILE;
    echo "delete from $REGISTRY_SCHEMA.$REGISTRY_TABLE;" >> $MYSQL_SEED_FILE;
    echo "alter table $REGISTRY_SCHEMA.$REGISTRY_TABLE delete where 1=1;" >> $CLICKHOUSE_SEED_FILE;
    find schemas -type f | while read fname; do
        SCHEMA=$(echo $fname | sed 's/schemas\///')
        CONTENTS=$(jq -c . $fname)
        # Postgres
        PG_CMD="insert into $REGISTRY_TABLE (created_at, updated_at, name, contents) values (now(), now(), '$SCHEMA','$CONTENTS');"
        echo $PG_CMD >> $PG_SEED_FILE;
        # Materialize
        echo $PG_CMD >> $MATERIALIZE_SEED_FILE;
        # Mysql
        MYSQL_CMD="insert into $REGISTRY_SCHEMA.$REGISTRY_TABLE (created_at, updated_at, name, contents) values (now(), now(), '$SCHEMA','$CONTENTS');"
        echo $MYSQL_CMD >> $MYSQL_SEED_FILE;
        # Clickhouse
        CLICKHOUSE_CMD="insert into $REGISTRY_SCHEMA.$REGISTRY_TABLE (created_at, updated_at, name, contents) values (now(), now(), '$SCHEMA', '$CONTENTS')"
        echo $CLICKHOUSE_CMD >> $CLICKHOUSE_SEED_FILE;
    done

    psql -h $REGISTRY_HOST -p 5432 -U $H_USER -f $PG_SEED_FILE;
    psql -h $REGISTRY_HOST -p 6875 -U materialize -f $PG_SEED_FILE;
    export MYSQL_PWD=honeypot; mysql -h 127.0.0.1 -u$H_USER < $MYSQL_SEED_FILE;
    clickhouse client -h $REGISTRY_HOST  --port 9000 -u $H_USER --password $H_PASS --queries-file $CLICKHOUSE_SEED_FILE
    cleanSeedFiles;
}

echo "\nHoneypotting...\n"
for run in {1..5}; do printf "🍯" && sleep 1; done

echo "\nSetting up clickhouse...\n"
# This is required due to the inability to pass env vars to clickhouse img.
clickhouse client -u $H_USER --password $H_PASS -q \"create database $REGISTRY_DB;\"

echo "\nSetting up Redpanda...\n";
rpk topic \
    create hpt-invalid \
    --brokers 127.0.0.1:9092;

rpk topic \
    create hpt-valid \
    --brokers 127.0.0.1:9092;


seedDatabaseBackends;