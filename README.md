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

## Some queries to run

```sql
\timing on
select max(total_cost) from orders;
```

```sql
select * from g_columnar_columns;
```

```sql
 SELECT google_columnar_engine_add(relation => 'acme.public.orders', columns => 'order_id,product_id,customer_id,units,total_cost,deleted_at');
```

# Example where the columnar engine comes into play on a small instance

```sql
SET google_columnar_engine.enable_columnar_scan=off;
select min(total_cost) from orders /*action='reporting-min',application='cmdline',route='non-columnar' */ ;
SELECT SUM(total_cost) AS total_spent FROM orders WHERE deleted_at is NULL /*action='reporting-min',application='cmdline',route='non-columnar' */;
SET google_columnar_engine.enable_columnar_scan=on;
select min(total_cost) from orders /*action='reporting-min',application='cmdline',route='columnar' */ ;
SELECT SUM(total_cost) AS total_spent FROM orders WHERE deleted_at is NULL /*action='reporting-min',application='cmdline',route='columnar' */;
```

# Then try again with the columnar engine on

```sql
SELECT customer_id, SUM(total_cost) AS total_spent FROM orders  WHERE deleted_at is NULL GROUP BY customer_id ORDER BY total_spent DESC LIMIT 10; 
```
