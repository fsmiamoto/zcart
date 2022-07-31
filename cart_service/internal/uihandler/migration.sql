CREATE TABLE IF NOT EXISTS products (
    id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    price REAL NOT NULL,
    description VARCHAR(255),
    image_url VARCHAR(255)
);

CREATE TABLE IF NOT EXISTS cart_products (
    cart_id VARCHAR(255),
    product_id VARCHAR(255),
    quantity INTEGER,
    PRIMARY KEY (cart_id, product_id),
    FOREIGN KEY (product_id) REFERENCES products (id)
);

INSERT INTO products (id,name,price,image_url) VALUES ('1','Coca Cola', 5.99, 'https://zcart-test-images.s3.amazonaws.com/coca2l.png');
INSERT INTO products (id,name,price,image_url) VALUES ('2','BomBril', 1.99, 'https://zcart-test-images.s3.amazonaws.com/bombril.png');
INSERT INTO products (id,name,price,image_url) VALUES ('3','Leite Longa Vida 1L', 4.99, 'https://zcart-test-images.s3.amazonaws.com/leite.png');
INSERT INTO products (id,name,price,image_url) VALUES ('4','Café', 8.99, 'https://zcart-test-images.s3.amazonaws.com/cafe.png');
INSERT INTO products (id,name,price,image_url) VALUES ('5','Chamyto', 10.99, 'https://zcart-test-images.s3.amazonaws.com/chamyto.png');
INSERT INTO products (id,name,price,image_url) VALUES ('6','Macarrão Instântaneo Nissin', 3.99, 'https://zcart-test-images.s3.amazonaws.com/lamen.png');

INSERT INTO cart_products (cart_id,product_id,quantity) VALUES ('2','1', 10);
INSERT INTO cart_products (cart_id,product_id,quantity) VALUES ('2','2', 5);
INSERT INTO cart_products (cart_id,product_id,quantity) VALUES ('2','3', 9);
INSERT INTO cart_products (cart_id,product_id,quantity) VALUES ('2','4', 1);
