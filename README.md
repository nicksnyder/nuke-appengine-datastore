# nuke-appengine-datastore
Implements a Google App Engine module that exposes a single endpoint to wipe your datastore.

## Run locally

```bash
$ goapp serve nuker.yaml
```

## Deploy module

```bash
$ goapp deploy -application YOURAPPID nuker.yaml
```

## Nuke your data

1. Open YOURAPPID.appspot.com/nuker/ in your browser
2. Type in the entity name you want to nuke (or leave empty to nuke all entities)
3. Click the nuke button
4. Wait.

## Why?

Because sometimes you just need to start over.

The benefit with this approach is that if you are operating within free quota (i.e. billing disabled), you can delete about 40x more data than you should be able to given the free quota limits because it takes a while for App Engine to realize you are over quota.

The implementation is pretty dumb. It queries the keys for all the entities in one batch, and then deletes them as fast as possible in batches of 500 (the max size of a delete batch in App Engine) concurrently. The initial query can take a while depending on how much data you have, so I just set the query timeout really long. It is probably also possible to run out of memory.
