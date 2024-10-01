package constants

const MOVIES_PER_PAGE = 20
const USERID_FOR_TESTING = 20000
const DISCUSSION_PER_PAGE = int(20)
const MOVIE_JSON_FOR_TESTING_POPULAR = `[
    {
        "id": 680,
        "adult": false,
        "belongs_to_collection": "",
        "budget": 8500000,
        "genres": [
            {
                "id": 53,
                "name": "Thriller"
            },
            {
                "id": 80,
                "name": "Crime"
            }
        ],
        "homepage": "https://www.miramax.com/movie/pulp-fiction/",
        "imdb_id": "tt0110912",
        "original_language": "en",
        "original_title": "Pulp Fiction",
        "overview": "A burger-loving hit man, his philosophical partner, a drug-addled gangster's moll and a washed-up boxer converge in this sprawling, comedic crime caper. Their adventures unfurl in three stories that ingeniously trip back and forth in time.",
        "poster_path": "/d5iIlFn5s0ImszYzBPb8JPIfbXD.jpg",
        "production_companies": [
            {
                "id": 14,
                "logo_path": "/m6AHu84oZQxvq7n1rsvMNJIAsMu.png",
                "name": "Miramax",
                "origin_country": "US"
            },
            {
                "id": 59,
                "logo_path": "/yH7OMeSxhfP0AVM6iT0rsF3F4ZC.png",
                "name": "A Band Apart",
                "origin_country": "US"
            },
            {
                "id": 216,
                "logo_path": "",
                "name": "Jersey Films",
                "origin_country": "US"
            }
        ],
        "production_countries": [
            {
                "iso_3166_1": "US",
                "name": "United States of America"
            }
        ],
        "release_date": "1994-09-10",
        "revenue": 213900000,
        "runtime": 154,
        "spoken_languages": [
            {
                "english_name": "English",
                "iso_639_1": "en",
                "name": "English"
            },
            {
                "english_name": "Spanish",
                "iso_639_1": "es",
                "name": "Español"
            },
            {
                "english_name": "French",
                "iso_639_1": "fr",
                "name": "Français"
            }
        ],
        "status": "Released",
        "tagline": "Just because you are a character doesn't mean you have character.",
        "title": "Pulp Fiction",
        "video": false,
        "popularity": 140.950236,
        "vote_average": 8.3,
        "vote_count": 8670,
        "backdrop_path": "/suaEOtk1N1sgg2MTM7oZd2cfVp3.jpg"
    }
]
	`
const MOVIE_JSON_FOR_TESTING = `{
	"adult": false,
	"belongs_to_collection": "Heat Collection",
	"budget": 60000000,
	"genres": [
	  {
		"id": 28,
		"name": "Action"
	  },
	  {
		"id": 80,
		"name": "Crime"
	  },
	  {
		"id": 18,
		"name": "Drama"
	  }
	]
	,"homepage": "https://www.20thcenturystudios.com/movies/heat",
	"id": 1,
	"imdb_id": "tt0113277",
	"original_language": "en",
	"original_title": "Heat",
	"overview": "Obsessive master thief Neil McCauley leads a top-notch crew on various daring heists throughout Los Angeles while determined detective Vincent Hanna pursues him without rest. Each man recognizes and respects the ability and the dedication of the other even though they are aware their cat-and-mouse game may end in violence.",
	"poster_path": "/umSVjVdbVwtx5ryCA2QXL44Durm.jpg",
	"production_companies": [
	  {
		"id": 508,
		"logo_path": "/7cxRWzi4LsVm4Utfpr1hfARNurT.png",
		"name": "Regency Enterprises",
		"origin_country": "US"
	  },
	  {
		"id": 675,
		"logo_path": "",
		"name": "Forward Pass",
		"origin_country": "US"
	  },
	  {
		"id": 174,
		"logo_path": "/IuAlhI9eVC9Z8UQWOIDdWRKSEJ.png",
		"name": "Warner Bros. Pictures",
		"origin_country": "US"
	  },
	  {
		"id": 4700,
		"logo_path": "/A32wmjrs9Psf4zw0uaixF0GXfxq.png",
		"name": "The Linson Company",
		"origin_country": "US"
	  },
	  {
		"id": 676,
		"logo_path": "",
		"name": "Monarchy Enterprises B.V.",
		"origin_country": ""
	  },
	  {
		"id": 10104,
		"logo_path": "/wRn5HnYMGeJHmItRypgOZOC6gwY.png",
		"name": "New Regency Pictures",
		"origin_country": "US"
	  }
	],
	"production_countries": [
	  {
		"iso_3166_1": "US",
		"name": "United States of America"
	  }
	],
	"release_date": "1995-12-15",
	"revenue": 187436818,
	"runtime": 170,
	"spoken_languages": [
	  {
		"english_name": "English",
		"iso_639_1": "en",
		"name": "English"
	  },
	  {
		"english_name": "Spanish",
		"iso_639_1": "es",
		"name": "Español"
	  }
	],
	"status": "Released",
	"tagline": "A Los Angeles crime saga.",
	"title": "Heat",
	"video": false,
	"popularity": 17.924927,
	"vote_average": 7.7,
	"vote_count": 1886,
	"backdrop_path": "/X7zKxmyrVmYCfcyvDgVLH8iNzA.jpg"
  }`
