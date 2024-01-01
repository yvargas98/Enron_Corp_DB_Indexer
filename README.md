# Enron_Corp_DB_Indexer

1. Download Enron Corp Database from http://www.cs.cmu.edu/~enron/enron_mail_20110402.tgz

2. Install ZincSearch from https://zincsearch.com/

3. Executes go build at terminal

4. Executes $ ./Enron_Corp_DB_Indexer path_of_database, then indexer starts. (% ./Enron_Corp_DB_Indexer ../enron_mail_20110402)

5. Start profiling:
- Verify if graphviz is installed (dot -V)
- Executes at terminal: go tool pprof -png cpu_profile.pprof > cpu_profile.png
- Executes at terminal: go tool pprof -png mem_profile.pprof > mem_profile.png

6. Set env variable:
- export SEARCH_SERVER_PASSWORD

