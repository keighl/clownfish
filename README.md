# Clownfish

[![Build Status](https://travis-ci.org/keighl/clownfish.png?branch=master)](https://travis-ci.org/keighl/clownfish)

Clownfish is a CLI tool for easily adding tables and indices in a RethinkDB database based on a simple YAML input. It traverses your YAML instructions and creates any tables and indices that don't exist yet.

### Usage

```yml
# clownfish.yml

conn:
  # RethinkDB cluster location
  host: localhost:28015

  # Name of the database.
  # Clownfish will create the DB if it doesn't exist already
  db: recipe_app

tables:
  # Name of a table
  users:
    indices:
      # Adds a secondary index on the `email` field
      email:
      api_token:

  recipes:
    # Specify a non-default primary key
    primary_key: name

    indices:
      user_id:
      flavor:
      # Adds a compound-secondary index named `user_flavor` that indexes both `user_id` and `flavor` fields
      user_flavor:
        fields: [user_id, flavor]
```

Run it!

```bash
# assumes `clownfish.yml`
$ clownfish

# or pass a specific filename
$ clownfish rethink_config.yml

# Ouputs
recipe_app: DB created
  users: table created
    indices:
      * email: index created
      * api_token: index created
  recipes: table created
    indices:
      * user_id: index created
      * flavor: index created
      * user_flavor: index created
```

Run it again (and again...) to pick up new tables and indices:

```bash
$ clownfish
recipe_app:
  users:
    indices:
      * email
      * api_token
      * gender: index create
  recipes:
    indices:
      * user_id
      * flavor
      * user_flavor
  ingredients: table created
    indices:
      name: index created
      price: index created
```

### Installation

**Quick install**

```bash
$ curl https://raw.githubusercontent.com/keighl/clownfish/master/install.sh | sudo bash
```

**From source:**

```bash
go install github.com/keighl/clownfish
```


