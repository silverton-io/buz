echo "\nHoneypotting...\n"
for run in {1..5}; do printf "🍯" && sleep 1; done

echo "\nSetting up clickhouse...\n"
docker exec clickhouse sh -c "clickhouse-client -u honeypot --password honeypot -q \"create database honeypot;\""

echo "\nSetting up materialize...\n"
docker exec materialize sh -c "psql -"
