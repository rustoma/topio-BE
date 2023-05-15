-- Rename the table to a temporary name
ALTER TABLE product RENAME TO temp_product;
ALTER TABLE page RENAME TO temp_page;


-- Create the new table with the updated schema
CREATE TABLE product (
                         id SERIAL PRIMARY KEY,
                         main_category VARCHAR(255),
                         sub_category VARCHAR(255),
                         description TEXT,
                         url VARCHAR(255),
                         order_weight INT,
                         score REAL,
                         main_image INT REFERENCES public.image(ID),
                         created_at TIMESTAMPTZ,
                         updated_at TIMESTAMPTZ
);

-- Copy the data from the temporary table to the new table
INSERT INTO product (id, main_category, sub_category, description, url, order_weight, score, main_image, created_at, updated_at)
SELECT id, main_category, sub_category, description, url, order_weight, score, main_image, created_at, updated_at
FROM temp_product;


-- Create the new table with the updated schema
CREATE TABLE page (
                      id SERIAL PRIMARY KEY,
                      title VARCHAR(255) NOT NULL,
                      intro TEXT,
                      body TEXT,
                      products INT REFERENCES public.product(ID),
                      slug VARCHAR(255) UNIQUE NOT NULL
);



-- Copy the data from the temporary table to the new table
INSERT INTO page (id, title, intro, body, products, slug)
SELECT id, title, intro, body, products, slug
FROM temp_page;

-- Drop the temporary table
DROP TABLE temp_page;
DROP TABLE temp_product