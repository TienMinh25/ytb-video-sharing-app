CREATE TABLE IF NOT EXISTS videos (
  id INT PRIMARY KEY AUTO_INCREMENT,
  description TEXT,
  upvote INT DEFAULT 0,
  downvote INT DEFAULT 0,
  thumbnail VARCHAR(1024) NOT NULL,
  video_url VARCHAR(1024) NOT NULL,
  account_id INT,
  FOREIGN KEY (account_id) REFERENCES accounts(id) ON DELETE CASCADE
);
