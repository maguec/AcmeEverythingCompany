# Acme Everything Company

This is an example of a Postgres compatable database for Acme Everything Company.
It allows you to create customers, a product catalog, and an order system to simulate database activity.

## Dependencies

This builds a Go binary that can run on a variety of platforms.

## Getting Started

### Setup the config file

```bash
cp sample-config.yml config.yaml
```

Edit the config file to reflect your database credentials.

### Seed the database with Customers and Products

```bash
./Acme --config=config.yaml dbseed
```

### Add Orders

```
./Acme --config=config.yaml addOrders
```

## All the database operations are tagged in Inisights
