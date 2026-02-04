package game

import "fmt"

// CreateAllCards returns all 106 cards for the game
func CreateAllCards() []Card {
	cards := make([]Card, 0, 106)

	// Property cards (28)
	cards = append(cards, createPropertyCards()...)

	// Wildcard cards (11)
	cards = append(cards, createWildcardCards()...)

	// Action cards (34)
	cards = append(cards, createActionCards()...)

	// Rent cards (13)
	cards = append(cards, createRentCards()...)

	// Money cards (20)
	cards = append(cards, createMoneyCards()...)

	return cards
}

func createPropertyCards() []Card {
	cards := make([]Card, 0, 28)

	// Brown (J/Z) - 2 cards, rent $1/$2, value $1
	cards = append(cards, &PropertyCard{
		BaseCard: BaseCard{ID: "prop_j", Type: CardTypeProperty, Name: "J", Value: 1},
		Color:    "brown", ColorHex: "#996633", SetSize: 2, Position: 1,
		Rent: map[int]int{1: 1, 2: 2},
	})
	cards = append(cards, &PropertyCard{
		BaseCard: BaseCard{ID: "prop_z", Type: CardTypeProperty, Name: "Z", Value: 1},
		Color:    "brown", ColorHex: "#996633", SetSize: 2, Position: 2,
		Rent: map[int]int{1: 1, 2: 2},
	})

	// Light Blue (A/C/E) - 3 cards, rent $1/$2/$3, value $1
	cards = append(cards, &PropertyCard{
		BaseCard: BaseCard{ID: "prop_a", Type: CardTypeProperty, Name: "A", Value: 1},
		Color:    "blue", ColorHex: "#0039A6", SetSize: 3, Position: 1,
		Rent: map[int]int{1: 1, 2: 2, 3: 3},
	})
	cards = append(cards, &PropertyCard{
		BaseCard: BaseCard{ID: "prop_c", Type: CardTypeProperty, Name: "C", Value: 1},
		Color:    "blue", ColorHex: "#0039A6", SetSize: 3, Position: 2,
		Rent: map[int]int{1: 1, 2: 2, 3: 3},
	})
	cards = append(cards, &PropertyCard{
		BaseCard: BaseCard{ID: "prop_e", Type: CardTypeProperty, Name: "E", Value: 1},
		Color:    "blue", ColorHex: "#0039A6", SetSize: 3, Position: 3,
		Rent: map[int]int{1: 1, 2: 2, 3: 3},
	})

	// Pink/Gray (Shuttles) - 3 cards, rent $1/$2/$4, value $2
	cards = append(cards, &PropertyCard{
		BaseCard: BaseCard{ID: "prop_42nd_shuttle", Type: CardTypeProperty, Name: "42nd St Shuttle", Value: 2},
		Color:    "pink", ColorHex: "#808183", SetSize: 3, Position: 1,
		Rent: map[int]int{1: 1, 2: 2, 3: 4},
	})
	cards = append(cards, &PropertyCard{
		BaseCard: BaseCard{ID: "prop_franklin_shuttle", Type: CardTypeProperty, Name: "Franklin Ave Shuttle", Value: 2},
		Color:    "pink", ColorHex: "#808183", SetSize: 3, Position: 2,
		Rent: map[int]int{1: 1, 2: 2, 3: 4},
	})
	cards = append(cards, &PropertyCard{
		BaseCard: BaseCard{ID: "prop_rockaway_shuttle", Type: CardTypeProperty, Name: "Rockaway Park Shuttle", Value: 2},
		Color:    "pink", ColorHex: "#808183", SetSize: 3, Position: 3,
		Rent: map[int]int{1: 1, 2: 2, 3: 4},
	})

	// Orange (B/D/F) - 3 cards, rent $1/$3/$5, value $2
	cards = append(cards, &PropertyCard{
		BaseCard: BaseCard{ID: "prop_b", Type: CardTypeProperty, Name: "B", Value: 2},
		Color:    "orange", ColorHex: "#FF6319", SetSize: 3, Position: 1,
		Rent: map[int]int{1: 1, 2: 3, 3: 5},
	})
	cards = append(cards, &PropertyCard{
		BaseCard: BaseCard{ID: "prop_d", Type: CardTypeProperty, Name: "D", Value: 2},
		Color:    "orange", ColorHex: "#FF6319", SetSize: 3, Position: 2,
		Rent: map[int]int{1: 1, 2: 3, 3: 5},
	})
	cards = append(cards, &PropertyCard{
		BaseCard: BaseCard{ID: "prop_f", Type: CardTypeProperty, Name: "F", Value: 2},
		Color:    "orange", ColorHex: "#FF6319", SetSize: 3, Position: 3,
		Rent: map[int]int{1: 1, 2: 3, 3: 5},
	})

	// Red (1/2/3) - 3 cards, rent $2/$3/$6, value $3
	cards = append(cards, &PropertyCard{
		BaseCard: BaseCard{ID: "prop_1", Type: CardTypeProperty, Name: "1", Value: 3},
		Color:    "red", ColorHex: "#EE352E", SetSize: 3, Position: 1,
		Rent: map[int]int{1: 2, 2: 3, 3: 6},
	})
	cards = append(cards, &PropertyCard{
		BaseCard: BaseCard{ID: "prop_2", Type: CardTypeProperty, Name: "2", Value: 3},
		Color:    "red", ColorHex: "#EE352E", SetSize: 3, Position: 2,
		Rent: map[int]int{1: 2, 2: 3, 3: 6},
	})
	cards = append(cards, &PropertyCard{
		BaseCard: BaseCard{ID: "prop_3", Type: CardTypeProperty, Name: "3", Value: 3},
		Color:    "red", ColorHex: "#EE352E", SetSize: 3, Position: 3,
		Rent: map[int]int{1: 2, 2: 3, 3: 6},
	})

	// Yellow (N/Q/R) - 3 cards, rent $2/$4/$6, value $3
	cards = append(cards, &PropertyCard{
		BaseCard: BaseCard{ID: "prop_n", Type: CardTypeProperty, Name: "N", Value: 3},
		Color:    "yellow", ColorHex: "#FCCC0A", SetSize: 3, Position: 1,
		Rent: map[int]int{1: 2, 2: 4, 3: 6},
	})
	cards = append(cards, &PropertyCard{
		BaseCard: BaseCard{ID: "prop_q", Type: CardTypeProperty, Name: "Q", Value: 3},
		Color:    "yellow", ColorHex: "#FCCC0A", SetSize: 3, Position: 2,
		Rent: map[int]int{1: 2, 2: 4, 3: 6},
	})
	cards = append(cards, &PropertyCard{
		BaseCard: BaseCard{ID: "prop_r", Type: CardTypeProperty, Name: "R", Value: 3},
		Color:    "yellow", ColorHex: "#FCCC0A", SetSize: 3, Position: 3,
		Rent: map[int]int{1: 2, 2: 4, 3: 6},
	})

	// Green (Hub Stations) - 3 cards, rent $2/$4/$7, value $4
	cards = append(cards, &PropertyCard{
		BaseCard: BaseCard{ID: "prop_penn", Type: CardTypeProperty, Name: "Penn Station", Value: 4},
		Color:    "green", ColorHex: "#00933C", SetSize: 3, Position: 1,
		Rent: map[int]int{1: 2, 2: 4, 3: 7},
	})
	cards = append(cards, &PropertyCard{
		BaseCard: BaseCard{ID: "prop_grand_central", Type: CardTypeProperty, Name: "Grand Central", Value: 4},
		Color:    "green", ColorHex: "#00933C", SetSize: 3, Position: 2,
		Rent: map[int]int{1: 2, 2: 4, 3: 7},
	})
	cards = append(cards, &PropertyCard{
		BaseCard: BaseCard{ID: "prop_atlantic", Type: CardTypeProperty, Name: "Atlantic Terminal", Value: 4},
		Color:    "green", ColorHex: "#00933C", SetSize: 3, Position: 3,
		Rent: map[int]int{1: 2, 2: 4, 3: 7},
	})

	// Dark Blue (Stadiums) - 2 cards, rent $3/$8, value $4
	cards = append(cards, &PropertyCard{
		BaseCard: BaseCard{ID: "prop_citi_field", Type: CardTypeProperty, Name: "Citi Field", Value: 4},
		Color:    "darkblue", ColorHex: "#2A344D", SetSize: 2, Position: 1,
		Rent: map[int]int{1: 3, 2: 8},
	})
	cards = append(cards, &PropertyCard{
		BaseCard: BaseCard{ID: "prop_yankee_stadium", Type: CardTypeProperty, Name: "Yankee Stadium", Value: 4},
		Color:    "darkblue", ColorHex: "#2A344D", SetSize: 2, Position: 2,
		Rent: map[int]int{1: 3, 2: 8},
	})

	// Railroad (Black) - 4 cards, rent $1/$2/$3/$4, value $2
	cards = append(cards, &PropertyCard{
		BaseCard: BaseCard{ID: "prop_lirr", Type: CardTypeProperty, Name: "LIRR", Value: 2},
		Color:    "railroad", ColorHex: "#000000", SetSize: 4, Position: 1,
		Rent: map[int]int{1: 1, 2: 2, 3: 3, 4: 4},
	})
	cards = append(cards, &PropertyCard{
		BaseCard: BaseCard{ID: "prop_metro_north", Type: CardTypeProperty, Name: "Metro-North", Value: 2},
		Color:    "railroad", ColorHex: "#000000", SetSize: 4, Position: 2,
		Rent: map[int]int{1: 1, 2: 2, 3: 3, 4: 4},
	})
	cards = append(cards, &PropertyCard{
		BaseCard: BaseCard{ID: "prop_nj_transit", Type: CardTypeProperty, Name: "NJ Transit", Value: 2},
		Color:    "railroad", ColorHex: "#000000", SetSize: 4, Position: 3,
		Rent: map[int]int{1: 1, 2: 2, 3: 3, 4: 4},
	})
	cards = append(cards, &PropertyCard{
		BaseCard: BaseCard{ID: "prop_path", Type: CardTypeProperty, Name: "PATH", Value: 2},
		Color:    "railroad", ColorHex: "#000000", SetSize: 4, Position: 4,
		Rent: map[int]int{1: 1, 2: 2, 3: 3, 4: 4},
	})

	// Utilities (G/L) - 2 cards, rent $1/$2, value $2
	cards = append(cards, &PropertyCard{
		BaseCard: BaseCard{ID: "prop_g", Type: CardTypeProperty, Name: "G", Value: 2},
		Color:    "utility", ColorHex: "#6CBE45", SetSize: 2, Position: 1,
		Rent: map[int]int{1: 1, 2: 2},
	})
	cards = append(cards, &PropertyCard{
		BaseCard: BaseCard{ID: "prop_l", Type: CardTypeProperty, Name: "L", Value: 2},
		Color:    "utility", ColorHex: "#6CBE45", SetSize: 2, Position: 2,
		Rent: map[int]int{1: 1, 2: 2},
	})

	return cards
}

