
.PHONY: all
all: download

.PHONY: download
download: data/cities15000.txt data/countryInfo.txt data/admin1CodesASCII.txt

data/cities15000.txt:
	mkdir -p data/
	wget http://download.geonames.org/export/dump/cities15000.zip -O data/cities15000.zip
	unzip data/cities15000.zip -d data/
	$(RM) data/cities15000.zip

data/countryInfo.txt:
	mkdir -p data/
	wget http://download.geonames.org/export/dump/countryInfo.txt -O data/countryInfo.txt

data/admin1CodesASCII.txt:
	mkdir -p data/
	wget https://download.geonames.org/export/dump/admin1CodesASCII.txt -O data/admin1CodesASCII.txt
