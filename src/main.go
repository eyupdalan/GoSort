package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	"sort"
	"strings"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome ", r.URL.Path[1:])
}

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/sort", sortHandler)
	http.HandleFunc("/multiplesort", multipleSortHandler)
	http.ListenAndServe(":8080", nil)
}

type AboutMe struct {
	Id          int
	Name        string
	Description string
}

func sortHandler(w http.ResponseWriter, r *http.Request) {
	about := prepareTestData()
	sort.Sort(about)

	b, err := json.Marshal(about)

	if (err != nil) {
		panic(err)
	}

	w.Write(b)
}

func parseSortQuery(sortValues string) []string {
	sorts := make([]string, 0)
	if (len(strings.TrimSpace(sortValues)) != 0) {
		sorts = strings.Split(sortValues, "|")
	}
	return sorts;
}

func prepareTestData() AboutMeList {
	aboutList := make([]AboutMe, 0)

	about1 := AboutMe{}
	about1.Id = 3
	about1.Name = "Ali"
	about1.Description = "test1"

	about2 := AboutMe{}
	about2.Id = 3
	about2.Name = "Veli"
	about2.Description = "test2"

	about3 := AboutMe{}
	about3.Id = 1
	about3.Name = "Ahmet"
	about3.Description = "test3"

	about4 := AboutMe{}
	about4.Id = 4
	about4.Name = "Mehmet"
	about4.Description = "test4"

	aboutList = append(aboutList, about1)
	aboutList = append(aboutList, about2)
	aboutList = append(aboutList, about3)
	aboutList = append(aboutList, about4)

	return aboutList
}

type AboutMeList []AboutMe

func (aboutMe AboutMeList) Len() int {
	return len(aboutMe)
}

func (aboutMe AboutMeList) Less(i, j int) bool {
	return aboutMe[i].Id < aboutMe[j].Id
}

func (aboutMe AboutMeList) Swap(i, j int) {
	aboutMe[i], aboutMe[j] = aboutMe[j], aboutMe[i]
}

func multipleSortHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("sort");
	sorts := parseSortQuery(query)

	about := prepareTestData()
	preparedSortData := prepareSortFunctions(sorts, about)
	sort.Sort(preparedSortData)

	b, err := json.Marshal(about)

	if (err != nil) {
		panic(err)
	}

	w.Write(b)
}

func prepareSortFunctions(sortFields []string, aboutMe []AboutMe) AboutMeListWithLessFunctions {
	sortData := AboutMeListWithLessFunctions{}
	sortData.Items = aboutMe
	sortData.Functions = make([]LessSortFunction, 0)
	for _, val := range sortFields {
		asc := string(val[0]) != "-"
		if (!asc) {
			val = val[1:]
		}
		switch val {
		case "id":
			if (asc) {
				function := func(q1, q2 AboutMe) bool {
					return q1.Id < q2.Id
				}
				sortData.Functions = append(sortData.Functions, function)
			} else {
				function := func(q1, q2 AboutMe) bool {
					return q1.Id > q2.Id
				}
				sortData.Functions = append(sortData.Functions, function)
			}

			break
		case "name":
			if (asc) {
				function := func(q1, q2 AboutMe) bool {
					return q1.Name < q2.Name
				}
				sortData.Functions = append(sortData.Functions, function)
			} else {
				function := func(q1, q2 AboutMe) bool {
					return q1.Name > q2.Name
				}
				sortData.Functions = append(sortData.Functions, function)
			}
			break
		}
	}
	return sortData
}

type LessSortFunction func(q1, q2 AboutMe) bool

type AboutMeListWithLessFunctions struct {
	Items     []AboutMe
	Functions []LessSortFunction
}

func (aboutMe AboutMeListWithLessFunctions) Len() int {
	return len(aboutMe.Items)
}

func (aboutMe AboutMeListWithLessFunctions) Less(i, j int) bool {
	q1, q2 := aboutMe.Items[i], aboutMe.Items[j]
	for _, value := range aboutMe.Functions {
		if (value(q1, q2)) {
			return true;
		}
		if (value(q2, q1)) {
			return false;
		}
	}
	return false
}

func (aboutMe AboutMeListWithLessFunctions) Swap(i, j int) {
	aboutMe.Items[i], aboutMe.Items[j] = aboutMe.Items[j], aboutMe.Items[i]
}
