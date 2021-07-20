*** Amazon-Product Details Scraper and Collector APIs ***

**Pre-requirements** 
1. Docker Version 20.10.0 or Higher     **link to download** = docker.com/products/docker-desktop
2. Docker-Compose Version 1.27.4 or Higher **link to download** = docs.docker.com/compose/install/
3. postman

**Steps for set up**
1. Clone the repo locally and open it in terminal
2. build code using docker-compose, run command "docker-compose build"
3. run the application in background using command "docker-compose up -d"

note: if above commands are not working then try with sudo

**Make API call Using postman**
1. use these URLs for GET/POST Methods
    1.1 "localhost:3030/scraper" 
        1.1.1 for POST Method (for scraping data) need to provide url in body as json obj
           for e.g. 
           {
               "url":"https://www.amazon.com/gp/aw/d/B00DUARBTA/ref=sspa_mw_detail_3?ie=UTF8&psc=1&pd_rd_i=B00DUARBTAp13NParams&smid=A15IBQIC022W2F"
           }
    1.2 "localhost:3031/collector"
        1.2.1 for GET Method (for viewing the scrapped data) use this url


**close services**
1. run command "docker-compose down"

**Used Techs**
1. Used "gocolly/colly" library for scrapping the data
2. golang
3. Mongodb
4. docker