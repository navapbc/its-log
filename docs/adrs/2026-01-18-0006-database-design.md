# Database design

Date: 2026-01-18

## Status

Accepted

## Context

Because of the choice to use SQLite as the primary database engine for `its-log`, it opens up the question of how the databases themselves are organized, along with the tables. `its-log` uses one database per application, which provides secure/inviolate separation of data between applications using the server.

## Design

API keys have metadata associated with them. Clients make their requests using only the key; that unique key is used to look up metadata, including an application ID. The application ID is therefore a property controlled by `its-log`, not the client.

When an event comes in---for example, to the `/v1/se` endpoint storing the values `cart` and `show` (e.g. `/v1/se/cart/show`)---then the `its-log` server will:

1. Look up the application ID associated with the key used (e.g. the key may be associated with the ID `bobs_spatula_store`)
2. Use that ID to construct a database filename based on the current day (e.g. `bobs_spatula_store_2026-01-18.sqlite`)
3. Store to the events table in that database.

The application ID is not chosen by the user of `its-log`; it is stored as a secret by the `its-log` administrators. In this way, it is not a vector for attack (e.g. the application ID cannot be named `/etc/network/rc.d/` or similar), and similarly, the date is not something the user chooses either. Therefore, the fact that each day's worth of data is stored in a separate file is unknown to the user of `its-log`.

```
     AppId: bob                          AppId: alice         
          │                                    │              
          │                                    │              
          ▼                                    ▼              
    ┌─────────────────────────────────────────────────┐       
    │                                                 │       
    │                its-log API server               │       
    │                                                 │       
    └─────┬────────────────────────────────────┬──────┘       
          │                                    │              
    ┌─────▼─────┐      ┌───────────┐      ┌────▼──────┐       
    │           │      │           │      │           │       
    │ SQLite DB │      │ SQLite DB │      │ SQLite DB │       
    │           │      │           │      │           │       
    └───────────┘      └───────────┘      └───────────┘       
  bob_20260118.sqlite                  alice_20260118.sqlite  
```

It is, therefore, *generally* safe to have multiple applications log to the same `its-log` host. However, in compliance-burdened environments, it may be inappropriate for a FISMA Low and FISMA Moderate system to log to the same `its-log` server. Although the data does not mix, it would be something that would need to be discussed with an ISSO/ISSM before proceeding.

## Decision

`its-log` will use a many-databases design, where:

1. Every application logs to its own database
2. A new database is created for every day of data

In this way, 10 applications logging for 10 days will create 100 database files. If this looks to be problematic, `its-log` can consider arranging files under subdirectories based on application IDs.