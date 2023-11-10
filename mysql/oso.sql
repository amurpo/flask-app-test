-- Crea la base de datos "oso" si no existe
CREATE DATABASE IF NOT EXISTS oso;

-- Utiliza la base de datos "oso"
USE oso;

-- Crea la tabla "images" si no existe
CREATE TABLE IF NOT EXISTS images (
  id INT AUTO_INCREMENT PRIMARY KEY,
  link VARCHAR(255) NOT NULL
);

-- Inserta los registros en la tabla "images"
INSERT INTO images (link) VALUES
  ('https://oso-flask.s3.sa-east-1.amazonaws.com/oso.jpg'),
  ('https://oso-flask.s3.sa-east-1.amazonaws.com/oso2.jpg'),
  ('https://oso-flask.s3.sa-east-1.amazonaws.com/oso33.jpg');
