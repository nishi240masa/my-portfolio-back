-- users テーブル
CREATE TABLE IF NOT EXISTS users (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(100) NOT NULL,
    role VARCHAR(50) DEFAULT 'user',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- apiSchema テーブル
CREATE TABLE IF NOT EXISTS apiSchema (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    field_id VARCHAR(100) NOT NULL, 
    view_name VARCHAR(100) NOT NULL,
    field_type VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- field_data テーブル
CREATE TABLE IF NOT EXISTS field_data (
    id SERIAL PRIMARY KEY,
    apiSchema_id INT NOT NULL REFERENCES apiSchema(id) ON DELETE CASCADE,
    field_type VARCHAR(50) NOT NULL,
    field_value JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- list_options テーブル
CREATE TABLE IF NOT EXISTS list_options (
    id SERIAL PRIMARY KEY,
    apiSchema_id INT NOT NULL REFERENCES apiSchema(id) ON DELETE CASCADE,
    value VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- api_kind_relation テーブル
CREATE TABLE IF NOT EXISTS api_kind_relation (
    id SERIAL PRIMARY KEY,
    apiSchema_id INT NOT NULL REFERENCES apiSchema(id) ON DELETE CASCADE,
    related_id INT NOT NULL REFERENCES apiSchema(id) ON DELETE CASCADE,
    relation_type VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- タイムスタンプ更新のトリガー関数
CREATE OR REPLACE FUNCTION update_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- 各テーブルにトリガーを設定
CREATE TRIGGER set_timestamp_users
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();

CREATE TRIGGER set_timestamp_apiSchema
BEFORE UPDATE ON apiSchema
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();

CREATE TRIGGER set_timestamp_field_data
BEFORE UPDATE ON field_data
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();

CREATE TRIGGER set_timestamp_list_options
BEFORE UPDATE ON list_options
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();

CREATE TRIGGER set_timestamp_api_kind_relation
BEFORE UPDATE ON api_kind_relation
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();

-- apiSchema 削除時のトリガー関数
CREATE OR REPLACE FUNCTION cascade_delete_apiSchema()
RETURNS TRIGGER AS $$
BEGIN
    DELETE FROM api_kind_relation WHERE related_id = OLD.id;
    RETURN OLD;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER delete_related_data
AFTER DELETE ON apiSchema
FOR EACH ROW
EXECUTE FUNCTION cascade_delete_apiSchema();

-- list_options 制約トリガー
CREATE OR REPLACE FUNCTION validate_list_options()
RETURNS TRIGGER AS $$
DECLARE
    schema_type VARCHAR(50);
BEGIN
    SELECT field_type INTO schema_type FROM apiSchema WHERE id = NEW.apiSchema_id;
    IF schema_type NOT IN ('select', 'dropdown') THEN
        RAISE EXCEPTION 'list_options can only be added to select or dropdown fields';
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER validate_options
BEFORE INSERT ON list_options
FOR EACH ROW
EXECUTE FUNCTION validate_list_options();

-- 循環リレーション防止トリガー
CREATE OR REPLACE FUNCTION check_cyclic_relation()
RETURNS TRIGGER AS $$
DECLARE
    is_cyclic BOOLEAN;
BEGIN
    WITH RECURSIVE relation_path AS (
        SELECT related_id FROM api_kind_relation WHERE apiSchema_id = NEW.related_id
        UNION ALL
        SELECT r.related_id FROM api_kind_relation r
        INNER JOIN relation_path rp ON rp.related_id = r.apiSchema_id
    )
    SELECT EXISTS (
        SELECT 1 FROM relation_path WHERE related_id = NEW.apiSchema_id
    ) INTO is_cyclic;

    IF is_cyclic THEN
        RAISE EXCEPTION 'Cyclic relation detected';
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER prevent_cyclic_relation
BEFORE INSERT ON api_kind_relation
FOR EACH ROW
EXECUTE FUNCTION check_cyclic_relation();
