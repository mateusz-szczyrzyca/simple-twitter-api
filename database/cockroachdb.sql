DROP DATABASE IF EXISTS twitter;
CREATE DATABASE twitter;
USE TWITTER;

DROP TABLE IF EXISTS messages;
DROP TABLE IF EXISTS users;
CREATE TABLE IF NOT EXISTS users (
  uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  username STRING NOT NULL,
  password STRING NOT NULL,
  status STRING NOT NULL,
  token STRING NULL
);

INSERT INTO users(username, password, status, token) VALUES('admin1', '2c0ca499381411013c69e56d708e14fce4995a5bc58e08c5d29aae07a0967ee8de5df8bb7c795c16727e86a6b22cedf706c7c668c0c035271ea82ed123c43485', 'u', NULL);
INSERT INTO users(username, password, status, token) VALUES('admin2', '8c48ccf0b18dacd736b2f8452d27697a8e2477aad71618b0d54b1eb9beebc06f86fc86849f8af5c6c2c3518228c36660b4735d05ebf97f95e6229fad1a6dfafc', 'a', NULL);
INSERT INTO users(username, password, status, token) VALUES('admin3', '1bb9d6ff045e018ee29e1d107e7ef8c12b0a419599408e1d3196d2725ada04223bce6b93671cafc585d123e12dce2baec8999d9133b0786d21171ca92bd494df', 'u', NULL);


CREATE TABLE IF NOT EXISTS messages (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  userid UUID NOT NULL REFERENCES users(uuid) ON DELETE CASCADE,
  datetime TIMESTAMP,
  tags STRING NULL,
  text STRING
);

INSERT INTO messages(userid, datetime, tags, text) VALUES((select uuid from users limit 1), TIMESTAMP '2015-01-25 10:10:10.555555', 
'tag1,tag2','Content of message 1');
INSERT INTO messages(userid, datetime, tags, text) VALUES((select uuid from users limit 1), TIMESTAMP '2016-01-26 10:10:10.555555', 
'tag1,tag2','Content of message 2');
INSERT INTO messages(userid, datetime, tags, text) VALUES((select uuid from users limit 1), TIMESTAMP '2017-07-26 10:10:10.555555', 
'tag1,tag2','Content of message 3');
INSERT INTO messages(userid, datetime, tags, text) VALUES((select uuid from users limit 1), TIMESTAMP '2017-01-15 10:10:10.555555', 
'tag3,tag4,tag5','Content of message 4');
INSERT INTO messages(userid, datetime, tags, text) VALUES((select uuid from users limit 1), TIMESTAMP '2018-07-26 10:10:10.555555', 
'tag3','Content of message 5');
INSERT INTO messages(userid, datetime, tags, text) VALUES((select uuid from users limit 1), TIMESTAMP '2016-10-25 10:10:10.555555', 
'','Content of message 6');
INSERT INTO messages(userid, datetime, tags, text) VALUES((select uuid from users limit 1), TIMESTAMP '2016-12-26 10:10:10.555555', 
'tag5','Content of message 7');
INSERT INTO messages(userid, datetime, tags, text) VALUES((select uuid from users limit 1), TIMESTAMP '2016-06-06 10:10:10.555555', 
'tag5','Content of message 8');
INSERT INTO messages(userid, datetime, tags, text) VALUES((select uuid from users limit 1), TIMESTAMP '2019-03-26 10:10:10.555555', 
'tag2019','Content of message 9');
INSERT INTO messages(userid, datetime, tags, text) VALUES((select uuid from users limit 1), TIMESTAMP '2019-01-12 10:10:10.555555', 
'tag2019','Content of message 10');
INSERT INTO messages(userid, datetime, tags, text) VALUES((select uuid from users limit 1), TIMESTAMP '2015-01-25 10:10:10.555555', 
'tag1','Content of message 11');
INSERT INTO messages(userid, datetime, tags, text) VALUES((select uuid from users limit 1), TIMESTAMP '2016-01-26 10:10:10.555555', 
'tag2','Content of message 12');
INSERT INTO messages(userid, datetime, tags, text) VALUES((select uuid from users limit 1), TIMESTAMP '2017-09-29 10:10:10.555555', 
'tag2','Content of message 13');
INSERT INTO messages(userid, datetime, tags, text) VALUES((select uuid from users limit 1), TIMESTAMP '2015-12-12 10:10:10.555555', 
'tag1,tag5','Content of message 14');
INSERT INTO messages(userid, datetime, tags, text) VALUES((select uuid from users limit 1), TIMESTAMP '2015-01-26 10:10:10.555555', 
'tag100','Content of message 15');
INSERT INTO messages(userid, datetime, tags, text) VALUES((select uuid from users limit 1), TIMESTAMP '2016-11-25 10:10:10.555555', 
'','Content of message 16');
INSERT INTO messages(userid, datetime, tags, text) VALUES((select uuid from users limit 1), TIMESTAMP '2016-12-24 10:10:10.555555', 
'tag1,tag2,tag3,tag4,tag5','Content of message 17');
INSERT INTO messages(userid, datetime, tags, text) VALUES((select uuid from users limit 1), TIMESTAMP '2018-01-01 10:10:10.555555', 
'tag2018','Content of message 18');
INSERT INTO messages(userid, datetime, tags, text) VALUES((select uuid from users limit 1), TIMESTAMP '2018-03-26 10:10:10.555555', 
'tag2018','Content of message 19');
INSERT INTO messages(userid, datetime, tags, text) VALUES((select uuid from users limit 1), TIMESTAMP '2018-01-12 10:10:10.555555', 
'tag2018','Content of message 20');


