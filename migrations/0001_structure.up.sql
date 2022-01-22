CREATE TABLE folder (
    id INT PRIMARY KEY,
    parent INT,
    name TEXT,
    priority INT,
    full_path TEXT,
    CONSTRAINT folder_parent
      FOREIGN KEY(parent) 
	  REFERENCES folder(id)
);