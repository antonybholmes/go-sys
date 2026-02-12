package query

import (
	"fmt"
	"testing"
)

func TestNormalizeQuery(t *testing.T) {
	x := normalizeQuery(`x+(y AND z) OR "a AND b"`)
	fmt.Println(x)

	x = normalizeQuery("( x+ y ) ,c")
	fmt.Println(x)

	// Test cases for normalizeQuery
}

func TestAdd(t *testing.T) {
	resp, _ := SqlBoolQuery("A+B,=C+(D+E)", func(placeholderIndex int, value string, addParens bool) string {

		// if exact {
		// 	return "(gex.gene_symbol = ? OR gex.ensembl_id = ?)"
		// } else {
		// 	return fmt.Sprintf("(gex.gene_symbol LIKE %s OR gex.ensembl_id LIKE %s)", placeholder, placeholder)
		// }

		return fmt.Sprintf("?%d", placeholderIndex) //fmt.Sprintf("(gex.gene_symbol LIKE %s OR gex.ensembl_id LIKE %s)", placeholder, placeholder)

	})

	//fmt.Println("A+B,=C+(D+E)")
	fmt.Println(resp.Sql)

	if resp.Sql != "((?1 AND ?2) OR (=?3 AND (?4 AND ?5)))" {
		t.Errorf("SqlBoolQuery = %s; want %s", resp.Sql, "((?1 AND ?2) OR (?3 AND (?4 AND ?5)))")
	}
}
