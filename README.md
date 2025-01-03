# kobo-readstat

![Build Status](https://github.com/timchurchard/readstat/workflows/Test/badge.svg)
![Coverage](https://img.shields.io/badge/Coverage-0-red)
[![License](https://img.shields.io/github/license/timchurchard/readstat)](/LICENSE)
[![Release](https://img.shields.io/github/release/timchurchard/readstat.svg)](https://github.com/timchurchard/readstat/releases/latest)
[![GitHub Releases Stats of readstat](https://img.shields.io/github/downloads/timchurchard/readstat/total.svg?logo=github)](https://somsubhra.github.io/github-release-stats/?username=timchurchard&repository=readstat)

Attempt to collect and show reading statistics from [Kobo](https://uk.kobobooks.com/collections/ereaders) e-reader devices.  This is a proof of concept using my two devices (Kobo Clara 2E and Libra 2 with database version 174).

## Usage

Use the `sync` command to read the Kobo database and write a local json file. And use the `stats` command to read the local json file to make statistics.

```shell
./kobo-readstat sync -d ./testfiles/20231219/libra2/KoboReader.sqlite -s tc_readstat.json
```

And the `stats` command to show stats in text or html

```text
./kobo-readstat stats -y 2023 -s ./tc_readstat.json --mode html --out tc_2023.html
```
![screenshot of html reading stats](.files/tc_2023_report.png "2023 report screenshot")

```text
./kobo-readstat stats -s ./tc_readstat.json -y 2023
Year: 2023
Finished books			: 18
Finished articles		: 286
Time reading books		: 21 days 22 hours 25 minutes 35 seconds
Time reading articles		: 1 days 21 hours 58 minutes 51 seconds
Total time reading		: 23 days 20 hours 24 minutes 26 seconds

----------

January 2023 - Finished books: 2, articles: 16, time spend reading books: 2 days 13 hours 47 minutes 22 seconds and articles: 2 hours 25 minutes 52 seconds
	 finished book: The Subtle Knife - Philip Pullman (Duration: 9h54m9s over 45 Sessions)
	 finished book: Northern Lights - Philip Pullman (Duration: 7h42m12s over 41 Sessions)

February 2023 - Finished books: 1, articles: 13, time spend reading books: 1 days 1 hours 46 minutes 16 seconds and articles: 1 hours 1 minutes 20 seconds
	 finished book: The Amber Spyglass - Philip Pullman (Duration: 16h44m15s over 113 Sessions)

March 2023 - Finished books: 0, articles: 36, time spend reading books: 8 hours 33 minutes 21 seconds and articles: 5 hours 48 minutes 16 seconds

April 2023 - Finished books: 1, articles: 26, time spend reading books: 1 days 9 hours 25 minutes 33 seconds and articles: 4 hours 36 minutes 10 seconds
	 finished book: Matilda - Roald Dahl & Quentin Blake (Duration: 28m40s over 6 Sessions)

May 2023 - Finished books: 1, articles: 26, time spend reading books: 1 days 12 hours 17 minutes 22 seconds and articles: 4 hours 19 minutes 30 seconds
	 finished book: Star's Reach: A Novel of the Deindustrial Future - John Michael Greer (Duration: 7h21m48s over 52 Sessions)

June 2023 - Finished books: 1, articles: 44, time spend reading books: 10 hours 54 minutes 24 seconds and articles: 8 hours 3 minutes 47 seconds
	 finished book: Everything's Eventual: 14 Dark Tales - Stephen King (Duration: 5h14m6s over 39 Sessions)

July 2023 - Finished books: 2, articles: 29, time spend reading books: 2 days 7 hours 2 minutes 24 seconds and articles: 4 hours 42 minutes 25 seconds
	 finished book: The Three-Body Problem, No. 2: Dark Forest - Liu Cixin (Duration: 22h3m16s over 128 Sessions)
	 finished book: The Gallows Pole - Benjamin Myers (Duration: 5h11m38s over 35 Sessions)

August 2023 - Finished books: 2, articles: 15, time spend reading books: 2 days 14 hours 10 minutes 38 seconds and articles: 2 hours 1 minutes 44 seconds
	 finished book: Molesworth - Geoffrey Willans (Duration: 5m39s over 1 Sessions)
	 finished book: Death's End - Cixin Liu (Duration: 26h35m39s over 134 Sessions)

September 2023 - Finished books: 2, articles: 16, time spend reading books: 2 days 19 hours 28 minutes 43 seconds and articles: 2 hours 48 minutes 1 seconds
	 finished book: The Outsider: A Novel - Stephen King (Duration: 21h39m37s over 102 Sessions)
	 finished book: Altered Carbon - Richard K. K. Morgan (Duration: 21h7m24s over 124 Sessions)

October 2023 - Finished books: 2, articles: 22, time spend reading books: 1 days 16 hours 59 minutes 54 seconds and articles: 3 hours 41 minutes 1 seconds
	 finished book: Eating Animals - Jonathan Safran Foer (Duration: 9h27m49s over 62 Sessions)
	 finished book: Lullaby - Chuck Palahniuk (Duration: 7h11m50s over 51 Sessions)

November 2023 - Finished books: 2, articles: 23, time spend reading books: 2 days 19 hours 38 minutes 26 seconds and articles: 3 hours 46 minutes 29 seconds
	 finished book: The Future - Naomi Alderman (Duration: 12h28m9s over 68 Sessions)
	 finished book: Mr Mercedes - Stephen King (Duration: 17h44m15s over 100 Sessions)

December 2023 - Finished books: 2, articles: 20, time spend reading books: 2 days 8 hours 21 minutes 12 seconds and articles: 2 hours 44 minutes 16 seconds
	 finished book: The Neverending Story - Michael Ende (Duration: 7h34m23s over 49 Sessions)
	 finished book: The Green Mile - Stephen King (Duration: 14h10m23s over 83 Sessions)
```