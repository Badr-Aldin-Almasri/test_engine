-- Migration: Separate nodes from flows table into flow_nodes table
-- This migration moves nodes from the JSONB column in flows to a separate table

-- Step 1: Create the flow_nodes table
CREATE TABLE IF NOT EXISTS flow_nodes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    flow_id UUID REFERENCES flows(id) ON DELETE CASCADE,
    node_id VARCHAR(255) NOT NULL, -- The node ID used in the flow
    type VARCHAR(50) NOT NULL,
    position_x DOUBLE PRECISION NOT NULL,
    position_y DOUBLE PRECISION NOT NULL,
    data_id VARCHAR(255) NOT NULL,
    data_type VARCHAR(50) NOT NULL,
    data_label VARCHAR(255),
    data_status VARCHAR(50),
    data_config JSONB DEFAULT '{}'::jsonb,
    data_output JSONB,
    data_error TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(flow_id, node_id) -- Ensure node_id is unique per flow
);

-- Step 2: Migrate existing nodes from flows.nodes JSONB to flow_nodes table
-- This handles existing data by extracting nodes from the JSONB column
DO $$
DECLARE
    flow_record RECORD;
    node_item JSONB;
    node_id_val TEXT;
    node_type_val TEXT;
    position_x_val DOUBLE PRECISION;
    position_y_val DOUBLE PRECISION;
    data_id_val TEXT;
    data_type_val TEXT;
    data_label_val TEXT;
    data_status_val TEXT;
    data_config_val JSONB;
    data_output_val JSONB;
    data_error_val TEXT;
BEGIN
    -- Loop through all flows
    FOR flow_record IN SELECT id, nodes FROM flows WHERE nodes IS NOT NULL AND nodes != '[]'::jsonb
    LOOP
        -- Extract each node from the JSONB array
        FOR node_item IN SELECT * FROM jsonb_array_elements(flow_record.nodes)
        LOOP
            -- Extract node fields
            node_id_val := node_item->>'id';
            node_type_val := node_item->>'type';
            
            -- Extract position
            IF node_item->'position' IS NOT NULL THEN
                position_x_val := (node_item->'position'->>'x')::DOUBLE PRECISION;
                position_y_val := (node_item->'position'->>'y')::DOUBLE PRECISION;
            ELSE
                position_x_val := 0;
                position_y_val := 0;
            END IF;
            
            -- Extract data fields
            IF node_item->'data' IS NOT NULL THEN
                data_id_val := node_item->'data'->>'id';
                data_type_val := node_item->'data'->>'type';
                data_label_val := node_item->'data'->>'label';
                data_status_val := node_item->'data'->>'status';
                
                -- Extract config (default to empty object if null)
                IF node_item->'data'->'config' IS NOT NULL THEN
                    data_config_val := node_item->'data'->'config';
                ELSE
                    data_config_val := '{}'::jsonb;
                END IF;
                
                -- Extract output (can be null)
                IF node_item->'data'->'output' IS NOT NULL THEN
                    data_output_val := node_item->'data'->'output';
                ELSE
                    data_output_val := NULL;
                END IF;
                
                -- Extract error (can be null)
                data_error_val := node_item->'data'->>'error';
                
                -- Set defaults for required fields if null
                IF data_id_val IS NULL THEN
                    data_id_val := node_id_val;
                END IF;
                IF data_type_val IS NULL THEN
                    data_type_val := node_type_val;
                END IF;
                IF data_label_val IS NULL THEN
                    data_label_val := '';
                END IF;
            ELSE
                -- If data is missing, use node fields as defaults
                data_id_val := node_id_val;
                data_type_val := node_type_val;
                data_label_val := '';
                data_status_val := NULL;
                data_config_val := '{}'::jsonb;
                data_output_val := NULL;
                data_error_val := NULL;
            END IF;
            
            -- Insert node into flow_nodes table
            INSERT INTO flow_nodes (
                flow_id, node_id, type, position_x, position_y,
                data_id, data_type, data_label, data_status,
                data_config, data_output, data_error,
                created_at, updated_at
            ) VALUES (
                flow_record.id,
                node_id_val,
                node_type_val,
                position_x_val,
                position_y_val,
                data_id_val,
                data_type_val,
                data_label_val,
                data_status_val,
                data_config_val,
                data_output_val,
                data_error_val,
                CURRENT_TIMESTAMP,
                CURRENT_TIMESTAMP
            ) ON CONFLICT (flow_id, node_id) DO NOTHING;
        END LOOP;
    END LOOP;
END $$;

-- Step 3: Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_flow_nodes_flow_id ON flow_nodes(flow_id);
CREATE INDEX IF NOT EXISTS idx_flow_nodes_node_id ON flow_nodes(flow_id, node_id);

-- Step 4: Create trigger for updated_at timestamp
CREATE TRIGGER update_flow_nodes_updated_at BEFORE UPDATE ON flow_nodes
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Step 5: Remove the nodes column from flows table
-- Note: This is the final step - uncomment only after verifying the migration worked
-- ALTER TABLE flows DROP COLUMN IF EXISTS nodes;

-- For safety, we'll keep the nodes column commented out above
-- You can manually verify the migration worked, then uncomment the line above
-- Or run it separately after verification

