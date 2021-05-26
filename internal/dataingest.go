package internal

import (
	"context"
	"encoding/json"
	"log"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
)

func DataIngest(issues []Issue) {
	log.Println("Ingestão de dados")

	ctx := context.Background()

	client, err := storage.NewClient(ctx)

	defer client.Close()

	if err != nil {
		panic(err.Error())
	}

	bkt := client.Bucket("dispatcher-opencalls")

	query := &storage.Query{Prefix: ""}

	var names []string
	it := bkt.Objects(ctx, query)

	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}

		if err != nil {
			panic(err.Error())
		}

		names = append(names, attrs.Name)
	}

	for _, issue := range issues {
		if len(names) == 0 {
			obj := bkt.Object(issue.ID)
			w := obj.NewWriter(ctx)

			b, errJson := json.Marshal(issue)

			if errJson != nil {
				panic(errJson)
			}

			result, errorWriter := w.Write(b)

			log.Println(result)
			log.Println(errorWriter)

			if errCloseWriter := w.Close(); errCloseWriter != nil {
				panic(errCloseWriter)
			}

			//log.Println(result, err)
		}
	}
}
