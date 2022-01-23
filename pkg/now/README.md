# Bhojpur ORM - Time Library

The `now` package is a time toolkit for Bhojpur ORM

## Installation

```
go get -u github.com/bhojpur/orm
```

## Usage

Calculating time based on the current time

```go
import "github.com/bhojpur/orm/pkg/now"

time.Now() // 2018-03-26 17:51:49.123456789 Mon

now.BeginningOfMinute()        // 2018-03-26 17:51:00 Mon
now.BeginningOfHour()          // 2018-03-26 17:00:00 Mon
now.BeginningOfDay()           // 2018-03-26 00:00:00 Mon
now.BeginningOfWeek()          // 2018-11-17 00:00:00 Sun
now.BeginningOfMonth()         // 2018-11-01 00:00:00 Fri
now.BeginningOfQuarter()       // 2018-10-01 00:00:00 Tue
now.BeginningOfYear()          // 2018-01-01 00:00:00 Tue

now.EndOfMinute()              // 2018-03-26 17:51:59.999999999 Mon
now.EndOfHour()                // 2018-03-26 17:59:59.999999999 Mon
now.EndOfDay()                 // 2018-03-26 23:59:59.999999999 Mon
now.EndOfWeek()                // 2018-11-23 23:59:59.999999999 Sat
now.EndOfMonth()               // 2018-11-30 23:59:59.999999999 Sat
now.EndOfQuarter()             // 2018-12-31 23:59:59.999999999 Tue
now.EndOfYear()                // 2018-12-31 23:59:59.999999999 Tue

now.WeekStartDay = time.Monday // Set Monday as first day, default is Sunday
now.EndOfWeek()                // 2018-11-24 23:59:59.999999999 Sun
```

Calculating time based on another time

```go
t := time.Date(2018, 02, 18, 17, 51, 49, 123456789, time.Now().Location())
now.With(t).EndOfMonth()   // 2018-02-28 23:59:59.999999999 Thu
```

Calculating time based on configuration

```go
location, err := time.LoadLocation("Asia/Kolkata")

myConfig := &now.Config{
	WeekStartDay: time.Monday,
	TimeLocation: location,
	TimeFormats: []string{"2006-01-02 15:04:05"},
}

t := time.Date(2018, 11, 18, 17, 51, 49, 123456789, time.Now().Location()) // // 2018-03-26 17:51:49.123456789 Mon
myConfig.With(t).BeginningOfWeek()         // 2018-03-26 00:00:00 Mon

myConfig.Parse("2002-10-12 22:14:01")     // 2002-10-12 22:14:01
myConfig.Parse("2002-10-12 22:14")        // returns error 'can't parse string as time: 2002-10-12 22:14'
```

### Monday/Sunday

Don't be bothered with the `WeekStartDay` setting, you can use `Monday`, `Sunday`

```go
now.Monday()              // 2018-03-26 00:00:00 Mon
now.Sunday()              // 2018-11-24 00:00:00 Sun (Next Sunday)
now.EndOfSunday()         // 2018-11-24 23:59:59.999999999 Sun (End of next Sunday)

t := time.Date(2018, 11, 24, 17, 51, 49, 123456789, time.Now().Location()) // 2018-11-24 17:51:49.123456789 Sun
now.With(t).Monday()       // 2018-03-26 00:00:00 Sun (Last Monday if today is Sunday)
now.With(t).Sunday()       // 2018-11-24 00:00:00 Sun (Beginning Of Today if today is Sunday)
now.With(t).EndOfSunday()  // 2018-11-24 23:59:59.999999999 Sun (End of Today if today is Sunday)
```

### Parse String to Time

```go
time.Now() // 2018-03-26 17:51:49.123456789 Mon

// Parse(string) (time.Time, error)
t, err := now.Parse("2017")                // 2017-01-01 00:00:00, nil
t, err := now.Parse("2017-10")             // 2017-10-01 00:00:00, nil
t, err := now.Parse("2017-10-13")          // 2017-10-13 00:00:00, nil
t, err := now.Parse("1999-12-12 12")       // 1999-12-12 12:00:00, nil
t, err := now.Parse("1999-12-12 12:20")    // 1999-12-12 12:20:00, nil
t, err := now.Parse("1999-12-12 12:20:21") // 1999-12-12 12:20:21, nil
t, err := now.Parse("10-13")               // 2018-10-13 00:00:00, nil
t, err := now.Parse("12:20")               // 2018-03-26 12:20:00, nil
t, err := now.Parse("12:20:13")            // 2018-03-26 12:20:13, nil
t, err := now.Parse("14")                  // 2018-03-26 14:00:00, nil
t, err := now.Parse("99:99")               // 2018-03-26 12:20:00, Can't parse string as time: 99:99

// MustParse must parse string to time or it will panic
now.MustParse("2018-03-26")             // 2018-03-26 00:00:00
now.MustParse("02-17")                  // 2018-02-17 00:00:00
now.MustParse("2-17")                   // 2018-02-17 00:00:00
now.MustParse("8")                      // 2018-03-26 08:00:00
now.MustParse("2002-10-12 22:14")       // 2002-10-12 22:14:00
now.MustParse("99:99")                  // panic: Can't parse string as time: 99:99
```

Extend `now` to support more formats is quite easy, just update `now.TimeFormats` with other time layouts, e.g:

```go
now.TimeFormats = append(now.TimeFormats, "26 Mar 2018 15:04")
```

Please send me pull requests if you want a format to be supported officially


## License

Released under the [MIT License](http://www.opensource.org/licenses/MIT).