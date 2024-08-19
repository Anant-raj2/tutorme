CREATE TABLE tutors(
  user_id UUID Unique,
  name varchar(25) NOT NULL,
  email varchar(25) NOT NULL Unique,
  grade_level int NOT NULL,
  gender varchar(10) NOT NULL,
  subject varchar(15) NOT NULL,
  PRIMARY KEY (user_id)
);
