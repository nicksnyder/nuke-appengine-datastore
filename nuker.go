package nuke

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"appengine"
	"appengine/datastore"
)

func init() {
	http.HandleFunc("/nuker/", handleNuke)
	http.HandleFunc("/nuker/put", handlePut)
}

func handleNuke(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		renderNukeConfirmation(w, r)
	} else {
		nukeDatastore(w, r)
	}
}

func renderNukeConfirmation(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`
	<!DOCTYPE html>
	<html>
	<head>
	<style>
	  input {
		  display: block;
			margin: 0px auto;
			font-size: 72px;
			text-align: center;
	  }
	</style>
	</head>
	<body>
	<form action="/nuker/nuke" method="POST">
	<input name="kind" type="text" placeholder="Entity name"/>
	<input value="Nuke" type="submit" />
	</form>
	</body>
	</html>
		`))
}

func nukeDatastore(w http.ResponseWriter, r *http.Request) {
	c := appengine.Timeout(appengine.NewContext(r), 1*time.Hour)
	kind := r.FormValue("kind")
	c.Debugf("starting query")
	q := datastore.NewQuery(kind).KeysOnly().Limit(-1)
	keys, err := q.GetAll(c, nil)
	if err != nil {
		c.Errorf("query failed: %s", err)
		return
	}
	c.Debugf("query finished %d", len(keys))
	var wg sync.WaitGroup
	count := len(keys)
	batchSize := 500
	for i := 0; i < count; i += batchSize {
		wg.Add(1)
		j := i + batchSize
		if j > count {
			j = count
		}
		c.Infof("will delete keys [%d:%d]", i, j)
		go func(i, j int) {
			defer wg.Done()
			c.Infof("starting delete [%d:%d]", i, j)
			if err := datastore.DeleteMulti(c, keys[i:j]); err != nil {
				c.Errorf("failed deleting keys [%d:%d]: %s", i, j, err)
			} else {
				c.Infof("deleted keys [%d:%d]", i, j)
			}
		}(i, j)
	}

	c.Debugf("waiting...")
	wg.Wait()
	w.Write([]byte(fmt.Sprintf("nuked %s", kind)))
}

func handlePut(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	var keys []*datastore.Key
	count := 2000
	for i := 0; i < count; i++ {
		keys = append(keys, datastore.NewKey(c, "Entity", "", 0, nil))
	}

	type Entity struct {
		Flag bool
		Name string
	}
	entities := make([]Entity, count)
	if _, err := datastore.PutMulti(c, keys, entities); err != nil {
		c.Errorf(err.Error())
		http.Error(w, err.Error(), 500)
	}
}
