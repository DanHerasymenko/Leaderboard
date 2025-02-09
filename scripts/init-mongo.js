use("leaderboard")

db.createCollection('users');
db.users.createIndex({"id": 1 }, {
    unique: true,
});
db.users.createIndex({"nickname": 1 }, {
    unique: true,
});

db.createCollection('scores');
db.scores.createIndex({"scoreID": 1 }, {
    unique: true,
});
db.scores.createIndex({
    "season": -1,
    "scoreDetails.rating": -1
});