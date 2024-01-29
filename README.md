# Enron_Corp_DB_Indexer

1. Download Enron Corp Database from http://www.cs.cmu.edu/~enron/enron_mail_20110402.tgz

2. Install ZincSearch (OpenObserve) from https://openobserve.ai/docs/quickstart/#openobserve-cloud

3. Set OpenObserve variables (just the first time) and executes OpenObserve:
- ZO_ROOT_USER_EMAIL="your user" ZO_ROOT_USER_PASSWORD="your password" ./openobserve

4. OpenObserve is run at http://localhost:5080/ 

5. cd Enron_Corp_DB_Indexer

6. Set env variables:
- export SEARCH_SERVER_URL=http://localhost:5080/api/default (OpenObserve API Url) 
- export SEARCH_SERVER_USERNAME (OpenObserve API username)
- export SEARCH_SERVER_PASSWORD (OpenObserve API password)
- export INDEX_NAME (stream name)

7. Executes go build at terminal

8. Executes ./Enron_Corp_DB_Indexer path_of_database, then indexer starts. (% ./Enron_Corp_DB_Indexer ../enron_mail_20110402)

9. When indexing finished start profiling:
- go tool pprof cpu_profile cpu_profile.pprof
- Verify if graphviz is installed (dot -V)
- Executes at terminal: go tool pprof -png cpu_profile.pprof > cpu_profile.png
- Executes at terminal: go tool pprof -png mem_profile.pprof > mem_profile.png


10. cd Visualizer/View

11. Executes npm run build

12. cd ..

13. Executes go build

14. Executes ./searcher -port 3000

15. Open http://localhost:3000/ and start to search

16. Optimization:
1. Envío de batches de tamaño 100, en lugar de enviar documento por documento a Open Observe
    - Comparando las imágenes cpu_profile_without-optimization.png y cpu_profile.png:
        a. Duración 949.52s vs 122.43s
        b. Para main aunque con la optimización se incrementó la proporción del tiempo de ejecución total, disminuyó el tiempo total que se gastó en la función, igualmente pasó con la función ProcessFile.
        c. Las llamadas al sistema (syscall) disminuyeron en tiempo y proporción del tiempo de ejecución total.
2. Uso de un map en la función formatEmailContent en lugar de crear http headers para formatear los emails.
3. Cambio MarshalIndent a simplemente Marshal. 
