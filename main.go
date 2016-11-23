package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	recipe "github.com/compybara/paprika2markdown/recipe"
	yaml "gopkg.in/yaml.v2"
)

// Use constants for default values
const defaultTemplate string = "./templates/template.md"
const defaultRecipes string = "./recipes"
const defaultOutput string = "./recipes.wiki"
const pathSeparator = string(os.PathSeparator)

// Instantiate variables to store command-line flag values.
var templateFile string
var recipesFolder string
var outputFolder string

func init() {
	// Setup the command-line flags.
	flag.StringVar(&templateFile, "t", defaultTemplate,
		"The template to use for generating output.")

	flag.StringVar(&recipesFolder, "r", defaultRecipes,
		"The folder containing the recipe YAML files.")

	flag.StringVar(&outputFolder, "o", defaultOutput, "Folder to store the"+
		" generated recipes. Existing files will be overwritten.")
}

func promptYesNo(question string) (answer bool) {
	// Ask the user for input for a Yes or No question. Returns true if the user
	// responds with "y, Y, yes, Yes", false for anything else.
	fmt.Printf("%v (y/N): ", question)
	var response string
	fmt.Scanln(&response)
	if strings.ToLower(string(response[0])) == "y" {
		return true
	}
	return false
}

func getYAMLFiles(path string, files []os.FileInfo) (yamlFiles []string) {
	// Take the list of files from ioutil.ReadDir and filter out only the
	// non-empty YAML files. Return a slice of strings with the full path.
	for _, file := range files {
		// Split the filename at the . and check if the last element is a correct
		// file extension (yml or yaml).
		split := strings.Split(file.Name(), ".")
		extension := strings.ToLower(split[len(split)-1])
		if extension == "yml" || extension == "yaml" {
			if file.Size() > 0 {
				fullPath := fmt.Sprintf("%s%s%s", path, pathSeparator, file.Name())
				yamlFiles = append(yamlFiles, fullPath)
			}
		}
	}
	return yamlFiles
}

func parseRecipes(yamlFiles []string) (recipes []recipe.Recipe) {
	// Read in each recipe file and parse to a recipe.Recipe struct.
	// Return a slice containing the Recipe objects.
	for _, file := range yamlFiles {
		data, err := ioutil.ReadFile(file)
		if err != nil {
			fmt.Printf("Cannot read recipe file %v, skipping.\n%v\n", file, err)
			continue

		}
		var r recipe.Recipe
		if err := yaml.Unmarshal(data, &r); err != nil {
			fmt.Printf("Error parsing YAML for %v\n%v\n", file, err)
			continue
		}
		r.ParseIngredientsList()
		recipes = append(recipes, r)
	}
	return recipes
}

func main() {
	flag.Parse()

	// Get the absolute path for each option and make sure they exist.
	templateFileFullPath, _ := filepath.Abs(templateFile)
	recipeTemplate := template.Must(template.ParseFiles(templateFileFullPath))
	fmt.Printf("Using template file %v\n", templateFile)

	recipesFolder, _ = filepath.Abs(recipesFolder)
	if _, err := os.Stat(recipesFolder); os.IsNotExist(err) {
		fmt.Printf("Recipes folder %v not found.\n", recipesFolder)
		os.Exit(1)
	}

	// If output folder doesn't exist ask the user if it can be created.
	outputFolderFullPath, _ := filepath.Abs(outputFolder)
	if _, err := os.Stat(outputFolderFullPath); os.IsNotExist(err) {
		q := "Output folder does not exist. Should I create it?"
		if promptYesNo(q) {
			fmt.Printf("Creating folder %v\n", outputFolder)
			mkdirErr := os.MkdirAll(outputFolderFullPath, 0755)
			if mkdirErr != nil {
				fmt.Printf("Error creating output folder: %v\n", mkdirErr)
				os.Exit(1)
			}
			fmt.Printf("Output folder created.\n")
		} else {
			fmt.Printf("Output folder not created. Quitting.")
			os.Exit(0)
		}
	}

	// Get all of the YAML files in the recipes folder.
	files, err := ioutil.ReadDir(recipesFolder)
	if err != nil {
		fmt.Printf("Problem reading recipes folder: %v\n", err)
	}
	fmt.Printf("Reading recipes from %v\n", recipesFolder)
	yamlFiles := getYAMLFiles(recipesFolder, files)
	recipes := parseRecipes(yamlFiles)

	// Generate a file from the template for each recipe.
	for _, recipe := range recipes {
		// Use recipe name .md for the output filename.
		outputFileName := fmt.Sprintf("%s.md", recipe.Name)
		outputFileFullPath := fmt.Sprintf("%s%s%s", outputFolder, pathSeparator, outputFileName)

		// Create the output file. Existing files will be truncated. Skip this recipe
		// if for some reason the output file cannot be opened. Maybe consider exiting
		// with an error here instead?
		outputFile, err := os.Create(outputFileFullPath)
		if err != nil {
			fmt.Printf("Error opening output file %v for writing: %v\n", outputFileName, err)
			continue
		}
		fmt.Printf("Saving recipe for %v to %v\n", recipe.Name, outputFileName)
		templateErr := recipeTemplate.Execute(outputFile, recipe)
		if templateErr != nil {
			fmt.Printf("Error executing template for %v\n%v\n", outputFileName, templateErr)
			continue
		}
	}
}
