package cmd

import (
	"fmt"
	"os"

	"github.com/benjamingeer/accicalc/internal/dataset"
	"github.com/spf13/cobra"
)

var rootCmd *cobra.Command = &cobra.Command{
	Use:   "accicalc",
	Short: "Process traffic accident data from data.gouv.fr",
	Long:  `Process traffic accident data from data.gouv.fr.`,
}

type Opts struct {
	dataPath  string
	startYear uint
	endYear   uint
}

var (
	// ldflags set by GoReleaser
	Version = "unknown"
	Commit  = "unknown"

	opts = Opts{}
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&opts.dataPath, "dataPath", "d", "./data", "path to directory of accident data")
	rootCmd.PersistentFlags().UintVarP(&opts.startYear, "startYear", "s", dataset.FirstYear, "first year to process")
	rootCmd.PersistentFlags().UintVarP(&opts.endYear, "endYear", "e", dataset.LastYear, "last year to process")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func handleError(operation func() error) {
	if err := operation(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func readAccidents() (accidents []*dataset.Accident, err error) {
	if opts.startYear < dataset.FirstYear || opts.startYear > dataset.LastYear {
		return nil, fmt.Errorf("invalid start year %v", opts.startYear)
	}

	if opts.endYear < dataset.FirstYear || opts.endYear > dataset.LastYear {
		return nil, fmt.Errorf("invalid end year %v", opts.endYear)
	}

	if opts.startYear > opts.endYear {
		return nil, fmt.Errorf("start year cannot be later than end year")
	}

	var allAccidents []*dataset.Accident

	for year := opts.startYear; year <= opts.endYear; year++ {
		if yearDatasetReader, ok := dataset.YearDatasetReaders[year]; ok {
			fmt.Fprintf(os.Stderr, "Reading data for %v...\n", year)

			accidents, err := yearDatasetReader.ReadCharacteristics(year, opts.dataPath)

			if err != nil {
				return nil, err
			}

			places, err := yearDatasetReader.ReadPlaces(year, opts.dataPath)

			if err != nil {
				return nil, err
			}

			vehicles, err := yearDatasetReader.ReadVehicles(year, opts.dataPath)

			if err != nil {
				return nil, err
			}

			users, err := yearDatasetReader.ReadUsers(year, opts.dataPath)

			if err != nil {
				return nil, err
			}

			accidentMap := make(map[string]*dataset.Accident)

			for _, accident := range accidents {
				if _, exists := accidentMap[accident.IdAccident]; exists {
					return nil, fmt.Errorf("in year %v, accident %v encountered twice", year, accident.IdAccident)
				}

				accidentMap[accident.IdAccident] = accident
			}

			vehicleMap := make(map[string]*dataset.Véhicule)

			for _, vehicle := range vehicles {
				uniqueVehicleId := makeUniqueVehicleId(vehicle.IdAccident, vehicle.IdVéhicule)

				if _, exists := vehicleMap[uniqueVehicleId]; exists {
					return nil, fmt.Errorf(
						"in year %v, in accident %v, vehicle %v encountered twice",
						year,
						vehicle.IdAccident,
						vehicle.IdVéhicule,
					)
				}

				vehicleMap[uniqueVehicleId] = vehicle
			}

			for _, user := range users {
				uniqueVehicleId := makeUniqueVehicleId(user.IdAccident, user.IdVéhicule)

				if vehicle, ok := vehicleMap[uniqueVehicleId]; ok {
					if vehicle.IdAccident != user.IdAccident {
						return nil, fmt.Errorf("in year %v, a user in accident %v has vehicle %v from accident %v",
							year,
							user.IdAccident,
							vehicle.IdVéhicule,
							vehicle.IdAccident,
						)
					}

					vehicle.Usagers = append(vehicle.Usagers, user)
				} else {
					if accident, ok := accidentMap[user.IdAccident]; ok {
						accident.AutresUsagers = append(accident.AutresUsagers, user)
					} else {
						return nil, fmt.Errorf("in year %v, a user has nonexistent accident %v",
							year,
							user.IdAccident,
						)
					}
				}
			}

			for _, vehicle := range vehicles {
				if accident, ok := accidentMap[vehicle.IdAccident]; ok {
					accident.Véhicules = append(accident.Véhicules, vehicle)
				} else {
					return nil, fmt.Errorf("in year %v, vehicle %v has nonexistent accident %v",
						year,
						vehicle.IdVéhicule,
						vehicle.IdAccident,
					)
				}
			}

			for _, place := range places {
				if accident, ok := accidentMap[place.IdAccident]; ok {
					accident.Lieu = place
				} else {
					return nil, fmt.Errorf("in year %v, a place has nonexistent accident %v",
						year,
						place.IdAccident,
					)
				}
			}

			allAccidents = append(allAccidents, accidents...)
		} else {
			return nil, fmt.Errorf("unsupported year %v", year)
		}
	}

	return allAccidents, nil
}

func makeUniqueVehicleId(accidentId string, vehicleId string) string {
	return fmt.Sprintf("%v|%v", accidentId, vehicleId)
}
