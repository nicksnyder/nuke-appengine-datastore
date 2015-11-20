# nuke-appengine-datastore
Implements a module for Google App Engine that exposes a single endpoint to wipe your datastore.

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
