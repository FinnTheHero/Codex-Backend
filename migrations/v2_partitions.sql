DO $$
DECLARE
    i INTEGER;
    partition_name TEXT;
    partitions_count INTEGER := 16;
BEGIN
    FOR i IN 0..partitions_count-1 LOOP
        partition_name := 'chapters_p' || i;

        EXECUTE format('CREATE TABLE IF NOT EXISTS %I PARTITION OF chapters
                       FOR VALUES WITH (MODULUS %s, REMAINDER %s)',
                      partition_name, partitions_count, i);

        EXECUTE format('CREATE UNIQUE INDEX IF NOT EXISTS idx_%I_id ON %I (id)',
                      partition_name, partition_name);
    END LOOP;
END $$;
