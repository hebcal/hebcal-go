package zmanim

// Hebcal - A Jewish Calendar Generator
// Copyright (c) 2022 Michael J. Radwin
// Derived from original C version, Copyright (C) 1994-2004 Danny Sadinoff
//
// This program is free software; you can redistribute it and/or
// modify it under the terms of the GNU General Public License
// as published by the Free Software Foundation; either version 2
// of the License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program. If not, see <https://www.gnu.org/licenses/>.

import "strings"

// Location represents a location for Zmanim
type Location struct {
	Name        string  // City name
	CountryCode string  // ISO 3166 two-letter in caps, e.g. "US", "GB", "IL"
	Latitude    float64 // In the range [-90,90]
	Longitude   float64 // In the range [-180,180]
	Elevation   int     // Elevation in meters above sea level (never negative)
	TimeZoneId  string  // timezone identifier such as "America/Los_Angeles" or "Asia/Jerusalem"
}

// NewLocation creates an instance of a Location object.
//
// elevation is the elevation in meters above sea level. It is used only when
// elevation-aware zmanim are requested (see [Zmanim.UseElevation]); negative
// values are clamped to 0.
//
// This function panics if the latitude or longitude are out of range.
func NewLocation(name string, countryCode string, latitude float64, longitude float64, elevation int, tzid string) Location {
	if latitude < -90 || latitude > 90 {
		panic("Latitude out of range [-90,90]")
	}
	if longitude < -180 || longitude > 180 {
		panic("Longitude out of range [-180,180]")
	}
	if elevation < 0 {
		elevation = 0
	}
	return Location{
		Name:        name,
		CountryCode: countryCode,
		Latitude:    latitude,
		Longitude:   longitude,
		Elevation:   elevation,
		TimeZoneId:  tzid,
	}
}

