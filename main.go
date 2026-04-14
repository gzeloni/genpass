package main

import (
	"crypto/rand"
	"fmt"
	"math"
	"math/big"
	"os"
	"strconv"
	"strings"
)

const (
	lower  = "abcdefghijklmnopqrstuvwxyz"
	upper  = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digits = "0123456789"
	symbol = "_-.:,!^"
)

var alnum = lower + upper + digits

type LengthError struct{ msg string }

func (e LengthError) Error() string { return e.msg }

type Args struct {
	length          int
	delimiter       string
	symbols         bool
	lengthOption    *int
	delimiterOption *string
	noSymbolsOption *bool
}

func randomInt(max int) int {
	n, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	if err != nil {
		panic(err)
	}
	return int(n.Int64())
}

func randomChar(charset string) byte {
	return charset[randomInt(len(charset))]
}

func validateLength(k int) (int, error) {
	if k < 36 {
		return 0, LengthError{"minimum password length is 36"}
	}
	root := int(math.Sqrt(float64(k)))
	if root*root != k {
		return 0, LengthError{"length must be a perfect square"}
	}
	return root, nil
}

func buildBodyCharset(delimiter string, allowSymbols bool) string {
	if !allowSymbols {
		return alnum
	}
	syms := symbol
	if len(delimiter) == 1 {
		syms = strings.ReplaceAll(syms, delimiter, "")
	}
	if syms == "" {
		return alnum
	}
	return alnum + syms
}

func fillRandom(chars []byte, charset string) {
	for i := range chars {
		chars[i] = randomChar(charset)
	}
}

func spiralOrder(n int) []int {
	order := make([]int, 0, n*n)
	top, bottom := 0, n-1
	left, right := 0, n-1

	for left <= right && top <= bottom {
		for c := left; c <= right; c++ {
			order = append(order, top*n+c)
		}
		top++

		for r := top; r <= bottom; r++ {
			order = append(order, r*n+right)
		}
		right--

		if top <= bottom {
			for c := right; c >= left; c-- {
				order = append(order, bottom*n+c)
			}
			bottom--
		}

		if left <= right {
			for r := bottom; r >= top; r-- {
				order = append(order, r*n+left)
			}
			left++
		}
	}

	return order
}

func primeIndices(limit int) []int {
	if limit <= 2 {
		return nil
	}
	sieve := make([]bool, limit)
	for i := 2; i < limit; i++ {
		sieve[i] = true
	}
	for p := 2; p*p < limit; p++ {
		if sieve[p] {
			for multiple := p * p; multiple < limit; multiple += p {
				sieve[multiple] = false
			}
		}
	}

	primes := make([]int, 0)
	for i := 2; i < limit; i++ {
		if sieve[i] {
			primes = append(primes, i)
		}
	}
	return primes
}

func chooseSymbol(delimiter string) byte {
	syms := symbol
	if len(delimiter) == 1 {
		syms = strings.ReplaceAll(syms, delimiter, "")
	}
	if syms == "" {
		return '_'
	}
	return randomChar(syms)
}

func ensureReservedDistinct(k int, positions ...int) {
	seen := make(map[int]struct{}, len(positions))
	for _, p := range positions {
		if p < 0 || p >= k {
			panic("reserved position out of bounds")
		}
		if _, ok := seen[p]; ok {
			panic("duplicate reserved position")
		}
		seen[p] = struct{}{}
	}
}

