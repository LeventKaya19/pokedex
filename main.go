package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	pokecache "github.com/LeventKaya19/pokedex/internal"
)

type CliCommand struct {
	name        string
	description string
	callback    func(*Config, string) error
}

type Config struct {
	Next     *string
	Previous *string
}

type AreaResults struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type EncounterData struct {
	EncounterMethodRates []struct {
		EncounterMethod struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"encounter_method"`
		VersionDetails []struct {
			Rate    int `json:"rate"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"encounter_method_rates"`
	GameIndex int `json:"game_index"`
	ID        int `json:"id"`
	Location  struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Name  string `json:"name"`
	Names []struct {
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
		Name string `json:"name"`
	} `json:"names"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			EncounterDetails []struct {
				Chance          int   `json:"chance"`
				ConditionValues []any `json:"condition_values"`
				MaxLevel        int   `json:"max_level"`
				Method          struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				} `json:"method"`
				MinLevel int `json:"min_level"`
			} `json:"encounter_details"`
			MaxChance int `json:"max_chance"`
			Version   struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
}

type Pokemon struct {
	Abilities []struct {
		Ability struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"ability"`
		IsHidden bool `json:"is_hidden"`
		Slot     int  `json:"slot"`
	} `json:"abilities"`
	BaseExperience int `json:"base_experience"`
	Cries          struct {
		Latest string `json:"latest"`
		Legacy string `json:"legacy"`
	} `json:"cries"`
	Forms []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"forms"`
	GameIndices []struct {
		GameIndex int `json:"game_index"`
		Version   struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"version"`
	} `json:"game_indices"`
	Height    int `json:"height"`
	HeldItems []struct {
		Item struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"item"`
		VersionDetails []struct {
			Rarity  int `json:"rarity"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"held_items"`
	ID                     int    `json:"id"`
	IsDefault              bool   `json:"is_default"`
	LocationAreaEncounters string `json:"location_area_encounters"`
	Moves                  []struct {
		Move struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"move"`
		VersionGroupDetails []struct {
			LevelLearnedAt  int `json:"level_learned_at"`
			MoveLearnMethod struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"move_learn_method"`
			Order        any `json:"order"`
			VersionGroup struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version_group"`
		} `json:"version_group_details"`
	} `json:"moves"`
	Name          string `json:"name"`
	Order         int    `json:"order"`
	PastAbilities []any  `json:"past_abilities"`
	PastTypes     []any  `json:"past_types"`
	Species       struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"species"`
	Sprites struct {
		BackDefault      string `json:"back_default"`
		BackFemale       string `json:"back_female"`
		BackShiny        string `json:"back_shiny"`
		BackShinyFemale  string `json:"back_shiny_female"`
		FrontDefault     string `json:"front_default"`
		FrontFemale      string `json:"front_female"`
		FrontShiny       string `json:"front_shiny"`
		FrontShinyFemale string `json:"front_shiny_female"`
		Other            struct {
			DreamWorld struct {
				FrontDefault string `json:"front_default"`
				FrontFemale  any    `json:"front_female"`
			} `json:"dream_world"`
			Home struct {
				FrontDefault     string `json:"front_default"`
				FrontFemale      string `json:"front_female"`
				FrontShiny       string `json:"front_shiny"`
				FrontShinyFemale string `json:"front_shiny_female"`
			} `json:"home"`
			OfficialArtwork struct {
				FrontDefault string `json:"front_default"`
				FrontShiny   string `json:"front_shiny"`
			} `json:"official-artwork"`
			Showdown struct {
				BackDefault      string `json:"back_default"`
				BackFemale       string `json:"back_female"`
				BackShiny        string `json:"back_shiny"`
				BackShinyFemale  any    `json:"back_shiny_female"`
				FrontDefault     string `json:"front_default"`
				FrontFemale      string `json:"front_female"`
				FrontShiny       string `json:"front_shiny"`
				FrontShinyFemale string `json:"front_shiny_female"`
			} `json:"showdown"`
		} `json:"other"`
		Versions struct {
			GenerationI struct {
				RedBlue struct {
					BackDefault      string `json:"back_default"`
					BackGray         string `json:"back_gray"`
					BackTransparent  string `json:"back_transparent"`
					FrontDefault     string `json:"front_default"`
					FrontGray        string `json:"front_gray"`
					FrontTransparent string `json:"front_transparent"`
				} `json:"red-blue"`
				Yellow struct {
					BackDefault      string `json:"back_default"`
					BackGray         string `json:"back_gray"`
					BackTransparent  string `json:"back_transparent"`
					FrontDefault     string `json:"front_default"`
					FrontGray        string `json:"front_gray"`
					FrontTransparent string `json:"front_transparent"`
				} `json:"yellow"`
			} `json:"generation-i"`
			GenerationIi struct {
				Crystal struct {
					BackDefault           string `json:"back_default"`
					BackShiny             string `json:"back_shiny"`
					BackShinyTransparent  string `json:"back_shiny_transparent"`
					BackTransparent       string `json:"back_transparent"`
					FrontDefault          string `json:"front_default"`
					FrontShiny            string `json:"front_shiny"`
					FrontShinyTransparent string `json:"front_shiny_transparent"`
					FrontTransparent      string `json:"front_transparent"`
				} `json:"crystal"`
				Gold struct {
					BackDefault      string `json:"back_default"`
					BackShiny        string `json:"back_shiny"`
					FrontDefault     string `json:"front_default"`
					FrontShiny       string `json:"front_shiny"`
					FrontTransparent string `json:"front_transparent"`
				} `json:"gold"`
				Silver struct {
					BackDefault      string `json:"back_default"`
					BackShiny        string `json:"back_shiny"`
					FrontDefault     string `json:"front_default"`
					FrontShiny       string `json:"front_shiny"`
					FrontTransparent string `json:"front_transparent"`
				} `json:"silver"`
			} `json:"generation-ii"`
			GenerationIii struct {
				Emerald struct {
					FrontDefault string `json:"front_default"`
					FrontShiny   string `json:"front_shiny"`
				} `json:"emerald"`
				FireredLeafgreen struct {
					BackDefault  string `json:"back_default"`
					BackShiny    string `json:"back_shiny"`
					FrontDefault string `json:"front_default"`
					FrontShiny   string `json:"front_shiny"`
				} `json:"firered-leafgreen"`
				RubySapphire struct {
					BackDefault  string `json:"back_default"`
					BackShiny    string `json:"back_shiny"`
					FrontDefault string `json:"front_default"`
					FrontShiny   string `json:"front_shiny"`
				} `json:"ruby-sapphire"`
			} `json:"generation-iii"`
			GenerationIv struct {
				DiamondPearl struct {
					BackDefault      string `json:"back_default"`
					BackFemale       string `json:"back_female"`
					BackShiny        string `json:"back_shiny"`
					BackShinyFemale  string `json:"back_shiny_female"`
					FrontDefault     string `json:"front_default"`
					FrontFemale      string `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale string `json:"front_shiny_female"`
				} `json:"diamond-pearl"`
				HeartgoldSoulsilver struct {
					BackDefault      string `json:"back_default"`
					BackFemale       string `json:"back_female"`
					BackShiny        string `json:"back_shiny"`
					BackShinyFemale  string `json:"back_shiny_female"`
					FrontDefault     string `json:"front_default"`
					FrontFemale      string `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale string `json:"front_shiny_female"`
				} `json:"heartgold-soulsilver"`
				Platinum struct {
					BackDefault      string `json:"back_default"`
					BackFemale       string `json:"back_female"`
					BackShiny        string `json:"back_shiny"`
					BackShinyFemale  string `json:"back_shiny_female"`
					FrontDefault     string `json:"front_default"`
					FrontFemale      string `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale string `json:"front_shiny_female"`
				} `json:"platinum"`
			} `json:"generation-iv"`
			GenerationV struct {
				BlackWhite struct {
					Animated struct {
						BackDefault      string `json:"back_default"`
						BackFemale       string `json:"back_female"`
						BackShiny        string `json:"back_shiny"`
						BackShinyFemale  string `json:"back_shiny_female"`
						FrontDefault     string `json:"front_default"`
						FrontFemale      string `json:"front_female"`
						FrontShiny       string `json:"front_shiny"`
						FrontShinyFemale string `json:"front_shiny_female"`
					} `json:"animated"`
					BackDefault      string `json:"back_default"`
					BackFemale       string `json:"back_female"`
					BackShiny        string `json:"back_shiny"`
					BackShinyFemale  string `json:"back_shiny_female"`
					FrontDefault     string `json:"front_default"`
					FrontFemale      string `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale string `json:"front_shiny_female"`
				} `json:"black-white"`
			} `json:"generation-v"`
			GenerationVi struct {
				OmegarubyAlphasapphire struct {
					FrontDefault     string `json:"front_default"`
					FrontFemale      string `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale string `json:"front_shiny_female"`
				} `json:"omegaruby-alphasapphire"`
				XY struct {
					FrontDefault     string `json:"front_default"`
					FrontFemale      string `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale string `json:"front_shiny_female"`
				} `json:"x-y"`
			} `json:"generation-vi"`
			GenerationVii struct {
				Icons struct {
					FrontDefault string `json:"front_default"`
					FrontFemale  any    `json:"front_female"`
				} `json:"icons"`
				UltraSunUltraMoon struct {
					FrontDefault     string `json:"front_default"`
					FrontFemale      string `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale string `json:"front_shiny_female"`
				} `json:"ultra-sun-ultra-moon"`
			} `json:"generation-vii"`
			GenerationViii struct {
				Icons struct {
					FrontDefault string `json:"front_default"`
					FrontFemale  string `json:"front_female"`
				} `json:"icons"`
			} `json:"generation-viii"`
		} `json:"versions"`
	} `json:"sprites"`
	Stats []struct {
		BaseStat int `json:"base_stat"`
		Effort   int `json:"effort"`
		Stat     struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Slot int `json:"slot"`
		Type struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"type"`
	} `json:"types"`
	Weight int `json:"weight"`
}

var commands = map[string]CliCommand{
	"exit": {
		name:        "exit",
		description: "Exit the Pokedex",
		callback:    commandExit,
	},
	"help": {
		name:        "help",
		description: "Help",
		callback:    commandHelp,
	},
	"map": {
		name:        "map",
		description: "Map",
		callback:    commandMap,
	},
	"mapb": {
		name:        "mapb",
		description: "Map Back",
		callback:    commandMapBack,
	},
	"explore": {
		name:        "explore",
		description: "Explore region",
		callback:    commandExplore,
	},
	"catch": {
		name:        "catch",
		description: "Catch a pokemon",
		callback:    commandCatch,
	},
	"inspect": {
		name:        "inspect",
		description: "Inspect a caught pokemon",
		callback:    commandInspect,
	},
	"pokedex": {
		name:        "pokedex",
		description: "Pokedex detail",
		callback:    commandPokedex,
	},
}

var catchedPokemons = map[string]Pokemon{}

const areaEndPoint = "https://pokeapi.co/api/v2/location-area"
const pokemonEndPoint = "https://pokeapi.co/api/v2/pokemon"

var cache *pokecache.Cache

func main() {
	s := bufio.NewScanner(os.Stdin)
	c := Config{}
	cache = pokecache.NewCache(time.Second * 15)
	initialUrl := areaEndPoint
	c.Next = &initialUrl
	for {
		fmt.Print("Pokedex > ")
		s.Scan()
		res := cleanInput(strings.ToLower(s.Text()))
		val, ok := commands[res[0]]
		if ok {
			param := ""
			if len(res) > 1 {
				param = res[1]
			}
			err := val.callback(&c, param)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

func commandExit(config *Config, parameter string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(config *Config, parameter string) error {
	fmt.Println("Welcome to the Pokedex!\nUsage:\n\nhelp: Displays a help message\nexit: Exit the Pokedex")
	return nil
}

func commandMap(config *Config, parameter string) error {
	if config.Next == nil {
		return errors.New("no next area to show")
	}
	url := *config.Next

	data, ok := cache.Get(url)
	var results AreaResults
	var err error
	if !ok {
		results, err = GetMap(url)
		if err != nil {
			return err
		}
	} else {
		results, err = GetCachedMap(data)
		if err != nil {
			return err
		}
	}

	for _, location := range results.Results {
		fmt.Println(location.Name)
	}
	nextUrl := results.Next
	prevUrl := results.Previous
	config.Previous = &prevUrl
	config.Next = &nextUrl

	return nil
}

func commandMapBack(config *Config, parameter string) error {
	if config.Previous == nil {
		return errors.New("no previous area to show")
	}
	url := *config.Previous

	data, ok := cache.Get(url)
	var results AreaResults
	var err error
	if !ok {
		results, err = GetMap(url)
		if err != nil {
			return err
		}
	} else {
		results, err = GetCachedMap(data)
		if err != nil {
			return err
		}
	}

	for _, location := range results.Results {
		fmt.Println(location.Name)
	}
	nextUrl := results.Next
	prevUrl := results.Previous
	config.Previous = &prevUrl
	config.Next = &nextUrl

	return nil
}

func commandExplore(config *Config, parameter string) error {
	exploreUrl := areaEndPoint + "/" + parameter

	data, ok := cache.Get(exploreUrl)
	var results EncounterData
	if ok {
		err := json.Unmarshal(data, &results)
		if err != nil {
			return err
		}
		listPokemonEncounters(results)
		return nil
	}
	res, err := http.Get(exploreUrl)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	cache.Add(exploreUrl, body)
	err = json.Unmarshal(body, &results)
	if err != nil {
		return err
	}

	listPokemonEncounters(results)

	return nil
}

func commandCatch(config *Config, paramter string) error {
	catchUrl := pokemonEndPoint + "/" + paramter
	fmt.Printf("Throwing a Pokeball at %s...\n", paramter)

	res, err := http.Get(catchUrl)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	var pokemon Pokemon
	err = json.Unmarshal(body, &pokemon)
	if err != nil {
		return err
	}

	dice := rand.Intn(pokemon.BaseExperience)

	if dice <= 50 {
		catchedPokemons[paramter] = pokemon
		fmt.Printf("%s was caught!\n", paramter)
	} else {
		fmt.Printf("%s escaped!\n", paramter)
	}

	return nil
}

func commandInspect(config *Config, paramter string) error {
	pokemon, ok := catchedPokemons[paramter]
	if !ok {
		fmt.Println("you have not caught that pokemon")
		return errors.New("you have not caught that pokemon")
	}
	fmt.Printf("Name: %s\n", pokemon.Name)
	fmt.Printf("Height: %d\n", pokemon.Height)
	fmt.Printf("Weight: %d\n", pokemon.Weight)
	fmt.Println("Stats:")
	fmt.Printf(" -hp: %d\n", pokemon.Stats[0].BaseStat)
	return nil
}

func commandPokedex(config *Config, parameter string) error {
	if len(catchedPokemons) == 0 {
		fmt.Println("No pokemon caught")
		return nil
	}
	fmt.Println("Your Pokedex:")
	for key := range catchedPokemons {
		fmt.Printf("- %s\n", key)
	}

	return nil
}

func cleanInput(text string) []string {
	res := []string{}
	slice := strings.Split(text, " ")

	for _, s := range slice {
		if s != "" {
			res = append(res, s)
		}
	}

	return res
}

func GetMap(url string) (AreaResults, error) {
	res, err := http.Get(url)
	if err != nil {
		return AreaResults{}, errors.New("fatal error")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return AreaResults{}, err
	}
	defer res.Body.Close()
	cache.Add(url, body)
	if res.StatusCode > 299 {
		return AreaResults{}, fmt.Errorf("bad response: %s", res.Status)
	}

	var result AreaResults
	err = json.Unmarshal(body, &result)
	if err != nil {
		return AreaResults{}, fmt.Errorf("could not unmarshal")
	}

	return result, nil
}

func GetCachedMap(data []byte) (AreaResults, error) {
	var result AreaResults
	err := json.Unmarshal(data, &result)
	if err != nil {
		return AreaResults{}, fmt.Errorf("could not unmarshal")
	}

	return result, nil
}

func listPokemonEncounters(data EncounterData) {
	for _, e := range data.PokemonEncounters {
		fmt.Println(e.Pokemon.Name)
	}
}
