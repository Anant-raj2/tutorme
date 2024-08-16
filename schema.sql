CREATE TABLE tutors(
  user_id UUID,
  name varchar(25) NOT NULL,
  grade_level int NOT NULL,
  role varchar(10) NOT NULL,
  gender varchar(10) NOT NULL,
  subject varchar(15) NOT NULL,
  PRIMARY KEY (user_id)
);
