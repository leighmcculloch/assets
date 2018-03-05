package main

import ishell "gopkg.in/abiosoft/ishell.v2"

func cmdFiles(c *ishell.Context, d Data) Data {
	names := d.Names()
	initIndexes := []int{}
	for i := range names {
		initIndexes = append(initIndexes, i)
	}
	choices := c.Checklist(names, "Files", initIndexes)
	return d.Subset(choices)
}
