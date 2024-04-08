package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-rod/rod"
)

func main() {

	// File
	data := "./wallet.txt"

	// Open File
	openFile, err := os.Open(data)
	if err != nil {
		fmt.Println("Error opening wallet file")
		return
	}
	defer openFile.Close()

	// Scanner for reading file
	scanner := bufio.NewScanner(openFile)

	// Get Request User for looping data
	fmt.Print("Enter the number of wallets to process: ")
	var input string
	fmt.Scanln(&input)

	// Convert input to integer
	maxCount, err := strconv.Atoi(input)
	if err != nil {
		fmt.Println("Invalid Input. Please try again")
		return
	}

	// Open output file in append mode, create if it doesn't exist
	outputFile, err := os.OpenFile("./result.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening result file:", err)
		return
	}

	// Launch Browser
	browser := rod.New().MustConnect()
	defer browser.MustClose()

	// Iteration looping
	loop := 0
	for scanner.Scan() {
		if loop >= maxCount {
			break
		}

		// Convert address to lowercase
		address := strings.ToLower(scanner.Text())

		// Fetching URL
		url := fmt.Sprintf("https://debank.com/profile/%s", address)

		// Open Page
		page := browser.MustPage(url).MustWindowFullscreen()

		// Wait for the page to load completely
		page.MustWaitLoad().MustWaitRequestIdle()

		//Wait 5 seconds
		time.Sleep(5 * time.Second)

		// element 1
		element, err := page.Element("#root div.HeaderInfo_totalAsset__dQnoy div")
		if err != nil {
			fmt.Println("Could not find the element, skipping:", err)
			continue
		}
		element.MustWaitVisible()
		text := element.MustText()

		split := strings.Split(text, "+")
		if len(split) > 0 {
			text = split[0]
		}

		// element 2
		elementData, err := page.Element("#root div.HeaderInfo_headerInfoTags__RWSTZ > div")
		if err != nil {
			fmt.Println("Could not find the element, skipping:", err)
			continue
		}

		elementData.MustWaitVisible()
		textData := elementData.MustText()

		splitData := strings.Split(textData, " ")
		if len(splitData) > 0 {
			textData = splitData[0]
		}

		// Format the result for output
		data := fmt.Sprintf("[%d] %s|%s|%s\n", loop+1, address, text, textData)
		fmt.Print(data)
		if _, err := outputFile.WriteString(data); err != nil {
			fmt.Println("Error writing to file:", err)
			continue
		}

		loop++
	}

	// Check error reading file
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading from file:", err)
	}
}
