package main

import (
	"fmt"
	"ieee-754-2008/ieee"
	"os"
)

func main() {
	var option int
	var input1, input2 float32
	var err error

	for {
		fmt.Printf("\n1- Soma\n2- Subtração\n3- Multiplicação\n4- Divisão\n5- Sair\n")
		fmt.Print("Digite o número da operação: ")
		_, err = fmt.Scanln(&option)
		if err != nil {
			os.Exit(1)
		}
		if option == 5 {
			os.Exit(0)
		}
		if option < 1 || option > 5 {
			fmt.Println("\033[31mErro:\033[0m Opção Inválida.")
			continue
		}

		fmt.Print("Digite o primeiro número: ")
		_, err = fmt.Scanln(&input1)
		if err != nil {
			fmt.Println("\033[31mErro:\033[0m Número Inválido.")
			continue
		}
		fmt.Print("Digite o segundo número: ")
		_, err = fmt.Scanln(&input2)
		if err != nil {
			fmt.Println("\033[31mErro:\033[0m Número Inválido.")
			continue
		}

		n1, errN1 := ieee.NewIEEENumber(input1)
		n2, errN2 := ieee.NewIEEENumber(input2)
		if errN1 != nil || errN2 != nil {
			fmt.Println("\033[31mErro:\033[0m Não foi possível criar o número IEEE")
			if errN1 != nil {
				fmt.Printf("Primeiro Número Error: %v", errN1)
			}
			if errN2 != nil {
				fmt.Printf("Primeiro Número Error: %v", errN2)
			}
			continue
		}
		var result *ieee.IEEENumber

		switch option {
		case 1:
			result, err = ieee.Sum(n1, n2)
		case 2:
			result, err = ieee.Sub(n1, n2)
		case 3:
			result, err = ieee.Mult(n1, n2)
		case 4:
			result, err = ieee.Div(n1, n2)
		}

		if err != nil {
			fmt.Printf("\033[33mAVISO:\033[0m %v\n", err)
			continue
		}

		fmt.Println("\nResultado:")
		fmt.Printf("\033[33mResultado Float32:\033[0m %.15f\n", result.ToFloat32())
		fmt.Printf("\033[33mResultado Int (Truncamento):\033[0m %d\n", result.ToInt(ieee.Truncate))
		fmt.Printf("\033[33mResultado Int (Arredondamento):\033[0m %d\n", result.ToInt(ieee.Nearest))
		fmt.Println("==================================================")
		result.Debug()
		fmt.Println(result)
		fmt.Println("==================================================")
	}
}
