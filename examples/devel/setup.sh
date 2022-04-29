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
