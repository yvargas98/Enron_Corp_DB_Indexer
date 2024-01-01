# Enron_Corp_DB_Indexer

1. Download Enron Corp Database from http://www.cs.cmu.edu/~enron/enron_mail_20110402.tgz

2. Install ZincSearch (OpenObserve) from https://zincsearch.com/

3. Set OpenObserve variables and executes OpenObserve:
- ZO_ROOT_USER_EMAIL="your user" ZO_ROOT_USER_PASSWORD="your password" ./openobserve

4. OpenObserve is run at http://localhost:5080/ 

5. cd Enron_Corp_DB_Indexer

6. Set env variables:
- export SEARCH_SERVER_URL (OpenObserve API Url)
- export SEARCH_SERVER_USERNAME (OpenObserve API username)
- export SEARCH_SERVER_PASSWORD (OpenObserve API password)
- export INDEXER_URL (OpenObserve url for create index)

7.. Executes go build at terminal

8. Executes $ ./Enron_Corp_DB_Indexer path_of_database, then indexer starts. (% ./Enron_Corp_DB_Indexer ../enron_mail_20110402)

9. Start profiling:
- Verify if graphviz is installed (dot -V)
- Executes at terminal: go tool pprof -png cpu_profile.pprof > cpu_profile.png
- Executes at terminal: go tool pprof -png mem_profile.pprof > mem_profile.png


10. cd Visualizer/View

11. Executes npm run build

12. cd ..

13. Executes go build

14. Executes ./searcher -port 3000

15. Open http://localhost:3000/ and start to search



