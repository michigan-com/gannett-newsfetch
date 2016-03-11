package commands

import (
	"time"

	"github.com/spf13/cobra"

	"github.com/michigan-com/gannett-newsfetch/config"
	api "github.com/michigan-com/gannett-newsfetch/gannettApi/fetch"
)

var articlesCmd = &cobra.Command{
	Use:   "articles",
	Short: "Get articles from the gannett API based on the news source",
	Run:   articleCmdRun,
}

func articleCmdRun(command *cobra.Command, args []string) {
	var apiConfig, _ = config.GetApiConfig()

	for _, code := range apiConfig.SiteCodes {
		api.GetArticlesByDay(code, time.Now())
	}

}
