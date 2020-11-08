GO = go

DOWNLOAD = wget --quiet --output-document
UNZIP = unzip -d
MKDIR = mkdir -p


.PHONY: all
all: download-icons download-libs data build

.PHONY: clean
clean:
	$(RM) -r datetime data/cities*.json third-party/ assets/js/third-party/ templates/icon_*.svg


.PHONY: build
build: datetime

datetime: *.go
	$(GO) build -tags "$(TAGS)" -v -o datetime

.PHONY: test
test:
	$(GO) test -cover -bench=. -v


DATASETS = \
	third-party/cities15000.txt \
	third-party/admin1CodesASCII.txt \
	third-party/countryInfo.txt

.PHONY: data
data: data/cities.json

data/cities.json: $(DATASETS) scripts/data.go
	cd scripts && $(GO) run data.go

third-party/cities15000.txt:
	$(MKDIR) third-party/
	$(DOWNLOAD) third-party/cities15000.zip http://download.geonames.org/export/dump/cities15000.zip
	$(UNZIP) third-party/ third-party/cities15000.zip
	$(RM) third-party/cities15000.zip

third-party/countryInfo.txt:
	$(MKDIR) third-party/
	$(DOWNLOAD) third-party/countryInfo.txt http://download.geonames.org/export/dump/countryInfo.txt

third-party/admin1CodesASCII.txt:
	$(MKDIR) third-party/
	$(DOWNLOAD) third-party/admin1CodesASCII.txt https://download.geonames.org/export/dump/admin1CodesASCII.txt

.PHONY: download-libs
download-libs: assets/js/third-party/luxon.min.js

assets/js/third-party/luxon.min.js:
	$(MKDIR) assets/js/third-party/
	$(DOWNLOAD) assets/js/third-party/luxon.min.js https://cdn.jsdelivr.net/npm/luxon@1.25.0/build/global/luxon.min.js


ICONS = \
	solid_sun \
	solid_moon \
	solid_adjust

.PHONY: download-icons
download-icons: $(foreach icon,$(ICONS),templates/icon_$(icon).svg)

.DELETE_ON_ERROR: templates/icon_%.svg
templates/icon_%.svg:
	$(DOWNLOAD) $@ https://github.com/FortAwesome/Font-Awesome/raw/5.15.1/svgs/$(subst _,/,$*).svg
