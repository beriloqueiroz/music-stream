db = db.getSiblingDB("music-stream");
db.createCollection("musics");
db.musics.createIndex({ title: "text", artist : "text", "lyrics.text" : "text"});
