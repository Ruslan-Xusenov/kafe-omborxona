CREATE TABLE recipes (
    id SERIAL PRIMARY KEY,
    product_id INT NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    name VARCHAR(200) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE recipe_ingredients (
    id SERIAL PRIMARY KEY,
    recipe_id INT NOT NULL REFERENCES recipes(id) ON DELETE CASCADE,
    ingredient_id INT NOT NULL REFERENCES products(id) ON DELETE RESTRICT,
    quantity DECIMAL(14,3) NOT NULL CHECK (quantity > 0)
);

CREATE INDEX idx_recipes_product ON recipes(product_id);
CREATE INDEX idx_recipe_ingredients_recipe ON recipe_ingredients(recipe_id);