var classicCities = []Location{
	{"Abuja", "NG", 9.05785, 7.49508, 476, "Africa/Lagos"},
	{"Acre", "IL", 32.92814, 35.07647, 8, "Asia/Jerusalem"},
	{"Adelaide", "AU", -34.92866, 138.59863, 61, "Australia/Adelaide"},
	{"Albany", "US", 42.65258, -73.75623, 57, "America/New_York"},
	{"Albuquerque", "US", 35.08449, -106.65114, 1516, "America/Denver"},
	{"Almaty", "KZ", 43.25, 76.91667, 795, "Asia/Almaty"},
	{"Amsterdam", "NL", 52.37403, 4.88969, 12, "Europe/Amsterdam"},
	{"Anaheim", "US", 33.83529, -117.9145, 53, "America/Los_Angeles"},
	{"Anchorage", "US", 61.21806, -149.90028, 0, "America/Anchorage"},
	{"Arad", "IL", 31.25882, 35.21282, 611, "Asia/Jerusalem"},
	{"Arlington TX", "US", 32.73569, -97.10807, 190, "America/Chicago"},
	{"Ashdod", "IL", 31.79213, 34.64966, 27, "Asia/Jerusalem"},
	{"Ashkelon", "IL", 31.66926, 34.57149, 42, "Asia/Jerusalem"},
	{"Ashqelon", "IL", 31.66926, 34.57149, 42, "Asia/Jerusalem"},
	{"Athens", "GR", 37.98376, 23.72784, 86, "Europe/Athens"},
	{"Atlanta", "US", 33.749, -84.38798, 336, "America/New_York"},
	{"Auckland", "NZ", -36.84853, 174.76349, 43, "Pacific/Auckland"},
	{"Aurora", "US", 39.72943, -104.83192, 1651, "America/Denver"},
	{"Austin", "US", 30.26715, -97.74306, 165, "America/Chicago"},
	{"Baghdad", "IQ", 33.34058, 44.40088, 41, "Asia/Baghdad"},
	{"Bakersfield", "US", 35.37329, -119.01871, 129, "America/Los_Angeles"},
	{"Baku", "AZ", 40.37767, 49.89201, 0, "Asia/Baku"},
	{"Baltimore", "US", 39.29038, -76.61219, 35, "America/New_York"},
	{"Bangkok", "TH", 13.75398, 100.50144, 11, "Asia/Bangkok"},
	{"Barcelona", "ES", 41.38879, 2.15899, 47, "Europe/Madrid"},
	{"Basel", "CH", 47.55839, 7.57327, 279, "Europe/Zurich"},
	{"Bat Yam", "IL", 32.02379, 34.75185, 34, "Asia/Jerusalem"},
	{"Baton Rouge|LA", "US", 30.44332, -91.18747, 20, "America/Chicago"},
	{"Beer Sheva", "IL", 31.25181, 34.7913, 285, "Asia/Jerusalem"},
	{"Beersheba", "IL", 31.25181, 34.7913, 286, "Asia/Jerusalem"},
	{"Beijing", "CN", 39.9075, 116.39723, 49, "Asia/Shanghai"},
	{"Berlin", "DE", 52.52437, 13.41053, 43, "Europe/Berlin"},
	{"Bet Shemesh", "IL", 31.73072, 34.99293, 285, "Asia/Jerusalem"},
	{"Birmingham", "GB", 52.48142, -1.89983, 150, "Europe/London"},
	{"Birobidzhan", "RU", 48.79284, 132.92386, 85, "Asia/Vladivostok"},
	{"Bnei Brak", "IL", 32.08074, 34.8338, 50, "Asia/Jerusalem"},
	{"Bogota", "CO", 4.60971, -74.08175, 2582, "America/Bogota"},
	{"Boise", "US", 43.6135, -116.20345, 835, "America/Boise"},
	{"Bolzano", "IT", 46.49067, 11.33982, 257, "Europe/Rome"},
	{"Boston", "US", 42.35843, -71.05977, 38, "America/New_York"},
	{"Bozen", "IT", 46.49067, 11.33982, 257, "Europe/Rome"},
	{"Brisbane", "AU", -27.46794, 153.02809, 40, "Australia/Brisbane"},
	{"Brussels", "BE", 50.85045, 4.34878, 28, "Europe/Brussels"},
	{"Bucharest", "RO", 44.43225, 26.10626, 83, "Europe/Bucharest"},
	{"Budapest", "HU", 47.49801, 19.03991, 104, "Europe/Budapest"},
	{"Buenos Aires", "AR", -34.61315, -58.37723, 31, "America/Argentina/Buenos_Aires"},
	{"Buffalo", "US", 42.88645, -78.87837, 191, "America/New_York"},
	{"Burlington", "US", 44.47588, -73.21207, 62, "America/New_York"},
	{"Cairo", "EG", 30.06263, 31.24967, 25, "Africa/Cairo"},
	{"Calgary", "CA", 51.05011, -114.08529, 1043, "America/Edmonton"},
	{"Cape Town", "ZA", -33.92584, 18.42322, 26, "Africa/Johannesburg"},
	{"Caracas", "VE", 10.48801, -66.87919, 893, "America/Caracas"},
	{"Casablanca", "MA", 33.58831, -7.61138, 26, "Africa/Casablanca"},
	{"Chandler", "US", 33.30616, -111.84125, 368, "America/Phoenix"},
	{"Chapel Hill", "US", 35.9132, -79.05584, 149, "America/New_York"},
	{"Charlotte", "US", 35.22709, -80.84313, 279, "America/New_York"},
	{"Chicago", "US", 41.85003, -87.65005, 180, "America/Chicago"},
	{"Chisinau", "MD", 47.00556, 28.8575, 53, "Europe/Chisinau"},
	{"Chula Vista", "US", 32.64005, -117.0842, 21, "America/Los_Angeles"},
	{"Cincinnati", "US", 39.162, -84.45689, 267, "America/New_York"},
	{"Cleveland", "US", 41.4995, -81.69541, 204, "America/New_York"},
	{"Colorado Springs", "US", 38.83388, -104.82136, 1837, "America/Denver"},
	{"Columbus", "US", 39.96118, -82.99879, 241, "America/New_York"},
	{"Copenhagen", "DK", 55.67594, 12.56553, 14, "Europe/Copenhagen"},
	{"Corpus Christi", "US", 27.80058, -97.39638, 9, "America/Chicago"},
	{"Dallas", "US", 32.78306, -96.80667, 139, "America/Chicago"},
	{"Delhi", "IN", 28.65195, 77.23149, 228, "Asia/Kolkata"},
	{"Denver", "US", 39.73915, -104.9847, 1636, "America/Denver"},
	{"Des Moines", "US", 41.60054, -93.60911, 265, "America/Chicago"},
	{"Detroit", "US", 42.33143, -83.04575, 192, "America/Detroit"},
	{"Dhaka", "BD", 23.7104, 90.40744, 23, "Asia/Dhaka"},
	{"Dimona", "IL", 31.07079, 35.03269, 551, "Asia/Jerusalem"},
	{"Dnipro", "UA", 48.46664, 35.04066, 61, "Europe/Kiev"},
	{"Dortmund", "DE", 51.51494, 7.466, 99, "Europe/Berlin"},
	{"Dresden", "DE", 51.05089, 13.73832, 117, "Europe/Berlin"},
	{"Dubai", "AE", 25.07725, 55.30927, 24, "Asia/Dubai"},
	{"Dublin", "IE", 53.33306, -6.24889, 17, "Europe/Dublin"},
	{"Dundee", "GB", 56.46913, -2.97489, 77, "Europe/London"},
	{"Durban", "ZA", -29.8579, 31.0292, 18, "Africa/Johannesburg"},
	{"Durham", "US", 35.99403, -78.89862, 122, "America/New_York"},
	{"Dusseldorf", "DE", 51.22172, 6.77616, 46, "Europe/Berlin"},
	{"Edmonton", "CA", 53.55014, -113.46871, 610, "America/Edmonton"},
	{"Eilat", "IL", 29.55805, 34.94821, 63, "Asia/Jerusalem"},
	{"El Paso", "US", 31.75872, -106.48693, 1136, "America/Denver"},
	{"Far Rockaway", "US", 40.60538, -73.75513, 6, "America/New_York"},
	{"Fort Wayne", "US", 41.1306, -85.12886, 250, "America/Indiana/Indianapolis"},
	{"Fort Worth", "US", 32.72541, -97.32085, 203, "America/Chicago"},
	{"Frankfurt", "DE", 50.11552, 8.68417, 114, "Europe/Berlin"},
	{"Fremont", "US", 37.54827, -121.98857, 18, "America/Los_Angeles"},
	{"Fresno", "US", 36.74773, -119.77237, 100, "America/Los_Angeles"},
	{"Gibraltar", "GI", 36.14474, -5.35257, 11, "Europe/Gibraltar"},
	{"Glasgow", "GB", 55.86515, -4.25763, 41, "Europe/London"},
	{"Great Neck", "US", 40.80066, -73.72846, 38, "America/New_York"},
	{"Greenlawn", "US", 40.86899, -73.36512, 71, "America/New_York"},
	{"Greensboro", "US", 36.07264, -79.79198, 259, "America/New_York"},
	{"Grenoble", "FR", 45.17869, 5.71479, 223, "Europe/Paris"},
	{"Guadalajara", "MX", 20.66682, -103.39182, 1598, "America/Mexico_City"},
	{"Guangzhou", "CN", 23.11667, 113.25, 15, "Asia/Shanghai"},
	{"Hadera", "IL", 32.44192, 34.9039, 12, "Asia/Jerusalem"},
	{"Haifa", "IL", 32.81841, 34.9885, 40, "Asia/Jerusalem"},
	{"Halifax", "CA", 44.6464, -63.57291, 28, "America/Halifax"},
	{"Hamburg", "DE", 53.55073, 9.99302, 11, "Europe/Berlin"},
	{"Hamilton", "CA", 43.25011, -79.84963, 94, "America/Toronto"},
	{"Hartford", "US", 41.76371, -72.68509, 25, "America/New_York"},
	{"Hawaii", "US", 21.30694, -157.85833, 18, "Pacific/Honolulu"},
	{"Helsinki", "FI", 60.16952, 24.93545, 26, "Europe/Helsinki"},
	{"Henderson", "US", 36.0397, -114.98194, 573, "America/Los_Angeles"},
	{"Herzliya", "IL", 32.16627, 34.82536, 28, "Asia/Jerusalem"},
	{"Holon", "IL", 32.01034, 34.77918, 31, "Asia/Jerusalem"},
	{"Hong Kong", "HK", 22.27832, 114.17469, 59, "Asia/Hong_Kong"},
	{"Honolulu", "US", 21.30694, -157.85833, 17, "Pacific/Honolulu"},
	{"Houston", "US", 29.76328, -95.36327, 30, "America/Chicago"},
	{"Indianapolis", "US", 39.76838, -86.15804, 241, "America/Indiana/Indianapolis"},
	{"Irkutsk", "RU", 52.29778, 104.29639, 428, "Asia/Irkutsk"},
	{"Irvine", "US", 33.66946, -117.82311, 22, "America/Los_Angeles"},
	{"Irving", "US", 32.81402, -96.94889, 151, "America/Chicago"},
	{"Istanbul", "TR", 41.01384, 28.94966, 32, "Europe/Istanbul"},
	{"Jacksonville", "US", 30.33218, -81.65565, 10, "America/New_York"},
	{"Jersey City", "US", 40.72816, -74.07764, 11, "America/New_York"},
	{"Jerusalem", "IL", 31.76904, 35.21633, 786, "Asia/Jerusalem"},
	{"Johannesburg", "ZA", -26.20227, 28.04363, 1767, "Africa/Johannesburg"},
	{"Kaifeng", "CN", 34.7986, 114.30742, 75, "Asia/Shanghai"},
	{"Kaliningrad", "RU", 54.70649, 20.51095, 2, "Europe/Kaliningrad"},
	{"Kansas City", "US", 39.09973, -94.57857, 285, "America/Chicago"},
	{"Karachi", "PK", 24.8608, 67.0104, 8, "Asia/Karachi"},
	{"Kathmandu", "NP", 27.70169, 85.3206, 1296, "Asia/Kathmandu"},
	{"Kazan", "RU", 55.78874, 49.12214, 67, "Europe/Moscow"},
	{"Kfar Saba", "IL", 32.175, 34.90694, 56, "Asia/Jerusalem"},
	{"Kharkiv", "UA", 49.98081, 36.25272, 113, "Europe/Kiev"},
	{"Kiev", "UA", 50.45466, 30.5238, 187, "Europe/Kiev"},
	{"Kiryas Joel", "US", 41.34204, -74.16792, 213, "America/New_York"},
	{"Kiryat Gat", "IL", 31.60998, 34.76422, 133, "Asia/Jerusalem"},
	{"Kyiv", "UA", 50.45466, 30.5238, 184, "Europe/Kiev"},
	{"Kyoto", "JP", 35.02107, 135.75385, 50, "Asia/Tokyo"},
	{"La Paz", "BO", -16.5, -68.15, 3782, "America/La_Paz"},
	{"Lagos", "NG", 6.45407, 3.39467, 10, "Africa/Lagos"},
	{"Lakewood", "US", 40.09789, -74.21764, 24, "America/New_York"},
	{"Las Vegas", "US", 36.17497, -115.13722, 613, "America/Los_Angeles"},
	{"Leeds", "GB", 53.79648, -1.54785, 47, "Europe/London"},
	{"Leipzig", "DE", 51.33962, 12.37129, 115, "Europe/Berlin"},
	{"Lexington KY", "US", 37.98869, -84.47772, 299, "America/New_York"},
	{"Lima", "PE", -12.04318, -77.02824, 153, "America/Lima"},
	{"Lincoln", "US", 40.8, -96.66696, 366, "America/Chicago"},
	{"Livingston", "US", 40.79593, -74.31487, 98, "America/New_York"},
	{"Llandudno", "GB", 53.32498, -3.83148, 11, "Europe/London"},
	{"Lod", "IL", 31.9467, 34.8903, 68, "Asia/Jerusalem"},
	{"London", "GB", 51.50853, -0.12574, 25, "Europe/London"},
	{"Long Beach", "US", 33.76696, -118.18923, 28, "America/Los_Angeles"},
	{"Los Angeles", "US", 34.05223, -118.24368, 96, "America/Los_Angeles"},
	{"Lyon", "FR", 45.74846, 4.84671, 173, "Europe/Paris"},
	{"Madison", "US", 43.07305, -89.40123, 272, "America/Chicago"},
	{"Madrid", "ES", 40.4165, -3.70256, 666, "Europe/Madrid"},
	{"Manchester", "GB", 53.48095, -2.23743, 50, "Europe/London"},
	{"Manila", "PH", 14.6042, 120.9822, 13, "Asia/Manila"},
	{"Marseilles", "FR", 43.29695, 5.38107, 28, "Europe/Paris"},
	{"Medzhybizh", "UA", 49.43658, 27.40907, 290, "Europe/Kiev"},
	{"Melbourne", "AU", -37.814, 144.96332, 25, "Australia/Melbourne"},
	{"Memphis", "US", 35.14953, -90.04898, 84, "America/Chicago"},
	{"Mercer Island", "US", 47.57065, -122.22207, 104, "America/Los_Angeles"},
	{"Mesa", "US", 33.42227, -111.82264, 380, "America/Phoenix"},
	{"Mexico City", "MX", 19.42847, -99.12766, 2240, "America/Mexico_City"},
	{"Miami", "US", 25.77427, -80.19366, 25, "America/New_York"},
	{"Milan", "IT", 45.46427, 9.18951, 127, "Europe/Rome"},
	{"Milwaukee", "US", 43.0389, -87.90647, 199, "America/Chicago"},
	{"Minneapolis", "US", 44.97997, -93.26384, 262, "America/Chicago"},
	{"Minsk", "BY", 53.9, 27.56667, 222, "Europe/Minsk"},
	{"Mississauga", "CA", 43.5789, -79.6583, 160, "America/Toronto"},
	{"Mitzpe Ramon", "IL", 30.60944, 34.80111, 855, "Asia/Jerusalem"},
	{"Modiin", "IL", 31.89825, 35.01051, 240, "Asia/Jerusalem"},
	{"Montevideo", "UY", -34.90328, -56.18816, 33, "America/Montevideo"},
	{"Montreal", "CA", 45.50884, -73.58781, 216, "America/Toronto"},
	{"Moscow", "RU", 55.75222, 37.61556, 144, "Europe/Moscow"},
	{"Mumbai", "IN", 19.07283, 72.88261, 8, "Asia/Kolkata"},
	{"Munich", "DE", 48.13743, 11.57549, 525, "Europe/Berlin"},
	{"Nashville", "US", 36.16589, -86.78444, 169, "America/Chicago"},
	{"Nazareth", "IL", 32.69925, 35.30483, 356, "Asia/Jerusalem"},
	{"Netanya", "IL", 32.33291, 34.85992, 38, "Asia/Jerusalem"},
	{"New Haven", "US", 41.30815, -72.92816, 18, "America/New_York"},
	{"New Orleans", "US", 29.95465, -90.07507, 18, "America/Chicago"},
	{"New York", "US", 40.71427, -74.00597, 57, "America/New_York"},
	{"Newark", "US", 40.73566, -74.17237, 21, "America/New_York"},
	{"Newton", "US", 42.33704, -71.20922, 37, "America/New_York"},
	{"Nice", "FR", 43.70313, 7.26608, 19, "Europe/Paris"},
	{"Norfolk", "US", 36.84681, -76.28522, 13, "America/New_York"},
	{"Oakland", "US", 37.80437, -122.2708, 23, "America/Los_Angeles"},
	{"Odessa", "UA", 46.48572, 30.74383, 36, "Europe/Kiev"},
	{"Oklahoma City", "US", 35.46756, -97.51643, 394, "America/Chicago"},
	{"Omaha", "US", 41.25861, -95.93779, 315, "America/Chicago"},
	{"Orlando", "US", 28.53834, -81.37924, 53, "America/New_York"},
	{"Osaka", "JP", 34.69374, 135.50218, 17, "Asia/Tokyo"},
	{"Ottawa", "CA", 45.41117, -75.69812, 71, "America/Toronto"},
	{"Panama City", "PA", 8.9936, -79.51973, 17, "America/Panama"},
	{"Paris", "FR", 48.85341, 2.3488, 42, "Europe/Paris"},
	{"Passaic", "US", 40.85677, -74.12848, 38, "America/New_York"},
	{"Pawtucket", "US", 41.87871, -71.38256, 0, "America/New_York"},
	{"Perth", "AU", -31.95224, 115.8614, 30, "Australia/Perth"},
	{"Petach Tikvah", "IL", 32.08707, 34.88747, 54, "Asia/Jerusalem"},
	{"Philadelphia", "US", 39.95233, -75.16379, 8, "America/New_York"},
	{"Phoenix", "US", 33.44838, -112.07404, 366, "America/Phoenix"},
	{"Pittsburgh", "US", 40.44062, -79.99589, 239, "America/New_York"},
	{"Plano", "US", 33.01984, -96.69889, 207, "America/Chicago"},
	{"Portland", "US", 45.52345, -122.67621, 15, "America/Los_Angeles"},
	{"Porto Alegre", "BR", -30.03283, -51.23019, 46, "America/Sao_Paulo"},
	{"Poway", "US", 32.96282, -117.03586, 156, "America/Los_Angeles"},
	{"Prague", "CZ", 50.08804, 14.42076, 201, "Europe/Prague"},
	{"Princeton", "US", 40.34872, -74.65905, 72, "America/New_York"},
	{"Providence", "US", 41.82399, -71.41283, 0, "America/New_York"},
	{"Ra'anana", "IL", 32.1836, 34.87386, 49, "Asia/Jerusalem"},
	{"Raleigh", "US", 35.7721, -78.63861, 99, "America/New_York"},
	{"Ramat Gan", "IL", 32.08227, 34.81065, 56, "Asia/Jerusalem"},
	{"Ramla", "IL", 31.92923, 34.86563, 84, "Asia/Jerusalem"},
	{"Regina", "CA", 50.45008, -104.6178, 577, "America/Regina"},
	{"Reno", "US", 39.52963, -119.8138, 1378, "America/Los_Angeles"},
	{"Richmond Hill", "CA", 43.87111, -79.43725, 233, "America/Toronto"},
	{"Richmond", "US", 37.55376, -77.46026, 67, "America/New_York"},
	{"Riga", "LV", 56.946, 24.10589, 6, "Europe/Riga"},
	{"Rio de Janeiro", "BR", -22.90642, -43.18223, 12, "America/Sao_Paulo"},
	{"Rishon LeZiyyon", "IL", 31.97102, 34.78939, 56, "Asia/Jerusalem"},
	{"Riverside", "US", 33.95335, -117.39616, 256, "America/Los_Angeles"},
	{"Rochester", "US", 43.15478, -77.61556, 157, "America/New_York"},
	{"Rome", "IT", 41.89193, 12.51133, 53, "Europe/Rome"},
	{"Rosario", "AR", -32.94682, -60.63932, 38, "America/Argentina/Cordoba"},
	{"Rotterdam", "NL", 51.9225, 4.47917, 11, "Europe/Amsterdam"},
	{"Sacramento", "US", 38.58157, -121.4944, 16, "America/Los_Angeles"},
	{"Safed", "IL", 32.96465, 35.496, 779, "Asia/Jerusalem"},
	{"Saint Louis", "US", 38.62727, -90.19789, 149, "America/Chicago"},
	{"Saint Paul", "US", 44.94441, -93.09327, 243, "America/Chicago"},
	{"Saint Petersburg", "RU", 59.93863, 30.31413, 11, "Europe/Moscow"},
	{"Salzburg", "AT", 47.79941, 13.04399, 435, "Europe/Vienna"},
	{"San Antonio", "US", 29.42412, -98.49363, 205, "America/Chicago"},
	{"San Diego", "US", 32.71533, -117.15726, 20, "America/Los_Angeles"},
	{"San Francisco", "US", 37.77493, -122.41942, 28, "America/Los_Angeles"},
	{"San Jose", "US", 37.33939, -121.89496, 23, "America/Los_Angeles"},
	{"San Juan", "PR", 18.46633, -66.10572, 9, "America/Puerto_Rico"},
	{"San Salvador", "SV", 13.68935, -89.18718, 651, "America/El_Salvador"},
	{"Santa Ana", "US", 33.74557, -117.86783, 40, "America/Los_Angeles"},
	{"Santiago", "CL", -33.45694, -70.64827, 555, "America/Santiago"},
	{"Sao Paulo", "BR", -23.5475, -46.63611, 769, "America/Sao_Paulo"},
	{"Saskatoon", "CA", 52.13238, -106.66892, 485, "America/Regina"},
	{"Scottsdale", "US", 33.50921, -111.89903, 382, "America/Phoenix"},
	{"Sderot", "IL", 31.525, 34.59693, 90, "Asia/Jerusalem"},
	{"Seattle", "US", 47.60621, -122.33207, 56, "America/Los_Angeles"},
	{"Shanghai", "CN", 31.22222, 121.45806, 13, "Asia/Shanghai"},
	{"Shenzhen", "CN", 22.54554, 114.0683, 4, "Asia/Shanghai"},
	{"Singapore", "SG", 1.28967, 103.85007, 20, "Asia/Singapore"},
	{"Spokane", "US", 47.65966, -117.42908, 539, "America/Los_Angeles"},
	{"Stanford", "US", 37.42411, -122.16608, 31, "America/Los_Angeles"},
	{"Stockholm", "SE", 59.32938, 18.06871, 19, "Europe/Stockholm"},
	{"Stockton", "US", 37.9577, -121.29078, 6, "America/Los_Angeles"},
	{"Strasbourg", "FR", 48.58392, 7.74553, 146, "Europe/Paris"},
	{"Stuttgart", "DE", 48.78232, 9.17702, 253, "Europe/Berlin"},
	{"Sudbury", "US", 42.38343, -71.41617, 64, "America/New_York"},
	{"Sydney", "AU", -33.86785, 151.20732, 58, "Australia/Sydney"},
	{"Tacoma", "US", 47.25288, -122.44429, 74, "America/Los_Angeles"},
	{"Tampa", "US", 27.94752, -82.45843, 38, "America/New_York"},
	{"Tashkent", "UZ", 41.26465, 69.21627, 423, "Asia/Tashkent"},
	{"Teaneck", "US", 40.8976, -74.01597, 42, "America/New_York"},
	{"Tehran", "IR", 35.69439, 51.42151, 1179, "Asia/Tehran"},
	{"Tel Aviv", "IL", 32.08088, 34.78057, 15, "Asia/Jerusalem"},
	{"The Hague", "NL", 52.07667, 4.29861, 5, "Europe/Amsterdam"},
	{"Tianjin", "CN", 39.14222, 117.17667, 10, "Asia/Shanghai"},
	{"Tiberias", "IL", 32.79221, 35.53124, 0, "Asia/Jerusalem"},
	{"Tijuana", "MX", 32.5027, -117.00371, 101, "America/Tijuana"},
	{"Tokyo", "JP", 35.6895, 139.69171, 42, "Asia/Tokyo"},
	{"Toledo", "US", 41.66394, -83.55521, 189, "America/New_York"},
	{"Toronto", "CA", 43.70011, -79.4163, 175, "America/Toronto"},
	{"Toulouse", "FR", 43.60426, 1.44367, 153, "Europe/Paris"},
	{"Tucson", "US", 32.22174, -110.92648, 758, "America/Phoenix"},
	{"Tulsa", "US", 36.15398, -95.99277, 228, "America/Chicago"},
	{"Tunis", "TN", 36.81897, 10.16579, 23, "Africa/Tunis"},
	{"Uman", "UA", 48.7501, 30.21944, 219, "Europe/Kiev"},
	{"Vancouver", "CA", 49.24966, -123.11934, 70, "America/Vancouver"},
	{"Vaughan", "CA", 43.8361, -79.49827, 218, "America/Toronto"},
	{"Venice", "IT", 45.43713, 12.33265, 6, "Europe/Rome"},
	{"Vienna", "AT", 48.20849, 16.37208, 191, "Europe/Vienna"},
	{"Virginia Beach", "US", 36.85293, -75.97799, 4, "America/New_York"},
	{"Volgograd", "RU", 48.71939, 44.50183, 66, "Europe/Volgograd"},
	{"Warsaw", "PL", 52.22977, 21.01178, 113, "Europe/Warsaw"},
	{"Washington DC", "US", 38.89511, -77.03637, 6, "America/New_York"},
	{"Wellington", "NZ", -41.28664, 174.77557, 32, "Pacific/Auckland"},
	{"White Plains", "US", 41.03399, -73.76291, 82, "America/New_York"},
	{"Wichita", "US", 37.69224, -97.33754, 402, "America/Chicago"},
	{"Willemstad", "CW", 12.1084, -68.93354, 2, "America/Curacao"},
	{"Windsor", "CA", 42.30008, -83.01654, 190, "America/Toronto"},
	{"Winnipeg", "CA", 49.8844, -97.14704, 242, "America/Winnipeg"},
	{"Woodmere", "US", 40.63205, -73.71263, 11, "America/New_York"},
	{"Worcester", "US", 42.26259, -71.80229, 164, "America/New_York"},
}

