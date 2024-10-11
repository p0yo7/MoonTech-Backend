DROP DATABASE IF EXISTS MoonTech;
CREATE DATABASE MoonTech;
USE MoonTech;

CREATE TABLE frameworks (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(40),
    language VARCHAR(40),
    licence VARCHAR(40),
    compatibility VARCHAR(40)
);


CREATE TABLE teams ( 
    id INT AUTO_INCREMENT PRIMARY KEY,
    team_name varchar(100)
);

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
-- descripcion

-- aregar tamano
-- agregar ubicacion
-- ubicacion 
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
    team INT,
    leaderId INT,
    position VARCHAR(100),
    role VARCHAR(50),
    FOREIGN KEY (team) REFERENCES teams(id),
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
    budget INT, -- budget en dolares
    status INT, -- active 1 inactive 0
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
    FOREIGN KEY (owner) REFERENCES users(id),
    FOREIGN KEY (parent) REFERENCES requirements(id)
);

-- id, requerimiento, createdBy, nombre, descripcion, area (equipo), lenguajes, frameworks, tiempo, costo, ajuste, 
-- Integrar la parte de ajuste (feedback), descripcion, lenguajes, frameworks 
-- Tabla de tareas
-- CREATE TABLE tasks (
--     id INT AUTO_INCREMENT PRIMARY KEY,
--     requirement INT, 
--     area INT, -- Hace referencia a Area
--     name VARCHAR(100),
--     createdBy INT, -- Hace referencia a User
--     description TEXT,
--     timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
--     estimatedTime INT,
--     FOREIGN KEY (area) REFERENCES areas(id),
--     FOREIGN KEY (createdBy) REFERENCES users(id)
-- );

CREATE TABLE tasks (
    id INT AUTO_INCREMENT PRIMARY KEY,
    requirement INT, -- referencia al requerimiento del que se creo
    team INT, -- referencia al equipo que se involucra
    createdBy INT, -- referencia al usuario que lo creo 
    name VARCHAR(100), -- nombre de la task
    description VARCHAR(255),
    language VARCHAR(40), -- referencia al lenguaje que se utiliza
    framework INT, -- referencia al framework que se utiliza
    estimated_time INT, -- tiempo que dura el desarrollo (horas)
    estimated_cost INT, -- costo en dolares del desarrollo
    ajuste DECIMAL(10,2), -- ajuste
    FOREIGN KEY (requirement) REFERENCES requirements(id),
    FOREIGN KEY (team) REFERENCES teams(id),
    FOREIGN KEY (createdBy) REFERENCES users(id),
    FOREIGN KEY (framework) REFERENCES frameworks(id)
);

-- get Areas
-- get Teams 
-- get 
CREATE VIEW getAreas AS 
    select id, name from areas a;

CREATE VIEW getTeams AS 
    select id, team_name from teams t;



-- Obtener los proyectos activos de un usuario
DELIMITER $$
CREATE PROCEDURE GetActiveProjectsForUser(
    IN userId int)
BEGIN
    select * from projects p
    inner join users u on u.id = p.owner
    where p.status = 1 and p.owner = userId;
END$$
DELIMITER ;


DELIMITER $$
CREATE PROCEDURE GetProjectInfo(
    IN projId int
)
BEGIN
    select projName, company from project p
    where p.id = projId;

END$$
DELIMITER ;



-- Obtener la informacion de un proyecto cuando se le da click por id
-- Se necesitan: requerimientos, tareas, comments
DELIMITER $$
CREATE PROCEDURE GetProjectRequirements(
    IN projectId INT
)
BEGIN 
    SELECT 
        r.id AS requirement_id,
        r.text AS requirement_text,
        r.approved AS requirement_approved,
        r.timestamp AS requirement_timestamp
    FROM 
        requirements r
    WHERE 
        r.projectId = projectId;
END$$
DELIMITER ;


-- get Tasks
DELIMITER $$
CREATE PROCEDURE GetProjectTasks(
    IN projectId INT
)
BEGIN 
    SELECT 
        t.id AS task_id,
        t.title AS task_title,
        t.description AS task_description,
        t.estimatedTime AS task_estimated_time
    FROM 
        tasks t
    WHERE 
        t.area = (SELECT area FROM projects WHERE id = projectId);
END$$
DELIMITER ; 
-- get ActiveProjectsForUser
-- calculate Progress based on stage and requirements
-- get ProjectInfo