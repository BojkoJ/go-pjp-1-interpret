package main

import (
	"bufio"   // Balíček pro načítání vstupu po řádcích
	"fmt"     // Balíček pro formátovaný výstup
	"os"      // Balíček pro vstup a výstup
	"strconv" // Balíček pro konverzi řetězců na čísla
	"unicode" // Balíček pro práci s Unicode
)

// Typy tokenů
const (
	NUMBER   = "NUMBER"
	PLUS     = "PLUS"
	MINUS    = "MINUS"
	MULTIPLY = "MULTIPLY"
	DIVIDE   = "DIVIDE"
	LPAREN   = "LPAREN"  // (
	RPAREN   = "RPAREN"  // )
	END      = "END"     // Konec vstupu
	INVALID  = "INVALID" // Neplatný token
)

// Token je jedna jednotka výrazu (např. číslo, operátor nebo závorka).
type Token struct {
	Type  string // Typ tokenu (např. NUMBER, PLUS, MINUS, ...)
	Value string // Hodnota tokenu (např. "42", "+", "-", ...)
}

// Lexer je zodpovědný za převod vstupního řetězce na tokeny.
type Lexer struct {
	input       string // Daný příklad - jako řetězec
	position    int    // Aktuální pozice v řetězci
	currentChar rune   // Aktuální zpracovávaný znak
}

// NewLexer inicializuje nový lexer pro daný vstupní výraz.
func NewLexer(input string) *Lexer {
	// Vytvoříme proměnnou a dáme do ní nový lexer, jako pointer
	lexer := &Lexer{input: input, position: 0}

	if len(input) > 0 {
		lexer.currentChar = rune(input[0]) // Převedeme první znak na rune (na byte)
	}

	return lexer
}

// Metoda advance posune lexer na další znak v řetězci.
func (l *Lexer) advance() {
	l.position++

	if l.position < len(l.input) { // Pokud je pozice menší než délka vstupu
		l.currentChar = rune(l.input[l.position]) // Do currentChar uložíme znak na nové pozici
	} else {
		l.currentChar = 0 // Pokud je větší, nebo rovno délce vstupu, nastavíme na 0 - KONEC vstupu
	}
}

// Metoda skipWhitespace se postará o přeskočení bílých znaků
func (l *Lexer) skipWhitespace() {
	for unicode.IsSpace(l.currentChar) { // while cyklus, dokud je znak mezera - posuneme se na další znak
		l.advance()
	}
}

// metoda readNumber načte více číslicové číslo ze vstupu
func (l *Lexer) readNumber() string {
	number := ""

	for unicode.IsDigit(l.currentChar) { // Pokud je aktuální znak číslo
		number += string(l.currentChar) // Přidáme ho do řetězce number
		l.advance()                     // Posuneme se na další znak
	}
	return number
}

// metoda getNextToken vrátí další token ze vstupu
func (l *Lexer) getNextToken() Token {
	for l.currentChar != 0 { // dokud není konec vstupu

		if unicode.IsSpace(l.currentChar) { // Pokud je aktuální znak mezera - přeskočíme ji
			l.skipWhitespace()
			continue
		}

		if unicode.IsDigit(l.currentChar) { // Pokud je aktuální znak číslo
			return Token{NUMBER, l.readNumber()} // Vrátíme token s číslem
		}

		switch l.currentChar {

		case '+': // Pokud je aktuální znak plus
			l.advance()             // Posuneme se na další znak
			return Token{PLUS, "+"} // vrátíme token s plus
		case '-': // Mínus
			l.advance()
			return Token{MINUS, "-"}
		case '*': // Krát
			l.advance()
			return Token{MULTIPLY, "*"}
		case '/': // Děleno
			l.advance()
			return Token{DIVIDE, "/"}
		case '(': // Otevírací závorka
			l.advance()
			return Token{LPAREN, "("}
		case ')': // Zavírací závorka
			l.advance()
			return Token{RPAREN, ")"}
		default: // Default - neznámý znak
			return Token{INVALID, string(l.currentChar)}

		}
	}
	return Token{END, ""} // Cyklus neproběhl ani jednou, vrátíme token END
}

// Parser je zodpovědný za analýzu tokenů a výpočet výrazu.
type Parser struct {
	lexer *Lexer // Lexer pro získání tokenů
	token Token  // Aktuálně zpracovávaný token
}

