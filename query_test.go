package sys

import (
	"fmt"
	"testing"
)

func TestAdd(t *testing.T) {
	sql, args := SqlBoolQuery("A+B,C+(D+E)", func(placeholder string, exact bool) string {

		// if exact {
		// 	return "(gex.gene_symbol = ? OR gex.ensembl_id = ?)"
		// } else {
		// 	return fmt.Sprintf("(gex.gene_symbol LIKE %s OR gex.ensembl_id LIKE %s)", placeholder, placeholder)
		// }

		return fmt.Sprintf("(gex.gene_symbol LIKE %s OR gex.ensembl_id LIKE %s)", placeholder, placeholder)
	})

	fmt.Println(sql)
	fmt.Println(args)

	if sql != "A" {
		t.Errorf("SqlBoolQuery = %s; want %s", sql, "A+B,C+(D+E)")
	}
}
