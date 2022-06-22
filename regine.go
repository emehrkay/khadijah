package khadijah

import "fmt"

func newRegine(matchClause M, rootMaxx *Maxine) *regine {
	return &regine{
		matchClause: matchClause,
		rootMaxx:    rootMaxx,
	}
}

type regine struct {
	matchClause M
	rootMaxx    *Maxine
}

func (r *regine) nodeWithProperties(entity interface{}, label *string) *Maxine {
	maxx := r.rootMaxx.Parse(entity)

	if label == nil {
		label = &maxx.EntityName
	}

	maxx.Query = fmt.Sprintf(`(%s:%s %s)`, maxx.Variable, *label, maxx.CreateQuery)

	return maxx
}

func (r *regine) matchNodeWithMatch(entity interface{}, label *string, matchClause M, withReturn bool) *Maxine {
	maxx := r.rootMaxx.Parse(entity)
	maxx.ParseMatchClause(matchClause)

	if label == nil {
		label = &maxx.EntityName
	}

	maxx.Query = fmt.Sprintf(`MATCH (%s:%s) WHERE %s`, maxx.Variable, *label, maxx.MatchClause)

	if withReturn {
		maxx.Query = fmt.Sprintf(`%s RETURN %s`, maxx.Query, maxx.Variable)
	}

	return maxx
}

func (r *regine) matchNode(entity interface{}, label *string, withReturn bool) *Maxine {
	return r.matchNodeWithMatch(entity, label, r.matchClause, withReturn)
}

// CREATE (x:Label {param: $param}) RETURN x
func (r *regine) createNode(entity interface{}, label *string, withReturn bool, excludes ...string) *Maxine {
	maxx := r.rootMaxx.Parse(entity, excludes...)

	if label == nil {
		label = &maxx.EntityName
	}

	maxx.Query = fmt.Sprintf(`CREATE (%s:%s %s)`, maxx.Variable, *label, maxx.CreateQuery)

	if withReturn {
		maxx.Query = fmt.Sprintf(`%s RETURN %s`, maxx.Query, maxx.Variable)
	}

	return maxx
}

// MERGE (x:Label {param: $param}) SET param1 = $param1 RETURN x
func (r *regine) updateNodeWithMatch(entity interface{}, label *string, matchClause M, withReturn bool, excludes ...string) *Maxine {
	maxx := r.rootMaxx.Parse(entity, excludes...)
	maxx.ParseMatchClause(matchClause)

	if label == nil {
		label = &maxx.EntityName
	}

	maxx.Query = fmt.Sprintf(`MERGE (%s:%s) WHERE %s SET %s`, maxx.Variable, *label, maxx.MatchClause, maxx.SetQuery)

	if withReturn {
		maxx.Query = fmt.Sprintf(`%s RETURN %s`, maxx.Query, maxx.Variable)
	}

	return maxx
}

// MERGE (x:Label {id: $id}) SET param1 = $param1 RETURN x
func (r *regine) updateNode(entity interface{}, label *string, withReturn bool, excludes ...string) *Maxine {
	return r.updateNodeWithMatch(entity, label, r.matchClause, withReturn, excludes...)
}

// MATCH (x {param: $param}) [DETACH] DELETE x
func (r *regine) deleteNodeWithMatch(entity interface{}, detach bool, matchClause M) *Maxine {
	detachClause := " "
	if detach {
		detachClause = " DETACH "
	}

	maxx := r.rootMaxx.Parse(entity)
	maxx.ParseMatchClause(matchClause)

	maxx.Query = fmt.Sprintf(`MATCH (%s) WHERE %s%sDELETE %s`, maxx.Variable, maxx.MatchClause, detachClause, maxx.Variable)
	return maxx
}

// MATCH (x {param: $param}) [DETACH] DELETE x
func (r *regine) detachDeleteNodeWithMatch(entity interface{}, matchClause M) *Maxine {
	return r.deleteNodeWithMatch(entity, true, matchClause)
}

// MATCH (x {param: $param}) DETACH DELETE x
func (r *regine) detachDeleteNode(entity interface{}) *Maxine {
	return r.deleteNodeWithMatch(entity, true, r.matchClause)
}

// MATCH (x {param: $param}) [DETACH] DELETE x
func (r *regine) deleteNode(entity interface{}, detach bool) *Maxine {
	return r.deleteNodeWithMatch(entity, detach, r.matchClause)
}
