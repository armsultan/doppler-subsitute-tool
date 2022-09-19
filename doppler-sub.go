package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
)

func main() {

	// cli help if arguments are not provided
	if len(os.Args) < 3 {
		fmt.Println("Usage:")
		fmt.Println("DOPPLER_TOKEN=\"dp.st.dev_xxxxxxxxxxxxxxxxx\" go run doppler-sub.go <READ Directory> <WRITE Directory> <VARIABLE EXPRESSION>")
		fmt.Println("e.g. DOPPLER_TOKEN=\"dp.st.dev_xxxxxxxxxxxxxxxxx\" go run doppler-sub.go ./files ./export dollar-curly")
		fmt.Println("\n")
		fmt.Println("Available Variable Expression Formats:")
		fmt.Println(" * dollar	i.e. $MYVAR,\n * dollar-curly	i.e. ${MYVAR} (default),\n * handlebars	i.e. {{MYVAR}}, and\n * dollar-handlebars	i.e. ${{MYVAR}}")
		os.Exit(0)
	}

	// Get Doppler Token from env
	var token string = os.Getenv("DOPPLER_TOKEN")

	// set variable from third argument or default to "dollar-curly"
	var expType string
	// check if third argument is provided
	if len(os.Args) > 3 {
		expType = os.Args[3]
	} else {
		// set default
		expType = "dollar-curly"
	}

	// Map of supported variable expression formats and their regex
	regexFormat := map[string]string{
		"dollar":            `\$[A-Z]{1,}[A-Z0-9].*?`,
		"dollar-curly":      `\${([A-Z_]{1,}[A-Z0-9_].*?)}`,
		"handlebars":        `{{[A-Z_]{1,}[A-Z0-9_].*?}}`,
		"dollar-handlebars": `\${{[A-Z_]{1,}[A-Z0-9_].*?}}`,
	}

	// Set regex for variable expressions
	regExpression, _ := regexp.Compile(regexFormat[expType])
	// Print expType and regExpression
	fmt.Println("Variable expression format to target: ", expType, "i.e.", regExpression)

	// Save secrets from Doppler to a map
	var secrets map[string]interface{} = getSecrets(token)
	// Print out the secrets
	// fmt.Println("\n\n-------------SECRETS-----------------------\n\n")
	// fmt.Println(secrets)

	// Loop through of all files in the directory provided
	err := filepath.Walk(os.Args[1], func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// if file is not a directory
		if !info.IsDir() {
			fmt.Println("Reading", path)
			fmt.Println("\tSecrets Matched:")

			// Read file
			var file string = (readFile(path))

			// Match all capture groups into a array
			matches := regExpression.FindAllStringSubmatch(file, -1)

			// Create a int to count the number of secrets matched
			var count int = 0

			//	For each capture groups look up the secret in the map and
			//	replace the variable expression with the secret
			for _, match := range matches {
				// Reset count
				count = 0
				// if secret exists in map replace variable expression with
				// secret
				if val, ok := secrets[match[1]]; ok {
					// print out matches
					fmt.Println("\t\t", match[1], "\t✔")
					// print out matches and corresponding secret
					// fmt.Println("\t\t", match[1], "\t✔", "replacing with", val.(string))

					// Create empty string variable for regex expression
					var targetRegex string
					// set the Regex subsitution
					switch {
					case expType == "dollar":
						// Set targetRegex with dollar varible
						// expression
						targetRegex = `\$` + match[1] + `.*?`
					case expType == "handlebars":
						// Set targetRegex with handlebars varible
						// expression
						targetRegex = `{{` + match[1] + `.*?}}`
					case expType == "dollar-handlebars":
						// Set targetRegex with dollar handlebars
						// varible expression
						targetRegex = `\${{` + match[1] + `.*?}}`
					default:
						// Set targetRegex with DEFAULT dollar
						// curly varible expression
						targetRegex = `\${` + match[1] + `.*?}`
					}

					// set regex for variable expression
					x, _ := regexp.Compile(targetRegex)
					// Print out regex
					// fmt.Println("\t\tusing regex", targetRegex)

					// replace variable expression with secret
					file = x.ReplaceAllString(file, val.(string))
					// Print file
					// fmt.Println(file)

					// increment count
					count++

				}
				// If secret does not exist in map print out not found
				if _, ok := secrets[match[1]]; !ok {
					fmt.Println("\t\t", match[1], "\t✗")
				}
			}
			// if count is greater than 0 write file
			if count > 0 {

				// write file and create directory if it does not exist
				err := os.MkdirAll(os.Args[2], 0755)
				if err != nil {
					log.Fatal(err)
				}
				// Write file with replaced secrets and capture errors
				err = ioutil.WriteFile(os.Args[2]+"/"+info.Name(), []byte(file), 0644)
				// If error print error
				if err != nil {
					log.Fatal(err)
				}
				// Print count
				fmt.Println("\tTotal variables matched:", count)
				fmt.Println("\tSecrets written to", os.Args[2]+"/"+info.Name())
			}
			// If count is 0 print no secrets matched
			if count == 0 {
				fmt.Println("\t--NONE--")
				fmt.Println("\tNo secrets matched, and file not written")
			}
		}

		return nil
	})
	if err != nil {
		log.Println(err)
	}

}

// Return map of secrets
func getSecrets(token string) map[string]interface{} {

	client := &http.Client{}
	// Make and API call to Doppler with token in Basic Auth username (No Password)
	req, err := http.NewRequest("GET", "https://api.doppler.com/v3/configs/config/secrets/download?format=json", nil)
	req.SetBasicAuth(token, "")
	resp, err := client.Do(req)
	// Check for errors
	if err != nil {
		fmt.Print("Error retrieving secrets: ", err)
	}

	// Print response body
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Print("Error reading response body: ", err)
	}

	// Convert JSON String to Map
	var result map[string]interface{}
	json.Unmarshal([]byte(body), &result)

	// return map
	return result
}

// Read file
func readFile(filename string) string {
	// Read file
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
	}
	x := string(file)
	return x
}
