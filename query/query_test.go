package query

import (
	"fmt"
	"testing"
)

func TestAdd(t *testing.T) {
	resp, _ := SqlBoolQuery("A+B,=C+(D+E)", func(placeholder int, matchType MatchType) string {

		// if exact {
		// 	return "(gex.gene_symbol = ? OR gex.ensembl_id = ?)"
		// } else {
		// 	return fmt.Sprintf("(gex.gene_symbol LIKE %s OR gex.ensembl_id LIKE %s)", placeholder, placeholder)
		// }

		return fmt.Sprintf("?%d", placeholder) //fmt.Sprintf("(gex.gene_symbol LIKE %s OR gex.ensembl_id LIKE %s)", placeholder, placeholder)

	})

	//fmt.Println("A+B,=C+(D+E)")
	fmt.Println(resp.Sql)
	fmt.Println(resp.Args...)

	if resp.Sql != "((?1 AND ?2) OR (=?3 AND (?4 AND ?5)))" {
		t.Errorf("SqlBoolQuery = %s; want %s", resp.Sql, "((?1 AND ?2) OR (?3 AND (?4 AND ?5)))")
	}
}
