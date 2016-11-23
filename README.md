# paprika2markdown

[![Build Status](https://travis-ci.org/compybara/paprika2markdown.svg?branch=master)](https://travis-ci.org/compybara/paprika2markdown)

This utility is designed to parse recipes stored in the YAML formats used by [Paprika](https://paprikaapp.com/help/mac/#importrecipes) or [Pyprika](http://pyprika.readthedocs.io/en/latest/yaml-spec.html) and translate them into human-readable Markdown documents using a template.

You can use the templates included in the `templates/` directory, or provide your own.

## Installation

To install this utility first ensure that Go is [installed](https://golang.org/doc/install) and configued on your machine. Then run the following command:

  `go get -u github.com/compybara/paprika2markdown`

## Usage

paprika2markdown reads all `.yml` or `.yaml` files from a given directory of recipes and outputs markdown files to the specified output folder. The recipes folder and the template must exist and be readable or the program will not run.

If the output folder provided does not already exist the program will prompt the user at runtime for permission to create it. Any existing files in the output folder will be overwritten.

    Usage of ./paprika2markdown:
    -o string
      Folder to store the generated recipes. Existing files will be overwritten. (default "./recipes.wiki")
    -r string
      The folder containing the recipe YAML files. (default "./recipes")
    -t string
      The template to use for generating output. (default "./templates/template.md")

## Testing

Tests are written using ginkgo / gomega. Run tests with `ginkgo` or `go testgit td`
