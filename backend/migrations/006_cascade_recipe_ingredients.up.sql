ALTER TABLE recipe_ingredients DROP CONSTRAINT recipe_ingredients_ingredient_id_fkey;
ALTER TABLE recipe_ingredients ADD CONSTRAINT recipe_ingredients_ingredient_id_fkey FOREIGN KEY (ingredient_id) REFERENCES products(id) ON DELETE CASCADE;