func buildPasswordChars(k int, n int, delimiter string, allowSymbols bool) []byte {
	grid := make([]byte, k)
	fillRandom(grid, buildBodyCharset(delimiter, allowSymbols))

	order := spiralOrder(n)
	primes := primeIndices(k)

	firstPrime := primes[0]
	middlePrime := primes[len(primes)/2]
	lastPrime := primes[len(primes)-1]

	start := order[0]
	end := order[k-1]
	digitPos := order[firstPrime]
	lowerPos := order[middlePrime]

	var symbolPos int
	useSymbol := allowSymbols
	if useSymbol {
		symbolPos = order[lastPrime]
		if symbolPos == start || symbolPos == end || symbolPos == digitPos || symbolPos == lowerPos {
			useSymbol = false
		}
	}

	if useSymbol {
		ensureReservedDistinct(k, start, end, digitPos, lowerPos, symbolPos)
	} else {
		ensureReservedDistinct(k, start, end, digitPos, lowerPos)
	}

	grid[start] = randomChar("Gg")
	grid[end] = randomChar("Zz")
	grid[digitPos] = randomChar(digits)
	grid[lowerPos] = randomChar(lower)

	upperPos := order[1]
	for upperPos == start || upperPos == end || upperPos == digitPos || upperPos == lowerPos || (useSymbol && upperPos == symbolPos) {
		upperPos++
		if upperPos >= k {
			panic("could not allocate uppercase position")
		}
	}
	grid[upperPos] = randomChar(upper)

	if useSymbol {
		grid[symbolPos] = chooseSymbol(delimiter)
	}

	seq := make([]byte, k)
	for i, idx := range order {
		seq[i] = grid[idx]
	}
	return seq
}

func splitBlocks(seq []byte, blockSize int) []string {
	blocks := make([]string, 0, len(seq)/blockSize)
	for i := 0; i < len(seq); i += blockSize {
		blocks = append(blocks, string(seq[i:i+blockSize]))
	}
	return blocks
}

func genPass(k int, delimiter string, allowSymbols bool) (string, error) {
	n, err := validateLength(k)
	if err != nil {
		return "", err
	}
	seq := buildPasswordChars(k, n, delimiter, allowSymbols)
	return strings.Join(splitBlocks(seq, n), delimiter), nil
}

func parseArgs(argv []string) (Args, error) {
	args := Args{
		length:    36,
		delimiter: "-",
		symbols:   true,
	}

	positionals := make([]string, 0, 2)

	for i := 0; i < len(argv); i++ {
		switch argv[i] {
		case "-k", "--length":
			if i+1 >= len(argv) {
				return args, fmt.Errorf("flag needs an argument: %s", argv[i])
			}
			v, err := strconv.Atoi(argv[i+1])
			if err != nil {
				return args, err
			}
			args.lengthOption = &v
			i++

		case "-d", "--delimiter":
			if i+1 >= len(argv) {
				return args, fmt.Errorf("flag needs an argument: %s", argv[i])
			}
			v := argv[i+1]
			args.delimiterOption = &v
			i++

		case "--no-symbols":
			v := true
			args.noSymbolsOption = &v

		default:
			if strings.HasPrefix(argv[i], "-") {
				return args, fmt.Errorf("unknown argument: %s", argv[i])
			}
			positionals = append(positionals, argv[i])
		}
	}

	if len(positionals) > 0 {
		v, err := strconv.Atoi(positionals[0])
		if err != nil {
			return args, err
		}
		args.length = v
	}
	if len(positionals) > 1 {
		args.delimiter = positionals[1]
	}
	if len(positionals) > 2 {
		return args, fmt.Errorf("too many positional arguments")
	}

	return args, nil
}

func resolveCLIArgs(args Args) (int, string, bool) {
	length := args.length
	if args.lengthOption != nil {
		length = *args.lengthOption
	}

	delimiter := args.delimiter
	if args.delimiterOption != nil {
		delimiter = *args.delimiterOption
	}

	symbols := args.symbols
	if args.noSymbolsOption != nil && *args.noSymbolsOption {
		symbols = false
	}

	return length, delimiter, symbols
}

func printUsageAndExit(message string) {
	fmt.Fprintf(os.Stderr, "usage: %s [length] [delimiter] [-k length] [-d delimiter] [--no-symbols]\n", os.Args[0])
	fmt.Fprintln(os.Stderr, message)
	os.Exit(2)
}

func main() {
	args, err := parseArgs(os.Args[1:])
	if err != nil {
		printUsageAndExit(err.Error())
	}

	length, delimiter, symbols := resolveCLIArgs(args)

	password, err := genPass(length, delimiter, symbols)
	if err != nil {
		printUsageAndExit(err.Error())
	}

	fmt.Println(password)
}
