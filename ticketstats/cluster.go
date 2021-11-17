package ticketstats

import (
	"fmt"
	"log"
)

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
				log.Println("DEBUG: cluster by clone", issue.Key, "->", clone.Key)
			} else {
				clone.Childs = append(issue.Childs, issue)
				log.Println("DEBUG: cluster by clone", clone.Key, "->", issue.Key)
			}
		}

		for _, duplicateKey := range issue.LinkDuplicates {
			duplicate, ok := keyIndex[duplicateKey]
			if !ok {
				continue
			}
			if issue.Status == "Closed" && duplicate.Status != "Closed" {
				duplicate.Childs = append(issue.Childs, issue)
				log.Println("DEBUG: cluster by duplicate closed 1", duplicate.Key, "->", issue.Key)
			} else if issue.Status != "Closed" && duplicate.Status == "Closed" {
				issue.Childs = append(issue.Childs, duplicate)
				log.Println("DEBUG: cluster by duplicate closed 2", issue.Key, "->", duplicate.Key)
			} else {
				if issue.Created.Before(duplicate.Created) {
					issue.Childs = append(issue.Childs, duplicate)
					log.Println("DEBUG: cluster by duplicate created 1", issue.Key, "->", duplicate.Key)
				} else {
					duplicate.Childs = append(issue.Childs, issue)
					log.Println("DEBUG: cluster by duplicate created 2", duplicate.Key, "->", issue.Key)
				}
			}
		}

		for _, splitKey := range issue.LinkIssueSplits {
			split, ok := keyIndex[splitKey]
			if !ok {
				continue
			}
			issue.Childs = append(issue.Childs, split)
			log.Println("DEBUG: cluster by splits", issue.Key, "->", split.Key)
		}

		for _, partKey := range issue.LinkParts {
			part, ok := keyIndex[partKey]
			if !ok {
				continue
			}
			issue.Childs = append(issue.Childs, part)
			log.Println("DEBUG: cluster by parts", issue.Key, "->", part.Key)
		}

		for _, childKey := range issue.LinkParents {
			child, ok := keyIndex[childKey]
			if !ok {
				continue
			}
			issue.Childs = append(issue.Childs, child)
			log.Println("DEBUG: cluster by parents", issue.Key, "->", child.Key)
		}
	}

	for _, issue := range issues {
		removeDuplicateChildsRecursive(issue, make(map[string]*Issue))
	}

	for _, issue := range issues {
		linkParentsRecursive(issue)
	}
}

func linkParentsRecursive(issue *Issue) {
	for _, child := range issue.Childs {
		child.Parents = append(child.Parents, issue)
		linkParentsRecursive(child)
	}
}

func removeDuplicateChildsRecursive(issue *Issue, set map[string]*Issue) {
	if len(issue.Childs) > 0 {
		for _, child := range issue.Childs {
			removeDuplicateChildsRecursive(child, set)
		}
		issue.Childs = removeDuplicates(issue.Childs, set)
	}
}

func removeDuplicates(issues []*Issue, set map[string]*Issue) []*Issue {
	filtered := make([]*Issue, 0)

	for _, issue := range issues {
		_, ok := set[issue.Key]
		if !ok {
			set[issue.Key] = issue
			filtered = append(filtered, issue)
		}
	}

	return filtered
}

func PrintClusters(issues []*Issue, shorten bool) {
	i := 1
	for _, issue := range issues {
		if len(issue.Parents) > 0 {
			continue
		}

		if len(issue.Childs) > 0 {
			fmt.Println("Cluster", i)
			i++
			fmt.Printf("%s %s [%s]\n", issue.Key, issue.Summary, issue.Status)
			printClustersRecursive(issue.Childs, shorten, "|-")
		}
	}
}

func printClustersRecursive(issues []*Issue, shorten bool, prefix string) {
	for _, issue := range issues {
		fmt.Printf("%s %s %s [%s]\n", prefix, issue.Key, issue.Summary, issue.Status)
		for j, child := range issue.Childs {
			fmt.Printf("%s %s %s [%s]\n", prefix, child.Key, child.Summary, child.Status)
			if shorten && j > 9 {
				fmt.Println(prefix, "   ...", len(issue.Childs), "childs")
				break
			}
			if len(child.Childs) > 0 {
				printClustersRecursive(child.Childs, shorten, "| "+prefix)
			}
		}
	}
}
