package shodan

import (
	"context"
	"log"
	"net/http"

	"github.com/shadowscatcher/shodan"
	"github.com/shadowscatcher/shodan/search"
	"github.com/shadowscatcher/shodan/search/ssl_versions"
)

func Find(domin string, key string) []string {
	var ips []string
	certSearch := search.Params{
		Page: 1,
		Query: search.Query{
			SSL: domin,
			SSLOpts: search.SSLOpts{
				Cert: search.CertOptions{
					Expired: true,
				},
				Version: ssl_versions.TLSv1_2,
			},
		},
	}

	client, _ := shodan.GetClient(key, http.DefaultClient, true)
	ctx := context.Background()
	result, err := client.Search(ctx, certSearch)
	if err != nil {
		log.Fatal(err)
	}

	for _, match := range result.Matches {

		ips = append(ips, match.IPstr)
	}
	for i := 1; i < 5; i++ {
		certSearch.Page++
		result, err := client.Search(ctx, certSearch)
		if err != nil {
			log.Fatal(err)
		}

		for _, match := range result.Matches {
			ips = append(ips, match.IPstr)
		}
	}
	return ips
}
