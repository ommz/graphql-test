package repository

import (
	"encoding/json"
	"strconv"

	log "github.com/sirupsen/logrus"
)

func GetGraphQLQueryJSON(n uint64) []byte {
	nStr := strconv.FormatUint(n, 16) // integer n to string conversion. n=5 by default

	queryMap := map[string]string{
		"query": `
			query last_projects($n: Int = ` + nStr + `) {
				projects(last:$n) {
					nodes {
						name,
						description,
						forksCount
					}
				}
			}`,
	}

	log.Info("Marshaling GraphQL query")
	queryJSON, err := json.Marshal(queryMap)
	if err != nil {
		log.Fatal(err)
	}

	return queryJSON
}
