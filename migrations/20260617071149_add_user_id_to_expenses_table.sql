-- +goose Up
-- Step 1: Add column as nullable first (safe with existing data)
ALTER TABLE expenses ADD COLUMN user_id UUID;

-- Step 2: Backfill existing expenses with the first available user
-- ⚠️ Run this ONLY after you've registered at least one user
UPDATE expenses
SET user_id = (SELECT id FROM users ORDER BY created_at LIMIT 1)
WHERE user_id IS NULL;

-- Step 3: Add foreign key constraint
ALTER TABLE expenses ADD CONSTRAINT fk_expenses_user
    FOREIGN KEY (user_id) REFERENCES users(id);

-- Step 4: Make column NOT NULL (now safe because all rows have a valid user_id)
ALTER TABLE expenses ALTER COLUMN user_id SET NOT NULL;

-- Step 5: Add index for dashboard performance
CREATE INDEX idx_expenses_user_date ON expenses(user_id, date DESC);

-- +goose Down
DROP INDEX IF EXISTS idx_expenses_user_date;
ALTER TABLE expenses DROP CONSTRAINT IF EXISTS fk_expenses_user;
ALTER TABLE expenses DROP COLUMN IF EXISTS user_id;
