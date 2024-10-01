package models

type Movie struct {
	Id                  int                 `json:"id" bson:"id"`
	Adult               bool                `json:"adult" bson:"adult"`
	BelongsToCollection string              `json:"belongs_to_collection" bson:"belongs_to_collection"`
	Budget              float64             `json:"budget" bson:"budget"`
	Genres              []Genre             `json:"genres" bson:"genres"`
	Homepage            string              `json:"homepage" bson:"homepage"`
	ImdbID              string              `json:"imdb_id" bson:"imdb_id"`
	OriginalLanguage    string              `json:"original_language" bson:"original_language"`
	OriginalTitle       string              `json:"original_title" bson:"original_title"`
	Overview            string              `json:"overview" bson:"overview"`
	PosterPath          string              `json:"poster_path" bson:"poster_path"`
	ProductionCompanies []ProductionCompany `json:"production_companies" bson:"production_companies"`
	ProductionCountries []struct {
		Iso3166_1 string `json:"iso_3166_1" bson:"iso_3166_1"`
		Name      string `json:"name" bson:"name"`
	} `json:"production_countries" bson:"production_countries"`
	ReleaseDate     string           `json:"release_date" bson:"release_date"`
	Revenue         float64          `json:"revenue" bson:"revenue"`
	Runtime         int              `json:"runtime" bson:"runtime"`
	SpokenLanguages []SpokenLanguage `json:"spoken_languages" bson:"spoken_languages"`
	Status          string           `json:"status" bson:"status"`
	Tagline         string           `json:"tagline" bson:"tagline"`
	Title           string           `json:"title" bson:"title"`
	Video           bool             `json:"video" bson:"video"`
	Popularity      float64          `json:"popularity" bson:"popularity"`
	VoteAverage     float64          `json:"vote_average" bson:"vote_average"`
	VoteCount       int              `json:"vote_count" bson:"vote_count"`
	BackdropPath    string           `json:"backdrop_path" bson:"backdrop_path"`
}