func createWildcardCards() []Card {
	cards := make([]Card, 0, 11)

	// Broadway Junction (Blue/Brown) - 1 card, value $1
	cards = append(cards, &WildcardCard{
		BaseCard:     BaseCard{ID: "wild_broadway", Type: CardTypeWildcard, Name: "Broadway Junction", Value: 1},
		Colors:       []string{"blue", "brown"},
		CurrentColor: "blue",
		Description:  "Can be used as Light Blue (A/C/E) or Brown (J/Z)",
	})

	// Jamaica Station (Blue/Railroad) - 1 card, value $4
	cards = append(cards, &WildcardCard{
		BaseCard:     BaseCard{ID: "wild_jamaica", Type: CardTypeWildcard, Name: "Jamaica Station", Value: 4},
		Colors:       []string{"blue", "railroad"},
		CurrentColor: "blue",
		Description:  "Can be used as Light Blue (A/C/E) or Railroad",
	})

	// Service Advisory (Pink/Orange) - 2 cards, value $2
	for i := 1; i <= 2; i++ {
		cards = append(cards, &WildcardCard{
			BaseCard:     BaseCard{ID: fmt.Sprintf("wild_service_advisory_%d", i), Type: CardTypeWildcard, Name: "Service Advisory", Value: 2},
			Colors:       []string{"pink", "orange"},
			CurrentColor: "pink",
			Description:  "Can be used as Pink (Shuttles) or Orange (B/D/F)",
		})
	}

	// Times Square (Red/Yellow) - 2 cards, value $3
	for i := 1; i <= 2; i++ {
		cards = append(cards, &WildcardCard{
			BaseCard:     BaseCard{ID: fmt.Sprintf("wild_times_square_%d", i), Type: CardTypeWildcard, Name: "Times Square", Value: 3},
			Colors:       []string{"red", "yellow"},
			CurrentColor: "red",
			Description:  "Can be used as Red (1/2/3) or Yellow (N/Q/R)",
		})
	}

	// Big Game (Dark Blue/Green) - 1 card, value $4
	cards = append(cards, &WildcardCard{
		BaseCard:     BaseCard{ID: "wild_big_game", Type: CardTypeWildcard, Name: "Big Game", Value: 4},
		Colors:       []string{"darkblue", "green"},
		CurrentColor: "darkblue",
		Description:  "Can be used as Dark Blue (Stadiums) or Green (Hubs)",
	})

	// Grand Central (Green/Railroad) - 1 card, value $4
	cards = append(cards, &WildcardCard{
		BaseCard:     BaseCard{ID: "wild_grand_central", Type: CardTypeWildcard, Name: "Grand Central", Value: 4},
		Colors:       []string{"green", "railroad"},
		CurrentColor: "green",
		Description:  "Can be used as Green (Hubs) or Railroad",
	})

	// Weekend Service (Utility/Railroad) - 1 card, value $2
	cards = append(cards, &WildcardCard{
		BaseCard:     BaseCard{ID: "wild_weekend_service", Type: CardTypeWildcard, Name: "Weekend Service", Value: 2},
		Colors:       []string{"utility", "railroad"},
		CurrentColor: "utility",
		Description:  "Can be used as Utilities (G/L) or Railroad",
	})

	// Fulton Center (Multi-color) - 2 cards, value $0
	for i := 1; i <= 2; i++ {
		cards = append(cards, &WildcardCard{
			BaseCard:     BaseCard{ID: fmt.Sprintf("wild_fulton_%d", i), Type: CardTypeWildcard, Name: "Fulton Center", Value: 0},
			Colors:       []string{"brown", "blue", "pink", "orange", "red", "yellow", "green", "darkblue", "railroad", "utility"},
			CurrentColor: "blue",
			Description:  "Can be used as any color (no cash value)",
		})
	}

	return cards
}

