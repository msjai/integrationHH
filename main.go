package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strconv"
)

type ResponseCompanyList struct {
	PerPage int `json:"per_page"`
	Page    int `json:"page"`
	Pages   int `json:"pages"`
	Found   int `json:"found"`
	Items   []struct {
		Id            string `json:"id"`
		Name          string `json:"name"`
		Url           string `json:"url"`
		AlternateUrl  string `json:"alternate_url"`
		VacanciesUrl  string `json:"vacancies_url"`
		OpenVacancies int    `json:"open_vacancies"`
		LogoUrls      struct {
			Field1 string `json:"90"`
		} `json:"logo_urls"`
	} `json:"items"`
}

type ResponseCompanyInfo struct {
	Name               string `json:"name"`
	Type               string `json:"type"`
	Id                 string `json:"id"`
	SiteUrl            string `json:"site_url"`
	Description        string `json:"description"`
	BrandedDescription string `json:"branded_description"`
	VacanciesUrl       string `json:"vacancies_url"`
	OpenVacancies      int    `json:"open_vacancies"`
	Trusted            bool   `json:"trusted"`
	AlternateUrl       string `json:"alternate_url"`
	InsiderInterviews  []struct {
		Url   string `json:"url"`
		Id    string `json:"id"`
		Title string `json:"title"`
	} `json:"insider_interviews"`
	LogoUrls struct {
		Field1   string `json:"90"`
		Field2   string `json:"240"`
		Original string `json:"original"`
	} `json:"logo_urls"`
	Area struct {
		Url  string `json:"url"`
		Id   string `json:"id"`
		Name string `json:"name"`
	} `json:"area"`
	Relations  []interface{} `json:"relations"`
	Industries []struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	} `json:"industries"`
}

func main() {

	var client = CreateClient()
	var resMaxVacancyCompany *http.Response
	var resMaxVacancyCompanyInfo *http.Response

	resMaxVacancyCompany = MakeGetRequest("https://api.hh.ru/employers?only_with_vacancies=true&text='IT'&per_page=100", client)

	TopVacancyCompanyID := MakeRequestTopVacancyCompanyIDCompany(resMaxVacancyCompany)

	resMaxVacancyCompanyInfo = MakeGetRequest("https://api.hh.ru/employers/"+TopVacancyCompanyID, client)

	TopCompanyInfo(resMaxVacancyCompanyInfo)
}

func CreateClient() http.Client {

	client := http.Client{}

	return client
}

func MakeGetRequest(UrlRequest string, client http.Client) *http.Response {

	//client := http.Client{}
	req, err := http.NewRequest("GET", UrlRequest, nil)
	if err != nil {
		log.Fatalln(err)
	}

	req.Header = http.Header{
		"HH-User-Agent": []string{"HH-User-Agent"},
	}

	res, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	//defer res.Body.Close()

	return res
}

func MakeRequestTopVacancyCompanyIDCompany(res *http.Response) string {

	defer res.Body.Close()

	//Чтение тела
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
	}

	//sb := string(body)
	//log.Printf(sb)

	var result ResponseCompanyList
	if err := json.Unmarshal(body, &result); err != nil { // Parse []byte to the go struct pointer
		fmt.Println("Can not unmarshal JSON")
	}

	sort.Slice(result.Items, func(i, j int) bool { return result.Items[i].OpenVacancies < result.Items[j].OpenVacancies })

	for _, rec := range result.Items {
		fmt.Println(rec.Name + " <<Количество вакансий открыто>>: " + strconv.Itoa(rec.OpenVacancies))
	}

	TopVacancyCompanyID := result.Items[len(result.Items)-1].Id

	return TopVacancyCompanyID

}

func TopCompanyInfo(res *http.Response) {

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
	}

	//sb := string(body)
	//log.Printf(sb)

	var result ResponseCompanyInfo
	if err := json.Unmarshal(body, &result); err != nil { // Parse []byte to the go struct pointer
		fmt.Println("Can not unmarshal JSON")
	}

	fmt.Println("")

	if result.Trusted == true {
		fmt.Println("Комания " + result.Name + " заслуживает доверия. Она провереннная!")
	} else {
		fmt.Println("Комания " + result.Name + " НЕ заслуживает доверия. Она не проверенная!")
	}

}
