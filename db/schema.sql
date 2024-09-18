DROP DATABASE IF EXISTS MoonTech;
CREATE DATABASE MoonTech;
USE MoonTech;
-- Tabla de usuarios
CREATE TABLE User (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) NOT NULL,
    firstname VARCHAR(50),
    lastname VARCHAR(50),
    workEmail VARCHAR(100) UNIQUE,
    workPhone VARCHAR(20),
    password VARCHAR(255) NOT NULL,
    area INT,
    leaderId INT,
    position VARCHAR(100),
    role VARCHAR(50)
);

-- Tabla de proyectos
CREATE TABLE Project (
    id INT AUTO_INCREMENT PRIMARY KEY,
    projName VARCHAR(100),
    owner INT, -- Hace referencia al userId
    company INT, -- Hace referencia a la tabla de Companies
    area INT, -- Hace referencia a la tabla de Areas
    startDate DATE,
    FOREIGN KEY (owner) REFERENCES User(id),
    FOREIGN KEY (company) REFERENCES Companies(id),
    FOREIGN KEY (area) REFERENCES Area(id)
);

-- Tabla de líderes
CREATE TABLE Leaders (
    id INT AUTO_INCREMENT PRIMARY KEY,
    projId INT, -- Hace referencia a Project
    userId INT, -- Hace referencia a User
    area INT, -- Hace referencia a Area
    FOREIGN KEY (projId) REFERENCES Project(id),
    FOREIGN KEY (userId) REFERENCES User(id),
    FOREIGN KEY (area) REFERENCES Area(id)
);

-- Tabla de requerimientos
CREATE TABLE Requirements (
    id INT AUTO_INCREMENT PRIMARY KEY,
    projectId INT, -- Hace referencia a Project
    owner INT, -- Hace referencia a User
    text TEXT,
    timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
    approved BOOLEAN,
    approverId INT, -- Hace referencia a User
    FOREIGN KEY (projectId) REFERENCES Project(id),
    FOREIGN KEY (owner) REFERENCES User(id),
    FOREIGN KEY (approverId) REFERENCES User(id)
);

-- Tabla de comentarios
CREATE TABLE Comment (
    id INT AUTO_INCREMENT PRIMARY KEY,
    owner INT, -- Hace referencia a User
    parent INT, -- Puede ser un requerimiento o una tarea
    text TEXT,
    timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (owner) REFERENCES User(id)
);

-- Tabla de tareas
CREATE TABLE Tasks (
    id INT AUTO_INCREMENT PRIMARY KEY,
    area INT, -- Hace referencia a Area
    title VARCHAR(100),
    createdBy INT, -- Hace referencia a User
    description TEXT,
    timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
    estimatedTime INT,
    approved BOOLEAN,
    approverId INT, -- Hace referencia a User
    FOREIGN KEY (area) REFERENCES Area(id),
    FOREIGN KEY (createdBy) REFERENCES User(id),
    FOREIGN KEY (approverId) REFERENCES User(id)
);

-- Tabla de compañías
CREATE TABLE Companies (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100),
    representativeId INT, -- Hace referencia a la tabla de Representantes
    businessType INT, -- Hace referencia a la tabla businessType
    FOREIGN KEY (representativeId) REFERENCES Representatives(id),
    FOREIGN KEY (businessType) REFERENCES businessType(id)
);

-- Tabla de tipos de negocio
CREATE TABLE businessType (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100),
    color VARCHAR(20)
);

-- Tabla de representantes
CREATE TABLE Representatives (
    id INT AUTO_INCREMENT PRIMARY KEY,
    firstname VARCHAR(50),
    lastname VARCHAR(50),
    workEmail VARCHAR(100),
    workPhone VARCHAR(20)
);

-- Tabla de áreas
CREATE TABLE Area (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100),
    description TEXT
);