func createActionCards() []Card {
	cards := make([]Card, 0, 34)

	// Swipe In (Pass Go) - 10 cards, value $1
	for i := 1; i <= 10; i++ {
		cards = append(cards, &ActionCard{
			BaseCard: BaseCard{ID: fmt.Sprintf("action_swipe_in_%d", i), Type: CardTypeAction, Name: "Swipe In", Value: 1},
			Effect:   "Draw 2 extra cards",
			MTATheme: "Pass Go",
		})
	}

	// Fare Evasion (Just Say No) - 3 cards, value $4
	for i := 1; i <= 3; i++ {
		cards = append(cards, &ActionCard{
			BaseCard: BaseCard{ID: fmt.Sprintf("action_fare_evasion_%d", i), Type: CardTypeAction, Name: "Fare Evasion", Value: 4},
			Effect:   "Cancel any action card played against you",
			MTATheme: "Just Say No",
		})
	}

	// Power Broker (Sly Deal) - 3 cards, value $3
	for i := 1; i <= 3; i++ {
		cards = append(cards, &ActionCard{
			BaseCard: BaseCard{ID: fmt.Sprintf("action_power_broker_%d", i), Type: CardTypeAction, Name: "Power Broker", Value: 3},
			Effect:   "Steal a property from any player (not from a complete set)",
			MTATheme: "Sly Deal",
		})
	}

	// Service Change (Forced Deal) - 3 cards, value $3
	for i := 1; i <= 3; i++ {
		cards = append(cards, &ActionCard{
			BaseCard: BaseCard{ID: fmt.Sprintf("action_service_change_%d", i), Type: CardTypeAction, Name: "Service Change", Value: 3},
			Effect:   "Swap one of your properties with another player's (not from complete sets)",
			MTATheme: "Forced Deal",
		})
	}

	// Line Closure (Deal Breaker) - 2 cards, value $5
	for i := 1; i <= 2; i++ {
		cards = append(cards, &ActionCard{
			BaseCard: BaseCard{ID: fmt.Sprintf("action_line_closure_%d", i), Type: CardTypeAction, Name: "Line Closure", Value: 5},
			Effect:   "Steal a complete property set from any player (includes improvements)",
			MTATheme: "Deal Breaker",
		})
	}

	// Missed Your Train (Debt Collector) - 3 cards, value $3
	for i := 1; i <= 3; i++ {
		cards = append(cards, &ActionCard{
			BaseCard: BaseCard{ID: fmt.Sprintf("action_missed_train_%d", i), Type: CardTypeAction, Name: "Missed Your Train", Value: 3},
			Effect:   "Force any player to pay you $5",
			MTATheme: "Debt Collector",
		})
	}

	// It's My Stop! (It's My Birthday) - 3 cards, value $2
	for i := 1; i <= 3; i++ {
		cards = append(cards, &ActionCard{
			BaseCard: BaseCard{ID: fmt.Sprintf("action_its_my_stop_%d", i), Type: CardTypeAction, Name: "It's My Stop!", Value: 2},
			Effect:   "All players pay you $2",
			MTATheme: "It's My Birthday",
		})
	}

	// Rush Hour (Double the Rent) - 2 cards, value $1
	for i := 1; i <= 2; i++ {
		cards = append(cards, &ActionCard{
			BaseCard: BaseCard{ID: fmt.Sprintf("action_rush_hour_%d", i), Type: CardTypeAction, Name: "Rush Hour", Value: 1},
			Effect:   "Play with a rent card to double the rent amount",
			MTATheme: "Double the Rent",
		})
	}

	// Express Service (House) - 3 cards, value $3
	for i := 1; i <= 3; i++ {
		cards = append(cards, &ActionCard{
			BaseCard: BaseCard{ID: fmt.Sprintf("action_express_service_%d", i), Type: CardTypeAction, Name: "Express Service", Value: 3},
			Effect:   "Add to a complete set to add $3 to rent (not on Railroads/Utilities)",
			MTATheme: "House",
		})
	}

	// New Station (Hotel) - 2 cards, value $4
	for i := 1; i <= 2; i++ {
		cards = append(cards, &ActionCard{
			BaseCard: BaseCard{ID: fmt.Sprintf("action_new_station_%d", i), Type: CardTypeAction, Name: "New Station", Value: 4},
			Effect:   "Add to a complete set with Express Service to add $4 to rent",
			MTATheme: "Hotel",
		})
	}

	return cards
}

