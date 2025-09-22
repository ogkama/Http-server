CREATE TABLE users (
    user_id UUID PRIMARY KEY,
    user_login VARCHAR(50) NOT NULL UNIQUE,
    user_password TEXT NOT NULL
);

CREATE TABLE tasks (
    task_id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    task_status TEXT NOT NULL CHECK (task_status IN ('in progress', 'done','ERROR')),
    task_data JSON NOT NULL,
    task_result JSON,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);



create or replace function set_updated_at()
returns trigger language plpgsql as $$
begin
  new.updated_at = now();
  return new;
end ;
$$;

create trigger trg_users_updated
before update on users
for each row execute function set_updated_at();
 
create trigger trg_tasks_updated
before update on tasks
for each row execute function set_updated_at();