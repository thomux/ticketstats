package ticketstats

import (
	"fmt"
	"log"
)

// ClusterIssues builds a tree for the tickets based on the Jira issue links.
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
			part.Childs = append(issue.Childs, issue)
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

// Clusters returns all issue clusters.
// If needsChilds is true, issue without childs are not part of the result.
func Clusters(issues []*Issue, needsChilds bool) []*Issue {
	clusters := make([]*Issue, 0)
	for _, issue := range issues {
		if needsChilds && len(issue.Childs) == 0 {
			continue
		}
		if len(issue.Parents) == 0 {
			clusters = append(clusters, issue)
		}
	}
	return clusters
}

// linkParentsRecursive creates the backward links for the child issues.
func linkParentsRecursive(issue *Issue) {
	for _, child := range issue.Childs {
		child.Parents = append(child.Parents, issue)
		linkParentsRecursive(child)
	}
}

// removeDuplicateChildsRecursive removes duplicate child links recursive
// childs which are linked form multiple tree levels will be linked to the
// highest tree node.
func removeDuplicateChildsRecursive(issue *Issue, set map[string]*Issue) {
	if len(issue.Childs) > 0 {
		issue.Childs = removeDuplicates(issue.Childs, set)
		for _, child := range issue.Childs {
			removeDuplicateChildsRecursive(child, set)
		}
	}
}

// removeDuplicates reduce the given issue list to contain all issues only
// once (list to set). All issues contained in the given set will be removed
// completely (list set minus given set).
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

// PrintClusters print the "clusters" contained in the given issue list,
// i.e. if a issue has childs, the whole tree is printed, if a issue has
// no childs, it is skipped.
func PrintClusters(issues []*Issue, config Config) {
	i := 1
	for _, issue := range Clusters(OpenTickets(issues, config), true) {
		if len(issue.Parents) > 0 {
			continue
		}

		if len(issue.Childs) > 0 {
			fmt.Println("Cluster", i)
			i++
			fmt.Printf("%s %s [%s]\n", issue.Key, issue.Summary, issue.Status)
			printClustersRecursive(issue.Childs, "|-")
		}
	}
}

// printClustersRecursive prints recursive the "clusters" contained in the given
// issue list.
func printClustersRecursive(issues []*Issue, prefix string) {
	for _, issue := range issues {
		fmt.Printf("%s %s %s [%s]\n", prefix, issue.Key, issue.Summary, issue.Status)
		for _, child := range issue.Childs {
			fmt.Printf("%s %s %s [%s]\n", prefix, child.Key, child.Summary, child.Status)
			if len(child.Childs) > 0 {
				printClustersRecursive(child.Childs, "| "+prefix)
			}
		}
	}
}
