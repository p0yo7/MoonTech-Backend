DROP DATABASE IF EXISTS MoonTech;
CREATE DATABASE MoonTech;
USE MoonTech;

-- Tabla de tipos de negocio
CREATE TABLE businessTypes (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100),
    color VARCHAR(20)
);

-- Tabla de representantes
CREATE TABLE representatives (
    id INT AUTO_INCREMENT PRIMARY KEY,
    firstname VARCHAR(50),
    lastname VARCHAR(50),
    work_email VARCHAR(100),
    work_phone VARCHAR(20)
);

-- Tabla de áreas
CREATE TABLE areas (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100),
    description TEXT
);

-- Tabla de compañías
CREATE TABLE companies (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100),
    representativeId INT, -- Hace referencia a la tabla de Representantes
    businessType INT, -- Hace referencia a la tabla businessType
    FOREIGN KEY (representativeId) REFERENCES representatives(id),
    FOREIGN KEY (businessType) REFERENCES businessTypes(id)
);

-- Tabla de usuarios
CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) NOT NULL,
    first_name VARCHAR(50),
    lastname VARCHAR(50),
    work_email VARCHAR(100) UNIQUE,
    work_phone VARCHAR(20),
    password VARCHAR(255) NOT NULL,
    area INT,
    leaderId INT,
    position VARCHAR(100),
    role VARCHAR(50),
    FOREIGN KEY (area) REFERENCES areas(id),
    FOREIGN KEY (leaderId) REFERENCES users(id) -- Referencia a sí mismo
);

-- Tabla de proyectos
CREATE TABLE projects (
    id INT AUTO_INCREMENT PRIMARY KEY,
    projName VARCHAR(100),
    owner INT, -- Hace referencia al userId
    company INT, -- Hace referencia a la tabla de Companies
    area INT, -- Hace referencia a la tabla de Areas
    startDate DATE,
    FOREIGN KEY (owner) REFERENCES users(id),
    FOREIGN KEY (company) REFERENCES companies(id),
    FOREIGN KEY (area) REFERENCES areas(id)
);

-- Tabla de líderes
CREATE TABLE leaders (
    id INT AUTO_INCREMENT PRIMARY KEY,
    projId INT, -- Hace referencia a Project
    userId INT, -- Hace referencia a User
    area INT, -- Hace referencia a Area
    FOREIGN KEY (projId) REFERENCES projects(id),
    FOREIGN KEY (userId) REFERENCES users(id),
    FOREIGN KEY (area) REFERENCES areas(id)
);

-- Tabla de requerimientos
CREATE TABLE requirements (
    id INT AUTO_INCREMENT PRIMARY KEY,
    projectId INT, -- Hace referencia a Project
    owner INT, -- Hace referencia a User
    text TEXT,
    timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
    approved BOOLEAN,
    approverId INT, -- Hace referencia a User
    FOREIGN KEY (projectId) REFERENCES projects(id),
    FOREIGN KEY (owner) REFERENCES users(id),
    FOREIGN KEY (approverId) REFERENCES users(id)
);

-- Tabla de comentarios
CREATE TABLE comments (
    id INT AUTO_INCREMENT PRIMARY KEY,
    owner INT, -- Hace referencia a User
    parent INT, -- Puede ser un requerimiento o una tarea
    text TEXT,
    timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (owner) REFERENCES users(id)
);

-- Tabla de tareas
CREATE TABLE tasks (
    id INT AUTO_INCREMENT PRIMARY KEY,
    area INT, -- Hace referencia a Area
    title VARCHAR(100),
    createdBy INT, -- Hace referencia a User
    description TEXT,
    timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
    estimatedTime INT,
    approved BOOLEAN,
    approverId INT, -- Hace referencia a User
    FOREIGN KEY (area) REFERENCES areas(id),
    FOREIGN KEY (createdBy) REFERENCES users(id),
    FOREIGN KEY (approverId) REFERENCES users(id)
);