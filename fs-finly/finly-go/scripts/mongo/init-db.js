db = db.getSiblingDB("finly");

db.createUser({
  user: "finly-user",
  pwd: "finly-password",
  roles: [
    {
      role: "readWrite",
      db: "finly",
    },
  ],
})