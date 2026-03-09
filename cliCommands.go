package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"pokedex/internal/pokecache"
)

type CliCommands struct {
	name        string
	description string
	callback    func(*config, map[string]Pokemon) error
}

type config struct {
	nextUrl *string
	prevUrl *string
	cache   *pokecache.Cache
	option  *string
}

type Pokemon struct {
	Id             int    `json:"id"`
	Name           string `json:"name"`
	BaseExperience int    `json:"base_experience"`
	Height         int    `json:"height"`
	IsDefault      bool   `json:"is_default"`
	Order          int    `json:"order"`
	Weight         int    `json:"weight"`

	Abilities []any          `json:"abilities"`
	Forms     []any          `json:"forms"`
	HeldItems []any          `json:"held_items"`
	Moves     []any          `json:"moves"`
	Species   map[string]any `json:"species"`
	Sprites   map[string]any `json:"sprites"`
	Cries     map[string]any `json:"cries"`
	Stats     []pokeStats    `json:"stats"`
	Types     []pokeTypes    `json:"types"`
}
type pokeStats struct {
	BaseStat int `json:"base_stat"`
	Effort   int `json:"effort"`
	Stat     struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	} `json:"stat"`
}
type pokeTypes struct {
	Slot int `json:"slot"`
	Type struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	} `json:"type"`
}

type ApiResponce struct {
	Count    int64          `json:"count"`
	Next     *string        `json:"next"`
	Previous *string        `json:"previous"`
	Results  []AreaResponce `json:"results"`
}