const MOVIE_DUPLICATE_FOR_TESTING = `{
	"adult": false,
	"belongs_to_collection": "Heat Collection",
	"budget": 60000000,
	"genres": [
	  {
		"id": 28,
		"name": "Action"
	  },
	  {
		"id": 80,
		"name": "Crime"
	  },
	  {
		"id": 18,
		"name": "Drama"
	  }
	]
	,"homepage": "https://www.20thcenturystudios.com/movies/heat",
	"id": 949,
	"imdb_id": "tt0113277",
	"original_language": "en",
	"original_title": "Heat",
	"overview": "Obsessive master thief Neil McCauley leads a top-notch crew on various daring heists throughout Los Angeles while determined detective Vincent Hanna pursues him without rest. Each man recognizes and respects the ability and the dedication of the other even though they are aware their cat-and-mouse game may end in violence.",
	"poster_path": "/umSVjVdbVwtx5ryCA2QXL44Durm.jpg",
	"production_companies": [
	  {
		"id": 508,
		"logo_path": "/7cxRWzi4LsVm4Utfpr1hfARNurT.png",
		"name": "Regency Enterprises",
		"origin_country": "US"
	  },
	  {
		"id": 675,
		"logo_path": "",
		"name": "Forward Pass",
		"origin_country": "US"
	  },
	  {
		"id": 174,
		"logo_path": "/IuAlhI9eVC9Z8UQWOIDdWRKSEJ.png",
		"name": "Warner Bros. Pictures",
		"origin_country": "US"
	  },
	  {
		"id": 4700,
		"logo_path": "/A32wmjrs9Psf4zw0uaixF0GXfxq.png",
		"name": "The Linson Company",
		"origin_country": "US"
	  },
	  {
		"id": 676,
		"logo_path": "",
		"name": "Monarchy Enterprises B.V.",
		"origin_country": ""
	  },
	  {
		"id": 10104,
		"logo_path": "/wRn5HnYMGeJHmItRypgOZOC6gwY.png",
		"name": "New Regency Pictures",
		"origin_country": "US"
	  }
	],
	"production_countries": [
	  {
		"iso_3166_1": "US",
		"name": "United States of America"
	  }
	],
	"release_date": "1995-12-15",
	"revenue": 187436818,
	"runtime": 170,
	"spoken_languages": [
	  {
		"english_name": "English",
		"iso_639_1": "en",
		"name": "English"
	  },
	  {
		"english_name": "Spanish",
		"iso_639_1": "es",
		"name": "Español"
	  }
	],
	"status": "Released",
	"tagline": "A Los Angeles crime saga.",
	"title": "Heat",
	"video": false,
	"popularity": 17.924927,
	"vote_average": 7.7,
	"vote_count": 1886,
	"backdrop_path": "/X7zKxmyrVmYCfcyvDgVLH8iNzA.jpg"
  }`
