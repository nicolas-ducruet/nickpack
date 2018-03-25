package nickpack

/********************************************************************
	0. VARIABLES AND INIT
	1. ARRAYS & SLICES
	2. DATABASE
	3. FILESYSTEM
	4. MONITORING
	5. STRINGS
*********************************************************************/

import (
	"encoding/csv"
	"bufio"
	"database/sql"
	"fmt"
	"go/build"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

/********************************************************************
	0. VARIABLES & INIT
*********************************************************************/
var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")


func init() {
}

/********************************************************************
	1. ARRAYS & SLICES
*********************************************************************/

// Returns an array of elements that pass a test function
func FilterArStr(ss []string, test func(string) bool) (ret []string) {
    for _, s := range ss {
        if test(s) {
            ret = append(ret, s)
        }
    }
    return
}

/********************************************************************
	2. DATABASE
*********************************************************************/
// Check if row exists in db
func RowExists(db *sql.DB, query string, args ...interface{}) (bool,error) {
	
	var exists bool
	query = fmt.Sprintf("SELECT exists (%s)", query)

	if err := db.QueryRow(query, args...).Scan(&exists); err != nil && err != sql.ErrNoRows {
		return false,err
	}
	return exists,nil
}

func UpdateRowDB(db *sql.DB, query string,v []interface{}) (bool,error) {
	query = CleanQuery(query)
	
	stmt, err := db.Prepare(query)
	if err != nil {
		return false,err
	}
	res, err2 := stmt.Exec(v...)
	if err2 != nil {
		return false,err2
	}
	if _, err3 := res.RowsAffected(); err3 != nil {
		return false,err3
	}
	return true, nil
}

func InsertRowDB(db *sql.DB, query string,v []interface{}) (bool,error) {
	query = CleanQuery(query)
	stmt, err := db.Prepare(query)
	if err != nil {
		return false,err
	}
	if _, err2 := stmt.Exec(v...); err2 != nil {
		return false,err2
	}
	return true, nil
}

func SeedFromCsv(db *sql.DB, filename string, query string) bool {
	csvIn, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	r := csv.NewReader(csvIn)
	defer csvIn.Close()

	for {
		line, err := r.Read()

		if(len(line) == 0){
			continue
		}

		if err != nil {
			/* End of file -> execute sql statement */
			if err == io.EOF {
				/* Trim the last , */
				query = TrimSuffix(query, ",")
				InsertRowDB(db,query,[]interface{}{})
				return true
			}
			fmt.Println(line)
			panic(err.Error())
			return false
		}

		row := strings.Split(line[0], "#")

		query += "("
		for _, value := range row {
			
			/* Escaping ' quotes */
			value = strings.Replace(value, "'", "", -1)

			if(value == "NULL"){
				query += "NULL,"
				
			} else{
				query += "'" + value + "',"
			}
		}
		query = TrimSuffix(query, ",")
		query += "),"
	}
}

/********************************************************************
	3. FILESYSTEM
*********************************************************************/
// Create new file and write content
func WriteToFile(content, destination string) error {

	f, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(content)
	if err != nil {
		return err
	}

	return err
}

// Log whatever text in whatever existing file
func AppendToFile(nameFile, content string) {
	
	dirPath := build.Default.GOPATH + "/logs"
	logPath := dirPath+"/"+nameFile+".log"

	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		CreateDirIfNotExist(dirPath)
	}

	f, err := os.OpenFile(logPath, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()

	log.SetOutput(f)
	log.Print(content + "\r\n")
}

// Check if file exists
func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

// Download file from url
func DownloadFromUrl(url,filename string) {

	response, e := http.Get(url)
	if e != nil {
		log.Fatal(e)
	}
	defer response.Body.Close()

	//open a file for writing
	file, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	
	_, err = io.Copy(file, response.Body)
	if err != nil {
		log.Fatal(err)
	}
	file.Close()
}

// Create a dir if it doesn't exist
func CreateDirIfNotExist(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
			err = os.MkdirAll(dir, 0755)
			if err != nil {
					panic(err)
			}
	}
}

/********************************************************************
	4. MONITORING
*********************************************************************/
func TimeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	fmt.Printf("%s took %s\n", name, elapsed)
}

/********************************************************************
	5. STRINGS
*********************************************************************/
// Clear mysql queries
func CleanQuery(query string) string{
	var re = regexp.MustCompile(`AND\s+GROUP BY`)
	query = re.ReplaceAllString(query, `GROUP BY`)
	re = regexp.MustCompile(`AND\s+ORDER BY`)
	query = re.ReplaceAllString(query, `ORDER BY`)
	re = regexp.MustCompile(`WHERE\s+GROUP BY`)
	query = re.ReplaceAllString(query, `GROUP BY`)
	re = regexp.MustCompile(`WHERE\s+ORDER BY`)
	query = re.ReplaceAllString(query, `ORDER BY`)
	re = regexp.MustCompile(`AND\s+OR`)
	query = re.ReplaceAllString(query, `AND `)
	re = regexp.MustCompile(`WHERE\s+AND`)
	query = re.ReplaceAllString(query, `WHERE `)
	re = regexp.MustCompile(`AND\s+AND`)
	query = re.ReplaceAllString(query, `AND `)
	re = regexp.MustCompile(`AND\s+()`)
	query = re.ReplaceAllString(query, `AND `)
	query = TrimSuffix(query, ",")
	return query
}

// Remove specific character from end of string
func TrimSuffix(s, suffix string) string {
	if strings.HasSuffix(s, suffix) {
		s = s[:len(s)-len(suffix)]
	}
	return s
}

func RegexReplace(sample,regex,replace string) string {
	var re = regexp.MustCompile(regex)
	s := re.ReplaceAllString(sample, replace)
	return s
}

// Generate random string of runes
func RandStringRunes(n int) string {
    b := make([]rune, n)
    for i := range b {
        b[i] = letterRunes[rand.Intn(len(letterRunes))]
    }
    return string(b)
}

// Generate random int between bounds
func RandInt(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}

// Split string based on regex condition
func RegSplit(text string, delimeter string) []string {
    reg := regexp.MustCompile(delimeter)
    indexes := reg.FindAllStringIndex(text, -1)
    laststart := 0
    result := make([]string, len(indexes) + 1)
    for i, element := range indexes {
            result[i] = text[laststart:element[0]]
            laststart = element[1]
    }
    result[len(indexes)] = text[laststart:len(text)]
    return result
}

// Remove all whitespaces from string
func WSTrim(str string) string {
	return strings.Join(strings.Fields(str), "")
}

// readLines reads a whole file into memory and returns a slice of its lines.
func ReadLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

// writeLines writes the lines to the given file.
func WriteLines(lines []string, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	for _, line := range lines {
		fmt.Fprintln(w, line)
	}
	return w.Flush()
}