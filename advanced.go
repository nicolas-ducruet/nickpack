package nickpack

import (
	//"github.com/sergi/go-diff/diffmatchpatch"
)

func init() {
}

/*
func findPercentMatch(){
	// Remove all whitespaces and split
	categoryString = strings.ToLower(strings.Trim(categoryString," "))
	categoriesInWords := nickpack.RegSplit(categoryString, "[!@#$%^&*]{1}")

	for _, category := range categories {
		dmp := diffmatchpatch.New()
		categorySlugInWords := nickpack.RegSplit(category.Slug, "[_]{1}")
		for _, wordDB := range categorySlugInWords {
			for _, wordAmazon := range categoriesInWords {
				diffs := dmp.DiffMain(wordDB, wordAmazon, false)
				percentageDiff := float64(dmp.DiffLevenshtein(diffs))/float64(len(wordDB))
				if percentageDiff < 0.25 {
					// UPDATE category for the item
					stmtIns, err := db.Prepare("UPDATE items SET category = ?, updated_at = NOW() WHERE id = ?")
					if err != nil {
						goto C
					}
					defer stmtIns.Close()
					_, err = stmtIns.Exec(category.Id, items[index].Id)
					if err != nil {
						goto C
					}
					items[index].Category = &category.Id
					break
				}
			}
		}
	}
}
*/
