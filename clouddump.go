/*
 * This is an unpublished work copyright 2018 Jens-Uwe Mager
 * 30177 Hannover, Germany, jum@anubis.han.de
 */

package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/datastore"
	"golang.org/x/net/context"
	"google.golang.org/api/iterator"
)

var (
	projectID = flag.String("projectid", "", "google cloud project id")
	kind      = flag.String("kind", "", "table kind to query")
	indent    = flag.Bool("indent", false, "indent json output")
)

type myPlist datastore.PropertyList

func main() {

	flag.Parse()

	ctx := context.Background()

	client, err := datastore.NewClient(ctx, *projectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	q := datastore.NewQuery(*kind)
	it := client.Run(ctx, q)
	for {
		var e myPlist
		var err error
		if _, err = it.Next((*datastore.PropertyList)(&e)); err != nil {
			if err == iterator.Done {
				break
			}
			log.Fatalf("query next: %v\n", err.Error())
		}
		enc := json.NewEncoder(os.Stdout)
		if *indent {
			enc.SetIndent("", "\t")
		}
		if err = enc.Encode(&e); err != nil {
			log.Fatalf("encode crash json: %v", err.Error())
		}
	}
}

func (pl *myPlist) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString("{")
	length := len(*pl)
	count := 0
	for _, e := range *pl {
		jsonValue, err := json.Marshal(e.Value)
		if err != nil {
			return nil, err
		}
		buffer.WriteString(fmt.Sprintf("\"%s\":%s", e.Name, string(jsonValue)))
		count++
		if count < length {
			buffer.WriteString(",")
		}
	}
	buffer.WriteString("}")
	return buffer.Bytes(), nil
}
