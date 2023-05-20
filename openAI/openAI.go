package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sashabaranov/go-openai"
	"log"
	"os"
	"regexp"
	"strings"
	"topio/scrapper"
	"unicode/utf8"
)

type ProductWithGeneratedDescription struct {
	scrapper.Product
	RelatedPage          int    `json:"related_page"`
	GeneratedDescription string `json:"generated_description"`
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
						"3. Ta pojedyncza array powinna być jednynym elementem zwracanym. Zwróć tylko i wyłącznie array bez dodatkowego opisu." +
						"4. Dane wejściowe znajdują się poniżej:" +
						ai.prepareBaseDescriptionForAI(productNames),
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

func (ai *AI) GenerateQuestionsForProduct(numberOfQuestionsToGenerate int, productName string) ([]string, error) {

	messages := []openai.ChatCompletionMessage{
		{
			Role: openai.ChatMessageRoleUser,
			Content: "Podaj " + fmt.Sprint(numberOfQuestionsToGenerate) + " najważniejsze pytania, które kupujący ten produkt powinien zadać, aby świadomie wybrać najlepszy produkt dla siebie." +
				"Produkt to " + productName +
				" Podaj tylko pytania, bez odpowiedzi." +
				" Pytania wyredaguj w stylu, który zawiera elementy emocjonalne i perswazyjne, które mają na celu zainteresowanie czytelnika i przekonanie go do zapoznania się z zawartością. Wykorzystuje także pytania retoryczne, aby wzbudzić ciekawość i skłonić czytelnika do refleksji nad swoimi preferencjami." +
				" Zwróć odpowiedź jako strukturę danych array. Zwróć tylko i wyłącznie array.",
		},
	}

	response, err := ai.AskOpenAI(messages)

	if err != nil {
		fmt.Printf("[GenerateBodyBasedOnQuestions ERROR]: %v\n", err)
		return nil, err
	}

	var arrayString string
	// Define the array delimiter
	delimiter := "]"

	// Find the start and end positions of the array within the long string
	startIndex := strings.Index(response, "[")
	endIndex := strings.Index(response, delimiter)

	if startIndex >= 0 && endIndex >= 0 {
		// Extract the array substring
		arrayString = response[startIndex : endIndex+len(delimiter)]
		arrayString = strings.ReplaceAll(arrayString, "'", "\"")
		re := regexp.MustCompile(`("[^"]+")|[\s]+`)
		arrayString = re.ReplaceAllString(arrayString, "$1")
		arrayString = strings.ReplaceAll(arrayString, "\n", "")
	} else {
		fmt.Println("No array found in the string.")
	}

	var questionsArray []string
	err = json.Unmarshal([]byte(strings.TrimSpace(arrayString)), &questionsArray)
	if err != nil {
		fmt.Println("Error when converting response from AI to array:", err)
		return nil, err
	}

	return questionsArray, nil
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
			0,
			resp.Choices[0].Message.Content,
		})

	}

	return productsWithDescription
}

func (ai *AI) AskOpenAI(messages []openai.ChatCompletionMessage) (string, error) {
	resp, err := ai.Client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:       openai.GPT3Dot5Turbo,
			Messages:    messages,
			Temperature: 1,
		},
	)

	if err != nil {
		fmt.Printf("[AskOpenAI ERROR]: %v\n", err)
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}

func (ai *AI) GenerateBodyBasedOnQuestions(questions []string) (string, error) {

	messages := []openai.ChatCompletionMessage{
		{
			Role: openai.ChatMessageRoleUser,
			Content: "Produkt który mnie interesuje to odkurzacze pionowe. +" +
				"Proszę stwórz tekst, w którym znajdą się odpowiedzi na poniższe pytania. +" +
				"Zwracaj odpowiedź jako <h2>,<h3>,<p>,<ul>,<li>,<b> tagi nic poza tym" +
				"Pod tagiem <h2> lub <h3> zawsze musi znaleźć się <p> odpowiedź nie może zakończyć się tagiem <h2> lub <h3>" +
				"Zawsze kończ odpowiedź tagiem <p>" +
				"Nie mogą istnieć dwa tagi <h2> lub <h3> jeden pod drugim" +
				"Tekst wyredaguj w stylu, który zawiera elementy emocjonalne i perswazyjne, " +
				"które mają na celu zainteresowanie czytelnika i przekonanie go do zapoznania się z zawartością. " +
				"Wykorzystuje także pytania retoryczne, aby wzbudzić ciekawość i skłonić czytelnika do refleksji nad swoimi preferencjami." +
				"Nie podawaj konkretnych modeli produktu w odpowiedziach." +
				"Dodaj podsumowanie dopiero na samym końcu, na podstawie dostarczonych odpowiedzi." +
				"Jeżeli odpowiesz całkowicie na pytanie na końcu dodaj ---END---, ---END--- nie powinno być w żadnym znaczniku html",
		},
	}

	var htmlResponse bytes.Buffer

	for _, question := range questions {

		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: question,
		})

		responseToQuestion, err := ai.AskOpenAI(messages)

		if err != nil {
			fmt.Printf("[GenerateBodyBasedOnQuestions ERROR]: %v\n", err)
			return "", err
		}

		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleSystem,
			Content: responseToQuestion,
		})

		htmlResponse.WriteString(responseToQuestion)

		for !strings.Contains(htmlResponse.String(), "---END---") {

			messages = append(messages, openai.ChatCompletionMessage{
				Role: openai.ChatMessageRoleUser,
				Content: "Dokończ odpowiadać na zadane pytania od miejsca w którym skończyłeś" +
					"Odpowiedź wyredaguj w sposób ciągły do poprzedniego tekstu. Bez dodawania żadnego wstępu.",
			})

			response, err := ai.AskOpenAI(messages)

			if err != nil {
				fmt.Printf("[GenerateBodyBasedOnQuestions ERROR]: %v\n", err)
				return "", err
			}

			messages = append(messages, openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleSystem,
				Content: response,
			})

			htmlResponse.WriteString(response)

		}

	}

	cleanText := strings.ReplaceAll(htmlResponse.String(), "\n", "")
	cleanText = strings.ReplaceAll(cleanText, "---END---", "")

	return cleanText, nil
}
