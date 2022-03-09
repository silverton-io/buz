-- Create valid events source
create source
    src_valid_events
from
    kafka
        broker 'redpanda-1:29092,redpanda-2:29093,redpanda-3:29094'
        topic 'hpt-valid'
        format bytes;


-- Create invalid events source
create source
    src_invalid_events
from
    kafka
        broker 'redpanda-1:29092,redpanda-2:29093,redpanda-3:29094'
        topic 'hpt-invalid'
        format bytes;


-- Create valid events view
create materialized view valid_events as
    select
        *
    from (
        select
            convert_from(data, 'utf8')::jsonb as event
        from
            src_valid_events 
    );

create materialized view invalid_events as
    select
        *
    from (
        select
            convert_from(data, 'utf8')::jsonb as event
        from
            src_invalid_events 
    );