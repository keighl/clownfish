# clownfish

[![Build Status](https://travis-ci.org/keighl/clownfish.png?branch=master)](https://travis-ci.org/keighl/clownfish) [![Coverage Status](https://coveralls.io/repos/keighl/clownfish/badge.svg)](https://coveralls.io/r/keighl/clownfish)

Clownfish is a CLI tool for quickly adding tables and indices to a RethinkDB database based on a simple YAML input. Useful for deployment, scaffolding and migrations.

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
      # Adds a comppund-secondary index named `user_flavor that indexes both `user_id` and `flavor` fields
      user_flavor:
        fields: [user_id, flavor]

  restaurants:
    # Specify indices you want REMOVED from a table
    absent_indices:
      - sanitation_rating
      - minimum_wage

# Specify tables you want REMOVED from the database
absent_tables:
  - cheeses
  - cakes
  - soups
```

Run it!

```bash
# assumes `clownfish.yml`
$ clownfish

# or pass a specific filename
$ clownfish rethink_config.yml

# Ouputs
+ recipe_app
  + users
      + email
      + api_token
  + recipes
      + user_id
      + flavor
      + user_flavor
  + restaurants
      - sanitation_rating
      - minimum_wage
  - cheeses
  - cakes
  - soups
```

Run it again (and again...) to pick up new tables and indices:

### Installation

Find the right binary for your system on the [releases page.](https://github.com/keighl/clownfish/releases/latest)

```bash
// Download it
curl -O -L https://github.com/keighl/clownfish/releases/download/0.1.0/clownfish-linux-amd64.tgz

// Extract it
tar xzvf clownfish-linux-amd64.tgz

// Install it
sudo mv clownfish-linux-amd64 /usr/local/bin/clownfish
```

**If you have `go` installed:**

```bash
go install github.com/keighl/clownfish
```


