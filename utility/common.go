package utility

import "sort"

//ContribV is a structure to keep users and contribution
type ContribV struct {
	User         string
	Contribution int32
}

//ListSize size of the table printed on command line
var ListSize = 50

//RankContributors sorts the list according to contribution
func RankContributors(in map[string]int32) []ContribV {
	var rankedContribV []ContribV
	for k, v := range in {
		rankedContribV = append(rankedContribV, ContribV{k, v})
	}
	sort.Slice(rankedContribV, func(a, b int) bool {
		return rankedContribV[a].Contribution > rankedContribV[b].Contribution
	})
	return rankedContribV
}
