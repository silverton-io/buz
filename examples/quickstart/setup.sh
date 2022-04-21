
# Create necessary kafka topics
echo "Creating Redpanda topics...\n";
rpk topic \
    create hpt-invalid \
    --brokers 127.0.0.1:9092;

rpk topic \
    create hpt-valid \
    --brokers 127.0.0.1:9092;

# Set up materialize
echo "\nCreating Materialize sources and views...\n";
psql \
    -h 127.0.0.1 \
    -p 6875 \
    -U materialize \
    -f examples/quickstart/materialize/create_sources_and_views.sql;

echo "\nHoneypotting...\n"
for run in {1..15}; do printf "🍯" && sleep 1; done
echo "\n\nOpening associated resources...\n";

open http://localhost:8080/;
sleep 2;
open http://localhost:8080/stats;
sleep 2;
open http://localhost:8080/schemas;
sleep 2;
open http://localhost:8080/schemas/com.silverton.io/snowplow/page_view/v1.0.json;
sleep 2;
open http://localhost:8081/topics;
sleep 2;