const MOVIE_NON_EXISTENT = `{
	"adult": false,
	"belongs_to_collection": "Heat Collection",
	"budget": 60000000,
	"genres": [
	  {
		"id": 28,
		"name": "Action"
	  },
	  {
		"id": 80,
		"name": "Crime"
	  },
	  {
		"id": 18,
		"name": "Drama"
	  }
	]
	,"homepage": "https://www.20thcenturystudios.com/movies/heat",
	"id": 20000,
	"imdb_id": "tt0113277",
	"original_language": "en",
	"original_title": "Heat",
	"overview": "Obsessive master thief Neil McCauley leads a top-notch crew on various daring heists throughout Los Angeles while determined detective Vincent Hanna pursues him without rest. Each man recognizes and respects the ability and the dedication of the other even though they are aware their cat-and-mouse game may end in violence.",
	"poster_path": "/umSVjVdbVwtx5ryCA2QXL44Durm.jpg",
	"production_companies": [
	  {
		"id": 508,
		"logo_path": "/7cxRWzi4LsVm4Utfpr1hfARNurT.png",
		"name": "Regency Enterprises",
		"origin_country": "US"
	  },
	  {
		"id": 675,
		"logo_path": "",
		"name": "Forward Pass",
		"origin_country": "US"
	  },
	  {
		"id": 174,
		"logo_path": "/IuAlhI9eVC9Z8UQWOIDdWRKSEJ.png",
		"name": "Warner Bros. Pictures",
		"origin_country": "US"
	  },
	  {
		"id": 4700,
		"logo_path": "/A32wmjrs9Psf4zw0uaixF0GXfxq.png",
		"name": "The Linson Company",
		"origin_country": "US"
	  },
	  {
		"id": 676,
		"logo_path": "",
		"name": "Monarchy Enterprises B.V.",
		"origin_country": ""
	  },
	  {
		"id": 10104,
		"logo_path": "/wRn5HnYMGeJHmItRypgOZOC6gwY.png",
		"name": "New Regency Pictures",
		"origin_country": "US"
	  }
	],
	"production_countries": [
	  {
		"iso_3166_1": "US",
		"name": "United States of America"
	  }
	],
	"release_date": "1995-12-15",
	"revenue": 187436818,
	"runtime": 170,
	"spoken_languages": [
	  {
		"english_name": "English",
		"iso_639_1": "en",
		"name": "English"
	  },
	  {
		"english_name": "Spanish",
		"iso_639_1": "es",
		"name": "Español"
	  }
	],
	"status": "Released",
	"tagline": "A Los Angeles crime saga.",
	"title": "Heat",
	"video": false,
	"popularity": 17.924927,
	"vote_average": 7.7,
	"vote_count": 1886,
	"backdrop_path": "/X7zKxmyrVmYCfcyvDgVLH8iNzA.jpg"
  }`

const MOVIE_EDIT = `{
	"adult": false,
	"belongs_to_collection": "Heat Collection",
	"budget": 60000000,
	"genres": [
	  {
		"id": 28,
		"name": "Action"
	  },
	  {
		"id": 80,
		"name": "Crime"
	  },
	  {
		"id": 18,
		"name": "Drama"
	  }
	]
	,"homepage": "https://www.20thcenturystudios.com/movies/heat",
	"id": 1,
	"imdb_id": "tt0113277",
	"original_language": "en",
	"original_title": "TESTEDIT",
	"overview": "Obsessive master thief Neil McCauley leads a top-notch crew on various daring heists throughout Los Angeles while determined detective Vincent Hanna pursues him without rest. Each man recognizes and respects the ability and the dedication of the other even though they are aware their cat-and-mouse game may end in violence.",
	"poster_path": "/umSVjVdbVwtx5ryCA2QXL44Durm.jpg",
	"production_companies": [
	  {
		"id": 508,
		"logo_path": "/7cxRWzi4LsVm4Utfpr1hfARNurT.png",
		"name": "Regency Enterprises",
		"origin_country": "US"
	  },
	  {
		"id": 675,
		"logo_path": "",
		"name": "Forward Pass",
		"origin_country": "US"
	  },
	  {
		"id": 174,
		"logo_path": "/IuAlhI9eVC9Z8UQWOIDdWRKSEJ.png",
		"name": "Warner Bros. Pictures",
		"origin_country": "US"
	  },
	  {
		"id": 4700,
		"logo_path": "/A32wmjrs9Psf4zw0uaixF0GXfxq.png",
		"name": "The Linson Company",
		"origin_country": "US"
	  },
	  {
		"id": 676,
		"logo_path": "",
		"name": "Monarchy Enterprises B.V.",
		"origin_country": ""
	  },
	  {
		"id": 10104,
		"logo_path": "/wRn5HnYMGeJHmItRypgOZOC6gwY.png",
		"name": "New Regency Pictures",
		"origin_country": "US"
	  }
	],
	"production_countries": [
	  {
		"iso_3166_1": "US",
		"name": "United States of America"
	  }
	],
	"release_date": "1995-12-15",
	"revenue": 187436818,
	"runtime": 170,
	"spoken_languages": [
	  {
		"english_name": "English",
		"iso_639_1": "en",
		"name": "English"
	  },
	  {
		"english_name": "Spanish",
		"iso_639_1": "es",
		"name": "Español"
	  }
	],
	"status": "Released",
	"tagline": "A Los Angeles crime saga.",
	"title": "Heat",
	"video": false,
	"popularity": 17.924927,
	"vote_average": 7.7,
	"vote_count": 1886,
	"backdrop_path": "/X7zKxmyrVmYCfcyvDgVLH8iNzA.jpg"
  }`
