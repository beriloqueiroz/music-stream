db = db.getSiblingDB("musicstream");
db.createCollection("musics");
db.musics.createIndex({ title: "text", artist : "text", "lyrics.text" : "text"});
