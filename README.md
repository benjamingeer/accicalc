# accicalc: assemble tables of data about traffic accidents in France

This command-line program uses the official open data available at
[Bases de données annuelles des accidents corporels de la circulation routière](https://www.data.gouv.fr/fr/datasets/bases-de-donnees-annuelles-des-accidents-corporels-de-la-circulation-routiere-annees-de-2005-a-2021/) to
assemble easy-to-use CSV files of data about traffic accidents. Currently it is mainly useful
for compiling tables of accidents in which pedestrians or cyclists were injured.
You can filter by *commune* and by the types of users involved (e.g. cyclists),
specify start and end years, and get one CSV file covering the years you are interested in.

To use it, first create directories `2005`, `2006`, etc., under `data`, and download
all the official data files into the corresponding directories.

To compile the program, you will need [Go](https://go.dev/). Type `make` to compile.
Then type `./accicalc help` for instructions.