// Vytvoříme nový parser s lexerem
func NewParser(lexer *Lexer) *Parser {
	return &Parser{lexer: lexer, token: lexer.getNextToken()} // Vrátíme pointer na nový parser, který má lexer a aktuální token
}

// metoda eat zkontroluje, zda aktuální token odpovídá očekávanému typu, a pak do aktuálního tokenu uloží další token
func (p *Parser) eat(tokenType string) bool {
	if p.token.Type == tokenType {
		p.token = p.lexer.getNextToken()
		return true
	}
	return false
}

// metoda parseFactor vrátí číslo, nebo vstoupí do závorek a zavolá parseExpression, která zpracuje výraz uvnitř závorek
// Závorky mají nejvyšší prioritu
func (p *Parser) parseFactor() (int, error) {

	if p.token.Type == NUMBER { // Pokud je aktuální token číslo
		value, _ := strconv.Atoi(p.token.Value) // Převedeme ho na číslo (je to typ string)
		p.eat(NUMBER)                           // Ověříme, že je to číslo a posuneme se na další token

		return value, nil // Vrátíme číslo
	} else if p.token.Type == LPAREN { // Pokud je aktuální token otevírací závorka
		p.eat(LPAREN)                      // Posuneme se na další token
		result, err := p.parseExpression() // Zpracujeme výraz uvnitř závorek

		if err != nil {
			return 0, err
		}

		if !p.eat(RPAREN) { // Pokud chybí zavírací závorka, vrátíme chybu
			return 0, fmt.Errorf("missing closing parenthesis")
		}

		return result, nil
	}
	return 0, fmt.Errorf("invalid syntax")
}

// řeší násobení a dělení (střední přednost)
func (p *Parser) parseTerm() (int, error) {
	result, err := p.parseFactor() // Získáme číslo, vyřešíme závorky
	if err != nil {
		return 0, err
	}

	for p.token.Type == MULTIPLY || p.token.Type == DIVIDE { // Dokud je token násobení nebo dělení
		operator := p.token // Uložíme si aktuální token do proměnné operator

		p.eat(operator.Type) // Posuneme se na další token

		// Zpracujeme případné závorky, pokud nejsou, jen získáme číslo
		// Např. 1 * (2 + 3)
		nextValue, err := p.parseFactor()
		if err != nil {
			return 0, err
		}

		// Zpracování násobení a dělení
		if operator.Type == MULTIPLY {
			result *= nextValue
		} else {
			if nextValue == 0 {
				return 0, fmt.Errorf("division by zero")
			}
			result /= nextValue
		}
	}
	return result, nil
}

// řeší sčítání a odčítání (nižší přednost)
func (p *Parser) parseExpression() (int, error) {
	result, err := p.parseTerm() // Násobení a dělení má vyšší přednost - zpracujeme ho první
	if err != nil {
		return 0, err
	}

	for p.token.Type == PLUS || p.token.Type == MINUS { // Dokud je token plus nebo mínus
		operator := p.token // Do proměnné operator uložíme aktuální token

		p.eat(operator.Type) // Posuneme se na další token

		// Zpracujeme násobení a dělení, které může být za plus nebo mínus
		// Např. 1 + 2 * 3
		nextValue, err := p.parseTerm()
		if err != nil {
			return 0, err
		}

		// Zpracování plus a mínus
		if operator.Type == PLUS {
			result += nextValue
		} else {
			result -= nextValue
		}
	}
	return result, nil
}

// funkce evaluateExpression vezme řetězec aritmetického výrazu a vypočítá výsledek.
func evaluateExpression(expression string) (int, error) {
	lexer := NewLexer(expression)
	parser := NewParser(lexer)
	return parser.parseExpression()
}

func main() {
	// Čteme input
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Enter number of expressions: ")
	scanner.Scan()
	numExpressions, _ := strconv.Atoi(scanner.Text())

	// Pro každý výraz zavoláme funkci evaluateExpression
	for i := 0; i < numExpressions; i++ {
		fmt.Print("Enter expression: ")
		scanner.Scan()
		expression := scanner.Text()
		result, err := evaluateExpression(expression)
		if err != nil {
			fmt.Println("ERROR")
		} else {
			fmt.Println(result)
		}
	}
}
