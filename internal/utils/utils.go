package utils

import "fmt"

func IntToWord(num int) (string, error) {
	// Define a map of integer to string representation
	numberWords := map[int]string{
		0:  "zero",
		1:  "one",
		2:  "two",
		3:  "three",
		4:  "four",
		5:  "five",
		6:  "six",
		7:  "seven",
		8:  "eight",
		9:  "nine",
		10: "ten",
	}

	// Check if the number is in the map
	if word, exists := numberWords[num]; exists {
		return word, nil
	}
	return "", fmt.Errorf("number %d is not supported", num)
}

func GetColumnsToUse(number int) string {
	if number == 0 {
		return "one"
	}
	fNumber := float64(12 / number)
	if fNumber <= 1 {
		return "one"
	}
	numberStr, err := IntToWord(int(fNumber))
	if err != nil {
		return "one"
	}
	return numberStr
}
