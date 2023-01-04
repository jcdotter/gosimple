// Copyright 2022 escend llc. All rights reserved.
// Any unlicensed use of source code is prohibited
// Author: James Dotter
// Last modified: James Dotter - 20221206

// converts currency amounts from and to strings
// stores the currency as a struct containing
// infomation related to the currency denomitaiton and format
// leverages currency standards stored in currencies.json and defaults.json

package currency

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"strconv"
	"strings"
)

type Amount struct {
	Currency struct {
		Number        string
		Code          string
		Symbol        string
		Name          string
		Decimals      int
		FractionDelim string
		Divisor       int
	}
	Value float64
}

func (a *Amount) New(value float64, denomination string) error {
	denomination = strings.ToUpper(denomination)
	currencies := getCurrencyInfo()
	if currency, ok := currencies[denomination]; ok {
		a.Currency.Number = currency["num"].(string)
		a.Currency.Code = currency["code"].(string)
		a.Currency.Symbol = currency["symbol"].(string)
		a.Currency.Name = currency["currency"].(string)
		a.Currency.Decimals = int(currency["decimals"].(float64))
		a.Currency.FractionDelim = currency["fractiondelim"].(string)
		a.Currency.Divisor = int(currency["divisor"].(float64))
		places := math.Max(1, math.Pow(10, float64(a.Currency.Decimals)))
		a.Value = math.Round(value*places) / places
		return nil
	} else {
		return fmt.Errorf("utlis.types.currency.New: '%v' is not a recognized currency denomination", denomination)
	}
}

func (a *Amount) NewFromStr(s string) error {
	num, info, e := stringToParts(s)
	if !e {
		a, e = buildAmount(num, info, a)
	} else {
		return fmt.Errorf("utlis.types.currency.NewFromStr: '%v' is not a recognized currency", s)
	}
	return nil
}

// stringToParts evaluates a string containing a currency
// returns the value as a float,
// info about the currency as a string, and
// whether there was an error in attempting to parse the string
func stringToParts(s string) (float64, string, bool) {
	info := ""
	num := float64(0)
	nState := "pending"
	prior := ""
	sign := 1
	dec := 0
	e := false
	if s[0] == 45 || (s[0] == 40 && s[len(s)-1] == 41) {
		// first char is '-' or first is '(' and last is ')'
		sign = -1
	}
	for _, ch := range s {
		c := string(ch)
		if d, err := strconv.Atoi(c); err == nil {
			if e = nState == "complete"; !e {
				if dec == 0 {
					num = num*float64(10) + float64(d)
				} else {
					places := math.Max(1, math.Pow(10, float64(dec)))
					num += (float64(d) / places)
					dec++
				}
				if prior == "-" && nState == "pending" {
					sign = -1
				}
				nState = "active"
			}
		} else if nState == "active" && strings.Contains(c, `.`) {
			if e = dec != 0; !e {
				dec = 1
			}
		} else if !(nState == "active" && strings.Contains(c, `,`)) {
			info += c
			if nState == "active" {
				nState = "complete"
			}
			if strings.Contains("$£¥₩€", c) && prior == "-" {
				sign = -1
			}
		}
		if e {
			break
		}
		prior = c
	}
	return num * float64(sign), info, e
}

func buildAmount(value float64, info string, a *Amount) (*Amount, bool) {
	currencies := getCurrencyInfo()
	curMatches := []map[string]any{}
	symMatches := []map[string]any{}
	e := false
	for _, currency := range currencies {
		if strings.Contains(info, currency["code"].(string)) {
			curMatches = append(curMatches, currency)
		}
		if strings.Contains(info, currency["symbol"].(string)) {
			symMatches = append(symMatches, currency)
		}
	}
	if len(curMatches) != 1 && len(symMatches) > 0 {
		s := ""
		c := ""
		q := 0
		for _, m := range symMatches {
			if m["symbol"].(string) == s {
				q++
			}
			if len(m["symbol"].(string)) > len(s) {
				s = m["symbol"].(string)
				c = m["code"].(string)
				q = 1
			}
		}
		if q > 1 {
			defaults := getCurrencyDefaults()
			if code, ok := defaults[s]; ok {
				c = code
			} else {
				e = true
			}
		}
		if !e {
			a.New(value, c)
		}
	} else if len(curMatches) == 1 {
		a.New(value, curMatches[0]["code"].(string))
	} else {
		e = true
	}
	return a, e
}

func (a *Amount) Format() (string, error) {
	if a.Currency.Code == "" {
		return "", fmt.Errorf("utils.types.currency.Amount.Format: unable to format currency, no amount provided")
	}
	factor := math.Max(1, math.Pow(10, float64(a.Currency.Decimals)))
	numStr := strconv.FormatFloat(math.Round(a.Value/float64(a.Currency.Divisor)*factor)/factor, 'f', a.Currency.Decimals, 64)
	numParts := strings.Split(numStr, `.`)
	if len(numParts) == 1 {
		numParts = append(numParts, strings.Repeat("0", a.Currency.Decimals))
	} else if dec := a.Currency.Decimals - len(numParts[1]); dec > 0 {
		i, e := strconv.Atoi(numParts[1])
		if e != nil {
			panic("could not read decimals")
		}
		numParts[1] = fmt.Sprint(float64(i) * math.Max(1, math.Pow(10, float64(dec))))
	}
	if l := len(numParts[0]); l > 3 {
		iStr := ""
		for _, d := range numParts[0] {
			iStr += string(d)
			l--
			if l != 0 && l%3 == 0 {
				mDelim := ","
				if a.Currency.FractionDelim == mDelim {
					mDelim = "."
				}
				iStr += mDelim
			}
		}
		numParts[0] = iStr
	}
	d := ""
	if a.Currency.Decimals > 0 {
		d = a.Currency.FractionDelim
	}
	s := fmt.Sprintf(`%v%v%v%v %v`, a.Currency.Symbol, numParts[0], d, numParts[1], a.Currency.Code)
	return s, nil
}

func getCurrencyInfo() map[string]map[string]any {
	file, err := ioutil.ReadFile("./types/currency/currencies.json")
	currencies := map[string]map[string]any{}
	err = json.Unmarshal(file, &currencies)
	if err != nil {
		panic("utils.types.currency: could not access stored currencies")
	}
	return currencies
}

func getCurrencyDefaults() map[string]string {
	file, err := ioutil.ReadFile("./types/currency/defaults.json")
	defaults := map[string]string{}
	err = json.Unmarshal(file, &defaults)
	if err != nil {
		panic("utils.types.currency: could not access stored default currency symbols")
	}
	return defaults
}
