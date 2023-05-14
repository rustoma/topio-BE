package ai

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sashabaranov/go-openai"
	"log"
	"os"
	"topio/scrapper"
	"unicode/utf8"
)

type ProductWithGeneratedDescription struct {
	scrapper.Product
	GeneratedDescription string
}

type AI struct {
	Client *openai.Client
}

func InitAI() *openai.Client {
	_, ok := os.LookupEnv("AI_KEY")

	if !ok {
		log.Fatal("App requires AI_KEY variable")
	}

	return openai.NewClient(os.Getenv("AI_KEY"))
}

func (ai *AI) prepareBaseDescriptionForAI(description string) string {
	var preparedDescription string
	maxLength := 8000

	textLength := utf8.RuneCountInString(description)

	if textLength > maxLength {
		preparedDescription = description[:maxLength]
	} else {
		preparedDescription = description
	}

	return preparedDescription
}

func (ai *AI) SortProductsFromBestToWorse(products []ProductWithGeneratedDescription) ([]ProductWithGeneratedDescription, error) {

	productNamesArray, err := ai.generateArrayOfProductNamesFromBestToWorse(products)

	if err != nil {
		fmt.Printf("generateListOfProductsFromBestToWorse error: %v\n", err)
	}

	var sortedProducts []ProductWithGeneratedDescription

	for _, name := range productNamesArray {
		index := -1

		for i, product := range products {
			if product.Name == name {
				index = i
				break
			}
		}

		if index != -1 {
			sortedProducts = append(sortedProducts, products[index])
		} else {
			fmt.Printf("Element with name '%s' not found\n", name)
		}
	}

	if len(sortedProducts) != len(products) {
		return nil, errors.New("[SortProductsFromBestToWorse ERROR]: the number of items after sorting is not equal to the number of input products")
	}

	return sortedProducts, nil
}

func (ai *AI) generateArrayOfProductNamesFromBestToWorse(products []ProductWithGeneratedDescription) ([]string, error) {

	productNames := ""

	for _, product := range products {
		productNames += product.Name + "\n"
	}

	resp, err := ai.Client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role: openai.ChatMessageRoleUser,
					Content: "Na podstawie ogólnie dostępnych informacji posortuj listę od najlepszego do najgorszego według kryteriów opisanych poniżej w punktach:" +
						"1. Nie modyfikuj nazw. Identyczne nazwy z listy powinny się znaleźć w danych wyjściowych." +
						"2. Zwróć dane jako strukturę danych array" +
						"3. Ta pojedyncza array powinna być jednynym elementem zwracanym." +
						"4. Dane wejściowe znajdują się pomiędzy potrójnym backtickiem" +
						"```" + ai.prepareBaseDescriptionForAI(productNames) + "```",
				},
			},
			Temperature: 0,
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletion in generateListOfProductsFromBestToWorse error: %v\n", err)
		return nil, err
	}

	var productNamesArray []string
	err = json.Unmarshal([]byte(resp.Choices[0].Message.Content), &productNamesArray)
	if err != nil {
		fmt.Println("Error when converting response from AI to array:", err)
		return nil, err
	}

	return productNamesArray, nil
}

func (ai *AI) GenerateProductsWithDescription(products []scrapper.Product) []ProductWithGeneratedDescription {

	var productsWithDescription []ProductWithGeneratedDescription

	for _, product := range products {

		resp, err := ai.Client.CreateChatCompletion(
			context.Background(),
			openai.ChatCompletionRequest{
				Model: openai.GPT3Dot5Turbo,
				Messages: []openai.ChatCompletionMessage{
					{
						Role: openai.ChatMessageRoleUser,
						Content: "Na podstawie opisu:" +
							ai.prepareBaseDescriptionForAI(product.Description) + "Wygeneruj recenzje produktu",
					},
				},
				Temperature: 0.2,
			},
		)

		if err != nil {
			fmt.Printf("ChatCompletion error: %v\n", err)
			return nil
		}

		productsWithDescription = append(productsWithDescription, ProductWithGeneratedDescription{
			product,
			resp.Choices[0].Message.Content,
		})

	}

	return productsWithDescription

}
