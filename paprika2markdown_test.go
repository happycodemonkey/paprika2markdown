package main_test

import (
	"io/ioutil"
	"path/filepath"
	"text/template"

	. "github.com/compybara/paprika2markdown/recipe"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Paprika2markdown", func() {
	var (
		templateFile         string
		templateFileFullPath string
		recipe               Recipe
	)

	BeforeEach(func() {
		templateFile = "./templates/template.md"
		templateFileFullPath, _ = filepath.Abs(templateFile)
		recipe = Recipe{
			Name:        "Chicken Nuggets",
			Source:      "The back of the packet",
			PrepTime:    "1 minute",
			CookTime:    "20 minutes",
			Servings:    []int{1, 2},
			Categories:  []string{"Frozen food", "Lazy dinners"},
			Notes:       "How to heat up frozen chicken nuggets.",
			Ingredients: []string{"One packet chicken nuggets", "(1/4 cup) Ketchup"},
			Directions: []string{"Preheat oven to 400 Degrees",
				"Open packet",
				"Put nuggets onto a pan",
				"Place in oven",
				"Bake for 20 minutes",
				"Put onto a plate with your bare hands",
				"Nurse burns caused by impatience",
				"Put ketchup onto plate in a pile next to the nuggets"},
		}
	})

	Describe("Default template", func() {
		Context("should load from file", func() {
			It("should exist", func() {
				Expect(templateFileFullPath).To(BeAnExistingFile())
			})
			It("should parse and execute", func() {
				recipeTemplate, parseErr := template.ParseFiles(templateFileFullPath)
				Ω(parseErr).ShouldNot(HaveOccurred())
				// Make sure it executes and just discard the output.
				execErr := recipeTemplate.Execute(ioutil.Discard, recipe)
				Ω(execErr).ShouldNot(HaveOccurred())
			})
		})
	})
})
