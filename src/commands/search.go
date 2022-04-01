package commands

import (
	"fmt"
	"strings"
	"bread/src/helpers/utils"
	"github.com/schollz/progressbar/v3"
	"github.com/microcosm-cc/bluemonday"
)

type SearchCmd struct {
	Name string `arg:"" name:"name" help:"name to search for." type:"string"`
}

func (cmd *SearchCmd) Run(debug bool) (error) {
	var err error
	fmt.Println("Updating Catalog Data...")
	err = utils.FetchAppImageCatalog()
	if err != nil {
		return err
	}

	cmd.Name = strings.ToLower(cmd.Name)

	jsonData, err := utils.ReadAppImageCatalog()
	if err != nil {
		return err
	}

	var foundItems []utils.AppImageFeedItem
	bar := progressbar.Default(
		int64(len(jsonData.Items)),
		"Searching List",
	)

	// This Loop Will Check if the name of description has our search target
	for index := range jsonData.Items {
		item := jsonData.Items[index]
		item.Name = strings.ToLower(item.Name)
		item.Description = strings.ToLower(item.Description)
		if strings.Contains(item.Name, cmd.Name) || strings.Contains(item.Description, cmd.Name) {
			// This loop will loop and check if the provider has a github link or not
			for providerIndex := range item.Links {
				if strings.ToLower(item.Links[providerIndex].Type) != "github" {
					// Finally remove all the html from the description
					p := bluemonday.StripTagsPolicy()
					item.Description = p.Sanitize(item.Description)
					// Get the first line of the description
					item.Description = strings.Split(item.Description, ".")[0]

					// Make the first element the github url
					item.Links[0].Type = "github"
					item.Links[0].Url = strings.ToLower(item.Links[providerIndex].Url)

					// Try to convert the URL to short user/repo format
					githubUserRepo, err := utils.GetUserRepoFromUrl(item.Links[0].Url)
					if err == nil {
						item.Links[0].Url = githubUserRepo
					}
				
					// append it to the foundItems
					foundItems = append(foundItems, item)
					break
				}
			}
		}
		bar.Add(1)
	}

	bar.Finish()

	if len(foundItems) == 0 {
		fmt.Println("Nothing Found in the catalog!")
	} else {
		for foundIndex := range foundItems {
			fmt.Println("\n" + foundItems[foundIndex].Name + " - " + foundItems[foundIndex].Links[0].Url)
			if foundItems[foundIndex].Description != "" {
				fmt.Println("  " + foundItems[foundIndex].Description)
			} else {
				fmt.Println("  No Description provided from Author!")
			}
		}
	}
	return nil
}