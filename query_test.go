package sys

import (
	"fmt"
	"testing"
)

func TestAdd(t *testing.T) {
	sql, args := SqlBoolQuery("A+B,=C+(D+E)", func(placeholder string, exact bool) string {

		// if exact {
		// 	return "(gex.gene_symbol = ? OR gex.ensembl_id = ?)"
		// } else {
		// 	return fmt.Sprintf("(gex.gene_symbol LIKE %s OR gex.ensembl_id LIKE %s)", placeholder, placeholder)
		// }

		return placeholder //fmt.Sprintf("(gex.gene_symbol LIKE %s OR gex.ensembl_id LIKE %s)", placeholder, placeholder)

	})

	//fmt.Println("A+B,=C+(D+E)")
	fmt.Println(sql)
	fmt.Println(args)

	if sql != "((?1 AND ?2) OR (=?3 AND (?4 AND ?5)))" {
		t.Errorf("SqlBoolQuery = %s; want %s", sql, "((?1 AND ?2) OR (?3 AND (?4 AND ?5)))")
	}
}