type AreaResponce struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type ExploreResponse struct {
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

func Exit(cnf *config, pokedex map[string]Pokemon) error {
	fmt.Println("\tClosing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}
func Help(cnf *config, pokedex map[string]Pokemon) error {
	fmt.Println("\tWelcome to the Pokedex! Usage: Help: Displays a Help message Exit: Exit the Pokedex")
	return nil
}

func Map(config *config, pokedex map[string]Pokemon) error {

	url := "https://pokeapi.co/api/v2/location-area"

	if config.nextUrl != nil {
		url = *config.nextUrl
	}

	cache := config.cache
	body, ok := cache.Get(url)

	if !ok {
		res, err := http.Get(url)
		if err != nil {
			fmt.Errorf("failed to get location, error: %v\n", err)
			return err
		}

		defer res.Body.Close()

		body, err = io.ReadAll(res.Body)
		if err != nil {
			return err
		}

		cache.Add(url, body)

	}

	var data ApiResponce
	json.Unmarshal(body, &data)

	config.nextUrl = data.Next
	config.prevUrl = data.Previous

	fmt.Print("Available areas:")
	for _, area := range data.Results {
		fmt.Printf("\n\t%s", area.Name)
	}
	fmt.Println("")

	return nil
}

func Mapb(conf *config, pokedex map[string]Pokemon) error {
	if conf.prevUrl == nil {
		fmt.Println("You're on the first page")
		return nil
	}

	cache := conf.cache
	body, ok := cache.Get(*conf.prevUrl)

	if !ok {
		res, err := http.Get(*conf.prevUrl)
		if err != nil {
			return err
		}

		defer res.Body.Close()

		body, err = io.ReadAll(res.Body)
		if err != nil {
			return err
		}

		cache.Add(*conf.prevUrl, body)
	}

	var data ApiResponce
	if err := json.Unmarshal(body, &data); err != nil {
		return err
	}

	conf.nextUrl = data.Next
	conf.prevUrl = data.Previous

	fmt.Print("Available areas:")
	for _, area := range data.Results {
		fmt.Printf("\n\t%s", area.Name)
	}
	fmt.Println("")

	return nil
}

func Explore(conf *config, pokedex map[string]Pokemon) error { //--/location-area/:id

	url := "https://pokeapi.co/api/v2/location-area/" + *conf.option

	cache := conf.cache
	body, ok := cache.Get(url)

	var data ExploreResponse

	if !ok {
		res, err := http.Get(url)
		if err != nil {
			return err
		}

		defer res.Body.Close()

		body, err = io.ReadAll(res.Body)
		if err != nil {
			return err
		}

		if err = json.Unmarshal(body, &data); err != nil {
			return err
		}

		for _, v := range data.PokemonEncounters {
			fmt.Printf("--%s\n", v.Pokemon.Name)
		}

		cache.Add(url, body)

		return nil
	}

	if err := json.Unmarshal(body, &data); err != nil {
		return err
	}

	for _, v := range data.PokemonEncounters {
		fmt.Printf("--%s\n", v.Pokemon.Name)
	}

	return nil
}

func Catch(conf *config, pokedex map[string]Pokemon) error {

	url := "https://pokeapi.co/api/v2/pokemon/" + *conf.option

	cache := conf.cache

	body, ok := cache.Get(url)

	if !ok {
		res, err := http.Get(url)
		if err != nil {
			return err
		}

		defer res.Body.Close()

		body, err = io.ReadAll(res.Body)

		var data Pokemon
		if err = json.Unmarshal(body, &data); err != nil {
			return err
		}

		roll := rand.Intn(data.BaseExperience)

		fmt.Printf("Throwing a Pokeball at %s...\n", data.Name)

		if roll < 100 {
			fmt.Printf("%s was caught!\n", data.Name)
			pokedex[data.Name] = data
		} else {
			fmt.Printf("%s escaped!\n", data.Name)
		}

		cache.Add(url, body)

		return nil
	}

	var data Pokemon
	if err := json.Unmarshal(body, &data); err != nil {
		return err
	}

	roll := rand.Intn(data.BaseExperience)

	fmt.Printf("Throwing a Pokeball at %s...\n", data.Name)

	if roll < 100 {
		fmt.Printf("%s was caught!\n", data.Name)
		pokedex[data.Name] = data
	} else {
		fmt.Printf("%s escaped!\n", data.Name)
	}

	return nil
}

func Inspect(conf *config, pokedex map[string]Pokemon) error {

	pokeToInspect := conf.option

	poke, ok := pokedex[*pokeToInspect]
	if !ok {
		fmt.Printf("You have not caught that pokemon: %s\n", *pokeToInspect)
		return nil
	}

	fmt.Printf("Name: %s\n", poke.Name)
	fmt.Printf("Height: %d\n", poke.Height)
	fmt.Printf("Weight: %d\n", poke.Weight)
	fmt.Printf("Stats:\n")

	stat := struct {
		hp             int
		attack         int
		defense        int
		specialAttack  int
		specialDefense int
		speed          int
	}{}

	for _, v := range poke.Stats {
		switch v.Stat.Name {
		case "hp":
			stat.hp = v.BaseStat
		case "attack":
			stat.attack = v.BaseStat
		case "defense":
			stat.defense = v.BaseStat
		case "special-attack":
			stat.specialAttack = v.BaseStat
		case "special-defense":
			stat.specialDefense = v.BaseStat
		case "speed":
			stat.speed = v.BaseStat
		default:
			continue
		}
	}

	fmt.Printf("  -hp: %d\n", stat.hp)
	fmt.Printf("  -attack: %d\n", stat.attack)
	fmt.Printf("  -defense: %d\n", stat.defense)
	fmt.Printf("  -special-attack: %d\n", stat.specialAttack)
	fmt.Printf("  -special-defense: %d\n", stat.specialDefense)
	fmt.Printf("  -speed: %d\n", stat.speed)

	fmt.Printf("Types:\n")

	for _, v := range poke.Types {
		fmt.Printf("  - %s\n", v.Type.Name)
	}

	return nil
}

func Pokedex(conf *config, pokedex map[string]Pokemon) error {

	if len(pokedex) < 1 {
		fmt.Println("You don't have pokemons")
		return errors.New("no pokemons")
	}
	for key, _ := range pokedex {
		fmt.Printf("  - %s\n", key)
	}

	return nil
}
