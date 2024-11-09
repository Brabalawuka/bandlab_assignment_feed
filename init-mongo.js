db = db.getSiblingDB(process.env.MONGO_INITDB_DATABASE); // 使用环境变量创建数据库
db.createCollection("posts");
db.createCollection("comments");

// A created index for composite key
// Composite key is a combination of commentCount and CommentTime and PostID
db.posts.createIndex({ compositeKey: -1 });