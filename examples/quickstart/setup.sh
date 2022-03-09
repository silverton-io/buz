
# Create necessary kafka topics
rpk topic \
    create hpt-invalid \
    --brokers 127.0.0.1:9092;

rpk topic \
    create hpt-valid \
    --brokers 127.0.0.1:9092;

psql \
    -h 127.0.0.1 \
    -p 6875 \
    -U materialize \
    -f create_sources_and_views.sql;
