package router

type TwitterInfo struct {
	Card string
	Site string
	Creator string
}

type OpenGraphInfo struct {
	Type string
}

type Info struct {
	Title string
	Description string
	Keywords string
	Author string
	Canonical string
	Image string
	URL string
	Twitter TwitterInfo
	OpenGraph OpenGraphInfo
}

type InfoParameters struct {
	Title string
	Description string
	Keywords string
	Author string
	Canonical string
	Image string
	Twitter TwitterInfo
	OpenGraph OpenGraphInfo
}