func createRentCards() []Card {
	cards := make([]Card, 0, 13)

	// Light Blue / Brown rent - 2 cards, value $1
	for i := 1; i <= 2; i++ {
		cards = append(cards, &RentCard{
			BaseCard:   BaseCard{ID: fmt.Sprintf("rent_blue_brown_%d", i), Type: CardTypeRent, Name: "Rent", Value: 1},
			Colors:     []string{"blue", "brown"},
			ColorNames: []string{"Light Blue", "Brown"},
			Target:     "all",
		})
	}

	// Pink / Orange rent - 2 cards, value $1
	for i := 1; i <= 2; i++ {
		cards = append(cards, &RentCard{
			BaseCard:   BaseCard{ID: fmt.Sprintf("rent_pink_orange_%d", i), Type: CardTypeRent, Name: "Rent", Value: 1},
			Colors:     []string{"pink", "orange"},
			ColorNames: []string{"Pink", "Orange"},
			Target:     "all",
		})
	}

	// Red / Yellow rent - 2 cards, value $1
	for i := 1; i <= 2; i++ {
		cards = append(cards, &RentCard{
			BaseCard:   BaseCard{ID: fmt.Sprintf("rent_red_yellow_%d", i), Type: CardTypeRent, Name: "Rent", Value: 1},
			Colors:     []string{"red", "yellow"},
			ColorNames: []string{"Red", "Yellow"},
			Target:     "all",
		})
	}

	// Dark Blue / Green rent - 2 cards, value $1
	for i := 1; i <= 2; i++ {
		cards = append(cards, &RentCard{
			BaseCard:   BaseCard{ID: fmt.Sprintf("rent_darkblue_green_%d", i), Type: CardTypeRent, Name: "Rent", Value: 1},
			Colors:     []string{"darkblue", "green"},
			ColorNames: []string{"Dark Blue", "Green"},
			Target:     "all",
		})
	}

	// Railroad / Utility rent - 2 cards, value $1
	for i := 1; i <= 2; i++ {
		cards = append(cards, &RentCard{
			BaseCard:   BaseCard{ID: fmt.Sprintf("rent_railroad_utility_%d", i), Type: CardTypeRent, Name: "Rent", Value: 1},
			Colors:     []string{"railroad", "utility"},
			ColorNames: []string{"Railroad", "Utility"},
			Target:     "all",
		})
	}

	// Wild Rent - 3 cards, value $3
	for i := 1; i <= 3; i++ {
		cards = append(cards, &RentCard{
			BaseCard:   BaseCard{ID: fmt.Sprintf("rent_wild_%d", i), Type: CardTypeRent, Name: "Wild Rent", Value: 3},
			Colors:     []string{}, // Any color
			ColorNames: []string{},
			Target:     "one", // Wild rent targets one player only
		})
	}

	return cards
}

