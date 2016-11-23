package recipe_test

import (
	. "github.com/compybara/paprika2markdown/recipe"
	yaml "gopkg.in/yaml.v2"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Recipe", func() {
	var (
		recipeDefine Recipe
		recipeImport Recipe
		recipeYAML   []byte
	)

	BeforeEach(func() {
		// Define a recipe object manually.
		recipeDefine = Recipe{
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

		recipeYAML = []byte("---\n" +
			"name: Ham and Swiss Sandwich\n" +
			"servings: [1]\n" +
			"source: The Earl of Sandwich\n" +
			"source_url:\n" +
			"prep_time: 5 minutes\n" +
			"cook_time: none\n" +
			"categories:\n" +
			"  - sandwiches\n" +
			"  - lunch\n" +
			"notes: A simple ham & cheese sandwich.\n" +
			"ingredients:\n" +
			"  - Two slices of bread\n" +
			"  - One slice swiss cheese\n" +
			"  - Two slices ham\n" +
			"  - (1 tsp) Mustard\n" +
			"directions:\n" +
			"  - Spread mustard onto the bread.\n" +
			"  - Place ham and cheese slices on top of one bread slice.\n" +
			"  - Place second bread slice on top of ham and cheese.")
	})

	Describe("Manually created", func() {
		Context("should have the correct data", func() {
			It("should be chicken nuggets", func() {
				Expect(recipeDefine.Name).To(Equal("Chicken Nuggets"))
			})

			It("should have an integer 1-2 Servings.", func() {
				Expect(recipeDefine.Servings).Should(ConsistOf([]int{1, 2}))
			})

			It("should be frozen food", func() {
				Expect(recipeDefine.Categories).To(ContainElement("Frozen food"))
			})

			It("should successfully parse ingredients", func() {
				recipeDefine.ParseIngredientsList()
				Expect(recipeDefine.IngredientsList).ShouldNot(BeEmpty())
				Expect(recipeDefine.IngredientsList[1].Amount).To(Equal("1/4"))
				Expect(recipeDefine.IngredientsList[1].Unit).To(Equal("cup"))
			})
		})
	})

	Describe("YAML import", func() {
		Context("Should import a recipe", func() {
			It("should be a ham and cheese sandwich", func() {
				err := yaml.Unmarshal(recipeYAML, &recipeImport)
				Ω(err).ShouldNot(HaveOccurred())
				Expect(recipeImport.Name).To(Equal("Ham and Swiss Sandwich"))
			})

			It("should only be one serving", func() {
				Expect(recipeImport.Servings[0]).To(Equal(1))
			})

			It("should not have a source url", func() {
				Expect(recipeImport.SourceURL).To(BeEmpty())
			})

			It("should have three steps in its directions.", func() {
				Expect(recipeImport.Directions).To(HaveLen(3))
			})

			It("should parse the ingredients list", func() {
				recipeImport.ParseIngredientsList()
				Expect(recipeImport.IngredientsList).To(HaveLen(4))
				Expect(recipeImport.IngredientsList[0].Label).To(Equal("Two slices of bread"))
				Expect(recipeImport.IngredientsList[3].Amount).To(Equal("1"))
				Expect(recipeImport.IngredientsList[3].Unit).To(Equal("tsp"))
			})

		})
	})
})
