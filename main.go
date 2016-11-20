package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

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
	fmt.Println(fmt.Sprintf("%s (y/N): ", question))
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

func main() {
	flag.Parse()

	// Get the absolute path for each option and make sure they exist.
	templateFile, _ = filepath.Abs(templateFile)
	if _, err := os.Stat(templateFile); os.IsNotExist(err) {
		fmt.Println(fmt.Sprintf("Template file %s does not exist.", templateFile))
		os.Exit(1)
	}

	recipesFolder, _ = filepath.Abs(recipesFolder)
	if _, err := os.Stat(recipesFolder); os.IsNotExist(err) {
		fmt.Println(fmt.Sprintf("Recipes folder %s not found.", recipesFolder))
		os.Exit(1)
	}

	// If output folder doesn't exist ask the user if it can be created.
	outputFolder, _ = filepath.Abs(outputFolder)
	if _, err := os.Stat(outputFolder); os.IsNotExist(err) {
		q := "Output folder does not exist. Should I create it?"
		if promptYesNo(q) {
			fmt.Println(fmt.Sprintf("Creating folder %s", outputFolder))
			mkdirErr := os.Mkdir(outputFolder, 0755)
			if mkdirErr != nil {
				fmt.Println(fmt.Sprintf("Error creating output folder: %s", mkdirErr))
				os.Exit(1)
			}
			fmt.Println("Output folder created.")
		} else {
			fmt.Println("Output folder not created. Quitting.")
			os.Exit(0)
		}
	}

	// Get all of the YAML files in the recipes folder.
	files, err := ioutil.ReadDir(recipesFolder)
	if err != nil {
		fmt.Println(fmt.Sprintf("Problem reading recipes folder: %s", err))
	}
	yamlFiles := getYAMLFiles(recipesFolder, files)
	for _, file := range yamlFiles {
		data, err := ioutil.ReadFile(file)
		fmt.Println(file)
		if err != nil {
			fmt.Println("Cannot read file:", err)
			os.Exit(1)
		}
		var r recipe.Recipe
		if err := yaml.Unmarshal(data, &r); err != nil {
			fmt.Println("Error parsing YAML file:", err)
			continue
		}
		r.ParseIngredientsList()
		for _, item := range r.IngredientsList {
			fmt.Println("amount: ", item.Amount)
			fmt.Println("unit: ", item.Unit)
			fmt.Println("label: ", item.Label)
		}
	}

}
