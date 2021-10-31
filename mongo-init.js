import * as schema from "./schema.json";

db = db.getSiblingDB("faceit");

db.createUser({
  user: "root",
  pwd: "test",
  roles: [
    "root"
  ],
});

db.createCollection("user", {
  validator: {$jsonSchema: schema}
});
