{{/* Go template for converting a pyprika Recipe object to Markdown */}}
# {{ .Name }}

{{ if .PrepTime }}
**Prep time:** {{ .PrepTime -}}
{{ end }}
{{ if .CookTime }}
**Cook time:** {{ .CookTime -}}
{{ end }}
{{ if .Servings }}
**servings:** {{ index .Servings 0 }}-{{ index .Servings 1 -}}
{{ end }}

{{ if .Notes }}
{{ .Notes }}
{{ end -}}

## Ingredients

{{ range $index, $Ingredient := .IngredientsList }}
{{- if $Ingredient.Amount -}}
- {{ $Ingredient.Amount }} {{ $Ingredient.Unit }} {{ $Ingredient.Label -}}
{{- else -}}
- {{ $Ingredient.Label -}}
{{ end }}
{{ end }}

## Directions

{{ range $index, $Direction := .Directions }}
{{ $index }}.  {{ $Direction -}}
{{ end }}

{{/* If source is provided with a URL, make it a link. */}}
{{ if and .Source .SourceURL }}
###### Source: [{{ .Source }}]({{ .SourceURL }})
{{ else if .Source }}
###### Source: {{ .Source }}
{{ end }}
