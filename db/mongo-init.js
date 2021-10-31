db = db.getSiblingDB("faceit");

db.createUser({
  user: "root",
  pwd: "test",
  roles: [
    {"role": "dbOwner", "db": "faceit"}
  ],
});

db.createCollection("user", {
  validator: {
    $jsonSchema:
    {
      "$schema": "http://json-schema.org/draft-04/schema#",
      "$id": "https://example.com/employee.schema.json",
      "title": "User",
      "description": "This is the schema for the user object.",
      "type": "object",
      "properties": {

        "id": {
          "description": "A UUID for the user",
          "type": "string"
        },

        "first_name": {
          "description": "The user's first name",
          "type": "string"
        },

        "last_name": {
          "description": "The user's last name",
          "type": "string"
        },

        "nickname": {
          "description": "The user's nickname",
          "type": "string"
        },

        "password": {
          "description": "The user's salted and hashed password",
          "type": "string"
        },

        "email": {
          "description": "The user's email address",
          "type": "string"
        },

        "country": {
          "description": "The user's country code",
          "type": "string"
        },

        "created_at": {
          "description": "Entry creation timestamp",
          "type": "string"
        },

        "updated_at": {
          "description": "Entry update timestamp",
          "type": "string"
        }
      }
    }
  }
});
