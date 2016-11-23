package recipe

import "regexp"

/*
 * Compile a regexp for parsing Ingredient details. The expected format is
 * similar to that used by the Pyprika library for python. A line like
 *
 * (1 cup) Heavy Cream
 *
 * Should match three groups: (1) (cup) (Heavy Cream) which can then be stored
 * in the Ingredient type with Amount: 1, Unit: cup, Label: Heavy Cream
 */
var ingredientExp = regexp.MustCompile(`^\((?P<Amount>.+)\s(?P<Unit>.+)\)\s+(?P<Label>.*)$`)

// Define a struct for Recipes.
type Recipe struct {
	Name            string   `yaml:"name"`
	Notes           string   `yaml:"notes"`
	Source          string   `yaml:"source"`
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
	/*
	 * Convert the raw list of ingredient strings Unmarshaled from the YAML file
	 * to an array of Ingredient types. If the provided does not match the expected
	 * pattern it will be stored in an Ingredient where the Amount and Unit are
	 * empty, and the whole line will be stored as the Label.
	 */

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
