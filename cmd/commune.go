package cmd

import (
	"encoding/json"
	"errors"
	"sort"
	"strconv"

	"github.com/benjamingeer/accicalc/internal/dataset"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var communeCmd *cobra.Command = &cobra.Command{
	Use:   "commune",
	Short: "Generate a CSV file of people involved in traffic accidents in a particular commune.",
	Long: `Generate a CSV file of people involved in traffic accidents in a particular commune.
Example:

accicalc commune --department 94 --commune 33 --cyclists
`,
	Run: func(cmd *cobra.Command, args []string) {
		handleError(pedestrians)
	},
	Args: cobra.NoArgs,
}

type CommuneOpts struct {
	flags                   *pflag.FlagSet
	département             string
	commune                 uint
	includePedestrians      bool
	includeCyclists         bool
	includeOthersInVehicles bool
	includeUnharmed         bool
	limitToMinors           bool
	outputFile              string
}

var communeOpts = CommuneOpts{}

type CatégoriePersonne int

const (
	CatégoriePersonnePiéton CatégoriePersonne = iota
	CatégoriePersonneCycliste
	CatégoriePersonneAutre
)

func (catégoriePersonne CatégoriePersonne) String() string {
	return [...]string{
		"Piéton",
		"Cycliste",
		"Autre",
	}[catégoriePersonne]
}

func (catégoriePersonne CatégoriePersonne) MarshalJSON() ([]byte, error) {
	return json.Marshal(catégoriePersonne.String())
}

func getCatégoriePersonne(usager *dataset.Usager, véhicule *dataset.Véhicule) CatégoriePersonne {
	if usager.CatégorieUsager == dataset.Piéton {
		return CatégoriePersonnePiéton
	} else if usager.CatégorieUsager == dataset.Conducteur &&
		véhicule != nil && véhicule.CatégorieVéhicule == dataset.Bicyclette {
		return CatégoriePersonneCycliste
	} else {
		return CatégoriePersonneAutre
	}
}

type Personne struct {
	Date                       string
	Adresse                    string
	Latitude                   string
	Longitude                  string
	CatégorieDePersonne        CatégoriePersonne
	Gravité                    dataset.Gravité
	AnnéeDeNaissance           int
	Sexe                       dataset.Sexe
	VéhiculeQuiAHeurtéLePiéton string
}

func (personne Personne) AsJson() (string, error) {
	return dataset.ToJson(personne)
}

type PersonneNonPiéton struct {
	Date                string
	Adresse             string
	Latitude            string
	Longitude           string
	CatégorieDePersonne CatégoriePersonne
	Gravité             dataset.Gravité
	AnnéeDeNaissance    int
	Sexe                dataset.Sexe
}

func (personneNonPiéton PersonneNonPiéton) AsJson() (string, error) {
	return dataset.ToJson(personneNonPiéton)
}

type ByDate []Personne

func (slice ByDate) Len() int                  { return len(slice) }
func (slice ByDate) Less(left, right int) bool { return slice[left].Date < slice[right].Date }
func (slice ByDate) Swap(left, right int)      { slice[left], slice[right] = slice[right], slice[left] }

func init() {
	communeCmd.Flags().StringVarP(&communeOpts.département, "department", "p", "", "department code")
	_ = communeCmd.MarkFlagRequired("department")
	communeCmd.Flags().UintVarP(&communeOpts.commune, "commune", "c", 0, "commune number")
	_ = communeCmd.MarkFlagRequired("commune")
	communeCmd.Flags().BoolVarP(&communeOpts.includePedestrians, "pedestrians", "r", false, "include pedestrians")
	communeCmd.Flags().BoolVarP(&communeOpts.includeCyclists, "cyclists", "y", false, "include cyclists")
	communeCmd.Flags().BoolVarP(&communeOpts.includeOthersInVehicles, "other", "t", false, "include other vehicle drivers/passengers")
	communeCmd.Flags().BoolVarP(&communeOpts.includeUnharmed, "unharmed", "i", false, "include unharmed persons")
	communeCmd.Flags().BoolVarP(&communeOpts.limitToMinors, "minors", "m", false, "minors only")
	communeCmd.Flags().StringVarP(&communeOpts.outputFile, "out", "o", "", "output file (defaults to standard out)")
	rootCmd.AddCommand(communeCmd)
	communeOpts.flags = communeCmd.Flags()
}

func pedestrians() error {
	var maybeOutputFile *string

	if !(communeOpts.includePedestrians || communeOpts.includeCyclists || communeOpts.includeOthersInVehicles) {
		return errors.New("no user categories selected")
	}

	if communeOpts.flags.Changed("out") {
		maybeOutputFile = &communeOpts.outputFile
	}

	accidents, err := readAccidents()

	if err != nil {
		return err
	}

	filteredAccidents := dataset.Filter(accidents, func(accident *dataset.Accident) bool {
		return accident.Département == communeOpts.département && *accident.Commune == int(communeOpts.commune)
	})

	var personnes []Personne

	for _, accident := range filteredAccidents {
		for _, véhicule := range accident.Véhicules {
			usagers := dataset.Filter(véhicule.Usagers, includePerson(accident, véhicule))

			for _, usager := range usagers {
				var véhiculeQuiAHeurtéLePiéton string

				if usager.CatégorieUsager == dataset.Piéton {
					véhiculeQuiAHeurtéLePiéton = véhicule.CatégorieVéhicule.String()
				} else {
					véhiculeQuiAHeurtéLePiéton = ""
				}

				personnes = append(personnes,
					Personne{
						Date:                       accident.Date,
						Adresse:                    accident.Adresse,
						Latitude:                   accident.Latitude,
						Longitude:                  accident.Longitude,
						CatégorieDePersonne:        getCatégoriePersonne(usager, véhicule),
						Gravité:                    usager.Gravité,
						AnnéeDeNaissance:           usager.AnnéeNaissance,
						Sexe:                       usager.Sexe,
						VéhiculeQuiAHeurtéLePiéton: véhiculeQuiAHeurtéLePiéton,
					},
				)
			}
		}

		autresUsagers := dataset.Filter(accident.AutresUsagers, includePerson(accident, nil))

		for _, usager := range autresUsagers {
			var véhiculeQuiAHeurtéLePiéton string

			if usager.CatégorieUsager == dataset.Piéton {
				véhiculeQuiAHeurtéLePiéton = dataset.CatégorieVéhiculeIndéterminable.String()
			} else {
				véhiculeQuiAHeurtéLePiéton = ""
			}

			personnes = append(personnes,
				Personne{
					Date:                       accident.Date,
					Adresse:                    accident.Adresse,
					Latitude:                   accident.Latitude,
					Longitude:                  accident.Longitude,
					CatégorieDePersonne:        getCatégoriePersonne(usager, nil),
					Gravité:                    usager.Gravité,
					AnnéeDeNaissance:           usager.AnnéeNaissance,
					Sexe:                       usager.Sexe,
					VéhiculeQuiAHeurtéLePiéton: véhiculeQuiAHeurtéLePiéton,
				},
			)
		}
	}

	sort.Sort(ByDate(personnes))
	var rows []any

	if communeOpts.includePedestrians {
		rows = dataset.ToSliceOfAny(personnes)
	} else {
		var nonPiétons []PersonneNonPiéton

		for _, personne := range personnes {
			nonPiétons = append(nonPiétons,
				PersonneNonPiéton{
					Date:                personne.Date,
					Adresse:             personne.Adresse,
					Latitude:            personne.Latitude,
					Longitude:           personne.Longitude,
					CatégorieDePersonne: personne.CatégorieDePersonne,
					Gravité:             personne.Gravité,
					AnnéeDeNaissance:    personne.AnnéeDeNaissance,
					Sexe:                personne.Sexe,
				},
			)
		}

		rows = dataset.ToSliceOfAny(nonPiétons)
	}

	return dataset.WriteCsv(rows, maybeOutputFile)
}

func includePerson(accident *dataset.Accident, véhicule *dataset.Véhicule) func(usager *dataset.Usager) bool {
	return func(usager *dataset.Usager) bool {
		if communeOpts.includeUnharmed || !(usager.Gravité == dataset.GravitéNonRenseignée || usager.Gravité == dataset.Indemne) {
			catégoriePersonne := getCatégoriePersonne(usager, véhicule)

			return ((communeOpts.includePedestrians && catégoriePersonne == CatégoriePersonnePiéton) ||
				(communeOpts.includeCyclists && catégoriePersonne == CatégoriePersonneCycliste) ||
				(communeOpts.includeOthersInVehicles && catégoriePersonne == CatégoriePersonneAutre)) &&
				(!communeOpts.limitToMinors || wasMinor(usager, accident))
		} else {
			return false
		}
	}
}

func wasMinor(usager *dataset.Usager, accident *dataset.Accident) bool {
	accidentYear, _ := strconv.Atoi(accident.Date[0:4])
	return accidentYear-usager.AnnéeNaissance < 18
}
