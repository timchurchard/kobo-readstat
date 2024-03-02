# readstat

![Build Status](https://github.com/timchurchard/readstat/workflows/Test/badge.svg)
![Coverage](https://img.shields.io/badge/Coverage-31.4%25-yellow)
[![License](https://img.shields.io/github/license/timchurchard/readstat)](/LICENSE)
[![Release](https://img.shields.io/github/release/timchurchard/readstat.svg)](https://github.com/timchurchard/readstat/releases/latest)
[![GitHub Releases Stats of readstat](https://img.shields.io/github/downloads/timchurchard/readstat/total.svg?logo=github)](https://somsubhra.github.io/github-release-stats/?username=timchurchard&repository=readstat)

Attempt to collect and show reading statistics from e-reader devicess, like [Kobo](https://uk.kobobooks.com/collections/ereaders).  This is a proof of concept using my two devices (Kobo Clara 2E and Libra 2 with database version 174).

## Usage

Use the `sync` command to read the Kobo database. Syncing updates a local json file. And use the `stats` command to make stats.

```text
./readstat sync -d ./testfiles/20231219/libra2/KoboReader.sqlite -s tc_readstat.json

./readstat stats -s ./tc_readstat.json -y 2023
Year: 2023
Finished books			: 19
Finished articles		: 280
Total finished words		: 3,111,402
Time reading books		: 10 days 14 hours 19 minutes 27 seconds
Time reading articles		: 1 days 21 hours 50 minutes 52 seconds
Total time reading		: 12 days 12 hours 10 minutes 19 seconds

----------

January 2023 - Finished books: 2, articles: 14, time spend reading books: 1 days 4 hours 33 minutes 43 seconds and articles: 2 hours 25 minutes 52 seconds
	 finished book: The Subtle Knife - Philip Pullman (Duration: 9h54m9s over 45 Sessions)
	 finished book: Northern Lights - Philip Pullman (Duration: 7h42m12s over 41 Sessions)

February 2023 - Finished books: 1, articles: 13, time spend reading books: 14 hours 43 minutes 14 seconds and articles: 1 hours 1 minutes 20 seconds
	 finished book: The Amber Spyglass - Philip Pullman (Duration: 14h12m26s over 99 Sessions)

March 2023 - Finished books: 0, articles: 36, time spend reading books: 4 hours 0 minutes 53 seconds and articles: 5 hours 42 minutes 56 seconds

April 2023 - Finished books: 2, articles: 26, time spend reading books: 15 hours 55 minutes 17 seconds and articles: 4 hours 36 minutes 10 seconds
	 finished book: Matilda - Roald Dahl & Quentin Blake (Duration: 10m27s over 3 Sessions)
	 finished book: Lullaby - Chuck Palahniuk (Duration: 7h11m50s over 51 Sessions)

May 2023 - Finished books: 2, articles: 26, time spend reading books: 11 hours 9 minutes 55 seconds and articles: 4 hours 19 minutes 30 seconds
	 finished book: Star's Reach: A Novel of the Deindustrial Future - John Michael Greer (Duration: 1h15m8s over 11 Sessions)
	 finished book: Eating Animals - Jonathan Safran Foer (Duration: 2h35m4s over 19 Sessions)

June 2023 - Finished books: 1, articles: 43, time spend reading books: 5 hours 55 minutes 3 seconds and articles: 8 hours 3 minutes 47 seconds
	 finished book: Everything's Eventual: 14 Dark Tales - Stephen King (Duration: 1h18m4s over 11 Sessions)

July 2023 - Finished books: 2, articles: 29, time spend reading books: 1 days 11 hours 0 minutes 45 seconds and articles: 4 hours 42 minutes 25 seconds
	 finished book: The Three-Body Problem, No. 2: Dark Forest - Liu Cixin (Duration: 22h3m16s over 128 Sessions)
	 finished book: The Gallows Pole - Benjamin Myers (Duration: 42m15s over 5 Sessions)

August 2023 - Finished books: 2, articles: 15, time spend reading books: 22 hours 47 minutes 51 seconds and articles: 2 hours 1 minutes 44 seconds
	 finished book: Molesworth - Geoffrey Willans (Duration: 5m39s over 1 Sessions)
	 finished book: Death's End - Cixin Liu (Duration: 15h32m16s over 81 Sessions)

September 2023 - Finished books: 2, articles: 16, time spend reading books: 1 days 21 hours 2 minutes 25 seconds and articles: 2 hours 45 minutes 22 seconds
	 finished book: The Outsider: A Novel - Stephen King (Duration: 21h39m37s over 102 Sessions)
	 finished book: Altered Carbon - Richard K. K. Morgan (Duration: 14h40m47s over 94 Sessions)

October 2023 - Finished books: 2, articles: 19, time spend reading books: 17 hours 0 minutes 28 seconds and articles: 3 hours 41 minutes 1 seconds
	 finished book: Eating Animals - Jonathan Safran Foer (Duration: 0s over 0 Sessions)
	 finished book: Lullaby - Chuck Palahniuk (Duration: 0s over 0 Sessions)

November 2023 - Finished books: 1, articles: 23, time spend reading books: 1 days 2 hours 49 minutes 36 seconds and articles: 3 hours 46 minutes 29 seconds
	 finished book: Mr Mercedes - Stephen King (Duration: 4h56m53s over 25 Sessions)

December 2023 - Finished books: 2, articles: 20, time spend reading books: 1 days 3 hours 20 minutes 17 seconds and articles: 2 hours 44 minutes 16 seconds
	 finished book: The Green Mile - Stephen King (Duration: 8h53m35s over 47 Sessions)
	 finished book: The Neverending Story - Michael Ende (Duration: 7h33m54s over 48 Sessions)
```
