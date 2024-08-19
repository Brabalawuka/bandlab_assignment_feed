db = db.getSiblingDB(process.env.MONGO_INITDB_DATABASE); // 使用环境变量创建数据库
db.createCollection("mycollection");
