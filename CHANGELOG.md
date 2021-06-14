# System Webservice

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/).

Please note, that this project, while following numbering syntax, it DOES NOT adhere
to [Semantic Versioning](http://semver.org/spec/v2.0.0.html) rules.

## Types of changes

* ```Added``` for new features.
* ```Changed``` for changes in existing functionality.
* ```Deprecated``` for soon-to-be removed features.
* ```Removed``` for now removed features.
* ```Fixed``` for any bug fixes.
* ```Security``` in case of vulnerabilities.

## [2021.2.3.14] - 2021-06-14

### Added
- copyright
- updated libraries
- added proper locale for index terminal data
- proper aligning top menu on workplaces and index page
- proper loading index page
- software textlogo on every page

### Changed
- Software name in the name of tab
- Buttons in top menu moved 2px down
- Name @ Company changes to 10px
- Durations in index left charts downloaded from backend


## [2021.2.3.11] - 2021-06-11

### Changed

- emphasis for charts set to null

## [2021.2.3.10] - 2021-06-10

### Changed

- cursor the same everywhere
- positioning in pages the same everywhere
- proper chart datetimes

## [2021.2.3.9] - 2021-06-09

### Changed

- left menu only compacted
- info charm showing only information about developer
- proper handling resizing charts
- icons moved around
- proper loading workplaces and index data

## [2021.2.3.8] - 2021-06-08

### Fixed

- proper consumption data recalculating

### Changed
- better/faster calculating production data and consumption data

## [2021.2.3.7] - 2021-06-07

### Added

- basic statistics menu, working with remembering selected

### Fixed

- proper assigned user for downtimes data table II.
- proper width changing for charts (and elsewhere)

## [2021.2.3.3] - 2021-06-03

### Fixed

- proper assigned user for downtimes data table
- proper refreshing index and workplaces page

### Changed

- port color handling, no color, when port color = #000000

## [2021.2.3.2] - 2021-06-02

### Fixed

- sorting of data for timeline chart in the backend

## [2021.2.3.1] - 2021-06-01

### Changed

- improved chart sizing on index page

## [2021.2.2.27] - 2021-05-27

### Changed

- removed magic numbers

## [2021.2.2.26] - 2021-05-26

### Changed

- improved tooltip for charts
- improved digital chart
- improved workplaces page

## [2021.2.2.25] - 2021-05-25

### Changed

- caching user settings improved


## [2021.2.2.24] - 2021-05-24

### Added

- combined chart

## [2021.2.2.20] - 2021-05-20

### Added

- production chart
- automatic chart change when chart switched


## [2021.2.2.19] - 2021-05-19

### Added

- logout button and functionality

## [2021.2.2.18] - 2021-05-18

### Added

- digital chart, complete
- charts: two buttons for faster loading and terminal data loading
- data: one button for compacting table size with proper saving to user config

## [2021.2.2.17] - 2021-05-17

### Changed

- analog chart completed with zooming and panning

## [2021.2.2.13] - 2021-05-13

### Changed

- updated to latest go 1.16.4
- updated to latest go libraries

## [2021.2.2.6] - 2021-05-06

### Changed

- index page updated with proper tooltips
- analog chart updated with terminal data
- analog chart updated with proper tooltips
- analog chart updated for twice as faster chart drawing

## [2021.2.2.5] - 2021-05-05

### Added

- basic digital chart

## [2021.2.2.4] - 2021-05-04

### Added

- complete analog chart

## [2021.2.2.3] - 2021-05-03

### Changed

- updated to latest go libraries
- tested drawing analog chart with more than 4 000 000 data: 10 seconds downloading data, 15 seconds drawing

## [2021.2.1.28] - 2021-04-28

### Added

- loading information

### Changed

- charts changed to plotly

## [2021.2.1.27] - 2021-04-27

### Added

- finished charm for displaying debug information

### Changed

- reloading workplaces silently in the background
- debug button moved to right top
- alarms button moved to right top

## [2021.2.1.26] - 2021-04-26

### Changed

- font Milliard

### Added

- power consumption overview chart
- charm for displaying debug information

## [2021.2.1.14] - 2021-04-14

### Added

- proper locales for index page
- 30 days overview table for index page
- workplace selection for index page

### Changed

- control buttons reverted back to top

## [2021.2.1.13] - 2021-04-13

### Fixed

- settings caching after save change
- export data from data table

### Changed

- gui structure for data

### Added

- workplace selection for index page

## [2021.2.1.7] - 2021-04-07

### Changed

- calendar data calculation changed for actual state_record table

## [2021.2.1.6] - 2021-04-06

### Added

- calendar data calculation make fastest possible

## [2021.2.1.1] - 2021-04-04

### Added

- calculating productivity for calendar overview

### Added

- main page fully working production overview
- main page fully working terminal data overview
- main page partially working calendar overview

## [2021.1.3.30] - 2021-03-30

### Changed

- speed-up: workplaces page load in under 100ms

### Added

- main page fully working production overview
- main page fully working terminal data overview
- main page partially working calendar overview

## [2021.1.3.26] - 2021-03-26

## [2021.1.3.29] - 2021-03-29

### Added

- main index.html positioning and dummy data

## [2021.1.3.26] - 2021-03-26

### Fixed

- workplaces port handling

### Changed

- rendering html data the same for all pages

## [2021.1.3.25] - 2021-03-25

### Changed

- workplaces page design

## [2021.1.3.24] - 2021-03-24

### Changed

- javascript settings code reformat
- go caching code reformat
- go settings (all files) code reformat
- better css formatting for all the tables

### Fixed

- proper saving new workplace ports (when deleted old are found, they are updated)

## [2021.1.3.23] - 2021-03-23

### Changed

- fully completed workplace settings page

## [2021.1.3.22] - 2021-03-22

### Changed

- color selection for states, breakdowns and downtimes
- almost completed workplace settings page
- added gravatar for user

## [2021.1.3.19] - 2021-03-19

### Added

- partially completed workplace settings page

## [2021.1.3.17] - 2021-03-17

### Added

- complete device settings page

## [2021.1.3.16] - 2021-03-16

### Added

- complete breakdowns settings page
- complete downtimes settings page
- complete faults settings page
- complete packages settings page
- complete users settings page
- complete system settings page
- complete user settings page

## [2021.1.3.15] - 2021-03-15

### Added

- complete product settings page
- complete parts settings page
- complete states settings page
- complete workshifts settings page

## [2021.1.3.11] - 2021-03-11

### Added

- complete loading all settings
- complete operation settings page
- complete order settings page

## [2021.1.3.10] - 2021-03-10

### Added

- complete loading alarm settings
- complete editing alarm
- complete saving new alarm

### Changed

- loading setting and data table immediately after change in selection
- code reformat

## [2021.1.3.9] - 2021-03-09

### Added

- loading proper table with proper data on page load
- loading alarm settings

### Changed

- controls moved to top menu

## [2021.1.3.8] - 2021-03-08

### Changed

- selection for tables updated

## [2021.1.3.4] - 2021-03-04

### Changed

- datetime picker for data and charts
- menu UI for data and charts
- proper language for charts
- proper "no data" for charts
- better loading menus
- faster work with locales in background

## [2021.1.3.3] - 2021-03-03

### Added

- chart menu
- basic analog chart

## [2021.1.3.2] - 2021-03-02

### Added

- complete table for alarm table
- complete table for breakdown table
- complete table for downtime table
- complete table for fault table
- complete table for package table
- complete table for part table
- complete table for state table
- complete table for user table
- complete table for system stats table

## [2021.1.2.26] - 2021-02-26

### Changed

- code refactoring
- logging refactoring

## [2021.1.2.25] - 2021-02-25

### Added

- locales for datetime picker
- locales for table
- complete table for orders, with proper UTC search
- remembering last selection from data selection menu
- remembering last chosen workplace from data selection

## [2021.1.2.24] - 2021-02-24

### Changed

- loading table and table data

## [2021.1.2.23] - 2021-02-23

### Added

- complete data page selection functionality (menu, datetime, workplaces), downloaded from backend
- first part of backend functionality for data page (what to download from database)

## [2021.1.2.21] - 2021-02-21

### Changed

- updated to latest libraries

## [2021.1.2.19] - 2021-02-19

### Added

- navigation menu for data

## [2021.1.2.17] - 2021-02-17

### Changed

- proper menu behavior
- proper workplace panel behavior
- mt4cloak for proper loading behavior

### Added

- remembering menu collapsed or expanded
- remembering workplace panels opened or closed
- realtime data for workplaces

## [2021.1.2.15] - 2021-02-15

### Added

- workplace overview

## [2021.1.2.12] - 2021-02-13

### Added

- responsive
- locales

## [2021.1.2.12] - 2021-02-12

### Added

- project structure
- basic authorization