// LookupCity returns a Location object of one of 60 "classic" Hebcal city names.
//
// If not found, returns nil.
//
// The following city names are supported:
//
// Ashdod, Atlanta, Austin, Baghdad, Beer Sheva,
// Berlin, Baltimore, Bogota, Boston, Budapest,
// Buenos Aires, Buffalo, Chicago, Cincinnati, Cleveland,
// Dallas, Denver, Detroit, Eilat, Gibraltar, Haifa,
// Hawaii, Helsinki, Houston, Jerusalem, Johannesburg,
// Kiev, La Paz, Livingston, Las Vegas, London, Los Angeles,
// Marseilles, Miami, Minneapolis, Melbourne, Mexico City,
// Montreal, Moscow, New York, Omaha, Ottawa, Panama City,
// Paris, Pawtucket, Petach Tikvah, Philadelphia, Phoenix,
// Pittsburgh, Providence, Portland, Saint Louis, Saint Petersburg,
// San Diego, San Francisco, Sao Paulo, Seattle, Sydney,
// Tel Aviv, Tiberias, Toronto, Vancouver, White Plains,
// Washington DC, Worcester
//
// City name lookup is case-insensitive.
func LookupCity(name string) *Location {
	str := strings.ToLower(name)
	for _, loc := range classicCities {
		candidate := strings.ToLower(loc.Name)
		if candidate == str {
			return &loc
		}
	}
	return nil
}

func AllCities() []Location {
	return classicCities
}
