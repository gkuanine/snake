package snake

import (
	"fmt"
	"testing"
)

func TestPageDelete(t *testing.T) {
	var w = make(map[string]interface{})
	w["id"] = 1
	w["stock"] = 10
	w["version"] = 2
	w["name like "] = "%手机%"
	updateSnake := NewDeleteSnake().Table("t1").Where(w).BuildSql()
	fmt.Println(updateSnake.lastSql)
	fmt.Println(updateSnake.lastParams)

}
