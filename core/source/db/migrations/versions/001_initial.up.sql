CREATE TABLE IF NOT EXISTS "grocery_items" ("uid" text NOT NULL DEFAULT null,"created_at" timestamptz,"updated_at" timestamptz,"user_uid" text,"item" text,"store" text,"is_collected" boolean,PRIMARY KEY ("uid"));
CREATE TABLE IF NOT EXISTS "inventory_items" ("uid" text NOT NULL DEFAULT null,"updated_at" timestamptz,"created_at" timestamptz,"user_uid" text,"item" text,"in_stock" boolean,"category" text DEFAULT 'cooking',PRIMARY KEY ("uid"));
CREATE TABLE IF NOT EXISTS "o_auth_states" ("ip_address" text NOT NULL DEFAULT null,"ttl" timestamptz NOT NULL,"provider" text NOT NULL,"state" text NOT NULL,PRIMARY KEY ("ip_address"));
CREATE TABLE IF NOT EXISTS "recipes" ("uid" text NOT NULL DEFAULT null,"created_at" timestamptz,"updated_at" timestamptz,"is_trashed" boolean DEFAULT false,"trash_at" bigint,"user_uid" text,"name" text,"description" text,"meal" text,"instructions" text,PRIMARY KEY ("uid"));
CREATE TABLE IF NOT EXISTS "recipe_ingredients" ("uid" text NOT NULL DEFAULT null,"recipe_uid" text,"user_uid" text,"created_at" timestamptz,"updated_at" timestamptz,"sort_order" bigint,"quantity" text,"unit" text,"item" text,"optional" boolean DEFAULT false,"modifier" text,PRIMARY KEY ("uid"),CONSTRAINT "fk_recipes_ingredients" FOREIGN KEY ("recipe_uid") REFERENCES "recipes"("uid") ON DELETE CASCADE ON UPDATE CASCADE);
CREATE TABLE IF NOT EXISTS "users" ("uid" text NOT NULL DEFAULT null,"created_at" timestamptz,"updated_at" timestamptz,"role" text,"email" text,"password" bytea,"github_token" text DEFAULT null,"google_token" text DEFAULT null,"facebook_token" text DEFAULT null,"microsoft_token" text DEFAULT null,PRIMARY KEY ("uid"),CONSTRAINT "uni_users_email" UNIQUE ("email"));

CREATE INDEX IF NOT EXISTS "idx_inventory_items_user_uid" ON "inventory_items" ("user_uid");
CREATE INDEX IF NOT EXISTS "idx_recipes_is_trashed" ON "recipes" ("is_trashed");
CREATE INDEX IF NOT EXISTS "idx_recipes_user_uid" ON "recipes" ("user_uid");
CREATE INDEX IF NOT EXISTS "idx_recipe_ingredients_item" ON "recipe_ingredients" ("item");
CREATE INDEX IF NOT EXISTS "idx_recipe_ingredients_user_uid" ON "recipe_ingredients" ("user_uid");
CREATE INDEX IF NOT EXISTS "idx_recipe_ingredients_recipe_uid" ON "recipe_ingredients" ("recipe_uid");
CREATE INDEX IF NOT EXISTS "idx_grocery_items_user_uid" ON "grocery_items" ("user_uid");
CREATE INDEX IF NOT EXISTS "idx_users_email" ON "users" ("email");