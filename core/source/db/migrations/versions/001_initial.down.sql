DROP INDEX "idx_inventory_items_user_uid";
DROP INDEX "idx_recipes_is_trashed";
DROP INDEX "idx_recipes_user_uid";
DROP INDEX "idx_recipe_ingredients_item";
DROP INDEX "idx_recipe_ingredients_user_uid";
DROP INDEX "idx_recipe_ingredients_recipe_uid";
DROP INDEX "idx_grocery_items_user_uid";
DROP INDEX "idx_users_email";

DROP TABLE "grocery_items";
DROP TABLE "inventory_items";
DROP TABLE "o_auth_states";
DROP TABLE "recipes";
DROP TABLE "recipe_ingredients";
DROP TABLE "users";
