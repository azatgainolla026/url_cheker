#  Concurrent URL Status Checker & Report Generator

A simple Go tool for checking a list of URLs using **concurrent requests**. The program sends GET requests to each URL, measures the response time, and saves the report in CSV format.

---

###  How to run it

1. **Clone the repository:**

```bash
git clone https://github.com/azatgainolla026/url_cheker.git
cd url-checker
```
2. **Create a file with a list of URLs (`urls.txt`):**
 Each line in the file will contain one URL
3. **Run the program:**
```bash
go run main.go checker.go -f urls.txt -c 5
```
4. **Expected output:**
![image](https://github.com/user-attachments/assets/b5a5fcc4-4c1b-48e4-a193-2d038dfa53f7)
![image](https://github.com/user-attachments/assets/7297eb4d-2a60-4682-8b3b-a9456d4bc8e3)


### Requirements
Ensure you have Go installed 

