package khadijah

import "fmt"

func newSynclaire(matchClause M, startVariable, endVariable string) *synclarie {
	return &synclarie{
		matchClause:   matchClause,
		startVariable: startVariable,
		endVariable:   endVariable,
	}
}

type synclarie struct {
	matchClause   M
	startVariable string
	endVariable   string
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
	maxx := RootMaxx.Parse(edge)
	maxx.Query = fmt.Sprintf(`MATCH %s, %s CREATE (%s)%s[%s:%s %s]%s(%s)`,
		nodeStart.Query,
		nodeEnd.Query,
		khadStart.StartVariable,
		dirStart,
		maxx.Variable,
		*edgeLabel,
		maxx.CreateQuery,
		dirEnd,
		khadEnd.EndVariable)

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
	maxx := RootMaxx.Parse(edge)
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

	if withReturn {
		maxx.Query = fmt.Sprintf(`%s RETURN %s, %s, %s`, maxx.Query, khadStart.Variable, maxx.Variable, khadEnd.Variable)
	}

	return maxx
}
