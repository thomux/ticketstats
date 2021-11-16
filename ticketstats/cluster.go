package ticketstats

import "fmt"

func ClusterIssues(issues []*Issue) {
	keyIndex := make(map[string]*Issue)
	idIndex := make(map[string]*Issue)

	for _, issue := range issues {
		keyIndex[issue.Key] = issue
		idIndex[issue.Id] = issue
	}

	for _, issue := range issues {
		if issue.Parent != "" {
			parent, ok := idIndex[issue.Parent]
			if ok {
				parent.Childs = append(parent.Childs, issue)
			}
		}
	}

	for _, issue := range issues {
		for _, cloneKey := range issue.LinkCloners {
			clone, ok := keyIndex[cloneKey]
			if !ok {
				continue
			}
			if issue.Created.Before(clone.Created) {
				issue.Childs = append(issue.Childs, clone)
			} else {
				clone.Childs = append(issue.Childs, issue)
			}
		}

		for _, duplicateKey := range issue.LinkDuplicates {
			duplicate, ok := keyIndex[duplicateKey]
			if !ok {
				continue
			}
			if issue.Status == "Closed" && duplicate.Status != "Closed" {
				duplicate.Childs = append(issue.Childs, issue)
			} else if issue.Status != "Closed" && duplicate.Status == "Closed" {
				issue.Childs = append(issue.Childs, duplicate)
			} else {
				if issue.Created.Before(duplicate.Created) {
					issue.Childs = append(issue.Childs, duplicate)
				} else {
					duplicate.Childs = append(issue.Childs, issue)
				}
			}
		}

		for _, splitKey := range issue.LinkIssueSplits {
			split, ok := keyIndex[splitKey]
			if !ok {
				continue
			}
			issue.Childs = append(issue.Childs, split)
		}

		for _, partKey := range issue.LinkParts {
			part, ok := keyIndex[partKey]
			if !ok {
				continue
			}
			issue.Childs = append(issue.Childs, part)
		}

		for _, parentKey := range issue.LinkParents {
			parent, ok := keyIndex[parentKey]
			if !ok {
				continue
			}
			parent.Childs = append(issue.Childs, issue)
		}
	}
}

func PrintClusters(issues []*Issue, shorten bool) {
	for i, issue := range issues {
		if len(issue.Childs) > 0 {
			fmt.Println("Cluster", i)
			fmt.Println(issue.Key, issue.Summary, issue.Type)
			for j, child := range issue.Childs {
				fmt.Println("|-", child.Key, child.Summary, child.Type)
				if shorten && j > 9 {
					fmt.Println("   ...", len(issue.Childs), "childs")
					break
				}
			}
			fmt.Println("-----------------")
		}
	}
}
