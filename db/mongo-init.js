db = db.getSiblingDB("faceit");

db.createUser({
  user: "owner",
  pwd: "test",
  roles: [
    {role: "dbOwner", db: "faceit"}
  ],
});

db.createCollection("user", {
  validator: {
    $jsonSchema:
    {
      bsonType: "object",
      required: ["id", "first_name", "last_name", "nickname", "password", "email", "country", "created_at", "updated_at"],
      properties: {

        id: {
          description: "A UUID for the user",
          bsonType: "string"
        },

        first_name: {
          description: "The user's first name",
          bsonType: "string"
        },

        last_name: {
          description: "The user's last name",
          bsonType: "string"
        },

        nickname: {
          description: "The user's nickname",
          bsonType: "string"
        },

        password: {
          description: "The user's salted and hashed password",
          bsonType: "string"
        },

        email: {
          description: "The user's email address",
          bsonType: "string"
        },

        country: {
          description: "The user's country code",
          bsonType: "string"
        },

        created_at: {
          description: "Entry creation timestamp",
          bsonType: "string"
        },

        updated_at: {
          description: "Entry update timestamp",
          bsonType: "string"
        }
      }
    }
  }
});
