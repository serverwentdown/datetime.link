GO = go

DOWNLOAD = wget --quiet --output-document
UNZIP = unzip -d
MKDIR = mkdir -p


.PHONY: all
all: download-icons data build

.PHONY: clean
clean:
	$(RM) -r datetime js/data.json data/ templates/icon_*.svg


.PHONY: build
build: datetime

datetime: *.go
	$(GO) build -v -o datetime

.PHONY: test
test:
	$(GO) test -v


DATASETS = \
	data/cities15000.txt \
	data/admin1CodesASCII.txt \
	data/countryInfo.txt

.PHONY: data
data: js/data.json

js/data.json: $(DATASETS) scripts/data.go
	cd scripts && $(GO) run data.go

data/cities15000.txt:
	$(MKDIR) data/
	$(DOWNLOAD) data/cities15000.zip http://download.geonames.org/export/dump/cities15000.zip
	$(UNZIP) data/ data/cities15000.zip
	$(RM) data/cities15000.zip

data/countryInfo.txt:
	$(MKDIR) data/
	$(DOWNLOAD) data/countryInfo.txt http://download.geonames.org/export/dump/countryInfo.txt

data/admin1CodesASCII.txt:
	$(MKDIR) data/
	$(DOWNLOAD) data/admin1CodesASCII.txt https://download.geonames.org/export/dump/admin1CodesASCII.txt


ICONS = \
	solid_sun \
	solid_moon \
	solid_adjust

.PHONY: download-icons
download-icons: $(foreach icon,$(ICONS),templates/icon_$(icon).svg)

.DELETE_ON_ERROR: templates/icon_%.svg
templates/icon_%.svg:
	$(DOWNLOAD) $@ https://github.com/FortAwesome/Font-Awesome/raw/5.15.1/svgs/$(subst _,/,$*).svg
