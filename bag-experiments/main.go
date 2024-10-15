package main

import (
	"fmt"

	"github.com/GopherML/bag"
)

func main() {
	bag := buildModel()
	testSample(bag, "HONDA PMT 8005426632 C24249")
	testSample(bag, "HIMALCHULI RESTAURANT  MADISON       WI")
	testSample(bag, "STATE FARM  INSURANCE  800-956-6310  IL")
	testSample(bag, "SPOTIFY")

}

func buildModel() *bag.Bag {

	var t bag.TrainingSet
	t.Samples = bag.SamplesByLabel{
		"Home": {
			"Mortgage",
			"Rent",
			"Homeowners Insurance",
			"Property Taxes",
			"MORTG",
		},
		"Auto": {
			"Honda",
			"Cadillac",
			"Mazda",
			"State Farm",
			"Progressive",
			"Geico",
			"Car Insurance",
			"STATE FARM  INSURANCE  800-956-6310  IL",
			"BP",
			"Mobil",
			"Shell",
			"Gas Station",
			"Oil Change",
			"Valvoline",
			"Jiffy Lube",
			"Car Wash",
			"Car",
		},
		"Utilities": {
			"Electricity",
			"MGE",
			"Madison Gas and Electric",
			"We Energies",
			"Water",
			"Internet",
			"Comcast",
			"Charter",
			"AT&T",
			"Verizon",
			"TDS",
			"Phone",
			"Cell Phone",
			"Cell Phone Bill",
			"Internet Service",
			"Cable",
			"Television",
		},
		"Food": {
			"Groceries",
			"Hy-Vee",
			"Woodman's",
			"Pick 'n Save",
			"Trader Joe's",
			"Whole Foods",
			"Costco",
			"Sam's Club",
			"Restaurants",
			"Mcdonald's",
			"Subway",
			"Chipotle",
			"Qdoba",
			"Panera",
			"Starbucks",
			"Starbucks Coffee",
			"Firehouse Subs",
			"Jimmy John's",
			"Milio's",
			"Domino's",
			"Pizza Hut",
			"Papa John's",
			"Little Caesars",
			"KFC",
			"Taco Bell",
			"Arby's",
			"Hardee's",
			"Carl's Jr",
			"Burger King",
			"Wendy's",
			"Chick-fil-A",
			"Fast Food",
			"Delivery",
			"Food Delivery",
			"DoorDash",
			"Grubhub",
			"Uber Eats",
			"Postmates",
			"Coffee shop",
			"Coffeehouse",
			"Dunkin'",
			"Dunkin' Donuts",
			"Starbucks",
			"Bar",
			"Alcohol",
		},
		"Health": {
			"Health Insurance",
			"Medical",
			"Dental",
			"Vision",
			"Prescriptions",
			"Pharmacy",
			"CVS",
			"Walgreens",
			"Doctor",
			"Chiropractor",
			"Therapist",
			"Therapy",
			"Medication",
			"Vitamins",
			"Supplements",
			"Health",
			"Prescription",
			"Medicine",
			"Medical Bill",
			"Medical Expense",
			"Medical Insurance",
			"Glasses",
			"Contacts",
			"Contact Lenses",
			"Eye Doctor",
			"Eye Exam",
			"Dentist",
			"Dental Bill",
		},
		"Shopping": {
			"Amazon",
			"Target",
			"Walmart",
			"Nordstrom",
			"Madewell",
			"J Crew",
		},
		"Entertainment": {
			"Netflix",
			"Spotify",
			"Apple Music",
			"HBO",
			"Hulu",
			"Disney+",
			"Amazon Prime",
			"Prime Video",
			"Prime",
			"YouTube",
			"Concert",
			"Movie",
			"Movie Theater",
			"Concert Ticket",
			"Music",
			"Entertainment",
			"Streaming",
			"Streaming Service",
			"Subscription",
			"Tickets",
			"Event",
			"Live Music",
			"Play",
		},
	}

	categoryBag, err := bag.NewFromTrainingSet(t)
	if err != nil {
		fmt.Println(err)
		panic("AAAH!")
	}

	return categoryBag
}

func testSample(bag *bag.Bag, sample string) {
	// Predict the category of a new sample
	exampleResults := bag.GetResults(sample)
	fmt.Println("Collection of results", exampleResults)
	match := exampleResults.GetHighestProbability()
	fmt.Println("Highest probability:", match)
}
