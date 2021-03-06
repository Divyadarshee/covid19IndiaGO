package data

var StateCodes = map[string]string{
	"Andhra Pradesh":              "AP",
	"Arunachal Pradesh":           "AR",
	"Assam":                       "AS",
	"Bihar":                       "BR",
	"Chhattisgarh":                "CT",
	"Goa":                         "GA",
	"Gujarat":                     "GJ",
	"Haryana":                     "HR",
	"Himachal Pradesh":            "HP",
	"Jharkhand":                   "JH",
	"Karnataka":                   "KA",
	"Kerala":                      "KL",
	"Madhya Pradesh":              "MP",
	"Maharashtra":                 "MH",
	"Manipur":                     "MN",
	"Meghalaya":                   "ML",
	"Mizoram":                     "MZ",
	"Nagaland":                    "NL",
	"Odisha":                      "OR",
	"Punjab":                      "PB",
	"Rajasthan":                   "RJ",
	"Sikkim":                      "SK",
	"Tamil Nadu":                  "TN",
	"Telangana":                   "TG",
	"Tripura":                     "TR",
	"Uttarakhand":                 "UT",
	"Uttar Pradesh":               "UP",
	"West Bengal":                 "WB",
	"Andaman and Nicobar Islands": "AN",
	"Chandigarh":                  "CH",
	"Dadra and Nagar Haveli and Daman and Diu": "DN",
	"Delhi":             "DL",
	"Jammu and Kashmir": "JK",
	"Ladakh":            "LA",
	"Lakshadweep":       "LD",
	"Puducherry":        "PY",
	"Total":             "TT",
}

var Covid19SourceUrl string = "https://data.covid19india.org/v4/min/data.min.json"

var ReverseGeocodingUrl string = "https://eu1.locationiq.com/v1/reverse.php?key=pk.c8772d5c9e1d9046e5995be8c9edcaa4&lat=%f&lon=%f&format=json"
