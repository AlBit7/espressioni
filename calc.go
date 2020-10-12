package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

/*	LOGICA PROGRAMMA

	fin che non ci sono parentesi:

		risolvere ()

			risolvere ^
			risolvere /
			risolvere *
			risolvere -
			risolvere +

	a questo punto non ci sono più parentesi quindi:

	risolvere ()

*/

// ESP = espressione (forma abbreviata)
// divide l'espressione in numeri separati da operazioni
type ESP struct {
	numeri     []float64
	operazioni []uint8
}

// inizializzazione espressione
//
// se la stringa è 65+67/23-1^4 viene creato un ESP con:
//
// 	- numeri: [65, 67, 23, 1, 4]
// 	- operazioni: [43, 47, 45, 94]

func nuovaEspressione(stringa string) *ESP {

	p := new(ESP)
	numeroTemporaneo := ""
	valoreNumericoTemporaneo := 0.0
	parteDiNumero := false

	for _, carattere := range stringa {
		if unaDelleQuattroOperazioni(uint8(carattere)) && !parteDiNumero {
			p.operazioni = append(p.operazioni, uint8(carattere))
			valoreNumericoTemporaneo, _ = strconv.ParseFloat(numeroTemporaneo, 64)
			p.numeri = append(p.numeri, valoreNumericoTemporaneo)
			numeroTemporaneo = ""
			parteDiNumero = true
		} else {
			if carattere == '.' {
				numeroTemporaneo += string(carattere)
			} else {
				if carattere == '-' {
					numeroTemporaneo += "-"
				} else {
					numeroTemporaneo += fmt.Sprint(int(carattere) - 48)
				}
			}
			parteDiNumero = false
		}
	}

	valoreNumericoTemporaneo, _ = strconv.ParseFloat(numeroTemporaneo, 64)
	p.numeri = append(p.numeri, valoreNumericoTemporaneo)

	return p

}

// ----------------------- MAIN --------------------------

func main() {

	// leggi espressione sottoforma di stringa
	espressione := togliSpazi(input("\ninserisci un'espressione: "))
	fmt.Printf("\n")

	// dopo questo for tutte le parentesi verranno risolte
	for presenzaDiParentesi(espressione) {

		// dopo questo for una parentesi (in ordine di risoluzione)
		// sarà risolta e la stringa verrà modificata col risultato
		// dell'espressione all'interno della parentesi
		tmp := 0
		for i, elemento := range espressione {

			if elemento == '(' {
				tmp = i
			} else if elemento == ')' {
				espressione = espressione[:tmp] + fmt.Sprint(risolvi(nuovaEspressione(espressione[tmp+1:i]))) + espressione[i+1:]
				fmt.Println(" - ", espressione)
				break
			}

		}

	}

	fmt.Println("\nRisultato:", risolvi(nuovaEspressione(espressione)))

}

// ----------------------- FUNZIONI PER LA RISOLUZIONE --------------------------

func risolvi(espressione *ESP) float64 {

	// for viene eseguito fino a che c'è solo un
	// numero all'interno dell'espressione ovvero
	// quando non ci sono più operazioni da fare

	for len(espressione.numeri) != 1 {
		for risolviPotenza(espressione) {
		}
		for risolviDivisione(espressione) {
		}
		for risolviMoltiplicazione(espressione) {
		}
		for risolviSottrazione(espressione) {
		}
		for risolviAddizione(espressione) {
		}
	}
	return espressione.numeri[0]
}

func risolviPotenza(espressione *ESP) bool {

	for i, op := range espressione.operazioni {
		if op == 94 {

			// assegna il risultato al posto dell'esponente
			espressione.numeri[i+1] = float64(math.Pow(float64(espressione.numeri[i]), float64(espressione.numeri[i+1])))

			// pop espressione.numeri[i]
			// pop espressione.operazioni[i]
			elimina(espressione, i)

			// risolviPotenza risolve solo la prima potenza che trova
			// la funzione risolvi richiamerà risolvi potenza
			return true

		}
	}

	// poi, se non trova potenze, lo segnala alla funzione risolvi
	// successivamente la funzione risolvi passa a chiamare risolviDivisione e così via...
	return false

}

func risolviDivisione(espressione *ESP) bool {

	for i, op := range espressione.operazioni {
		if op == 47 {

			espressione.numeri[i+1] = espressione.numeri[i] / espressione.numeri[i+1]
			elimina(espressione, i)
			return true

		}
	}

	return false

}

func risolviMoltiplicazione(espressione *ESP) bool {

	for i, op := range espressione.operazioni {
		if op == 42 {

			espressione.numeri[i+1] = espressione.numeri[i] * espressione.numeri[i+1]
			elimina(espressione, i)
			return true

		}
	}

	return false

}

func risolviSottrazione(espressione *ESP) bool {

	for i, op := range espressione.operazioni {
		if op == 45 {

			espressione.numeri[i+1] = espressione.numeri[i] - espressione.numeri[i+1]
			elimina(espressione, i)
			return true

		}
	}

	return false

}

func risolviAddizione(espressione *ESP) bool {

	for i, op := range espressione.operazioni {
		if op == 43 {

			espressione.numeri[i+1] = espressione.numeri[i] + espressione.numeri[i+1]
			elimina(espressione, i)
			return true

		}
	}

	return false

}

// ----------------------- FUNZIONI AUSILIARI --------------------------

func unaDelleQuattroOperazioni(carattere uint8) bool {
	if carattere == '+' || carattere == '-' || carattere == '*' || carattere == '/' || carattere == '^' {
		return true
	}
	return false
}

func elimina(espressione *ESP, i int) {

	copy(espressione.numeri[i:], espressione.numeri[i+1:])
	espressione.numeri = espressione.numeri[:len(espressione.numeri)-1]

	copy(espressione.operazioni[i:], espressione.operazioni[i+1:])
	espressione.operazioni = espressione.operazioni[:len(espressione.operazioni)-1]

}

func presenzaDiParentesi(espressione string) bool {
	for _, elemento := range espressione {
		if elemento == '(' {
			return true
		}
	}
	return false
}

func togliSpazi(stringa string) string {
	tmp := strings.NewReplacer(" ", "")
	return tmp.Replace(stringa)
}

func input(messaggio string) string {

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("%s", messaggio)
	scanner.Scan()

	return scanner.Text()

}