func createMoneyCards() []Card {
	cards := make([]Card, 0, 20)

	// $1 - 6 cards
	for i := 1; i <= 6; i++ {
		cards = append(cards, &MoneyCard{
			BaseCard:     BaseCard{ID: fmt.Sprintf("money_1_%d", i), Type: CardTypeMoney, Name: "$1", Value: 1},
			Denomination: 1,
			DisplayValue: "$1",
			Theme:        "Base Fare",
		})
	}

	// $2 - 5 cards
	for i := 1; i <= 5; i++ {
		cards = append(cards, &MoneyCard{
			BaseCard:     BaseCard{ID: fmt.Sprintf("money_2_%d", i), Type: CardTypeMoney, Name: "$2", Value: 2},
			Denomination: 2,
			DisplayValue: "$2",
		})
	}

	// $3 - 3 cards
	for i := 1; i <= 3; i++ {
		cards = append(cards, &MoneyCard{
			BaseCard:     BaseCard{ID: fmt.Sprintf("money_3_%d", i), Type: CardTypeMoney, Name: "$3", Value: 3},
			Denomination: 3,
			DisplayValue: "$3",
		})
	}

	// $4 - 3 cards
	for i := 1; i <= 3; i++ {
		cards = append(cards, &MoneyCard{
			BaseCard:     BaseCard{ID: fmt.Sprintf("money_4_%d", i), Type: CardTypeMoney, Name: "$4", Value: 4},
			Denomination: 4,
			DisplayValue: "$4",
		})
	}

	// $5 - 2 cards
	for i := 1; i <= 2; i++ {
		cards = append(cards, &MoneyCard{
			BaseCard:     BaseCard{ID: fmt.Sprintf("money_5_%d", i), Type: CardTypeMoney, Name: "$5", Value: 5},
			Denomination: 5,
			DisplayValue: "$5",
		})
	}

	// $10 - 1 card
	cards = append(cards, &MoneyCard{
		BaseCard:     BaseCard{ID: "money_10_1", Type: CardTypeMoney, Name: "$10", Value: 10},
		Denomination: 10,
		DisplayValue: "$10",
		Theme:        "Unlimited MetroCard",
	})

	return cards
}
