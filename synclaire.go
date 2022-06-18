package khadijah

import "fmt"

func newSynclaire(matchClause M, startVariable, endVariable string, rootMaxx *Maxine) *synclarie {
	return &synclarie{
		matchClause:   matchClause,
		startVariable: startVariable,
		endVariable:   endVariable,
		rootMaxx:      rootMaxx,
	}
}

type synclarie struct {
	matchClause   M
	startVariable string
	endVariable   string
	rootMaxx      *Maxine
}

func (s *synclarie) getDirection(direction string) (dirStart, dirEnd string) {
	switch direction {
	case "in":
		dirStart = "<-"
		dirEnd = "-"

	case "out":
		dirStart = "-"
		dirEnd = "->"

	default:
		dirStart = "-"
		dirEnd = "-"
	}

	return dirStart, dirEnd
}

func (s *synclarie) createEdgeWithMatches(start interface{}, startLabel *string, startMatchClause M, direction string, end interface{}, endLabel *string, endMatchClause M, edge interface{}, edgeLabel *string, withReturn bool, excldues ...string) *Maxine {
	khadStart := New(
		SetVariable("start"),
		SetParamPrefix("start_"),
		SetMatchClause(startMatchClause),
	)
	khadEnd := New(
		SetVariable("end"),
		SetParamPrefix("end_"),
		SetMatchClause(endMatchClause),
	)
	nodeStart := khadStart.MatchNode(start, startLabel, false)
	nodeEnd := khadEnd.MatchNode(end, endLabel, false)
	dirStart, dirEnd := s.getDirection(direction)
	maxx := s.rootMaxx.Parse(edge)
	labelEdge := maxx.EntityName
	if edgeLabel != nil {
		labelEdge = *edgeLabel
	}

	maxx.Query = fmt.Sprintf(`%s %s CREATE (%s)%s[%s:%s %s]%s(%s)`,
		nodeStart.Query,
		nodeEnd.Query,
		khadStart.StartVariable,
		dirStart,
		maxx.Variable,
		labelEdge,
		maxx.CreateQuery,
		dirEnd,
		khadEnd.EndVariable)

	maxx.MergeParams(nodeStart.Params, nodeEnd.Params)

	if withReturn {
		maxx.Query = fmt.Sprintf(`%s RETURN %s, %s, %s`, maxx.Query, khadStart.Variable, maxx.Variable, khadEnd.Variable)
	}

	return maxx
}

func (s *synclarie) updateEdgeWithMatches(start interface{}, startLabel *string, startMatchClause M, direction string, end interface{}, endLabel *string, endMatchClause M, edge interface{}, edgeLabel *string, edgeMatchClause M, withReturn bool, excldues ...string) *Maxine {
	khadStart := New(
		SetVariable("start"),
		SetParamPrefix("start"),
		SetMatchClause(startMatchClause),
	)
	khadEnd := New(
		SetVariable("end"),
		SetParamPrefix("end"),
		SetMatchClause(endMatchClause),
	)
	nodeStart := khadStart.NodeWithProperties(start, startLabel)
	nodeEnd := khadEnd.NodeWithProperties(start, endLabel)
	dirStart, dirEnd := s.getDirection(direction)
	maxx := s.rootMaxx.Parse(edge)
	maxx.Query = fmt.Sprintf(`MATCH %s, %s MERGE (%s)%s[%s:%s %s]%s(%s) SET %s`,
		nodeStart.Query,
		nodeEnd.Query,
		khadStart.StartVariable,
		dirStart,
		maxx.Variable,
		*edgeLabel,
		edgeMatchClause,
		dirEnd,
		khadEnd.EndVariable,
		maxx.SetQuery,
	)

	maxx.MergeParams(nodeStart.Params, nodeEnd.Params)

	if withReturn {
		maxx.Query = fmt.Sprintf(`%s RETURN %s, %s, %s`, maxx.Query, khadStart.Variable, maxx.Variable, khadEnd.Variable)
	}

	return maxx
}
