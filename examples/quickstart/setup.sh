
# Create necessary kafka topics
rpk topic \
    create hpt-invalid \
    --brokers 127.0.0.1:9092;

rpk topic \
    create hpt-valid \
    --brokers 127.0.0.1:9092;

# Set up materialize
psql \
    -h 127.0.0.1 \
    -p 6875 \
    -U materialize \
    -f examples/quickstart/materialize/create_sources_and_views.sql;
