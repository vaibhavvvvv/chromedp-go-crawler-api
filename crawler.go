package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/chromedp/chromedp"
)

type PartsCSV struct {
	Part        string `json:"part,omitempty"`
	Description string `json:"description,omitempty"`
	List        string `json:"list,omitempty"`
	CorePrice   string `json:"corePricing,omitempty"`
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		var res1, res2, res3, res4 string

		key := strings.TrimSpace(r.URL.Query().Get("key"))

		opts := append(

			chromedp.DefaultExecAllocatorOptions[3:],
			chromedp.NoFirstRun,
			chromedp.NoDefaultBrowserCheck,
		)

		parentCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
		defer cancel()

		ctx, cancel := chromedp.NewContext(parentCtx)
		defer cancel()

		tasks := chromedp.Tasks{
			chromedp.Navigate("http://www.iautoparts.biz/pronto/entrepot/WAW"),
			chromedp.SendKeys("input[name=username]", "WAW21948"),
			chromedp.SendKeys("input[name=password]", "FH4MULUG"),
			chromedp.Submit("input[type=submit]"),
			chromedp.Click("td[class=NavSubTab]", chromedp.NodeVisible),
			chromedp.SendKeys("input[name=pn0]", key),
			chromedp.Submit("input[name=pn0]"),
			chromedp.Text("#idC2C_1100001", &res1, chromedp.NodeVisible),
			chromedp.Text(".PartPrice", &res2, chromedp.NodeVisible),
			chromedp.Text(".CorePrice", &res3, chromedp.NodeVisible),
			chromedp.Text("#idPN_1100001", &res4, chromedp.NodeVisible),
		}

		if err := chromedp.Run(ctx, tasks); err != nil {
			panic(err)
		}

		PartCSV := PartsCSV{}

		PartCSV.Description = strings.TrimSpace(res1)
		PartCSV.List = strings.TrimSpace(res2)
		PartCSV.CorePrice = strings.TrimSpace(res3)
		PartCSV.Part = strings.TrimSpace(res4)

		jsonData, err := json.Marshal(PartCSV)
		if err != nil {
			panic(err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonData)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
