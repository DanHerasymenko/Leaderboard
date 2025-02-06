use("leaderboard")

db.createCollection('users');
db.users.createIndex({"id": 1 }, {
    unique: true,
});
db.users.createIndex({"name": 1 }, {
    unique: true,
});

db.createCollection('scores');
db.messages.createIndex({"id": 1 }, {
    unique: true,
});
db.messages.createIndex({
    "channel_name": -1,
    "id": -1
})