package snake

import (
	"fmt"
	"reflect"
	"strings"
)

// sql build where
func BuildAndOrWhere(where map[string]interface{}, andOr string) (whereSQL string, vals []interface{}, err error) {
	for k, v := range where {
		k = strings.Trim(k, " ")
		ks := strings.Split(k, " ")
		if len(ks) > 2 {
			return "", nil, fmt.Errorf("Error in query condition: %s. ", k)
		}
		if whereSQL != "" {
			whereSQL += " " + andOr + " "
		}
		strings.Join(ks, ",")
		switch len(ks) {
		case 1:
			fmt.Println(reflect.TypeOf(v))
			switch v := v.(type) {
			case NullType:
				if v == IsNotNull {
					whereSQL += fmt.Sprint(k, " IS NOT NULL")
				} else {
					whereSQL += fmt.Sprint(k, " IS NULL")
				}
			default:
				t := reflect.TypeOf(v)
				reflectType := fmt.Sprint(t)
				switch reflectType {
				case "string":
					whereSQL += fmt.Sprint(k, "=?")
					vals = append(vals, v)
				case "uint", "uint8", "uint16", "uint32", "uint64", "int", "int8", "int16", "int32", "int64", "float32", "float64", "complex64", "complex128":
					whereSQL += fmt.Sprint(k, "=?")
					vals = append(vals, v)
				case "decimal.Decimal":
					whereSQL += fmt.Sprint(k, "=?")
					vals = append(vals, v)
				}
			}
			break
		case 2:
			k = ks[0]
			switch ks[1] {
			case "=":
				whereSQL += fmt.Sprint(k, "=?")
				vals = append(vals, v)
				break
			case ">":
				whereSQL += fmt.Sprint(k, ">?")
				vals = append(vals, v)
				break
			case ">=":
				whereSQL += fmt.Sprint(k, ">=?")
				vals = append(vals, v)
				break
			case "<":
				whereSQL += fmt.Sprint(k, "<?")
				vals = append(vals, v)
				break
			case "<=":
				whereSQL += fmt.Sprint(k, "<=?")
				vals = append(vals, v)
				break
			case "!=":
				whereSQL += fmt.Sprint(k, "!=?")
				vals = append(vals, v)
				break
			case "<>":
				whereSQL += fmt.Sprint(k, "!=?")
				vals = append(vals, v)
				break
			case "in":
				whereSQL += fmt.Sprint(k, " in (?)")
				vals = append(vals, v)
				break
			case "like":
				whereSQL += fmt.Sprint(k, " like ?")
				vals = append(vals, v)
			}
			break
		}
	}
	whereSQL = "(" + whereSQL + ")"
	return
}

// sql build where
func BuildUpdateField(field map[string]interface{}) (fieldSQL string, vals []interface{}, err error) {
	for k, v := range field {
		switch v := v.(type) {
		case NullType:
			if v == IsNull {
				fieldSQL += fmt.Sprint(k, "=null ")
			}
		default:
			t := reflect.TypeOf(v)
			reflectType := fmt.Sprint(t)
			switch reflectType {
			case "string":
				fieldSQL += fmt.Sprint(k, "=? ")
				vals = append(vals, v)
			case "uint", "uint8", "uint16", "uint32", "uint64", "int", "int8", "int16", "int32", "int64", "float32", "float64", "complex64", "complex128":
				fieldSQL += fmt.Sprint(k, "=? ")
				vals = append(vals, v)
			case "decimal.Decimal":
				fieldSQL += fmt.Sprint(k, "=? ")
				vals = append(vals, v)
			}
		}

	}
	fieldSQL = "(" + fieldSQL + ")"
	return
}
