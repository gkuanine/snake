package snake

import (
	"fmt"
	"testing"
)

func TestPageUpdate(t *testing.T) {
	var stock = 10
	var field = make(map[string]interface{})
	field["name"] = "xiaoming"
	field["stock"] = stock - 1
	field["version"] = 3
	field["ff"] = IsNull
	var w = make(map[string]interface{})
	w["id"] = 1
	w["stock"] = stock
	w["version"] = 2
	updateSnake := NewUpdateSnake().Table("t1").Update(field).Where(w).BuildSql()
	fmt.Println(updateSnake.lastSql)
	fmt.Println(updateSnake.lastParams)

}
