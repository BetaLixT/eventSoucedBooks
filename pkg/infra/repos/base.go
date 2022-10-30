package repos

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"eventSourcedBooks/pkg/domain/base"

	"github.com/lib/pq"
)

func expandKey(key string, typ string) (string, error) {
	if fields := strings.Split(key, "."); len(fields) < 2 {
		return key, nil
	} else {
		last := len(fields) - 1
		for i := 1; i < last; i++ {
			fields[i] = "'" + fields[i] + "'"
		}
		fields[last] = ">'" + fields[last] + "'"
		nKey := strings.Join(fields, "->")

		if typ == "string" || typ == "[]string" {
			return nKey, nil
		} else if typ == "bool" || typ == "[]bool" {
			return "(" + nKey + ")::bool", nil
		} else if typ == "int" || typ == "[]int" || typ == "float" || typ == "[]float" {
			return "(" + nKey + ")::numeric", nil
		} else if typ == "time" || typ == "[]time" {
			return "(" + nKey + ")::timestamp with time zone", nil
		}

		return "", fmt.Errorf("failed to expand key with type %s", typ)
	}
}

type operationFunc func(key string, p string) string
type parserFunc func(val string) (any, error)

var operationStore = map[string]operationFunc{
	"eq": func(key string, p string) string {
		return key + " = $" + p
	},
	"neq": func(key string, p string) string {
		return key + " != $" + p
	},
	"lt": func(key string, p string) string {
		return key + " < $" + p
	},
	"le": func(key string, p string) string {
		return key + " <= $" + p
	},
	"gt": func(key string, p string) string {
		return key + " > $" + p
	},
	"ge": func(key string, p string) string {
		return key + " >= $" + p
	},
	"like": func(key string, p string) string {
		return key + " LIKE $" + p
	},
	"any": func(key string, p string) string {
		return "$" + p + " <@ " + key
	},
	"all": func(key string, p string) string {
		return "$" + p + " @> " + key
	},
	"in": func(key string, p string) string {
		return key + " = Any($" + p + ")"
	},
	"nin": func(key string, p string) string {
		return key + " != Any($" + p + ")"
	},
}

var parserStore = map[string]parserFunc{
	"bool": func(val string) (any, error) {
		return val == "true", nil
	},
	"int": func(val string) (any, error) {
		if v, err := strconv.Atoi(val); err != nil {
			return nil, fmt.Errorf("failed to convert %s to integer: %w", val, err)
		} else {
			return v, nil
		}
	},
	"float": func(val string) (any, error) {
		if v, err := strconv.ParseFloat(val, 64); err != nil {
			return nil, fmt.Errorf("failed to convert %s to float: %w", val, err)
		} else {
			return v, nil
		}
	},
	"string": func(val string) (any, error) {
		return val, nil
	},
	"time": func(val string) (any, error) {
		if v, err := time.Parse(time.RFC3339, val); err != nil {
			return nil, fmt.Errorf("failed to convert %s to time: %w", val, err)
		} else {
			return v, nil
		}
	},
	"[]bool": func(val string) (any, error) {
		var vs []bool
		if err := json.Unmarshal([]byte(val), &vs); err != nil {
			return nil, fmt.Errorf("failed to convert %s to []bool: %w", val, err)
		} else {
			return vs, nil
		}
	},
	"[]int": func(val string) (any, error) {
		var vs []int32
		if err := json.Unmarshal([]byte(val), &vs); err != nil {
			return nil, fmt.Errorf("failed to convert %s to []int: %w", val, err)
		} else {
			return pq.Int32Array(vs), nil
		}
	},
	"[]float": func(val string) (any, error) {
		var vs []float64
		if err := json.Unmarshal([]byte(val), &vs); err != nil {
			return nil, fmt.Errorf("failed to convert %s to []float: %w", val, err)
		} else {
			return pq.Float64Array(vs), nil
		}
	},
	"[]string": func(val string) (any, error) {
		var vs []string
		if err := json.Unmarshal([]byte(val), &vs); err != nil {
			return nil, fmt.Errorf("failed to convert %s to []string: %w", val, err)
		} else {
			return pq.StringArray(vs), nil
		}
	},
	// This type is very questionable
	"[]time": func(val string) (any, error) {
		var vs []time.Time
		if err := json.Unmarshal([]byte(val), &vs); err != nil {
			return nil, fmt.Errorf("failed to convert %s to []time: %w", val, err)
		} else {
			return val, nil
		}
	},
}

// Builds a condition for a WHERE clause from a Search Query specs
func searchQueryToCond(
	pidx int,
	sq *base.SearchQuery,
) (string, any, error) {

	key, err := expandKey(sq.Key, sq.ValueType)
	if err != nil {
		return "", nil, fmt.Errorf("could not expand key %w", err)
	}
	opFnc, ex := operationStore[sq.Operation]
	if !ex {
		return "", nil, fmt.Errorf("query operation \"%s\" does not exist", sq.Operation)
	}
	prFnc, ex := parserStore[sq.ValueType]
	if !ex {
		return "", nil, fmt.Errorf("query value type \"%s\" does not exist", sq.ValueType)
	}
	param, err := prFnc(sq.Value)
	if err != nil {
		return "", nil, fmt.Errorf("could not parse: %w", err)
	}
	clause := opFnc(key, strconv.Itoa(pidx))
	return clause, param, nil
}

// Builds a WHERE clause from Search Query specs
// Starts placeholder count from pbeg
func buildWhereClause(
	pbeg int,
	sqs []base.SearchQuery,
) (string, []any, error) {
	terms := make([]string, len(sqs))
	params := make([]any, len(sqs))
	var err error

	if len(sqs) == 0 {
		return "", params, nil
	}

	for i := 0; i < len(sqs); i++ {
		terms[i], params[i], err = searchQueryToCond((i + pbeg), &sqs[i])
		if err != nil {
			return "", nil, fmt.Errorf("failed to compile where clause: %w", err)
		}
	}

	return strings.Join(terms, " AND "), params, nil
}

// TODO: Move this somewhere common
func strArrEq(a, b []string) bool {
	if len(a) != len(b) {
		return false
	} else {
		for i, e := range a {
			if b[i] != e {
				return false
			}
		}
	}
	return true
}
