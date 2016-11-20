package recipe

import "regexp"

// Compile a regexp for parsing Ingredient details.
var ingredientExp = regexp.MustCompile(`^\((?P<Amount>.+)\s(?P<Unit>.+)\)\s+(?P<Label>.*)$`)

// Define a struct for Recipes.
type Recipe struct {
	Name            string   `yaml:"name"`
	Notes           string   `yaml:"notes"`
	SourceURL       string   `yaml:"source_url"`
	Servings        []int    `yaml:"servings"`
	PrepTime        string   `yaml:"prep_time"`
	CookTime        string   `yaml:"cook_time"`
	Categories      []string `yaml:"categories"`
	Ingredients     []string `yaml:"ingredients"`
	Directions      []string `yaml:"directions"`
	IngredientsList []Ingredient
}

type Ingredient struct {
	Amount string
	Unit   string
	Label  string
}

func (r *Recipe) ParseIngredientsList() {
	var parsed_ingredients []Ingredient

	for _, ingredient := range r.Ingredients {
		matches := ingredientExp.FindStringSubmatch(ingredient)
		if len(matches) < 4 {
			if len(ingredient) > 0 {
				/*
				 * If the line contains text, but doesn't parse the correct amount of
				 * matches to be in the standard Ingredient format just set the Amount
				 * to empty and store the string as just a Label.
				 */
				parsed_ingredients = append(parsed_ingredients, Ingredient{
					Amount: "",
					Unit:   "",
					Label:  ingredient,
				})
			} else {
				// Empty string, so skip this ingredient.
				continue
			}
		} else {
			parsed_ingredients = append(parsed_ingredients, Ingredient{
				Amount: matches[1],
				Unit:   matches[2],
				Label:  matches[3],
			})
		}
	}
	// Store the parsed ingredient slice in the Ingredients field.
	r.IngredientsList = parsed_ingredients
}
